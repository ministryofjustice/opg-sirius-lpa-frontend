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
				Source:   "PHONE",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor with a digital LPA exists").
					UponReceiving("A request to create an additional draft LPA with minimal data").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/donors/234/digital-lpas"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"types":  []string{"personal-welfare"},
							"source": "PHONE",
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: []map[string]interface{}{
							{
								"caseSubtype": dsl.String("personal-welfare"),
								"uId":         dsl.Regex("M-GHIJ-7890-ABCD", `^M(-[0-9A-Z]{4}){3}$`),
							},
						},
					})
			},
			expectedResponse: map[string]string{
				"personal-welfare": "M-GHIJ-7890-ABCD",
			},
		},
		{
			name: "All possible details",
			additionalDraftData: AdditionalDraft{
				CaseType:                []string{"personal-welfare", "property-and-affairs"},
				Source:                  "PHONE",
				CorrespondentFirstNames: "Leon Marius",
				CorrespondentLastName:   "Selden",
				CorrespondentAddress: &Address{
					Line1:    "Nitzsche, Nader And Schuppe",
					Line2:    "6064 Alessandro Plain",
					Line3:    "Pittsfield",
					Town:     "Concord",
					Postcode: "JN2 7UO",
					Country:  "GB",
				},
				CorrespondenceByWelsh:     true,
				CorrespondenceLargeFormat: false,
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor with a digital LPA exists").
					UponReceiving("A request to create an additional draft LPA with all possible data").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/donors/234/digital-lpas"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"types":                   []string{"personal-welfare", "property-and-affairs"},
							"source":                  "PHONE",
							"correspondentFirstNames": "Leon Marius",
							"correspondentLastName":   "Selden",
							"correspondentAddress": map[string]string{
								"addressLine1": "Nitzsche, Nader And Schuppe",
								"addressLine2": "6064 Alessandro Plain",
								"addressLine3": "Pittsfield",
								"town":         "Concord",
								"postcode":     "JN2 7UO",
								"country":      "GB",
							},
							"correspondenceByWelsh": true,
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
							{
								"caseSubtype": dsl.String("property-and-affairs"),
								"uId":         dsl.Regex("M-ABCD-1234-EF56", `^M(-[0-9A-Z]{4}){3}$`),
							},
						},
					})
			},
			expectedResponse: map[string]string{
				"personal-welfare":     "M-GHIJ-7890-KLMN",
				"property-and-affairs": "M-ABCD-1234-EF56",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				person, err := client.CreateAdditionalDraft(Context{Context: context.Background()}, 234, tc.additionalDraftData)
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
