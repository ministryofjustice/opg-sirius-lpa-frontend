package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name               string
		setup              func()
		expectedResponse   SearchResponse
		expectedPagination *Pagination
		expectedError      func(int) error
		searchTerm         string
	}{
		{
			name:       "OK",
			searchTerm: "bob",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists to be referenced by name").
					UponReceiving("A search request for a donor not related to a case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/search/persons"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"term":        "bob",
							"personTypes": AllPersonTypes,
							"size":        PageLimit,
							"from":        0,
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"aggregations": matchers.Like(map[string]interface{}{
								"personType": map[string]int{
									"Donor": 1,
								},
							}),
							"total": matchers.Like(map[string]interface{}{
								"count": matchers.Like(1),
							}),
							"results": matchers.EachLike(map[string]interface{}{
								"id":           matchers.Like(36),
								"uId":          matchers.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
								"firstname":    matchers.Like("Bob"),
								"surname":      matchers.Like("Smith"),
								"dob":          matchers.Like("17/03/1990"),
								"addressLine1": matchers.Like("123 Somewhere Road"),
								"personType":   matchers.Like("Donor"),
								"cases": matchers.EachLike(map[string]interface{}{
									"id":          matchers.Like(23),
									"uId":         matchers.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": matchers.Term("pfa", "hw|pfa"),
									"status":      matchers.Like("Perfect"),
									"caseType":    matchers.Like("LPA"),
								}, 1),
							}, 1),
						}),
					})
			},
			expectedResponse: SearchResponse{
				Results: []Person{
					{
						ID:           36,
						UID:          "7000-8548-8461",
						Firstname:    "Bob",
						Surname:      "Smith",
						DateOfBirth:  DateString("1990-03-17"),
						AddressLine1: "123 Somewhere Road",
						PersonType:   "Donor",
						Cases: []*Case{
							{
								ID:       23,
								UID:      "7000-5382-4438",
								CaseType: "LPA",
								SubType:  "pfa",
								Status:   "Perfect",
							},
						},
					},
				},
				Aggregations: Aggregations{
					PersonType: map[string]int{
						"Donor": 1,
					},
				},
				Total: SearchTotal{
					Count: 1,
				},
			},
			expectedPagination: &Pagination{
				TotalItems:  1,
				CurrentPage: 1,
				TotalPages:  1,
				PageSize:    PageLimit,
			},
		},
		{
			name:       "Deleted case",
			searchTerm: "700000005555",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have deleted a case").
					UponReceiving("A search request for the deleted case uid").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/search/persons"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"term":        "700000005555",
							"personTypes": AllPersonTypes,
							"size":        PageLimit,
							"from":        0,
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"total": matchers.Like(map[string]interface{}{
								"count": matchers.Like(0),
							}),
						}),
					})
			},
			expectedResponse: SearchResponse{
				Total: SearchTotal{
					Count: 0,
				},
			},
			expectedPagination: &Pagination{
				TotalItems:  0,
				CurrentPage: 1,
				TotalPages:  0,
				PageSize:    PageLimit,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				results, pagination, err := client.Search(Context{Context: context.Background()}, tc.searchTerm, 1, AllPersonTypes)
				assert.Equal(t, tc.expectedResponse, results)
				assert.Equal(t, tc.expectedPagination, pagination)
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
