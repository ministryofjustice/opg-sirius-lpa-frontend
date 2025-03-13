package sirius

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProgressIndicatorsHttpClient struct {
	mock.Mock
}

func (m *mockProgressIndicatorsHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestProgressIndicatorsForDigitalLpaSuccess(t *testing.T) {
	pact, err := newPact2()
	assert.NoError(t, err)

	expectedResponse := []ProgressIndicator{
		{"FEES", "IN_PROGRESS"},
	}

	pact.
		AddInteraction().
		Given("A digital LPA with UID LPA M-QEQE-EEEE-WERT and a fees progress indicator with status 'In progress' exists").
		UponReceiving("A request for the progress indicators for a digital LPA").
		WithCompleteRequest(consumer.Request{
			Method: http.MethodGet,
			Path:   matchers.String("/lpa-api/v1/digital-lpas/M-QEQE-EEEE-WERT/progress-indicators"),
		}).
		WithCompleteResponse(consumer.Response{
			Status: http.StatusOK,
			Body: matchers.Like(map[string]interface{}{
				"uid": matchers.Like("M-QEQE-EEEE-WERT"),
				"progressIndicators": matchers.EachLike(map[string]interface{}{
					"indicator": matchers.Like("FEES"),
					"status":    matchers.Like("IN_PROGRESS"),
				}, 0),
			}),
			Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
		})

	assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
		client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

		progressIndicators, err := client.ProgressIndicatorsForDigitalLpa(Context{Context: context.Background()}, "M-QEQE-EEEE-WERT")

		assert.Equal(t, expectedResponse, progressIndicators)
		assert.Nil(t, err)

		return nil
	}))
}

func TestGetApplicationProgressProgressIndicatorsFail(t *testing.T) {
	mockClient := &mockProgressIndicatorsHttpClient{}
	mockClient.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Networking issue"))

	client := NewClient(mockClient, "http://localhost")
	_, err := client.ProgressIndicatorsForDigitalLpa(Context{Context: context.Background()}, "M-QEQE-EEEE-QQQE")

	assert.Equal(t, "Networking issue", err.Error())
}
