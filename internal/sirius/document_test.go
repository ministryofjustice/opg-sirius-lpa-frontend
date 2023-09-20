package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestDocument(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/lpas/800/documents"),
						Query: dsl.MapMatcher{
							"type[]":    dsl.EachLike("Draft", 1),
							"type[-][]": dsl.EachLike("Preview", 1),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"id":                  dsl.Like(1),
							"uuid":                dsl.String("dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
							"type":                dsl.String("Draft"),
							"friendlyDescription": dsl.String("Dr Consuela Aysien - LPA perfect + reg due date: applicant"),
							"createdDate":         dsl.String(`15/12/2022 13:41:04`),
							"direction":           dsl.String("Outgoing"),
							"filename":            dsl.String("LP-A.pdf"),
							"mimeType":            dsl.String(`application\/pdf`),
							"correspondent": dsl.Like(map[string]interface{}{
								"id":        dsl.Like(189),
								"firstname": dsl.String("Consuela"),
								"surname":   dsl.String("Aysien"),
							}),
							"childCount": dsl.Like(0),
							"systemType": dsl.String("LP-A"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []Document{
				{
					ID:                  1,
					UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
					Type:                "Draft",
					FriendlyDescription: "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
					CreatedDate:         `15/12/2022 13:41:04`,
					Direction:           "Outgoing",
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/lpas/800/documents"),
						Query: dsl.MapMatcher{
							"type[]":    dsl.EachLike("Draft", 1),
							"type[-][]": dsl.EachLike("Preview", 1),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"id":                  dsl.Like(1),
							"uuid":                dsl.String("dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
							"type":                dsl.String("Draft"),
							"friendlyDescription": dsl.String("Dr Consuela Aysien - LPA perfect + reg due date: applicant"),
							"createdDate":         dsl.String(`15/12/2022 13:41:04`),
							"direction":           dsl.String("Outgoing"),
							"filename":            dsl.String("LP-A.pdf"),
							"mimeType":            dsl.String(`application\/pdf`),
							"correspondent": dsl.Like(map[string]interface{}{
								"id":        dsl.Like(189),
								"firstname": dsl.String("Consuela"),
								"surname":   dsl.String("Aysien"),
							}),
							"childCount": dsl.Like(0),
							"systemType": dsl.String("LP-A"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []Document{
				{
					ID:                  1,
					UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
					Type:                "Draft",
					FriendlyDescription: "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
					CreatedDate:         `15/12/2022 13:41:04`,
					Direction:           "Outgoing",
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				documents, err := client.Documents(Context{Context: context.Background()}, tc.caseType, 800, []string{TypeDraft}, []string{TypePreview})

				assert.Equal(t, tc.expectedResponse, documents)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestDocumentByUuid(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/documents/dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":                  dsl.Like(1),
							"uuid":                dsl.String("dfef6714-b4fe-44c2-b26e-90dfe3663e95"),
							"type":                dsl.String("Draft"),
							"friendlyDescription": dsl.String("Dr Consuela Aysien - LPA perfect + reg due date: applicant"),
							"createdDate":         dsl.String(`15/12/2022 13:41:04`),
							"direction":           dsl.String("Outgoing"),
							"filename":            dsl.String("LP-A.pdf"),
							"mimeType":            dsl.String(`application\/pdf`),
							"correspondent": dsl.Like(map[string]interface{}{
								"id": dsl.Like(189),
							}),
							"childCount": dsl.Like(0),
							"systemType": dsl.String("LP-A"),
							"content":    dsl.String("Test content"),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: Document{
				ID:                  1,
				UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				Type:                "Draft",
				FriendlyDescription: "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
				CreatedDate:         "15/12/2022 13:41:04",
				Direction:           "Outgoing",
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				document, err := client.DocumentByUUID(Context{Context: context.Background()}, "dfef6714-b4fe-44c2-b26e-90dfe3663e95")

				assert.Equal(t, tc.expectedResponse, document)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
