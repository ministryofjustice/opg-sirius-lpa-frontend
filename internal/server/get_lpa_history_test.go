package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGetLpaHistory struct {
	mock.Mock
}

func (m *mockGetLpaHistory) GetEvents(ctx sirius.Context, donorId string, caseIds []string, sourceTypes []string, sortBy string) (sirius.LpaEventsResponse, error) {
	args := m.Called(ctx, donorId, caseIds, sourceTypes, sortBy)
	return args.Get(0).(sirius.LpaEventsResponse), args.Error(1)
}

func TestGetLpaHistory(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		caseIdsArgs []string
		response    sirius.LpaEventsResponse
	}{
		{
			name:        "all history by default when no case ID provided",
			path:        "/lpa-api/v1/persons/123/events",
			caseIdsArgs: []string(nil),
			response: sirius.LpaEventsResponse{
				Events: []sirius.LpaEvent{
					{
						ID:         99,
						CreatedOn:  "2024-01-01T12:00:00Z",
						Type:       "Save",
						Hash:       "a",
						SourceType: shared.LpaEventSourceTypeLpa,
						OwningCase: sirius.OwningCase{
							ID:       1,
							CaseType: "Lpa",
						},
					},
				},
				Limit: 999,
				Total: 1,
				Pages: sirius.Pages{},
				Metadata: sirius.EventMetaData{
					CaseIds: nil,
					SourceTypes: []sirius.SourceType{
						{
							SourceType: "Lpa",
							Total:      1,
						},
					},
				},
			},
		},
		{
			name:        "case histories when multiple IDs provided",
			path:        "/lpa-api/v1/persons/123/events?id[]=1&id[]=2",
			caseIdsArgs: []string{"1", "2"},
			response: sirius.LpaEventsResponse{
				Events: []sirius.LpaEvent{
					{
						ID:         99,
						CreatedOn:  "2024-01-01T12:00:00Z",
						Type:       "Save",
						Hash:       "a",
						SourceType: shared.LpaEventSourceTypeLpa,
						OwningCase: sirius.OwningCase{
							ID:       1,
							CaseType: "Lpa",
						},
					},
					{
						ID:         98,
						CreatedOn:  "2024-02-01T12:00:00Z",
						Type:       "Save",
						Hash:       "b",
						SourceType: shared.LpaEventSourceTypeLpa,
						OwningCase: sirius.OwningCase{
							ID:       2,
							CaseType: "Lpa",
						},
					},
				},
				Limit: 999,
				Total: 2,
				Pages: sirius.Pages{},
				Metadata: sirius.EventMetaData{
					CaseIds: nil,
					SourceTypes: []sirius.SourceType{
						{
							SourceType: "Lpa",
							Total:      2,
						},
					},
				},
			},
		},
	}

	client := &mockGetLpaHistory{}
	var err error

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client.On("GetEvents", mock.Anything, "123", tc.caseIdsArgs, []string{}, "desc").Return(tc.response, err)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, getLpaHistory{
					Events:              tc.response.Events,
					EventFilterData:     tc.response.Metadata.SourceTypes,
					TotalEvents:         tc.response.Total,
					TotalFilteredEvents: 0,
					Form: FilterLpaEventsForm{
						Sort: "desc",
					},
				}).
				Return(nil)

			server := newMockServer("/lpa-api/v1/persons/{donorId}/events", GetLpaHistory(client, template.Func))

			req, _ := http.NewRequest(http.MethodGet, tc.path, nil)
			resp, err := server.serve(req)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetLpaHistoryWhenFailureOnGetEvents(t *testing.T) {
	client := &mockGetLpaHistory{}
	client.On("GetEvents", mock.Anything, "123", []string(nil), []string{}, "desc").Return(sirius.LpaEventsResponse{}, errExample)

	server := newMockServer("/lpa-api/v1/persons/{donorId}/events", GetLpaHistory(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/lpa-api/v1/persons/123/events", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}
