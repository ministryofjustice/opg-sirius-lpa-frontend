package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type mockManageAttorneysClient struct {
	mock.Mock
}

func (m *mockManageAttorneysClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func TestManageAttorneys(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-1111-2222-3333",
		},
	}

	client := &mockManageAttorneysClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(caseSummary, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, manageAttorneysData{
			CaseSummary: caseSummary,
			Error:       sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/manage-attorneys", ManageAttorneys(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/manage-attorneys", nil)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestManageAttorneyGetCaseSummaryFails(t *testing.T) {
	caseSummary := sirius.CaseSummary{}

	client := &mockManageAttorneysClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(caseSummary, expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, manageAttorneysData{}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/manage-attorneys", ManageAttorneys(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/manage-attorneys", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
}

func TestManageAttorneyTemplateErrors(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-1111-2222-3333",
		},
	}

	client := &mockManageAttorneysClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(caseSummary, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, manageAttorneysData{
			CaseSummary: caseSummary,
			Error:       sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(expectedError)

	server := newMockServer("/lpa/{uid}/manage-attorneys", ManageAttorneys(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/manage-attorneys", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
}
