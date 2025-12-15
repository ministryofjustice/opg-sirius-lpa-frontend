package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"

	"github.com/stretchr/testify/assert"
)

func TestDocument(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		caseType         CaseType
		setup            func()
		expectedError    func(int) error
		expectedResponse []Document
	}{
		{
			name:     "LPA",
			caseType: CaseTypeLpa,
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa with a draft document").
					UponReceiving("A request for the documents by case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/lpas/800/documents"),
						Query: matchers.MapMatcher{
							"type[]":    matchers.Like("Draft"),
							"type[-][]": matchers.Like("Preview"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"id":                  matchers.Like(1),
							"uuid":                matchers.String("dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
							"type":                matchers.String("Draft"),
							"friendlyDescription": matchers.String("Dr Consuela Aysien - LPA perfect + reg due date: applicant"),
							"createdDate":         matchers.String(`15/12/2022 13:41:04`),
							"direction":           matchers.String("Outgoing"),
							"filename":            matchers.String("LP-A.pdf"),
							"mimeType":            matchers.String(`application\/pdf`),
							"correspondent": matchers.Like(map[string]interface{}{
								"id":        matchers.Like(189),
								"firstname": matchers.String("Consuela"),
								"surname":   matchers.String("Aysien"),
							}),
							"childCount": matchers.Like(0),
							"systemType": matchers.String("LP-A"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []Document{
				{
					ID:                  1,
					UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
					Type:                "Draft",
					FriendlyDescription: "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
					CreatedDate:         `15/12/2022 13:41:04`,
					Direction:           shared.DocumentDirectionOut,
					MimeType:            `application\/pdf`,
					SystemType:          "LP-A",
					FileName:            "LP-A.pdf",
					Correspondent:       Person{ID: 189, Firstname: "Consuela", Surname: "Aysien"},
					ChildCount:          0,
				},
			},
		},
		{
			name:     "Digital LPA",
			caseType: CaseTypeDigitalLpa,
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa with a draft document").
					UponReceiving("A request for the documents for a digital LPA").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/lpas/800/documents"),
						Query: matchers.MapMatcher{
							"type[]":    matchers.Like("Draft"),
							"type[-][]": matchers.Like("Preview"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"id":                  matchers.Like(1),
							"uuid":                matchers.String("dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
							"type":                matchers.String("Draft"),
							"friendlyDescription": matchers.String("Dr Consuela Aysien - LPA perfect + reg due date: applicant"),
							"createdDate":         matchers.String(`15/12/2022 13:41:04`),
							"direction":           matchers.String("Outgoing"),
							"filename":            matchers.String("LP-A.pdf"),
							"mimeType":            matchers.String(`application\/pdf`),
							"correspondent": matchers.Like(map[string]interface{}{
								"id":        matchers.Like(189),
								"firstname": matchers.String("Consuela"),
								"surname":   matchers.String("Aysien"),
							}),
							"childCount": matchers.Like(0),
							"systemType": matchers.String("LP-A"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []Document{
				{
					ID:                  1,
					UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
					Type:                "Draft",
					FriendlyDescription: "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
					CreatedDate:         `15/12/2022 13:41:04`,
					Direction:           shared.DocumentDirectionOut,
					MimeType:            `application\/pdf`,
					SystemType:          "LP-A",
					FileName:            "LP-A.pdf",
					Correspondent:       Person{ID: 189, Firstname: "Consuela", Surname: "Aysien"},
					ChildCount:          0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				documents, err := client.Documents(Context{Context: context.Background()}, tc.caseType, 800, []string{TypeDraft}, []string{TypePreview})

				assert.Equal(t, tc.expectedResponse, documents)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestDocumentByUuid(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedError    func(int) error
		expectedResponse Document
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa with a draft document").
					UponReceiving("A request for a document by uuid").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/documents/dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":                  matchers.Like(1),
							"uuid":                matchers.String("dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
							"type":                matchers.String("Draft"),
							"friendlyDescription": matchers.String("Dr Consuela Aysien - LPA perfect + reg due date: applicant"),
							"createdDate":         matchers.String(`15/12/2022 13:41:04`),
							"direction":           matchers.String("Outgoing"),
							"filename":            matchers.String("LP-A.pdf"),
							"mimeType":            matchers.String(`application\/pdf`),
							"correspondent": matchers.Like(map[string]interface{}{
								"id": matchers.Like(189),
							}),
							"childCount": matchers.Like(0),
							"systemType": matchers.String("LP-A"),
							"content":    matchers.String("Test content"),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Document{
				ID:                  1,
				UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:                "Draft",
				FriendlyDescription: "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
				CreatedDate:         "15/12/2022 13:41:04",
				Direction:           shared.DocumentDirectionOut,
				MimeType:            `application\/pdf`,
				SystemType:          "LP-A",
				FileName:            "LP-A.pdf",
				Content:             "Test content",
				Correspondent:       Person{ID: 189},
				ChildCount:          0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				document, err := client.DocumentByUUID(Context{Context: context.Background()}, "dfef6714-b4fe-44c2-b26e-90dfe3663e95")

				assert.Equal(t, tc.expectedResponse, document)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestGetPersonDocumentBy(t *testing.T) {
	t.Parallel()

	pactResponse := consumer.Response{
		Status: http.StatusOK,
		Body: matchers.Like(map[string]interface{}{
			"limit": matchers.Like(999),
			"total": matchers.Like(1),
			"pages": matchers.Like(map[string]interface{}{
				"current": matchers.Like(1),
				"total":   matchers.Like(1),
			}),
			"documents": matchers.EachLike(map[string]interface{}{
				"id":                  matchers.Like(4056),
				"uuid":                matchers.Like("dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
				"type":                matchers.Like("Save"),
				"friendlyDescription": matchers.Like("Dr Consuela Aysien - LPA perfect + reg due date: applicant"),
				"createdDate":         matchers.Like(`15/12/2022 13:41:04`),
				"direction":           matchers.Like("Outgoing"),
				"filename":            matchers.Like("LP-A.pdf"),
				"mimeType":            matchers.Like(`application\/pdf`),
				"childCount":          matchers.Like(0),
				"systemType":          matchers.Like("LP-A"),
			}, 1),
		}),
		Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
	}

	expectedResponse := DocumentList{
		Limit: 999,
		Pages: Pages{
			Current: 1,
			Total:   1,
		},
		Total: 1,
		Documents: []Document{
			{
				ID:                  4056,
				UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:                "Save",
				FriendlyDescription: "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
				CreatedDate:         `15/12/2022 13:41:04`,
				Direction:           shared.DocumentDirectionOut,
				MimeType:            `application\/pdf`,
				SystemType:          "LP-A",
				FileName:            "LP-A.pdf",
				ChildCount:          0,
			},
		},
	}

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		caseIDs          []string
		setup            func()
		expectedError    func(int) error
		expectedResponse DocumentList
	}{
		{
			name:    "Can get documents on a person",
			caseIDs: []string(nil),
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists with a published document").
					UponReceiving("A request for the documents on a person").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/400/documents"),
						Query: matchers.MapMatcher{
							"filter": matchers.String("draft:0,preview:0"),
							"limit":  matchers.Like("999"),
						},
					}).
					WithCompleteResponse(pactResponse)
			},
			expectedResponse: expectedResponse,
		},
		{
			name:    "Can get documents on a case",
			caseIDs: []string{"4056"},
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists with a published document").
					UponReceiving("A request for the documents on a case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/400/documents"),
						Query: matchers.MapMatcher{
							"filter": matchers.String("draft:0,preview:0,case:4056"),
							"limit":  matchers.Like("999"),
						},
					}).
					WithCompleteResponse(pactResponse)
			},
			expectedResponse: expectedResponse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				documents, err := client.GetPersonDocuments(Context{Context: context.Background()}, 400, tc.caseIDs)

				assert.Equal(t, tc.expectedResponse, documents)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

// non-pact test
func TestDocumentIsViewable(t *testing.T) {
	d := Document{}
	assert.True(t, d.IsViewable())

	d.Direction = shared.DocumentDirectionIn
	assert.True(t, d.IsViewable())

	d.SubType = "Reduced fee request evidence"
	assert.False(t, d.IsViewable())

	d.ReceivedDateTime = "11/12/2023"
	assert.True(t, d.IsViewable())
}
