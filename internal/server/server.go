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

	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Server struct {
	Templates map[string]*template.Template
	Client    *sirius.Client
}

func getContext(r *http.Request) sirius.Context {
	token := ""

	if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
		token, _ = url.QueryUnescape(cookie.Value)
	}

	return sirius.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
		IsPartial: r.Header.Get("HX-Request") == "true",
	}
}

type Client interface {
	ActionPanelClient
	CasePageClient
	CreateComplaintClient
	CreateFeeDecisionClient
	CreateObjectionClient
	CreatePaymentClient
	AllocateCasesClient
	ApplyFeeReductionClient
	AssignTaskClient
	AttorneyDecisionsClient
	UpdateAttorneyDetailsClient
	UpdateCaseStatusClient
	UpdateCertificateProviderDetailsClient
	UpdateDonorDetailsClient
	UpdateDraftClient
	UpdateStatusClient
	UpdateTrustCorporationDetailsClient
	ClearTaskClient
	CompareDocsClient
	CreateAdditionalDraftClient
	CreateAttorneyClient
	CreateCertificateProviderClient
	CreateCorrespondentClient
	CreateDocumentClient
	CreateDocumentDigitalLpaClient
	CreateDonorClient
	CreateDraftClient
	CreateEpaClient
	CreateInvestigationClient
	CreateLpaClient
	CreateNotifiedPersonClient
	CreateReplacementAttorneyClient
	CreateTrustCorporationClient
	DeleteDocumentClient
	DeleteNoteClient
	DeletePaymentClient
	DeleteRelationshipClient
	DocumentListClient
	DonorDetailsClient
	UpdateCertificateProviderClient
	UpdateComplaintClient
	UpdateCorrespondentClient
	UpdateDatesClient
	UpdateDocumentClient
	UpdateDonorClient
	UpdateFeeReductionClient
	UpdateInvestigationClient
	UpdatePaymentClient
	EventClient
	GetApplicationProgressClient
	GetDocumentsClient
	GetHistoryClient
	GetLpaDetailsClient
	GetLpaHistoryClient
	GetPaymentsClient
	InvestigationHoldClient
	LinkPersonClient
	ManageAttorneysClient
	ManageFeesClient
	ManageRestrictionsClient
	MiReportingClient
	ObjectionOutcomeClient
	PostcodeLookupClient
	RelationshipClient
	RemoveAnAttorneyClient
	ResolveObjectionClient
	SearchClient
	SearchDonorsClient
	SearchUsersClient
	SelectOrCreateCorrespondentClient
	SiriusHeaderCaseInfoClient
	SiriusHeaderCalendarClient
	SiriusHeaderPeopleInfoClient
	TaskClient
	UnlinkPersonClient
	UpdateAttorneyClient
	UpdateDecisionsClient
	UpdateObjectionClient
	UpdateNotifiedPersonClient
	ViewDocumentClient
	WarningClient
}

var decoder = form.NewDecoder()

