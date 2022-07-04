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
	"golang.org/x/sync/errgroup"
)

const Megabyte = 1024 * 1024

type EventClient interface {
	NoteTypes(ctx sirius.Context) ([]string, error)
	CreateNote(ctx sirius.Context, entityID int, entityType sirius.EntityType, noteType, name, description string, file *sirius.NoteFile) error
	Person(ctx sirius.Context, id int) (sirius.Person, error)
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type eventData struct {
	XSRFToken string
	NoteTypes []string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

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

		entityType, err := sirius.ParseEntityType(r.FormValue("entity"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := eventData{XSRFToken: ctx.XSRFToken}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			noteTypes, err := client.NoteTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			data.NoteTypes = noteTypes
			return nil
		})

		group.Go(func() error {
			switch entityType {
			case sirius.EntityTypePerson:
				person, err := client.Person(ctx.With(groupCtx), entityID)
				if err != nil {
					return err
				}
				data.Entity = fmt.Sprintf("%s %s", person.Firstname, person.Surname)
			case sirius.EntityTypeLpa, sirius.EntityTypeEpa:
				caseitem, err := client.Case(ctx.With(groupCtx), entityID)
				if err != nil {
					return err
				}
				data.Entity = fmt.Sprintf("%s %s", caseitem.CaseType, caseitem.UID)
			}

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			if err := r.ParseMultipartForm(64 * Megabyte); err != nil {
				return err
			}

			var (
				noteType    = postFormString(r, "type")
				name        = postFormString(r, "name")
				description = postFormString(r, "description")
				file, err   = findNoteFile(r.MultipartForm, "file")
			)
			if err != nil {
				return err
			}

			err = client.CreateNote(ctx, entityID, entityType, noteType, name, description, file)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
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
