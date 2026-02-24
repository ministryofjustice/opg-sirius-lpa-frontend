package server

import (
	"net/http"
	"net/url"
	"strings"
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
						Type:       "INS",
						Hash:       "a",
						SourceType: shared.LpaEventSourceTypeLpa,
						OwningCase: sirius.OwningCase{
							ID:       1,
							CaseType: "LPA",
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
							SourceType: shared.LpaEventSourceTypeLpa,
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
						Type:       "INS",
						Hash:       "a",
						SourceType: shared.LpaEventSourceTypeLpa,
						OwningCase: sirius.OwningCase{
							ID:       1,
							CaseType: "LPA",
						},
					},
					{
						ID:         98,
						CreatedOn:  "2024-02-01T12:00:00Z",
						Type:       "INS",
						Hash:       "b",
						SourceType: shared.LpaEventSourceTypeLpa,
						OwningCase: sirius.OwningCase{
							ID:       2,
							CaseType: "LPA",
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
							SourceType: shared.LpaEventSourceTypeLpa,
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
			client.
				On("GetEvents", mock.Anything, "123", tc.caseIdsArgs, []string{}, "desc").
				Return(tc.response, err)

			eventsWithContext := make([]LpaEventWithContext, len(tc.response.Events))
			for i, event := range tc.response.Events {
				eventsWithContext[i] = LpaEventWithContext{
					LpaEvent: event,
					DonorID:  "123",
				}
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, getLpaHistory{
					DonorID:             "123",
					Events:              eventsWithContext,
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

func TestPostFiltersLpaHistory(t *testing.T) {

	unfilteredResponse := sirius.LpaEventsResponse{
		Events: []sirius.LpaEvent{
			{
				ID:         103,
				CreatedOn:  "2025-03-03T12:00:00Z",
				Type:       "INS",
				Hash:       "E",
				SourceType: shared.LpaEventSourceTypeAddress,
				OwningCase: sirius.OwningCase{
					ID:       3,
					CaseType: "LPA",
				},
			},
			{
				ID:         102,
				CreatedOn:  "2024-03-02T12:00:00Z",
				Type:       "UPD",
				Hash:       "D",
				SourceType: shared.LpaEventSourceTypePayment,
				OwningCase: sirius.OwningCase{
					ID:       4,
					CaseType: "LPA",
				},
			},
			{
				ID:         101,
				CreatedOn:  "2023-03-01T12:00:00Z",
				Type:       "INS",
				Hash:       "C",
				SourceType: shared.LpaEventSourceTypeLpa,
				OwningCase: sirius.OwningCase{
					ID:       3,
					CaseType: "LPA",
				},
			},
		},
		Limit: 999,
		Total: 3,
		Pages: sirius.Pages{},
		Metadata: sirius.EventMetaData{
			CaseIds: nil,
			SourceTypes: []sirius.SourceType{
				{SourceType: shared.LpaEventSourceTypeLpa, Total: 1},
				{SourceType: shared.LpaEventSourceTypePayment, Total: 1},
				{SourceType: shared.LpaEventSourceTypeAddress, Total: 1},
			},
		},
	}

	filteredResponse := sirius.LpaEventsResponse{
		Events: []sirius.LpaEvent{
			{
				ID:         101,
				CreatedOn:  "2024-03-01T12:00:00Z",
				Type:       "INS",
				Hash:       "c",
				SourceType: shared.LpaEventSourceTypeLpa,
				OwningCase: sirius.OwningCase{
					ID:       3,
					CaseType: "LPA",
				},
			},
			{
				ID:         102,
				CreatedOn:  "2024-03-02T12:00:00Z",
				Type:       "UPD",
				Hash:       "d",
				SourceType: shared.LpaEventSourceTypePayment,
				OwningCase: sirius.OwningCase{
					ID:       4,
					CaseType: "LPA",
				},
			},
		},
		Limit: 999,
		Total: 2,
		Pages: sirius.Pages{},
		Metadata: sirius.EventMetaData{
			CaseIds: nil,
			SourceTypes: []sirius.SourceType{
				{SourceType: shared.LpaEventSourceTypeLpa, Total: 1},
				{SourceType: shared.LpaEventSourceTypePayment, Total: 1},
			},
		},
	}

	client := &mockGetLpaHistory{}
	client.On("GetEvents", mock.Anything, "123", []string(nil), []string{}, "desc").
		Return(unfilteredResponse, nil)
	client.
		On("GetEvents", mock.Anything, "123", []string(nil), []string{"Lpa", "Payment"}, "asc").
		Return(filteredResponse, nil)

	eventsWithContext := make([]LpaEventWithContext, len(filteredResponse.Events))
	for i, event := range filteredResponse.Events {
		eventsWithContext[i] = LpaEventWithContext{
			LpaEvent: event,
			DonorID:  "123",
		}
	}

	template := &mockTemplate{}
	template.On("Func", mock.Anything, getLpaHistory{
		DonorID:             "123",
		Events:              eventsWithContext,
		EventFilterData:     unfilteredResponse.Metadata.SourceTypes,
		TotalEvents:         unfilteredResponse.Total,
		TotalFilteredEvents: filteredResponse.Total,
		IsFiltered:          true,
		Form: FilterLpaEventsForm{
			Types: []string{"Lpa", "Payment"},
			Sort:  "asc",
		},
	}).Return(nil)

	server := newMockServer("/lpa-api/v1/persons/{donorId}/events", GetLpaHistory(client, template.Func))

	form := url.Values{
		"sort": {"asc"},
		"type": {"Lpa", "Payment"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/lpa-api/v1/persons/123/events", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	r.PostForm = form
	_, err := server.serve(r)

	assert.Equal(t, nil, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostFiltersLpaHistoryWhenFailureOnGetEvents(t *testing.T) {
	client := &mockGetLpaHistory{}

	client.On("GetEvents", mock.Anything, "123", []string(nil), []string{}, "desc").
		Return(sirius.LpaEventsResponse{}, nil)
	client.
		On("GetEvents", mock.Anything, "123", []string(nil), []string{"Lpa", "Payment"}, "asc").
		Return(sirius.LpaEventsResponse{}, errExample)

	server := newMockServer("/lpa-api/v1/persons/{donorId}/events", GetLpaHistory(client, nil))

	form := url.Values{
		"sort": {"asc"},
		"type": {"Lpa", "Payment"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/lpa-api/v1/persons/123/events", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}