func New(logger *slog.Logger, client Client, templates template.Templates, prefix, siriusPublicURL, webDir string) http.Handler {
	wrap := errorHandler(templates.Get("error.gohtml"), prefix, siriusPublicURL)
	mux := http.NewServeMux()

	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {})

	//search
	mux.Handle("/search-users", wrap(SearchUsers(client)))
	mux.Handle("/search-persons", wrap(SearchDonors(client)))
	mux.Handle("/search-postcode", wrap(SearchPostcode(client)))
	mux.Handle("/search", wrap(Search(client, templates.Get("search.gohtml"))))

	//shared templates (Used in both modernise and LPA)
	mux.Handle("/create-payment", wrap(CreatePayment(client, templates.Get("create-payment.gohtml"))))
	mux.Handle("/apply-fee-reduction", wrap(ApplyFeeReduction(client, templates.Get("apply-fee-reduction.gohtml"))))
	mux.Handle("/assign-task", wrap(AssignTask(client, templates.Get("assign-task.gohtml"))))
	mux.Handle("/create-event", wrap(Event(client, templates.Get("event.gohtml"))))
	mux.Handle("/create-task", wrap(Task(client, templates.Get("create-task.gohtml"))))
	mux.Handle("/create-warning", wrap(Warning(client, templates.Get("warning.gohtml"))))
	mux.Handle("/update-document", wrap(UpdateDocument(client, templates.Get("update-document.gohtml"))))
	mux.Handle("/edit-document", wrap(UpdateDocument(client, templates.Get("update-document.gohtml"))))

	//modernise
	mux.Handle("/create-fee-decision", wrap(CreateFeeDecision(client, templates.Get("create-fee-decision.gohtml"))))
	mux.Handle("/create-objection", wrap(CreateObjection(client, templates.Get("objection.gohtml"))))
	mux.Handle("/update-case-status", wrap(UpdateCaseStatus(client, templates.Get("update-case-status.gohtml"))))
	mux.Handle("/update-donor-details", wrap(UpdateDonorDetails(client, templates.Get("update-donor-details.gohtml"))))
	mux.Handle("/clear-task", wrap(ClearTask(client, templates.Get("clear_task.gohtml"))))
	mux.Handle("/create-additional-draft-lpa", wrap(CreateAdditionalDraft(client, templates.Get("create_additional_draft.gohtml"))))
	mux.Handle("/digital-lpa/create", wrap(CreateDraft(client, templates.Get("create_draft.gohtml"))))
	mux.Handle("/lpa/{uid}", wrap(GetApplicationProgressDetails(client, templates.Get("mlpa-application-progress.gohtml"))))
	mux.Handle("/lpa/{uid}/attorney/{attorneyUID}/update-details", wrap(UpdateAttorneyDetails(client, templates.Get("update-attorney-details.gohtml"))))
	mux.Handle("/lpa/{uid}/certificate-provider/update-details", wrap(UpdateCertificateProviderDetails(client, templates.Get("update-certificate-provider-details.gohtml"))))
	mux.Handle("/lpa/{uid}/update-draft", wrap(UpdateDraft(client, templates.Get("update-draft.gohtml"))))
	mux.Handle("/lpa/{uid}/documents", wrap(GetDocuments(client, templates.Get("mlpa-documents.gohtml"))))
	mux.Handle("/lpa/{uid}/documents/new", wrap(CreateDocumentDigitalLpa(client, templates.Get("mlpa-create_document.gohtml"))))
	mux.Handle("/lpa/{uid}/history", wrap(GetHistory(client, templates.Get("mlpa-history.gohtml"))))
	mux.Handle("/lpa/{uid}/lpa-details", wrap(GetLpaDetails(client, templates.Get("mlpa-details.gohtml"))))
	mux.Handle("/lpa/{uid}/objection/{id}", wrap(UpdateObjection(client, templates.Get("objection.gohtml"), templates.Get("confirm-objection.gohtml"))))
	mux.Handle("/lpa/{uid}/objection/{id}/outcome", wrap(ObjectionOutcome(client, templates.Get("objection-outcome.gohtml"))))
	mux.Handle("/lpa/{uid}/objection/{id}/resolve", wrap(ResolveObjection(client, templates.Get("resolve-objection.gohtml"))))
	mux.Handle("/lpa/{uid}/manage-attorney-decisions", wrap(AttorneyDecisions(client, templates.Get("mlpa-attorney-decisions.gohtml"), templates.Get("mlpa-confirm-attorney-decisions.gohtml"))))
	mux.Handle("/lpa/{uid}/manage-attorneys", wrap(ManageAttorneys(client, templates.Get("mlpa-manage-attorneys.gohtml"))))
	mux.Handle("/lpa/{uid}/manage-restrictions", wrap(ManageRestrictions(client, templates.Get("manage-restrictions.gohtml"), templates.Get("confirm-restrictions.gohtml"))))
	mux.Handle("/lpa/{uid}/payments", wrap(GetPayments(client, templates.Get("mlpa-payments.gohtml"))))
	mux.Handle("/lpa/{uid}/remove-an-attorney", wrap(RemoveAnAttorney(client, templates.Get("mlpa-remove-attorney.gohtml"), templates.Get("mlpa-confirm-attorney-removal.gohtml"), templates.Get("mlpa-attorney-decisions.gohtml"))))
	mux.Handle("/lpa/{uid}/trust-corporation/{trustCorporationUID}/update-details", wrap(UpdateTrustCorporationDetails(client, templates.Get("update-trust-corporation-details.gohtml"))))
	mux.Handle("/lpa/{uid}/update-decisions", wrap(UpdateDecisions(client, templates.Get("mlpa-update-decisions.gohtml"))))
	mux.Handle("/manage-fees", wrap(CreateFeeDecision(client, templates.Get("manage_fees.gohtml"))))

	//LPA
	mux.Handle("/action-panel", wrap(ActionPanel(client, templates.Get("action-panel-wrapper.gohtml"))))
	mux.Handle("/create-complaint", wrap(CreateComplaint(client, templates.Get("create-complaint.gohtml"))))
	mux.Handle("/add-complaint", wrap(CreateComplaint(client, templates.Get("create-complaint.gohtml"))))
	mux.Handle("/allocate-cases", wrap(AllocateCases(client, templates.Get("allocate-cases.gohtml"))))
	mux.Handle("/update-status", wrap(UpdateStatus(client, templates.Get("update-status.gohtml"))))
	mux.Handle("/change-status", wrap(UpdateStatus(client, templates.Get("update-status.gohtml"))))
	mux.Handle("/create-attorney", wrap(CreateAttorney(client, templates.Get("create-attorney.gohtml"))))
	mux.Handle("/update-attorney", wrap(UpdateAttorney(client, templates.Get("create-attorney.gohtml"))))
	mux.Handle("/create-certificate-provider", wrap(CreateCertificateProvider(client, templates.Get("certificate-provider.gohtml"))))
	mux.Handle("/create-correspondent", wrap(CreateCorrespondent(client, templates.Get("create-correspondent.gohtml"))))
	mux.Handle("/create-donor", wrap(CreateDonor(client, templates.Get("donor.gohtml"))))
	mux.Handle("/create-document", wrap(CreateDocument(client, templates.Get("create-document.gohtml"))))
	mux.Handle("/create-epa", wrap(CreateEpa(client, templates.Get("create-epa.gohtml"))))
	mux.Handle("/create-decline-fee-decision", wrap(CreateFeeDecision(client, templates.Get("decline-fee-reduction.gohtml"))))
	mux.Handle("/create-investigation", wrap(CreateInvestigation(client, templates.Get("create-investigation.gohtml"))))
	mux.Handle("/create-lpa", wrap(CreateLpa(client, templates.Get("create-lpa.gohtml"))))
	mux.Handle("/create-relationship", wrap(Relationship(client, templates.Get("create-relationship.gohtml"))))
	mux.Handle("/create-notified-person", wrap(CreateNotifiedPerson(client, templates.Get("create-notified-person.gohtml"))))
	mux.Handle("/create-replacement-attorney", wrap(CreateReplacementAttorney(client, templates.Get("create-replacement-attorney.gohtml"))))
	mux.Handle("/create-trust-corporation", wrap(CreateTrustCorporation(client, templates.Get("create-trust-corporation.gohtml"))))
	mux.Handle("/compare/{id}/{caseUid}", wrap(CompareDocs(client, templates.Get("compare-docs.gohtml"))))
	mux.Handle("/delete-fee-reduction", wrap(DeletePayment(client, templates.Get("delete-fee-reduction.gohtml"))))
	mux.Handle("/delete-note", wrap(DeleteNote(client, templates.Get("delete-note.gohtml"))))
	mux.Handle("/delete-payment", wrap(DeletePayment(client, templates.Get("delete-payment.gohtml"))))
	mux.Handle("/delete-relationship", wrap(DeleteRelationship(client, templates.Get("delete-relationship.gohtml"))))
	mux.Handle("/donor/{id}", wrap(CasePage(client, templates.Get("case-page.gohtml"))))
	mux.Handle("/donor/{donorId}/details", wrap(DonorDetails(client, templates.Get("donor-details.gohtml"))))
	mux.Handle("/donor/{id}/documents", wrap(DocumentList(client, templates.Get("documents.gohtml"))))
	mux.Handle("/donor/{donorId}/history", wrap(GetLpaHistory(client, templates.Get("lpa-history.gohtml"))))
	mux.Handle("/view-document/{uuid}/{id}", wrap(ViewDocument(client, templates.Get("view-document.gohtml"))))
	mux.Handle("/delete-document/{uuid}", wrap(DeleteDocument(client, templates.Get("delete-document.gohtml"))))
	mux.Handle("/update-certificate-provider", wrap(UpdateCertificateProvider(client, templates.Get("certificate-provider.gohtml"))))
	mux.Handle("/update-complaint", wrap(UpdateComplaint(client, templates.Get("edit_complaint.gohtml"))))
	mux.Handle("/edit-complaint", wrap(UpdateComplaint(client, templates.Get("edit_complaint.gohtml"))))
	mux.Handle("/edit-correspondent", wrap(UpdateCorrespondent(client, templates.Get("create-correspondent.gohtml"))))
	mux.Handle("/update-dates", wrap(UpdateDates(client, templates.Get("update-dates.gohtml"))))
	mux.Handle("/edit-dates", wrap(UpdateDates(client, templates.Get("update-dates.gohtml"))))
	mux.Handle("/update-donor", wrap(UpdateDonor(client, templates.Get("donor.gohtml"))))
	mux.Handle("/edit-donor", wrap(UpdateDonor(client, templates.Get("donor.gohtml"))))
	mux.Handle("/update-fee-reduction", wrap(UpdateFeeReduction(client, templates.Get("update-fee-reduction.gohtml"))))
	mux.Handle("/update-investigation", wrap(UpdateInvestigation(client, templates.Get("edit_investigation.gohtml"))))
	mux.Handle("/edit-investigation", wrap(UpdateInvestigation(client, templates.Get("edit_investigation.gohtml"))))
	mux.Handle("/update-payment", wrap(UpdatePayment(client, templates.Get("update-payment.gohtml"))))
	mux.Handle("/investigation-hold", wrap(InvestigationHold(client, templates.Get("investigation_hold.gohtml"))))
	mux.Handle("/link-person", wrap(LinkPerson(client, templates.Get("link-person.gohtml"))))
	mux.Handle("/mi-reporting", wrap(MiReporting(client, templates.Get("mi-reporting.gohtml"))))
	mux.Handle("/notify-error/{donorId}/{createdDocumentUUID}/{caseType}/{caseID}/{caseUUID}", wrap(NotifyPageLimit(templates.Get("notify-page-limit.gohtml"))))
	mux.Handle("/payments/{id}", wrap(GetPayments(client, templates.Get("payments.gohtml"))))
	mux.Handle("/select-or-create-correspondent", wrap(SelectOrCreateCorrespondent(client, templates.Get("select-or-create-correspondent.gohtml"))))
	mux.Handle("/sirius-header-calendars", wrap(SiriusHeaderCalendars(client, templates.Get("sirius-header-partial-calendars.gohtml"))))
	mux.Handle("/sirius-header-case-info", wrap(SiriusHeaderCaseInfo(client, templates.Get("sirius-header-partial-case-info.gohtml"))))
	mux.Handle("/sirius-header-people-info", wrap(SiriusHeaderPeopleInfo(client, templates.Get("sirius-header-partial-people-info.gohtml"))))
	mux.Handle("/update-notified-person", wrap(UpdateNotifiedPerson(client, templates.Get("create-notified-person.gohtml"))))
	mux.Handle("/unlink-person", wrap(UnlinkPerson(client, templates.Get("unlink-person.gohtml"))))
	mux.Handle("/view-document/{uuid}", wrap(ViewDocument(client, templates.Get("view-document.gohtml"))))
	mux.Handle("/working-days", wrap(WorkingDays(client, templates.Get("working-days-partial.gohtml"))))
	mux.Handle("/calendar-month", wrap(CalendarMonthPartial(client, templates.Get("calendar-month-partial.gohtml"))))

	static := http.FileServer(http.Dir("web/static"))
	mux.Handle("/assets/{path...}", static)
	mux.Handle("/javascript/{path...}", static)
	mux.Handle("/stylesheets/{path...}", static)

	muxWithHeaders := securityheaders.Use(setCSPHeader(mux))

	loggerMiddleware := telemetry.Middleware(logger)
	xsrfMiddleware := xsrfHandler(logger, templates.Get("error.gohtml"), siriusPublicURL)

	return otelhttp.NewHandler(http.StripPrefix(prefix, xsrfMiddleware(loggerMiddleware(muxWithHeaders))), "lpa-frontend")
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

