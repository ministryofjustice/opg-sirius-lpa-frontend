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

func (m *mockGetDocuments) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func TestGetDocuments(t *testing.T) {
	documents := []sirius.Document{
		{
			ID:   1,
			UUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
			Type: sirius.TypeSave,
		},
	}

	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      676,
				Subtype: "hw",
			},
		},
		TaskList: []sirius.Task{},
	}

	client := &mockGetDocuments{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(caseSummary, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseType("lpa"), 676, []string{}, []string{sirius.TypeDraft, sirius.TypePreview}).
		Return(documents, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getDocumentsData{
			CaseSummary: caseSummary,
			Documents:   documents,
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/documents", GetDocuments(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876/documents", nil)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetDocumentsWhenFailureOnGetDigitalLpa(t *testing.T) {
	client := &mockGetDocuments{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(sirius.CaseSummary{}, expectedError)

	server := newMockServer("/lpa/{uid}/documents", GetDocuments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDocumentsWhenFailureOnGetDocuments(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      1532,
				Subtype: "hw",
			},
		},
		TaskList: []sirius.Task{},
	}

	client := &mockGetDocuments{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(caseSummary, nil)
	client.
		On("Documents", mock.Anything, sirius.CaseType("lpa"), 1532, []string{}, []string{sirius.TypeDraft, sirius.TypePreview}).
		Return([]sirius.Document{}, expectedError)

	server := newMockServer("/lpa/{uid}/documents", GetDocuments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDocumentsWhenFailureOnGetCaseSummary(t *testing.T) {
	client := &mockGetDocuments{}
	client.
		On("CaseSummary", mock.Anything, "M-A876-A876-A876").
		Return(sirius.CaseSummary{}, expectedError)

	server := newMockServer("/lpa/{uid}/documents", GetDocuments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-A876-A876-A876/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
