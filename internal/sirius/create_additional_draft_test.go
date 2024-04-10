package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateAdditionalDraft(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name                string
		additionalDraftData AdditionalDraft
		setup               func()
		expectedResponse    map[string]string
		expectedError       func(int) error
	}{
		{
			name: "Minimal",
			additionalDraftData: AdditionalDraft{
				CaseType: []string{"personal-welfare"},
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a Modernised LPA user").
					UponReceiving("A request to create an additional draft LPA with minimal data").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/donors/188/digital-lpas"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"types": []string{"personal-welfare"},
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: []map[string]interface{}{
							{
								"caseSubtype": dsl.String("personal-welfare"),
								"uId":         dsl.Regex("M-GHIJ-7890-KLMN", `^M(-[0-9A-Z]{4}){3}$`),
							},
						},
					})
			},
			expectedResponse: map[string]string{
				"personal-welfare": "M-GHIJ-7890-KLMN",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				person, err := client.CreateAdditionalDraft(Context{Context: context.Background()}, 188, tc.additionalDraftData)
				if (tc.expectedError) == nil {
					assert.Equal(t, tc.expectedResponse, person)
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
