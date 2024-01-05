package sirius

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCaseSummaryHttpClient struct {
	mock.Mock
}

func (m *mockCaseSummaryHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

type testCase struct {
	Description          string
	DigitalLpaResponse   http.Response
	DigitalLpaError      error
	TasksForCaseResponse http.Response
	TasksForCaseError    error
	WarningsResponse     http.Response
	WarningsError        error
	ExpectedError        error
}

func setupTestCase(
	t *testing.T,
	description string,
	digitalLpaError error,
	tasksForCaseError error,
	warningsError error,
	expectedError error,
) testCase {
	// digital LPA mock
	var digitalLpaBody bytes.Buffer
	err := json.NewEncoder(&digitalLpaBody).Encode(DigitalLpa{
		SiriusData: SiriusData{
			ID: 1,
		},
	})
	if err != nil {
		t.Fatal("Could not compile digital LPA JSON")
	}
	respForDigitalLpa := http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(digitalLpaBody.Bytes())),
	}

	// tasks for case mock
	var tasksForCaseBody bytes.Buffer
	err = json.NewEncoder(&tasksForCaseBody).Encode(map[string][]Task{
		"tasks": []Task{},
	})
	if err != nil {
		t.Fatal("Could not compile tasks for case JSON")
	}
	respForTasksForCase := http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(tasksForCaseBody.Bytes())),
	}

	// warnings mock
	var warningsBody bytes.Buffer
	err = json.NewEncoder(&warningsBody).Encode([]Warning{})
	if err != nil {
		t.Fatal("Could not compile warnings JSON")
	}
	respForWarnings := http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(warningsBody.Bytes())),
	}

	return testCase{
		Description:          description,
		DigitalLpaResponse:   respForDigitalLpa,
		DigitalLpaError:      digitalLpaError,
		TasksForCaseResponse: respForTasksForCase,
		TasksForCaseError:    tasksForCaseError,
		WarningsResponse:     respForWarnings,
		WarningsError:        warningsError,
		ExpectedError:        expectedError,
	}
}

func setupTestCases(t *testing.T) []testCase {
	digitalLpaErr := errors.New("Unable to fetch digital LPA")
	tasksForCaseErr := errors.New("Unable to fetch tasks for case")
	warningsErr := errors.New("Unable to fetch warnings")

	return []testCase{
		setupTestCase(t, "Case summary: all requests successful", nil, nil, nil, nil),
		setupTestCase(t, "Case summary: digital LPA request failure", digitalLpaErr, nil, nil, digitalLpaErr),
		setupTestCase(t, "Case summary: tasks for case request failure", nil, tasksForCaseErr, nil, tasksForCaseErr),
		setupTestCase(t, "Case summary: warnings request failure", nil, nil, warningsErr, warningsErr),
	}
}

// not tested by pact as this is just an amalgam of other pact-tested methods
func TestCaseSummary(t *testing.T) {
	reqForDigitalLpaMatcher := mock.MatchedBy(func(r *http.Request) bool {
		digitalLpaUrl, _ := url.Parse("http://localhost:8888/lpa-api/v1/digital-lpas/M-QWER-TY34-3434")
		return digitalLpaUrl.String() == r.URL.String()
	})

	reqForTasksForCaseMatcher := mock.MatchedBy(func(r *http.Request) bool {
		tasksForCaseUrl, _ := url.Parse("http://localhost:8888/lpa-api/v1/cases/1/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC")
		return tasksForCaseUrl.String() == r.URL.String()
	})

	reqForWarningsMatcher := mock.MatchedBy(func(r *http.Request) bool {
		warningsUrl, _ := url.Parse("http://localhost:8888/lpa-api/v1/cases/1/warnings")
		return warningsUrl.String() == r.URL.String()
	})

	for _, testCase := range setupTestCases(t) {
		mockHttpClient := mockCaseSummaryHttpClient{}
		client := NewClient(&mockHttpClient, "http://localhost:8888")

		mockHttpClient.On("Do", reqForDigitalLpaMatcher).Return(&testCase.DigitalLpaResponse, testCase.DigitalLpaError)
		mockHttpClient.On("Do", reqForTasksForCaseMatcher).Return(&testCase.TasksForCaseResponse, testCase.TasksForCaseError)
		mockHttpClient.On("Do", reqForWarningsMatcher).Return(&testCase.WarningsResponse, testCase.WarningsError)

		_, err := client.CaseSummary(Context{Context: context.Background()}, "M-QWER-TY34-3434")
		assert.Equal(t, testCase.ExpectedError, err)
	}
}
