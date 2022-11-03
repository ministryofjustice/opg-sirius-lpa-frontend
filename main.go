package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/server"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/mod/sumdb/dirhash"
	"google.golang.org/grpc"
)

func initTracerProvider(ctx context.Context, logger *logging.Logger) func() {
	resource, err := ecs.NewResourceDetector().Detect(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("0.0.0.0:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		logger.Fatal(err)
	}

	idg := xray.NewIDGenerator()
	tp := trace.NewTracerProvider(
		trace.WithResource(resource),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(traceExporter),
		trace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Fatal(err)
		}
	}
}

func main() {
	logger := logging.New(os.Stdout, "opg-sirius-lpa-frontend")

	port := env.Get("PORT", "8080")
	webDir := env.Get("WEB_DIR", "web")
	siriusURL := env.Get("SIRIUS_URL", "http://localhost:9001")
	siriusPublicURL := env.Get("SIRIUS_PUBLIC_URL", "")
	prefix := env.Get("PREFIX", "")

	staticHash, err := dirhash.HashDir(webDir+"/static", webDir, dirhash.DefaultHash)
	if err != nil {
		logger.Fatal(err)
	}

	tmpls, err := template.Parse(webDir+"/template", map[string]interface{}{
		"sirius": func(s string) string {
			return siriusPublicURL + s
		},
		"prefix": func(s string) string {
			return prefix + s
		},
		"prefixAsset": func(s string) string {
			if len(staticHash) >= 11 {
				return prefix + s + "?" + url.QueryEscape(staticHash[3:11])
			} else {
				return prefix + s
			}
		},
		"today": func() string {
			return time.Now().Format("2006-01-02")
		},
		"field": func(name, label string, value interface{}, error map[string]string, attrs ...interface{}) map[string]interface{} {
			field := map[string]interface{}{
				"name":  name,
				"label": label,
				"value": value,
				"error": error,
			}

			if len(attrs)%2 != 0 {
				panic("must have even number of attrs")
			}

			for i := 0; i < len(attrs); i += 2 {
				field[attrs[i].(string)] = attrs[i+1]
			}

			return field
		},
		"fee": func(amount int) string {
			float := float64(amount)
			return fmt.Sprintf("%.2f", float/100)
		},
		"formatDate": func(s sirius.DateString) (string, error) {
			if s != "" {
				return s.ToSirius()
			}
			return "", nil
		},
		"translateRefData": func(types []sirius.RefDataItem, tmplHandle string) string {
			for _, refDataType := range types {
				if refDataType.Handle == tmplHandle {
					return refDataType.Label
				}
			}
			return tmplHandle
		},
		"ToLower": strings.ToLower,
	})
	if err != nil {
		logger.Fatal(err)
	}

	if env.Get("TRACING_ENABLED", "0") == "1" {
		shutdown := initTracerProvider(context.Background(), logger)
		defer shutdown()
	}

	httpClient := http.DefaultClient
	httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)

	client := sirius.NewClient(httpClient, siriusURL)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           server.New(logger, client, tmpls, prefix, siriusPublicURL, webDir),
		ReadHeaderTimeout: 20 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Print("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Print("signal received: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		logger.Print(err)
	}
}
