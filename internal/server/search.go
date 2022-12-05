package server

import (
	"errors"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"net/url"
	"regexp"
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
		if searchTerm == "" {
			return errors.New("search term required")
		}
		sv := url.Values{}
		sv.Add("term", r.FormValue("term"))

		filters := newSearchFilters(r.Form)

		data := searchData{
			SearchTerm: searchTerm,
			Filters:    filters,
		}

		results, pagination, err := client.Search(ctx, searchTerm, getPage(r), filters.PersonType)
		if err != nil {
			return err
		}

		if results.Total.Count == 0 {
			re := regexp.MustCompile(`\D+`)
			formattedTerm := re.ReplaceAllString(searchTerm, "")
			isUid, err := regexp.MatchString(`^\d{12}$`, re.ReplaceAllString(formattedTerm, ""))
			if err != nil {
				return err
			}

			if isUid {
				data.DeletedCases, err = client.DeletedCases(ctx, formattedTerm)
				if err != nil {
					return err
				}
			}
		}

		data.Results = results.Results
		data.Total = results.Total.Count
		data.Aggregations = results.Aggregations
		data.Pagination = newPagination(pagination, sv.Encode(), filters.Encode())

		return tmpl(w, data)
	}
}
