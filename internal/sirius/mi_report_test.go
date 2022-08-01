package sirius

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestMiReport(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name           string
		setup          func()
		cookies        []*http.Cookie
		expectedResult *MiReportResponse
		expectedError  func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request for an MI report").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/reporting/applications"),
						Query: dsl.MapMatcher{
							"reportType": dsl.String("epasReceived"),
						},
						Headers: dsl.MapMatcher{
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"data": dsl.Like(map[string]interface{}{
								"result_count":       dsl.Like(10),
								"report_type":        dsl.String("epasReceived"),
								"report_description": dsl.String("Number of EPAs received"),
							}),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResult: &MiReportResponse{
				ResultCount:       10,
				ReportType:        "epasReceived",
				ReportDescription: "Number of EPAs received",
			},
		},
		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request for an MI report without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/reporting/applications"),
						Query: dsl.MapMatcher{
							"reportType": dsl.String("epasReceived"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:   http.StatusUnauthorized,
					URL:    fmt.Sprintf("http://localhost:%d/api/reporting/applications?reportType=epasReceived", port),
					Method: http.MethodGet,
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				form := url.Values{
					"reportType": {"epasReceived"},
				}
				result, err := client.MiReport(getContext(tc.cookies), form)

				assert.Equal(t, tc.expectedResult, result)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
