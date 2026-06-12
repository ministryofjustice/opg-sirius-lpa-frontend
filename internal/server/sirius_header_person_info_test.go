package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSiriusHeaderPersonInfoClient struct {
	mock.Mock
}

func (m *mockSiriusHeaderPersonInfoClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func TestGetSiriusHeaderPersonInfo(t *testing.T) {
	donor := sirius.Person{ID: 123}
	caseItem := sirius.Case{Donor: &donor}
	client := &mockSiriusHeaderPersonInfoClient{}
	client.
		On("Case", mock.Anything, 1).
		Return(caseItem, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderPeopleInfoData{
			CaseID:         1,
			Case:           caseItem,
			SelectedPerson: donor,
			SelectedID:     123,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1", nil)
	w := httptest.NewRecorder()

	err := SiriusHeaderPeopleInfo(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSiriusHeaderPersonInfoWithSelected(t *testing.T) {
	attorney := sirius.Attorney{Person: sirius.Person{ID: 2}}
	caseItem := sirius.Case{Donor: &sirius.Person{ID: 123}, Attorneys: []sirius.Attorney{attorney}}
	client := &mockSiriusHeaderPersonInfoClient{}
	client.
		On("Case", mock.Anything, 1).
		Return(caseItem, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderPeopleInfoData{
			CaseID:         1,
			Case:           caseItem,
			SelectedPerson: attorney,
			SelectedID:     2,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&selected=2", nil)
	w := httptest.NewRecorder()

	err := SiriusHeaderPeopleInfo(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSiriusHeaderPersonInfoBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":           "/",
		"bad-id":          "/?id=test",
		"bad-selected-id": "/?id=1&selected=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			donor := sirius.Person{ID: 123}
			client := &mockSiriusHeaderPersonInfoClient{}
			client.
				On("Case", mock.Anything, 1).
				Return(sirius.Case{Donor: &donor}, nil)

			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := SiriusHeaderPeopleInfo(client, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetSiriusHeaderPersonInfoWhenCaseErrors(t *testing.T) {
	client := &mockSiriusHeaderPersonInfoClient{}
	client.
		On("Case", mock.Anything, 1).
		Return(sirius.Case{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1", nil)
	w := httptest.NewRecorder()

	err := SiriusHeaderPeopleInfo(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetSiriusHeaderPersonInfoGetSelectedPerson(t *testing.T) {
	caseItem := sirius.Case{
		Donor: &sirius.Person{ID: 1},
		Attorneys: []sirius.Attorney{
			{Person: sirius.Person{ID: 2}},
		},
		ReplacementAttorneys: []sirius.Attorney{
			{Person: sirius.Person{ID: 3}},
		},
		TrustCorporations: []sirius.Attorney{
			{Person: sirius.Person{ID: 4}},
		},
		CertificateProviders: []sirius.Person{{ID: 5}},
		NotifiedPersons:      []sirius.Person{{ID: 6}},
		Correspondent:        &sirius.Correspondent{Person: sirius.Person{ID: 7}},
	}

	testCases := []struct {
		Name           string
		SelectedID     string
		SelectedPerson sirius.Recipient
	}{
		{
			Name:           "attorney",
			SelectedID:     "2",
			SelectedPerson: caseItem.Attorneys[0],
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			client := &mockSiriusHeaderPersonInfoClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseItem, nil)

			selectedIdInt, _ := strconv.Atoi(tc.SelectedID)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, siriusHeaderPeopleInfoData{
					CaseID:         123,
					Case:           caseItem,
					SelectedPerson: tc.SelectedPerson,
					SelectedID:     selectedIdInt,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/?id=123&selected=%s", tc.SelectedID), nil)
			w := httptest.NewRecorder()

			err := SiriusHeaderPeopleInfo(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
