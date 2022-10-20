package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateInvestigation(t *testing.T) {
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
					Given("I am an investigations user with a pending case assigned").
					UponReceiving("A request to create an investigation on the case").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/lpas/800/investigations"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: dsl.Like(map[string]interface{}{
							"investigationTitle":        dsl.String("Test Investigation"),
							"additionalInformation":     dsl.String("This is an investigation"),
							"type":                      dsl.String("Priority"),
							"investigationReceivedDate": dsl.String("05/04/2022"),
						}),
					}).
					WillRespondWith(dsl.Response{Status: http.StatusCreated})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.CreateInvestigation(Context{Context: context.Background()}, 800, CaseTypeLpa, Investigation{
					Title:        "Test Investigation",
					Information:  "This is an investigation",
					Type:         "Priority",
					DateReceived: DateString("2022-04-05"),
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
