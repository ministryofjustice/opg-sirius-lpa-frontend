package server

import (
	"html/template"
	"io"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type Logger interface {
	Request(*http.Request, error)
}

type Server struct {
	Templates map[string]*template.Template
	Client    *sirius.Client
}

func getContext(r *http.Request) sirius.Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	return sirius.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

type Client interface {
	WarningClient
}

func New(logger Logger, client Client, templates map[string]*template.Template, prefix, siriusPublicURL, webDir string) http.Handler {
	wrap := errorHandler(logger, templates["error.gohtml"], prefix, siriusPublicURL)

	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {})
	mux.Handle("/create-warning", wrap(Warning(client, templates["warning.gohtml"])))
	mux.Handle("/assets/", http.FileServer(http.Dir("web/static")))
	mux.Handle("/javascript/", http.FileServer(http.Dir("web/static")))
	mux.Handle("/stylesheets/", http.FileServer(http.Dir("web/static")))
	return mux
}

func securityHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Security-Policy", "default-src 'self'")
		w.Header().Add("Referrer-Policy", "same-origin")
		w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-Frame-Options", "SAMEORIGIN")
		w.Header().Add("X-XSS-Protection", "1; mode=block")

		h.ServeHTTP(w, r)
	}
}

type Handler func(w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	SiriusURL string
	Path      string

	Code  int
	Error string
}

type unauthorizedError interface {
	IsUnauthorized() bool
}

func errorHandler(logger Logger, tmplError Template, prefix, siriusURL string) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := next(w, r); err != nil {
				if v, ok := err.(unauthorizedError); ok && v.IsUnauthorized() {
					http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
					return
				}

				logger.Request(r, err)

				code := http.StatusInternalServerError

				w.WriteHeader(code)
				err = tmplError.ExecuteTemplate(w, "page", errorVars{
					SiriusURL: siriusURL,
					Path:      "",
					Code:      code,
					Error:     err.Error(),
				})

				if err != nil {
					logger.Request(r, err)
					http.Error(w, "Could not generate error template", http.StatusInternalServerError)
				}
			}
		})
	}
}
