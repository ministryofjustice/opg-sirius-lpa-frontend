package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type mockObjectionOutcomeClient struct {
	mock.Mock
}

func (m *mockObjectionOutcomeClient) GetObjection(ctx sirius.Context, objectionId string) (sirius.Objection, error) {
	args := m.Called(ctx, objectionId)
	return args.Get(0).(sirius.Objection), args.Error(1)
}

var testResolvedObjection = sirius.Objection{
	ID:            12,
	Notes:         "Test",
	ObjectionType: "factual",
	ReceivedDate:  "2025-03-12",
	LpaUids:       []string{"M-3333-3333-3333"},
	Resolutions: []sirius.ObjectionResolution{{
		Uid:             "M-3333-3333-3333",
		Resolution:      "Upheld",
		ResolutionNotes: "Test",
		ResolutionDate:  "2025-04-04",
	}},
}

func TestGetObjectionOutcomeTemplate(t *testing.T) {
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
			client := &mockObjectionOutcomeClient{}
			client.
				On("GetObjection", mock.Anything, "12").
				Return(testResolvedObjection, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, objectionOutcomeData{
					Objection:  testResolvedObjection,
					Resolution: testResolvedObjection.Resolutions[0],
				}).
				Return(tc.errorReturned)

			server := newMockServer("/lpa/{uid}/objection/{id}/outcome", ObjectionOutcome(client, template.Func))

			r, _ := http.NewRequest(http.MethodGet, "/lpa/M-3333-3333-3333/objection/12/outcome", nil)
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

func TestGetObjectionOutcomeWhenGetObjectionErrors(t *testing.T) {
	client := &mockResolveObjectionClient{}
	client.
		On("GetObjection", mock.Anything, "12").
		Return(sirius.Objection{}, errExample)

	server := newMockServer("/lpa/{uid}/objection/{id}/resolve", ResolveObjection(client, nil))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/M-3333-3333-3333/objection/12/resolve", nil)
	_, err := server.serve(r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}
