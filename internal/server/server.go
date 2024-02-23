package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

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
	AddFeeDecisionClient
	AddPaymentClient
	AllocateCasesClient
	AssignTaskClient
	ApplyFeeReductionClient
	ChangeStatusClient
	CreateDonorClient
	CreateDocumentClient
	CreateDraftClient
	CreateInvestigationClient
	EditInvestigationClient
	InvestigationHoldClient
	DeletePaymentClient
	DeleteRelationshipClient
	EditComplaintClient
	EditDatesClient
	EditDocumentClient
	EditDonorClient
	EditPaymentClient
	EventClient
	GetPaymentsClient
	GetDocumentsClient
	LinkPersonClient
	LpaClient
	ManageFeesClient
	MiReportingClient
	RelationshipClient
	SearchDonorsClient
	SearchClient
	PostcodeLookupClient
	SearchUsersClient
	TaskClient
	UnlinkPersonClient
	WarningClient
}

var decoder *form.Decoder

func New(logger *slog.Logger, client Client, templates template.Templates, prefix, siriusPublicURL, webDir string) http.Handler {
	wrap := errorHandler(templates.Get("error.gohtml"), prefix, siriusPublicURL)

	mux := chi.NewRouter()

	mux.Use(telemetry.AttachMiddleware(logger))

	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {})

	mux.Handle("/create-warning", wrap(Warning(client, templates.Get("warning.gohtml"))))
	mux.Handle("/create-event", wrap(Event(client, templates.Get("event.gohtml"))))
	mux.Handle("/create-task", wrap(Task(client, templates.Get("task.gohtml"))))
	mux.Handle("/create-relationship", wrap(Relationship(client, templates.Get("relationship.gohtml"))))
	mux.Handle("/create-donor", wrap(CreateDonor(client, templates.Get("donor.gohtml"))))
	mux.Handle("/create-investigation", wrap(CreateInvestigation(client, templates.Get("create_investigation.gohtml"))))
	mux.Handle("/create-document", wrap(CreateDocument(client, templates.Get("create_document.gohtml"))))
	mux.Handle("/edit-document", wrap(EditDocument(client, templates.Get("edit_document.gohtml"))))
	mux.Handle("/investigation-hold", wrap(InvestigationHold(client, templates.Get("investigation_hold.gohtml"))))
	mux.Handle("/edit-investigation", wrap(EditInvestigation(client, templates.Get("edit_investigation.gohtml"))))
	mux.Handle("/edit-donor", wrap(EditDonor(client, templates.Get("donor.gohtml"))))
	mux.Handle("/delete-relationship", wrap(DeleteRelationship(client, templates.Get("delete_relationship.gohtml"))))
	mux.Handle("/edit-dates", wrap(EditDates(client, templates.Get("edit_dates.gohtml"))))
	mux.Handle("/link-person", wrap(LinkPerson(client, templates.Get("link_person.gohtml"))))
	mux.Handle("/add-complaint", wrap(AddComplaint(client, templates.Get("add_complaint.gohtml"))))
	mux.Handle("/edit-complaint", wrap(EditComplaint(client, templates.Get("edit_complaint.gohtml"))))
	mux.Handle("/unlink-person", wrap(UnlinkPerson(client, templates.Get("unlink_person.gohtml"))))
	mux.Handle("/change-status", wrap(ChangeStatus(client, templates.Get("change_status.gohtml"))))
	mux.Handle("/allocate-cases", wrap(AllocateCases(client, templates.Get("allocate_cases.gohtml"))))
	mux.Handle("/assign-task", wrap(AssignTask(client, templates.Get("assign_task.gohtml"))))
	mux.Handle("/mi-reporting", wrap(MiReporting(client, templates.Get("mi_reporting.gohtml"))))
	mux.Handle("/add-payment", wrap(AddPayment(client, templates.Get("add_payment.gohtml"))))
	mux.Handle("/delete-payment", wrap(DeletePayment(client, templates.Get("delete_payment.gohtml"))))
	mux.Handle("/manage-fees", wrap(AddFeeDecision(client, templates.Get("manage_fees.gohtml"))))
	mux.Handle("/add-fee-decision", wrap(AddFeeDecision(client, templates.Get("add_fee_decision.gohtml"))))
	mux.Handle("/apply-fee-reduction", wrap(ApplyFeeReduction(client, templates.Get("apply_fee_reduction.gohtml"))))
	mux.Handle("/delete-fee-reduction", wrap(DeletePayment(client, templates.Get("delete_fee_reduction.gohtml"))))
	mux.Handle("/edit-payment", wrap(EditPayment(client, templates.Get("edit_payment.gohtml"))))
	mux.Handle("/edit-fee-reduction", wrap(EditFeeReduction(client, templates.Get("edit_fee_reduction.gohtml"))))
	mux.Handle("/payments/{id}", wrap(GetPayments(client, templates.Get("payments.gohtml"))))
	mux.Handle("/lpa/{uid}/lpa-details", wrap(GetLpaDetails(client, templates.Get("mlpa-details.gohtml"))))
	mux.Handle("/lpa/{uid}/payments", wrap(GetPayments(client, templates.Get("mlpa-payments.gohtml"))))
	mux.Handle("/lpa/{uid}/documents", wrap(GetDocuments(client, templates.Get("mlpa-documents.gohtml"))))
	mux.Handle("/lpa/{uid}/documents/new", wrap(CreateDocumentDigitalLpa(client, templates.Get("mlpa-create_document.gohtml"))))
	mux.Handle("/search-users", wrap(SearchUsers(client)))
	mux.Handle("/search-persons", wrap(SearchDonors(client)))
	mux.Handle("/search-postcode", wrap(SearchPostcode(client)))
	mux.Handle("/search", wrap(Search(client, templates.Get("search.gohtml"))))
	mux.Handle("/digital-lpa/create", wrap(CreateDraft(client, templates.Get("create-draft.gohtml"))))
	mux.Handle("/lpa/{uid}", wrap(Lpa(client, templates.Get("lpa.gohtml"))))

	static := http.FileServer(http.Dir("web/static"))
	mux.Handle("/assets/*", static)
	mux.Handle("/javascript/*", static)
	mux.Handle("/stylesheets/*", static)

	muxWithHeaders := securityheaders.Use(setCSPHeader(mux))

	return otelhttp.NewHandler(http.StripPrefix(prefix, muxWithHeaders), "lpa-frontend")
}

