package server

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateLpaClient struct {
	mock.Mock
}

func (m *mockCreateLpaClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockCreateLpaClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func TestGetCreateLpa(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createLpaData{
					DonorId:   123,
					DonorName: "Firstname Surname",
					Title:     "Create an LPA",
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
			w := httptest.NewRecorder()

			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}

			err := CreateLpa(client, template.Func, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateLpaEdit(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
	client.
		On("Case", mock.Anything, 456).
		Return(sirius.Case{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createLpaData{
			DonorId:   123,
			DonorName: "Firstname Surname",
			Title:     "Edit LPA",
			CaseId:    456,
			CaseItem:  sirius.Case{},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=456", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateLpaBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{}, nil)

			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := CreateLpa(client, nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestCreateLpaWhenPersonErrors(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil, nil)(w, r)

	assert.Equal(t, err, errExample)
	mock.AssertExpectationsForObjects(t, client)
}

func TestCreateLpaWhenCaseErrors(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Case", mock.Anything, 456).
		Return(sirius.Case{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=456", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil, nil)(w, r)

	assert.Equal(t, err, errExample)
	mock.AssertExpectationsForObjects(t, client)
}
