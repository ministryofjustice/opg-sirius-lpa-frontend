package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDigitalLpa(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/digital-lpas/M-1234-9876-4567"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: map[string]interface{}{
							"uId": matchers.Like("M-1234-9876-4567"),
							"opg.poas.sirius": map[string]interface{}{
								"id":                 matchers.Like(789),
								"caseSubtype":        matchers.Like("property-and-affairs"),
								"status":             matchers.Like("Draft"),
								"createdDate":        matchers.Term("26/03/2018", `^\d{1,2}/\d{1,2}/\d{4}$`),
								"complaintCount":     matchers.Like(1),
								"investigationCount": matchers.Like(2),
								"taskCount":          matchers.Like(3),
								"warningCount":       matchers.Like(4),
								"objectionCount":     matchers.Like(5),
								"application": map[string]interface{}{
									"donorFirstNames": matchers.Like("Zoraida"),
									"donorLastName":   matchers.Like("Swanberg"),
									"donorDob":        matchers.Term("27/05/1978", `^\d{1,2}/\d{1,2}/\d{4}$`),
									"donorPhone":      matchers.Like("073456249524"),
									"donorEmail":      matchers.Like("zswanberg@host.example"),
									"donorAddress": map[string]interface{}{
										"addressLine1": matchers.Like("Apartment 24"),
										"addressLine2": matchers.Like("Navigation Building"),
										"addressLine3": matchers.Like("90 London Road"),
										"town":         matchers.Like("Birmingham"),
										"postcode":     matchers.Like("B15 4TA"),
										"country":      matchers.Term("GB", `^[A-Z]{2}$`),
									},
									"correspondentFirstNames": matchers.Like("Heath"),
									"correspondentLastName":   matchers.Like("Enstad"),
									"correspondentAddress": map[string]interface{}{
										"addressLine1": matchers.Like("Main Line Bungalow"),
										"addressLine2": matchers.Like("Himmerton Lane"),
										"addressLine3": matchers.Like("Sutton"),
										"town":         matchers.Like("Scarsdale"),
										"postcode":     matchers.Like("S24 7DJ"),
										"country":      matchers.Term("GB", `^[A-Z]{2}$`),
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
		{
			name: "OK2",
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA in statutory waiting period").
					UponReceiving("A request for the digital LPA in statutory waiting period").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/digital-lpas/M-1111-2222-3333"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: map[string]interface{}{
							"uId": matchers.Like("M-1111-2222-3333"),
							"opg.poas.sirius": map[string]interface{}{
								"id":                 matchers.Like(111),
								"caseSubtype":        matchers.Like("property-and-affairs"),
								"status":             matchers.Like("Statutory waiting period"),
								"createdDate":        matchers.Term("26/03/2018", `^\d{1,2}/\d{1,2}/\d{4}$`),
								"dueDate":            matchers.Term("09/04/2018", `^\d{1,2}/\d{1,2}/\d{4}$`),
								"complaintCount":     matchers.Like(0),
								"investigationCount": matchers.Like(0),
								"taskCount":          matchers.Like(0),
								"warningCount":       matchers.Like(0),
								"objectionCount":     matchers.Like(0),
								"application": map[string]interface{}{
									"donorFirstNames": matchers.Like("Lonnie"),
									"donorLastName":   matchers.Like("Jakubowski"),
									"donorDob":        matchers.Term("22/03/1949", `^\d{1,2}/\d{1,2}/\d{4}$`),
									"donorPhone":      matchers.Like("07123456789"),
									"donorEmail":      matchers.Like("Lonnie.Jakubowski@example.com"),
									"donorAddress": map[string]interface{}{
										"addressLine1": matchers.Like("528 Fourth Avenue"),
										"addressLine2": matchers.Like("Lower Kozey Cross"),
										"addressLine3": matchers.Like("East Thiel"),
										"town":         matchers.Like("Ahlen"),
										"postcode":     matchers.Like("YL06 6GF"),
										"country":      matchers.Term("GB", `^[A-Z]{2}$`),
									},
								},
							},
						},
					})
			},
			expectedResponse: DigitalLpa{
				UID: "M-1111-2222-3333",
				SiriusData: SiriusData{
					ID:                 111,
					Subtype:            "property-and-affairs",
					Status:             "Statutory waiting period",
					CreatedDate:        DateString("2018-03-26"),
					DueDate:            DateString("2018-04-09"),
					ComplaintCount:     0,
					InvestigationCount: 0,
					TaskCount:          0,
					WarningCount:       0,
					ObjectionCount:     0,
					Application: Draft{
						DonorFirstNames: "Lonnie",
						DonorLastName:   "Jakubowski",
						DonorDob:        DateString("1949-03-22"),
						DonorAddress: Address{
							Line1:    "528 Fourth Avenue",
							Line2:    "Lower Kozey Cross",
							Line3:    "East Thiel",
							Town:     "Ahlen",
							Postcode: "YL06 6GF",
							Country:  "GB",
						},
						PhoneNumber: "07123456789",
						Email:       "Lonnie.Jakubowski@example.com",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				digitalLpa, err := client.DigitalLpa(Context{Context: context.Background()}, tc.expectedResponse.UID, false)

				assert.Equal(t, tc.expectedResponse, digitalLpa)
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
