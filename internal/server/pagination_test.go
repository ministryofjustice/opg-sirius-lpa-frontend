package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPagination(t *testing.T) {
	testCases := map[string]struct {
		TotalItems, CurrentPage, TotalPages, PageSize int
		Start, End, PreviousPage, NextPage            int
		HasPrevious, HasNext                          bool
		Pages                                         []int
	}{
		"empty": {
			TotalItems:  0,
			CurrentPage: 1,
			TotalPages:  0,
			PageSize:    25,
			Start:       1,
			End:         0,
			HasPrevious: false,
			HasNext:     false,
			Pages:       []int{},
		},
		"one-item": {
			TotalItems:  1,
			CurrentPage: 1,
			TotalPages:  1,
			PageSize:    25,
			Start:       1,
			End:         1,
			HasPrevious: false,
			HasNext:     false,
			Pages:       []int{1},
		},
		"one-page": {
			TotalItems:  25,
			CurrentPage: 1,
			TotalPages:  1,
			PageSize:    25,
			Start:       1,
			End:         25,
			HasPrevious: false,
			HasNext:     false,
			Pages:       []int{1},
		},
		"many-pages": {
			TotalItems:   76,
			CurrentPage:  2,
			TotalPages:   4,
			PageSize:     25,
			Start:        26,
			End:          50,
			HasPrevious:  true,
			PreviousPage: 1,
			HasNext:      true,
			NextPage:     3,
			Pages:        []int{1, 2, 3, 4},
		},
		"first-of-many-pages": {
			TotalItems:  76,
			CurrentPage: 1,
			TotalPages:  4,
			PageSize:    25,
			Start:       1,
			End:         25,
			HasPrevious: false,
			HasNext:     true,
			NextPage:    2,
			Pages:       []int{1, 2, 3, 4},
		},
		"last-of-many-pages": {
			TotalItems:   76,
			CurrentPage:  4,
			TotalPages:   4,
			PageSize:     25,
			Start:        76,
			End:          76,
			HasPrevious:  true,
			PreviousPage: 3,
			HasNext:      false,
			Pages:        []int{1, 2, 3, 4},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			pagination := newPagination(&sirius.Pagination{
				TotalItems:  tc.TotalItems,
				CurrentPage: tc.CurrentPage,
				TotalPages:  tc.TotalPages,
				PageSize:    tc.PageSize,
			}, "term=bob", "")

			assert.Equal("?term=bob", pagination.SearchTerm)
			assert.Equal(tc.Start, pagination.Start())
			assert.Equal(tc.End, pagination.End())
			assert.Equal(tc.HasPrevious, pagination.HasPrevious())
			if tc.HasPrevious {
				assert.Equal(tc.PreviousPage, pagination.PreviousPage())
			}
			assert.Equal(tc.HasNext, pagination.HasNext())
			if tc.HasNext {
				assert.Equal(tc.NextPage, pagination.NextPage())
			}
			assert.Equal(tc.Pages, pagination.Pages())
		})
	}
}

func TestPaginationPagesWhenOverflow(t *testing.T) {
	testCases := map[int][]int{
		1:  []int{1, 2, -1, 10},
		2:  []int{1, 2, 3, -1, 10},
		3:  []int{1, 2, 3, 4, -1, 10},
		4:  []int{1, -1, 3, 4, 5, -1, 10},
		5:  []int{1, -1, 4, 5, 6, -1, 10},
		6:  []int{1, -1, 5, 6, 7, -1, 10},
		7:  []int{1, -1, 6, 7, 8, -1, 10},
		8:  []int{1, -1, 7, 8, 9, 10},
		9:  []int{1, -1, 8, 9, 10},
		10: []int{1, -1, 9, 10},
	}

	for current, pages := range testCases {
		t.Run(fmt.Sprintf("Page%d", current), func(t *testing.T) {
			pagination := newPagination(&sirius.Pagination{
				TotalItems:  250,
				CurrentPage: current,
				TotalPages:  10,
				PageSize:    25,
			}, "term=bob", "")

			assert.Equal(t, pages, pagination.Pages())
		})
	}
}

func TestPaginationWithFilters(t *testing.T) {
	assert := assert.New(t)

	pagination := newPagination(&sirius.Pagination{}, "term=bob", "person-type=Donor&person-type=Trust+Corporation")

	assert.Equal("?term=bob", pagination.SearchTerm)
	assert.Equal("&person-type=Donor&person-type=Trust+Corporation", pagination.Filters)
}
