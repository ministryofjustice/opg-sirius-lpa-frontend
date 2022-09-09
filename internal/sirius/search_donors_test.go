package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestSearchDonors(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []Person
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists to be referenced").
					UponReceiving("A search for donors").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/search/persons"),
						Body: dsl.Like(map[string]interface{}{
							"term":        "7000-9999-0001",
							"personTypes": []string{"Donor"},
						}),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"results": dsl.EachLike(map[string]interface{}{
								"id":        dsl.Like(47),
								"uId":       dsl.Like("7000-0000-0003"),
								"firstname": dsl.Like("John"),
								"surname":   dsl.Like("Doe"),
							}, 1),
						}),
					})
			},
			expectedResponse: []Person{
				{
					ID:        47,
					UID:       "7000-0000-0003",
					Firstname: "John",
					Surname:   "Doe",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				donors, err := client.SearchDonors(Context{Context: context.Background()}, "7000-9999-0001")
				assert.Equal(t, tc.expectedResponse, donors)
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

func TestSearchDonorsTooShort(t *testing.T) {
	client := NewClient(http.DefaultClient, "")

	donors, err := client.SearchDonors(Context{Context: context.Background()}, "ad")
	assert.Nil(t, donors)
	assert.Equal(t, fmt.Errorf("Search term must be at least three characters"), err)
}
