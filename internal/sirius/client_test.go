package sirius

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func newPact() *dsl.Pact {
	return &dsl.Pact{
		Consumer:          "sirius-lpa-frontend",
		Provider:          "sirius",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		LogDir:            "../../logs",
		PactDir:           "../../pacts",
	}
}

func TestStatusError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/some/url", nil)

	resp := &http.Response{
		StatusCode: http.StatusTeapot,
		Request:    req,
	}

	err := newStatusError(resp)

	assert.Equal(t, "POST /some/url returned 418", err.Error())
	assert.Equal(t, "unexpected response from Sirius", err.Title())
	assert.Equal(t, err, err.Data())
	assert.False(t, err.IsUnauthorized())
}

func TestStatusErrorIsUnauthorized(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/some/url", nil)

	resp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Request:    req,
	}

	err := newStatusError(resp)

	assert.Equal(t, "POST /some/url returned 401", err.Error())
	assert.Equal(t, "unexpected response from Sirius", err.Title())
	assert.Equal(t, err, err.Data())
	assert.True(t, err.IsUnauthorized())
}

func TestToFieldErrors(t *testing.T) {
	var unformattedErr flexibleFieldErrors
	err := json.Unmarshal([]byte(`{"riskAssessmentDate":["This field is required"],"reportApprovalDate":["This field is required"]}`), &unformattedErr)
	if err != nil {
		return
	}
	result, err := unformattedErr.toFieldErrors()
	formattedErr := FieldErrors{"riskAssessmentDate": {"": "This field is required"}, "reportApprovalDate": {"": "This field is required"}}

	assert.Equal(t, formattedErr, result)
	assert.Nil(t, err)
}

func TestToFieldErrorsThrowsError(t *testing.T) {
	var unformattedErr flexibleFieldErrors
	err := json.Unmarshal([]byte(`{"test":123}`), &unformattedErr)
	if err != nil {
		return
	}
	result, err := unformattedErr.toFieldErrors()

	assert.Equal(t, err, errors.New("could not parse field validation_errors"))
	assert.Nil(t, result)
}
