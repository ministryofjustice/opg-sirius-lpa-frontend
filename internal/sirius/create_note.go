package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *Client) CreateNote(ctx Context, entityID int, entityType, noteType, name, description string, file *NoteFile) error {
	data, err := json.Marshal(noteRequest{
		Name:        name,
		Description: description,
		Type:        noteType,
		File:        file,
	})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/%ss/%d/notes", entityType, entityID), bytes.NewReader(data))
	if err != nil {
		return err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if res.StatusCode != http.StatusCreated {
		return newStatusError(res)
	}

	return nil
}
