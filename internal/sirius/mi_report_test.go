package sirius

import (
	"context"
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
			expectedResult: &MiReportResponse{
				ResultCount:       10,
				ReportType:        "epasReceived",
				ReportDescription: "Number of EPAs received",
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
				result, err := client.MiReport(Context{Context: context.Background()}, form)

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
