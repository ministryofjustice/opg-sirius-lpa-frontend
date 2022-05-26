package server

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEventClient struct {
	mock.Mock
}

func (m *mockEventClient) CreateNote(ctx sirius.Context, entityID int, entityType sirius.EntityType, noteType, name, description string, file *sirius.NoteFile) error {
	args := m.Called(ctx, entityID, entityType, noteType, name, description, file)
	return args.Error(0)
}

func (m *mockEventClient) NoteTypes(ctx sirius.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockEventClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockEventClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func TestGetEvent(t *testing.T) {
	testCases := map[string]struct {
		url            string
		clientSetup    func(*mockEventClient)
		expectedEntity string
	}{
		"person": {
			url: "/?id=123&entity=person",
			clientSetup: func(client *mockEventClient) {
				client.
					On("Person", mock.Anything, 123).
					Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
			},
			expectedEntity: "John Doe",
		},
		"lpa": {
			url: "/?id=123&entity=lpa",
			clientSetup: func(client *mockEventClient) {
				client.
					On("Case", mock.Anything, 123).
					Return(sirius.Case{UID: "7000-0000-0001", CaseType: "LPA"}, nil)
			},
			expectedEntity: "LPA 7000-0000-0001",
		},
		"epa": {
			url: "/?id=123&entity=epa",
			clientSetup: func(client *mockEventClient) {
				client.
					On("Case", mock.Anything, 123).
					Return(sirius.Case{UID: "7000-0000-0001", CaseType: "EPA"}, nil)
			},
			expectedEntity: "EPA 7000-0000-0001",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockEventClient{}
			client.
				On("NoteTypes", mock.Anything).
				Return([]string{"a", "b"}, nil)
			tc.clientSetup(client)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, eventData{
					NoteTypes: []string{"a", "b"},
					Entity:    tc.expectedEntity,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()

			err := Event(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestGetEventBadQueryString(t *testing.T) {
	testCases := map[string]string{
		"no-id":      "/?entity=person",
		"no-entity":  "/?id=123",
		"bad-entity": "/?id=123&entity=what",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockEventClient{}
			client.
				On("NoteTypes", mock.Anything).
				Return([]string{"a", "b"}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, eventData{
					NoteTypes: []string{"a", "b"},
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := Event(client, template.Func)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetEventWhenNoteTypeErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockEventClient{}
	client.
		On("NoteTypes", mock.Anything).
		Return([]string{}, expectedError)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&entity=person", nil)
	w := httptest.NewRecorder()

	err := Event(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetEventWhenTemplateErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockEventClient{}
	client.
		On("NoteTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&entity=person", nil)
	w := httptest.NewRecorder()

	err := Event(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
}

func TestPostEvent(t *testing.T) {
	client := &mockEventClient{}
	client.
		On("NoteTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("CreateNote", mock.Anything, 123, sirius.EntityTypePerson, "Application processing", "Something", "More words", (*sirius.NoteFile)(nil)).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, eventData{
			Success:   true,
			NoteTypes: []string{"a", "b"},
			Entity:    "John Doe",
		}).
		Return(nil)

	var buf bytes.Buffer
	form := multipart.NewWriter(&buf)
	_ = form.WriteField("type", "Application processing")
	_ = form.WriteField("name", "Something")
	_ = form.WriteField("description", "More words")
	_ = form.Close()

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&entity=person", &buf)
	r.Header.Add("Content-Type", form.FormDataContentType())
	w := httptest.NewRecorder()

	err := Event(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostEventWithFile(t *testing.T) {
	client := &mockEventClient{}
	client.
		On("NoteTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("CreateNote", mock.Anything, 123, sirius.EntityTypePerson, "Application processing", "Something", "More words",
			&sirius.NoteFile{Name: "test.txt", Type: "application/octet-stream", Source: "SGV5"}).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, eventData{
			Success:   true,
			NoteTypes: []string{"a", "b"},
			Entity:    "John Doe",
		}).
		Return(nil)

	var buf bytes.Buffer
	form := multipart.NewWriter(&buf)
	_ = form.WriteField("type", "Application processing")
	_ = form.WriteField("name", "Something")
	_ = form.WriteField("description", "More words")
	part, _ := form.CreateFormFile("file", "test.txt")
	_, _ = part.Write([]byte("Hey"))
	_ = form.Close()

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&entity=person", &buf)
	r.Header.Add("Content-Type", form.FormDataContentType())
	w := httptest.NewRecorder()

	err := Event(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostEventWithBadForm(t *testing.T) {
	client := &mockEventClient{}
	client.
		On("NoteTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

	var buf bytes.Buffer
	form := multipart.NewWriter(&buf)
	_ = form.WriteField("type", "Application processing")
	_ = form.WriteField("name", "Something")
	_ = form.WriteField("description", "More words")
	_ = form.Close()

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&entity=person", &buf)
	w := httptest.NewRecorder()

	err := Event(client, nil)(w, r)

	assert.NotNil(t, err)
}

func TestPostEventWhenCreateNoteFails(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockEventClient{}
	client.
		On("NoteTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("CreateNote", mock.Anything, 123, sirius.EntityTypePerson, "Application processing", "Something", "More words", (*sirius.NoteFile)(nil)).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, eventData{
			Success:   true,
			NoteTypes: []string{"a", "b"},
			Entity:    "John Doe",
		}).
		Return(nil)

	var buf bytes.Buffer
	form := multipart.NewWriter(&buf)
	_ = form.WriteField("type", "Application processing")
	_ = form.WriteField("name", "Something")
	_ = form.WriteField("description", "More words")
	_ = form.Close()

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&entity=person", &buf)
	r.Header.Add("Content-Type", form.FormDataContentType())
	w := httptest.NewRecorder()

	err := Event(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
}

func TestPostEventWhenValidationError(t *testing.T) {
	expectedErrors := sirius.ValidationError{
		Field: sirius.FieldErrors{
			"field": {"reason": "Description"},
		},
	}

	client := &mockEventClient{}
	client.
		On("NoteTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("CreateNote", mock.Anything, 123, sirius.EntityTypePerson, "Application processing", "Something", "More words", (*sirius.NoteFile)(nil)).
		Return(expectedErrors)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, eventData{
			Success:     false,
			NoteTypes:   []string{"a", "b"},
			Error:       expectedErrors,
			Entity:      "John Doe",
			Type:        "Application processing",
			Name:        "Something",
			Description: "More words",
		}).
		Return(nil)

	var buf bytes.Buffer
	form := multipart.NewWriter(&buf)
	_ = form.WriteField("type", "Application processing")
	_ = form.WriteField("name", "Something")
	_ = form.WriteField("description", "More words")
	_ = form.Close()

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&entity=person", &buf)
	r.Header.Add("Content-Type", form.FormDataContentType())
	w := httptest.NewRecorder()

	err := Event(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
