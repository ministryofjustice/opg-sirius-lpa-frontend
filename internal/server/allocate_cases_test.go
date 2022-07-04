package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAllocateCasesClient struct {
	mock.Mock
}

func (m *mockAllocateCasesClient) AllocateCases(ctx sirius.Context, assigneeID int, allocations []sirius.CaseAllocation) error {
	args := m.Called(ctx, assigneeID, allocations)
	return args.Error(0)
}

func (m *mockAllocateCasesClient) Teams(ctx sirius.Context) ([]sirius.Team, error) {
	args := m.Called(ctx)
	return args.Get(0).([]sirius.Team), args.Error(1)
}

func (m *mockAllocateCasesClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func TestGetAllocateCases(t *testing.T) {
	client := &mockAllocateCasesClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, allocateCasesData{
			Teams:    []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			Entities: []string{"LPA 7000-0000-0000"},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AllocateCases(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetAllocateCasesMultiple(t *testing.T) {
	client := &mockAllocateCasesClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
	client.
		On("Case", mock.Anything, 456).
		Return(sirius.Case{UID: "7000-1111-1111", CaseType: "EPA"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(d allocateCasesData) bool {
			sort.Strings(d.Entities)

			return len(d.Teams) == 1 && d.Teams[0] == sirius.Team{ID: 1, DisplayName: "A Team"} &&
				len(d.Entities) == 2 &&
				d.Entities[0] == "EPA 7000-1111-1111" && d.Entities[1] == "LPA 7000-0000-0000"
		})).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&id=456", nil)
	w := httptest.NewRecorder()

	err := AllocateCases(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetAllocateCasesBadQueryString(t *testing.T) {
	testCases := map[string]string{
		"no-id":      "/",
		"bad-id":     "/?id=what",
		"one-bad-id": "/?id=1&id=bad&id=2",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := AllocateCases(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetAllocateCasesWhenTeamsErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockAllocateCasesClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, expectedError)
	client.
		On("Case", mock.Anything, mock.Anything).
		Return(sirius.Case{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AllocateCases(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetAllocateCasesWhenCaseErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockAllocateCasesClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, nil)
	client.
		On("Case", mock.Anything, mock.Anything).
		Return(sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AllocateCases(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetAllocateCasesWhenTemplateErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockAllocateCasesClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, nil)
	client.
		On("Case", mock.Anything, mock.Anything).
		Return(sirius.Case{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AllocateCases(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostAllocateCases(t *testing.T) {
	client := &mockAllocateCasesClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
	client.
		On("AllocateCases", mock.Anything, 66, []sirius.CaseAllocation{{ID: 123, CaseType: "LPA"}}).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, allocateCasesData{
			Success:          true,
			Teams:            []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			AssigneeUserName: "System user",
			Entities:         []string{"LPA 7000-0000-0000"},
		}).
		Return(nil)

	form := url.Values{
		"assignTo":     {"user"},
		"assigneeUser": {"66:System user"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AllocateCases(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostAllocateCasesMultiple(t *testing.T) {
	client := &mockAllocateCasesClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
	client.
		On("Case", mock.Anything, 456).
		Return(sirius.Case{UID: "7000-1111-1111", CaseType: "EPA"}, nil)
	client.
		On("AllocateCases", mock.Anything, 66, mock.MatchedBy(func(a []sirius.CaseAllocation) bool {
			lpa := sirius.CaseAllocation{ID: 123, CaseType: "LPA"}
			epa := sirius.CaseAllocation{ID: 456, CaseType: "EPA"}

			return len(a) == 2 && ((a[0] == lpa && a[1] == epa) || (a[1] == lpa && a[0] == epa))
		})).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(d allocateCasesData) bool {
			sort.Strings(d.Entities)

			return d.Success == true &&
				d.AssigneeUserName == "System user" &&
				len(d.Teams) == 1 && d.Teams[0] == sirius.Team{ID: 1, DisplayName: "A Team"} &&
				len(d.Entities) == 2 &&
				d.Entities[0] == "EPA 7000-1111-1111" && d.Entities[1] == "LPA 7000-0000-0000"
		})).
		Return(nil)

	form := url.Values{
		"assignTo":     {"user"},
		"assigneeUser": {"66:System user"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&id=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AllocateCases(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostAllocateCasesWhenAllocateCasesFails(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockAllocateCasesClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
	client.
		On("AllocateCases", mock.Anything, 66, []sirius.CaseAllocation{{ID: 123, CaseType: "LPA"}}).
		Return(expectedError)

	form := url.Values{
		"assignTo":     {"user"},
		"assigneeUser": {"66:System user"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AllocateCases(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostAllocateCasesWhenAssignToNotSet(t *testing.T) {
	client := &mockAllocateCasesClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
	client.
		On("AllocateCases", mock.Anything, mock.Anything, mock.Anything).
		Return(sirius.ValidationError{
			Field: sirius.FieldErrors{
				"assigneeId": {"empty": "Not set"},
			},
		})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, allocateCasesData{
			Teams:    []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			Entities: []string{"LPA 7000-0000-0000"},
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"assignTo": {"": "Assignee not set"},
				},
			},
		}).
		Return(nil)

	form := url.Values{
		"assigneeUser": {"66"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AllocateCases(client, template.Func)(w, r)
	assert.Nil(t, err)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostAllocateCasesWhenValidationError(t *testing.T) {
	testCases := map[string]struct {
		field            string
		value            string
		assigneeUserName string
	}{
		"team": {
			field: "assigneeTeam",
			value: "66",
		},
		"user": {
			field:            "assigneeUser",
			value:            "66:Some user",
			assigneeUserName: "Some user",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockAllocateCasesClient{}
			client.
				On("Teams", mock.Anything).
				Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
			client.
				On("Case", mock.Anything, 123).
				Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
			client.
				On("AllocateCases", mock.Anything, mock.Anything, mock.Anything).
				Return(sirius.ValidationError{Field: sirius.FieldErrors{
					"field":      {"reason": "Description"},
					"assigneeId": {"problem": "Because"},
				}})

			expectedErrors := sirius.FieldErrors{
				"field": {"reason": "Description"},
			}
			expectedErrors[tc.field] = map[string]string{"problem": "Because"}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, allocateCasesData{
					AssignTo:         name,
					Teams:            []sirius.Team{{ID: 1, DisplayName: "A Team"}},
					Entities:         []string{"LPA 7000-0000-0000"},
					Error:            sirius.ValidationError{Field: expectedErrors},
					AssigneeUserName: tc.assigneeUserName,
				}).
				Return(nil)

			form := url.Values{
				"assignTo": {name},
			}
			form.Add(tc.field, tc.value)

			r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := AllocateCases(client, template.Func)(w, r)
			assert.Nil(t, err)

			resp := w.Result()
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}
