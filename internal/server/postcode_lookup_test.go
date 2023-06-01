package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPostcodeLookupClient struct {
	mock.Mock
}

func (m *mockPostcodeLookupClient) PostcodeLookup(ctx sirius.Context, postcode string) ([]sirius.PostcodeLookupAddress, error) {
	args := m.Called(ctx, postcode)
	return args.Get(0).([]sirius.PostcodeLookupAddress), args.Error(1)
}

func TestGetPostcodeLookup(t *testing.T) {
	expectedAddresses := []sirius.PostcodeLookupAddress{{
		Line1:       "17 Some Road",
		Line2:       "Testingsborough",
		Town:        "Teston",
		Postcode:    "SW1A 0AA",
		Description: "17 Some Road, Testingsborough, Teston",
	}}

	client := &mockPostcodeLookupClient{}
	client.
		On("PostcodeLookup", mock.Anything, "SW1A 0AA").
		Return(expectedAddresses, nil)

	req, _ := http.NewRequest(http.MethodGet, "/?postcode=SW1A 0AA", nil)

	w := httptest.NewRecorder()
	err := SearchPostcode(client)(w, req)

	assert.Nil(t, err)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var addresses []sirius.PostcodeLookupAddress
	_ = json.NewDecoder(resp.Body).Decode(&addresses)

	assert.Equal(t, expectedAddresses, addresses)
}

func TestGetPostcodeLookupWhenError(t *testing.T) {
	client := &mockPostcodeLookupClient{}
	client.
		On("PostcodeLookup", mock.Anything, "SW1A 0AA").
		Return([]sirius.PostcodeLookupAddress{}, expectedError)

	req, _ := http.NewRequest(http.MethodGet, "/?postcode=SW1A 0AA", nil)

	w := httptest.NewRecorder()
	err := SearchPostcode(client)(w, req)

	assert.Equal(t, expectedError, err)
}
