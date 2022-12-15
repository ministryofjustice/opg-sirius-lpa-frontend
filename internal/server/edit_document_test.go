package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEditDocumentClient struct {
	mock.Mock
}

func (m *mockEditDocumentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockEditDocumentClient) Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int) ([]sirius.Document, error) {
	args := m.Called(ctx, caseType, caseId)
	return args.Get(0).([]sirius.Document), args.Error(1)
}

func (m *mockEditDocumentClient) DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func TestGetEditDocument(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "7000"}

			document := sirius.Document{
				ID:   1,
				UUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
			}

			documents := []sirius.Document{
				document,
			}

			client := &mockEditDocumentClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseItem, nil)
			client.
				On("Documents", mock.Anything, sirius.CaseType(caseType), 123).
				Return(documents, nil)
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, editDocumentData{
					Case:      caseItem,
					Documents: documents,
					Document:  document,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&case="+caseType, nil)
			w := httptest.NewRecorder()

			err := EditDocument(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
