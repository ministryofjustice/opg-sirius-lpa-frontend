package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockViewDocumentClient struct {
	mock.Mock
}

func (m *mockViewDocumentClient) DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func TestGetViewDocument(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			document := sirius.Document{
				ID:         1,
				UUID:       "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				SystemType: "LP-LETTER",
				Type:       sirius.TypeSave,
			}

			client := &mockViewDocumentClient{}
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)

			template := &mockTemplate{}
			templateData := viewDocumentData{
				Document: document,
			}

			template.
				On("Func", mock.Anything, templateData).
				Return(nil)

			server := newMockServer("/view-document/{uuid}", ViewDocument(client, template.Func))

			req, _ := http.NewRequest(http.MethodGet, "/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
			_, err := server.serve(req)

			assert.Nil(t, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetViewDocumentWhenCaseErrors(t *testing.T) {
	client := &mockViewDocumentClient{}

	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(sirius.Document{}, errExample)

	server := newMockServer("/view-document/{uuid}", ViewDocument(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}
