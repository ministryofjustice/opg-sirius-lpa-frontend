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

type mockEditComplaintClient struct {
	mock.Mock
}

func (m *mockEditComplaintClient) EditComplaint(ctx sirius.Context, id int, complaint sirius.Complaint) error {
	args := m.Called(ctx, id, complaint)
	return args.Error(0)
}

func (m *mockEditComplaintClient) Complaint(ctx sirius.Context, id int) (sirius.Complaint, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Complaint), args.Error(1)
}

func (m *mockEditComplaintClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockEditComplaintClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetEditComplaint(t *testing.T) {
	complaint := sirius.Complaint{
		Category: "01",
	}

	client := &mockEditComplaintClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil)
	client.
		On("Complaint", mock.Anything, 123).
		Return(complaint, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editComplaintData{
			Complaint:             complaint,
			Categories:            complaintCategories,
			ComplainantCategories: demoComplainantCategories,
			Origins:               demoComplaintOrigins,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditComplaint(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetEditComplaintNoID(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	err := EditComplaint(nil, nil)(w, r)

	assert.NotNil(t, err)
}

func TestGetEditComplaintWhenRefDataErrors(t *testing.T) {
	client := &mockEditComplaintClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditComplaint(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetEditComplaintWhenComplaintErrors(t *testing.T) {
	client := &mockEditComplaintClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil)
	client.
		On("Complaint", mock.Anything, 123).
		Return(sirius.Complaint{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditComplaint(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetEditComplaintWhenTemplateErrors(t *testing.T) {
	client := &mockEditComplaintClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil)
	client.
		On("Complaint", mock.Anything, 123).
		Return(sirius.Complaint{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editComplaintData{
			Categories:            complaintCategories,
			ComplainantCategories: demoComplainantCategories,
			Origins:               demoComplaintOrigins,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditComplaint(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditComplaint(t *testing.T) {
	complaint := sirius.Complaint{
		Category:       "01",
		Description:    "This is a complaint",
		ReceivedDate:   sirius.DateString("2022-04-05"),
		Severity:       "Minor",
		SubCategory:    "07",
		Summary:        "In summary...",
		Resolution:     "complaint upheld",
		ResolutionInfo: "This is what we did",
		ResolutionDate: sirius.DateString("2022-05-06"),
	}

	client := &mockEditComplaintClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil)
	client.
		On("Complaint", mock.Anything, 123).
		Return(complaint, nil)
	client.
		On("EditComplaint", mock.Anything, 123, complaint).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editComplaintData{
			Success:               true,
			Complaint:             complaint,
			Categories:            complaintCategories,
			ComplainantCategories: demoComplainantCategories,
			Origins:               demoComplaintOrigins,
		}).
		Return(nil)

	form := url.Values{
		"category":       {"01"},
		"description":    {"This is a complaint"},
		"receivedDate":   {"2022-04-05"},
		"severity":       {"Minor"},
		"subCategory":    {"07"},
		"summary":        {"In summary..."},
		"resolution":     {"complaint upheld"},
		"resolutionInfo": {"This is what we did"},
		"resolutionDate": {"2022-05-06"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditComplaint(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditComplaintWhenEditComplaintValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	complaint := sirius.Complaint{Description: "This is a complaint"}

	client := &mockEditComplaintClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil)
	client.
		On("Complaint", mock.Anything, 123).
		Return(complaint, nil)
	client.
		On("EditComplaint", mock.Anything, 123, complaint).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editComplaintData{
			Success:               false,
			Error:                 expectedError,
			Complaint:             complaint,
			Categories:            complaintCategories,
			ComplainantCategories: demoComplainantCategories,
			Origins:               demoComplaintOrigins,
		}).
		Return(nil)

	form := url.Values{
		"description": {"This is a complaint"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditComplaint(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditComplaintWhenEditComplaintOtherError(t *testing.T) {
	complaint := sirius.Complaint{Description: "This is a complaint"}

	client := &mockEditComplaintClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplainantCategory).
		Return(demoComplainantCategories, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.ComplaintOrigin).
		Return(demoComplaintOrigins, nil)
	client.
		On("EditComplaint", mock.Anything, 123, complaint).
		Return(expectedError)

	form := url.Values{
		"description": {"This is a complaint"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditComplaint(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
