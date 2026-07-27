package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSiriusHeaderCaseInfoClient struct {
	mock.Mock
}

func (m *mockSiriusHeaderCaseInfoClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func TestGetSiriusHeaderCaseInfo(t *testing.T) {
	caseItem := sirius.Case{
		ID:            1,
		UID:           "7000-0000-0001",
		CaseRecNumber: "CR123",
	}
	client := &mockSiriusHeaderCaseInfoClient{}
	client.
		On("Case", mock.Anything, 1).
		Return(caseItem, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderCaseInfoData{
			CaseID: 1,
			Case:   caseItem,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1", nil)
	w := httptest.NewRecorder()

	err := SiriusHeaderCaseInfo(client, template.Func)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSiriusHeaderCaseInfoBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockSiriusHeaderCaseInfoClient{}

			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := SiriusHeaderCaseInfo(client, nil)(PageVars{}, w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetSiriusHeaderCaseInfoWhenCaseErrors(t *testing.T) {
	client := &mockSiriusHeaderCaseInfoClient{}
	client.
		On("Case", mock.Anything, 1).
		Return(sirius.Case{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1", nil)
	w := httptest.NewRecorder()

	err := SiriusHeaderCaseInfo(client, nil)(PageVars{}, w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}
