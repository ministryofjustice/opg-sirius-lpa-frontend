package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type mockChangeDraftClient struct {
	mock.Mock
}

func (m *mockChangeDraftClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockChangeDraftClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockChangeDraftClient) ChangeDraft(ctx sirius.Context, caseUID string, draftData sirius.ChangeDraft) error {
	return m.Called(ctx, caseUID, draftData).Error(0)
}

var testChangeDraftCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-EEEE-EEEE-EEEE",
		SiriusData: sirius.SiriusData{
			ID:  12,
			UID: "M-EEEE-EEEE-EEEE",
			Application: sirius.Draft{
				DonorFirstNames: "Jack",
				DonorLastName:   "Black",
				DonorDob:        "1990-02-22",
				DonorAddress: sirius.Address{
					Line1:    "9 Mount Pleasant Drive",
					Town:     "East Harling",
					Postcode: "NR16 2GB",
					Country:  "UK",
				},
				PhoneNumber:               "077577575757",
				Email:                     "a@example.com",
				CorrespondenceByWelsh:     false,
				CorrespondenceLargeFormat: false,
			},
			Donor: sirius.Donor{
				ID:           1,
				Firstname:    "Jack",
				Surname:      "Black",
				DateOfBirth:  "1990-02-22",
				AddressLine1: "9 Mount Pleasant Drive",
				Town:         "East Harling",
				Postcode:     "NR16 2GB",
				Country:      "UK",
				Phone:        "077577575757",
				Email:        "a@example.com",
			},
		},
	},
}

