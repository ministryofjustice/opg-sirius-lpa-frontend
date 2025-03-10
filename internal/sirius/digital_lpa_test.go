package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
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
		{
			name: "OK2",
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA in statutory waiting period").
					UponReceiving("A request for the digital LPA in statutory waiting period").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/digital-lpas/M-1111-2222-3333"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: map[string]interface{}{
							"uId": dsl.Like("M-1111-2222-3333"),
							"opg.poas.sirius": map[string]interface{}{
								"id":                 dsl.Like(111),
								"caseSubtype":        dsl.Like("property-and-affairs"),
								"status":             dsl.Like("Statutory waiting period"),
								"createdDate":        dsl.Term("26/03/2018", `^\d{1,2}/\d{1,2}/\d{4}$`),
								"dueDate":            dsl.Term("09/04/2018", `^\d{1,2}/\d{1,2}/\d{4}$`),
								"complaintCount":     dsl.Like(0),
								"investigationCount": dsl.Like(0),
								"taskCount":          dsl.Like(0),
								"warningCount":       dsl.Like(0),
								"objectionCount":     dsl.Like(0),
								"application": map[string]interface{}{
									"donorFirstNames": dsl.Like("Lonnie"),
									"donorLastName":   dsl.Like("Jakubowski"),
									"donorDob":        dsl.Term("22/03/1949", `^\d{1,2}/\d{1,2}/\d{4}$`),
									"donorPhone":      dsl.Like("07123456789"),
									"donorEmail":      dsl.Like("Lonnie.Jakubowski@example.com"),
									"donorAddress": map[string]interface{}{
										"addressLine1": dsl.Like("528 Fourth Avenue"),
										"addressLine2": dsl.Like("Lower Kozey Cross"),
										"addressLine3": dsl.Like("East Thiel"),
										"town":         dsl.Like("Ahlen"),
										"postcode":     dsl.Like("YL06 6GF"),
										"country":      dsl.Term("GB", `^[A-Z]{2}$`),
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				digitalLpa, err := client.DigitalLpa(Context{Context: context.Background()}, tc.expectedResponse.UID, false)

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
