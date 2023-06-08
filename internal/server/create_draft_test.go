package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateDraftClient struct {
	mock.Mock
}

func (m *mockCreateDraftClient) CreateDraft(ctx sirius.Context, draft sirius.Draft) (map[string]string, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *mockCreateDraftClient) GetUserDetails(ctx sirius.Context) (sirius.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.User), args.Error(1)
}

func TestGetCreateDraft(t *testing.T) {
	client := &mockCreateDraftClient{}
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{Roles: []string{"private-mlpa"}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDraftData{}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/digital-lpas/create", nil)
	w := httptest.NewRecorder()

	err := CreateDraft(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateDraftForbidden(t *testing.T) {
	client := &mockCreateDraftClient{}
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{}, nil)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/digital-lpas/create", nil)
	w := httptest.NewRecorder()

	err := CreateDraft(client, template.Func)(w, r)

	assert.Equal(t, sirius.StatusError{Code: 403}, err)

	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDraft(t *testing.T) {
	client := &mockCreateDraftClient{}
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{Roles: []string{"private-mlpa"}}, nil)
	client.
		On("CreateDraft", mock.Anything, sirius.Draft{
			CaseType:          []string{"pfa", "hw"},
			Source:            "PHONE",
			DonorName:         "Gerald Ryan Sandel",
			CorrespondentName: "Rosalinda Langdale",
			DonorDob:          sirius.DateString("1943-03-06"),
			Email:             "gerald.sandel@somehost.example",
			PhoneNumber:       "01638294820",
			DonorAddress: sirius.Address{
				Line1:    "Bradtke",
				Line2:    "Zipper House",
				Line3:    "Mills Ports",
				Town:     "Deerfield Beach",
				Postcode: "QY9 9QW",
				Country:  "GB",
			},
			CorrespondentAddress: &sirius.Address{
				Line1:    "Intensity Office",
				Line2:    "Lind Run",
				Line3:    "Hendersonville",
				Town:     "Moline",
				Postcode: "OE6 2DV",
				Country:  "GB",
			},
		}).
		Return(map[string]string{
			"pfa": "M-0123-4567-8901",
		}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDraftData{
			Draft: draft{
				SubTypes:        []string{"pfa", "hw"},
				DonorFirstname:  "Gerald",
				DonorMiddlename: "Ryan",
				DonorSurname:    "Sandel",
				Dob:             dob{Day: 6, Month: 3, Year: 1943},
				Email:           "gerald.sandel@somehost.example",
				Phone:           "01638294820",
				DonorAddress: sirius.Address{
					Line1:    "Bradtke",
					Line2:    "Zipper House",
					Line3:    "Mills Ports",
					Town:     "Deerfield Beach",
					Postcode: "QY9 9QW",
					Country:  "GB",
				},
				Recipient:               "other",
				CorrespondentFirstname:  "Rosalinda",
				CorrespondentMiddlename: "",
				CorrespondentSurname:    "Langdale",
				CorrespondentAddress: sirius.Address{
					Line1:    "Intensity Office",
					Line2:    "Lind Run",
					Line3:    "Hendersonville",
					Town:     "Moline",
					Postcode: "OE6 2DV",
					Country:  "GB",
				},
			},
			Success: true,
			Uids: []createDraftResult{
				{Subtype: "pfa", Uid: "M-0123-4567-8901"},
			},
		}).
		Return(nil)

	form := url.Values{
		"subtype":                   {"pfa", "hw"},
		"donorFirstname":            {"Gerald"},
		"donorMiddlename":           {"Ryan"},
		"donorSurname":              {"Sandel"},
		"dobDay":                    {"6"},
		"dobMonth":                  {"3"},
		"dobYear":                   {"1943"},
		"donorEmail":                {"gerald.sandel@somehost.example"},
		"donorPhone":                {"01638294820"},
		"donorAddressLine1":         {"Bradtke"},
		"donorAddressLine2":         {"Zipper House"},
		"donorAddressLine3":         {"Mills Ports"},
		"donorTown":                 {"Deerfield Beach"},
		"donorPostcode":             {"QY9 9QW"},
		"donorCountry":              {"GB"},
		"recipient":                 {"other"},
		"correspondentFirstname":    {"Rosalinda"},
		"correspondentMiddlename":   {""},
		"correspondentSurname":      {"Langdale"},
		"correspondentAddressLine1": {"Intensity Office"},
		"correspondentAddressLine2": {"Lind Run"},
		"correspondentAddressLine3": {"Hendersonville"},
		"correspondentTown":         {"Moline"},
		"correspondentPostcode":     {"OE6 2DV"},
		"correspondentCountry":      {"GB"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/digital-lpas/create", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateDraft(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDraftWhenAPIFails(t *testing.T) {
	client := &mockCreateDraftClient{}
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{Roles: []string{"private-mlpa"}}, nil)
	client.
		On("CreateDraft", mock.Anything, sirius.Draft{
			Source:    "PHONE",
			DonorName: "Gerald Sandel",
			DonorAddress: sirius.Address{
				Country: "GB",
			},
		}).
		Return(map[string]string{}, expectedError)

	template := &mockTemplate{}

	form := url.Values{
		"donorFirstname": {"Gerald"},
		"donorSurname":   {"Sandel"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/digital-lpas/create", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateDraft(client, template.Func)(w, r)
	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDraftWhenValidationError(t *testing.T) {
	client := &mockCreateDraftClient{}
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{Roles: []string{"private-mlpa"}}, nil)
	client.
		On("CreateDraft", mock.Anything, sirius.Draft{
			Source:    "PHONE",
			DonorName: "Gerald",
			DonorAddress: sirius.Address{
				Country: "GB",
			},
		}).
		Return(map[string]string{}, sirius.ValidationError{Field: sirius.FieldErrors{
			"surname": {"required": "This field is required"},
		}})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDraftData{
			Draft: draft{
				DonorFirstname: "Gerald",
			},
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"surname": {"required": "This field is required"},
				},
			},
		}).
		Return(nil)

	form := url.Values{
		"donorFirstname": {"Gerald"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/digital-lpas/create", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateDraft(client, template.Func)(w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
