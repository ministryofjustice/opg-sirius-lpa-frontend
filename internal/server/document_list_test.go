package server

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDocumentListClient struct {
	mock.Mock
}

func (m *mockDocumentListClient) CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.Case), args.Error(1)
}

func (m *mockDocumentListClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockDocumentListClient) GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error) {
	args := m.Called(ctx, personID, caseIDs)
	return args.Get(0).(sirius.DocumentList), args.Error(1)
}

func (m *mockDocumentListClient) DownloadMultiple(ctx sirius.Context, docIDs []string) (*http.Response, error) {
	args := m.Called(ctx, docIDs)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *mockDocumentListClient) GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.Permissions), args.Error(1)
}

func (m *mockDocumentListClient) GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error) {
	args := m.Called(ctx, caseType, caseId)
	return args.Get(0).(sirius.DocumentDraftCount), args.Error(1)
}

func (m *mockDocumentListClient) PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.PersonReference), args.Error(1)
}

var singleDocumentList = sirius.DocumentList{
	Limit: 999,
	Pages: sirius.Pages{
		Current: 1,
		Total:   1,
	},
	Total: 1,
	Documents: []sirius.Document{
		{
			ID:                  639,
			UUID:                "31e6f4c2-5f8b-47c3-bc98-64b47c938e52",
			Type:                "Save",
			FriendlyDescription: "Letter",
			CreatedDate:         "25/07/2022 14:17:13",
			Direction:           shared.DocumentDirectionOut,
			FileName:            "LP-NA-3A.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-1234-0000",
					SubType:  "pfa",
					CaseType: "LPA",
				},
			},
			SystemType: "LP-NA-3A",
		},
	},
}

var twoCasesDocumentList = sirius.DocumentList{
	Limit: 999,
	Pages: sirius.Pages{
		Current: 1,
		Total:   1,
	},
	Total: 2,
	Documents: []sirius.Document{
		{
			ID:                  443,
			UUID:                "c8e3a1df-7b9b-4d45-94d9-2b8fc0d9e0fd",
			Type:                "LPA",
			FriendlyDescription: "LP1H - Health Instrument",
			CreatedDate:         "01/06/2022 15:39:01",
			Direction:           shared.DocumentDirectionIn,
			FileName:            "LP1H.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-9876-0000",
					SubType:  "hw",
					CaseType: "LPA",
				},
			},
			SubType: "hw",
		},
		{
			ID:                  639,
			UUID:                "31e6f4c2-5f8b-47c3-bc98-64b47c938e52",
			Type:                "Save",
			FriendlyDescription: "Letter",
			CreatedDate:         "25/07/2022 14:17:13",
			Direction:           shared.DocumentDirectionEmpty,
			FileName:            "LP-NA-3A.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-1234-0000",
					SubType:  "pfa",
					CaseType: "LPA",
				},
			},
			SystemType: "LP-NA-3A",
		},
	},
}

var allDocumentList = sirius.DocumentList{
	Limit: 999,
	Pages: sirius.Pages{
		Current: 1,
		Total:   1,
	},
	Total: 3,
	Documents: []sirius.Document{
		{
			ID:                  443,
			UUID:                "c8e3a1df-7b9b-4d45-94d9-2b8fc0d9e0fd",
			Type:                "LPA",
			FriendlyDescription: "LP1H - Health Instrument",
			CreatedDate:         "01/06/2022 15:39:01",
			Direction:           shared.DocumentDirectionIn,
			FileName:            "LP1H.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-9876-0000",
					SubType:  "hw",
					CaseType: "LPA",
				},
			},
			SubType: "hw",
		},
		{
			ID:                  639,
			UUID:                "31e6f4c2-5f8b-47c3-bc98-64b47c938e52",
			Type:                "Save",
			FriendlyDescription: "Letter",
			CreatedDate:         "25/07/2022 14:17:13",
			Direction:           shared.DocumentDirectionOut,
			FileName:            "LP-NA-3A.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-1234-0000",
					SubType:  "pfa",
					CaseType: "LPA",
				},
			},
			SystemType: "LP-NA-3A",
		},
		{
			ID:                  928,
			UUID:                "d9e12f73-3ab2-4d24-9a63-6b0b3e49b1c5",
			Type:                "Application Related",
			FriendlyDescription: "EPA.pdf",
			CreatedDate:         "08/01/2025 10:36:41",
			Direction:           shared.DocumentDirectionIn,
			FileName:            "EPA.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-5678-0000",
					SubType:  "pfa",
					CaseType: "EPA",
				},
			},
			SubType: "pfa",
		},
	},
}

var expectedDonor = sirius.Person{
	ID:        82,
	Firstname: "Jane",
	Surname:   "Doe",
}

