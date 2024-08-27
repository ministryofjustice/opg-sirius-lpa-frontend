package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDigitalLpa(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse DigitalLpa
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA exists").
					UponReceiving("A request for the digital LPA").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/digital-lpas/M-1234-9876-4567"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: map[string]interface{}{
							"uId": dsl.Like("M-1234-9876-4567"),
							"opg.poas.sirius": map[string]interface{}{
								"id":                 dsl.Like(789),
								"caseSubtype":        dsl.Like("property-and-affairs"),
								"status":             dsl.Like("Draft"),
								"createdDate":        dsl.Term("26/03/2018", `^\d{1,2}/\d{1,2}/\d{4}$`),
								"complaintCount":     dsl.Like(1),
								"investigationCount": dsl.Like(2),
								"taskCount":          dsl.Like(3),
								"warningCount":       dsl.Like(4),
								"objectionCount":     dsl.Like(5),
								"application": map[string]interface{}{
									"donorFirstNames": dsl.Like("Zoraida"),
									"donorLastName":   dsl.Like("Swanberg"),
									"donorDob":        dsl.Term("27/05/1978", `^\d{1,2}/\d{1,2}/\d{4}$`),
									"donorPhone":      dsl.Like("073456249524"),
									"donorEmail":      dsl.Like("zswanberg@host.example"),
									"donorAddress": map[string]interface{}{
										"addressLine1": dsl.Like("Apartment 24"),
										"addressLine2": dsl.Like("Navigation Building"),
										"addressLine3": dsl.Like("90 London Road"),
										"town":         dsl.Like("Birmingham"),
										"postcode":     dsl.Like("B15 4TA"),
										"country":      dsl.Term("GB", `^[A-Z]{2}$`),
									},
									"correspondentFirstNames": dsl.Like("Heath"),
									"correspondentLastName":   dsl.Like("Enstad"),
									"correspondentAddress": map[string]interface{}{
										"addressLine1": dsl.Like("Main Line Bungalow"),
										"addressLine2": dsl.Like("Himmerton Lane"),
										"addressLine3": dsl.Like("Sutton"),
										"town":         dsl.Like("Scarsdale"),
										"postcode":     dsl.Like("S24 7DJ"),
										"country":      dsl.Term("GB", `^[A-Z]{2}$`),
									},
								},
							},
						},
					})
			},
			expectedResponse: DigitalLpa{
				UID: "M-1234-9876-4567",
				SiriusData: SiriusData{
					ID:                 789,
					Subtype:            "property-and-affairs",
					Status:             "Draft",
					CreatedDate:        DateString("2018-03-26"),
					ComplaintCount:     1,
					InvestigationCount: 2,
					TaskCount:          3,
					WarningCount:       4,
					ObjectionCount:     5,
					Application: Draft{
						DonorFirstNames: "Zoraida",
						DonorLastName:   "Swanberg",
						DonorDob:        DateString("1978-05-27"),
						DonorAddress: Address{
							Line1:    "Apartment 24",
							Line2:    "Navigation Building",
							Line3:    "90 London Road",
							Town:     "Birmingham",
							Postcode: "B15 4TA",
							Country:  "GB",
						},
						CorrespondentFirstNames: "Heath",
						CorrespondentLastName:   "Enstad",
						CorrespondentAddress: &Address{
							Line1:    "Main Line Bungalow",
							Line2:    "Himmerton Lane",
							Line3:    "Sutton",
							Town:     "Scarsdale",
							Postcode: "S24 7DJ",
							Country:  "GB",
						},
						PhoneNumber: "073456249524",
						Email:       "zswanberg@host.example",
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

				digitalLpa, err := client.DigitalLpa(Context{Context: context.Background()}, "M-1234-9876-4567")

				assert.Equal(t, tc.expectedResponse, digitalLpa)
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
