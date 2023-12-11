package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGetDocuments struct {
	mock.Mock
}

func (m *mockGetDocuments) Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int, docTypes []string, notDocTypes []string) ([]sirius.Document, error) {
	args := m.Called(ctx, caseType, caseId, docTypes, notDocTypes)
	return args.Get(0).([]sirius.Document), args.Error(1)
}

func (m *mockGetDocuments) DigitalLpa(ctx sirius.Context, uid string) (sirius.DigitalLpa, error) {
	args := m.Called(ctx, uid)

	return args.Get(0).(sirius.DigitalLpa), args.Error(1)
}

func (m *mockGetDocuments) TasksForCase(ctx sirius.Context, id int) ([]sirius.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.Task), args.Error(1)
}

func TestGetDocuments(t *testing.T) {
	documents := []sirius.Document{
		{
			ID:   1,
			UUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
			Type: sirius.TypeSave,
		},
	}

	digitalLpa := sirius.DigitalLpa{
		ID:      676,
		UID:     "M-9876-9876-9876",
		Subtype: "hw",
	}

	client := &mockGetDocuments{}
	client.
		On("DigitalLpa", mock.Anything, "M-9876-9876-9876").
		Return(digitalLpa, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseType("lpa"), 676, []string{}, []string{sirius.TypeDraft, sirius.TypePreview}).
		Return(documents, nil)
	client.
		On("TasksForCase", mock.Anything, 676).
		Return([]sirius.Task{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getDocumentsData{
			Lpa:       digitalLpa,
			Documents: documents,
			TaskList:  []sirius.Task{},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/documents", GetDocuments(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876/documents", nil)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetPaymentsWhenFailureOnGetDigitalLpa(t *testing.T) {
	client := &mockGetDocuments{}
	client.
		On("DigitalLpa", mock.Anything, "M-9876-9876-9876").
		Return(sirius.DigitalLpa{}, expectedError)

	server := newMockServer("/lpa/{uid}/documents", GetDocuments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPaymentsWhenFailureOnGetDocuments(t *testing.T) {
	digitalLpa := sirius.DigitalLpa{
		ID:      1532,
		UID:     "M-9876-9876-9876",
		Subtype: "hw",
	}

	client := &mockGetDocuments{}
	client.
		On("DigitalLpa", mock.Anything, "M-9876-9876-9876").
		Return(digitalLpa, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseType("lpa"), 1532, []string{}, []string{sirius.TypeDraft, sirius.TypePreview}).
		Return([]sirius.Document{}, expectedError)
	client.
		On("TasksForCase", mock.Anything, 1532).
		Return([]sirius.Task{}, nil)

	server := newMockServer("/lpa/{uid}/documents", GetDocuments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