type Handler func(w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	SiriusURL     string
	Path          string
	Code          int
	Error         string
	CorrelationId string
}

type unauthorizedError interface {
	IsUnauthorized() bool
}

type RedirectError string

func (e RedirectError) Error() string {
	return "redirect to " + string(e)
}

func (e RedirectError) To() string {
	return string(e)
}

type ProblemError struct {
	Title            string             `json:"title"`
	Detail           string             `json:"detail"`
	ValidationErrors sirius.FieldErrors `json:"validationErrors"`
}

func errorHandler(tmplError template.Template, prefix, siriusURL string) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := next(w, r); err != nil {
				if errors.Is(err, context.Canceled) {
					w.WriteHeader(499)
					return
				}

				if v, ok := err.(unauthorizedError); ok && v.IsUnauthorized() {
					http.Redirect(w, r, fmt.Sprintf("%s/auth?redirect=%s", siriusURL, url.QueryEscape(prefix+r.URL.Path)), http.StatusFound)
					return
				}

				if redirect, ok := err.(RedirectError); ok {
					http.Redirect(w, r, prefix+redirect.To(), http.StatusFound)
					return
				}

				code := http.StatusInternalServerError
				correlationId := ""
				logger := r.Context().Value("logger").(*slog.Logger)

				if statusError, ok := err.(sirius.StatusError); ok {
					code = statusError.Code
					correlationId = statusError.CorrelationId
				}

				if r.Header.Get("Accept") == "application/json" {
					rfcErr := ProblemError{
						Title: err.Error(),
					}

					if ve, ok := err.(sirius.ValidationError); ok {
						code = http.StatusBadRequest
						rfcErr.Detail = ve.Detail
						rfcErr.ValidationErrors = ve.Field
					}

					if code == http.StatusInternalServerError {
						logger.Error(err.Error())
					}

					w.Header().Add("Content-Type", "application/problem+json")
					w.WriteHeader(code)

					err = json.NewEncoder(w).Encode(rfcErr)

					if err != nil {
						logger.Error(err.Error())
						http.Error(w, "Could not generate error JSON", http.StatusInternalServerError)
					}

					return
				}

				if code == http.StatusInternalServerError {
					logger.Error(err.Error())
				}

				w.WriteHeader(code)
				err = tmplError(w, errorVars{
					SiriusURL:     siriusURL,
					Path:          "",
					Code:          code,
					Error:         err.Error(),
					CorrelationId: correlationId,
				})

				if err != nil {
					logger.Error(err.Error())
					http.Error(w, "Could not generate error template", http.StatusInternalServerError)
				}
			}
		})
	}
}

func postFormString(r *http.Request, name string) string {
	return strings.TrimSpace(r.PostFormValue(name))
}

func postFormCheckboxChecked(r *http.Request, name string, value string) bool {
	for _, val := range r.PostForm[name] {
		if val == value {
			return true
		}
	}

	return false
}

func postFormInt(r *http.Request, name string) (int, error) {
	return strconv.Atoi(postFormString(r, name))
}

func postFormDateString(r *http.Request, name string) sirius.DateString {
	return sirius.DateString(postFormString(r, name))
}

func translateRefData(types []sirius.RefDataItem, tmplHandle string) string {
	for _, refDataType := range types {
		if refDataType.Handle == tmplHandle {
			return refDataType.Label
		}
	}
	return tmplHandle
}

func setCSPHeader(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:")

		h.ServeHTTP(w, r)
	}
}
