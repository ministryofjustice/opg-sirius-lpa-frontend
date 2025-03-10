package sirius

import (
	"fmt"
)

type NoteFile struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

type noteRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	File        *NoteFile `json:"file,omitempty"`
}

func (c *Client) CreateNote(ctx Context, entityID int, entityType EntityType, noteType, name, description string, file *NoteFile) error {
	data := noteRequest{
		Name:        name,
		Description: description,
		Type:        noteType,
		File:        file,
	}

	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/%ss/%d/notes", entityType, entityID), data, nil)
}
