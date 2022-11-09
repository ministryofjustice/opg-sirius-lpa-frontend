package sirius

import (
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

func TestFormatToValidationErr(t *testing.T) {
	unformattedErr := SiriusValidationError{Errors: SiriusFieldErrors{"riskAssessmentDate": []string{"This field is required"}, "reportApprovalDate": []string{"This field is required"}}}
	result := FormatToValidationError(unformattedErr)
	formattedErr := ValidationError{Field: FieldErrors{"riskAssessmentDate": {"": "This field is required"}, "reportApprovalDate": {"": "This field is required"}}}

	assert.Equal(t, formattedErr, result)
}
