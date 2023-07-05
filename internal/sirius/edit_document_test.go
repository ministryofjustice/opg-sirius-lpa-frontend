package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEditDocument(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					Given("I have an lpa with a draft document").
					UponReceiving("A request to edit the document").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/documents/dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: dsl.Like(map[string]interface{}{
							"content": dsl.String("<p>Edited test content</p>"),
						}),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"id":                  dsl.Like(1),
							"uuid":                dsl.String("dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
							"type":                dsl.String("Draft"),
							"friendlyDescription": dsl.String("Dr Consuela Aysien - LPA perfect + reg due date: applicant"),
							"createdDate":         dsl.String(`15/12/2022 13:41:04`),
							"direction":           dsl.String("Outgoing"),
							"filename":            dsl.String("LP-A.pdf"),
							"mimeType":            dsl.String(`application\/pdf`),
							"childCount":          dsl.Like(0),
							"systemType":          dsl.String("LP-A"),
							"content":             dsl.String("<p>Edited test content</p>"),
						}),
					})
			},
			expectedResponse: Document{
				ID:                  1,
				UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:                "Draft",
				FriendlyDescription: "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
				CreatedDate:         `15/12/2022 13:41:04`,
				Direction:           "Outgoing",
				MimeType:            `application\/pdf`,
				SystemType:          "LP-A",
				FileName:            "LP-A.pdf",
				Content:             "<p>Edited test content</p>",
				ChildCount:          0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				document, err := client.EditDocument(Context{Context: context.Background()}, "dfef6714-b4fe-44c2-b26e-90dfe3663e95", "<p>Edited test content</p>")

				assert.Equal(t, tc.expectedResponse, document)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
