package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvailableStatuses(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		caseType         CaseType
		expectedResponse []string
		expectedError    func(int) error
	}{
		{
			name: "OK LPA",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request for the LPA's available statuses").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/lpas/800/available-statuses"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Body:    matchers.EachLike(matchers.String("Perfect"), 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			caseType:         CaseTypeLpa,
			expectedResponse: []string{"Perfect"},
		},
		{
			name: "OK EPA",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending EPA assigned").
					UponReceiving("A request for the EPA's available statuses").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/epas/800/available-statuses"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Body:    matchers.EachLike(matchers.String("Perfect"), 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			caseType:         CaseTypeEpa,
			expectedResponse: []string{"Perfect"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				caseitem, err := client.AvailableStatuses(Context{Context: context.Background()}, 800, tc.caseType)

				assert.Equal(t, tc.expectedResponse, caseitem)
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