func TestGetChangeDraft(t *testing.T) {
	tests := []struct {
		name          string
		caseUID       string
		form          formDraftDetails
		errorReturned error
	}{
		{
			name:    "Change Draft Details",
			caseUID: "M-EEEE-EEEE-EEEE",
			form: formDraftDetails{
				FirstNames:  "Jack",
				LastName:    "Black",
				DateOfBirth: dob{Day: 22, Month: 2, Year: 1990},
				Address: sirius.Address{
					Line1:    "9 Mount Pleasant Drive",
					Town:     "East Harling",
					Postcode: "NR16 2GB",
					Country:  "UK",
				},
				Email:       "a@example.com",
				PhoneNumber: "077577575757",
			},
			errorReturned: nil,
		},
		{
			name:    "Template Error Returned",
			caseUID: "M-EEEE-EEEE-EEEE",
			form: formDraftDetails{
				FirstNames:  "Jack",
				LastName:    "Black",
				DateOfBirth: dob{Day: 22, Month: 2, Year: 1990},
				Address: sirius.Address{
					Line1:    "9 Mount Pleasant Drive",
					Town:     "East Harling",
					Postcode: "NR16 2GB",
					Country:  "UK",
				},
				Email:       "a@example.com",
				PhoneNumber: "077577575757",
			},
			errorReturned: expectedError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockChangeDraftClient{}
			client.
				On("CaseSummary", mock.Anything, tc.caseUID).
				Return(testChangeDraftCaseSummary, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
				Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything,
					changeDraftData{
						Countries: []sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}},
						CaseUID:   tc.caseUID,
						Form:      tc.form,
					}).
				Return(tc.errorReturned)

			server := newMockServer("/lpa/{uid}/change-draft", ChangeDraft(client, template.Func))

			r, _ := http.NewRequest(http.MethodGet, "/lpa/"+tc.caseUID+"/change-draft", nil)
			_, err := server.serve(r)

			if tc.errorReturned != nil {
				assert.Equal(t, tc.errorReturned, err)
			} else {
				assert.Nil(t, err)
			}

			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetChangeDraftWhenCaseSummaryErrors(t *testing.T) {
	client := &mockChangeDraftClient{}
	client.
		On("CaseSummary", mock.Anything, "M-FFFF-FFFF-FFFF").
		Return(sirius.CaseSummary{}, expectedError)

	assertChangeDraftErrors(t, client, "M-FFFF-FFFF-FFFF", expectedError)
}

func TestGetChangeDraftWhenRefDataByCategoryErrors(t *testing.T) {
	client := &mockChangeDraftClient{}
	client.
		On("CaseSummary", mock.Anything, "M-FFFF-FFFF-FFFF").
		Return(testChangeDraftCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	assertChangeDraftErrors(t, client, "M-FFFF-FFFF-FFFF", expectedError)
}

func assertChangeDraftErrors(t *testing.T, client *mockChangeDraftClient, uid string, expectedError error) {
	server := newMockServer("/lpa/{uid}/change-draft", ChangeDraft(client, nil))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/"+uid+"/change-draft", nil)
	_, err := server.serve(r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostChangeDraft(t *testing.T) {
	tests := []struct {
		name          string
		apiError      error
		expectedError error
	}{
		{
			name:          "Post form successfully submits",
			apiError:      nil,
			expectedError: RedirectError("/lpa/M-EEEE-EEEE-EEEE/lpa-details"),
		},
		{
			name:          "Post form returns an API failure",
			apiError:      expectedError,
			expectedError: expectedError,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockChangeDraftClient{}
			client.
				On("CaseSummary", mock.Anything, "M-EEEE-EEEE-EEEE").
				Return(testChangeDraftCaseSummary, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
				Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)
			client.
				On("ChangeDraft", mock.Anything, "M-EEEE-EEEE-EEEE", sirius.ChangeDraft{
					FirstNames:  "Samuel",
					LastName:    "Smith",
					DateOfBirth: "1991-01-01",
					Address: sirius.Address{
						Line1:    "9 Mount",
						Line2:    "Pleasant Drive",
						Line3:    "Norwich",
						Town:     "East Harling",
						Postcode: "NR16 2GB",
						Country:  "GB",
					},
					Phone: "123456789",
					Email: "test@test.com",
				}).
				Return(tc.apiError)

			template := &mockTemplate{}

			server := newMockServer("/lpa/{uid}/change-draft", ChangeDraft(client, template.Func))

			form := url.Values{
				"firstNames":       {"Samuel"},
				"lastName":         {"Smith"},
				"dob.day":          {"1"},
				"dob.month":        {"1"},
				"dob.year":         {"1991"},
				"address.Line1":    {"9 Mount"},
				"address.Line2":    {"Pleasant Drive"},
				"address.Line3":    {"Norwich"},
				"address.Town":     {"East Harling"},
				"address.Postcode": {"NR16 2GB"},
				"address.Country":  {"GB"},
				"phoneNumber":      {"123456789"},
				"email":            {"test@test.com"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/lpa/M-EEEE-EEEE-EEEE/change-draft", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			_, err := server.serve(r)

			assert.Equal(t, tc.expectedError, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostChangeDraftWhenValidationError(t *testing.T) {
	client := &mockChangeDraftClient{}
	client.
		On("CaseSummary", mock.Anything, "M-EEEE-EEEE-EEEE").
		Return(testChangeDraftCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)
	client.
		On("ChangeDraft", mock.Anything, "M-EEEE-EEEE-EEEE", sirius.ChangeDraft{
			LastName:    "Black",
			DateOfBirth: "1990-02-22",
			Address: sirius.Address{
				Line1:    "9 Mount Pleasant Drive",
				Town:     "East Harling",
				Postcode: "NR16 2GB",
				Country:  "UK",
			},
			Phone: "077577575757",
			Email: "a@example.com",
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"firstNames": {"required": "Value required and can't be empty"},
		}})

	template := &mockTemplate{}

	template.
		On("Func", mock.Anything,
			mock.MatchedBy(func(data changeDraftData) bool {
				return data.Error.Field["firstNames"]["required"] == "Value required and can't be empty"
			}),
		).
		Return(nil)

	server := newMockServer("/lpa/{uid}/change-draft", ChangeDraft(client, template.Func))

	form := url.Values{
		"firstNames": {""},
	}

	r, _ := http.NewRequest(http.MethodPost, "/lpa/M-EEEE-EEEE-EEEE/change-draft", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
