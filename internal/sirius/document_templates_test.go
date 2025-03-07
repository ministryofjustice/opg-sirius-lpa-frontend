package sirius

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestDocumentTypes(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/templates/lpa"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"DD": matchers.Like(map[string]interface{}{
								"label": matchers.Like("Donor deceased: Blank template"),
								"inserts": matchers.Like(map[string]interface{}{
									"all": matchers.Like(map[string]interface{}{
										"DD1": matchers.Like(map[string]interface{}{
											"label": matchers.Like("DD1 - Case complete"),
											"order": matchers.Like(0),
										}),
									}),
								}),
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				documentTemplateTypes, err := client.DocumentTemplates(Context{Context: context.Background()}, CaseTypeLpa)

				assert.Equal(t, tc.expectedResponse, documentTemplateTypes)
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
