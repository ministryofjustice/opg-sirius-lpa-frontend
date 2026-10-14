package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUpdateDocumentClient struct {
	mock.Mock
}

func (m *mockUpdateDocumentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockUpdateDocumentClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockUpdateDocumentClient) Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int, docTypes []string, notDocTypes []string) ([]sirius.Document, error) {
	args := m.Called(ctx, caseType, caseId, docTypes, notDocTypes)
	return args.Get(0).([]sirius.Document), args.Error(1)
}

func (m *mockUpdateDocumentClient) DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func (m *mockUpdateDocumentClient) UpdateDocument(ctx sirius.Context, uuid string, content string) (sirius.Document, error) {
	args := m.Called(ctx, uuid, content)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func (m *mockUpdateDocumentClient) DeleteDocument(ctx sirius.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *mockUpdateDocumentClient) AddDocument(ctx sirius.Context, caseID int, document sirius.Document, docType string, blankSections []string) (sirius.Document, error) {
	args := m.Called(ctx, caseID, document, docType, blankSections)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func (m *mockUpdateDocumentClient) DocumentTemplates(ctx sirius.Context, caseType sirius.CaseType) ([]sirius.DocumentTemplateData, error) {
	args := m.Called(ctx, caseType)
	return args.Get(0).([]sirius.DocumentTemplateData), args.Error(1)
}

func (m *mockUpdateDocumentClient) FeatureToggles(ctx sirius.Context) (sirius.FeatureToggles, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.FeatureToggles), args.Error(1)
}

func TestGetUpdateDocument(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa", "digital_lpa", "lpa_htmx"} {
		t.Run(caseType, func(t *testing.T) {
			htmx := false
			if caseType == "lpa_htmx" {
				htmx = true
				caseType = "lpa"
			}
			caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: caseType, UID: "7000"}

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

			client := &mockUpdateDocumentClient{}
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
				On("FeatureToggles", mock.Anything).
				Return(sirius.FeatureToggles{}, nil)

			template := &mockTemplate{}
			templateData := updateDocumentData{
				IsPartial:           htmx,
				Case:                caseItem,
				DonorId:             1,
				Documents:           documents,
				Document:            document,
				UsesNotify:          true,
				Section11:           Section11,
				Section10And11:      Section10And11,
				Section15:           Section15,
				Section11And15:      Section11And15,
				Section10And11And15: Section10And11And15,
				Section12And13And15: Section12And13And15,
			}

			if caseType == "digital_lpa" {
				caseSummary := sirius.CaseSummary{DigitalLpa: sirius.DigitalLpa{}, TaskList: []sirius.Task{}}

				client.
					On("CaseSummary", mock.Anything, "7000").
					Return(caseSummary, nil)

				templateData.CaseSummary = caseSummary
			}

			template.
				On("Func", mock.Anything, templateData).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=155&case="+caseType, nil)
			w := httptest.NewRecorder()

			if htmx {
				r.Header.Add("HX-Request", "true")
			}

			err := UpdateDocument(client, template.Func)(w, r)
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
			caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: caseType, UID: "7000"}

			document := sirius.Document{
				ID:      1,
				UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:    sirius.TypeDraft,
				Content: "Test content",
			}

			documents := []sirius.Document{
				document,
			}

			client := &mockUpdateDocumentClient{}
			client.
				On("Case", mock.Anything, 144).
				Return(caseItem, nil)
			client.
				On("Documents", mock.Anything, sirius.CaseType(caseType), 144, []string{sirius.TypeDraft}, []string{}).
				Return(documents, nil)
			client.
				On("UpdateDocument", mock.Anything, document.UUID, "Edited test content").
				Return(document, nil)
			client.
				On("DocumentTemplates", mock.Anything, sirius.CaseType(caseType)).
				Return([]sirius.DocumentTemplateData{}, nil)
			client.
				On("FeatureToggles", mock.Anything).
				Return(sirius.FeatureToggles{}, nil)

			template := &mockTemplate{}
			templateData := updateDocumentData{
				Case:                caseItem,
				DonorId:             1,
				Documents:           documents,
				Document:            document,
				Section11:           Section11,
				Section10And11:      Section10And11,
				Section15:           Section15,
				Section11And15:      Section11And15,
				Section10And11And15: Section10And11And15,
				Section12And13And15: Section12And13And15,
			}

			if caseType == "digital_lpa" {
				caseSummary := sirius.CaseSummary{DigitalLpa: sirius.DigitalLpa{}, TaskList: []sirius.Task{}}

				client.
					On("CaseSummary", mock.Anything, "7000").
					Return(caseSummary, nil)

				templateData.CaseSummary = caseSummary
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

			err := UpdateDocument(client, template.Func)(w, r)
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
			caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: caseType, UID: "700700"}

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

			client := &mockUpdateDocumentClient{}
			client.
				On("Case", mock.Anything, 288).
				Return(caseItem, nil)
			client.
				On("DeleteDocument", mock.Anything, document.UUID).
				Return(nil)
			client.
				On("FeatureToggles", mock.Anything).
				Return(sirius.FeatureToggles{}, nil)

			template := &mockTemplate{}
			errExample = nil

			if caseType == "digital_lpa" {
				errExample = RedirectError("/lpa/700700/documents")
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

				template.
					On("Func", mock.Anything, updateDocumentData{
						Case:                caseItem,
						DonorId:             1,
						Documents:           documents,
						Document:            document,
						Section11:           Section11,
						Section10And11:      Section10And11,
						Section15:           Section15,
						Section11And15:      Section11And15,
						Section10And11And15: Section10And11And15,
						Section12And13And15: Section12And13And15,
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

			err := UpdateDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Equal(t, errExample, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostPublishDocument(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa", "digital_lpa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: caseType, UID: "700700"}

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
			client := &mockUpdateDocumentClient{}
			client.
				On("UpdateDocument", mock.Anything, document.UUID, "Test content").
				Return(document, nil)
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)
			client.
				On("AddDocument", mock.Anything, 544, document, sirius.TypeSave, []string{}).
				Return(publishedDocument, nil)
			client.
				On("DeleteDocument", mock.Anything, document.UUID).
				Return(nil)
			client.
				On("Case", mock.Anything, 544).
				Return(caseItem, nil)
			client.
				On("FeatureToggles", mock.Anything).
				Return(sirius.FeatureToggles{}, nil)

			template := &mockTemplate{}
			errExample = nil

			if caseType == "digital_lpa" {
				errExample = RedirectError("/lpa/700700/documents")
			} else {
				client.
					On("Documents", mock.Anything, sirius.CaseType(caseType), 544, []string{sirius.TypeDraft}, []string{}).
					Return(documents, nil)
				client.
					On("DocumentTemplates", mock.Anything, sirius.CaseType(caseType)).
					Return([]sirius.DocumentTemplateData{}, nil)

				template.
					On("Func", mock.Anything, updateDocumentData{
						Case:                caseItem,
						DonorId:             1,
						Documents:           documents,
						Document:            document,
						Success:             true,
						Section11:           Section11,
						Section10And11:      Section10And11,
						Section15:           Section15,
						Section11And15:      Section11And15,
						Section10And11And15: Section10And11And15,
						Section12And13And15: Section12And13And15,
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

			err := UpdateDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Equal(t, errExample, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostPreviewDocument(t *testing.T) {
	caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: "lpa", UID: "7000"}

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

	client := &mockUpdateDocumentClient{}
	client.
		On("UpdateDocument", mock.Anything, document.UUID, "Test content").
		Return(document, nil)
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("AddDocument", mock.Anything, 888, document, sirius.TypePreview, []string{}).
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
		On("FeatureToggles", mock.Anything).
		Return(sirius.FeatureToggles{}, nil)

	template := &mockTemplate{}

	template.
		On("Func", mock.Anything, updateDocumentData{
			Case:                caseItem,
			DonorId:             1,
			Documents:           documents,
			Document:            document,
			PreviewDraft:        true,
			DownloadUUID:        "efef6714-b4fe-44c2-b26e-90dfe3663e96",
			Section11:           Section11,
			Section10And11:      Section10And11,
			Section15:           Section15,
			Section11And15:      Section11And15,
			Section10And11And15: Section10And11And15,
			Section12And13And15: Section12And13And15,
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

	err := UpdateDocument(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostSaveDocumentAndExit(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa", "digital_lpa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: caseType, UID: "700700"}
			document := sirius.Document{
				ID:      1,
				UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:    sirius.TypeDraft,
				Content: "Test content",
			}

			client := &mockUpdateDocumentClient{}
			client.
				On("UpdateDocument", mock.Anything, document.UUID, "Test content").
				Return(document, nil)
			client.
				On("Case", mock.Anything, 987).
				Return(caseItem, nil)
			client.
				On("FeatureToggles", mock.Anything).
				Return(sirius.FeatureToggles{}, nil)

			template := &mockTemplate{}

			if caseType == "digital_lpa" {
				errExample = RedirectError("/lpa/700700/documents")
			} else {
				errExample = nil
				template.
					On("Func", mock.Anything, updateDocumentData{
						SaveAndExit:         true,
						Case:                caseItem,
						DonorId:             1,
						Section11:           Section11,
						Section10And11:      Section10And11,
						Section15:           Section15,
						Section11And15:      Section11And15,
						Section10And11And15: Section10And11And15,
						Section12And13And15: Section12And13And15,
					}).
					Return(nil)
			}

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

			err := UpdateDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Equal(t, errExample, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetUpdateDocumentBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":    "/?case=lpa",
		"no-case":  "/?id=1111",
		"bad-case": "/?id=1111&case=person",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := UpdateDocument(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetUpdateDocumentWhenCaseErrors(t *testing.T) {
	client := &mockUpdateDocumentClient{}
	client.
		On("Case", mock.Anything, 222).
		Return(sirius.Case{}, errExample)
	client.
		On("FeatureToggles", mock.Anything).
		Return(sirius.FeatureToggles{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=222&case=lpa", nil)
	w := httptest.NewRecorder()

	err := UpdateDocument(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetCreateDocumentWhenFailureOnDocuments(t *testing.T) {
	caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: "lpa", UID: "7000"}

	client := &mockUpdateDocumentClient{}
	client.
		On("Case", mock.Anything, 535).
		Return(caseItem, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseTypeLpa, 535, []string{sirius.TypeDraft}, []string{}).
		Return([]sirius.Document{}, errExample)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return([]sirius.DocumentTemplateData{}, nil)
	client.
		On("FeatureToggles", mock.Anything).
		Return(sirius.FeatureToggles{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=535&case=lpa", nil)
	w := httptest.NewRecorder()

	err := UpdateDocument(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetCreateDocumentWhenFailureOnDocumentByUUID(t *testing.T) {
	caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: "lpa", UID: "7000"}

	document := sirius.Document{
		ID:   1,
		UUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		Type: sirius.TypeDraft,
	}

	documents := []sirius.Document{
		document,
	}

	client := &mockUpdateDocumentClient{}
	client.
		On("Case", mock.Anything, 843).
		Return(caseItem, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseTypeLpa, 843, []string{sirius.TypeDraft}, []string{}).
		Return(documents, nil)
	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(sirius.Document{}, errExample)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return([]sirius.DocumentTemplateData{}, nil)
	client.
		On("FeatureToggles", mock.Anything).
		Return(sirius.FeatureToggles{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=843&case=lpa", nil)
	w := httptest.NewRecorder()

	err := UpdateDocument(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetUpdateDocumentWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: "lpa", UID: "7000"}

	document := sirius.Document{
		ID:   1,
		UUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		Type: sirius.TypeDraft,
	}

	documents := []sirius.Document{
		document,
	}

	client := &mockUpdateDocumentClient{}
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
		On("FeatureToggles", mock.Anything).
		Return(sirius.FeatureToggles{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, updateDocumentData{
			Case:                caseItem,
			DonorId:             1,
			Document:            document,
			Documents:           documents,
			Section11:           Section11,
			Section10And11:      Section10And11,
			Section15:           Section15,
			Section11And15:      Section11And15,
			Section10And11And15: Section10And11And15,
			Section12And13And15: Section12And13And15,
		}).
		Return(errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := UpdateDocument(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostPublishDocumentWithBlankSections(t *testing.T) {
	tests := []struct {
		name             string
		selectedSections string
		blankSections    []string
		section11Count1  string
		section11Count2  string
		section11Count3  string
		section11Count4  string
	}{
		{
			name:             "Section 11",
			selectedSections: Section11,
			blankSections:    []string{"pfa-11", "pfa-11"},
			section11Count1:  "2",
		},
		{
			name:             "Section 10 and 11",
			selectedSections: Section10And11,
			blankSections:    []string{"pfa-10", "pfa-11", "pfa-11", "pfa-11"},
			section11Count2:  "3",
		},
		{
			name:             "Section 11 and 15",
			selectedSections: "11+15",
			blankSections:    []string{"pfa-11", "pfa-11", "pfa-11", "pfa-15"},
			section11Count3:  "3",
		},
		{
			name:             "Section 10, 11 and 15",
			selectedSections: "10+11+15",
			blankSections:    []string{"pfa-10", "pfa-11", "pfa-11", "pfa-15"},
			section11Count4:  "2",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: "lpa", UID: "700700", SubType: "pfa"}

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

			client := &mockUpdateDocumentClient{}
			client.
				On("UpdateDocument", mock.Anything, document.UUID, "Test content").
				Return(document, nil)
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)
			client.
				On("AddDocument", mock.Anything, 544, document, sirius.TypeSave, tc.blankSections).
				Return(publishedDocument, nil).
				Return(publishedDocument, nil)
			client.
				On("DeleteDocument", mock.Anything, document.UUID).
				Return(nil)
			client.
				On("Case", mock.Anything, 544).
				Return(caseItem, nil)
			client.
				On("Documents", mock.Anything, sirius.CaseType("lpa"), 544, []string{sirius.TypeDraft}, []string{}).
				Return(documents, nil)
			client.
				On("DocumentTemplates", mock.Anything, sirius.CaseType("lpa")).
				Return([]sirius.DocumentTemplateData{}, nil)
			client.
				On("FeatureToggles", mock.Anything).
				Return(sirius.FeatureToggles{}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, updateDocumentData{
					Case:                  caseItem,
					DonorId:               1,
					Documents:             documents,
					Document:              document,
					Success:               true,
					HasBlankSections:      "true",
					SelectedBlankSections: tc.selectedSections,
					Section11Count1:       tc.section11Count1,
					Section11Count2:       tc.section11Count2,
					Section11Count3:       tc.section11Count3,
					Section11Count4:       tc.section11Count4,
					Section11:             Section11,
					Section10And11:        Section10And11,
					Section15:             Section15,
					Section11And15:        Section11And15,
					Section10And11And15:   Section10And11And15,
					Section12And13And15:   Section12And13And15,
				}).
				Return(nil)

			form := url.Values{
				"id":                 {"544"},
				"case":               {"lpa"},
				"documentControls":   {"publish"},
				"documentTextEditor": {"Test content"},
				"documentUUID":       {"dfef6714-b4fe-44c2-b26e-90dfe3663e95"},
				"hasBlankSections":   {"true"},
				"blankSections":      {tc.selectedSections},
				"section11Count1":    {tc.section11Count1},
				"section11Count2":    {tc.section11Count2},
				"section11Count3":    {tc.section11Count3},
				"section11Count4":    {tc.section11Count4},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=544&case=lpa", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostDocumentWithBlankSectionsWhenNoSectionsSelected(t *testing.T) {
	for _, postType := range []string{"publish", "preview"} {
		t.Run(postType, func(t *testing.T) {
			caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: "lpa", UID: "700700", SubType: "pfa"}

			document := sirius.Document{
				ID:      1,
				UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:    sirius.TypeDraft,
				Content: "Test content",
			}

			documents := []sirius.Document{
				document,
			}

			client := &mockUpdateDocumentClient{}
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)
			client.
				On("Case", mock.Anything, 544).
				Return(caseItem, nil)
			client.
				On("Documents", mock.Anything, sirius.CaseType("lpa"), 544, []string{sirius.TypeDraft}, []string{}).
				Return(documents, nil)
			client.
				On("DocumentTemplates", mock.Anything, sirius.CaseType("lpa")).
				Return([]sirius.DocumentTemplateData{}, nil)
			client.
				On("FeatureToggles", mock.Anything).
				Return(sirius.FeatureToggles{}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, updateDocumentData{
					Case:                caseItem,
					DonorId:             1,
					Documents:           documents,
					Document:            document,
					HasBlankSections:    "true",
					Section11:           Section11,
					Section10And11:      Section10And11,
					Section15:           Section15,
					Section11And15:      Section11And15,
					Section10And11And15: Section10And11And15,
					Section12And13And15: Section12And13And15,
					Error: sirius.ValidationError{
						Field: sirius.FieldErrors{
							"blankSections": {"reason": "Please select sections to insert"},
						},
					},
				}).
				Return(nil)

			form := url.Values{
				"id":                 {"544"},
				"case":               {"lpa"},
				"documentControls":   {postType},
				"documentTextEditor": {"Test content"},
				"documentUUID":       {"dfef6714-b4fe-44c2-b26e-90dfe3663e95"},
				"hasBlankSections":   {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=544&case=lpa", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostDocumentWithBlankSectionsWhenSection11SelectedButNoCount(t *testing.T) {
	for _, postType := range []string{"publish", "preview"} {
		t.Run(postType, func(t *testing.T) {
			caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: "lpa", UID: "700700", SubType: "pfa"}

			document := sirius.Document{
				ID:      1,
				UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:    sirius.TypeDraft,
				Content: "Test content",
			}

			documents := []sirius.Document{
				document,
			}

			client := &mockUpdateDocumentClient{}
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)
			client.
				On("Case", mock.Anything, 544).
				Return(caseItem, nil)
			client.
				On("Documents", mock.Anything, sirius.CaseType("lpa"), 544, []string{sirius.TypeDraft}, []string{}).
				Return(documents, nil)
			client.
				On("DocumentTemplates", mock.Anything, sirius.CaseType("lpa")).
				Return([]sirius.DocumentTemplateData{}, nil)
			client.
				On("FeatureToggles", mock.Anything).
				Return(sirius.FeatureToggles{}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, updateDocumentData{
					Case:                  caseItem,
					DonorId:               1,
					Documents:             documents,
					Document:              document,
					HasBlankSections:      "true",
					SelectedBlankSections: "10+11+15",
					Section11:             Section11,
					Section10And11:        Section10And11,
					Section15:             Section15,
					Section11And15:        Section11And15,
					Section10And11And15:   Section10And11And15,
					Section12And13And15:   Section12And13And15,
					Error: sirius.ValidationError{
						Field: sirius.FieldErrors{
							"section11Count4": {"reason": "Please select how many section 11 to insert"},
						},
					},
				}).
				Return(nil)

			form := url.Values{
				"id":                 {"544"},
				"case":               {"lpa"},
				"documentControls":   {postType},
				"documentTextEditor": {"Test content"},
				"documentUUID":       {"dfef6714-b4fe-44c2-b26e-90dfe3663e95"},
				"hasBlankSections":   {"true"},
				"blankSections":      {"10+11+15"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=544&case=lpa", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostDocumentWithBlankSectionsWhenSection11And15SelectedButNoCount(t *testing.T) {
	for _, postType := range []string{"publish", "preview"} {
		t.Run(postType, func(t *testing.T) {
			caseItem := sirius.Case{Donor: &sirius.Person{ID: 1}, CaseType: "lpa", UID: "700700", SubType: "pfa"}

			document := sirius.Document{
				ID:      1,
				UUID:    "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:    sirius.TypeDraft,
				Content: "Test content",
			}

			documents := []sirius.Document{
				document,
			}

			client := &mockUpdateDocumentClient{}
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)
			client.
				On("Case", mock.Anything, 544).
				Return(caseItem, nil)
			client.
				On("Documents", mock.Anything, sirius.CaseType("lpa"), 544, []string{sirius.TypeDraft}, []string{}).
				Return(documents, nil)
			client.
				On("DocumentTemplates", mock.Anything, sirius.CaseType("lpa")).
				Return([]sirius.DocumentTemplateData{}, nil)
			client.
				On("FeatureToggles", mock.Anything).
				Return(sirius.FeatureToggles{}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, updateDocumentData{
					Case:                  caseItem,
					DonorId:               1,
					Documents:             documents,
					Document:              document,
					HasBlankSections:      "true",
					SelectedBlankSections: "11+15",
					Section11:             Section11,
					Section10And11:        Section10And11,
					Section15:             Section15,
					Section11And15:        Section11And15,
					Section10And11And15:   Section10And11And15,
					Section12And13And15:   Section12And13And15,
					Error: sirius.ValidationError{
						Field: sirius.FieldErrors{
							"section11Count3": {"reason": "Please select how many section 11 to insert"},
						},
					},
				}).
				Return(nil)

			form := url.Values{
				"id":                 {"544"},
				"case":               {"lpa"},
				"documentControls":   {postType},
				"documentTextEditor": {"Test content"},
				"documentUUID":       {"dfef6714-b4fe-44c2-b26e-90dfe3663e95"},
				"hasBlankSections":   {"true"},
				"blankSections":      {"11+15"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=544&case=lpa", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestParseBlankSectionsRequiresSection11Count(t *testing.T) {
	tests := []struct {
		selectedSections string
		field            string
	}{
		{selectedSections: Section11, field: "section11Count1"},
		{selectedSections: Section10And11, field: "section11Count2"},
		{selectedSections: Section11And15, field: "section11Count3"},
		{selectedSections: Section10And11And15, field: "section11Count4"},
	}

	for _, tc := range tests {
		t.Run(tc.selectedSections, func(t *testing.T) {
			_, err := parseBlankSections(true, tc.selectedSections, section11Counts{}, "pfa")

			assert.Equal(t, sirius.ValidationError{
				Field: sirius.FieldErrors{
					tc.field: {"reason": "Please select how many section 11 to insert"},
				},
			}, err)
		})
	}
}

func TestGetAttorneyCount(t *testing.T) {
	caseItem := sirius.Case{
		Donor:    &sirius.Person{ID: 1},
		CaseType: "lpa",
		UID:      "7000",
		Attorneys: []sirius.Attorney{
			{Person: sirius.Person{ID: 1}},
			{Person: sirius.Person{ID: 2}, SystemStatus: shared.BoolPtr(false)},
			{Person: sirius.Person{ID: 3}, SystemStatus: shared.BoolPtr(true)},
		},
		ReplacementAttorneys: []sirius.Attorney{{Person: sirius.Person{ID: 4}}},
		TrustCorporations: []sirius.TrustCorporation{
			{Attorney: sirius.Attorney{Person: sirius.Person{ID: 5}}},
			{Attorney: sirius.Attorney{Person: sirius.Person{ID: 6}, SystemStatus: shared.BoolPtr(false)}},
			{Attorney: sirius.Attorney{Person: sirius.Person{ID: 7}, SystemStatus: shared.BoolPtr(true)}},
			{Attorney: sirius.Attorney{Person: sirius.Person{ID: 8}}, IsReplacementAttorney: true},
		},
	}

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

	client := &mockUpdateDocumentClient{}
	client.
		On("Case", mock.Anything, 155).
		Return(caseItem, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseType("lpa"), 155, []string{sirius.TypeDraft}, []string{}).
		Return(documents, nil)
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseType("lpa")).
		Return(documentTemplates, nil)
	client.
		On("FeatureToggles", mock.Anything).
		Return(sirius.FeatureToggles{}, nil)

	template := &mockTemplate{}
	templateData := updateDocumentData{
		Case:                caseItem,
		DonorId:             1,
		Documents:           documents,
		Document:            document,
		UsesNotify:          true,
		AttorneyCount:       4,
		Section11:           Section11,
		Section10And11:      Section10And11,
		Section15:           Section15,
		Section11And15:      Section11And15,
		Section10And11And15: Section10And11And15,
		Section12And13And15: Section12And13And15,
	}

	template.
		On("Func", mock.Anything, templateData).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=155&case=lpa", nil)
	w := httptest.NewRecorder()

	err := UpdateDocument(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
