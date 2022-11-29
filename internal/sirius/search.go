package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

const PageLimit = 10

type searchRequest struct {
	Term        string   `json:"term"`
	PersonTypes []string `json:"personTypes"`
	Limit       int      `json:"size"`
	From        int      `json:"from"`
}

type AggregationsResult struct {
	PersonType map[string]interface{} `json:"personType"`
}

type SearchTotal struct {
	Count int `json:"count"`
}

type SearchResponse struct {
	Results            []Person           `json:"results"`
	RawAggregations    json.RawMessage    `json:"aggregations"`
	AggregationsObject AggregationsResult `json:"-"`
	Total              SearchTotal        `json:"total"`
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
		return v, nil, fmt.Errorf("Search term must be at least three characters")
	}

	if len(personTypeFilters) == 0 {
		personTypeFilters = AllPersonTypes
	}

	data, err := json.Marshal(searchRequest{Term: term, PersonTypes: personTypeFilters, Limit: PageLimit, From: PageLimit * (page - 1)})
	if err != nil {
		return v, nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/lpa-api/v1/search/persons", bytes.NewReader(data))
	if err != nil {
		return v, nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return v, nil, newStatusError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return v, nil, err
	}

	if v.RawAggregations[0] != '[' {
		err = json.Unmarshal(v.RawAggregations, &v.AggregationsObject)
	}

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
