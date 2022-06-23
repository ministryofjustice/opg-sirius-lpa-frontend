package server

import (
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/template"
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

type Client interface {
	AddComplaintClient
	ChangeStatusClient
	DeleteRelationshipClient
	EditComplaintClient
	EditDatesClient
	EventClient
	LinkPersonClient
	RelationshipClient
	SearchDonorsClient
	SearchUsersClient
	TaskClient
	WarningClient
}

func New(logger Logger, client Client, templates template.Templates, prefix, siriusPublicURL, webDir string) http.Handler {
	wrap := errorHandler(logger, templates.Get("error.gohtml"), prefix, siriusPublicURL)

	mux := http.NewServeMux()

	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/create-warning", wrap(Warning(client, templates.Get("warning.gohtml"))))
	mux.Handle("/create-event", wrap(Event(client, templates.Get("event.gohtml"))))
	mux.Handle("/create-task", wrap(Task(client, templates.Get("task.gohtml"))))
	mux.Handle("/create-relationship", wrap(Relationship(client, templates.Get("relationship.gohtml"))))
	mux.Handle("/delete-relationship", wrap(DeleteRelationship(client, templates.Get("delete_relationship.gohtml"))))
	mux.Handle("/edit-dates", wrap(EditDates(client, templates.Get("edit_dates.gohtml"))))
	mux.Handle("/link-person", wrap(LinkPerson(client, templates.Get("link_person.gohtml"))))
	mux.Handle("/add-complaint", wrap(AddComplaint(client, templates.Get("add_complaint.gohtml"))))
	mux.Handle("/edit-complaint", wrap(EditComplaint(client, templates.Get("edit_complaint.gohtml"))))
	mux.Handle("/change-status", wrap(ChangeStatus(client, templates.Get("change_status.gohtml"))))

	mux.Handle("/search-users", wrap(SearchUsers(client)))
	mux.Handle("/search-persons", wrap(SearchDonors(client)))

	static := http.FileServer(http.Dir("web/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return http.StripPrefix(prefix, securityheaders.Use(mux))
}

type Handler func(w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	SiriusURL string
	Path      string
	Code      int
	Error     string
}

type unauthorizedError interface {
	IsUnauthorized() bool
}

func errorHandler(logger Logger, tmplError template.Template, prefix, siriusURL string) func(next Handler) http.Handler {
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
				err = tmplError(w, errorVars{
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
