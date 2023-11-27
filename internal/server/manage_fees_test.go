package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockManageFeesClient struct {
	mock.Mock
}

func (m *mockManageFeesClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func TestGetManageFees(t *testing.T) {
	tests := []struct {
		scenario          string
		caseItem          sirius.Case
		expectedReturnUrl string
	}{
		{
			scenario: "Digital LPA",
			caseItem: sirius.Case{
				ID:       71,
				UID:      "M-AAAA-BBBB-DDDD",
				CaseType: "DIGITAL_LPA",
			},
			expectedReturnUrl: "/lpa/M-AAAA-BBBB-DDDD/payments",
		},
		{
			scenario: "Non-digital LPA",
			caseItem: sirius.Case{
				ID:  81,
				UID: "7000-0000-0021",
			},
			expectedReturnUrl: "/payments/81",
		},
	}

	for _, test := range tests {
		client := &mockManageFeesClient{}
		client.
			On("Case", mock.Anything, test.caseItem.ID).
			Return(test.caseItem, nil)

		template := &mockTemplate{}
		template.
			On("Func", mock.Anything, manageFeesData{
				Case:      test.caseItem,
				ReturnUrl: test.expectedReturnUrl,
			}).
			Return(nil)

		r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/?id=%d", test.caseItem.ID), nil)
		w := httptest.NewRecorder()

		err := ManageFees(client, template.Func)(w, r)
		resp := w.Result()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mock.AssertExpectationsForObjects(t, client, template)
	}
}
