package sirius

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAnomaliesHttpClient struct {
	mock.Mock
}

func (m *mockAnomaliesHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestAnomaliesForDigitalLpaSuccess(t *testing.T) {
	pact := newPact()

	defer pact.Teardown()

	expectedResponse := []Anomaly{}

	pact.
		AddInteraction().
		Given("A digital LPA with UID LPA M-QWQW-QTQT-WERT and anomalies exists").
		UponReceiving("A request for the anomalies for a digital LPA").
		WithRequest(dsl.Request{
			Method: http.MethodGet,
			Path:   dsl.String("/lpa-api/v1/digital-lpas/M-QWQW-QTQT-WERT/anomalies"),
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusOK,
			Body: dsl.Like(map[string]interface{}{
				"uid":       dsl.Like("M-QWQW-QTQT-WERT"),
				"anomalies": dsl.EachLike(map[string]interface{}{}, 0),
			}),
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
		})

	assert.Nil(t, pact.Verify(func() error {
		client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

		anomalies, err := client.AnomaliesForDigitalLpa(Context{Context: context.Background()}, "M-QWQW-QTQT-WERT")

		assert.Equal(t, expectedResponse, anomalies)
		assert.Nil(t, err)

		return nil
	}))
}

func TestAnomaliesForDigitalLpaFail(t *testing.T) {
	mockClient := &mockAnomaliesHttpClient{}
	mockClient.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Networking issue"))

	client := NewClient(mockClient, "http://localhost")
	_, err := client.ProgressIndicatorsForDigitalLpa(Context{Context: context.Background()}, "M-QEQE-EEEE-QQQE")

	assert.Equal(t, "Networking issue", err.Error())
}
