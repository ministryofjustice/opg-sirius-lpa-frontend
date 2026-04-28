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

func TestCreateCorrespondent(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		correspondent Correspondent
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			correspondent: Correspondent{
				Person: Person{
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
				}},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending EPA assigned").
					UponReceiving("A request to create a correspondent").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/persons"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: []map[string]interface{}{{
							"caseId":                800,
							"personType":            "Correspondent",
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
						}},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
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

				err := client.CreateCorrespondent(Context{Context: context.Background()}, 800, tc.correspondent)
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

func TestCorrespondentSummary(t *testing.T) {
	correspondent := Correspondent{
		Person: Person{
			Firstname: "Melanie",
			Surname:   "Vanvolkenburg",
		},
	}
	assert.Equal(t, "Melanie Vanvolkenburg", correspondent.Summary())
}

func TestCorrespondentAddressSummary(t *testing.T) {
	correspondent := Correspondent{
		Person: Person{
			AddressLine1: "29737 Andrew Plaza",
			AddressLine2: "Apt. 814",
			Town:         "Hirthehaven",
			County:       "Saskatchewan",
			Postcode:     "S7R 9F9",
			Country:      "Canada",
		},
	}
	assert.Equal(t, "29737 Andrew Plaza, Apt. 814, Hirthehaven, Saskatchewan, S7R 9F9, Canada", correspondent.AddressSummary())
}
