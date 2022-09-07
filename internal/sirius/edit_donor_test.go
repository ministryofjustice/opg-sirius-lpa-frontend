package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestEditDonor(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name          string
		personData    Person
		setup         func()
		cookies       []*http.Cookie
		expectedError func(int) error
	}{
		{
			name: "OK",
			personData: Person{
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
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/donors/189"),
						Body: map[string]interface{}{
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
						},
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.EditDonor(getContext(tc.cookies), 189, tc.personData)
				if (tc.expectedError) == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
