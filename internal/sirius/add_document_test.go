package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddDocument(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedError    func(int) error
		expectedResponse Document
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to add a document to the case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/lpas/800/documents"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"caseId":          matchers.Like(800),
							"type":            matchers.String("Save"),
							"filename":        matchers.String("LP-A.pdf"),
							"correspondentId": matchers.Like(189),
							"systemType":      matchers.String("LP-A"),
							"content":         matchers.String("Test content"),
						}),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"id":                  matchers.Like(1),
							"uuid":                matchers.String("dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
							"type":                matchers.String("Save"),
							"friendlyDescription": matchers.String("Dr Consuela Aysien - LPA perfect + reg due date: applicant"),
							"createdDate":         matchers.String(`15/12/2022 13:41:04`),
							"direction":           matchers.String("Outgoing"),
							"filename":            matchers.String("LP-A.pdf"),
							"mimeType":            matchers.String(`application\/pdf`),
							"childCount":          matchers.Like(0),
							"systemType":          matchers.String("LP-A"),
							"content":             matchers.String("Test content"),
						}),
					})
			},
			expectedResponse: Document{
				ID:                  1,
				UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:                "Save",
				FriendlyDescription: "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
				CreatedDate:         `15/12/2022 13:41:04`,
				Direction:           "Outgoing",
				MimeType:            `application\/pdf`,
				SystemType:          "LP-A",
				FileName:            "LP-A.pdf",
				Content:             "Test content",
				ChildCount:          0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				document, err := client.AddDocument(Context{Context: context.Background()}, 800, Document{
					ID:                  1,
					UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
					Type:                "Draft",
					FriendlyDescription: "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
					CreatedDate:         `15/12/2022 13:41:04`,
					Direction:           "Outgoing",
					MimeType:            `application\/pdf`,
					SystemType:          "LP-A",
					FileName:            "LP-A.pdf",
					Content:             "Test content",
					ChildCount:          0,
				}, TypeSave)

				assert.Equal(t, tc.expectedResponse, document)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}
