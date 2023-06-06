package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCreateContact(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/non-case-contacts"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
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
							"otherField":            true,
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: map[string]interface{}{
							"id":  dsl.Like(771),
							"uId": dsl.Like("7000-0000-2688"),
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				person, err := client.CreateContact(Context{Context: context.Background()}, tc.personData)
				if (tc.expectedError) == nil {
					assert.Equal(t, tc.expectedPerson, person)
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
