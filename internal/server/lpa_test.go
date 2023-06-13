package server

import (
	"net/http"
	"net/http/httptest"
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

	req, _ := http.NewRequest(http.MethodGet, "/lpa?uid=M-9876-9876-9876", nil)
	w := httptest.NewRecorder()

	err := Lpa(client, template.Func)(w, req)
	assert.Nil(t, err)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
