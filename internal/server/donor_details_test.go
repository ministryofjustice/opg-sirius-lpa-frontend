package server

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDonorDetailsClient struct {
	mock.Mock
}

func (m *mockDonorDetailsClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func TestDonorDetailsFail(t *testing.T) {
	expectedError := errors.New("network error")

	client := &mockDonorDetailsClient{}
	client.
		On("Person", mock.Anything, 1).
		Return(sirius.Person{}, expectedError)

	template := &mockTemplate{}

	server := newMockServer("/donor/{donorId}/details", DonorDetails(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/donor/1/details", nil)
	_, err := server.serve(req)

	assert.Error(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestDonorDetailsSuccess(t *testing.T) {
	expectedDonor := sirius.Person{
		ID:                    123,
		UID:                   "7000-0000-0001",
		Salutation:            "Mr",
		Firstname:             "John",
		Middlenames:           "Paul",
		Surname:               "Smith",
		DateOfBirth:           sirius.DateString("1950-01-01"),
		AddressLine1:          "123 Main Street",
		AddressLine2:          "Flat 1",
		Town:                  "London",
		County:                "Greater London",
		Postcode:              "SW1A 1AA",
		Country:               "England",
		PhoneNumber:           "020 7946 0958",
		Email:                 "john.smith@example.com",
		CorrespondenceByPost:  true,
		CorrespondenceByEmail: true,
		CorrespondenceByPhone: false,
	}

	client := &mockDonorDetailsClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(expectedDonor, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, DonorDetailsData{
			Donor: expectedDonor,
		}).
		Return(nil)

	server := newMockServer("/donor/{donorId}/details", DonorDetails(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/donor/123/details", nil)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, template)
}
