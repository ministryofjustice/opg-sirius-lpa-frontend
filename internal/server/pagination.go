package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type Pagination struct {
	SearchTerm  string
	Filters     string
	TotalItems  int
	CurrentPage int
	TotalPages  int
	PageSize    int
}

func newPagination(p *sirius.Pagination, s, f string) *Pagination {
	if p == nil {
		return nil
	}

	if s != "" {
		s = "?" + s
	}

	if f != "" {
		f = "&" + f
	}

	return &Pagination{
		SearchTerm:  s,
		Filters:     f,
		TotalItems:  p.TotalItems,
		CurrentPage: p.CurrentPage,
		TotalPages:  p.TotalPages,
		PageSize:    p.PageSize,
	}
}

func (p *Pagination) Start() int {
	return (p.CurrentPage-1)*p.PageSize + 1
}

func (p *Pagination) End() int {
	end := p.CurrentPage * p.PageSize
	if end < p.TotalItems {
		return end
	}

	return p.TotalItems
}

func (p *Pagination) HasPrevious() bool {
	return p.CurrentPage > 1
}

func (p *Pagination) PreviousPage() int {
	return p.CurrentPage - 1
}

func (p *Pagination) HasNext() bool {
	return p.TotalItems > p.CurrentPage*p.PageSize
}

func (p *Pagination) NextPage() int {
	return p.CurrentPage + 1
}

func (p *Pagination) Pages() []int {
	if p.TotalPages <= 7 {
		pages := make([]int, p.TotalPages)
		for i := 0; i < p.TotalPages; i++ {
			pages[i] = i + 1
		}
		return pages
	}

	pages := make([]int, 0, 7)

	if p.CurrentPage > 1 {
		prev := p.CurrentPage - 1

		switch prev {
		case 1:
			pages = append(pages, 1)
		case 2:
			pages = append(pages, 1, 2)
		default:
			pages = append(pages, 1, -1, prev)
		}
	}

	pages = append(pages, p.CurrentPage)

	if p.CurrentPage < p.TotalPages {
		next := p.CurrentPage + 1

		switch next {
		case p.TotalPages:
			pages = append(pages, p.TotalPages)
		case p.TotalPages - 1:
			pages = append(pages, next, p.TotalPages)
		default:
			pages = append(pages, next, -1, p.TotalPages)
		}
	}

	return pages
}
