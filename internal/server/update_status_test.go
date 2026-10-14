package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUpdateStatusClient struct {
	mock.Mock
}

func (m *mockUpdateStatusClient) CreateNote(ctx sirius.Context, entityID int, entityType sirius.EntityType, noteType, name, description string, file *sirius.NoteFile) error {
	return m.Called(ctx, entityID, entityType, noteType, name, description, file).Error(0)
}

func (m *mockUpdateStatusClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockUpdateStatusClient) EditCase(ctx sirius.Context, caseID int, caseType sirius.CaseType, caseData sirius.Case) error {
	return m.Called(ctx, caseID, caseType, caseData).Error(0)
}

func (m *mockUpdateStatusClient) AvailableStatuses(ctx sirius.Context, caseID int, caseType sirius.CaseType) ([]string, error) {
	args := m.Called(ctx, caseID, caseType)
	return args.Get(0).([]string), args.Error(1)
}

func TestGetUpdateStatus(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "700700"}

			client := &mockUpdateStatusClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseItem, nil)

			client.
				On("AvailableStatuses", mock.Anything, 123, sirius.CaseType(caseType)).
				Return([]string{"Cancelled", "Withdrawn"}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, updateStatusData{
					Entity:            caseType + " 700700",
					AvailableStatuses: []string{"Cancelled", "Withdrawn"},
					CaseID:            123,
					CaseType:          caseType,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&case="+caseType, nil)
			w := httptest.NewRecorder()

			err := UpdateStatus(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetUpdateStatusNoID(t *testing.T) {
	testCases := map[string]string{
		"no-id":    "/?case=lpa",
		"no-case":  "/?id=123",
		"bad-case": "/?id=123&case=person",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := UpdateStatus(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetUpdateStatusWhenCaseErrors(t *testing.T) {
	caseItem := sirius.Case{CaseType: "PFA", UID: "700700"}

	client := &mockUpdateStatusClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := UpdateStatus(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetUpdateStatusWhenAvailableStatusesErrors(t *testing.T) {
	caseItem := sirius.Case{CaseType: "PFA", UID: "700700"}

	client := &mockUpdateStatusClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)

	client.
		On("AvailableStatuses", mock.Anything, 123, sirius.CaseTypeLpa).
		Return([]string{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := UpdateStatus(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetUpdateStatusWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{CaseType: "PFA", UID: "700700"}

	client := &mockUpdateStatusClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)

	client.
		On("AvailableStatuses", mock.Anything, 123, sirius.CaseTypeLpa).
		Return([]string{"Cancelled", "Withdrawn"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, updateStatusData{
			Entity:            "PFA 700700",
			AvailableStatuses: []string{"Cancelled", "Withdrawn"},
			CaseID:            123,
			CaseType:          "lpa",
		}).
		Return(errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := UpdateStatus(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateStatus(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "700700"}

			client := &mockUpdateStatusClient{}
			client.
				On("EditCase", mock.Anything, 123, sirius.CaseType(caseType), sirius.Case{
					Status: shared.CaseStatusTypeWithdrawn,
				}).
				Return(nil)

			client.
				On("Case", mock.Anything, 123).
				Return(caseItem, nil)

			client.
				On("AvailableStatuses", mock.Anything, 123, sirius.CaseType(caseType)).
				Return([]string{"Cancelled", "Withdrawn"}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, updateStatusData{
					Success:           true,
					Entity:            caseType + " 700700",
					AvailableStatuses: []string{"Cancelled", "Withdrawn"},
					NewStatus:         "Withdrawn",
					CaseID:            123,
					CaseType:          caseType,
				}).
				Return(nil)

			form := url.Values{
				"status": {"Withdrawn"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateStatus(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			client.AssertNotCalled(t, "CreateNote", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostUpdateStatusWithNotes(t *testing.T) {
	testCases := map[string]sirius.EntityType{
		"lpa": sirius.EntityTypeLpa,
		"epa": sirius.EntityTypeEpa,
	}

	for caseType, noteEntityType := range testCases {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "700700"}

			client := &mockUpdateStatusClient{}
			client.
				On("EditCase", mock.Anything, 123, sirius.CaseType(caseType), sirius.Case{
					Status: shared.CaseStatusTypeWithdrawn,
				}).
				Return(nil)

			client.
				On("CreateNote", mock.Anything, 123, noteEntityType, "Status change - Notes", "Status changed to Withdrawn", "Case note details", (*sirius.NoteFile)(nil)).
				Return(nil)

			client.
				On("Case", mock.Anything, 123).
				Return(caseItem, nil)

			client.
				On("AvailableStatuses", mock.Anything, 123, sirius.CaseType(caseType)).
				Return([]string{"Cancelled", "Withdrawn"}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, updateStatusData{
					Success:           true,
					Entity:            caseType + " 700700",
					AvailableStatuses: []string{"Cancelled", "Withdrawn"},
					NewStatus:         "Withdrawn",
					CaseID:            123,
					CaseType:          caseType,
					Notes:             "Case note details",
				}).
				Return(nil)

			form := url.Values{
				"status": {"Withdrawn"},
				"notes":  {"Case note details"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateStatus(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostUpdateStatusWhenUpdateStatusErrors(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "700700"}

	client := &mockUpdateStatusClient{}
	client.
		On("EditCase", mock.Anything, 123, sirius.CaseTypeLpa, sirius.Case{
			Status: shared.CaseStatusTypeWithdrawn,
		}).
		Return(errExample)

	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)

	client.
		On("AvailableStatuses", mock.Anything, 123, sirius.CaseTypeLpa).
		Return([]string{"Cancelled", "Withdrawn"}, nil)

	form := url.Values{
		"status": {"Withdrawn"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&case=lpa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateStatus(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetUpdateStatusHtmx(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "700700"}

	client := &mockUpdateStatusClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)

	client.
		On("AvailableStatuses", mock.Anything, 123, sirius.CaseTypeLpa).
		Return([]string{"Cancelled", "Withdrawn"}, nil)

	partialTemplate := &mockTemplate{}
	partialTemplate.
		On("Func", mock.Anything, updateStatusData{
			Entity:            "lpa 700700",
			AvailableStatuses: []string{"Cancelled", "Withdrawn"},
			CaseID:            123,
			CaseType:          "lpa",
			IsPartial:         true,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	r.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()

	err := UpdateStatus(client, partialTemplate.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, partialTemplate)
}

func TestPostUpdateStatusHtmx(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "700700"}

	client := &mockUpdateStatusClient{}
	client.
		On("EditCase", mock.Anything, 123, sirius.CaseTypeLpa, sirius.Case{
			Status: shared.CaseStatusTypeWithdrawn,
		}).
		Return(nil)

	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)

	client.
		On("AvailableStatuses", mock.Anything, 123, sirius.CaseTypeLpa).
		Return([]string{"Cancelled", "Withdrawn"}, nil)

	partialTemplate := &mockTemplate{}
	partialTemplate.
		On("Func", mock.Anything, updateStatusData{
			Success:           true,
			Entity:            "lpa 700700",
			AvailableStatuses: []string{"Cancelled", "Withdrawn"},
			NewStatus:         "Withdrawn",
			CaseID:            123,
			CaseType:          "lpa",
			IsPartial:         true,
		}).
		Return(nil)

	form := url.Values{
		"status": {"Withdrawn"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&case=lpa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	r.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()

	err := UpdateStatus(client, partialTemplate.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, partialTemplate)
}
