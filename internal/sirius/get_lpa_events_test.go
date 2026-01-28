package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestGetLpaEvents(t *testing.T) {
	t.Parallel()

	pactResponse := consumer.Response{
		Status: http.StatusOK,
		Body: matchers.Like(map[string]interface{}{
			"limit": matchers.Like(999),
			"total": matchers.Like(1),
			"pages": matchers.Like(map[string]interface{}{
				"current": matchers.Like(1),
				"total":   matchers.Like(1),
			}),
			"events": matchers.EachLike(map[string]interface{}{
				"id":         matchers.Like(4056),
				"sourceType": matchers.Like("Lpa"),
				"type":       matchers.Like("INS"),
				"createdOn":  matchers.Like("2026-01-16T04:10:55+00:00"),
				"hash":       matchers.Like("JIG"),
			}, 1),
		}),
		Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
	}

	expectedResponse := LpaEventsResponse{
		Limit: 999,
		Pages: Pages{
			Current: 1,
			Total:   1,
		},
		Total: 1,
		Events: []LpaEvent{
			{
				ID:         4056,
				Type:       "INS",
				SourceType: shared.LpaEventSourceTypeLpa,
				CreatedOn:  "2026-01-16T04:10:55+00:00",
				Hash:       "JIG",
			},
		},
	}

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		caseIDs          []string
		setup            func()
		expectedError    func(int) error
		expectedResponse LpaEventsResponse
	}{
		{
			name:    "Can get events on a person",
			caseIDs: []string(nil),
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists with multiple cases and event history").
					UponReceiving("A request for the events on a person").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/189/events"),
						Query: matchers.MapMatcher{
							"sort":  matchers.String("id:desc"),
							"limit": matchers.Like("999"),
						},
					}).
					WithCompleteResponse(pactResponse)
			},
			expectedResponse: expectedResponse,
		},
		{
			name:    "Can get events on a single case",
			caseIDs: []string{"111111"},
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists with multiple cases and event history").
					UponReceiving("A request for the events on a case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/189/events"),
						Query: matchers.MapMatcher{
							"filter": matchers.String("case:111111"),
							"sort":   matchers.String("id:desc"),
							"limit":  matchers.Like("999"),
						},
					}).
					WithCompleteResponse(pactResponse)
			},
			expectedResponse: expectedResponse,
		},
		{
			name:    "Can get events on a multiple case",
			caseIDs: []string{"111111", "222222"},
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists with multiple cases and event history").
					UponReceiving("A request for the events on a case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/189/events"),
						Query: matchers.MapMatcher{
							"filter": matchers.String("case:111111,case:222222"),
							"sort":   matchers.String("id:desc"),
							"limit":  matchers.Like("999"),
						},
					}).
					WithCompleteResponse(pactResponse)
			},
			expectedResponse: expectedResponse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				documents, err := client.GetEvents(Context{Context: context.Background()}, "189", tc.caseIDs, []string{}, "desc")

				assert.Equal(t, tc.expectedResponse, documents)
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
