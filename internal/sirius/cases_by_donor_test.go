package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCasesByDonor(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/persons/400/cases"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"cases": []map[string]interface{}{
								{
									"id":          dsl.Like(405),
									"uId":         dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": dsl.Term("pfa", "hw|pfa"),
								},
								{
									"id":          dsl.Like(406),
									"uId":         dsl.Term("7000-5382-8764", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": dsl.Term("hw", "hw|pfa"),
								},
							},
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []Case{
				{
					ID:       405,
					UID:      "7000-5382-4438",
					CaseType: "LPA",
					SubType:  "pfa",
				},
				{
					ID:       406,
					UID:      "7000-5382-8764",
					CaseType: "LPA",
					SubType:  "hw",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				caseitem, err := client.CasesByDonor(Context{Context: context.Background()}, 400)

				assert.Equal(t, tc.expectedResponse, caseitem)
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
