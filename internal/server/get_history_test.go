package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGetHistoryClient struct {
	mock.Mock
}

func (m *mockGetHistoryClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockGetHistoryClient) GetEvents(ctx sirius.Context, donorId int, caseId int) (any, error) {
	args := m.Called(ctx, donorId, caseId)
	return args.Get(0), args.Error(1)
}

func (m *mockGetHistoryClient) GetCombinedEvents(ctx sirius.Context, uid string) (sirius.APIEvent, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.APIEvent), args.Error(1)
}

func TestGetHistorySuccessForDigitalLpa(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9999",
			SiriusData: sirius.SiriusData{
				ID:      12,
				Subtype: "hw",
				Donor: sirius.Donor{
					ID: 8,
				},
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: shared.ParseCaseStatusType("processing"), // Non-empty status indicates digital LPA
			},
		},
	}

	client := &mockGetHistoryClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9999").
		Return(caseSummary, nil)
	client.
		On("GetCombinedEvents", mock.Anything, "M-9876-9876-9999").
		Return(sirius.APIEvent{{
			ChangeSet:     nil,
			CreatedOn:     "08/08/1999",
			Entity:        nil,
			ID:            2,
			Source:        "mustard",
			SourceType:    "french",
			Type:          "LPA",
			User:          sirius.EventUser{DisplayName: "Bear Ghost"},
			UUID:          "654de60e-446d-4b2f-b2a7-321bf03b37df",
			FormattedUUID: "",
			Applied:       "08/08/1999",
			DateTime:      "08/08/1999",
		}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getHistory{
			CaseSummary: caseSummary,
			EventData: sirius.APIEvent{
				sirius.Event{
					ChangeSet:     nil,
					CreatedOn:     "08/08/1999",
					Entity:        nil,
					ID:            2,
					Source:        "mustard",
					SourceType:    "french",
					Type:          "LPA",
					User:          sirius.EventUser{DisplayName: "Bear Ghost"},
					UUID:          "654de60e-446d-4b2f-b2a7-321bf03b37df",
					FormattedUUID: "MVG6MDSE",
					Applied:       "08/08/1999",
					DateTime:      "08/08/1999",
				},
			},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/history", GetHistory(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9999/history", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetHistoryWhenFailureOnGetCombinedEventsForDigitalLpa(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9999",
			SiriusData: sirius.SiriusData{
				ID:      12,
				Subtype: "hw",
				Donor: sirius.Donor{
					ID: 8,
				},
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: shared.ParseCaseStatusType("processing"),
			},
		},
	}

	client := &mockGetHistoryClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9999").
		Return(caseSummary, nil)
	client.
		On("GetCombinedEvents", mock.Anything, "M-9876-9876-9999").
		Return(sirius.APIEvent{}, errExample)

	server := newMockServer("/lpa/{uid}/history", GetHistory(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9999/history", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetHistoryWhenFailureOnGetCaseSummary(t *testing.T) {
	client := &mockGetHistoryClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9999").
		Return(sirius.CaseSummary{}, errExample)

	template := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/history", GetHistory(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9999/history", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestLPAEventIDFromUUIDReturnFormattedUUID(t *testing.T) {
	result, _ := LPAEventIDFromUUID("654de60e-446d-4b2f-b2a7-321bf03b37df")
	assert.Equal(t, "MVG6MDSE", result)
}

func TestLPAEventIDFromUUIDReturnError(t *testing.T) {
	tests := []struct {
		name        string
		uuidStr     string
		expectErr   bool
		expectedLen int
	}{
		{
			name:      "Invalid hex characters",
			uuidStr:   "550e8400-e29b-41d4-a716-44665544ZZZZ",
			expectErr: true,
		},
		{
			name:      "Too short UUID decode fails",
			uuidStr:   "1234",
			expectErr: true,
		},
		{
			name:      "Empty string",
			uuidStr:   "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LPAEventIDFromUUID(tt.uuidStr)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil (output: %v)", got)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(got) != tt.expectedLen {
				t.Errorf("expected Base32 length %d, got %d (%s)", tt.expectedLen, len(got), got)
			}
		})
	}
}
