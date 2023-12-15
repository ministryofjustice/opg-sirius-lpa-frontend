package sirius

import (
	"bytes"
	"context"
	"encoding/json"
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

// not tested by pact as this is just an amalgam of two other pact-tested methods
func TestCaseSummary(t *testing.T) {
	mockHttpClient := mockCaseSummaryHttpClient{}
	client := NewClient(&mockHttpClient, "http://localhost:8888")

	reqForDigitalLpaMatcher := mock.MatchedBy(func (r *http.Request) bool {
		digitalLpaUrl, _ := url.Parse("http://localhost:8888/lpa-api/v1/digital-lpas/M-QWER-TY34-3434")
		return digitalLpaUrl.String() == r.URL.String()
	})
	var digitalLpaBody bytes.Buffer
    json.NewEncoder(&digitalLpaBody).Encode(DigitalLpa{
    	ID: 1,
    })
	respForDigitalLpa := http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewReader(digitalLpaBody.Bytes())),
	}
	mockHttpClient.On("Do", reqForDigitalLpaMatcher).Return(&respForDigitalLpa, nil)

	reqForTasksForCaseMatcher := mock.MatchedBy(func (r *http.Request) bool {
		tasksForCaseUrl, _ := url.Parse("http://localhost:8888/lpa-api/v1/cases/1/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC")
		return tasksForCaseUrl.String() == r.URL.String()
	})
	var tasksForCaseBody bytes.Buffer
    json.NewEncoder(&tasksForCaseBody).Encode(map[string][]Task{
    	"tasks": []Task{
    		Task{},
    		Task{},
    		Task{},
    	},
    })
	respForTasksForCase := http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewReader(tasksForCaseBody.Bytes())),
	}
	mockHttpClient.On("Do", reqForTasksForCaseMatcher).Return(&respForTasksForCase, nil)

	_, err := client.CaseSummary(Context{Context: context.Background()}, "M-QWER-TY34-3434")
	assert.Nil(t, err)
}
