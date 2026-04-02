package server

import (
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type SearchClient interface {
	Search(ctx sirius.Context, term string, page int, personTypeFilters []string) (sirius.SearchResponse, *sirius.Pagination, error)
	DeletedCases(ctx sirius.Context, uid string) ([]sirius.DeletedCase, error)
}

type searchData struct {
	Results      []sirius.Person
	Total        int
	Aggregations sirius.Aggregations
	Filters      searchFilters
	SearchTerm   string
	Pagination   *Pagination
	DeletedCases []sirius.DeletedCase
}

type searchFilters struct {
	Set        bool
	PersonType []string
}

func (f searchFilters) Encode() string {
	if !f.Set {
		return ""
	}

	form := url.Values{}
	for _, v := range f.PersonType {
		form.Add("person-type", v)
	}

	return form.Encode()
}

func newSearchFilters(form url.Values) searchFilters {
	filters := searchFilters{}

	if selectedPersonType, ok := form["person-type"]; ok {
		for _, spt := range selectedPersonType {
			for _, pt := range sirius.AllPersonTypes {
				if spt == pt {
					filters.PersonType = append(filters.PersonType, spt)
					filters.Set = true
				}
			}
		}
	}

	return filters
}

func Search(client SearchClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		searchTerm := r.FormValue("term")
		search := url.Values{}
		search.Add("term", r.FormValue("term"))

		filters := newSearchFilters(r.Form)

		data := searchData{
			SearchTerm: searchTerm,
			Filters:    filters,
		}

		// If no search term, just render the template (front-end will handle the empty state)
		if searchTerm == "" {
			return tmpl(w, data)
		}

		results, pagination, err := client.Search(ctx, searchTerm, getPage(r), filters.PersonType)
		if err != nil {
			return err
		}

		if results.Total.Count == 0 {
			re := regexp.MustCompile(`\D+`)
			input := re.ReplaceAllString(searchTerm, "")
			isUid, err := regexp.MatchString(`^\d{12}$`, input)
			if err != nil {
				return err
			}

			if isUid {
				data.DeletedCases, err = client.DeletedCases(ctx, input)
				for i := range data.DeletedCases {
					data.DeletedCases[i].DeletedStatus = shared.CaseStatusTypeDeleted
				}
				if err != nil {
					return err
				}
			}
		}

		total := 0
		hasFilters := len(filters.PersonType) > 0
		for personType, personTypeCount := range results.Aggregations.PersonType {
			if !hasFilters || slices.Contains(filters.PersonType, personType) {
				total += personTypeCount
			}
		}

		data.Total = total
		data.Results = results.Results
		data.Aggregations = results.Aggregations
		data.Pagination = newPagination(pagination, search.Encode(), filters.Encode())

		return tmpl(w, data)
	}
}

func getPage(r *http.Request) int {
	page := r.FormValue("page")
	if page == "" {
		return 1
	}

	v, err := strconv.Atoi(page)
	if err != nil {
		return 1
	}

	return v
}
