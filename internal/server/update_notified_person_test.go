package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUpdateNotifiedPersonClient struct {
	mock.Mock
}

func (m *mockUpdateNotifiedPersonClient) Lpa(ctx sirius.Context, caseId int) (sirius.Lpa, error) {
	args := m.Called(ctx, caseId)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockUpdateNotifiedPersonClient) UpdateNotifiedPerson(ctx sirius.Context, notifiedPersonId int, notifiedPerson sirius.NotifiedPerson) error {
	args := m.Called(ctx, notifiedPersonId, notifiedPerson)
	return args.Error(0)
}

func TestGetEditNotifiedPerson(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			existingNotifiedPerson := sirius.NotifiedPerson{
				Person: sirius.Person{
					ID:        4,
					Firstname: "Rudolph",
					Surname:   "Stotesbury",
				},
				NoticeGivenDate: sirius.DateString("2022-04-05"),
			}

			client := &mockUpdateNotifiedPersonClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{NotifiedPersons: []sirius.NotifiedPerson{existingNotifiedPerson}}}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, updateNotifiedPersonData{
					IsPartial:      isHtmx,
					DonorId:        1,
					CaseId:         2,
					NotifiedPerson: existingNotifiedPerson,
					IsEditing:      true,
					Title:          "Update notified person details",
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&notifiedPersonId=4", nil)
			w := httptest.NewRecorder()

			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}

			err := UpdateNotifiedPerson(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostEditNotifiedPerson(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			dateString := "2022-04-05"
			existingNotifiedPerson := sirius.NotifiedPerson{Person: sirius.Person{ID: 4}}
			updatedNotifiedPerson := sirius.NotifiedPerson{
				Person: sirius.Person{
					Salutation:        "Rev",
					Firstname:         "Rudolph",
					Middlenames:       "Modesto",
					Surname:           "Stotesbury",
					AddressLine1:      "Rotonda Gerardo 769",
					AddressLine2:      "Appartamento 94",
					AddressLine3:      "Augusto terme",
					Town:              "San Sabazio",
					County:            "Benevento",
					Postcode:          "57797",
					Country:           "Italy",
					IsAirmailRequired: true,
				},
				NoticeGivenDate: sirius.DateString(dateString),
			}

			client := &mockUpdateNotifiedPersonClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{NotifiedPersons: []sirius.NotifiedPerson{existingNotifiedPerson}}}, nil).
				On("UpdateNotifiedPerson", mock.Anything, 4, updatedNotifiedPerson).
				Return(nil)

			template := &mockTemplate{}
			if isHtmx {
				template.
					On("Func", mock.Anything, updateNotifiedPersonData{
						IsPartial:      true,
						DonorId:        1,
						CaseId:         2,
						NotifiedPerson: updatedNotifiedPerson,
						IsEditing:      true,
						Title:          "Update notified person details",
						HtmxRedirect:   "/create-lpa?id=1&caseId=2",
						HtmxSwap:       "innerHTML show:#scroll-to-notified-person:top",
					}).
					Return(nil)
			}

			form := url.Values{
				"salutation":        {"Rev"},
				"firstname":         {"Rudolph"},
				"middlenames":       {"Modesto"},
				"surname":           {"Stotesbury"},
				"noticeGivenDate":   {dateString},
				"addressLine1":      {"Rotonda Gerardo 769"},
				"addressLine2":      {"Appartamento 94"},
				"addressLine3":      {"Augusto terme"},
				"town":              {"San Sabazio"},
				"county":            {"Benevento"},
				"postcode":          {"57797"},
				"country":           {"Italy"},
				"isAirmailRequired": {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&notifiedPersonId=4", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := UpdateNotifiedPerson(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedError := RedirectError("/create-lpa?id=1&caseId=2#scroll-to-notified-person")
				assert.Equal(t, expectedError, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostEditNotifiedPersonNextAnother(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			dateString := "2022-04-05"
			existingNotifiedPerson := sirius.NotifiedPerson{Person: sirius.Person{ID: 4}}
			updatedNotifiedPerson := sirius.NotifiedPerson{
				Person: sirius.Person{
					Salutation:        "Rev",
					Firstname:         "Rudolph",
					Middlenames:       "Modesto",
					Surname:           "Stotesbury",
					AddressLine1:      "Rotonda Gerardo 769",
					AddressLine2:      "Appartamento 94",
					AddressLine3:      "Augusto terme",
					Town:              "San Sabazio",
					County:            "Benevento",
					Postcode:          "57797",
					Country:           "Italy",
					IsAirmailRequired: true,
				},
				NoticeGivenDate: sirius.DateString(dateString),
			}

			client := &mockUpdateNotifiedPersonClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{NotifiedPersons: []sirius.NotifiedPerson{
					existingNotifiedPerson,
					{Person: sirius.Person{ID: 5}},
				}}}, nil).
				On("UpdateNotifiedPerson", mock.Anything, 4, updatedNotifiedPerson).
				Return(nil)

			template := &mockTemplate{}
			if isHtmx {
				template.
					On("Func", mock.Anything, updateNotifiedPersonData{
						IsPartial:            true,
						DonorId:              1,
						CaseId:               2,
						NotifiedPerson:       updatedNotifiedPerson,
						IsEditing:            true,
						Title:                "Update notified person details",
						NextNotifiedPersonId: 5,
						HtmxRedirect:         "/update-notified-person?id=1&caseId=2&notifiedPersonId=5",
						HtmxSwap:             "innerHTML scroll:.action-panel__content:top",
					}).
					Return(nil)
			}

			form := url.Values{
				"salutation":           {"Rev"},
				"firstname":            {"Rudolph"},
				"middlenames":          {"Modesto"},
				"surname":              {"Stotesbury"},
				"noticeGivenDate":      {dateString},
				"addressLine1":         {"Rotonda Gerardo 769"},
				"addressLine2":         {"Appartamento 94"},
				"addressLine3":         {"Augusto terme"},
				"town":                 {"San Sabazio"},
				"county":               {"Benevento"},
				"postcode":             {"57797"},
				"country":              {"Italy"},
				"isAirmailRequired":    {"true"},
				"next-notified-person": {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&notifiedPersonId=4", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := UpdateNotifiedPerson(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedRedirect := RedirectError("/update-notified-person?id=1&caseId=2&notifiedPersonId=5")
				assert.Equal(t, expectedRedirect, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
