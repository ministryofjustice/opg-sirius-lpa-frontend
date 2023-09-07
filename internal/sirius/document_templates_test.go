package sirius

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestDocumentTypes(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []DocumentTemplateData
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for document templates").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/templates/lpa"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"DD": dsl.Like(map[string]interface{}{
								"label": dsl.Like("Donor deceased: Blank template"),
								"inserts": dsl.Like(map[string]interface{}{
									"all": dsl.Like(map[string]interface{}{
										"DD1": dsl.Like(map[string]interface{}{
											"label": dsl.Like("DD1 - Case complete"),
											"order": dsl.Like(0),
										}),
									}),
								}),
							}),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []DocumentTemplateData{
				{
					Inserts: []Insert{
						{
							Key:      "all",
							InsertId: "DD1",
							Label:    "DD1 - Case complete",
							Order:    0,
						},
					},
					TemplateId: "DD",
					Label:      "Donor deceased: Blank template",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				documentTemplateTypes, err := client.DocumentTemplates(Context{Context: context.Background()}, CaseTypeLpa)

				assert.Equal(t, tc.expectedResponse, documentTemplateTypes)
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

func TestHandlingNotifyField(t *testing.T) {
	config := "{\"label\":\"An example letter sent with notify\",\"govukNotify\":true}"
	response := documentTemplateApiResponse{
		"WITH-NOTIFY": json.RawMessage([]byte(config)),
	}

	data, err := response.toDocumentData()

	assert.Nil(t, err)
	assert.Equal(t, "An example letter sent with notify", data[0].Label)
	assert.Equal(t, true, data[0].UsesNotify)
}
