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

func TestBankHolidays(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	pact.
		AddInteraction().
		Given("").
		UponReceiving("A request to get bank holidays").
		WithCompleteRequest(consumer.Request{
			Method: http.MethodGet,
			Path:   matchers.String("/lpa-api/v1/dates/bank-holidays"),
		}).
		WithCompleteResponse(consumer.Response{
			Status:  http.StatusOK,
			Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
			Body: matchers.Like(map[string]interface{}{
				"2025": map[string]interface{}{
					"New Year": matchers.String("2025-01-01T00:00:00+00:00"),
				},
			}),
		})

	expectedResponse := BankHolidays{
		"2025": {
			"New Year": "2025-01-01T00:00:00+00:00",
		},
	}

	t.Run("Get bank holidays", func(t *testing.T) {
		assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

			bankHolidays, err := client.BankHolidays(Context{Context: context.Background()})

			assert.Equal(t, expectedResponse, bankHolidays)
			assert.Nil(t, err)

			return nil
		}))
	})

}
