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

type mockAddComplaintClient struct {
	mock.Mock
}

func (m *mockAddComplaintClient) AddComplaint(ctx sirius.Context, caseID int, caseType sirius.CaseType, complaint sirius.Complaint) error {
	args := m.Called(ctx, caseID, caseType, complaint)
	return args.Error(0)
}

func (m *mockAddComplaintClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockAddComplaintClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

var demoComplainantCategories = []sirius.RefDataItem{
	{
		Handle: "LPA_DONOR",
		Label:  "LPA Donor",
	},
}

var demoComplaintOrigins = []sirius.RefDataItem{
	{
		Handle: "PHONE",
		Label:  "Phone call",
	},
}

var demoComplaintCategories = []sirius.RefDataItem{
	{
		Handle: "02",
		Label:  "OPG Decisions",
		Subcategories: []sirius.RefDataItem{
			{
				Handle: "18",
				Label:  "Fee Decision",
			},
		},
	},
}

func TestGetAddComplaint(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			client := &mockAddComplaintClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(sirius.Case{CaseType: caseType, UID: "7000"}, nil).
				On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
				Return(demoComplainantCategories, nil).
				On("RefDataByCategory", mock.Anything, sirius.ComplaintCategory).
				Return(demoComplaintCategories, nil).
				On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
				Return(demoComplaintOrigins, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, addComplaintData{
					Entity:                caseType + " 7000",
					Categories:            demoComplaintCategories,
					ComplainantCategories: demoComplainantCategories,
					Origins:               demoComplaintOrigins,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&case="+caseType, nil)
			w := httptest.NewRecorder()

			err := AddComplaint(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetAddComplaintBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":    "/?case=lpa",
		"no-case":  "/?id=123",
		"bad-case": "/?id=123&case=person",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := AddComplaint(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetAddComplaintWhenCaseErrors(t *testing.T) {
	client := &mockAddComplaintClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{}, errExample)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil).
		On("RefDataByCategory", mock.Anything, sirius.ComplaintCategory).
		Return(demoComplaintCategories, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := AddComplaint(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetAddComplaintWhenRefDataErrors(t *testing.T) {
	client := &mockAddComplaintClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{CaseType: "LPA", UID: "7000"}, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return([]sirius.RefDataItem{}, errExample)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil).
		On("RefDataByCategory", mock.Anything, sirius.ComplaintCategory).
		Return(demoComplaintCategories, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := AddComplaint(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetAddComplaintWhenTemplateErrors(t *testing.T) {
	client := &mockAddComplaintClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{CaseType: "LPA", UID: "7000"}, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil).
		On("RefDataByCategory", mock.Anything, sirius.ComplaintCategory).
		Return(demoComplaintCategories, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addComplaintData{
			Entity:                "LPA 7000",
			Categories:            demoComplaintCategories,
			ComplainantCategories: demoComplainantCategories,
			Origins:               demoComplaintOrigins,
		}).
		Return(errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := AddComplaint(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostAddComplaint(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			client := &mockAddComplaintClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(sirius.Case{CaseType: caseType, UID: "7000"}, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
				Return(demoComplainantCategories, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
				Return(demoComplaintOrigins, nil).
				On("RefDataByCategory", mock.Anything, sirius.ComplaintCategory).
				Return(demoComplaintCategories, nil)
			client.
				On("AddComplaint", mock.Anything, 123, sirius.CaseType(caseType), sirius.Complaint{
					Category:             "01",
					Description:          "This is a complaint",
					ReceivedDate:         sirius.DateString("2022-04-05"),
					Severity:             "Minor",
					InvestigatingOfficer: "Test Officer",
					SubCategory:          "07",
					Summary:              "In summary...",
				}).
				Return(nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, addComplaintData{
					Success:               true,
					Entity:                caseType + " 7000",
					Categories:            demoComplaintCategories,
					ComplainantCategories: demoComplainantCategories,
					Origins:               demoComplaintOrigins,
				}).
				Return(nil)

			form := url.Values{
				"category":             {"01"},
				"description":          {"This is a complaint"},
				"receivedDate":         {"2022-04-05"},
				"severity":             {"Minor"},
				"investigatingOfficer": {"Test Officer"},
				"subCategory":          {"07"},
				"summary":              {"In summary..."},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := AddComplaint(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostAddComplaintWhenAddComplaintValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	complaint := sirius.Complaint{Description: "This is a complaint"}

	client := &mockAddComplaintClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{CaseType: "LPA", UID: "7000"}, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil).
		On("RefDataByCategory", mock.Anything, sirius.ComplaintCategory).
		Return(demoComplaintCategories, nil)
	client.
		On("AddComplaint", mock.Anything, 123, sirius.CaseTypeLpa, complaint).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addComplaintData{
			Success:               false,
			Error:                 expectedError,
			Entity:                "LPA 7000",
			Complaint:             complaint,
			Categories:            demoComplaintCategories,
			ComplainantCategories: demoComplainantCategories,
			Origins:               demoComplaintOrigins,
		}).
		Return(nil)

	form := url.Values{
		"description": {"This is a complaint"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&case=lpa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AddComplaint(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostAddComplaintWhenAddComplaintOtherError(t *testing.T) {
	complaint := sirius.Complaint{Description: "This is a complaint"}

	client := &mockAddComplaintClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{CaseType: "LPA", UID: "7000"}, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil).
		On("RefDataByCategory", mock.Anything, sirius.ComplaintCategory).
		Return(demoComplaintCategories, nil)
	client.
		On("AddComplaint", mock.Anything, 123, sirius.CaseTypeLpa, complaint).
		Return(errExample)

	form := url.Values{
		"description": {"This is a complaint"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&case=lpa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AddComplaint(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}
