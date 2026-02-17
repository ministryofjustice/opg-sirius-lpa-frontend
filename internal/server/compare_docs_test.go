package server

import (
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

func TestGetCompareDocsPane1(t *testing.T) {
	document := sirius.Document{
		ID:        1,
		UUID:      "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		CaseItems: []sirius.Case{{ID: 456}},
	}
	documentList := sirius.DocumentList{
		Documents: []sirius.Document{document},
	}

	client := &mockCompareDocsClient{}
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 77, []string{"456"}).
		Return(documentList, nil)

	template := &mockTemplate{}
	templateData := compareDocsData{
		DocListPane1Data: documentPageData{
			DocumentList:  documentList,
			SelectedCases: document.CaseItems,
			Comparing:     true,
			CompareURLs: map[string]string{
				"dfef6714-b4fe-44c2-b26e-90dfe3663e95": "/compare/77/456?pane1=dfef6714-b4fe-44c2-b26e-90dfe3663e95",
			},
		},
		DocListPane2Data: documentPageData{
			DocumentList:  documentList,
			SelectedCases: document.CaseItems,
			Comparing:     true,
			CompareURLs: map[string]string{
				"dfef6714-b4fe-44c2-b26e-90dfe3663e95": "/compare/77/456?pane2=dfef6714-b4fe-44c2-b26e-90dfe3663e95&pane1=dfef6714-b4fe-44c2-b26e-90dfe3663e95",
			},
		},
		Pane1: "doc",
		Pane2: "list",
		View1: &viewingDocumentData{
			Document: document,
			Pane:     1,
			BackURL:  "/compare/77/456",
		},
		View2: nil,
	}

	template.
		On("Func", mock.Anything, templateData).
		Return(nil)

	server := newMockServer("/compare/{id}/{caseId}", CompareDocs(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/compare/77/456?pane1=dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCompareDocsPane2(t *testing.T) {
	document := sirius.Document{
		ID:        1,
		UUID:      "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		CaseItems: []sirius.Case{{ID: 456}},
	}
	documentList := sirius.DocumentList{
		Documents: []sirius.Document{document},
	}

	client := &mockCompareDocsClient{}
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 77, []string{"456"}).
		Return(documentList, nil)

	template := &mockTemplate{}
	templateData := compareDocsData{
		DocListPane1Data: documentPageData{
			DocumentList:  documentList,
			SelectedCases: document.CaseItems,
			Comparing:     true,
			CompareURLs: map[string]string{
				"dfef6714-b4fe-44c2-b26e-90dfe3663e95": "/compare/77/456?pane1=dfef6714-b4fe-44c2-b26e-90dfe3663e95&pane2=dfef6714-b4fe-44c2-b26e-90dfe3663e95",
			},
		},
		DocListPane2Data: documentPageData{
			DocumentList:  documentList,
			SelectedCases: document.CaseItems,
			Comparing:     true,
			CompareURLs: map[string]string{
				"dfef6714-b4fe-44c2-b26e-90dfe3663e95": "/compare/77/456?pane2=dfef6714-b4fe-44c2-b26e-90dfe3663e95",
			},
		},
		Pane1: "list",
		Pane2: "doc",
		View1: nil,
		View2: &viewingDocumentData{
			Document: document,
			Pane:     2,
			BackURL:  "/compare/77/456",
		},
	}

	template.
		On("Func", mock.Anything, templateData).
		Return(nil)

	server := newMockServer("/compare/{id}/{caseId}", CompareDocs(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/compare/77/456?pane2=dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCompareDocs(t *testing.T) {
	document1 := sirius.Document{
		ID:        1,
		UUID:      "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		CaseItems: []sirius.Case{{ID: 456}},
	}
	document2 := sirius.Document{
		ID:        2,
		UUID:      "agyq4619-j3yq-88b1-p95d-03hes5772k46",
		CaseItems: []sirius.Case{{ID: 456}},
	}
	documentList := sirius.DocumentList{
		Documents: []sirius.Document{document1, document2},
	}

	client := &mockCompareDocsClient{}
	client.
		On("GetPersonDocuments", mock.Anything, 77, []string{"456"}).
		Return(documentList, nil)

	template := &mockTemplate{}
	templateData := compareDocsData{
		DocListPane1Data: documentPageData{
			DocumentList:  documentList,
			SelectedCases: document1.CaseItems,
			Comparing:     true,
			CompareURLs: map[string]string{
				"agyq4619-j3yq-88b1-p95d-03hes5772k46": "/compare/77/456?pane1=agyq4619-j3yq-88b1-p95d-03hes5772k46",
				"dfef6714-b4fe-44c2-b26e-90dfe3663e95": "/compare/77/456?pane1=dfef6714-b4fe-44c2-b26e-90dfe3663e95",
			},
		},
		DocListPane2Data: documentPageData{
			DocumentList:  documentList,
			SelectedCases: document1.CaseItems,
			Comparing:     true,
			CompareURLs: map[string]string{
				"agyq4619-j3yq-88b1-p95d-03hes5772k46": "/compare/77/456?pane2=agyq4619-j3yq-88b1-p95d-03hes5772k46",
				"dfef6714-b4fe-44c2-b26e-90dfe3663e95": "/compare/77/456?pane2=dfef6714-b4fe-44c2-b26e-90dfe3663e95",
			},
		},
		Pane1: "list",
		Pane2: "list",
		View1: nil,
		View2: nil,
	}

	template.
		On("Func", mock.Anything, templateData).
		Return(nil)

	server := newMockServer("/compare/{id}/{caseId}", CompareDocs(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/compare/77/456", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
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

func TestGetCompareDocsWhenCase1Errors(t *testing.T) {
	document := sirius.Document{
		ID:        1,
		UUID:      "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		CaseItems: []sirius.Case{{ID: 456}},
	}
	documentList := sirius.DocumentList{
		Documents: []sirius.Document{document},
	}

	client := &mockCompareDocsClient{}
	client.
		On("GetPersonDocuments", mock.Anything, 77, []string{"456"}).
		Return(documentList, nil)
	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(sirius.Document{}, errExample)

	server := newMockServer("/compare/{id}/{caseId}", CompareDocs(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/compare/77/456?pane1=dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetCompareDocsWhenCase2Errors(t *testing.T) {
	document := sirius.Document{
		ID:        1,
		UUID:      "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		CaseItems: []sirius.Case{{ID: 456}},
	}
	documentList := sirius.DocumentList{
		Documents: []sirius.Document{document},
	}

	client := &mockCompareDocsClient{}
	client.
		On("GetPersonDocuments", mock.Anything, 77, []string{"456"}).
		Return(documentList, nil)
	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(sirius.Document{}, errExample)

	server := newMockServer("/compare/{id}/{caseId}", CompareDocs(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/compare/77/456?pane2=dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetCompareDocsBadID(t *testing.T) {
	client := &mockCompareDocsClient{}
	template := &mockTemplate{}

	server := newMockServer("/compare/{id}/{caseId}", CompareDocs(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/compare/bad-id/456", nil)
	_, err := server.serve(req)

	assert.NotNil(t, err)
}
