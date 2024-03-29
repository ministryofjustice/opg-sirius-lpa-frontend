package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestEditDates(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/lpas/800/edit-dates"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"rejectedDate":     "04/03/2022",
							"cancellationDate": "04/03/2022",
							"dispatchDate":     "04/03/2022",
							"dueDate":          "04/03/2022",
							"invalidDate":      "04/03/2022",
							"receiptDate":      "04/03/2022",
							"registrationDate": "04/03/2022",
							"withdrawnDate":    "04/03/2022",
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.EditDates(Context{Context: context.Background()}, 800, "lpa", Dates{
					RejectedDate:     DateString("2022-03-04"),
					CancellationDate: DateString("2022-03-04"),
					DispatchDate:     DateString("2022-03-04"),
					DueDate:          DateString("2022-03-04"),
					InvalidDate:      DateString("2022-03-04"),
					ReceiptDate:      DateString("2022-03-04"),
					RegistrationDate: DateString("2022-03-04"),
					WithdrawnDate:    DateString("2022-03-04"),
				})

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
