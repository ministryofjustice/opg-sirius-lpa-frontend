package server

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SiriusClientMock struct {
	mock.Mock
}

type TemplateMock struct {
	mock.Mock
}

func (s *SiriusClientMock) CreateWarning(ctx sirius.Context, personId int, warningType, warningNote string) error {
	args := s.Called(ctx, personId, warningType, warningNote)
	return args.Error(0)
}

func (t *TemplateMock) ExecuteTemplate(w io.Writer, temp string, tempData interface{}) error {
	args := t.Called(w, temp, tempData)
	return args.Error(0)
}

func (s *SiriusClientMock) WarningTypes(ctx sirius.Context) ([]sirius.RefDataItem, error) {
	args := s.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetWarning(t *testing.T) {
	siriusClient := new(SiriusClientMock)
	siriusClient.On("WarningTypes", mock.Anything).Return(
		[]sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		nil,
	)

	template := new(TemplateMock)
	template.On("ExecuteTemplate", mock.Anything, "page", WarningData{
		WasWarningCreated: false,
		XSRFToken:         "",
		WarningTypes: []sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
	}).Return(nil)

	req, _ := http.NewRequest(http.MethodGet, "/?id=89", nil)

	w := httptest.NewRecorder()
	err := Warning(siriusClient, template)(w, req)

	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestPostWarning(t *testing.T) {
	siriusClient := new(SiriusClientMock)
	siriusClient.On("WarningTypes", mock.Anything).Return(
		[]sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		nil,
	)

	template := new(TemplateMock)
	template.On("ExecuteTemplate", mock.Anything, "page", WarningData{
		WasWarningCreated: true,
		XSRFToken:         "",
		WarningTypes: []sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
	}).Return(nil)

	siriusClient.On("CreateWarning", mock.Anything, 89, "Complaint Recieved", "Some random warning notes").Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/?id=89", strings.NewReader(url.Values{
		"warning-type":  {"Complaint Recieved"},
		"warning-notes": {"Some random warning notes"},
	}.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := Warning(siriusClient, template)(w, req)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestPostWarningValidationErrors(t *testing.T) {
	siriusClient := new(SiriusClientMock)
	siriusClient.On("WarningTypes", mock.Anything).Return(
		[]sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		nil,
	)

	v := sirius.ValidationError{}

	siriusClient.On("CreateWarning", mock.Anything, 89, "Complaint Recieved", "").Return(v)

	template := new(TemplateMock)
	template.On("ExecuteTemplate", mock.Anything, "page", WarningData{
		WasWarningCreated: false,
		XSRFToken:         "",
		WarningTypes: []sirius.RefDataItem{
			{
				Handle: "Complaint Received",
				Label:  "Complaint Received",
			},
		},
		ValidationErr: v,
	}).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/?id=89", strings.NewReader(url.Values{
		"warning-type":  {"Complaint Recieved"},
		"warning-notes": {""},
	}.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := Warning(siriusClient, template)(w, req)
	assert.Nil(t, err)
	result := w.Result()
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestCreateWarningReturnsError(t *testing.T) {
	siriusClient := new(SiriusClientMock)
	siriusClient.On("WarningTypes", mock.Anything).Return(
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
		"warning-type":  {"Complaint Recieved"},
		"warning-notes": {"Some notes"},
	}.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := Warning(siriusClient, nil)(w, req)
	assert.Equal(t, e, err)
}

func TestGetWarningTypesFail(t *testing.T) {
	expectedErr := errors.New("Failed to get warning types")
	siriusClient := new(SiriusClientMock)
	siriusClient.On("WarningTypes", mock.Anything).Return(nil, expectedErr)

	req, _ := http.NewRequest(http.MethodPost, "/?id=89", strings.NewReader(url.Values{
		"warning-type":  {"Complaint Recieved"},
		"warning-notes": {"Some notes"},
	}.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	err := Warning(siriusClient, nil)(w, req)

	assert.Equal(t, expectedErr, err)
}
