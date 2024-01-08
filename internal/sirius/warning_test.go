package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWarningHttpClient struct {
	mock.Mock
}

func (m *mockWarningHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestWarningsForCase(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []Warning
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with a warning").
					UponReceiving("A request for the warnings for a case").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/cases/990/warnings"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"id":           dsl.Like(9901),
							"dateAdded":    dsl.String("08/01/2023"),
							"warningType":  dsl.String("Donor Deceased"),
							"warningText":  dsl.String("Donor died"),
							"caseItems":    dsl.EachLike(map[string]interface{}{
								"uId":      dsl.String("M-TTTT-RRRR-EEEE"),
								"caseType": dsl.String("DIGITAL_LPA"),
							}, 1),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []Warning{
				Warning{
					ID: 9901,
					DateAdded: "08/01/2023",
					WarningType: "Donor Deceased",
					WarningText: "Donor died",
					CaseItems: []Case{
						Case{
							UID: "M-TTTT-RRRR-EEEE",
							CaseType: "DIGITAL_LPA",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				warnings, err := client.WarningsForCase(Context{Context: context.Background()}, 990)

				assert.Equal(t, tc.expectedResponse, warnings)
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
