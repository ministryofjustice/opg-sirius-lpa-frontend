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

func TestEditDates(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to edit the dates").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/lpas/800/edit-dates"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"rejectedDate":     "04/03/2022",
							"cancellationDate": "04/03/2022",
							"dispatchDate":     "04/03/2022",
							"dueDate":          "04/03/2022",
							"invalidDate":      "04/03/2022",
							"paymentDate":      "08/02/2022",
							"receiptDate":      "04/03/2022",
							"registrationDate": "04/03/2022",
							"revokedDate":      "01/01/2023",
							"withdrawnDate":    "04/03/2022",
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

				err := client.EditDates(Context{Context: context.Background()}, 800, "lpa", Dates{
					RejectedDate:     DateString("2022-03-04"),
					CancellationDate: DateString("2022-03-04"),
					DispatchDate:     DateString("2022-03-04"),
					DueDate:          DateString("2022-03-04"),
					InvalidDate:      DateString("2022-03-04"),
					PaymentDate:      DateString("2022-02-08"),
					ReceiptDate:      DateString("2022-03-04"),
					RegistrationDate: DateString("2022-03-04"),
					RevokedDate:      DateString("2023-01-01"),
					WithdrawnDate:    DateString("2022-03-04"),
				})

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
