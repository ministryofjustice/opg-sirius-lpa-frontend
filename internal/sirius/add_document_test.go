package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddDocument(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name          string
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to add a document to the case").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/lpas/800/documents"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: dsl.Like(map[string]interface{}{
							"id":                  dsl.Like(1),
							"uuid":                dsl.String("dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
							"type":                dsl.String("Draft"),
							"friendlyDescription": dsl.String("Dr Consuela Aysien - __LPAONSCREENSUMMARY__"),
							"createdDate":         dsl.String(`15\/12\/2022 13:41:04`),
							"direction":           dsl.String("Outgoing"),
							"fileName":            dsl.String("LP-A.pdf"),
							"mimeType":            dsl.String(`application\/pdf`),
							"correspondent": dsl.Like(map[string]interface{}{
								"id":                    dsl.Like(1),
								"salutation":            dsl.String("Mrs"),
								"firstname":             dsl.String("Consuela"),
								"middlenames":           dsl.String(""),
								"surname":               dsl.String("Aysien"),
								"sageId":                dsl.String(""),
								"isAirmailRequired":     true,
								"phoneNumber":           dsl.String("072345678"),
								"email":                 dsl.String("c.aysien@ca.test"),
								"correspondenceByPost":  false,
								"correspondenceByEmail": true,
								"correspondenceByPhone": true,
								"correspondenceByWelsh": false,
							}),
							"childCount": dsl.Like(0),
							"systemType": dsl.String("LP-A"),
							"content":    dsl.String("Test content"),
						}),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body:    dsl.Like(map[string]interface{}{"id": dsl.Integer()}),
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.AddDocument(Context{Context: context.Background()}, 800, Document{
					ID:                  1,
					UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
					Type:                "Draft",
					FriendlyDescription: "Dr Consuela Aysien - __LPAONSCREENSUMMARY__",
					CreatedDate:         `15\/12\/2022 13:41:04`,
					Direction:           "Outgoing",
					MimeType:            `application\/pdf`,
					SystemType:          "LP-A",
					FileName:            "LP-A.pdf",
					Content:             "Test content",
					Correspondent:       Person{ID: 1, Firstname: "Consuela", Surname: "Aysien"},
					ChildCount:          0,
				})

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
