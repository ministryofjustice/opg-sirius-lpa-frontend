package sirius

import (
	"context"
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
					Given("Some document templates exist").
					UponReceiving("A request for document templates").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/templates/lpa"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"DD": dsl.Like(map[string]interface{}{
								"onScreenSummary": dsl.Like("DDONSCREENSUMMARY"),
								"location":        dsl.Like(`lpa\/DD.html.twig`),
								"inserts": dsl.Like(map[string]interface{}{
									"all": dsl.Like(map[string]interface{}{
										"DD1": dsl.Like(map[string]interface{}{
											"onScreenSummary": dsl.Like("DD1LPAINSERTONSCREENSUMMARY"),
											"location":        dsl.Like(`lpa\/inserts\/DD1.html.twig`),
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
							Key:             "all",
							InsertId:        "DD1",
							Location:        `lpa\/inserts\/DD1.html.twig`,
							OnScreenSummary: "DD1LPAINSERTONSCREENSUMMARY",
						},
					},
					Location:        `lpa\/DD.html.twig`,
					OnScreenSummary: "DDONSCREENSUMMARY",
					TemplateId:      "DD",
				},
			},
		},
		{
			name: "OK - template with no inserts",
			setup: func() {
				pact.
					AddInteraction().
					Given("Some document templates exist").
					UponReceiving("A request for document templates (no inserts)").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/templates/lpa"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"CT-bb": dsl.Like(map[string]interface{}{
								"onScreenSummary": dsl.Like("CTBBONSCREENSUMMARY"),
								"location":        dsl.Like(`complaints/CT-bb.html.twig`),
								"inserts":         dsl.Like([]interface{}{}),
							}),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []DocumentTemplateData{
				{
					Inserts:         []Insert(nil),
					Location:        `complaints/CT-bb.html.twig`,
					OnScreenSummary: "CTBBONSCREENSUMMARY",
					TemplateId:      "CT-bb",
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
