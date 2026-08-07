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

func TestCreateNotifiedPerson(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		notifiedPerson NotifiedPerson
		setup          func()
		expectedError  func(int) error
	}{
		{
			name: "OK",
			notifiedPerson: NotifiedPerson{
				Person: Person{
					Salutation:        "Prof",
					Firstname:         "Melanie",
					Middlenames:       "Josefina",
					Surname:           "Vanvolkenburg",
					AddressLine1:      "29737 Andrew Plaza",
					AddressLine2:      "Apt. 814",
					AddressLine3:      "Gislasonside",
					Town:              "Hirthehaven",
					County:            "Saskatchewan",
					Postcode:          "S7R 9F9",
					Country:           "Canada",
					IsAirmailRequired: true,
				},
				NoticeGivenDate: DateString("2022-04-05"),
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to create a notified person").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/persons"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: []map[string]interface{}{{
							"salutation":            "Prof",
							"firstname":             "Melanie",
							"middlenames":           "Josefina",
							"surname":               "Vanvolkenburg",
							"addressLine1":          "29737 Andrew Plaza",
							"addressLine2":          "Apt. 814",
							"addressLine3":          "Gislasonside",
							"otherNames":            "",
							"companyName":           "",
							"companyReference":      "",
							"correspondenceByEmail": false,
							"correspondenceByPhone": false,
							"correspondenceByPost":  false,
							"correspondenceByWelsh": false,
							"town":                  "Hirthehaven",
							"county":                "Saskatchewan",
							"postcode":              "S7R 9F9",
							"country":               "Canada",
							"dob":                   nil,
							"dateOfDeath":           nil,
							"email":                 "",
							"isAirmailRequired":     true,
							"noticeGivenDate":       "05/04/2022",
							"personType":            "NotifiedPerson",
							"phoneNumber":           "",
							"previousNames":         "",
							"researchOptOut":        false,
							"sageId":                "",
							"caseId":                800,
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

				err := client.CreateNotifiedPerson(Context{Context: context.Background()}, 800, tc.notifiedPerson)
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

func TestUpdateNotifiedPerson(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		notifiedPerson NotifiedPerson
		setup          func()
		expectedError  func(int) error
	}{
		{
			name: "OK",
			notifiedPerson: NotifiedPerson{
				Person: Person{
					ID:                123,
					Salutation:        "Prof",
					Firstname:         "Melanie",
					Middlenames:       "Josefina",
					Surname:           "Vanvolkenburg",
					AddressLine1:      "29737 Andrew Plaza",
					AddressLine2:      "Apt. 814",
					AddressLine3:      "Gislasonside",
					Town:              "Hirthehaven",
					County:            "Saskatchewan",
					Postcode:          "S7R 9F9",
					Country:           "Canada",
					IsAirmailRequired: true,
				},
				NoticeGivenDate: DateString("2022-04-05"),
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa with a notified person").
					UponReceiving("A request to update that notified person").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/notified-persons/123"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"id":                    123,
							"salutation":            "Prof",
							"firstname":             "Melanie",
							"middlenames":           "Josefina",
							"surname":               "Vanvolkenburg",
							"addressLine1":          "29737 Andrew Plaza",
							"addressLine2":          "Apt. 814",
							"addressLine3":          "Gislasonside",
							"otherNames":            "",
							"companyName":           "",
							"companyReference":      "",
							"correspondenceByEmail": false,
							"correspondenceByPhone": false,
							"correspondenceByPost":  false,
							"correspondenceByWelsh": false,
							"town":                  "Hirthehaven",
							"county":                "Saskatchewan",
							"postcode":              "S7R 9F9",
							"country":               "Canada",
							"dob":                   nil,
							"dateOfDeath":           nil,
							"email":                 "",
							"isAirmailRequired":     true,
							"noticeGivenDate":       "05/04/2022",
							"phoneNumber":           "",
							"previousNames":         "",
							"researchOptOut":        false,
							"sageId":                "",
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

				err := client.UpdateNotifiedPerson(Context{Context: context.Background()}, 123, tc.notifiedPerson)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}
