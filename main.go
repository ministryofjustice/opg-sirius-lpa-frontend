package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/server"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/telemetry"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/templatefn"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/mod/sumdb/dirhash"
)

func fatal(err error, logger *slog.Logger) {
	logger.Error(fmt.Sprintf("a fatal error occurred: %s", err.Error()))
	os.Exit(1)
}

func main() {
	logger := telemetry.NewLogger()

	port := env.Get("PORT", "8080")
	webDir := env.Get("WEB_DIR", "web")
	siriusURL := env.Get("SIRIUS_URL", "http://localhost:9001")
	siriusPublicURL := env.Get("SIRIUS_PUBLIC_URL", "")
	prefix := env.Get("PREFIX", "")
	exportTraces := env.Get("TRACING_ENABLED", "0") == "1"

	staticHash, err := dirhash.HashDir(webDir+"/static", webDir, dirhash.DefaultHash)
	if err != nil {
		fatal(err, logger)
	}

	tmpls, err := template.Parse(webDir+"/template", templatefn.All(siriusPublicURL, prefix, staticHash))
	if err != nil {
		fatal(err, logger)
	}

	shutdown := telemetry.InitTracerProvider(context.Background(), logger, exportTraces)
	defer shutdown()

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
			fatal(err, logger)
		}
	}()

	logger.Info("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Info("signal received: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		logger.Error(err.Error())
	}
}
