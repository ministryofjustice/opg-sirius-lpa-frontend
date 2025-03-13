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

func TestCreateDonor(t *testing.T) {
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
				DateOfBirth:           DateString("1978-04-19"),
				PreviouslyKnownAs:     "Colton Bacman",
				AlsoKnownAs:           "Mel",
				AddressLine1:          "29737 Andrew Plaza",
				AddressLine2:          "Apt. 814",
				AddressLine3:          "Gislasonside",
				Town:                  "Hirthehaven",
				County:                "Saskatchewan",
				Postcode:              "S7R 9F9",
				Country:               "Canada",
				IsAirmailRequired:     true,
				PhoneNumber:           "072345678",
				Email:                 "m.vancolkenburg@ca.test",
				CorrespondenceByPost:  false,
				CorrespondenceByEmail: true,
				CorrespondenceByPhone: true,
				CorrespondenceByWelsh: false,
				ResearchOptOut:        true,
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a Lay Team user").
					UponReceiving("A request to create a donor").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/donors"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"salutation":            "Prof",
							"firstname":             "Melanie",
							"middlenames":           "Josefina",
							"surname":               "Vanvolkenburg",
							"dob":                   "19/04/1978",
							"previousNames":         "Colton Bacman",
							"otherNames":            "Mel",
							"addressLine1":          "29737 Andrew Plaza",
							"addressLine2":          "Apt. 814",
							"addressLine3":          "Gislasonside",
							"town":                  "Hirthehaven",
							"county":                "Saskatchewan",
							"postcode":              "S7R 9F9",
							"country":               "Canada",
							"sageId":                "",
							"isAirmailRequired":     true,
							"phoneNumber":           "072345678",
							"email":                 "m.vancolkenburg@ca.test",
							"correspondenceByPost":  false,
							"correspondenceByEmail": true,
							"correspondenceByPhone": true,
							"correspondenceByWelsh": false,
							"researchOptOut":        true,
							"companyName":           "",
							"companyReference":      "",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: map[string]interface{}{
							"id":  matchers.Like(771),
							"uId": matchers.Like("7000-0290-0192"),
						},
					})
			},
			expectedPerson: Person{
				ID:  771,
				UID: "7000-0290-0192",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				person, err := client.CreateDonor(Context{Context: context.Background()}, tc.personData)
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
