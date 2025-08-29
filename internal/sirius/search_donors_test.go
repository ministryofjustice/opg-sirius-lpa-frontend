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

func TestSearchDonors(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []Person
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists to be referenced").
					UponReceiving("A search for donors").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/search/persons"),
						Body: map[string]interface{}{
							"term":        "7000-0000-0003",
							"personTypes": []string{"Donor"},
							"size":        PageLimit,
							"from":        0,
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"results": matchers.EachLike(map[string]interface{}{
								"id":           matchers.Like(47),
								"uId":          matchers.Like("7000-0000-0003"),
								"firstname":    matchers.Like("John"),
								"surname":      matchers.Like("Doe"),
								"addressLine1": matchers.Like("123 Somewhere Road"),
								"personType":   matchers.Like("Donor"),
								"cases": matchers.EachLike(map[string]interface{}{
									"id":          matchers.Like(23),
									"uId":         matchers.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": matchers.Term("pfa", "hw|pfa"),
									"status":      matchers.Like("Perfect"),
									"caseType":    matchers.Like("LPA"),
								}, 1),
							}, 1),
							"aggregations": matchers.Like(map[string]interface{}{
								"personType": map[string]int{
									"Donor": 1,
								},
							}),
							"total": matchers.Like(map[string]interface{}{
								"count": matchers.Like(1),
							}),
						}),
					})
			},
			expectedResponse: []Person{
				{
					ID:           47,
					UID:          "7000-0000-0003",
					Firstname:    "John",
					Surname:      "Doe",
					AddressLine1: "123 Somewhere Road",
					PersonType:   "Donor",
					Cases: []*Case{
						{
							ID:       23,
							UID:      "7000-8548-8461",
							CaseType: "LPA",
							SubType:  "pfa",
							Status:   shared.CaseStatusTypePerfect,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				donors, err := client.SearchDonors(Context{Context: context.Background()}, "7000-0000-0003")
				assert.Equal(t, tc.expectedResponse, donors)
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

func TestSearchDonorsTooShort(t *testing.T) {
	client := NewClient(http.DefaultClient, "")

	expectedErr := ValidationError{
		Detail: "Search term must be at least three characters",
		Field: FieldErrors{
			"term": {"reason": "Search term must be at least three characters"},
		},
	}

	donors, err := client.SearchDonors(Context{Context: context.Background()}, "ad")
	assert.Nil(t, donors)
	assert.Equal(t, expectedErr, err)
}
