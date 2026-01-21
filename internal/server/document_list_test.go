package server

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDocumentListClient struct {
	mock.Mock
}

func (m *mockDocumentListClient) CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.Case), args.Error(1)
}

func (m *mockDocumentListClient) GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error) {
	args := m.Called(ctx, personID, caseIDs)
	return args.Get(0).(sirius.DocumentList), args.Error(1)
}

func (m *mockDocumentListClient) DownloadMultiple(ctx sirius.Context, docIDs []string) (*http.Response, error) {
	args := m.Called(ctx, docIDs)
	return args.Get(0).(*http.Response), args.Error(1)
}

var singleDocumentList = sirius.DocumentList{
	Limit: 999,
	Pages: sirius.Pages{
		Current: 1,
		Total:   1,
	},
	Total: 1,
	Documents: []sirius.Document{
		{
			ID:                  639,
			UUID:                "31e6f4c2-5f8b-47c3-bc98-64b47c938e52",
			Type:                "Save",
			FriendlyDescription: "Letter",
			CreatedDate:         "25/07/2022 14:17:13",
			Direction:           shared.DocumentDirectionOut,
			FileName:            "LP-NA-3A.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-1234-0000",
					SubType:  "pfa",
					CaseType: "LPA",
				},
			},
			SystemType: "LP-NA-3A",
		},
	},
}

var twoCasesDocumentList = sirius.DocumentList{
	Limit: 999,
	Pages: sirius.Pages{
		Current: 1,
		Total:   1,
	},
	Total: 2,
	Documents: []sirius.Document{
		{
			ID:                  443,
			UUID:                "c8e3a1df-7b9b-4d45-94d9-2b8fc0d9e0fd",
			Type:                "LPA",
			FriendlyDescription: "LP1H - Health Instrument",
			CreatedDate:         "01/06/2022 15:39:01",
			Direction:           shared.DocumentDirectionIn,
			FileName:            "LP1H.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-9876-0000",
					SubType:  "hw",
					CaseType: "LPA",
				},
			},
			SubType: "hw",
		},
		{
			ID:                  639,
			UUID:                "31e6f4c2-5f8b-47c3-bc98-64b47c938e52",
			Type:                "Save",
			FriendlyDescription: "Letter",
			CreatedDate:         "25/07/2022 14:17:13",
			Direction:           shared.DocumentDirectionEmpty,
			FileName:            "LP-NA-3A.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-1234-0000",
					SubType:  "pfa",
					CaseType: "LPA",
				},
			},
			SystemType: "LP-NA-3A",
		},
	},
}

var allDocumentList = sirius.DocumentList{
	Limit: 999,
	Pages: sirius.Pages{
		Current: 1,
		Total:   1,
	},
	Total: 3,
	Documents: []sirius.Document{
		{
			ID:                  443,
			UUID:                "c8e3a1df-7b9b-4d45-94d9-2b8fc0d9e0fd",
			Type:                "LPA",
			FriendlyDescription: "LP1H - Health Instrument",
			CreatedDate:         "01/06/2022 15:39:01",
			Direction:           shared.DocumentDirectionIn,
			FileName:            "LP1H.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-9876-0000",
					SubType:  "hw",
					CaseType: "LPA",
				},
			},
			SubType: "hw",
		},
		{
			ID:                  639,
			UUID:                "31e6f4c2-5f8b-47c3-bc98-64b47c938e52",
			Type:                "Save",
			FriendlyDescription: "Letter",
			CreatedDate:         "25/07/2022 14:17:13",
			Direction:           shared.DocumentDirectionOut,
			FileName:            "LP-NA-3A.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-1234-0000",
					SubType:  "pfa",
					CaseType: "LPA",
				},
			},
			SystemType: "LP-NA-3A",
		},
		{
			ID:                  928,
			UUID:                "d9e12f73-3ab2-4d24-9a63-6b0b3e49b1c5",
			Type:                "Application Related",
			FriendlyDescription: "EPA.pdf",
			CreatedDate:         "08/01/2025 10:36:41",
			Direction:           shared.DocumentDirectionIn,
			FileName:            "EPA.pdf",
			MimeType:            "application/pdf",
			CaseItems: []sirius.Case{
				{
					UID:      "7000-5678-0000",
					SubType:  "pfa",
					CaseType: "EPA",
				},
			},
			SubType: "pfa",
		},
	},
}

