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
							"uId":                dsl.String("M-1234-9876-4567"),
							"caseSubtype":        dsl.String("pfa"),
							"status":             dsl.String("Draft"),
							"createdDate":        dsl.Term("26/03/2018", `^\d{1,2}/\d{1,2}/\d{4}$`),
							"complaintCount":     dsl.Like(1),
							"investigationCount": dsl.Like(2),
							"taskCount":          dsl.Like(3),
							"warningCount":       dsl.Like(4),
							"application": map[string]interface{}{
								"donorName":  dsl.String("Zoraida Swanberg"),
								"donorDob":   dsl.Term("27/05/1978", `^\d{1,2}/\d{1,2}/\d{4}$`),
								"donorPhone": dsl.String("073456249524"),
								"donorEmail": dsl.String("zswanberg@host.example"),
								"donorAddress": map[string]interface{}{
									"addressLine1": dsl.String("Apartment 24"),
									"addressLine2": dsl.String("Navigation Building"),
									"addressLine3": dsl.String("90 London Road"),
									"town":         dsl.String("Birmingham"),
									"postcode":     dsl.String("B15 4TA"),
									"country":      dsl.Term("GB", `^[A-Z]{2}$`),
								},
								"correspondentName": dsl.String("Heath Enstad"),
								"correspondentAddress": map[string]interface{}{
									"addressLine1": dsl.String("Main Line Bungalow"),
									"addressLine2": dsl.String("Himmerton Lane"),
									"addressLine3": dsl.String("Sutton"),
									"town":         dsl.String("Scarsdale"),
									"postcode":     dsl.String("S24 7DJ"),
									"country":      dsl.Term("GB", `^[A-Z]{2}$`),
								},
							},
						},
					})
			},
			expectedResponse: DigitalLpa{
				UID:                "M-1234-9876-4567",
				Subtype:            "pfa",
				Status:             "Draft",
				CreatedDate:        DateString("2018-03-26"),
				ComplaintCount:     1,
				InvestigationCount: 2,
				TaskCount:          3,
				WarningCount:       4,
				Application: Draft{
					DonorName: "Zoraida Swanberg",
					DonorDob:  DateString("1978-05-27"),
					DonorAddress: Address{
						Line1:    "Apartment 24",
						Line2:    "Navigation Building",
						Line3:    "90 London Road",
						Town:     "Birmingham",
						Postcode: "B15 4TA",
						Country:  "GB",
					},
					CorrespondentName: "Heath Enstad",
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
