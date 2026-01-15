package sirius

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestMiReport(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		setup          func()
		expectedResult *MiReportResponse
		expectedError  func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a System Admin").
					UponReceiving("A request for an MI report").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/reporting/applications"),
						Query: matchers.MapMatcher{
							"reportType": matchers.String("epasReceived"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"data": matchers.Like(map[string]interface{}{
								"result_count":       matchers.Like(10),
								"report_type":        matchers.String("epasReceived"),
								"report_description": matchers.String("Number of EPAs received"),
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				form := url.Values{
					"reportType": {"epasReceived"},
				}
				result, err := client.MiReport(Context{Context: context.Background()}, form)

				assert.Equal(t, tc.expectedResult, result)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}
