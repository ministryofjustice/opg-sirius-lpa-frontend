package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetNotifyPageLimitDocument(t *testing.T) {
	notifyWarning := "Document has been published but has not been posted via GOV.UK,Notify.\n" +
		"The published document has more than 10 pages and will not be posted automatically via GOV.UK,Notify. It must be manually uploaded to Xerox for posting."

	template := &mockTemplate{}
	templateData := notifyPageLimitData{
		CaseID:              8,
		CaseType:            "LPA",
		CaseUUID:            "700000000021",
		CreatedDocumentUUID: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		DefaultDocumentUUID: "ad8bbdc6-163d-49f9-9cac-a8141310e409",
		DonorID:             678,
		IsPartial:           false,
		NotifyWarning:       notifyWarning,
	}

	template.
		On("Func", mock.Anything, templateData).
		Return(nil)

	server := newMockServer("/notify-error/{donorId}/{createdDocumentUUID}/{caseType}/{caseID}/{caseUUID}", NotifyPageLimit(template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/notify-error/678/dfef6714-b4fe-44c2-b26e-90dfe3663e95/LPA/8/700000000021?defaultDocumentUUID=ad8bbdc6-163d-49f9-9cac-a8141310e409", nil)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, template)
}
