package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestEditComplaint(t *testing.T) {
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
					Given("A complaint exists").
					UponReceiving("A request to edit the complaint").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/complaints/986"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: dsl.Like(map[string]interface{}{
							"category":            "01",
							"description":         "This is seriously bad",
							"receivedDate":        "05/04/2022",
							"severity":            "Major",
							"subCategory":         "07",
							"complainantCategory": "LPA_DONOR",
							"origin":              "PHONE",
							"summary":             "This and that",
							"resolution":          "complaint upheld",
							"resolutionInfo":      "We did stuff",
							"resolutionDate":      "07/06/2022",
						}),
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

				err := client.EditComplaint(Context{Context: context.Background()}, 986, Complaint{
					Category:            "01",
					Description:         "This is seriously bad",
					ReceivedDate:        DateString("2022-04-05"),
					Severity:            "Major",
					SubCategory:         "07",
					ComplainantCategory: "LPA_DONOR",
					Origin:              "PHONE",
					Summary:             "This and that",
					Resolution:          "complaint upheld",
					ResolutionInfo:      "We did stuff",
					ResolutionDate:      DateString("2022-06-07"),
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
