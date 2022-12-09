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
					Given("An lpa exists to be referenced by donor name").
					UponReceiving("A search request for bob").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Body: dsl.Like(map[string]interface{}{
							"term":        dsl.String("bob"),
							"personTypes": dsl.Like(AllPersonTypes),
							"size":        dsl.Like(PageLimit),
							"from":        dsl.Like(0),
						}),
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
							"results": dsl.Like([]map[string]interface{}{
								{
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
									}, 1),
								},
							}),
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
								ID:      23,
								UID:     "7000-5382-4438",
								SubType: "pfa",
								Status:  "Perfect",
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

				results, pagination, err := client.Search(Context{Context: context.Background()}, "bob", 1, []string{})
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

func TestSearchDonorNotAssociatedWithCase(t *testing.T) {
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
					UponReceiving("A search request for daniel").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Body: dsl.Like(map[string]interface{}{
							"term":        dsl.String("daniel"),
							"personTypes": dsl.Like(AllPersonTypes),
							"size":        dsl.Like(PageLimit),
							"from":        dsl.Like(0),
						}),
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
								"id":           dsl.Like(33),
								"uId":          dsl.Term("7000-8554-1211", `\d{4}-\d{4}-\d{4}`),
								"firstname":    dsl.Like("Daniel"),
								"surname":      dsl.Like("Jones"),
								"addressLine1": dsl.Like("22 Test Road"),
								"personType":   dsl.Like("Donor"),
							}, 1),
						}),
					})
			},
			expectedResponse: SearchResponse{
				Results: []Person{
					{
						ID:           33,
						UID:          "7000-8554-1211",
						Firstname:    "Daniel",
						Surname:      "Jones",
						AddressLine1: "22 Test Road",
						PersonType:   "Donor",
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

				results, pagination, err := client.Search(Context{Context: context.Background()}, "daniel", 1, []string{})
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

func TestSearchFilteredResults(t *testing.T) {
	t.Parallel()

	pact := newIgnoredPact()
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
					Given("A donor and attorney exist with the same name").
					UponReceiving("A search request for an attorney named jack").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Body: dsl.Like(map[string]interface{}{
							"term":        dsl.String("jack"),
							"personTypes": dsl.Like([]string{"Attorney"}),
							"size":        dsl.Like(PageLimit),
							"from":        dsl.Like(0),
						}),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"aggregations": dsl.Like(map[string]interface{}{
								"personType": map[string]int{
									"Attorney": 1,
									"Donor":    1,
								},
							}),
							"total": dsl.Like(map[string]interface{}{
								"count": dsl.Like(1),
							}),
							"results": dsl.Like([]map[string]interface{}{
								{
									"id":           dsl.Like(36),
									"uId":          dsl.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
									"firstname":    dsl.Like("jack"),
									"surname":      dsl.Like("smith"),
									"addressLine1": dsl.Like("123 Somewhere Road"),
									"personType":   dsl.Like("Donor"),
									"cases": dsl.EachLike(map[string]interface{}{
										"id":          dsl.Like(23),
										"uId":         dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
										"caseSubtype": dsl.Term("pfa", "hw|pfa"),
										"status":      dsl.Like("Perfect"),
									}, 1),
								},
								{
									"id":           dsl.Like(37),
									"uId":          dsl.Term("7000-8512-5673", `\d{4}-\d{4}-\d{4}`),
									"firstname":    dsl.Like("jack"),
									"surname":      dsl.Like("taylor"),
									"addressLine1": dsl.Like("44 Unknown Road"),
									"personType":   dsl.Like("Attorney"),
									"cases": dsl.EachLike(map[string]interface{}{
										"id":          dsl.Like(24),
										"uId":         dsl.Term("7000-4382-1165", `\d{4}-\d{4}-\d{4}`),
										"caseSubtype": dsl.Term("hw", "hw|pfa"),
										"status":      dsl.Like("Pending"),
									}, 1),
								},
							}),
						}),
					})
			},
			expectedResponse: SearchResponse{
				Results: []Person{
					{
						ID:           36,
						UID:          "7000-8548-8461",
						Firstname:    "jack",
						Surname:      "smith",
						AddressLine1: "123 Somewhere Road",
						PersonType:   "Donor",
						Cases: []*Case{
							{
								ID:      23,
								UID:     "7000-5382-4438",
								SubType: "pfa",
								Status:  "Perfect",
							},
						},
					},
					{
						ID:           37,
						UID:          "7000-8512-5673",
						Firstname:    "jack",
						Surname:      "taylor",
						AddressLine1: "44 Unknown Road",
						PersonType:   "Attorney",
						Cases: []*Case{
							{
								ID:      24,
								UID:     "7000-4382-1165",
								SubType: "hw",
								Status:  "Pending",
							},
						},
					},
				},
				Aggregations: Aggregations{
					PersonType: map[string]int{
						"Donor":    1,
						"Attorney": 1,
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

				results, pagination, err := client.Search(Context{Context: context.Background()}, "jack", 1, []string{"Attorney"})
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
						Body: dsl.Like(map[string]interface{}{
							"term":        dsl.String("700000005555"),
							"personTypes": dsl.Like(AllPersonTypes),
							"size":        dsl.Like(PageLimit),
							"from":        dsl.Like(0),
						}),
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

				results, pagination, err := client.Search(Context{Context: context.Background()}, "700000005555", 1, []string{})
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
