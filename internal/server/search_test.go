package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockSearchClient struct {
	mock.Mock
}

func (m *mockSearchClient) Search(ctx sirius.Context, term string, page int, personTypeFilters []string) (sirius.SearchResponse, *sirius.Pagination, error) {
	args := m.Called(ctx, term, page, personTypeFilters)
	if v, ok := args.Get(1).(*sirius.Pagination); ok {
		return args.Get(0).(sirius.SearchResponse), v, args.Error(2)
	}
	return args.Get(0).(sirius.SearchResponse), nil, args.Error(2)
}

func TestGetSearch(t *testing.T) {
	persons := []sirius.Person{
		{ID: 1, Firstname: "John"},
		{ID: 2, Firstname: "Jane"},
		{ID: 3, Firstname: "Bob"},
	}

	expectedResponse := sirius.SearchResponse{
		Results: persons,
		Aggregations: sirius.Aggregations{
			PersonType: map[string]int{
				"Donor": 3,
			},
		},
		Total: sirius.SearchTotal{
			Count: 3,
		},
	}

	expectedPagination := &sirius.Pagination{TotalItems: 3}
	var filters []string

	client := &mockSearchClient{}
	client.
		On("Search", mock.Anything, "bob", 1, filters).
		Return(expectedResponse, expectedPagination, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, searchData{
			Results:      persons,
			Total:        expectedResponse.Total.Count,
			Aggregations: expectedResponse.Aggregations,
			Filters:      searchFilters{},
			SearchTerm:   "bob",
			Pagination:   newPagination(expectedPagination, "term=bob", ""),
		}).
		Return(nil)

	req, _ := http.NewRequest(http.MethodGet, "/search?term=bob", nil)
	w := httptest.NewRecorder()

	err := Search(client, template.Func)(w, req)
	assert.Nil(t, err)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSearchFiltered(t *testing.T) {
	persons := []sirius.Person{
		{ID: 1, Firstname: "John"},
		{ID: 2, Firstname: "Jane"},
		{ID: 3, Firstname: "Bob"},
	}

	expectedResponse := sirius.SearchResponse{
		Results: persons,
		Aggregations: sirius.Aggregations{
			PersonType: map[string]int{
				"Donor":    2,
				"Attorney": 1,
			},
		},
		Total: sirius.SearchTotal{
			Count: 3,
		},
	}

	expectedPagination := &sirius.Pagination{TotalItems: 3}
	filters := []string{"Donor", "Attorney"}

	client := &mockSearchClient{}
	client.
		On("Search", mock.Anything, "bob", 1, filters).
		Return(expectedResponse, expectedPagination, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, searchData{
			Results:      persons,
			Total:        expectedResponse.Total.Count,
			Aggregations: expectedResponse.Aggregations,
			Filters:      searchFilters{Set: true, PersonType: filters},
			SearchTerm:   "bob",
			Pagination: newPagination(
				expectedPagination,
				"term=bob",
				"person-type=Donor&person-type=Attorney",
			),
		}).
		Return(nil)

	req, _ := http.NewRequest(http.MethodGet, "/search?term=bob&person-type=Donor&person-type=Attorney", nil)
	w := httptest.NewRecorder()

	err := Search(client, template.Func)(w, req)
	assert.Nil(t, err)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSearchPaginationCalculations(t *testing.T) {
	persons := []sirius.Person{{ID: 1, Firstname: "John"}}

	expectedResponse := sirius.SearchResponse{
		Results: persons,
		Aggregations: sirius.Aggregations{
			PersonType: map[string]int{
				"Donor": 80,
			},
		},
		Total: sirius.SearchTotal{
			Count: 80,
		},
	}

	expectedPagination := &sirius.Pagination{
		TotalItems:  80,
		CurrentPage: 2,
		TotalPages:  4,
		PageSize:    sirius.PageLimit,
	}

	var filters []string

	client := &mockSearchClient{}
	client.
		On("Search", mock.Anything, "bob", 2, filters).
		Return(expectedResponse, expectedPagination, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, searchData{
			Results:      persons,
			Total:        expectedResponse.Total.Count,
			Aggregations: expectedResponse.Aggregations,
			Filters:      searchFilters{},
			SearchTerm:   "bob",
			Pagination:   newPagination(expectedPagination, "term=bob", ""),
		}).
		Return(nil)

	req, _ := http.NewRequest(http.MethodGet, "/search?term=bob&page=2", nil)
	w := httptest.NewRecorder()

	err := Search(client, template.Func)(w, req)
	assert.Nil(t, err)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSearchBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-search-term": "/search?term=",
		"bad-query":      "/search?abc=hello",
	}

	for name, urlParams := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, urlParams, nil)
			w := httptest.NewRecorder()

			err := Search(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetSearchErrors(t *testing.T) {
	var filters []string

	client := &mockSearchClient{}
	client.
		On("Search", mock.Anything, "bob", 1, filters).
		Return(sirius.SearchResponse{}, &sirius.Pagination{}, expectedError)

	req, _ := http.NewRequest(http.MethodGet, "/search?term=bob", nil)
	w := httptest.NewRecorder()

	err := Search(client, nil)(w, req)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetSearchTemplateErrors(t *testing.T) {
	persons := []sirius.Person{{ID: 1, Firstname: "John"}}

	expectedResponse := sirius.SearchResponse{
		Results: persons,
		Aggregations: sirius.Aggregations{
			PersonType: map[string]int{
				"Donor": 3,
			},
		},
		Total: sirius.SearchTotal{
			Count: 3,
		},
	}

	expectedPagination := &sirius.Pagination{TotalItems: 3}
	var filters []string

	client := &mockSearchClient{}
	client.
		On("Search", mock.Anything, "bob", 1, filters).
		Return(expectedResponse, expectedPagination, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, searchData{
			Results:      persons,
			Total:        expectedResponse.Total.Count,
			Aggregations: expectedResponse.Aggregations,
			Filters:      searchFilters{},
			SearchTerm:   "bob",
			Pagination:   newPagination(expectedPagination, "term=bob", ""),
		}).
		Return(expectedError)

	req, _ := http.NewRequest(http.MethodGet, "/search?term=bob", nil)
	w := httptest.NewRecorder()

	err := Search(client, template.Func)(w, req)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
