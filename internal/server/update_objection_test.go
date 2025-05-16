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

type mockUpdateObjectionClient struct {
	mock.Mock
}

func (m *mockUpdateObjectionClient) UpdateObjection(ctx sirius.Context, objectionID string, objection sirius.ObjectionRequest) error {
	args := m.Called(ctx, objectionID, objection)
	return args.Error(0)
}

func (m *mockUpdateObjectionClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockUpdateObjectionClient) GetObjection(ctx sirius.Context, objectionId string) (sirius.Objection, error) {
	args := m.Called(ctx, objectionId)
	return args.Get(0).(sirius.Objection), args.Error(1)
}

var testUpdateObjectionsCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-7777-8888-9999",
		SiriusData: sirius.SiriusData{
			ID:      676,
			UID:     "M-7777-8888-9999",
			Subtype: "personal-welfare",
			Status:  "Draft",
			LinkedCases: []sirius.SiriusData{
				{
					UID:     "M-9999-9999-9999",
					Subtype: "personal-welfare",
					Status:  "In progress",
				},
				{
					UID:     "M-8888-8888-8888",
					Subtype: "personal-welfare",
					Status:  "Registered",
				},
			},
		},
	},
	Objections: []sirius.ObjectionForCase{
		testObjection,
	},
}

var testObjection = sirius.ObjectionForCase{
	ID:            3,
	Notes:         "Test",
	ObjectionType: "factual",
	ReceivedDate:  "2025-03-12",
	LpaUids:       []string{"M-7777-8888-9999"},
}

var testObjection2 = sirius.Objection{
	ID:            3,
	Notes:         "Test",
	ObjectionType: "factual",
	ReceivedDate:  "2025-03-12",
	LpaUids:       []string{"M-7777-8888-9999"},
	Resolutions: []sirius.ObjectionResolution{
		{
			Resolution:      "not upheld",
			ResolutionNotes: "Everything is fine",
			ResolutionDate:  "2025-01-01",
		},
	},
}