func xsrfHandler(logger *slog.Logger, tmplError template.Template, siriusURL string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				cookieToken := ""

				if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
					cookieToken, _ = url.QueryUnescape(cookie.Value)
				}

				postToken := postFormString(r, "xsrfToken")

				if cookieToken != postToken {
					errorMessage := "Post request was not valid. Please refresh the page and try again."

					w.WriteHeader(http.StatusForbidden)
					_ = tmplError(w, errorVars{
						SiriusURL: siriusURL,
						Path:      r.URL.Path,
						Code:      http.StatusForbidden,
						Error:     errorMessage,
					})
					logger.Warn(errorMessage)

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
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
					urlString := prefix + redirect.To()

					if r.Header.Get("HX-Request") == "true" {
						u, parseErr := url.Parse(urlString)
						if parseErr == nil && u.Fragment != "" {
							q := u.Query()
							q.Set("scrollTo", u.Fragment)
							u.RawQuery = q.Encode()
							urlString = u.String()
						}
					}

					http.Redirect(w, r, urlString, http.StatusFound)
					return
				}

				code := http.StatusInternalServerError
				correlationId := ""
				logger := telemetry.LoggerFromContext(r.Context())

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

func postFormKeySet(r *http.Request, name string) bool {
	if _, val := r.PostForm[name]; val {
		return true
	}
	return false
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

func strToIntOrStatusError(val string) (int, error) {
	if val == "" {
		return 0, sirius.StatusError{Code: http.StatusNotFound}
	}

	i, err := strconv.Atoi(strings.TrimSpace(val))

	if err != nil {
		return 0, sirius.StatusError{Code: http.StatusBadRequest}
	}

	return i, nil
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
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data: s3.eu-west-1.amazonaws.com; frame-src 'self' blob:; object-src 'none'")

		h.ServeHTTP(w, r)
	}
}
