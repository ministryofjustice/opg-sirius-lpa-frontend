package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type mockManageAttorneysClient struct {
	mock.Mock
}

func (m *mockManageAttorneysClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func TestGetManageAttorneys(t *testing.T) {
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

func TestGetManageAttorneysGetCaseSummaryFails(t *testing.T) {
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

func TestGetManageAttorneysTemplateErrors(t *testing.T) {
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

func TestPostManageAttorneysInvalidData(t *testing.T) {
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
			AttorneyAction: "",
			CaseSummary:    caseSummary,
			Error: sirius.ValidationError{Field: sirius.FieldErrors{
				"attorneyAction": {"reason": "Please select an option to manage attorneys."},
			}},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/manage-attorneys", ManageAttorneys(client, template.Func))

	form := url.Values{}

	req, _ := http.NewRequest(
		http.MethodPost,
		"/lpa/M-1111-2222-3333/manage-attorneys",
		strings.NewReader(form.Encode()),
	)
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostManageAttorneysValidData(t *testing.T) {
	testCases := map[string]struct {
		attorneyAction      string
		expectedRedirectUrl string
	}{
		"Remove an attorney": {
			attorneyAction:      "remove-an-attorney",
			expectedRedirectUrl: "/lpa/M-1111-2222-3333/remove-an-attorney",
		},
		"Enable replacement attorney": {
			attorneyAction:      "enable-replacement-attorney",
			expectedRedirectUrl: "/lpa/M-1111-2222-3333/enable-replacement-attorney",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
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

			server := newMockServer("/lpa/{uid}/manage-attorneys", ManageAttorneys(client, template.Func))

			form := url.Values{
				"attorneyAction": {tc.attorneyAction},
			}

			req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/manage-attorneys", strings.NewReader(form.Encode()))
			req.Header.Add("Content-Type", formUrlEncoded)
			_, err := server.serve(req)

			assert.Equal(t, RedirectError(tc.expectedRedirectUrl), err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
