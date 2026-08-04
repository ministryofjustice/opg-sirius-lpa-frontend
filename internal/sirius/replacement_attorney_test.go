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

func TestCreateReplacementAttorney(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		attorney         Attorney
		expectedResponse Attorney
	}{
		{
			name: "OK",
			attorney: Attorney{
				Person: Person{
					Salutation:   "Mrs",
					Firstname:    "Rosalind",
					Surname:      "Achebe",
					Email:        "r.achebe@example.test",
					AddressLine1: "14 Kestrel Way",
					Town:         "Hirthehaven",
					County:       "Saskatchewan",
					Postcode:     "S7R 9F9",
					Country:      "Canada",
				},
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A case exists for a replacement attorney to be added to").
					UponReceiving("A request to create a replacement attorney").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/persons"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: []interface{}{
							map[string]interface{}{
								"personType":   matchers.String("ReplacementAttorney"),
								"caseId":       matchers.Integer(800),
								"salutation":   matchers.String("Mrs"),
								"firstname":    matchers.String("Rosalind"),
								"surname":      matchers.String("Achebe"),
								"email":        matchers.String("r.achebe@example.test"),
								"addressLine1": matchers.String("14 Kestrel Way"),
								"town":         matchers.String("Hirthehaven"),
								"county":       matchers.String("Saskatchewan"),
								"postcode":     matchers.String("S7R 9F9"),
								"country":      matchers.String("Canada"),
							},
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: []interface{}{
							map[string]interface{}{
								"id":         matchers.Like(321),
								"firstname":  matchers.String("Rosalind"),
								"surname":    matchers.String("Achebe"),
								"personType": matchers.String("Replacement Attorney"),
							},
						},
					})
			},
			expectedResponse: Attorney{Person: Person{ID: 321, Firstname: "Rosalind", Surname: "Achebe"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				attorney, err := client.CreateReplacementAttorney(Context{Context: context.Background()}, 800, tc.attorney)
				assert.Equal(t, tc.expectedResponse, attorney)
				assert.Nil(t, err)
				return nil
			}))
		})
	}
}
