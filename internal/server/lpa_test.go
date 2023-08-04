package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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

	var err error
	mux := chi.NewRouter()
	mux.HandleFunc("/lpa/{uid}", func(w http.ResponseWriter, r *http.Request) {
		err = Lpa(client, template.Func)(w, r)
	})

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876", nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, template)
}
