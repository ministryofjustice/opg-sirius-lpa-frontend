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

func TestCasesByDonor(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []Case
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists with more than 1 case").
					UponReceiving("A request for the donor's cases").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/400/cases"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"cases": []map[string]interface{}{
								{
									"id":          matchers.Like(405),
									"uId":         matchers.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": matchers.Term("pfa", "hw|pfa"),
								},
								{
									"id":          matchers.Like(406),
									"uId":         matchers.Term("7000-5382-8764", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": matchers.Term("hw", "hw|pfa"),
								},
							},
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []Case{
				{
					ID:      405,
					UID:     "7000-5382-4438",
					SubType: "pfa",
				},
				{
					ID:      406,
					UID:     "7000-5382-8764",
					SubType: "hw",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				caseitem, err := client.CasesByDonor(Context{Context: context.Background()}, 400)

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
