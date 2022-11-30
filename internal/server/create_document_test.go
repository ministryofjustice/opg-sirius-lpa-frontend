package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateDocumentClient struct {
	mock.Mock
}

func (m *mockCreateDocumentClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockCreateDocumentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockCreateDocumentClient) DocumentTemplates(ctx sirius.Context, caseType sirius.CaseType) ([]sirius.DocumentTemplateData, error) {
	args := m.Called(ctx, caseType)
	return args.Get(0).([]sirius.DocumentTemplateData), args.Error(1)
}

func TestGetCreateDocument(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "7000"}

			documentTemplates := []sirius.RefDataItem{
				{
					Handle: "DD",
					Label:  "DD Template Label",
				},
			}

			documentTemplateData := []sirius.DocumentTemplateData{
				{
					Inserts:    nil,
					TemplateId: "DD",
					UniversalTemplateData: sirius.UniversalTemplateData{
						Location:        "DD.html.twig",
						OnScreenSummary: "DDONSCREENSUMMARY",
					},
				},
			}

			client := &mockCreateDocumentClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseItem, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.DocumentTemplateIdCategory).
				Return(documentTemplates, nil)
			client.
				On("DocumentTemplates", mock.Anything, sirius.CaseType(caseType)).
				Return(documentTemplateData, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createDocumentData{
					Case:                    caseItem,
					DocumentTemplateRefData: documentTemplates,
					DocumentTemplates:       documentTemplateData,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&case="+caseType, nil)
			w := httptest.NewRecorder()

			err := CreateDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateDocumentBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":    "/?case=lpa",
		"no-case":  "/?id=123",
		"bad-case": "/?id=123&case=person",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := CreateDocument(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetCreateDocumentWhenCaseErrors(t *testing.T) {
	client := &mockCreateDocumentClient{}
	documentTemplates := []sirius.RefDataItem{
		{
			Handle: "DD",
			Label:  "Donor deceased: Blank template",
		},
	}

	documentTemplateData := []sirius.DocumentTemplateData{
		{
			Inserts:    nil,
			TemplateId: "DD",
			UniversalTemplateData: sirius.UniversalTemplateData{
				Location:        `lpa\/DD.html.twig`,
				OnScreenSummary: "DDONSCREENSUMMARY",
			},
		},
	}
	client.
		On("RefDataByCategory", mock.Anything, sirius.DocumentTemplateIdCategory).
		Return(documentTemplates, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return(documentTemplateData, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := CreateDocument(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetCreateDocumentWhenFailureOnGetDocumentRefData(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000"}

	documentTemplateData := []sirius.DocumentTemplateData{
		{
			Inserts:    nil,
			TemplateId: "DD",
			UniversalTemplateData: sirius.UniversalTemplateData{
				Location:        `lpa\/DD.html.twig`,
				OnScreenSummary: "DDONSCREENSUMMARY",
			},
		},
	}

	client := &mockCreateDocumentClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return(documentTemplateData, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.DocumentTemplateIdCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := CreateDocument(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetCreateDocumentWhenFailureOnGetDocumentTemplates(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000"}

	documentTemplates := []sirius.RefDataItem{
		{
			Handle: "DD",
			Label:  "Donor deceased: Blank template",
		},
	}

	client := &mockCreateDocumentClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.DocumentTemplateIdCategory).
		Return(documentTemplates, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return([]sirius.DocumentTemplateData{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := CreateDocument(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetCreateDocumentWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000"}

	documentTemplates := []sirius.RefDataItem{
		{
			Handle: "DD",
			Label:  "Donor deceased: Blank template",
		},
	}

	documentTemplateData := []sirius.DocumentTemplateData{
		{
			Inserts:    nil,
			TemplateId: "DD",
			UniversalTemplateData: sirius.UniversalTemplateData{
				Location:        `lpa\/DD.html.twig`,
				OnScreenSummary: "DDONSCREENSUMMARY",
			},
		},
	}

	client := &mockCreateDocumentClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.DocumentTemplateIdCategory).
		Return(documentTemplates, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeLpa).
		Return(documentTemplateData, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDocumentData{
			Case:                    caseItem,
			DocumentTemplateRefData: documentTemplates,
			DocumentTemplates:       documentTemplateData,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := CreateDocument(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestTranslateDocumentData(t *testing.T) {
	documentTemplateRefData := []sirius.RefDataItem{
		{
			Handle: "DDONSCREENSUMMARY",
			Label:  "DD Template Label",
		},
	}

	documentTemplateData := []sirius.DocumentTemplateData{
		{
			Inserts:    nil,
			TemplateId: "DD",
			UniversalTemplateData: sirius.UniversalTemplateData{
				Location:        `DD.html.twig`,
				OnScreenSummary: "DDONSCREENSUMMARY",
			},
		},
	}

	documentTemplateTypes := translateDocumentData(documentTemplateData, documentTemplateRefData)
	assert.Equal(t, documentTemplateTypes[0].Label, "DD Template Label")
	assert.Equal(t, documentTemplateTypes[0].Handle, "DD")
	assert.Equal(t, documentTemplateTypes[0].UserSelectable, false)
}

func TestTranslateInsertData(t *testing.T) {
	documentTemplateRefData := []sirius.RefDataItem{
		{
			Handle: "DDINSERTONSCREENSUMMARY",
			Label:  "DD Insert label",
		},
	}

	selectedTemplateInserts := []sirius.Insert{
		{
			Key:      "All",
			InsertId: "DDINSERT",
			UniversalTemplateData: sirius.UniversalTemplateData{
				Location:        `lpa\/DD.html.twig`,
				OnScreenSummary: "DDINSERTONSCREENSUMMARY",
			},
		},
	}

	translatedInsert := translateInsertData(selectedTemplateInserts, documentTemplateRefData)
	assert.Equal(t, translatedInsert[0].Label, "DD Insert label")
	assert.Equal(t, translatedInsert[0].Handle, "DDINSERT")
	assert.Equal(t, translatedInsert[0].Key, "All")
}
