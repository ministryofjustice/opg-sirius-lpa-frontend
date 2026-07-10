package server

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCompareDocsClient struct {
	mock.Mock
}

func (m *mockCompareDocsClient) DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func (m *mockCompareDocsClient) GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error) {
	args := m.Called(ctx, personID, caseIDs)
	return args.Get(0).(sirius.DocumentList), args.Error(1)
}

func (m *mockCompareDocsClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockCompareDocsClient) GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.Permissions), args.Error(1)
}

func (m *mockCompareDocsClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockCompareDocsClient) GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error) {
	args := m.Called(ctx, caseType, caseId)
	return args.Get(0).(sirius.DocumentDraftCount), args.Error(1)
}

func (m *mockCompareDocsClient) PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.PersonReference), args.Error(1)
}

func TestGetCompareDocsPanes(t *testing.T) {
	document1 := sirius.Document{
		ID:        1,
		UUID:      "doc1-uuid",
		CaseItems: []sirius.Case{{ID: 456, UID: "case-uid"}},
	}
	document2 := sirius.Document{
		ID:        2,
		UUID:      "doc2-uuid",
		CaseItems: []sirius.Case{{ID: 456, UID: "case-uid"}},
	}
	documentList := sirius.DocumentList{
		Documents: []sirius.Document{document1, document2},
	}
	person := sirius.Person{ID: 77}
	selectedCase := document1.CaseItems[0]
	permissions := sirius.Permissions{
		"v1-persons":       {Permissions: []string{"GET"}},
		"v1-persons-cases": {Permissions: []string{"GET"}},
	}

	tests := []struct {
		name                     string
		query                    string
		pane1                    string
		pane2                    string
		view1                    *viewingDocumentData
		view2                    *viewingDocumentData
		compareURLs1             map[string]string
		compareURLs2             map[string]string
		getDocuments             []sirius.Document
		closeURLToDocumentPanel1 string
		closeURLToDocumentPanel2 string
	}{
		{
			name:  "Pane 1 and Pane 2 shows document lists",
			query: "",
			pane1: "list",
			pane2: "list",
			view1: nil,
			view2: nil,
			compareURLs1: map[string]string{
				"doc1-uuid": "/compare/77/456?pane1=doc1-uuid",
				"doc2-uuid": "/compare/77/456?pane1=doc2-uuid",
			},
			compareURLs2: map[string]string{
				"doc1-uuid": "/compare/77/456?pane2=doc1-uuid",
				"doc2-uuid": "/compare/77/456?pane2=doc2-uuid",
			},
			getDocuments:             nil,
			closeURLToDocumentPanel1: "/donor/77/documents?uid[]=case-uid",
			closeURLToDocumentPanel2: "/donor/77/documents?uid[]=case-uid",
		},
		{
			name:  "Pane 1 shows a document, Pane 2 shows a doc list",
			query: "?pane1=doc1-uuid",
			pane1: "doc",
			pane2: "list",
			view1: &viewingDocumentData{
				Document: document1,
				Pane:     1,
				BackURL:  "/compare/77/456",
				CloseURL: "/donor/77/documents?uid[]=case-uid",
			},
			view2: nil,
			compareURLs1: map[string]string{
				"doc1-uuid": "/compare/77/456?pane1=doc1-uuid",
				"doc2-uuid": "/compare/77/456?pane1=doc2-uuid",
			},
			compareURLs2: map[string]string{
				"doc1-uuid": "/compare/77/456?pane2=doc1-uuid&pane1=doc1-uuid",
				"doc2-uuid": "/compare/77/456?pane2=doc2-uuid&pane1=doc1-uuid",
			},
			getDocuments:             []sirius.Document{document1},
			closeURLToDocumentPanel1: "",
			closeURLToDocumentPanel2: "/view-document/doc1-uuid/77?case=456&pane=1",
		},
		{
			name:  "Pane 1 shows a doc list, Pane 2 shows a document",
			query: "?pane2=doc1-uuid",
			pane1: "list",
			pane2: "doc",
			view1: nil,
			view2: &viewingDocumentData{
				Document: document1,
				Pane:     2,
				BackURL:  "/compare/77/456",
				CloseURL: "/donor/77/documents?uid[]=case-uid",
			},
			compareURLs1: map[string]string{
				"doc1-uuid": "/compare/77/456?pane1=doc1-uuid&pane2=doc1-uuid",
				"doc2-uuid": "/compare/77/456?pane1=doc2-uuid&pane2=doc1-uuid",
			},
			compareURLs2: map[string]string{
				"doc1-uuid": "/compare/77/456?pane2=doc1-uuid",
				"doc2-uuid": "/compare/77/456?pane2=doc2-uuid",
			},
			getDocuments:             []sirius.Document{document1},
			closeURLToDocumentPanel1: "/view-document/doc1-uuid/77?case=456&pane=2",
			closeURLToDocumentPanel2: "",
		},
		{
			name:  "Pane 1 and Pane 2 shows a document each",
			query: "?pane1=doc1-uuid&pane2=doc2-uuid",
			pane1: "doc",
			pane2: "doc",
			view1: &viewingDocumentData{
				Document: document1,
				Pane:     1,
				BackURL:  "/compare/77/456?pane2=doc2-uuid",
				CloseURL: "/view-document/doc2-uuid/77?case=456&pane=2",
			},
			view2: &viewingDocumentData{
				Document: document2,
				Pane:     2,
				BackURL:  "/compare/77/456?pane1=doc1-uuid",
				CloseURL: "/view-document/doc1-uuid/77?case=456&pane=1",
			},
			compareURLs1: map[string]string{
				"doc1-uuid": "/compare/77/456?pane1=doc1-uuid&pane2=doc2-uuid",
				"doc2-uuid": "/compare/77/456?pane1=doc2-uuid&pane2=doc2-uuid",
			},
			compareURLs2: map[string]string{
				"doc1-uuid": "/compare/77/456?pane2=doc1-uuid&pane1=doc1-uuid",
				"doc2-uuid": "/compare/77/456?pane2=doc2-uuid&pane1=doc1-uuid",
			},
			getDocuments:             []sirius.Document{document1, document2},
			closeURLToDocumentPanel1: "",
			closeURLToDocumentPanel2: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockCompareDocsClient{}
			client.
				On("GetPersonDocuments", mock.Anything, 77, []string{"456"}).
				Return(documentList, nil)
			client.
				On("Person", mock.Anything, 77).
				Return(person, nil)
			client.
				On("PersonReferences", mock.Anything, 77).
				Return([]sirius.PersonReference{}, nil)
			client.
				On("Case", mock.Anything, 456).
				Return(selectedCase, nil)
			client.
				On("GetDraftCount", mock.Anything, "", 456).
				Return(sirius.DocumentDraftCount{DraftCount: 0}, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(permissions, nil)
			for _, doc := range tc.getDocuments {
				client.
					On("DocumentByUUID", mock.Anything, doc.UUID).
					Return(doc, nil)
			}

			template := &mockTemplate{}
			templateData := compareDocsData{
				DocListPane1Data: documentPageData{
					DocumentList:  documentList,
					SelectedCases: []sirius.Case{selectedCase},
					Comparing:     true,
					CompareURLs:   tc.compareURLs1,
					CloseURL:      tc.closeURLToDocumentPanel1,
					DonorID:       77,
				},
				DocListPane2Data: documentPageData{
					DocumentList:  documentList,
					SelectedCases: []sirius.Case{selectedCase},
					Comparing:     true,
					CompareURLs:   tc.compareURLs2,
					CloseURL:      tc.closeURLToDocumentPanel2,
					DonorID:       77,
				},
				Pane1:                          tc.pane1,
				Pane2:                          tc.pane2,
				View1:                          tc.view1,
				View2:                          tc.view2,
				DonorID:                        77,
				SelectedCaseIds:                "456",
				Person:                         person,
				CaseUids:                       "&uid[]=case-uid",
				HasV1PersonsGetPermission:      true,
				HasV1PersonsCasesGetPermission: true,
				SelectedCases:                  []sirius.Case{selectedCase},
				ActionPanelButtons:             GetActionPanelButtons([]sirius.Case{selectedCase}, 77, "&uid[]=case-uid", false, false, false, permissions),
				HeaderButtons: SiriusHeaderButtons{
					BackToTimeline: true,
					CaseInfo:       true,
					PersonInfo:     true,
					Calendar:       true,
				},
			}

			template.
				On("Func", mock.Anything, templateData).
				Return(nil)

			server := newMockServer("/compare/{id}/{caseId}", CompareDocs(client, template.Func))

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/compare/77/456%s", tc.query), nil)
			_, err := server.serve(req)

			assert.Nil(t, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCompareDocsWhenGetUserDetailsErrors(t *testing.T) {
	client := &mockCompareDocsClient{}
	client.
		On("GetPersonDocuments", mock.Anything, 77, []string{"456"}).
		Return(sirius.DocumentList{}, errExample)

	server := newMockServer("/compare/{id}/{caseId}", CompareDocs(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/compare/77/456", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetCompareDocsWhenCaseErrors(t *testing.T) {
	document := sirius.Document{
		ID:        1,
		UUID:      "doc-uuid",
		CaseItems: []sirius.Case{{ID: 456}},
	}
	documentList := sirius.DocumentList{
		Documents: []sirius.Document{document},
	}
	person := sirius.Person{ID: 77}

	tests := []struct {
		name string
		pane string
	}{
		{
			name: "pane 1 errors",
			pane: "1",
		},
		{
			name: "pane 2 errors",
			pane: "2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockCompareDocsClient{}
			client.
				On("GetPersonDocuments", mock.Anything, 77, []string{"456"}).
				Return(documentList, nil)
			client.
				On("Person", mock.Anything, 77).
				Return(person, nil)
			client.
				On("PersonReferences", mock.Anything, 77).
				Return([]sirius.PersonReference{}, nil)
			client.
				On("Case", mock.Anything, 456).
				Return(document.CaseItems[0], nil)
			client.
				On("GetDraftCount", mock.Anything, "", 456).
				Return(sirius.DocumentDraftCount{DraftCount: 0}, nil)
			client.
				On("DocumentByUUID", mock.Anything, "abcd").
				Return(sirius.Document{}, errExample)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{}, nil)

			server := newMockServer("/compare/{id}/{caseId}", CompareDocs(client, nil))

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/compare/77/456?pane%s=abcd", tc.pane), nil)
			_, err := server.serve(req)

			assert.Equal(t, errExample, err)
		})
	}
}

func TestGetCompareDocsBadID(t *testing.T) {
	client := &mockCompareDocsClient{}
	template := &mockTemplate{}

	server := newMockServer("/compare/{id}/{caseId}", CompareDocs(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/compare/bad-id/456", nil)
	_, err := server.serve(req)

	assert.NotNil(t, err)
}
