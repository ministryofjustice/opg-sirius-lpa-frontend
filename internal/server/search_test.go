package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/mock"
)

type mockSearchClient struct {
	mock.Mock
}

func (m *mockSearchClient) Search(ctx sirius.Context, term string, page int, personTypeFilters []string) (sirius.SearchResponse, *sirius.Pagination, error) {
	args := m.Called(ctx, term, page, personTypeFilters)
	return args.Get(0).(sirius.SearchResponse), args.Get(1).(*sirius.Pagination), args.Error(1)
}

//func TestGetSearch(t *testing.T) {
//	persons := []sirius.Person{
//		{ID: 1, Firstname: "John"},
//		{ID: 2, Firstname: "Jane"},
//		{ID: 3, Firstname: "Bob"},
//	}
//
//	expectedResponse := sirius.SearchResponse{
//		Results:            persons,
//		RawAggregations:    nil,
//		AggregationsObject: sirius.AggregationsResult{},
//		Total: sirius.SearchTotal{
//			Count: 2,
//		},
//	}
//
//	p := &sirius.Pagination{TotalItems: 2}
//	var personTypeFilters []string
//
//	client := &mockSearchClient{}
//	client.
//		On("Search", mock.Anything, "bob", 1, personTypeFilters).
//		Return(expectedResponse, p, nil)
//
//	template := &mockTemplate{}
//	template.
//		On("Func", mock.Anything, searchData{
//			Results:    persons,
//			Total:      expectedResponse.Total.Count,
//			Filters:    searchFilters{},
//			SearchTerm: "bob",
//			Pagination: newPaginationWithQuery(p, "term=bob", ""),
//		}).
//		Return(nil)
//
//	req, _ := http.NewRequest(http.MethodGet, "/?term=bob", nil)
//	w := httptest.NewRecorder()
//
//	err := Search(client, template.Func)(w, req)
//	resp := w.Result()
//
//	assert.Nil(t, err)
//	assert.Equal(t, http.StatusOK, resp.StatusCode)
//
//	//var persons []sirius.Person
//	//_ = json.NewDecoder(resp.Body).Decode(&persons)
//	//
//	//assert.Equal(t, expectedResponse, persons)
//}
