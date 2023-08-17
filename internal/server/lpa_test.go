package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLpaClient struct {
	mock.Mock
}

func (m *mockLpaClient) DigitalLpa(ctx sirius.Context, uid string) (sirius.DigitalLpa, error) {
	args := m.Called(ctx, uid)

	return args.Get(0).(sirius.DigitalLpa), args.Error(1)
}

func TestGetLpa(t *testing.T) {
	digitalLpa := sirius.DigitalLpa{
		UID:     "M-9876-9876-9876",
		Subtype: "hw",
	}

	client := &mockLpaClient{}
	client.
		On("DigitalLpa", mock.Anything, "M-9876-9876-9876").
		Return(digitalLpa, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, lpaData{
			Lpa: digitalLpa,
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}", Lpa(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
