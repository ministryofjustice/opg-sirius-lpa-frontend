package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/logging"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/server"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

func main() {
	logger := logging.New(os.Stdout, "opg-sirius-lpa-frontend")

	port := getEnv("PORT", "8080")
	webDir := getEnv("WEB_DIR", "web")
	siriusURL := getEnv("SIRIUS_URL", "http://localhost:8080")
	siriusPublicURL := getEnv("SIRIUS_PUBLIC_URL", "")
	prefix := getEnv("PREFIX", "")

	layouts, err := template.New("").Funcs(map[string]interface{}{
		"sirius": func(s string) string {
			return s
		},
		"prefix": func(s string) string {
			return s
		},
	}).ParseGlob("./web/template/layouts/*")

	if err != nil {
		logger.Fatal(err)
	}

	files, _ := filepath.Glob("./web/template/*.gohtml")
	tmpls := map[string]*template.Template{}

	for _, file := range files {
		tmpls[filepath.Base(file)] = template.Must(template.Must(layouts.Clone()).ParseFiles(file))
	}

	client := sirius.NewClient(http.DefaultClient, siriusURL)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: server.New(logger, client, tmpls, prefix, siriusPublicURL, webDir),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