func TestGetDocumentList(t *testing.T) {
	cases := []sirius.Case{
		{
			ID:       1,
			CaseType: "LPA",
			SubType:  "PFA",
			UID:      "7000-1234-0000",
		},
		{
			ID:       2,
			CaseType: "LPA",
			SubType:  "HW",
			UID:      "7000-9876-0000",
		},
		{
			ID:       3,
			CaseType: "EPA",
			SubType:  "PFA",
			UID:      "7000-5678-0000",
		},
	}

	tests := []struct {
		name             string
		cases            []sirius.Case
		documentList     sirius.DocumentList
		expectedMultiple bool
		expectedCases    []sirius.Case
		caseIDs          []string
		path             string
	}{
		{
			name:             "on person with multiple cases",
			cases:            cases,
			documentList:     allDocumentList,
			expectedMultiple: true,
			expectedCases:    cases,
			caseIDs:          []string(nil),
			path:             "/donor/82/documents",
		},
		{
			name:             "on person with one case",
			cases:            cases[:1],
			documentList:     singleDocumentList,
			expectedMultiple: false,
			expectedCases:    cases[:1],
			caseIDs:          []string(nil),
			path:             "/donor/82/documents",
		},
		{
			name:             "one case specified",
			cases:            cases,
			documentList:     singleDocumentList,
			expectedMultiple: false,
			expectedCases:    []sirius.Case{cases[0]},
			caseIDs:          []string{"1"},
			path:             "/donor/82/documents?uid[]=7000-1234-0000",
		},
		{
			name:             "multiple cases specified",
			cases:            cases,
			documentList:     twoCasesDocumentList,
			expectedMultiple: true,
			expectedCases:    []sirius.Case{cases[0], cases[1]},
			caseIDs:          []string{"1", "2"},
			path:             "/donor/82/documents?uid[]=7000-1234-0000&uid[]=7000-9876-0000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockDocumentListClient{}
			client.
				On("CasesByDonor", mock.Anything, 82).
				Return(tc.cases, nil)
			client.
				On("GetPersonDocuments", mock.Anything, 82, tc.caseIDs).
				Return(tc.documentList, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything,
					documentListData{
						SelectedCases:         tc.expectedCases,
						DocumentList:          tc.documentList,
						MultipleCasesSelected: tc.expectedMultiple,
					},
				).
				Return(nil)

			server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

			r, _ := http.NewRequest(http.MethodGet, tc.path, nil)
			_, err := server.serve(r)

			assert.Nil(t, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestDocumentListDownloadMultipleSuccess(t *testing.T) {
	cases := []sirius.Case{{ID: 1, UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)

	downloadResp := &http.Response{
		StatusCode: http.StatusCreated,
		Header: http.Header{
			"Content-Type": []string{"content/octet-stream"},
		},
		Body: io.NopCloser(strings.NewReader("document-download-bytes")),
	}

	client.
		On("DownloadMultiple", mock.Anything, []string{"doc-uuid"}).
		Return(downloadResp, nil)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	form := url.Values{}
	form.Add("document", "doc-uuid")
	form.Add("actionDownload", "true")
	req, _ := http.NewRequest(http.MethodPost, "/donor/82/documents", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formUrlEncoded)

	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, "content/octet-stream", resp.Header().Get("Content-Type"))
	assert.Equal(t, "document-download-bytes", resp.Body.String())

	mock.AssertExpectationsForObjects(t, client)
	client.AssertNotCalled(t, "GetPersonDocuments")
}

func TestDocumentListDownloadMultipleError(t *testing.T) {
	cases := []sirius.Case{{ID: 1, UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("DownloadMultiple", mock.Anything, []string{"doc-uuid"}).
		Return((*http.Response)(nil), errExample)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	form := url.Values{}
	form.Add("document", "doc-uuid")
	form.Add("actionDownload", "true")
	req, _ := http.NewRequest(http.MethodPost, "/donor/82/documents", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formUrlEncoded)

	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
	client.AssertNotCalled(t, "GetPersonDocuments")
}

func TestDocumentListShowsValidationErrorWhenNoDocumentsSelected(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-1234-0000"},
		{ID: 2, UID: "7000-9876-0000"},
	}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
		Return(allDocumentList, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			documentListData{
				SelectedCases:         cases,
				DocumentList:          allDocumentList,
				MultipleCasesSelected: true,
				Error: sirius.ValidationError{
					Detail: "Select one or more documents and try again.",
				},
			},
		).
		Return(nil)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

	form := url.Values{}
	form.Add("actionDownload", "true")
	req, _ := http.NewRequest(http.MethodPost, "/donor/82/documents", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formUrlEncoded)

	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	mock.AssertExpectationsForObjects(t, client, template)
	client.AssertNotCalled(t, "DownloadMultiple")
}

func TestDocumentListDismissValidation(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-1234-0000"},
		{ID: 2, UID: "7000-9876-0000"},
		{ID: 3, UID: "7000-5678-0000"},
	}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string{"1", "2"}).
		Return(twoCasesDocumentList, nil)

	expectedCases := []sirius.Case{cases[0], cases[1]}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			documentListData{
				SelectedCases:         expectedCases,
				DocumentList:          twoCasesDocumentList,
				MultipleCasesSelected: true,
			},
		).
		Return(nil)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

	form := url.Values{}
	form.Add("uid[]", "7000-1234-0000")
	form.Add("uid[]", "7000-9876-0000")
	form.Add("dismissValidation", "true")
	req, _ := http.NewRequest(http.MethodPost, "/donor/82/documents", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", formUrlEncoded)

	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	mock.AssertExpectationsForObjects(t, client, template)
	client.AssertNotCalled(t, "DownloadMultiple")
}

func TestDocumentListInvalidDonorID(t *testing.T) {
	client := &mockDocumentListClient{}
	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/donor/abc/documents", nil)
	_, err := server.serve(req)

	assert.Error(t, err)
}

func TestGetDocumentListWhenCasesByDonorErrors(t *testing.T) {
	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return([]sirius.Case{}, errExample)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDocumentListWhenGetPersonDocumentsErrors(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", SubType: "PFA", UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
		Return(sirius.DocumentList{}, errExample)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDocumentListWhenTemplateErrors(t *testing.T) {
	cases := []sirius.Case{{ID: 1, CaseType: "LPA", SubType: "PFA", UID: "7000-1234-0000"}}

	client := &mockDocumentListClient{}
	client.
		On("CasesByDonor", mock.Anything, 82).
		Return(cases, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
		Return(singleDocumentList, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			documentListData{
				SelectedCases:         cases,
				DocumentList:          singleDocumentList,
				MultipleCasesSelected: false,
			},
		).
		Return(errExample)

	server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestSuccessMessageFormatter(t *testing.T) {
	tests := []struct {
		name            string
		docFriendlyName string
		docCreatedTime  string
		layout          string
		format          string
		expectedResult  string
	}{
		{
			name:            "valid date format",
			docFriendlyName: "Letter",
			docCreatedTime:  "02/07/2025 14:17:13",
			layout:          "02/01/2006 15:04:05",
			format:          "02/01/2006",
			expectedResult:  "02/07/2025 Letter",
		},
		{
			name:            "different valid date",
			docFriendlyName: "LP1H - Health Instrument",
			docCreatedTime:  "01/06/2022 15:39:01",
			layout:          "02/01/2006 15:04:05",
			format:          "02/01/2006",
			expectedResult:  "01/06/2022 LP1H - Health Instrument",
		},
		{
			name:            "invalid date format",
			docFriendlyName: "Document",
			docCreatedTime:  "invalid-date",
			layout:          "02/01/2006 15:04:05",
			format:          "02/01/2006",
			expectedResult:  "invalid date",
		},
		{
			name:            "empty friendly name",
			docFriendlyName: "",
			docCreatedTime:  "08/01/2025 10:36:41",
			layout:          "02/01/2006 15:04:05",
			format:          "02/01/2006",
			expectedResult:  "08/01/2025 ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := successMessageFormatter(tc.docFriendlyName, tc.docCreatedTime, tc.layout, tc.format)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestDocumentListSuccessMessage(t *testing.T) {
	cases := []sirius.Case{{ID: 1, UID: "7000-1234-0000"}}

	tests := []struct {
		name            string
		queryParams     string
		formValue       string
		expectedSuccess bool
		expectedMessage string
	}{
		{
			name:            "success query with no dismiss notification",
			queryParams:     "?success=true&documentFriendlyName=Letter&documentCreatedTime=02/07/2025%2014:17:13",
			formValue:       "",
			expectedSuccess: true,
			expectedMessage: "02/07/2025 Letter",
		},
		{
			name:            "success query with dismiss notification",
			queryParams:     "?success=true&documentFriendlyName=Letter&documentCreatedTime=02/07/2025%2014:17:13",
			formValue:       "true",
			expectedSuccess: false,
			expectedMessage: "",
		},
		{
			name:            "no success query",
			queryParams:     "?documentFriendlyName=Letter",
			formValue:       "",
			expectedSuccess: false,
			expectedMessage: "",
		},
		{
			name:            "success with invalid date",
			queryParams:     "?success=true&documentFriendlyName=Letter&documentCreatedTime=invalid",
			formValue:       "",
			expectedSuccess: true,
			expectedMessage: "invalid date",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockDocumentListClient{}
			client.
				On("CasesByDonor", mock.Anything, 82).
				Return(cases, nil)
			client.
				On("GetPersonDocuments", mock.Anything, 82, []string(nil)).
				Return(singleDocumentList, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, mock.MatchedBy(func(data documentListData) bool {
					return data.Success == tc.expectedSuccess && data.SuccessMessage == tc.expectedMessage
				})).
				Return(nil)

			server := newMockServer("/donor/{id}/documents", DocumentList(client, template.Func))

			form := url.Values{}
			if tc.formValue != "" {
				form.Add("dismissNotification", tc.formValue)
			}
			req, _ := http.NewRequest(http.MethodGet, "/donor/82/documents"+tc.queryParams, nil)
			if len(form) > 0 {
				req.Form = form
			}

			resp, err := server.serve(req)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
