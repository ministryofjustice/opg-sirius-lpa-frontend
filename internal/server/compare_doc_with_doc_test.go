package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockComparingDocumentClient struct {
	mock.Mock
}

func (m *mockComparingDocumentClient) DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func TestGetComparingDocument(t *testing.T) {
	document := sirius.Document{
		ID:   1,
		UUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
	}
	documentComparing := sirius.Document{
		ID:   564,
		UUID: "fweq6452-q4az-44c2-b26e-90dfe3679u22",
	}

	client := &mockComparingDocumentClient{}
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("DocumentByUUID", mock.Anything, documentComparing.UUID).
		Return(documentComparing, nil)

	template := &mockTemplate{}
	templateData := comparingDocumentsData{
		Document:          document,
		DocumentComparing: documentComparing,
	}

	template.
		On("Func", mock.Anything, templateData).
		Return(nil)

	server := newMockServer("/comparing-document", CompareDocWithDoc(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/comparing-document?docUid[]=dfef6714-b4fe-44c2-b26e-90dfe3663e95&docUid[]=fweq6452-q4az-44c2-b26e-90dfe3679u22", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetComparingDocumentErrors(t *testing.T) {
	client := &mockComparingDocumentClient{}

	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(sirius.Document{}, errExample)

	server := newMockServer("/comparing-document", CompareDocWithDoc(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/comparing-document?docUid[]=dfef6714-b4fe-44c2-b26e-90dfe3663e95&docUid[]=fweq6452-q4az-44c2-b26e-90dfe3679u22", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetComparingDocumentWhenGetSecondDocumentErrors(t *testing.T) {
	client := &mockComparingDocumentClient{}
	document := sirius.Document{
		ID:        1,
		UUID:      "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		CaseItems: []sirius.Case{{ID: 456}},
	}

	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("DocumentByUUID", mock.Anything, "fweq6452-q4az-44c2-b26e-90dfe3679u22").
		Return(sirius.Document{}, errExample)

	server := newMockServer("/comparing-document", CompareDocWithDoc(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/comparing-document?docUid[]=dfef6714-b4fe-44c2-b26e-90dfe3663e95&docUid[]=fweq6452-q4az-44c2-b26e-90dfe3679u22", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}
