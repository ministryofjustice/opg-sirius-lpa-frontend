package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWarningClient struct {
	mock.Mock
}

func (s *mockWarningClient) CreateWarning(ctx sirius.Context, personId int, warningType, warningNote string, caseIDs []int) error {
	args := s.Called(ctx, personId, warningType, warningNote, caseIDs)
	return args.Error(0)
}

func (s *mockWarningClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := s.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (s *mockWarningClient) CasesByDonor(ctx sirius.Context, personID int) ([]sirius.Case, error) {
	args := s.Called(ctx, personID)
	return args.Get(0).([]sirius.Case), args.Error(1)
}

func TestGetWarning(t *testing.T) {
	warningTypes := []sirius.RefDataItem{
		{
			Handle: "Complaint Received",
			Label:  "Complaint Received",
		},
	}

	cases := []sirius.Case{
		{
			CaseType: "PFA",
			UID:      "700700",
		},
	}

	siriusClient := &mockWarningClient{}
	siriusClient.On("RefDataByCategory", mock.Anything, sirius.WarningTypeCategory).Return(warningTypes, nil)
	siriusClient.On("CasesByDonor", mock.Anything, 89).Return(cases, nil)

	template := &mockTemplate{}
	template.On("Func", mock.Anything, warningData{
		Success:      false,
		XSRFToken:    "",
		WarningTypes: warningTypes,
		Cases:        cases,
	}).Return(nil)

	req, _ := http.NewRequest(http.MethodGet, "/?id=89", nil)

	w := httptest.NewRecorder()
	err := Warning(siriusClient, template.Func)(w, req)

	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestPostWarningWithOneCase(t *testing.T) {
	cases := []sirius.Case{
		{
			CaseType: "PFA",
			UID:      "7000-8888-0000",
		},
	}

	siriusClient := &mockWarningClient{}
	siriusClient.On("RefDataByCategory", mock.Anything, sirius.WarningTypeCategory).Return(
		[]sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		nil,
	)

	siriusClient.On("CasesByDonor", mock.Anything, 89).Return(cases, nil)

	template := &mockTemplate{}
	template.On("Func", mock.Anything, warningData{
		Success:   true,
		XSRFToken: "",
		WarningTypes: []sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		Cases: cases,
	}).Return(nil)

	siriusClient.On("CreateWarning", mock.Anything, 89, "Complaint Recieved", "Some random warning notes", []int{0}).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/?id=89", strings.NewReader(url.Values{
		"warningType": {"Complaint Recieved"},
		"warningText": {"Some random warning notes"},
	}.Encode()))

	req.Header.Add("content-type", formUrlEncoded)

	w := httptest.NewRecorder()
	err := Warning(siriusClient, template.Func)(w, req)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestPostWarningWithMultipleCases(t *testing.T) {
	cases := []sirius.Case{
		{
			ID:       1,
			CaseType: "PFA",
			UID:      "7000-1234-0000",
		},
		{
			ID:       2,
			CaseType: "HW",
			UID:      "7000-9876-0000",
		},
	}

	siriusClient := &mockWarningClient{}
	siriusClient.On("RefDataByCategory", mock.Anything, sirius.WarningTypeCategory).Return(
		[]sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		nil,
	)

	siriusClient.On("CasesByDonor", mock.Anything, 89).Return(cases, nil)

	template := &mockTemplate{}
	template.On("Func", mock.Anything, warningData{
		Success:   true,
		XSRFToken: "",
		WarningTypes: []sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		Cases: cases,
	}).Return(nil)

	siriusClient.On("CreateWarning", mock.Anything, 89, "Complaint Recieved", "Some random warning notes", []int{1, 2}).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/?id=89", strings.NewReader(url.Values{
		"case-id":     {"1", "2"},
		"warningType": {"Complaint Recieved"},
		"warningText": {"Some random warning notes"},
	}.Encode()))

	req.Header.Add("content-type", formUrlEncoded)

	w := httptest.NewRecorder()
	err := Warning(siriusClient, template.Func)(w, req)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestPostWarningWithNoCases(t *testing.T) {
	cases := []sirius.Case{}

	siriusClient := &mockWarningClient{}
	siriusClient.On("RefDataByCategory", mock.Anything, sirius.WarningTypeCategory).Return(
		[]sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		nil,
	)

	siriusClient.On("CasesByDonor", mock.Anything, 89).Return(cases, nil)

	template := &mockTemplate{}
	template.On("Func", mock.Anything, warningData{
		Success:   true,
		XSRFToken: "",
		WarningTypes: []sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		Cases: cases,
	}).Return(nil)

	siriusClient.On("CreateWarning", mock.Anything, 89, "Complaint Recieved", "Some random warning notes", []int{}).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/?id=89", strings.NewReader(url.Values{
		"warningType": {"Complaint Recieved"},
		"warningText": {"Some random warning notes"},
	}.Encode()))

	req.Header.Add("content-type", formUrlEncoded)

	w := httptest.NewRecorder()
	err := Warning(siriusClient, template.Func)(w, req)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestPostWarningValidationErrors(t *testing.T) {
	cases := []sirius.Case{}

	siriusClient := &mockWarningClient{}
	siriusClient.On("RefDataByCategory", mock.Anything, sirius.WarningTypeCategory).Return(
		[]sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		nil,
	)

	siriusClient.On("CasesByDonor", mock.Anything, 89).Return(cases, nil)

	ve := sirius.ValidationError{
		Field: sirius.FieldErrors{
			"x": {"y": "z"},
		},
	}

	siriusClient.On("CreateWarning", mock.Anything, 89, "Complaint Received", "", []int{}).Return(ve)

	template := &mockTemplate{}
	template.On("Func", mock.Anything, warningData{
		Success:   false,
		XSRFToken: "",
		WarningTypes: []sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		Error:       ve,
		WarningType: "Complaint Received",
		Cases:       cases,
	}).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/?id=89", strings.NewReader(url.Values{
		"warningType": {"Complaint Received"},
		"warningText": {""},
	}.Encode()))

	req.Header.Add("content-type", formUrlEncoded)

	w := httptest.NewRecorder()
	err := Warning(siriusClient, template.Func)(w, req)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestCreateWarningReturnsError(t *testing.T) {
	siriusClient := &mockWarningClient{}
	siriusClient.On("RefDataByCategory", mock.Anything, sirius.WarningTypeCategory).Return(
		[]sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		nil,
	)
	siriusClient.On("CasesByDonor", mock.Anything, 89).Return([]sirius.Case{}, nil)

	e := errors.New("Some error")

	siriusClient.On("CreateWarning", mock.Anything, 89, "Complaint Recieved", "Some notes", []int{}).Return(e)

	req, _ := http.NewRequest(http.MethodPost, "/?id=89", strings.NewReader(url.Values{
		"warningType": {"Complaint Recieved"},
		"warningText": {"Some notes"},
	}.Encode()))

	req.Header.Add("content-type", formUrlEncoded)

	w := httptest.NewRecorder()
	err := Warning(siriusClient, nil)(w, req)
	assert.Equal(t, e, err)
}

func TestGetWarningTypesFail(t *testing.T) {
	expectedErr := errors.New("Failed to get warning types")
	siriusClient := &mockWarningClient{}
	siriusClient.
		On("RefDataByCategory", mock.Anything, sirius.WarningTypeCategory).Return(nil, expectedErr)

	req, _ := http.NewRequest(http.MethodPost, "/?id=89", strings.NewReader(url.Values{
		"warningType": {"Complaint Recieved"},
		"warningText": {"Some notes"},
	}.Encode()))

	req.Header.Add("content-type", formUrlEncoded)

	w := httptest.NewRecorder()
	err := Warning(siriusClient, nil)(w, req)

	assert.Equal(t, expectedErr, err)
}
