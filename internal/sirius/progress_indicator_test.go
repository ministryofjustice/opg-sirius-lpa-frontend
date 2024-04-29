package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestProgressIndicatorsForDigitalLpa(t *testing.T) {
	pact := newPact()

	defer pact.Teardown()

	expectedResponse := []ProgressIndicator{}

	pact.
		AddInteraction().
		Given("A digital LPA with UID LPA M-QEQE-EEEE-WERT and a fees progress indicator with status 'In progress' exists").
		UponReceiving("A request for the progress indicators for a digital LPA").
		WithRequest(dsl.Request{
			Method: http.MethodGet,
			Path:   dsl.String("/lpa-api/v1/digital-lpas/M-QEQE-EEEE-WERT/progress-indicators"),
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusOK,
			Body: dsl.Like(map[string]interface{}{
				"uid": dsl.Like("M-QEQE-EEEE-WERT"),
				"progressIndicators": dsl.EachLike(map[string]interface{}{
					"indicator": dsl.Like("FEES"),
					"status":    dsl.Like("IN_PROGRESS"),
				}, 0),
			}),
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
		})

	assert.Nil(t, pact.Verify(func() error {
		client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

		progressIndicators, err := client.ProgressIndicatorsForDigitalLpa(Context{Context: context.Background()}, "M-QEQE-EEEE-WERT")

		assert.Equal(t, expectedResponse, progressIndicators)
		assert.Nil(t, err)

		return nil
	}))
}
