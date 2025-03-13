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

func TestCreateContact(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		personData     Person
		setup          func()
		expectedPerson Person
		expectedError  func(int) error
	}{
		{
			name: "OK",
			personData: Person{
				Salutation:            "Prof",
				Firstname:             "Melanie",
				Middlenames:           "Josefina",
				Surname:               "Vanvolkenburg",
				PreviouslyKnownAs:     "",
				AlsoKnownAs:           "",
				AddressLine1:          "29737 Andrew Plaza",
				AddressLine2:          "Apt. 814",
				AddressLine3:          "Gislasonside",
				Town:                  "Hirthehaven",
				County:                "Saskatchewan",
				Postcode:              "S7R 9F9",
				Country:               "",
				SageId:                "",
				IsAirmailRequired:     true,
				PhoneNumber:           "072345678",
				Email:                 "m.vancolkenburg@ca.test",
				CorrespondenceByPost:  false,
				CorrespondenceByEmail: true,
				CorrespondenceByPhone: true,
				CorrespondenceByWelsh: false,
				CompanyName:           "",
				CompanyReference:      "",
				ResearchOptOut:        false,
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a Lay Team user").
					UponReceiving("A request to create a contact").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/non-case-contacts"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"salutation":            "Prof",
							"firstname":             "Melanie",
							"middlenames":           "Josefina",
							"surname":               "Vanvolkenburg",
							"dob":                   nil,
							"previousNames":         "",
							"otherNames":            "",
							"addressLine1":          "29737 Andrew Plaza",
							"addressLine2":          "Apt. 814",
							"addressLine3":          "Gislasonside",
							"town":                  "Hirthehaven",
							"county":                "Saskatchewan",
							"postcode":              "S7R 9F9",
							"country":               "",
							"sageId":                "",
							"isAirmailRequired":     true,
							"phoneNumber":           "072345678",
							"email":                 "m.vancolkenburg@ca.test",
							"correspondenceByPost":  false,
							"correspondenceByEmail": true,
							"correspondenceByPhone": true,
							"correspondenceByWelsh": false,
							"companyName":           "",
							"companyReference":      "",
							"researchOptOut":        false,
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: map[string]interface{}{
							"id":  matchers.Like(771),
							"uId": matchers.Like("7000-0000-2688"),
						},
					})
			},
			expectedPerson: Person{
				ID:  771,
				UID: "7000-0000-2688",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				person, err := client.CreateContact(Context{Context: context.Background()}, tc.personData)
				if (tc.expectedError) == nil {
					assert.Equal(t, tc.expectedPerson, person)
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}