func TestGetUpdateObjectionsTemplate(t *testing.T) {
	tests := []struct {
		name           string
		errorReturned  error
		expectedStatus int
	}{
		{
			name:          "update objections template successfully loads",
			errorReturned: nil,
		},
		{
			name:          "returns template error",
			errorReturned: errExample,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockUpdateObjectionClient{}
			client.
				On("CaseSummary", mock.Anything, "M-7777-8888-9999").
				Return(testUpdateObjectionsCaseSummary, nil)
			client.
				On("GetObjection", mock.Anything, "3").
				Return(testObjection2, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, updateObjectionData{
					Title:   "Update Objection",
					CaseUID: "M-7777-8888-9999",
					Form: formObjection{
						LpaUids:       []string{"M-7777-8888-9999"},
						ReceivedDate:  dob{12, 03, 2025},
						ObjectionType: "factual",
						Notes:         "Test",
					},
					LinkedLpas: []sirius.SiriusData{
						{
							UID:     "M-7777-8888-9999",
							ID:      676,
							Subtype: "personal-welfare",
							Status:  "Draft",
							LinkedCases: []sirius.SiriusData{
								{
									UID:     "M-9999-9999-9999",
									Subtype: "personal-welfare",
									Status:  "In progress",
								},
								{
									UID:     "M-8888-8888-8888",
									Subtype: "personal-welfare",
									Status:  "Registered",
								},
							},
						},
						{
							UID:     "M-9999-9999-9999",
							Subtype: "personal-welfare",
							Status:  "In progress",
						},
					},
				}).
				Return(tc.errorReturned)

			server := newMockServer("/lpa/{uid}/objection/{id}", UpdateObjection(client, template.Func, template.Func))

			r, _ := http.NewRequest(http.MethodGet, "/lpa/M-7777-8888-9999/objection/3", nil)
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

func TestGetUpdateObjectionWhenCaseSummaryErrors(t *testing.T) {
	client := &mockUpdateObjectionClient{}
	client.
		On("CaseSummary", mock.Anything, "M-7777-8888-9999").
		Return(sirius.CaseSummary{}, errExample)
	client.
		On("GetObjection", mock.Anything, "3").
		Return(testObjection2, nil)

	server := newMockServer("/lpa/{uid}/objection/{id}", UpdateObjection(client, nil, nil))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/M-7777-8888-9999/objection/3", nil)
	_, err := server.serve(r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetUpdateObjectionWhenGetObjectionErrors(t *testing.T) {
	client := &mockUpdateObjectionClient{}
	client.
		On("CaseSummary", mock.Anything, "M-7777-8888-9999").
		Return(testUpdateObjectionsCaseSummary, nil)
	client.
		On("GetObjection", mock.Anything, "3").
		Return(sirius.Objection{}, errExample)

	server := newMockServer("/lpa/{uid}/objection/{id}", UpdateObjection(client, nil, nil))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/M-7777-8888-9999/objection/3", nil)
	_, err := server.serve(r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostUpdateObjectionToAPI(t *testing.T) {
	tests := []struct {
		name          string
		apiError      error
		expectedError error
	}{
		{
			name:          "Update objection post form success",
			apiError:      nil,
			expectedError: RedirectError("/lpa/M-7777-8888-9999"),
		},
		{
			name:          "Update objection returns an API failure",
			apiError:      errExample,
			expectedError: errExample,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockUpdateObjectionClient{}
			client.
				On("CaseSummary", mock.Anything, "M-7777-8888-9999").
				Return(testUpdateObjectionsCaseSummary, nil)
			client.
				On("GetObjection", mock.Anything, "3").
				Return(testObjection2, nil)
			client.
				On("UpdateObjection", mock.Anything, "3", sirius.ObjectionRequest{
					LpaUids:       []string{"M-7777-8888-9999", "M-9999-9999-9999"},
					ReceivedDate:  "2025-01-01",
					ObjectionType: "prescribed",
					Notes:         "Test",
				}).
				Return(tc.apiError)

			template := &mockTemplate{}

			form := url.Values{
				"lpaUids":            {"M-7777-8888-9999", "M-9999-9999-9999"},
				"receivedDate.day":   {"1"},
				"receivedDate.month": {"1"},
				"receivedDate.year":  {"2025"},
				"objectionType":      {"prescribed"},
				"notes":              {"Test"},
				"step":               {"confirm"},
			}

			server := newMockServer("/lpa/{uid}/objection/{id}", UpdateObjection(client, template.Func, template.Func))

			r, _ := http.NewRequest(http.MethodPost, "/lpa/M-7777-8888-9999/objection/3", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			_, err := server.serve(r)

			assert.Equal(t, tc.expectedError, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostUpdateObjectionWhenValidationError(t *testing.T) {
	client := &mockUpdateObjectionClient{}
	client.
		On("CaseSummary", mock.Anything, "M-7777-8888-9999").
		Return(testUpdateObjectionsCaseSummary, nil)
	client.
		On("GetObjection", mock.Anything, "3").
		Return(testObjection2, nil)
	client.
		On("UpdateObjection", mock.Anything, "3", sirius.ObjectionRequest{
			LpaUids:      []string{"M-7777-8888-9999"},
			ReceivedDate: "2025-01-01",
			Notes:        "Test",
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"objectionType": {"required": "Value required and can't be empty"},
		}})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			mock.MatchedBy(func(data updateObjectionData) bool {
				return data.Error.Field["objectionType"]["required"] == "Value required and can't be empty"
			}),
		).
		Return(nil)

	form := url.Values{
		"lpaUids":            {"M-7777-8888-9999"},
		"receivedDate.day":   {"1"},
		"receivedDate.month": {"1"},
		"receivedDate.year":  {"2025"},
		"objectionType":      {""},
		"notes":              {"Test"},
		"step":               {"confirm"},
	}

	server := newMockServer("/lpa/{uid}/objection/{id}", UpdateObjection(client, template.Func, template.Func))
	r, _ := http.NewRequest(http.MethodPost, "/lpa/M-7777-8888-9999/objection/3", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateObjectionToConfirmScreen(t *testing.T) {
	client := &mockUpdateObjectionClient{}
	client.
		On("CaseSummary", mock.Anything, "M-7777-8888-9999").
		Return(testUpdateObjectionsCaseSummary, nil)
	client.
		On("GetObjection", mock.Anything, "3").
		Return(testObjection2, nil)

	confirmTemplate := &mockTemplate{}
	confirmTemplate.
		On("Func", mock.Anything, struct {
			XSRFToken    string
			CaseUID      string
			ObjectionID  string
			ReceivedDate sirius.DateString
			Form         formObjection
		}{
			XSRFToken:    "",
			CaseUID:      "M-7777-8888-9999",
			ObjectionID:  "3",
			ReceivedDate: "2025-01-01",
			Form: formObjection{
				LpaUids:       []string{"M-7777-8888-9999", "M-9999-9999-9999"},
				ReceivedDate:  dob{1, 1, 2025},
				ObjectionType: "prescribed",
				Notes:         "Test",
			},
		}).
		Return(nil).Once()

	formTemplate := &mockTemplate{}

	form := url.Values{
		"lpaUids":            {"M-7777-8888-9999", "M-9999-9999-9999"},
		"receivedDate.day":   {"1"},
		"receivedDate.month": {"1"},
		"receivedDate.year":  {"2025"},
		"objectionType":      {"prescribed"},
		"notes":              {"Test"},
		"step":               {"review"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/lpa/M-7777-8888-9999/objection/3", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)

	server := newMockServer("/lpa/{uid}/objection/{id}", UpdateObjection(client, formTemplate.Func, confirmTemplate.Func))

	_, err := server.serve(r)
	assert.NoError(t, err)

	mock.AssertExpectationsForObjects(t, client, confirmTemplate)
}
