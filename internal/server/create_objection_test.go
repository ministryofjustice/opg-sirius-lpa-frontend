package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateObjectionClient struct {
	mock.Mock
}

func (m *mockCreateObjectionClient) CreateObjection(ctx sirius.Context, objection sirius.ObjectionRequest) error {
	args := m.Called(ctx, objection)
	return args.Error(0)
}

func (m *mockCreateObjectionClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

var testCreateObjectionsCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-9898-9898-9898",
		SiriusData: sirius.SiriusData{
			ID:      676,
			UID:     "M-9898-9898-9898",
			Subtype: "personal-welfare",
			Status:  shared.CaseStatusTypeDraft,
			LinkedCases: []sirius.SiriusData{
				{
					UID:     "M-9999-9999-9999",
					Subtype: "personal-welfare",
					Status:  shared.CaseStatusTypeInProgress,
				},
				{
					UID:     "M-8888-8888-8888",
					Subtype: "personal-welfare",
					Status:  shared.CaseStatusTypeRegistered,
				},
			},
		},
	},
}

func TestGetCreateObjectionsTemplate(t *testing.T) {
	tests := []struct {
		name           string
		errorReturned  error
		expectedStatus int
	}{
		{
			name:          "add objections template successfully loads",
			errorReturned: nil,
		},
		{
			name:          "returns template error",
			errorReturned: errExample,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockCreateObjectionClient{}
			client.
				On("CaseSummary", mock.Anything, "M-9898-9898-9898").
				Return(testCreateObjectionsCaseSummary, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createObjectionData{
					Title:   "Add Objection",
					CaseUID: "M-9898-9898-9898",
					LinkedLpas: []sirius.SiriusData{
						{
							UID:     "M-9898-9898-9898",
							ID:      676,
							Subtype: "personal-welfare",
							Status:  shared.CaseStatusTypeDraft,
							LinkedCases: []sirius.SiriusData{
								{
									UID:     "M-9999-9999-9999",
									Subtype: "personal-welfare",
									Status:  shared.CaseStatusTypeInProgress,
								},
								{
									UID:     "M-8888-8888-8888",
									Subtype: "personal-welfare",
									Status:  shared.CaseStatusTypeRegistered,
								},
							},
						},
						{
							UID:     "M-9999-9999-9999",
							Subtype: "personal-welfare",
							Status:  shared.CaseStatusTypeInProgress,
						},
					},
				}).
				Return(tc.errorReturned)

			r, _ := http.NewRequest(http.MethodGet, "add-objection?uid=M-9898-9898-9898", nil)
			w := httptest.NewRecorder()

			err := CreateObjection(client, template.Func)(w, r)

			if tc.errorReturned != nil {
				assert.Equal(t, tc.errorReturned, err)
			} else {
				resp := w.Result()
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}

			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateObjectionWhenCaseSummaryErrors(t *testing.T) {
	client := &mockCreateObjectionClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9898-9898-9898").
		Return(sirius.CaseSummary{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "add-objection?uid=M-9898-9898-9898", nil)
	w := httptest.NewRecorder()

	err := CreateObjection(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateObjection(t *testing.T) {
	tests := []struct {
		name          string
		apiError      error
		expectedError error
	}{
		{
			name:          "Add objection post form success",
			apiError:      nil,
			expectedError: RedirectError("/lpa/M-9898-9898-9898"),
		},
		{
			name:          "Add objection returns an API failure",
			apiError:      errExample,
			expectedError: errExample,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockCreateObjectionClient{}
			client.
				On("CaseSummary", mock.Anything, "M-9898-9898-9898").
				Return(testCreateObjectionsCaseSummary, nil)
			client.
				On("CreateObjection", mock.Anything, sirius.ObjectionRequest{
					LpaUids:       []string{"M-9898-9898-9898", "M-9999-9999-9999"},
					ReceivedDate:  "2025-01-01",
					ObjectionType: "factual",
					Notes:         "Test",
				}).
				Return(tc.apiError)

			template := &mockTemplate{}

			form := url.Values{
				"lpaUids":            {"M-9898-9898-9898", "M-9999-9999-9999"},
				"receivedDate.day":   {"1"},
				"receivedDate.month": {"1"},
				"receivedDate.year":  {"2025"},
				"objectionType":      {"factual"},
				"notes":              {"Test"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/create-objection?uid=M-9898-9898-9898", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateObjection(client, template.Func)(w, r)
			assert.Equal(t, tc.expectedError, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateObjectionWhenValidationError(t *testing.T) {
	client := &mockCreateObjectionClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9898-9898-9898").
		Return(testCreateObjectionsCaseSummary, nil)
	client.
		On("CreateObjection", mock.Anything, sirius.ObjectionRequest{
			LpaUids:      []string{"M-9898-9898-9898"},
			ReceivedDate: "2025-01-01",
			Notes:        "Test",
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"objectionType": {"required": "Value required and can't be empty"},
		}})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			mock.MatchedBy(func(data createObjectionData) bool {
				return data.Error.Field["objectionType"]["required"] == "Value required and can't be empty"
			}),
		).
		Return(nil)

	form := url.Values{
		"lpaUids":            {"M-9898-9898-9898"},
		"receivedDate.day":   {"1"},
		"receivedDate.month": {"1"},
		"receivedDate.year":  {"2025"},
		"objectionType":      {""},
		"notes":              {"Test"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/create-objection?uid=M-9898-9898-9898", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateObjection(client, template.Func)(w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
