package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

const Megabyte = 1024 * 1024

type EventClient interface {
	NoteTypes(ctx sirius.Context) ([]string, error)
	CreateNote(ctx sirius.Context, entityID int, entityType, noteType, name, description string, file *sirius.NoteFile) error
}

type eventData struct {
	XSRFToken string
	NoteTypes []string
	Success   bool
	Errors    sirius.ValidationErrors

	Type        string
	Name        string
	Description string
}

func Event(client EventClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		entityID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		entityType := r.FormValue("entity")
		if entityType != "person" && entityType != "lpa" && entityType != "epa" {
			return fmt.Errorf("entity must be one of 'person', 'lpa' or 'epa'")
		}

		ctx := getContext(r)

		noteTypes, err := client.NoteTypes(ctx)
		if err != nil {
			return err
		}

		data := eventData{
			Success:   false,
			XSRFToken: ctx.XSRFToken,
			NoteTypes: noteTypes,
		}

		if r.Method == http.MethodPost {
			// TODO figure out what this limit should be
			if err := r.ParseMultipartForm(10 * Megabyte); err != nil {
				return err
			}

			var (
				noteType    = r.FormValue("type")
				name        = r.FormValue("name")
				description = r.FormValue("description")
				file, err   = findNoteFile(r.MultipartForm, "file")
			)
			if err != nil {
				return err
			}

			err = client.CreateNote(ctx, entityID, entityType, noteType, name, description, file)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Errors = ve.Errors
				data.Type = noteType
				data.Name = name
				data.Description = description
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}

func findNoteFile(form *multipart.Form, key string) (*sirius.NoteFile, error) {
	files := form.File["file"]
	if len(files) != 1 {
		return nil, nil
	}

	f, err := files[0].Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)

	if _, err := io.Copy(enc, f); err != nil {
		return nil, err
	}

	if err := enc.Close(); err != nil {
		return nil, err
	}

	return &sirius.NoteFile{
		Name:   files[0].Filename,
		Type:   files[0].Header.Get("Content-Type"),
		Source: buf.String(),
	}, nil
}
