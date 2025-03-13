package sirius

import (
	"math"
)

const PageLimit = 25

type searchRequest struct {
	Term        string   `json:"term"`
	PersonTypes []string `json:"personTypes"`
	Limit       int      `json:"size"`
	From        int      `json:"from"`
}

type Aggregations struct {
	PersonType map[string]int `json:"personType"`
}

type SearchTotal struct {
	Count int `json:"count"`
}

type SearchResponse struct {
	Results      []Person     `json:"results,omitempty"`
	Aggregations Aggregations `json:"aggregations,omitempty"`
	Total        SearchTotal  `json:"total"`
}

var AllPersonTypes = []string{
	"Donor",
	"Client",
	"Attorney",
	"Deputy",
	"Replacement Attorney",
	"Trust Corporation",
	"Notified Person",
	"Certificate Provider",
	"Correspondent",
}

func (c *Client) Search(ctx Context, term string, page int, personTypeFilters []string) (SearchResponse, *Pagination, error) {
	var v SearchResponse
	if len(term) < 3 {
		err := ValidationError{
			Detail: "Search term must be at least three characters",
			Field: FieldErrors{
				"term": {"reason": "Search term must be at least three characters"},
			},
		}

		return v, nil, err
	}

	if len(personTypeFilters) == 0 {
		personTypeFilters = AllPersonTypes
	}

	data := searchRequest{Term: term, PersonTypes: personTypeFilters, Limit: PageLimit, From: PageLimit * (page - 1)}

	err := c.post(ctx, "/lpa-api/v1/search/persons", data, &v)
	if err != nil {
		return v, nil, err
	}

	return v, &Pagination{
		TotalItems:  v.Total.Count,
		CurrentPage: page,
		TotalPages:  int(math.Ceil(float64(v.Total.Count) / float64(PageLimit))),
		PageSize:    PageLimit,
	}, nil
}
