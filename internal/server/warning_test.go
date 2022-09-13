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

func (s *mockWarningClient) CreateWarning(ctx sirius.Context, personId int, warningType, warningNote string) error {
	args := s.Called(ctx, personId, warningType, warningNote)
	return args.Error(0)
}

func (s *mockWarningClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := s.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetWarning(t *testing.T) {
	warningTypes := []sirius.RefDataItem{
		{
			Handle: "Complaint Received",
			Label:  "Complaint Received",
		},
	}

	siriusClient := &mockWarningClient{}
	siriusClient.On("RefDataByCategory", mock.Anything, sirius.WarningTypeCategory).Return(warningTypes, nil)

	template := &mockTemplate{}
	template.On("Func", mock.Anything, warningData{
		Success:      false,
		XSRFToken:    "",
		WarningTypes: warningTypes,
	}).Return(nil)

	req, _ := http.NewRequest(http.MethodGet, "/?id=89", nil)

	w := httptest.NewRecorder()
	err := Warning(siriusClient, template.Func)(w, req)

	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestPostWarning(t *testing.T) {
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
	}).Return(nil)

	siriusClient.On("CreateWarning", mock.Anything, 89, "Complaint Recieved", "Some random warning notes").Return(nil)

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

	ve := sirius.ValidationError{
		Field: sirius.FieldErrors{
			"x": {"y": "z"},
		},
	}

	siriusClient.On("CreateWarning", mock.Anything, 89, "Complaint Received", "").Return(ve)

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

	e := errors.New("Some error")

	siriusClient.On("CreateWarning", mock.Anything, 89, "Complaint Recieved", "Some notes").Return(e)

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
