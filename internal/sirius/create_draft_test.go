package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDraft(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		draftData        Draft
		setup            func()
		expectedResponse map[string]string
		expectedError    func(int) error
	}{
		{
			name: "Minimal",
			draftData: Draft{
				CaseType:        []string{"personal-welfare"},
				Source:          "PHONE",
				DonorFirstNames: "Coleen",
				DonorLastName:   "Morneault",
				DonorDob:        DateString("1952-04-08"),
				DonorAddress: Address{
					Line1:    "Fluke House",
					Town:     "South Bend",
					Postcode: "AI1 6VW",
					Country:  "GB",
				},
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a Modernised LPA user").
					UponReceiving("A request to create a draft LPA with minimal data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/digital-lpas"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"types":           []string{"personal-welfare"},
							"source":          "PHONE",
							"donorFirstNames": "Coleen",
							"donorLastName":   "Morneault",
							"donorDob":        "08/04/1952",
							"donorAddress": map[string]string{
								"addressLine1": "Fluke House",
								"town":         "South Bend",
								"postcode":     "AI1 6VW",
								"country":      "GB",
							},
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: []map[string]interface{}{
							{
								"caseSubtype": matchers.String("personal-welfare"),
								"uId":         matchers.Regex("M-GHIJ-7890-KLMN", `^M(-[0-9A-Z]{4}){3}$`),
							},
						},
					})
			},
			expectedResponse: map[string]string{
				"personal-welfare": "M-GHIJ-7890-KLMN",
			},
		},
		{
			name: "All possible details",
			draftData: Draft{
				CaseType:        []string{"personal-welfare", "property-and-affairs"},
				Source:          "PHONE",
				DonorFirstNames: "Coleen Stephanie",
				DonorLastName:   "Morneault",
				DonorDob:        DateString("1952-04-08"),
				DonorAddress: Address{
					Line1:    "Fluke House",
					Line2:    "Summit",
					Line3:    "Houston",
					Town:     "South Bend",
					Postcode: "AI1 6VW",
					Country:  "GB",
				},
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
				PhoneNumber:               "07893932118",
				Email:                     "c.morneault@example.com",
				CorrespondenceByWelsh:     true,
				CorrespondenceLargeFormat: false,
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a Modernised LPA user").
					UponReceiving("A request to create a draft LPA with all possible data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/digital-lpas"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"types":           []string{"personal-welfare", "property-and-affairs"},
							"source":          "PHONE",
							"donorFirstNames": "Coleen Stephanie",
							"donorLastName":   "Morneault",
							"donorDob":        "08/04/1952",
							"donorAddress": map[string]string{
								"addressLine1": "Fluke House",
								"addressLine2": "Summit",
								"addressLine3": "Houston",
								"town":         "South Bend",
								"postcode":     "AI1 6VW",
								"country":      "GB",
							},
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
							"donorPhone":            "07893932118",
							"donorEmail":            "c.morneault@example.com",
							"correspondenceByWelsh": true,
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: []map[string]interface{}{
							{
								"caseSubtype": matchers.String("personal-welfare"),
								"uId":         matchers.Regex("M-GHIJ-7890-KLMN", `^M(-[0-9A-Z]{4}){3}$`),
							},
							{
								"caseSubtype": matchers.String("property-and-affairs"),
								"uId":         matchers.Regex("M-ABCD-1234-EF56", `^M(-[0-9A-Z]{4}){3}$`),
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				person, err := client.CreateDraft(Context{Context: context.Background()}, tc.draftData)
				if (tc.expectedError) == nil {
					assert.Equal(t, tc.expectedResponse, person)
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}