func TestGetDocumentList(t *testing.T) {
	cases := []sirius.Case{
		{
			ID:       1,
			CaseType: "LPA",
			SubType:  "PFA",
			UID:      "7000-1234-0000",
		},
		{
			ID:       2,
			CaseType: "LPA",
			SubType:  "HW",
			UID:      "7000-9876-0000",
		},
		{
			ID:       3,
			CaseType: "EPA",
			SubType:  "PFA",
			UID:      "7000-5678-0000",
		},
	}

	tests := []struct {
		name                      string
		cases                     []sirius.Case
		documentList              sirius.DocumentList
		expectedMultiple          bool
		expectedCases             []sirius.Case
		caseIDs                   []string
		caseUids                  string
		selectedCaseIds           string
		path                      string
		actionPanelButtons        []ActionPanelButton
		hasV1PersonsGetPermission bool
	}{
		{
			name:                      "on person with multiple cases",
			cases:                     cases,
			documentList:              allDocumentList,
			expectedMultiple:          true,
			expectedCases:             cases,
			selectedCaseIds:           "1+2+3",
			caseIDs:                   []string(nil),
			path:                      "/donor/82/documents",
			hasV1PersonsGetPermission: true,
			actionPanelButtons: []ActionPanelButton{
				{
					Label:    "Create warning",
					URL:      "/create-warning?id=82&entity=person",
					IconName: "aw-create-warning",
					Disabled: false,
				},
				{
					Label:    "Create event",
					URL:      "/create-event?id=82&entity=person",
					IconName: "aw-new-event",
					Disabled: false,
				},
				{
					Label:    "Add complaint",
					URL:      "",
					IconName: "aw-log-complaint",
					Disabled: true,
				},
				{
					Label:    "Create document",
					URL:      "",
					IconName: "aw-new-template",
					Disabled: true,
				},
				{
					Label:    "Retrieve draft",
					URL:      "",
					IconName: "aw-new-template",
					Disabled: true,
				},
				{
					Label:    "Change status",
					URL:      "",
					IconName: "aw-change-status",
					Disabled: true,
				},
				{
					Label:    "Fees",
					URL:      "",
					IconName: "aw-fees",
					Disabled: true,
				},
				{
					Label:    "New task",
					URL:      "",
					IconName: "aw-new-task",
					Disabled: true,
				},
				{
					Label:    "Create donor",
					URL:      "/create-donor?id=82&entity=person",
					IconName: "aw-create-person",
					Disabled: false,
				},
				{
					Label:    "Edit donor",
					URL:      "/edit-donor?id=82&entity=person",
					IconName: "aw-edit-person",
					Disabled: false,
				},
				{
					Label:    "Edit dates",
					URL:      "",
					IconName: "calendar-open",
					Disabled: true,
				},
				{
					Label:    "MI reporting",
					URL:      "/mi-reporting?donorId=82",
					IconName: "aw-mi",
					Disabled: false,
				},
				{
					Label:    "Allocate Case",
					URL:      "/allocate-cases?id=1&id=2&id=3&entity=lpa",
					IconName: "aw-allocate-case",
					Disabled: false,
				},
				{
					Label:    "Link record",
					URL:      "/link-person?id=82",
					IconName: "aw-link",
					Disabled: false,
				},
				{
					Label:    "Delete relationship",
					URL:      "/delete-relationship?id=82",
					IconName: "icon-minus",
					Disabled: false,
				},
			},
		},
		{
			name:                      "on person with one case",
			cases:                     cases[:1],
			documentList:              singleDocumentList,
			expectedMultiple:          false,
			expectedCases:             cases[:1],
			caseIDs:                   []string(nil),
			selectedCaseIds:           "1",
			path:                      "/donor/82/documents",
			hasV1PersonsGetPermission: false,
			actionPanelButtons: []ActionPanelButton{
				{
					Label:    "Create warning",
					URL:      "/create-warning?id=82&entity=lpa",
					IconName: "aw-create-warning",
					Disabled: false,
				},
				{
					Label:    "Create event",
					URL:      "/create-event?id=82&entity=person",
					IconName: "aw-new-event",
					Disabled: false,
				},
				{
					Label:    "Add complaint",
					URL:      "/add-complaint?id=1&case=lpa",
					IconName: "aw-log-complaint",
					Disabled: false,
				},
				{
					Label:    "Create document",
					URL:      "/create-document?id=1&case=lpa",
					IconName: "aw-new-template",
					Disabled: false,
				},
				{
					Label:    "Retrieve draft",
					URL:      "/edit-document?id=1&case=lpa",
					IconName: "aw-new-template",
					Disabled: false,
				},
				{
					Label:    "Change status",
					URL:      "/change-status?id=1&case=lpa&donorId=82",
					IconName: "aw-change-status",
					Disabled: false,
				},
				{
					Label:    "Fees",
					URL:      "/payments/1",
					IconName: "aw-fees",
					Disabled: false,
				},
				{
					Label:    "New task",
					URL:      "/create-task?id=1&entity=lpa",
					IconName: "aw-new-task",
					Disabled: false,
				},
				{
					Label:    "Create donor",
					URL:      "/create-donor?id=82&entity=person",
					IconName: "aw-create-person",
					Disabled: false,
				},
				{
					Label:    "Edit donor",
					URL:      "/edit-donor?id=82&entity=person",
					IconName: "aw-edit-person",
					Disabled: false,
				},
				{
					Label:    "Edit dates",
					URL:      "/edit-dates?id=1&case=lpa",
					IconName: "calendar-open",
					Disabled: false,
				},
				{
					Label:    "MI reporting",
					URL:      "/mi-reporting?donorId=82",
					IconName: "aw-mi",
					Disabled: false,
				},
				{
					Label:    "Allocate Case",
					URL:      "/allocate-cases?id=1&entity=lpa",
					IconName: "aw-allocate-case",
					Disabled: false,
				},
				{
					Label:    "Link record",
					URL:      "/link-person?id=82",
					IconName: "aw-link",
					Disabled: false,
				},
				{
					Label:    "Delete relationship",
					URL:      "/delete-relationship?id=82",
					IconName: "icon-minus",
					Disabled: false,
				},
			},
		},
		{
			name:                      "one case specified",
			cases:                     cases,
			documentList:              singleDocumentList,
			expectedMultiple:          false,
			expectedCases:             []sirius.Case{cases[0]},
			caseIDs:                   []string{"1"},
			caseUids:                  "&uid[]=7000-1234-0000",
			selectedCaseIds:           "1",
			path:                      "/donor/82/documents?uid[]=7000-1234-0000",
			hasV1PersonsGetPermission: false,
			actionPanelButtons: []ActionPanelButton{
				{
					Label:    "Create warning",
					URL:      "/create-warning?id=82&entity=lpa&uid[]=7000-1234-0000",
					IconName: "aw-create-warning",
					Disabled: false,
				},
				{
					Label:    "Create event",
					URL:      "/create-event?id=82&entity=person&uid[]=7000-1234-0000",
					IconName: "aw-new-event",
					Disabled: false,
				},
				{
					Label:    "Add complaint",
					URL:      "/add-complaint?id=1&case=lpa",
					IconName: "aw-log-complaint",
					Disabled: false,
				},
				{
					Label:    "Create document",
					URL:      "/create-document?id=1&case=lpa",
					IconName: "aw-new-template",
					Disabled: false,
				},
				{
					Label:    "Retrieve draft",
					URL:      "/edit-document?id=1&case=lpa",
					IconName: "aw-new-template",
					Disabled: false,
				},
				{
					Label:    "Change status",
					URL:      "/change-status?id=1&case=lpa&donorId=82&uid[]=7000-1234-0000",
					IconName: "aw-change-status",
					Disabled: false,
				},
				{
					Label:    "Fees",
					URL:      "/payments/1",
					IconName: "aw-fees",
					Disabled: false,
				},
				{
					Label:    "New task",
					URL:      "/create-task?id=1&entity=lpa&uid[]=7000-1234-0000",
					IconName: "aw-new-task",
					Disabled: false,
				},
				{
					Label:    "Create donor",
					URL:      "/create-donor?id=82&entity=person&uid[]=7000-1234-0000",
					IconName: "aw-create-person",
					Disabled: false,
				},
				{
					Label:    "Edit donor",
					URL:      "/edit-donor?id=82&entity=person&uid[]=7000-1234-0000",
					IconName: "aw-edit-person",
					Disabled: false,
				},
				{
					Label:    "Edit dates",
					URL:      "/edit-dates?id=1&case=lpa",
					IconName: "calendar-open",
					Disabled: false,
				},
				{
					Label:    "MI reporting",
					URL:      "/mi-reporting?donorId=82&uid[]=7000-1234-0000",
					IconName: "aw-mi",
					Disabled: false,
				},
				{
					Label:    "Allocate Case",
					URL:      "/allocate-cases?id=1&entity=lpa&uid[]=7000-1234-0000",
					IconName: "aw-allocate-case",
					Disabled: false,
				},
				{
					Label:    "Link record",
					URL:      "/link-person?id=82&uid[]=7000-1234-0000",
					IconName: "aw-link",
					Disabled: false,
				},
				{
					Label:    "Delete relationship",
					URL:      "/delete-relationship?id=82&uid[]=7000-1234-0000",
					IconName: "icon-minus",
					Disabled: false,
				},
			},
		},
		{
			name:                      "multiple cases specified",
			cases:                     cases,
			documentList:              twoCasesDocumentList,
			expectedMultiple:          true,
			expectedCases:             []sirius.Case{cases[0], cases[1]},
			caseIDs:                   []string{"1", "2"},
			caseUids:                  "&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
			selectedCaseIds:           "1+2",
			path:                      "/donor/82/documents?uid[]=7000-1234-0000&uid[]=7000-9876-0000",
			hasV1PersonsGetPermission: false,
			actionPanelButtons: []ActionPanelButton{
				{
					Label:    "Create warning",
					URL:      "/create-warning?id=82&entity=person&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
					IconName: "aw-create-warning",
					Disabled: false,
				},
				{
					Label:    "Create event",
					URL:      "/create-event?id=82&entity=person&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
					IconName: "aw-new-event",
					Disabled: false,
				},
				{
					Label:    "Add complaint",
					URL:      "",
					IconName: "aw-log-complaint",
					Disabled: true,
				},
				{
					Label:    "Create document",
					URL:      "",
					IconName: "aw-new-template",
					Disabled: true,
				},
				{
					Label:    "Retrieve draft",
					URL:      "",
					IconName: "aw-new-template",
					Disabled: true,
				},
				{
					Label:    "Change status",
					URL:      "",
					IconName: "aw-change-status",
					Disabled: true,
				},
				{
					Label:    "Fees",
					URL:      "",
					IconName: "aw-fees",
					Disabled: true,
				},
				{
					Label:    "New task",
					URL:      "",
					IconName: "aw-new-task",
					Disabled: true,
				},
				{
					Label:    "Create donor",
					URL:      "/create-donor?id=82&entity=person&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
					IconName: "aw-create-person",
					Disabled: false,
				},
				{
					Label:    "Edit donor",
					URL:      "/edit-donor?id=82&entity=person&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
					IconName: "aw-edit-person",
					Disabled: false,
				},
				{
					Label:    "Edit dates",
					URL:      "",
					IconName: "calendar-open",
					Disabled: true,
				},
				{
					Label:    "MI reporting",
					URL:      "/mi-reporting?donorId=82&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
					IconName: "aw-mi",
					Disabled: false,
				},
				{
					Label:    "Allocate Case",
					URL:      "/allocate-cases?id=1&id=2&entity=lpa&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
					IconName: "aw-allocate-case",
					Disabled: false,
				},
				{
					Label:    "Link record",
					URL:      "/link-person?id=82&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
					IconName: "aw-link",
					Disabled: false,
				},
				{
					Label:    "Delete relationship",
					URL:      "/delete-relationship?id=82&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
					IconName: "icon-minus",
					Disabled: false,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			permissions := sirius.Permissions{}
			if tc.hasV1PersonsGetPermission {
				permissions = sirius.Permissions{"v1-persons": sirius.PermissionType{Permissions: []string{"GET"}}}
			}
			client := &mockDocumentListClient{}
			client.
				On("CasesByDonor", mock.Anything, 82).
				Return(tc.cases, nil)
			client.
				On("GetPersonDocuments", mock.Anything, 82, tc.caseIDs).
				Return(tc.documentList, nil)
			client.
				On("Person", mock.Anything, 82).
				Return(expectedDonor, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(permissions, nil)
			client.
				On("PersonReferences", mock.Anything, 82).
				Return([]sirius.PersonReference{{ID: 987}}, nil)

			if len(tc.expectedCases) == 1 {
				client.
					On("GetDraftCount", mock.Anything, "lpa", 1).
					Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything,
					documentPageData{
						SelectedCases:             tc.expectedCases,
						Person:                    expectedDonor,
						DocumentList:              tc.documentList,
						MultipleCasesSelected:     tc.expectedMultiple,
						DonorID:                   82,
						CaseUids:                  tc.caseUids,
						ActionPanelButtons:        tc.actionPanelButtons,
						HasV1PersonsGetPermission: tc.hasV1PersonsGetPermission,
						SelectedCaseIds:           tc.selectedCaseIds,
					},
				).
				Return(nil)

			server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

			r, _ := http.NewRequest(http.MethodGet, tc.path, nil)
			_, err := server.serve(r)

			assert.Nil(t, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetDocumentListHasV1PersonsCasesGetPermission(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", UID: "7000-1234-0000"}}

	testCases := []struct {
		name                                   string
		permissions                            sirius.Permissions
		expectedHasV1PersonsCasesGetPermission bool
	}{
		{
			name:                                   "with v1-persons-cases GET permission",
			permissions:                            sirius.Permissions{"v1-persons-cases": sirius.PermissionType{Permissions: []string{"GET"}}},
			expectedHasV1PersonsCasesGetPermission: true,
		},
		{
			name:                                   "without v1-persons-cases GET permission",
			permissions:                            sirius.Permissions{},
			expectedHasV1PersonsCasesGetPermission: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockDocumentListClient{}
			client.
				On("CasesByDonor", mock.Anything, 82).
				Return(cases, nil)
			client.
				On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
				Return(singleDocumentList, nil)
			client.
				On("Person", mock.Anything, 82).
				Return(sirius.Person{}, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(tc.permissions, nil)
			client.
				On("GetDraftCount", mock.Anything, "lpa", 1).
				Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
			client.
				On("PersonReferences", mock.Anything, 82).
				Return([]sirius.PersonReference{{ID: 987}}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, mock.MatchedBy(func(data documentPageData) bool {
					return data.HasV1PersonsCasesGetPermission == tc.expectedHasV1PersonsCasesGetPermission
				})).
				Return(nil)

			server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

			req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
			_, err := server.serve(req)

			assert.Nil(t, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestDocumentListDownloadMultipleSuccess(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("PersonReferences", mock.Anything, 82).
		Return([]sirius.PersonReference{{ID: 987}}, nil)

	downloadResp := &http.Response{
		StatusCode: http.StatusCreated,
		Header: http.Header{
			"Content-Type": []string{"content/octet-stream"},
		},
		Body: io.NopCloser(strings.NewReader("document-download-bytes")),
	}

	client.
		On("DownloadMultiple", mock.Anything, []string{"doc-uuid"}).
		Return(downloadResp, nil)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	form := url.Values{}
	form.Add("document", "doc-uuid")
	form.Add("actionDownload", "true")
	req, _ := http.NewRequest(http.MethodPost, "/donor/82/documents", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formUrlEncoded)

	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, "content/octet-stream", resp.Header().Get("Content-Type"))
	assert.Equal(t, "document-download-bytes", resp.Body.String())

	mock.AssertExpectationsForObjects(t, client)
	client.AssertNotCalled(t, "GetPersonDocuments")
}

func TestDocumentListDownloadMultipleError(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("DownloadMultiple", mock.Anything, []string{"doc-uuid"}).
		Return((*http.Response)(nil), errExample)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("PersonReferences", mock.Anything, 82).
		Return([]sirius.PersonReference{{ID: 987}}, nil)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	form := url.Values{}
	form.Add("document", "doc-uuid")
	form.Add("actionDownload", "true")
	req, _ := http.NewRequest(http.MethodPost, "/donor/82/documents", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formUrlEncoded)

	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
	client.AssertNotCalled(t, "GetPersonDocuments")
}

func TestDocumentListShowsValidationErrorWhenNoDocumentsSelected(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-1234-0000"},
		{ID: 2, UID: "7000-9876-0000"},
	}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
		Return(allDocumentList, nil)
	client.
		On("Person", mock.Anything, 82).
		Return(expectedDonor, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("PersonReferences", mock.Anything, 82).
		Return([]sirius.PersonReference{{ID: 987}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			documentPageData{
				SelectedCases:         cases,
				SelectedCaseIds:       "1+2",
				Person:                expectedDonor,
				DocumentList:          allDocumentList,
				MultipleCasesSelected: true,
				DonorID:               82,
				Error: sirius.ValidationError{
					Detail: "Select one or more documents and try again.",
				},
				HasV1PersonsGetPermission: false,
				ActionPanelButtons: []ActionPanelButton{
					{
						Label:    "Create warning",
						URL:      "/create-warning?id=82&entity=person",
						IconName: "aw-create-warning",
						Disabled: false,
					},
					{
						Label:    "Create event",
						URL:      "/create-event?id=82&entity=person",
						IconName: "aw-new-event",
						Disabled: false,
					},
					{
						Label:    "Add complaint",
						URL:      "",
						IconName: "aw-log-complaint",
						Disabled: true,
					},
					{
						Label:    "Create document",
						URL:      "",
						IconName: "aw-new-template",
						Disabled: true,
					},
					{
						Label:    "Retrieve draft",
						URL:      "",
						IconName: "aw-new-template",
						Disabled: true,
					},
					{
						Label:    "Change status",
						URL:      "",
						IconName: "aw-change-status",
						Disabled: true,
					},
					{
						Label:    "Fees",
						URL:      "",
						IconName: "aw-fees",
						Disabled: true,
					},
					{
						Label:    "New task",
						URL:      "",
						IconName: "aw-new-task",
						Disabled: true,
					},
					{
						Label:    "Create donor",
						URL:      "/create-donor?id=82&entity=person",
						IconName: "aw-create-person",
						Disabled: false,
					},
					{
						Label:    "Edit donor",
						URL:      "/edit-donor?id=82&entity=person",
						IconName: "aw-edit-person",
						Disabled: false,
					},
					{
						Label:    "Edit dates",
						URL:      "",
						IconName: "calendar-open",
						Disabled: true,
					},
					{
						Label:    "MI reporting",
						URL:      "/mi-reporting?donorId=82",
						IconName: "aw-mi",
						Disabled: false,
					},
					{
						Label:    "Allocate Case",
						URL:      "/allocate-cases?id=1&id=2&entity=",
						IconName: "aw-allocate-case",
						Disabled: false,
					},
					{
						Label:    "Link record",
						URL:      "/link-person?id=82",
						IconName: "aw-link",
						Disabled: false,
					},
					{
						Label:    "Delete relationship",
						URL:      "/delete-relationship?id=82",
						IconName: "icon-minus",
						Disabled: false,
					},
				},
			},
		).
		Return(nil)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

	form := url.Values{}
	form.Add("actionDownload", "true")
	req, _ := http.NewRequest(http.MethodPost, "/donor/82/documents", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formUrlEncoded)

	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	mock.AssertExpectationsForObjects(t, client, template)
	client.AssertNotCalled(t, "DownloadMultiple")
}

func TestDocumentListReturnsNoContentWhenComparingAndNoDocumentsSelected(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-1234-0000"},
		{ID: 2, UID: "7000-9876-0000"},
	}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
		Return(allDocumentList, nil)
	client.
		On("Person", mock.Anything, 82).
		Return(expectedDonor, nil)
	client.
		On("PersonReferences", mock.Anything, 82).
		Return([]sirius.PersonReference{{ID: 987}}, nil)

	template := &mockTemplate{}

	server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

	form := url.Values{}
	form.Add("actionDownload", "true")
	form.Add("comparing", "true")
	req, _ := http.NewRequest(http.MethodPost, "/donor/82/documents", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formUrlEncoded)

	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	client.AssertNotCalled(t, "DownloadMultiple")
	mock.AssertExpectationsForObjects(t, client)
}

func TestDocumentListDismissValidation(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-1234-0000"},
		{ID: 2, UID: "7000-9876-0000"},
		{ID: 3, UID: "7000-5678-0000"},
	}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string{"1", "2"}).
		Return(twoCasesDocumentList, nil)
	client.
		On("Person", mock.Anything, 82).
		Return(expectedDonor, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("PersonReferences", mock.Anything, 82).
		Return([]sirius.PersonReference{{ID: 987}}, nil)

	expectedCases := []sirius.Case{cases[0], cases[1]}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			documentPageData{
				SelectedCases:             expectedCases,
				SelectedCaseIds:           "1+2",
				Person:                    expectedDonor,
				DocumentList:              twoCasesDocumentList,
				MultipleCasesSelected:     true,
				DonorID:                   82,
				CaseUids:                  "&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
				HasV1PersonsGetPermission: false,
				ActionPanelButtons: []ActionPanelButton{
					{
						Label:    "Create warning",
						URL:      "/create-warning?id=82&entity=person&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
						IconName: "aw-create-warning",
						Disabled: false,
					},
					{
						Label:    "Create event",
						URL:      "/create-event?id=82&entity=person&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
						IconName: "aw-new-event",
						Disabled: false,
					},
					{
						Label:    "Add complaint",
						URL:      "",
						IconName: "aw-log-complaint",
						Disabled: true,
					},
					{
						Label:    "Create document",
						URL:      "",
						IconName: "aw-new-template",
						Disabled: true,
					},
					{
						Label:    "Retrieve draft",
						URL:      "",
						IconName: "aw-new-template",
						Disabled: true,
					},
					{
						Label:    "Change status",
						URL:      "",
						IconName: "aw-change-status",
						Disabled: true,
					},
					{
						Label:    "Fees",
						URL:      "",
						IconName: "aw-fees",
						Disabled: true,
					},
					{
						Label:    "New task",
						URL:      "",
						IconName: "aw-new-task",
						Disabled: true,
					},
					{
						Label:    "Create donor",
						URL:      "/create-donor?id=82&entity=person&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
						IconName: "aw-create-person",
						Disabled: false,
					},
					{
						Label:    "Edit donor",
						URL:      "/edit-donor?id=82&entity=person&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
						IconName: "aw-edit-person",
						Disabled: false,
					},
					{
						Label:    "Edit dates",
						URL:      "",
						IconName: "calendar-open",
						Disabled: true,
					},
					{
						Label:    "MI reporting",
						URL:      "/mi-reporting?donorId=82&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
						IconName: "aw-mi",
						Disabled: false,
					},
					{
						Label:    "Allocate Case",
						URL:      "/allocate-cases?id=1&id=2&entity=&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
						IconName: "aw-allocate-case",
						Disabled: false,
					},
					{
						Label:    "Link record",
						URL:      "/link-person?id=82&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
						IconName: "aw-link",
						Disabled: false,
					},
					{
						Label:    "Delete relationship",
						URL:      "/delete-relationship?id=82&uid[]=7000-1234-0000&uid[]=7000-9876-0000",
						IconName: "icon-minus",
						Disabled: false,
					},
				},
			},
		).
		Return(nil)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

	form := url.Values{}
	form.Add("uid[]", "7000-1234-0000")
	form.Add("uid[]", "7000-9876-0000")
	form.Add("dismissValidation", "true")
	req, _ := http.NewRequest(http.MethodPost, "/donor/82/documents", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formUrlEncoded)

	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	mock.AssertExpectationsForObjects(t, client, template)
	client.AssertNotCalled(t, "DownloadMultiple")
}

func TestDocumentListInvalidDonorID(t *testing.T) {
	client := &mockDocumentListClient{}
	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/donor/abc/documents", nil)
	_, err := server.serve(req)

	assert.Error(t, err)
}

func TestGetDocumentListWhenCasesByDonorErrors(t *testing.T) {
	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return([]sirius.Case{}, errExample)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDocumentListWhenGetPersonDocumentsErrors(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", SubType: "PFA", UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
		Return(sirius.DocumentList{}, errExample)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("PersonReferences", mock.Anything, 82).
		Return([]sirius.PersonReference{{ID: 987}}, nil)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDocumentListWhenPersonErrors(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", SubType: "PFA", UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
		Return(singleDocumentList, nil)
	client.
		On("Person", mock.Anything, 82).
		Return(sirius.Person{}, errExample)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("PersonReferences", mock.Anything, 82).
		Return([]sirius.PersonReference{{ID: 987}}, nil)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDocumentListWhenPermissionsErrors(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", SubType: "PFA", UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
		Return(singleDocumentList, nil)
	client.
		On("Person", mock.Anything, 82).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, errExample)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("PersonReferences", mock.Anything, 82).
		Return([]sirius.PersonReference{{ID: 987}}, nil)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDocumentListWhenGetDraftCountErrors(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", SubType: "PFA", UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{}, errExample)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDocumentListWhenGetPersonReferencesErrors(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", SubType: "PFA", UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("PersonReferences", mock.Anything, 82).
		Return([]sirius.PersonReference{}, errExample)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDocumentListWhenTemplateErrors(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", SubType: "PFA", UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
		Return(singleDocumentList, nil)
	client.
		On("Person", mock.Anything, 82).
		Return(expectedDonor, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("PersonReferences", mock.Anything, 82).
		Return([]sirius.PersonReference{{ID: 987}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			documentPageData{
				SelectedCases:             cases,
				SelectedCaseIds:           "1",
				Person:                    expectedDonor,
				DocumentList:              singleDocumentList,
				MultipleCasesSelected:     false,
				DonorID:                   82,
				HasV1PersonsGetPermission: false,
				ActionPanelButtons: []ActionPanelButton{
					{
						Label:    "Create warning",
						URL:      "/create-warning?id=82&entity=lpa",
						IconName: "aw-create-warning",
						Disabled: false,
					},
					{
						Label:    "Create event",
						URL:      "/create-event?id=82&entity=person",
						IconName: "aw-new-event",
						Disabled: false,
					},
					{
						Label:    "Add complaint",
						URL:      "/add-complaint?id=1&case=lpa",
						IconName: "aw-log-complaint",
						Disabled: false,
					},
					{
						Label:    "Create document",
						URL:      "/create-document?id=1&case=lpa",
						IconName: "aw-new-template",
						Disabled: false,
					},
					{
						Label:    "Retrieve draft",
						URL:      "/edit-document?id=1&case=lpa",
						IconName: "aw-new-template",
						Disabled: false,
					},
					{
						Label:    "Change status",
						URL:      "/change-status?id=1&case=lpa&donorId=82",
						IconName: "aw-change-status",
						Disabled: false,
					},
					{
						Label:    "Fees",
						URL:      "/payments/1",
						IconName: "aw-fees",
						Disabled: false,
					},
					{
						Label:    "New task",
						URL:      "/create-task?id=1&entity=lpa",
						IconName: "aw-new-task",
						Disabled: false,
					},
					{
						Label:    "Create donor",
						URL:      "/create-donor?id=82&entity=person",
						IconName: "aw-create-person",
						Disabled: false,
					},
					{
						Label:    "Edit donor",
						URL:      "/edit-donor?id=82&entity=person",
						IconName: "aw-edit-person",
						Disabled: false,
					},
					{
						Label:    "Edit dates",
						URL:      "/edit-dates?id=1&case=lpa",
						IconName: "calendar-open",
						Disabled: false,
					},
					{
						Label:    "MI reporting",
						URL:      "/mi-reporting?donorId=82",
						IconName: "aw-mi",
						Disabled: false,
					},
					{
						Label:    "Allocate Case",
						URL:      "/allocate-cases?id=1&entity=lpa",
						IconName: "aw-allocate-case",
						Disabled: false,
					},
					{
						Label:    "Link record",
						URL:      "/link-person?id=82",
						IconName: "aw-link",
						Disabled: false,
					},
					{
						Label:    "Delete relationship",
						URL:      "/delete-relationship?id=82",
						IconName: "icon-minus",
						Disabled: false,
					},
				},
			},
		).
		Return(errExample)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestSuccessMessageFormatter(t *testing.T) {
	tests := []struct {
		name            string
		docFriendlyName string
		docCreatedTime  string
		layout          string
		format          string
		expectedResult  string
	}{
		{
			name:            "valid date format",
			docFriendlyName: "Letter",
			docCreatedTime:  "02/07/2025 14:17:13",
			layout:          "02/01/2006 15:04:05",
			format:          "02/01/2006",
			expectedResult:  "02/07/2025 Letter",
		},
		{
			name:            "different valid date",
			docFriendlyName: "LP1H - Health Instrument",
			docCreatedTime:  "01/06/2022 15:39:01",
			layout:          "02/01/2006 15:04:05",
			format:          "02/01/2006",
			expectedResult:  "01/06/2022 LP1H - Health Instrument",
		},
		{
			name:            "invalid date format",
			docFriendlyName: "Document",
			docCreatedTime:  "invalid-date",
			layout:          "02/01/2006 15:04:05",
			format:          "02/01/2006",
			expectedResult:  "invalid date",
		},
		{
			name:            "empty friendly name",
			docFriendlyName: "",
			docCreatedTime:  "08/01/2025 10:36:41",
			layout:          "02/01/2006 15:04:05",
			format:          "02/01/2006",
			expectedResult:  "08/01/2025 ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := successMessageFormatter(tc.docFriendlyName, tc.docCreatedTime, tc.layout, tc.format)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestDocumentListSuccessMessage(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", UID: "7000-1234-0000"}}

	tests := []struct {
		name            string
		queryParams     string
		formValue       string
		expectedSuccess bool
		expectedMessage string
	}{
		{
			name:            "success query with no dismiss notification",
			queryParams:     "?success=true&documentFriendlyName=Letter&documentCreatedTime=02/07/2025%2014:17:13",
			formValue:       "",
			expectedSuccess: true,
			expectedMessage: "02/07/2025 Letter",
		},
		{
			name:            "success query with dismiss notification",
			queryParams:     "?success=true&documentFriendlyName=Letter&documentCreatedTime=02/07/2025%2014:17:13",
			formValue:       "true",
			expectedSuccess: false,
			expectedMessage: "",
		},
		{
			name:            "no success query",
			queryParams:     "?documentFriendlyName=Letter",
			formValue:       "",
			expectedSuccess: false,
			expectedMessage: "",
		},
		{
			name:            "success with invalid date",
			queryParams:     "?success=true&documentFriendlyName=Letter&documentCreatedTime=invalid",
			formValue:       "",
			expectedSuccess: true,
			expectedMessage: "invalid date",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockDocumentListClient{}
			client.
				On("CasesByDonor", mock.Anything, 82).
				Return(cases, nil)
			client.
				On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
				Return(singleDocumentList, nil)
			client.
				On("Person", mock.Anything, 82).
				Return(expectedDonor, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{}, nil)
			client.
				On("GetDraftCount", mock.Anything, "lpa", 1).
				Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
			client.
				On("PersonReferences", mock.Anything, 82).
				Return([]sirius.PersonReference{{ID: 987}}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, mock.MatchedBy(func(data documentPageData) bool {
					return data.Success == tc.expectedSuccess && data.SuccessMessage == tc.expectedMessage
				})).
				Return(nil)

			server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

			form := url.Values{}
			if tc.formValue != "" {
				form.Add("dismissNotification", tc.formValue)
			}
			req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents"+tc.queryParams, nil)
			if len(form) > 0 {
				req.Form = form
			}

			resp, err := server.serve(req)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
