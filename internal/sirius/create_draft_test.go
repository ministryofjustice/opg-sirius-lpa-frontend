package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCreateDraft(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/digital-lpas"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
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
				Email:                     "c.morneault@somehost.example",
				CorrespondenceByWelsh:     true,
				CorrespondenceLargeFormat: false,
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a Modernised LPA user").
					UponReceiving("A request to create a draft LPA with all possible data").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/digital-lpas"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
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
							"donorEmail":            "c.morneault@somehost.example",
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
				"personal-welfare":  "M-GHIJ-7890-KLMN",
				"property-and-affairs": "M-ABCD-1234-EF56",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				person, err := client.CreateDraft(Context{Context: context.Background()}, tc.draftData)
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
