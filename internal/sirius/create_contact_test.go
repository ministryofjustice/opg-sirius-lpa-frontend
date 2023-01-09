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
				Salutation:            "Mrs",
				Firstname:             "Pauline",
				Middlenames:           "Suzanne",
				Surname:               "Price",
				CompanyName:           "",
				CompanyReference:      "",
				AddressLine1:          "278 Nicole Lock",
				AddressLine2:          "Toby Court",
				AddressLine3:          "",
				Town:                  "Russellstad",
				County:                "Cumbria",
				Postcode:              "HP19 9BW",
				Country:               "United Kingdom",
				IsAirmailRequired:     true,
				PhoneNumber:           "072345678",
				Email:                 "p.price@uk.test",
				CorrespondenceByPost:  false,
				CorrespondenceByEmail: true,
				CorrespondenceByPhone: true,
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
							"salutation":            "Mrs",
							"firstname":             "Pauline",
							"middlenames":           "Suzanne",
							"surname":               "Price",
							"addressLine1":          "278 Nicole Lock",
							"addressLine2":          "Toby Court",
							"town":                  "Russellstad",
							"county":                "Cumbria",
							"postcode":              "HP19 9BW",
							"country":               "United Kingdom",
							"sageId":                "",
							"isAirmailRequired":     true,
							"phoneNumber":           "072345678",
							"email":                 "p.price@uk.test",
							"correspondenceByPost":  false,
							"correspondenceByEmail": true,
							"correspondenceByPhone": true,
							"correspondenceByWelsh": false,
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
