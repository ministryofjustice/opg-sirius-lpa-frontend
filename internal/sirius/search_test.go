package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSearch(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name               string
		setup              func()
		expectedResponse   SearchResponse
		expectedPagination *Pagination
		expectedError      func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists to be referenced by name").
					UponReceiving("A search request for bob").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"term":        "bob",
							"personTypes": dsl.Like(AllPersonTypes),
							"size":        dsl.Like(PageLimit),
							"from":        dsl.Like(0),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"aggregations": dsl.Like(map[string]interface{}{
								"personType": map[string]int{
									"Donor": 1,
								},
							}),
							"total": dsl.Like(map[string]interface{}{
								"count": dsl.Like(1),
							}),
							"results": dsl.EachLike(map[string]interface{}{
								"id":           dsl.Like(36),
								"uId":          dsl.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
								"firstname":    dsl.Like("bob"),
								"surname":      dsl.Like("smith"),
								"addressLine1": dsl.Like("123 Somewhere Road"),
								"personType":   dsl.Like("Donor"),
								"cases": dsl.EachLike(map[string]interface{}{
									"id":          dsl.Like(23),
									"uId":         dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": dsl.Term("pfa", "hw|pfa"),
									"status":      dsl.Like("Perfect"),
									"caseType":    dsl.Like("LPA"),
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
						Firstname:    "bob",
						Surname:      "smith",
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				results, pagination, err := client.Search(Context{Context: context.Background()}, "bob", 1, AllPersonTypes)
				assert.Equal(t, tc.expectedResponse, results)
				assert.Equal(t, tc.expectedPagination, pagination)
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

func TestSearchDeletedCase(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name               string
		setup              func()
		expectedResponse   SearchResponse
		expectedPagination *Pagination
		expectedError      func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have deleted a case").
					UponReceiving("A search request for the deleted case uid").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"term":        "700000005555",
							"personTypes": AllPersonTypes,
							"size":        PageLimit,
							"from":        0,
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"total": dsl.Like(map[string]interface{}{
								"count": dsl.Like(0),
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				results, pagination, err := client.Search(Context{Context: context.Background()}, "700000005555", 1, AllPersonTypes)
				assert.Equal(t, tc.expectedResponse, results)
				assert.Equal(t, tc.expectedPagination, pagination)
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

func TestSearchTooShort(t *testing.T) {
	client := NewClient(http.DefaultClient, "")
	var res []Person

	results, pagination, err := client.Search(Context{Context: context.Background()}, "ad", 1, []string{})
	assert.Equal(t, results.Results, res)
	assert.Nil(t, pagination)
	assert.Equal(t, fmt.Errorf("Search term must be at least three characters"), err)
}
