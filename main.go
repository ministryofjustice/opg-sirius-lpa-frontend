package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/logging"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/server"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

func main() {
	layouts, err := template.New("").Funcs(map[string]interface{}{
		"sirius": func(s string) string {
			return s
		},
		"prefix": func(s string) string {
			return s
		},
	}).ParseGlob("./web/template/layouts/*")
	if err != nil {
		fmt.Print(err)
	}

	files, _ := filepath.Glob("./web/template/*.gohtml")
	tmpls := map[string]*template.Template{}

	for _, file := range files {
		tmpls[filepath.Base(file)] = template.Must(template.Must(layouts.Clone()).ParseFiles(file))
	}

	logger := logging.New(os.Stdout, "opg-sirius-lpa-frontend")

	response := server.New(logger, sirius.NewClient(http.DefaultClient, os.Getenv("SIRIUS_URL")), tmpls, "", os.Getenv("SIRIUS_PUBLIC_URL"), "web")
	http.ListenAndServe(":8888", response)
}
