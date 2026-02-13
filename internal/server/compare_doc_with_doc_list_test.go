package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCompareDocumentClient struct {
	mock.Mock
}

func (m *mockCompareDocumentClient) DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func (m *mockCompareDocumentClient) GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error) {
	args := m.Called(ctx, personID, caseIDs)
	return args.Get(0).(sirius.DocumentList), args.Error(1)
}

func TestGetCompareDocument(t *testing.T) {
	document := sirius.Document{
		ID:        1,
		UUID:      "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		CaseItems: []sirius.Case{{ID: 456}},
	}
	documentList := sirius.DocumentList{
		Documents: []sirius.Document{document},
	}

	client := &mockCompareDocumentClient{}
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 77, []string{"456"}).
		Return(documentList, nil)

	template := &mockTemplate{}
	templateData := documentPageData{
		DocumentList:  documentList,
		Document:      document,
		SelectedCases: document.CaseItems,
		Comparing:     true,
	}

	template.
		On("Func", mock.Anything, templateData).
		Return(nil)

	server := newMockServer("/comparing-document/{id}", CompareDocWithDocList(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/comparing-document/77?uid[]=dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCompareDocumentsBadID(t *testing.T) {
	client := &mockCompareDocumentClient{}
	template := &mockTemplate{}

	server := newMockServer("/compare-document/{id}", CompareDocWithDocList(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/compare-document/bad-id", nil)
	_, err := server.serve(req)

	assert.NotNil(t, err)
}

func TestGetCompareDocumentWhenCaseErrors(t *testing.T) {
	client := &mockCompareDocumentClient{}

	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(sirius.Document{}, errExample)

	server := newMockServer("/compare-document/{id}", CompareDocWithDocList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/compare-document/81?uid[]=dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetCompareDocumentWhenGetUserDetailsErrors(t *testing.T) {
	client := &mockCompareDocumentClient{}
	document := sirius.Document{
		ID:        1,
		UUID:      "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		CaseItems: []sirius.Case{{ID: 456}},
	}

	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 77, []string{"456"}).
		Return(sirius.DocumentList{}, errExample)

	server := newMockServer("/comparing-document/{id}", CompareDocWithDocList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/comparing-document/77?uid[]=dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetCompareDocumentWhenNoCasesOnDoc(t *testing.T) {
	document := sirius.Document{
		ID:   1,
		UUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
	}
	documentList := sirius.DocumentList{
		Documents: []sirius.Document{document},
	}

	client := &mockCompareDocumentClient{}
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 77, []string{}).
		Return(documentList, nil)

	template := &mockTemplate{}
	templateData := documentPageData{
		DocumentList:  documentList,
		Document:      document,
		SelectedCases: []sirius.Case{},
		Comparing:     true,
	}

	template.
		On("Func", mock.Anything, templateData).
		Return(nil)

	server := newMockServer("/comparing-document/{id}", CompareDocWithDocList(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/comparing-document/77?uid[]=dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
