package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

// TODO: fix this cause its wrong when you're adding the final notified person
const maxNotifiedPersons = 4

func allowNewNotifiedPerson(currentCount int) bool {
	return currentCount < (maxNotifiedPersons)
}

type CreateNotifiedPersonClient interface {
	CreateNotifiedPerson(ctx sirius.Context, caseId int, notifiedPerson sirius.NotifiedPerson) error
	Lpa(ctx sirius.Context, caseId int) (sirius.Lpa, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	UpdateNotifiedPerson(ctx sirius.Context, notifiedPersonId int, notifiedPerson sirius.NotifiedPerson) error
}

type createNotifiedPersonData struct {
	XSRFToken              string
	IsPartial              bool
	NotifiedPerson         sirius.NotifiedPerson
	RelationshipToDonors   []sirius.RefDataItem
	Error                  sirius.ValidationError
	DonorId                int
	CaseId                 int
	IsEditing              bool
	Title                  string
	NextNotifiedPersonId   int
	AllowNewNotifiedPerson bool
	HtmxRedirect           string
	HtmxSwap               string
}

func CreateNotifiedPerson(client CreateNotifiedPersonClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorId, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		data := createNotifiedPersonData{
			XSRFToken: ctx.XSRFToken,
			IsPartial: r.Header.Get("HX-Request") == "true",
			DonorId:   donorId,
			CaseId:    caseId,
			Title:     "Add a notified person",
		}

		data.RelationshipToDonors, err = client.RefDataByCategory(ctx, sirius.RelationshipToDonorCategory)
		if err != nil {
			return err
		}

		lpa, err := client.Lpa(ctx, caseId)
		if err != nil {
			return err
		}

		data.AllowNewNotifiedPerson = allowNewNotifiedPerson(len(lpa.NotifiedPersons))

		var notifiedPersonId int
		notifiedPersonIdStr := r.FormValue("notifiedPersonId")
		isEditing := notifiedPersonIdStr != ""
		if isEditing {
			notifiedPersonId, err = strToIntOrStatusError(notifiedPersonIdStr)
			if err != nil {
				return err
			}
			for _, notifiedPerson := range lpa.NotifiedPersons {
				if notifiedPerson.ID == notifiedPersonId {
					data.NotifiedPerson = notifiedPerson
					break
				}
			}

			data.Title = "Update notified person details"
			data.IsEditing = true
			data.NextNotifiedPersonId = GetNextNotifiedPersonId(notifiedPersonId, lpa.NotifiedPersons)
		}

		if r.Method == http.MethodPost {
			notifiedPerson := sirius.NotifiedPerson{
				Person: sirius.Person{
					Salutation:        postFormString(r, "salutation"),
					Firstname:         postFormString(r, "firstname"),
					Middlenames:       postFormString(r, "middlenames"),
					Surname:           postFormString(r, "surname"),
					AddressLine1:      postFormString(r, "addressLine1"),
					AddressLine2:      postFormString(r, "addressLine2"),
					AddressLine3:      postFormString(r, "addressLine3"),
					Town:              postFormString(r, "town"),
					County:            postFormString(r, "county"),
					Country:           postFormString(r, "country"),
					Postcode:          postFormString(r, "postcode"),
					IsAirmailRequired: postFormString(r, "isAirmailRequired") == "true",
				},
				NoticeGivenDate: postFormDateString(r, "noticeGivenDate"),
			}
			data.NotifiedPerson = notifiedPerson

			if isEditing {
				err = client.UpdateNotifiedPerson(ctx, notifiedPersonId, notifiedPerson)
			} else {
				err = client.CreateNotifiedPerson(ctx, caseId, notifiedPerson)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve

				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("add-another-notified-person") != "" {
				if data.IsPartial {
					data.HtmxRedirect = fmt.Sprintf("/create-notified-person?id=%d&caseId=%d", donorId, caseId)
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-notified-person?id=%d&caseId=%d", donorId, caseId))
			}

			if r.FormValue("next-notified-person") != "" {
				if data.IsPartial {
					data.HtmxRedirect = fmt.Sprintf("/create-notified-person?id=%d&caseId=%d&notifiedPersonId=%d", donorId, caseId, data.NextNotifiedPersonId)
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-notified-person?id=%d&caseId=%d&notifiedPersonId=%d", donorId, caseId, data.NextNotifiedPersonId))
			}

			if data.IsPartial {
				data.HtmxRedirect = fmt.Sprintf("/create-lpa?id=%d&caseId=%d", donorId, caseId)
				data.HtmxSwap = "innerHTML show:#accordion-create-lpa-heading-2:top"
				return tmpl(w, data)
			}
			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#accordion-create-lpa-heading-2", donorId, caseId))

		}

		return tmpl(w, data)
	}
}

func GetNextNotifiedPersonId(id int, notifiedPersons []sirius.NotifiedPerson) int {
	nextNotifiedPersonId := 0
	for _, notifiedPerson := range notifiedPersons {
		if notifiedPerson.ID > id && (nextNotifiedPersonId == 0 || notifiedPerson.ID < nextNotifiedPersonId) {
			nextNotifiedPersonId = notifiedPerson.ID
		}
	}
	return nextNotifiedPersonId
}
