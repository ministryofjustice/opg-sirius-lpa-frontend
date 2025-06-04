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

type mockResolveObjectionClient struct {
	mock.Mock
}

func (m *mockResolveObjectionClient) ResolveObjection(ctx sirius.Context, objectionID, lpaUid string, resolution sirius.ResolutionRequest) error {
	args := m.Called(ctx, objectionID, lpaUid, resolution)
	return args.Error(0)
}

func (m *mockResolveObjectionClient) GetObjection(ctx sirius.Context, objectionId string) (sirius.Objection, error) {
	args := m.Called(ctx, objectionId)
	return args.Get(0).(sirius.Objection), args.Error(1)
}

var testObjection1 = sirius.Objection{
	ID:            6,
	Notes:         "Test",
	ObjectionType: "factual",
	ReceivedDate:  "2025-03-12",
	LpaUids:       []string{"M-4444-4444-4444"},
}

func TestGetResolveObjectionsTemplate(t *testing.T) {
	tests := []struct {
		name          string
		errorReturned error
	}{
		{
			name:          "resolve objections template successfully loads",
			errorReturned: nil,
		},
		{
			name:          "returns template error",
			errorReturned: errExample,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockResolveObjectionClient{}
			client.
				On("GetObjection", mock.Anything, "6").
				Return(testObjection1, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, resolveObjectionData{
					CaseUID:     "M-4444-4444-4444",
					ObjectionId: "6",
					Objection:   testObjection1,
					LpaUids:     testObjection1.LpaUids,
					Form: formResolveObjection{
						Resolution:      []string{""},
						ResolutionNotes: []string{""},
					},
				}).
				Return(tc.errorReturned)

			server := newMockServer("/lpa/{uid}/objection/{id}/resolve", ResolveObjection(client, template.Func))

			r, _ := http.NewRequest(http.MethodGet, "/lpa/M-4444-4444-4444/objection/6/resolve", nil)
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

func TestGetResolveObjectionWhenGetObjectionErrors(t *testing.T) {
	client := &mockResolveObjectionClient{}
	client.
		On("GetObjection", mock.Anything, "6").
		Return(sirius.Objection{}, errExample)

	server := newMockServer("/lpa/{uid}/objection/{id}/resolve", ResolveObjection(client, nil))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/M-4444-4444-4444/objection/6/resolve", nil)
	_, err := server.serve(r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostResolveObjection(t *testing.T) {
	tests := []struct {
		name          string
		apiError      error
		expectedError error
	}{
		{
			name:          "Resolve objection post form success",
			apiError:      nil,
			expectedError: RedirectError("/lpa/M-4444-4444-4444"),
		},
		{
			name:          "Resolve objection returns an API failure",
			apiError:      errExample,
			expectedError: errExample,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockResolveObjectionClient{}
			client.
				On("GetObjection", mock.Anything, "6").
				Return(testObjection1, nil)
			client.
				On("ResolveObjection", mock.Anything, "6", "M-4444-4444-4444", sirius.ResolutionRequest{
					Resolution: "upheld",
					Notes:      "Test",
				}).
				Return(tc.apiError)

			template := &mockTemplate{}

			form := url.Values{
				"caseUid":           {"M-4444-4444-4444"},
				"resolution-0":      {"upheld"},
				"resolutionNotes-0": {"Test"},
			}

			server := newMockServer("/lpa/{uid}/objection/{id}/resolve", ResolveObjection(client, template.Func))

			r, _ := http.NewRequest(http.MethodPost, "/lpa/M-4444-4444-4444/objection/6/resolve", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			_, err := server.serve(r)

			assert.Equal(t, tc.expectedError, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostResolveObjectionWhenValidationError(t *testing.T) {
	client := &mockResolveObjectionClient{}
	client.
		On("GetObjection", mock.Anything, "6").
		Return(testObjection1, nil)
	client.
		On("ResolveObjection", mock.Anything, "6", "M-4444-4444-4444", sirius.ResolutionRequest{
			Notes: "Test",
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"resolution": {"required": "Value required and can't be empty"},
		}})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			mock.MatchedBy(func(data resolveObjectionData) bool {
				return data.Error.Field["resolution"]["required"] == "Value required and can't be empty"
			}),
		).
		Return(nil)

	form := url.Values{
		"caseUid":           {"M-4444-4444-4444"},
		"resolution-0":      {""},
		"resolutionNotes-0": {"Test"},
	}

	server := newMockServer("/lpa/{uid}/objection/{id}/resolve", ResolveObjection(client, template.Func))
	r, _ := http.NewRequest(http.MethodPost, "/lpa/M-4444-4444-4444/objection/6/resolve", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
