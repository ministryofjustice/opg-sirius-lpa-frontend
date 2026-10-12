package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

const maxNotifiedPersons = 5

func allowNewNotifiedPerson(currentCount int) bool {
	return currentCount < (maxNotifiedPersons)
}

type CreateNotifiedPersonClient interface {
	CreateNotifiedPerson(ctx sirius.Context, caseId int, notifiedPerson sirius.NotifiedPerson) error
	Lpa(ctx sirius.Context, caseId int) (sirius.Lpa, error)
}

type createNotifiedPersonData struct {
	XSRFToken              string
	IsPartial              bool
	NotifiedPerson         sirius.NotifiedPerson
	Error                  sirius.ValidationError
	DonorId                int
	CaseId                 int
	IsEditing              bool
	Title                  string
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
			IsPartial: ctx.IsPartial,
			DonorId:   donorId,
			CaseId:    caseId,
			Title:     "Add a notified person",
			IsEditing: false,
		}

		lpa, err := client.Lpa(ctx, caseId)
		if err != nil {
			return err
		}

		// +1 because this is before we save the final allowed notified person
		data.AllowNewNotifiedPerson = allowNewNotifiedPerson(len(lpa.NotifiedPersons) + 1)

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

			err = client.CreateNotifiedPerson(ctx, caseId, notifiedPerson)

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

			if data.IsPartial {
				data.HtmxRedirect = fmt.Sprintf("/create-lpa?id=%d&caseId=%d", donorId, caseId)
				data.HtmxSwap = "innerHTML show:#scroll-to-notified-person:top"
				return tmpl(w, data)
			}
			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#scroll-to-notified-person", donorId, caseId))

		}

		return tmpl(w, data)
	}
}
