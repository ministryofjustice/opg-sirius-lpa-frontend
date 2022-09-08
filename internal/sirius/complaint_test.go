package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestComplaint(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Complaint
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A complaint exists").
					UponReceiving("A request for the complaint").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/complaints/986"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"category":     dsl.String("01"),
							"description":  dsl.String("This is seriously bad"),
							"receivedDate": dsl.String("05/04/2022"),
							"severity":     dsl.String("Major"),
							"subCategory":  dsl.String("07"),
							"summary":      dsl.String("This and that"),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: Complaint{
				Category:     "01",
				Description:  "This is seriously bad",
				ReceivedDate: DateString("2022-04-05"),
				Severity:     "Major",
				SubCategory:  "07",
				Summary:      "This and that",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				complaint, err := client.Complaint(Context{Context: context.Background()}, 986)

				assert.Equal(t, tc.expectedResponse, complaint)
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
