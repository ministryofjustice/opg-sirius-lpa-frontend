package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestCreateLpa(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedError    func(int) error
		lpa              Lpa
		expectedResponse Lpa
	}{
		{
			name: "OK",
			lpa: Lpa{
				Case: Case{
					ReceiptDate:          DateString("2015-03-05"),
					SubType:              "pfa",
					ExpectedPaymentTotal: 0,
					Status:               shared.CaseStatusTypePending,
				},
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists").
					UponReceiving("A request to create an LPA").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/donors/189/lpas"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"receiptDate":          matchers.Term("05/03/2015", `^\d{1,2}/\d{1,2}/\d{4}$`),
							"caseSubtype":          matchers.String("pfa"),
							"expectedPaymentTotal": matchers.Integer(0),
							"status":               matchers.String("Pending"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: map[string]interface{}{
							"id":  matchers.Like(1),
							"uId": matchers.Like("7000-0000-0000"),
						},
					})
			},
			expectedResponse: Lpa{Case: Case{ID: 1, UID: "7000-0000-0000"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				lpa, err := client.CreateLpa(Context{Context: context.Background()}, 189, tc.lpa)
				assert.Equal(t, lpa, tc.expectedResponse)
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

func TestUpdateLpa(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		expectedError func(int) error
		lpa           Lpa
	}{
		{
			name: "OK",
			lpa: Lpa{
				Case: Case{
					ReceiptDate:          DateString("2015-03-05"),
					SubType:              "pfa",
					ExpectedPaymentTotal: 0,
					Status:               shared.CaseStatusTypePending,
				},
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending LPA assigned").
					UponReceiving("A request to update the LPA").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/lpas/800"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"receiptDate":          matchers.Term("05/03/2015", `^\d{1,2}/\d{1,2}/\d{4}$`),
							"caseSubtype":          matchers.String("pfa"),
							"expectedPaymentTotal": matchers.Integer(0),
							"status":               matchers.String("Pending"),
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

				err := client.UpdateLpa(Context{Context: context.Background()}, 800, tc.lpa)
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

func TestLpa(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Lpa
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending LPA assigned").
					UponReceiving("A request for the LPA").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/800"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":          800,
							"uId":         matchers.String("7000-0000-0000"),
							"caseType":    matchers.String("LPA"),
							"caseSubtype": matchers.String("pfa"),
							"status":      matchers.String("Pending"),
							"donor": matchers.Like(map[string]interface{}{
								"id": matchers.Like(189),
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Lpa{Case: Case{ID: 800, UID: "7000-0000-0000", CaseType: "LPA", SubType: "pfa", Status: shared.CaseStatusTypePending, Donor: &Person{ID: 189}}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				lpa, err := client.Lpa(Context{Context: context.Background()}, 800)

				assert.Equal(t, tc.expectedResponse, lpa)
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
