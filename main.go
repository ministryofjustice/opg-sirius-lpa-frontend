package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/server"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/templatefn"
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

	tmpls, err := template.Parse(webDir+"/template", templatefn.All(siriusPublicURL, prefix, staticHash))
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
