package main

import (
	"context"
	html "html/template"
	"log/slog"
	"maps"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/server"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/templatefn"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/mod/sumdb/dirhash"
)

func main() {
	ctx := context.Background()
	logger := telemetry.NewLogger("opg-sirius-lpa-frontend")

	if err := run(ctx, logger); err != nil {
		logger.Error("fatal startup error", slog.Any("err", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	port := env.Get("PORT", "8080")
	webDir := env.Get("WEB_DIR", "web")
	siriusURL := env.Get("SIRIUS_URL", "http://localhost:9001")
	siriusPublicURL := env.Get("SIRIUS_PUBLIC_URL", "")
	prefix := env.Get("PREFIX", "")
	exportTraces := env.Get("TRACING_ENABLED", "0") == "1"

	staticHash, err := dirhash.HashDir(webDir+"/static", webDir, dirhash.DefaultHash)
	if err != nil {
		return err
	}

	layouts, err := parseLayoutTemplates(webDir+"/template/layout", templatefn.All(siriusPublicURL, prefix, staticHash))
	if err != nil {
		return err
	}
	tmpls, err := parseTemplates(webDir+"/template", layouts)
	if err != nil {
		return err
	}
	mlpaTmpls, err := parseTemplates(webDir+"/template/mlpa", layouts)
	if err != nil {
		return err
	}
	poaTmpls, err := parseTemplates(webDir+"/template/poa", layouts)
	if err != nil {
		return err
	}

	templates := combineAllLayouts(tmpls, mlpaTmpls, poaTmpls)

	shutdown, err := telemetry.StartTracerProvider(ctx, logger, exportTraces)
	defer shutdown()
	if err != nil {
		return err
	}

	httpClient := http.DefaultClient
	httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)

	client := sirius.NewClient(httpClient, siriusURL)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           server.New(logger, client, templates, prefix, siriusPublicURL, webDir),
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("listen and server error", slog.Any("err", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Info("signal received: ", "sig", sig)

	tc, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return server.Shutdown(tc)
}

func combineAllLayouts(tmpls template.Templates, mlpaTmpls template.Templates, poaTmpls template.Templates) template.Templates {
	templates := make(template.Templates)

	maps.Copy(templates, tmpls)
	maps.Copy(templates, mlpaTmpls)
	maps.Copy(templates, poaTmpls)
	return templates
}

func parseLayoutTemplates(layoutDir string, funcs html.FuncMap) (*html.Template, error) {
	return html.New("").Funcs(funcs).ParseGlob(filepath.Join(layoutDir, "*.*"))
}

func parseTemplates(templateDir string, layouts *html.Template) (template.Templates, error) {
	files, err := filepath.Glob(filepath.Join(templateDir, "*.*"))
	if err != nil {
		return nil, err
	}

	tmpls := map[string]*html.Template{}
	for _, file := range files {
		clone, err := layouts.Clone()
		if err != nil {
			return nil, err
		}

		tmpl, err := clone.ParseFiles(file)
		if err != nil {
			return nil, err
		}

		tmpls[filepath.Base(file)] = tmpl
	}

	return tmpls, nil
}
