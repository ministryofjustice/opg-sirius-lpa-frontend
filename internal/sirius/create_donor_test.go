package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCreateDonor(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name           string
		personData     Person
		setup          func()
		cookies        []*http.Cookie
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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/donors"),
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
							"isAirmailRequired":     true,
							"phoneNumber":           "072345678",
							"email":                 "m.vancolkenburg@ca.test",
							"correspondenceByPost":  false,
							"correspondenceByEmail": true,
							"correspondenceByPhone": true,
							"correspondenceByWelsh": false,
							"researchOptOut":        true,
						},
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: map[string]interface{}{
							"id":  dsl.Like(771),
							"uId": dsl.Like("7000-0290-0192"),
						},
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedPerson: Person{
				ID:  771,
				UID: "7000-0290-0192",
			},
		},
		{
			name: "Unauthorized",
			personData: Person{
				Firstname: "Guillermo",
				Surname:   "Prothero",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a Lay Team user").
					UponReceiving("A request to create a donor without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/donors"),
						Body: map[string]interface{}{
							"firstname":             "Guillermo",
							"surname":               "Prothero",
							"isAirmailRequired":     false,
							"correspondenceByPost":  false,
							"correspondenceByEmail": false,
							"correspondenceByPhone": false,
							"correspondenceByWelsh": false,
							"researchOptOut":        false},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:   http.StatusUnauthorized,
					URL:    fmt.Sprintf("http://localhost:%d/lpa-api/v1/donors", port),
					Method: http.MethodPost,
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				person, err := client.CreateDonor(getContext(tc.cookies), tc.personData)
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
