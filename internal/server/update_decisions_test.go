package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUpdateDecisionsClient struct {
	mock.Mock
}

func (m *mockUpdateDecisionsClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockUpdateDecisionsClient) UpdateDecisions(ctx sirius.Context, uid string, decisions sirius.UpdateDecisions) error {
	args := m.Called(ctx, uid, decisions)
	return args.Error(0)
}

var updateDecisionsCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-1111-2222-3333",
		LpaStoreData: sirius.LpaStoreData{
			Attorneys: []sirius.LpaStoreAttorney{
				{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "302b05c7-896c-4290-904e-2005e4f1e81e",
						FirstNames: "Jack",
						LastName:   "Black",
						Address: sirius.LpaStoreAddress{
							Line1:    "9 Mount Pleasant Drive",
							Town:     "East Harling",
							Postcode: "NR16 2GB",
							Country:  "UK",
						},
					},
					DateOfBirth:     "1990-02-22",
					Status:          shared.ActiveAttorneyStatus.String(),
					AppointmentType: shared.OriginalAppointmentType.String(),
					Email:           "a@example.com",
					Mobile:          "077577575757",
					SignedAt:        "2024-01-12T10:09:09Z",
				},
				{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "123a01b1-456d-5391-813d-2010d3e2d72d",
						FirstNames: "Jack",
						LastName:   "White",
						Address: sirius.LpaStoreAddress{
							Line1:    "29 Grange Road",
							Town:     "Birmingham",
							Postcode: "B29 6BL",
							Country:  "UK",
						},
					},
					DateOfBirth:     "1990-02-22",
					Status:          shared.InactiveAttorneyStatus.String(),
					AppointmentType: shared.ReplacementAppointmentType.String(),
					Email:           "c@example.com",
					Mobile:          "07122121212",
					SignedAt:        "2024-11-28T19:22:11Z",
				},
			},
		},
	},
}

func TestGetUpdateDecisionsGet(t *testing.T) {
	client := &mockUpdateDecisionsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(updateDecisionsCaseSummary, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			updateDecisionsData{
				CaseSummary:         updateDecisionsCaseSummary,
				Form:                formDecisionsDetails{},
				ActiveAttorneyCount: 1,
			}).
		Return(errExample)

	server := newMockServer("/lpa/{uid}/update-decisions", UpdateDecisions(client, template.Func))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/update-decisions", nil)
	_, err := server.serve(r)

	assert.Equal(t, errExample, err)
}

func TestGetUpdateDecisionsGetWhenCaseSummaryErrors(t *testing.T) {
	client := &mockUpdateDecisionsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(sirius.CaseSummary{}, errExample)

	server := newMockServer("/lpa/{uid}/update-decisions", UpdateDecisions(client, nil))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/update-decisions", nil)
	_, err := server.serve(r)

	assert.Equal(t, errExample, err)
}

func TestPostUpdateDecisions(t *testing.T) {
	testcases := map[string]struct {
		updateError   error
		expectedError error
	}{
		"success": {
			expectedError: RedirectError("/lpa/M-1111-2222-3333/lpa-details"),
		},
		"failure": {
			updateError:   errExample,
			expectedError: errExample,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			form := url.Values{
				"howAttorneysMakeDecisions": {"jointly"},
			}

			client := &mockUpdateDecisionsClient{}
			client.
				On("CaseSummary", mock.Anything, "M-1111-2222-3333").
				Return(sirius.CaseSummary{}, nil)
			client.
				On("UpdateDecisions", mock.Anything, "M-1111-2222-3333", sirius.UpdateDecisions{
					HowAttorneysMakeDecisions: "jointly",
				}).
				Return(tc.updateError)

			server := newMockServer("/lpa/{uid}/update-decisions", UpdateDecisions(client, nil))

			r, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/update-decisions", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			_, err := server.serve(r)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestPostUpdateDecisionsWhenValidationError(t *testing.T) {
	form := url.Values{
		"howAttorneysMakeDecisions": {"not a real option"},
	}

	client := &mockUpdateDecisionsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(sirius.CaseSummary{}, nil)
	client.
		On("UpdateDecisions", mock.Anything, "M-1111-2222-3333", mock.Anything).
		Return(sirius.ValidationError{
			Field: sirius.FieldErrors{
				"howAttorneysMakeDecisions": {"invalid": "Value not a valid option"},
			},
		})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			mock.MatchedBy(func(data updateDecisionsData) bool {
				return data.Error.Field["howAttorneysMakeDecisions"]["invalid"] == "Value not a valid option"
			}),
		).
		Return(nil)

	server := newMockServer("/lpa/{uid}/update-decisions", UpdateDecisions(client, template.Func))

	r, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/update-decisions", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(r)

	assert.Nil(t, err)
}
