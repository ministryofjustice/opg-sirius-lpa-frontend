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

func (m *mockCreateDocumentClient) CreateContact(ctx sirius.Context, contact sirius.Person) (sirius.Person, error) {
	args := m.Called(ctx, contact)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockCreateDocumentClient) CreateDocument(ctx sirius.Context, caseID, correspondentID int, templateID string, inserts []string) (sirius.DocumentData, error) {
	args := m.Called(ctx, caseID, correspondentID, templateID, inserts)
	return args.Get(0).(sirius.DocumentData), args.Error(1)
}

func (m *mockCreateDocumentClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
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
					Inserts:         nil,
					TemplateId:      "DD",
					Location:        "DD.html.twig",
					OnScreenSummary: "DDONSCREENSUMMARY",
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

func TestPostCreateDocument(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "7000"}

			client := &mockCreateDocumentClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseItem, nil)
			client.
				On("CreateDocument", mock.Anything, 123, 1, "DD", []string{"DDINSERT"}).
				Return(sirius.DocumentData{}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createDocumentData{
					Case:    caseItem,
					Success: true,
				}).
				Return(nil)

			form := url.Values{
				"id":                {"123"},
				"case":              {caseType},
				"templateId":        {"DD"},
				"selectRecipients":  {"1"},
				"recipientControls": {"select"},
				"insert":            {"DDINSERT"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateDocumentGenerateNewRecipient(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			donor := sirius.Person{ID: 1}
			caseItem := sirius.Case{CaseType: caseType, UID: "7000", Donor: &donor}

			selectedTemplate := sirius.DocumentTemplateData{
				TemplateId: "DD",
			}

			contact := sirius.Person{
				CompanyName:           "Test Company Name",
				CompanyReference:      "Test Company Reference",
				AddressLine1:          "278 Nicole Lock",
				AddressLine2:          "Toby Court",
				AddressLine3:          "",
				Town:                  "Russellstad",
				County:                "Cumbria",
				Postcode:              "HP19 9BW",
				Country:               "",
				IsAirmailRequired:     false,
				PhoneNumber:           "072345678",
				Email:                 "test.company@uk.test",
				CorrespondenceByPost:  true,
				CorrespondenceByEmail: true,
				CorrespondenceByPhone: false,
				CorrespondenceByWelsh: false,
			}

			client := &mockCreateDocumentClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseItem, nil)
			client.
				On("CreateContact", mock.Anything, contact).
				Return(contact, nil)
			client.
				On("Person", mock.Anything, 1).
				Return(donor, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createDocumentData{
					Case:                caseItem,
					TemplateSelected:    selectedTemplate,
					Success:             true,
					SelectedInserts:     []string{"DDINSERT"},
					HasViewedInsertPage: true,
					Recipients:          []sirius.Person{donor, contact},
				}).
				Return(nil)

			form := url.Values{
				"id":                {"123"},
				"case":              {caseType},
				"templateId":        {"DD"},
				"recipientControls": {"generate"},
				"insert":            {"DDINSERT"},
				"companyName":       {"Test Company Name"},
				"companyReference":  {"Test Company Reference"},
				"addressLine1":      {"278 Nicole Lock"},
				"addressLine2":      {"Toby Court"},
				"addressLine3":      {""},
				"town":              {"Russellstad"},
				"county":            {"Cumbria"},
				"postcode":          {"HP19 9BW"},
				"isAirmailRequired": {"false"},
				"phoneNumber":       {"072345678"},
				"email":             {"test.company@uk.test"},
				"correspondenceBy":  {"post", "email"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
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
			Inserts:         nil,
			TemplateId:      "DD",
			Location:        `lpa\/DD.html.twig`,
			OnScreenSummary: "DDONSCREENSUMMARY",
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
			Inserts:         nil,
			TemplateId:      "DD",
			Location:        `lpa\/DD.html.twig`,
			OnScreenSummary: "DDONSCREENSUMMARY",
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

func TestGetCreateDocumentWhenPersonErrors(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000", Donor: &sirius.Person{ID: 1}}

	documentTemplates := []sirius.RefDataItem{
		{
			Handle: "DD",
			Label:  "Donor deceased: Blank template",
		},
	}

	documentTemplateData := []sirius.DocumentTemplateData{
		{
			Inserts:         nil,
			TemplateId:      "DD",
			Location:        `lpa\/DD.html.twig`,
			OnScreenSummary: "DDONSCREENSUMMARY",
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
	client.
		On("Person", mock.Anything, 1).
		Return(sirius.Person{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa&hasViewedInserts=true", nil)
	w := httptest.NewRecorder()

	err := CreateDocument(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
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
			Inserts:         nil,
			TemplateId:      "DD",
			Location:        `DD.html.twig`,
			OnScreenSummary: "DDONSCREENSUMMARY",
		},
	}

	documentTemplateTypes := translateDocumentData(documentTemplateData, documentTemplateRefData)
	assert.Equal(t, "DD Template Label", documentTemplateTypes[0].Label)
	assert.Equal(t, "DD", documentTemplateTypes[0].Handle)
	assert.Equal(t, false, documentTemplateTypes[0].UserSelectable)
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
			Key:             "All",
			InsertId:        "DDINSERT",
			Location:        `lpa\/DD.html.twig`,
			OnScreenSummary: "DDINSERTONSCREENSUMMARY",
		},
	}

	translatedInsert := translateInsertData(selectedTemplateInserts, documentTemplateRefData)
	assert.Equal(t, "DD Insert label", translatedInsert[0].Label)
	assert.Equal(t, "DDINSERT", translatedInsert[0].Handle)
	assert.Equal(t, "All", translatedInsert[0].Key)
}

func TestGetRecipients(t *testing.T) {
	donor := sirius.Person{ID: 1}
	trustCorp := sirius.Person{ID: 2}
	attorney := sirius.Person{ID: 3}
	caseItem := sirius.Case{Donor: &donor, TrustCorporations: []sirius.Person{trustCorp}, Attorneys: []sirius.Person{attorney}}

	r, _ := http.NewRequest("GET", "/", nil)
	ctx := getContext(r)
	client := &mockCreateDocumentClient{}
	client.
		On("Person", mock.Anything, 1).
		Return(donor, nil)
	client.
		On("Person", mock.Anything, 2).
		Return(trustCorp, nil)
	client.
		On("Person", mock.Anything, 3).
		Return(attorney, nil)

	recipients, _ := getRecipients(ctx, client, caseItem)
	assert.Equal(t, 3, len(recipients))
}

func TestSliceAtoi(t *testing.T) {
	testSliceStr := []string{"1", "2", "3"}
	result, err := sliceAtoi(testSliceStr)

	assert.Equal(t, nil, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestRemoveDuplicateStr(t *testing.T) {
	testSliceStr := []string{"test", "test", "testing"}
	result := removeDuplicateStr(testSliceStr)

	assert.Equal(t, []string{"test", "testing"}, result)
}

func TestGetSortedInsertKeys(t *testing.T) {
	selectedTemplateInsertsWithAll := []sirius.Insert{
		{
			Key:             "all",
			InsertId:        "DDINSERT",
			Location:        `lpa\/DD.html.twig`,
			OnScreenSummary: "DDINSERTONSCREENSUMMARY",
		},
		{
			Key:             "imperfect",
			InsertId:        "IM1INSERT",
			Location:        `lpa\/IM1.html.twig`,
			OnScreenSummary: "IM1INSERTONSCREENSUMMARY",
		},
		{
			Key:             "perfect",
			InsertId:        "P1INSERT",
			Location:        `lpa\/P1.html.twig`,
			OnScreenSummary: "P1INSERTONSCREENSUMMARY",
		},
	}

	selectedTemplateInsertsWithoutAll := []sirius.Insert{
		{
			Key:             "imperfect",
			InsertId:        "IM1INSERT",
			Location:        `lpa\/IM1.html.twig`,
			OnScreenSummary: "IM1INSERTONSCREENSUMMARY",
		},
		{
			Key:             "perfect",
			InsertId:        "P1INSERT",
			Location:        `lpa\/P1.html.twig`,
			OnScreenSummary: "P1INSERTONSCREENSUMMARY",
		},
	}

	result1 := getSortedInsertKeys(selectedTemplateInsertsWithAll)
	result2 := getSortedInsertKeys(selectedTemplateInsertsWithoutAll)

	assert.Equal(t, []string{"all", "imperfect", "perfect"}, result1)
	assert.Equal(t, []string{"all", "imperfect", "perfect"}, result2)
}
