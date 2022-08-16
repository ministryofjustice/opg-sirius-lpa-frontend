package sirius

import (
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPayment(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		cookies          []*http.Cookie
		expectedError    func(int) error
		expectedResponse []Payment
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa which has been paid for").
					UponReceiving("A request for the payments").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/cases/9/payments"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"id":      dsl.Like(2),
							"case_id": dsl.Like(9),
							"source": dsl.Like(map[string]interface{}{
								"name":  "Phone",
								"value": "PHONE",
							}),
							"amount":      dsl.Like(4100),
							"paymentdate": dsl.String("2022-08-23T00:00:00+00:00"),
							"type": dsl.Like(map[string]interface{}{
								"name":  "Card",
								"value": "CARD",
							}),
							"createddate":  dsl.String("2022-08-24T00:00:00+00:00"),
							"locked":       dsl.Like(false),
							"createdby_id": dsl.Like(123),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: []Payment{
				{
					ID:     2,
					CaseID: 9,
					Source: PaymentSource{
						Name:  "Phone",
						Value: "PHONE",
					},
					Amount:      FeeString("41.00"),
					PaymentDate: DateString("2022-08-23"),
					Type: TypeOfPayment{
						Name:  "Card",
						Value: "CARD",
					},
					CreatedDate: DateString("2022-08-24"),
					Locked:      false,
					CreatedByID: 123,
				},
			},
		},
		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa which has been paid for").
					UponReceiving("A request for the payments without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/cases/9/payments"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:   http.StatusUnauthorized,
					URL:    fmt.Sprintf("http://localhost:%d/lpa-api/v1/cases/9/payments", port),
					Method: http.MethodGet,
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				payments, err := client.Payments(getContext(tc.cookies), 9)

				assert.Equal(t, tc.expectedResponse, payments)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
