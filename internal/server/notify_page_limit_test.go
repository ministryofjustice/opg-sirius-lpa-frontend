package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockNotifyPageLimitClient struct {
	mock.Mock
}

func (m *mockNotifyPageLimitClient) DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func TestGetNotifyPageLimitDocument(t *testing.T) {
	document := sirius.Document{
		FriendlyDescription: "LP Letter",
		ID:                  1,
		UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		SystemType:          "LP-LETTER",
		Type:                sirius.TypeSave,
	}

	notifyWarning := "Document has been published named " + document.FriendlyDescription +
		". The published document has more than 10 pages and will not be sent automatically via GOV.UK,Notify. It must be manually uploaded to Xerox for sending."

	client := &mockNotifyPageLimitClient{}
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)

	template := &mockTemplate{}
	templateData := notifyPageLimitData{
		CaseID:              8,
		CaseType:            "LPA",
		CaseUUID:            "700000000021",
		CreatedDocumentUUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		DefaultDocumentUUID: "ad8bbdc6-163d-49f9-9cac-a8141310e409",
		DonorID:             678,
		Error:               sirius.ValidationError{},
		IsPartial:           false,
		NotifyWarning:       notifyWarning,
		XSRFToken:           "",
	}

	template.
		On("Func", mock.Anything, templateData).
		Return(nil)

	server := newMockServer("/notify-error/{donorId}/{createdDocumentUUID}/{caseType}/{caseID}/{caseUUID}", NotifyPageLimit(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/notify-error/678/dfef6714-b4fe-44c2-b26e-90dfe3663e95/LPA/8/700000000021?defaultDocumentUUID=ad8bbdc6-163d-49f9-9cac-a8141310e409", nil)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, template)
}
