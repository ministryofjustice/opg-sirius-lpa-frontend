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

func TestUpdateFeeDecision(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	pact.
		AddInteraction().
		Given("I have a case with declined fee reductions").
		UponReceiving("a request to edit a fee decision").
		WithCompleteRequest(consumer.Request{
			Method: http.MethodPut,
			Path:   matchers.String("/lpa-api/v1/cases/8521/fee-decisions/1245"),
			Headers: matchers.MapMatcher{
				"Content-Type": matchers.String("application/json"),
			},
			Body: matchers.Like(map[string]interface{}{
				"decisionType":   "DECLINED_EXEMPTION",
				"decisionReason": "Insufficient evidence",
				"decisionDate":   "18/10/2023",
			}),
		}).
		WithCompleteResponse(consumer.Response{
			Status: http.StatusOK,
			Body: matchers.Like(map[string]interface{}{
				"id":             matchers.Like(1245),
				"caseId":         matchers.Like(8521),
				"decisionType":   matchers.String("DECLINED_EXEMPTION"),
				"decisionDate":   matchers.String("18/10/2023"),
				"decisionReason": matchers.String("Insufficient evidence"),
			}),
			Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
		})

	assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
		client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

		err := client.UpdateFeeDecision(
			Context{Context: context.Background()},
			8521,
			1245,
			"DECLINED_EXEMPTION",
			"Insufficient evidence",
			DateString("2023-10-18"),
		)

		assert.Nil(t, err)
		return nil
	}))
}
