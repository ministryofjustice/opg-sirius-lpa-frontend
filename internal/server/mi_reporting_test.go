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

type mockMiReportingClient struct {
	mock.Mock
}

func (m *mockMiReportingClient) MiConfig(ctx sirius.Context) (map[string]sirius.MiConfigProperty, error) {
	args := m.Called(ctx)

	if v, ok := args.Get(0).(map[string]sirius.MiConfigProperty); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockMiReportingClient) MiReport(ctx sirius.Context, form url.Values) (*sirius.MiReportResponse, error) {
	args := m.Called(ctx, form)
	if v, ok := args.Get(0).(*sirius.MiReportResponse); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetMiReporting(t *testing.T) {
	reportTypes := []sirius.MiConfigEnum{
		{Name: "epaReport", Description: "EPA Report"},
		{Name: "lpaReport", Description: "LPA Report"},
	}

	client := &mockMiReportingClient{}
	client.
		On("MiConfig", mock.Anything).
		Return(map[string]sirius.MiConfigProperty{
			"reportType": {
				Enum: reportTypes,
			},
		}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, miReportingData{
			ReportTypes: reportTypes,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	err := MiReporting(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetMiReportingWhenReportTypeSelected(t *testing.T) {
	reportTypes := []sirius.MiConfigEnum{
		{Name: "epaReport", Description: "EPA Report"},
		{Name: "lpaReport", Description: "LPA Report"},
	}

	dateControl := sirius.MiConfigProperty{
		Description: "date",
		DependsOn: sirius.MiConfigDependsOn{
			ReportType: []sirius.MiConfigReportType{{Name: "lpaReport"}},
		},
	}

	client := &mockMiReportingClient{}
	client.
		On("MiConfig", mock.Anything).
		Return(map[string]sirius.MiConfigProperty{
			"reportType": {
				Enum: reportTypes,
			},
			"endDate":   dateControl,
			"startDate": dateControl,
			"otherThing": {
				Description: "date",
			},
		}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, miReportingData{
			ReportTypes: reportTypes,
			ReportType:  "lpaReport",
			ReportName:  "LPA Report",
			Controls: []namedControl{
				{
					Name:       "startDate",
					Label:      "From",
					Properties: dateControl,
				},
				{
					Name:       "endDate",
					Label:      "To",
					Properties: dateControl,
				},
			},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?reportType=lpaReport", nil)
	w := httptest.NewRecorder()

	err := MiReporting(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetMiReportingWhenConfigErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockMiReportingClient{}
	client.
		On("MiConfig", mock.Anything).
		Return(nil, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	err := MiReporting(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetMiReportingWhenTemplateErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	reportTypes := []sirius.MiConfigEnum{
		{Name: "epaReport", Description: "EPA Report"},
		{Name: "lpaReport", Description: "LPA Report"},
	}

	client := &mockMiReportingClient{}
	client.
		On("MiConfig", mock.Anything).
		Return(map[string]sirius.MiConfigProperty{
			"reportType": {
				Enum: reportTypes,
			},
		}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, miReportingData{
			ReportTypes: reportTypes,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	err := MiReporting(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostMiReporting(t *testing.T) {
	client := &mockMiReportingClient{}
	client.
		On("MiReport", mock.Anything, url.Values{
			"reportType":        {"lpaReport"},
			"applicationType[]": {"HW", "PFA"},
			"startDate":         {"02/01/2020"},
		}).
		Return(&sirius.MiReportResponse{
			ResultCount:       10,
			ReportType:        "lpaReport",
			ReportDescription: "LPA Report",
		}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, miReportingData{
			ResultCount: 10,
			ReportType:  "lpaReport",
			ReportName:  "LPA Report",
			Download:    "/api/reporting/export?OPG-Bypass-Membrane=1&applicationType%5B%5D=HW&applicationType%5B%5D=PFA&reportType=lpaReport&startDate=02%2F01%2F2020",
		}).
		Return(nil)

	form := url.Values{
		"reportType":      {"lpaReport"},
		"applicationType": {"HW", "PFA"},
		"startDate":       {"2020-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := MiReporting(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostMiReportingWhenError(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockMiReportingClient{}
	client.
		On("MiReport", mock.Anything, url.Values{}).
		Return(nil, expectedError)

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := MiReporting(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
