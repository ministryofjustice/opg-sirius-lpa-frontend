package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"term":        "bob",
							"personTypes": AllPersonTypes,
							"size":        PageLimit,
							"from":        0,
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
								"firstname":    dsl.Like("Bob"),
								"surname":      dsl.Like("Smith"),
								"dob":          dsl.Like("17/03/1990"),
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

				results, pagination, err := client.Search(Context{Context: context.Background()}, tc.searchTerm, 1, AllPersonTypes)
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

func TestSearchDifferentResults(t *testing.T) {
	pact := newIgnoredPact()
	defer pact.Teardown()

	testCases := []struct {
		name               string
		setup              func()
		expectedResponse   SearchResponse
		expectedPagination *Pagination
		expectedError      func(int) error
		searchTerm         string
		filtersApplied     []string
	}{
		{
			name:           "A donor with more than one case",
			searchTerm:     "harry",
			filtersApplied: AllPersonTypes,
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists").
					UponReceiving("A search request for a donor who has more than one case").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"term":        "harry",
							"personTypes": AllPersonTypes,
							"size":        PageLimit,
							"from":        0,
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
								"firstname":    dsl.Like("Harry"),
								"surname":      dsl.Like("Jones"),
								"addressLine1": dsl.Like("123 Somewhere Road"),
								"personType":   dsl.Like("Donor"),
								"cases": []map[string]interface{}{
									{
										"id":          dsl.Like(23),
										"uId":         dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
										"caseSubtype": dsl.Term("pfa", "hw|pfa"),
										"status":      dsl.Like("Perfect"),
										"caseType":    dsl.Like("LPA"),
									},
									{
										"id":          dsl.Like(24),
										"uId":         dsl.Term("7000-5382-8764", `\d{4}-\d{4}-\d{4}`),
										"caseSubtype": dsl.Term("hw", "hw|pfa"),
										"status":      dsl.Like("Pending"),
										"caseType":    dsl.Like("LPA"),
									},
								},
							}, 1),
						}),
					})
			},
			expectedResponse: SearchResponse{
				Results: []Person{
					{
						ID:           36,
						UID:          "7000-8548-8461",
						Firstname:    "Harry",
						Surname:      "Jones",
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
							{
								ID:       24,
								UID:      "7000-5382-8764",
								CaseType: "LPA",
								SubType:  "hw",
								Status:   "Pending",
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
			name:           "No search results",
			searchTerm:     "abcdefg",
			filtersApplied: AllPersonTypes,
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists").
					UponReceiving("A search request for a person who doesn't exist").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"term":        "abcdefg",
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
		{
			name:           "No search results with filters",
			searchTerm:     "abcdefg",
			filtersApplied: []string{"Attorney", "Trust Corporation"},
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists").
					UponReceiving("A search request for a person who doesn't exist with filters applied").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"term":        "abcdefg",
							"personTypes": []string{"Attorney", "Trust Corporation"},
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
		{
			name:           "Searching a deputy",
			searchTerm:     "fred",
			filtersApplied: AllPersonTypes,
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A search request for a deputy").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"term":        "fred",
							"personTypes": AllPersonTypes,
							"size":        PageLimit,
							"from":        0,
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"aggregations": dsl.Like(map[string]interface{}{
								"personType": map[string]int{
									"Deputy": 1,
								},
							}),
							"total": dsl.Like(map[string]interface{}{
								"count": dsl.Like(1),
							}),
							"results": dsl.EachLike(map[string]interface{}{
								"id":           dsl.Like(65),
								"uId":          dsl.Term("7000-6509-8813", `\d{4}-\d{4}-\d{4}`),
								"firstname":    dsl.Like("Fred"),
								"surname":      dsl.Like("Jones"),
								"addressLine1": dsl.Like("100 Random Road"),
								"personType":   dsl.Like("Deputy"),
								"cases": []map[string]interface{}{
									{
										"id":          dsl.Like(48),
										"uId":         dsl.Term("7000-5113-1871", `\d{4}-\d{4}-\d{4}`),
										"caseSubtype": dsl.Term("hw", "hw|pfa"),
										"caseType":    dsl.Like("ORDER"),
									},
								},
							}, 1),
						}),
					})
			},
			expectedResponse: SearchResponse{
				Results: []Person{
					{
						ID:           65,
						UID:          "7000-6509-8813",
						Firstname:    "Fred",
						Surname:      "Jones",
						AddressLine1: "100 Random Road",
						PersonType:   "Deputy",
						Cases: []*Case{
							{
								ID:       48,
								UID:      "7000-5113-1871",
								CaseType: "ORDER",
								SubType:  "hw",
							},
						},
					},
				},
				Aggregations: Aggregations{
					PersonType: map[string]int{
						"Deputy": 1,
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

				results, pagination, err := client.Search(Context{Context: context.Background()}, tc.searchTerm, 1, tc.filtersApplied)
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
