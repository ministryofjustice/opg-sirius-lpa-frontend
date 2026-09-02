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

type mockCreateReplacementAttorneyClient struct {
	mock.Mock
}

func (m *mockCreateReplacementAttorneyClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockCreateReplacementAttorneyClient) CreateReplacementAttorney(ctx sirius.Context, caseId int, attorney sirius.Attorney) (sirius.Attorney, error) {
	args := m.Called(ctx, caseId, attorney)
	return args.Get(0).(sirius.Attorney), args.Error(1)
}

func (m *mockCreateReplacementAttorneyClient) UpdateReplacementAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error {
	args := m.Called(ctx, attorneyId, attorney)
	return args.Error(0)
}

var lpaWithReceiptDate = sirius.Lpa{
	Case: sirius.Case{
		ReceiptDate: sirius.DateString("2026-08-01"),
	},
}

func TestGetCreateReplacementAttorney(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			client := &mockCreateReplacementAttorneyClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(lpaWithReceiptDate, nil)

			expectedData := createReplacementAttorneyData{
				DonorId:   1,
				CaseId:    2,
				Title:     "Add a replacement attorney",
				IsPartial: isHtmx,
			}
			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, expectedData).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
			w := httptest.NewRecorder()

			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}

			err := CreateReplacementAttorney(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateReplacementAttorneyNoReceiptDate(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			client := &mockCreateReplacementAttorneyClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{}, nil)

			expectedData := createReplacementAttorneyData{
				DonorId: 1,
				CaseId:  2,
				Title:   "Add a replacement attorney",
				Error: sirius.ValidationError{
					Field: sirius.FieldErrors{
						"receiptDate": {"receiptDate": "A receipt date must be added to the LPA before a replacement attorney can be added"},
					},
				},
				IsPartial: isHtmx,
			}
			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, expectedData).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
			w := httptest.NewRecorder()

			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}

			err := CreateReplacementAttorney(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetEditReplacementAttorney(t *testing.T) {
	existingAttorney := sirius.Attorney{
		Person: sirius.Person{
			ID:        4,
			Firstname: "Rudolph",
			Surname:   "Stotesbury",
		},
	}

	client := &mockCreateReplacementAttorneyClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{
			Case: sirius.Case{
				ReceiptDate:          sirius.DateString("2026-08-01"),
				ReplacementAttorneys: []sirius.Attorney{existingAttorney},
			},
		}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createReplacementAttorneyData{
			DonorId:   1,
			CaseId:    2,
			Attorney:  existingAttorney,
			IsEditing: true,
			Title:     "Update replacement attorney details",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&attorneyId=4", nil)
	w := httptest.NewRecorder()

	err := CreateReplacementAttorney(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateReplacementAttorneyBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := CreateReplacementAttorney(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetCreateReplacementAttorneyWhenLpaErrors(t *testing.T) {
	client := &mockCreateReplacementAttorneyClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := CreateReplacementAttorney(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateReplacementAttorney(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			dateString := "2022-04-05"
			attorney := sirius.Attorney{
				Person: sirius.Person{
					Salutation:   "Rev",
					Firstname:    "Rudolph",
					Middlenames:  "Modesto",
					Surname:      "Stotesbury",
					DateOfBirth:  sirius.DateString(dateString),
					AddressLine1: "Rotonda Gerardo 769",
					AddressLine2: "Appartamento 94",
					AddressLine3: "Augusto terme",
					Town:         "San Sabazio",
					County:       "Benevento",
					Postcode:     "57797",
					Country:      "Italy",
					PhoneNumber:  "079876543345",
					Email:        "rm2@email.test",
				},
			}
			client := &mockCreateReplacementAttorneyClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(lpaWithReceiptDate, nil).
				On("CreateReplacementAttorney", mock.Anything, 2, attorney).
				Return(sirius.Attorney{Person: sirius.Person{ID: 88}}, nil)

			template := &mockTemplate{}

			if isHtmx {
				template.
					On("Func", mock.Anything, createReplacementAttorneyData{
						DonorId:      1,
						CaseId:       2,
						Attorney:     attorney,
						Title:        "Add a replacement attorney",
						HtmxRedirect: "/create-lpa?id=1&caseId=2",
						HtmxSwap:     "innerHTML show:#scroll-to-replacement-attorneys:top",
						IsPartial:    true,
					}).
					Return(nil)
			}

			form := url.Values{
				"salutation":   {"Rev"},
				"firstname":    {"Rudolph"},
				"middlenames":  {"Modesto"},
				"surname":      {"Stotesbury"},
				"dob":          {dateString},
				"addressLine1": {"Rotonda Gerardo 769"},
				"addressLine2": {"Appartamento 94"},
				"addressLine3": {"Augusto terme"},
				"town":         {"San Sabazio"},
				"county":       {"Benevento"},
				"postcode":     {"57797"},
				"country":      {"Italy"},
				"phoneNumber":  {"079876543345"},
				"email":        {"rm2@email.test"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateReplacementAttorney(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedError := RedirectError("/create-lpa?id=1&caseId=2#scroll-to-replacement-attorneys")
				assert.Equal(t, err, expectedError)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostEditReplacementAttorney(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			existingAttorney := sirius.Attorney{Person: sirius.Person{ID: 4}}
			updatedAttorney := sirius.Attorney{
				Person: sirius.Person{
					Firstname: "Rudolph",
					Surname:   "Stotesbury",
				},
			}

			client := &mockCreateReplacementAttorneyClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{
					Case: sirius.Case{
						ReceiptDate:          sirius.DateString("2026-08-01"),
						ReplacementAttorneys: []sirius.Attorney{existingAttorney},
					},
				}, nil).
				On("UpdateReplacementAttorney", mock.Anything, 4, updatedAttorney).
				Return(nil)

			template := &mockTemplate{}

			if isHtmx {
				template.
					On("Func", mock.Anything, createReplacementAttorneyData{
						DonorId:      1,
						CaseId:       2,
						Attorney:     updatedAttorney,
						IsEditing:    true,
						Title:        "Update replacement attorney details",
						HtmxRedirect: "/create-lpa?id=1&caseId=2",
						HtmxSwap:     "innerHTML show:#scroll-to-replacement-attorneys:top",
						IsPartial:    true,
					}).
					Return(nil)
			}

			form := url.Values{
				"firstname": {"Rudolph"},
				"surname":   {"Stotesbury"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&attorneyId=4", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateReplacementAttorney(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedError := RedirectError("/create-lpa?id=1&caseId=2#scroll-to-replacement-attorneys")
				assert.Equal(t, err, expectedError)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateReplacementAttorneyWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	attorney := sirius.Attorney{
		Person: sirius.Person{
			Firstname: "Rudolph",
			Surname:   "Stotesbury",
		},
	}

	client := &mockCreateReplacementAttorneyClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(lpaWithReceiptDate, nil).
		On("CreateReplacementAttorney", mock.Anything, 2, attorney).
		Return(sirius.Attorney{}, expectedError)

	template := &mockTemplate{}

	expectedData := createReplacementAttorneyData{
		Attorney: attorney,
		DonorId:  1,
		CaseId:   2,
		Error:    expectedError,
		Title:    "Add a replacement attorney",
	}

	template.
		On("Func", mock.Anything, expectedData).
		Return(nil)

	form := url.Values{
		"firstname": {"Rudolph"},
		"surname":   {"Stotesbury"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateReplacementAttorney(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateReplacementAttorneyAddAnother(t *testing.T) {
	attorney := sirius.Attorney{
		Person: sirius.Person{
			Firstname: "Rudolph",
			Surname:   "Stotesbury",
		},
	}
	client := &mockCreateReplacementAttorneyClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(lpaWithReceiptDate, nil).
		On("CreateReplacementAttorney", mock.Anything, 2, attorney).
		Return(sirius.Attorney{Person: sirius.Person{ID: 88}}, nil)

	form := url.Values{
		"firstname":   {"Rudolph"},
		"surname":     {"Stotesbury"},
		"add-another": {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateReplacementAttorney(client, nil)(w, r)
	resp := w.Result()

	expectedRedirect := RedirectError("/create-replacement-attorney?id=1&caseId=2")
	assert.Equal(t, err, expectedRedirect)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client)
}
