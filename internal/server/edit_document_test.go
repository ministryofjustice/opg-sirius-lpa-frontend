package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEditDocumentClient struct {
	mock.Mock
}

func (m *mockEditDocumentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockEditDocumentClient) DigitalLpa(ctx sirius.Context, uid string) (sirius.DigitalLpa, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.DigitalLpa), args.Error(1)
}

func (m *mockEditDocumentClient) TasksForCase(ctx sirius.Context, id int) ([]sirius.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.Task), args.Error(1)
}

func (m *mockEditDocumentClient) Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int, docTypes []string, notDocTypes []string) ([]sirius.Document, error) {
	args := m.Called(ctx, caseType, caseId, docTypes, notDocTypes)
	return args.Get(0).([]sirius.Document), args.Error(1)
}

func (m *mockEditDocumentClient) DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func (m *mockEditDocumentClient) EditDocument(ctx sirius.Context, uuid string, content string) (sirius.Document, error) {
	args := m.Called(ctx, uuid, content)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func (m *mockEditDocumentClient) DeleteDocument(ctx sirius.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *mockEditDocumentClient) AddDocument(ctx sirius.Context, caseID int, document sirius.Document, docType string) (sirius.Document, error) {
	args := m.Called(ctx, caseID, document, docType)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func (m *mockEditDocumentClient) DocumentTemplates(ctx sirius.Context, caseType sirius.CaseType) ([]sirius.DocumentTemplateData, error) {
	args := m.Called(ctx, caseType)
	return args.Get(0).([]sirius.DocumentTemplateData), args.Error(1)
}

func TestGetEditDocument(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa", "digital_lpa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "7000"}

			document := sirius.Document{
				ID:         1,
				UUID:       "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				SystemType: "LP-LETTER",
				Type:       sirius.TypeDraft,
			}

			documents := []sirius.Document{
				document,
			}

			documentTemplates := []sirius.DocumentTemplateData{
				{
					TemplateId: "LP-LETTER",
					UsesNotify: true,
				},
			}

			client := &mockEditDocumentClient{}
			client.
				On("Case", mock.Anything, 155).
				Return(caseItem, nil)
			client.
				On("Documents", mock.Anything, sirius.CaseType(caseType), 155, []string{sirius.TypeDraft}, []string{}).
				Return(documents, nil)
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)
			client.
				On("DocumentTemplates", mock.Anything, sirius.CaseType(caseType)).
				Return(documentTemplates, nil)
			client.
				On("TasksForCase", mock.Anything, 155).
				Return([]sirius.Task{}, nil)

			template := &mockTemplate{}
			templateData := editDocumentData{
				Case:       caseItem,
				Documents:  documents,
				Document:   document,
				UsesNotify: true,
				TaskList:   []sirius.Task{},
			}

			if caseType == "digital_lpa" {
				lpa := sirius.DigitalLpa{}

				client.
					On("DigitalLpa", mock.Anything, "7000").
					Return(lpa, nil)

				templateData.Lpa = lpa
			}

			template.
				On("Func", mock.Anything, templateData).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=155&case="+caseType, nil)
			w := httptest.NewRecorder()

			err := EditDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostSaveDocument(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa", "digital_lpa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "7000"}

			document := sirius.Document{
				ID:      1,
				UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:    sirius.TypeDraft,
				Content: "Test content",
			}

			documents := []sirius.Document{
				document,
			}

			client := &mockEditDocumentClient{}
			client.
				On("Case", mock.Anything, 144).
				Return(caseItem, nil)
			client.
				On("Documents", mock.Anything, sirius.CaseType(caseType), 144, []string{sirius.TypeDraft}, []string{}).
				Return(documents, nil)
			client.
				On("EditDocument", mock.Anything, document.UUID, "Edited test content").
				Return(document, nil)
			client.
				On("DocumentTemplates", mock.Anything, sirius.CaseType(caseType)).
				Return([]sirius.DocumentTemplateData{}, nil)
			client.
				On("TasksForCase", mock.Anything, 144).
				Return([]sirius.Task{}, nil)

			template := &mockTemplate{}
			templateData := editDocumentData{
				Case:      caseItem,
				Documents: documents,
				Document:  document,
				TaskList:  []sirius.Task{},
			}

			if caseType == "digital_lpa" {
				lpa := sirius.DigitalLpa{}

				client.
					On("DigitalLpa", mock.Anything, "7000").
					Return(lpa, nil)

				templateData.Lpa = lpa
			}

			template.
				On("Func", mock.Anything, templateData).
				Return(nil)

			form := url.Values{
				"id":                 {"144"},
				"case":               {caseType},
				"documentControls":   {"save"},
				"documentTextEditor": {"Edited test content"},
				"documentUUID":       {"dfef6714-b4fe-44c2-b26e-90dfe3663e95"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=144&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := EditDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostDeleteDocument(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa", "digital_lpa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "700700"}

			document := sirius.Document{
				ID:      1,
				UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:    sirius.TypeDraft,
				Content: "Test content",
			}

			documents := []sirius.Document{
				document,
				{
					ID:      2,
					UUID:    "efef6714-b4fe-44c2-b26e-90dfe3663e96",
					Type:    sirius.TypeDraft,
					Content: "Some more content",
				},
			}

			client := &mockEditDocumentClient{}
			client.
				On("Case", mock.Anything, 288).
				Return(caseItem, nil)
			client.
				On("DeleteDocument", mock.Anything, document.UUID).
				Return(nil)

			template := &mockTemplate{}
			expectedError = nil

			if caseType == "digital_lpa" {
				expectedError = RedirectError("/lpa/700700/documents")
			} else {
				client.
					On("Documents", mock.Anything, sirius.CaseType(caseType), 288, []string{sirius.TypeDraft}, []string{}).
					Return(documents, nil)
				client.
					On("DocumentByUUID", mock.Anything, document.UUID).
					Return(document, nil)
				client.
					On("DocumentTemplates", mock.Anything, sirius.CaseType(caseType)).
					Return([]sirius.DocumentTemplateData{}, nil)
				client.
					On("TasksForCase", mock.Anything, 288).
					Return([]sirius.Task{}, nil)

				template.
					On("Func", mock.Anything, editDocumentData{
						Case:      caseItem,
						Documents: documents,
						Document:  document,
						TaskList:  []sirius.Task{},
					}).
					Return(nil)
			}

			form := url.Values{
				"id":                 {"288"},
				"case":               {caseType},
				"documentControls":   {"delete"},
				"documentTextEditor": {"Test content"},
				"documentUUID":       {"dfef6714-b4fe-44c2-b26e-90dfe3663e95"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=288&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := EditDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Equal(t, expectedError, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostPublishDocument(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa", "digital_lpa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "700700"}

			document := sirius.Document{
				ID:      1,
				UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:    sirius.TypeDraft,
				Content: "Test content",
			}

			documents := []sirius.Document{
				document,
			}

			publishedDocument := sirius.Document{
				ID:      1,
				UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:    sirius.TypeSave,
				Content: "Test content",
			}
			client := &mockEditDocumentClient{}
			client.
				On("EditDocument", mock.Anything, document.UUID, "Test content").
				Return(document, nil)
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)
			client.
				On("AddDocument", mock.Anything, 544, document, sirius.TypeSave).
				Return(publishedDocument, nil)
			client.
				On("DeleteDocument", mock.Anything, document.UUID).
				Return(nil)
			client.
				On("Case", mock.Anything, 544).
				Return(caseItem, nil)

			template := &mockTemplate{}
			expectedError = nil

			if caseType == "digital_lpa" {
				expectedError = RedirectError("/lpa/700700/documents")
			} else {
				client.
					On("Documents", mock.Anything, sirius.CaseType(caseType), 544, []string{sirius.TypeDraft}, []string{}).
					Return(documents, nil)
				client.
					On("DocumentTemplates", mock.Anything, sirius.CaseType(caseType)).
					Return([]sirius.DocumentTemplateData{}, nil)
				client.
					On("TasksForCase", mock.Anything, 544).
					Return([]sirius.Task{}, nil)

				template.
					On("Func", mock.Anything, editDocumentData{
						Case:      caseItem,
						Documents: documents,
						Document:  document,
						Success:   true,
						TaskList:  []sirius.Task{},
					}).
					Return(nil)

			}

			form := url.Values{
				"id":                 {"544"},
				"case":               {caseType},
				"documentControls":   {"publish"},
				"documentTextEditor": {"Test content"},
				"documentUUID":       {"dfef6714-b4fe-44c2-b26e-90dfe3663e95"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=544&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := EditDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Equal(t, expectedError, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostPreviewDocument(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "700700"}

	document := sirius.Document{
		ID:      1,
		UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		Type:    sirius.TypeDraft,
		Content: "Test content",
	}

	documents := []sirius.Document{
		document,
	}

	previewDocument := sirius.Document{
		ID:      1,
		UUID:    "efef6714-b4fe-44c2-b26e-90dfe3663e96",
		Type:    sirius.TypePreview,
		Content: "Test content",
	}

	client := &mockEditDocumentClient{}
	client.
		On("EditDocument", mock.Anything, document.UUID, "Test content").
		Return(document, nil)
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("AddDocument", mock.Anything, 888, document, sirius.TypePreview).
		Return(previewDocument, nil)
	client.
		On("Case", mock.Anything, 888).
		Return(caseItem, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseType("lpa"), 888, []string{sirius.TypeDraft}, []string{}).
		Return(documents, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return([]sirius.DocumentTemplateData{}, nil)
	client.
		On("TasksForCase", mock.Anything, 888).
		Return([]sirius.Task{}, nil)

	template := &mockTemplate{}

	template.
		On("Func", mock.Anything, editDocumentData{
			Case:         caseItem,
			Documents:    documents,
			Document:     document,
			PreviewDraft: true,
			DownloadUUID: "efef6714-b4fe-44c2-b26e-90dfe3663e96",
			TaskList:     []sirius.Task{},
		}).
		Return(nil)

	form := url.Values{
		"id":                 {"888"},
		"case":               {"lpa"},
		"documentControls":   {"preview"},
		"documentTextEditor": {"Test content"},
		"documentUUID":       {"dfef6714-b4fe-44c2-b26e-90dfe3663e95"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=888&case=lpa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditDocument(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostSaveDocumentAndExit(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa", "digital_lpa"} {
		t.Run(caseType, func(t *testing.T) {
			document := sirius.Document{
				ID:      1,
				UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:    sirius.TypeDraft,
				Content: "Test content",
			}

			client := &mockEditDocumentClient{}
			client.
				On("EditDocument", mock.Anything, document.UUID, "Test content").
				Return(document, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, editDocumentData{
					SaveAndExit: true,
				}).
				Return(nil)

			form := url.Values{
				"id":                 {"987"},
				"case":               {caseType},
				"documentControls":   {"saveAndExit"},
				"documentTextEditor": {"Test content"},
				"documentUUID":       {"dfef6714-b4fe-44c2-b26e-90dfe3663e95"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=987&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := EditDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetEditDocumentBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":    "/?case=lpa",
		"no-case":  "/?id=1111",
		"bad-case": "/?id=1111&case=person",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := EditDocument(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetEditDocumentWhenCaseErrors(t *testing.T) {
	client := &mockEditDocumentClient{}
	client.
		On("Case", mock.Anything, 222).
		Return(sirius.Case{}, expectedError)
	client.
		On("Documents", mock.Anything, sirius.CaseTypeLpa, 222, []string{sirius.TypeDraft}, []string{}).
		Return([]sirius.Document{}, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return([]sirius.DocumentTemplateData{}, nil)
	client.
		On("TasksForCase", mock.Anything, 222).
		Return([]sirius.Task{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=222&case=lpa", nil)
	w := httptest.NewRecorder()

	err := EditDocument(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetCreateDocumentWhenFailureOnDocuments(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000"}

	client := &mockEditDocumentClient{}
	client.
		On("Case", mock.Anything, 535).
		Return(caseItem, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseTypeLpa, 535, []string{sirius.TypeDraft}, []string{}).
		Return([]sirius.Document{}, expectedError)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return([]sirius.DocumentTemplateData{}, nil)
	client.
		On("TasksForCase", mock.Anything, 535).
		Return([]sirius.Task{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=535&case=lpa", nil)
	w := httptest.NewRecorder()

	err := EditDocument(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetCreateDocumentWhenFailureOnDocumentByUUID(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000"}

	document := sirius.Document{
		ID:   1,
		UUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		Type: sirius.TypeDraft,
	}

	documents := []sirius.Document{
		document,
	}

	client := &mockEditDocumentClient{}
	client.
		On("Case", mock.Anything, 843).
		Return(caseItem, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseTypeLpa, 843, []string{sirius.TypeDraft}, []string{}).
		Return(documents, nil)
	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(sirius.Document{}, expectedError)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return([]sirius.DocumentTemplateData{}, nil)
	client.
		On("TasksForCase", mock.Anything, 843).
		Return([]sirius.Task{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=843&case=lpa", nil)
	w := httptest.NewRecorder()

	err := EditDocument(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetEditDocumentWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000"}

	document := sirius.Document{
		ID:   1,
		UUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		Type: sirius.TypeDraft,
	}

	documents := []sirius.Document{
		document,
	}

	client := &mockEditDocumentClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseTypeLpa, 123, []string{sirius.TypeDraft}, []string{}).
		Return(documents, nil)
	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(document, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return([]sirius.DocumentTemplateData{}, nil)
	client.
		On("TasksForCase", mock.Anything, 123).
		Return([]sirius.Task{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editDocumentData{
			Case:      caseItem,
			Document:  document,
			Documents: documents,
			TaskList:  []sirius.Task{},
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := EditDocument(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
