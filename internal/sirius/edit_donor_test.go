package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestEditDonor(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		personData    Person
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			personData: Person{
				ID:                    188,
				UID:                   "7000-0000-0007",
				Salutation:            "Dr",
				Firstname:             "Will",
				Middlenames:           "Oswald",
				Surname:               "Niesborella",
				DateOfBirth:           DateString("1995-07-01"),
				PreviouslyKnownAs:     "Will Macphail",
				AlsoKnownAs:           "Bill",
				AddressLine1:          "47209 Stacey Plain",
				AddressLine2:          "Suite 113",
				AddressLine3:          "Devonburgh",
				Town:                  "Marquardtville",
				County:                "North Carolina",
				Postcode:              "40936",
				Country:               "United States",
				IsAirmailRequired:     true,
				PhoneNumber:           "0841781784",
				Email:                 "docniesborella@mail.test",
				CorrespondenceByPost:  true,
				CorrespondenceByEmail: true,
				CorrespondenceByPhone: false,
				CorrespondenceByWelsh: true,
				ResearchOptOut:        true,
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists").
					UponReceiving("A request to edit a donor").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/donors/188"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"id":                    188,
							"uId":                   "7000-0000-0007",
							"salutation":            "Dr",
							"firstname":             "Will",
							"middlenames":           "Oswald",
							"surname":               "Niesborella",
							"dob":                   "01/07/1995",
							"previousNames":         "Will Macphail",
							"otherNames":            "Bill",
							"addressLine1":          "47209 Stacey Plain",
							"addressLine2":          "Suite 113",
							"addressLine3":          "Devonburgh",
							"town":                  "Marquardtville",
							"county":                "North Carolina",
							"postcode":              "40936",
							"country":               "United States",
							"isAirmailRequired":     true,
							"phoneNumber":           "0841781784",
							"email":                 "docniesborella@mail.test",
							"sageId":                "",
							"correspondenceByPost":  true,
							"correspondenceByEmail": true,
							"correspondenceByPhone": false,
							"correspondenceByWelsh": true,
							"researchOptOut":        true,
							"companyName":           "",
							"companyReference":      "",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.EditDonor(Context{Context: context.Background()}, 188, tc.personData)
				if (tc.expectedError) == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}
