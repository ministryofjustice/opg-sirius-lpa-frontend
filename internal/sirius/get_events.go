package sirius

import (
	"fmt"
)

type APIEvent []Event

type Event struct {
	ChangeSet     []interface{} `json:"changeSet"`
	CreatedOn     string        `json:"createdOn"`
	Entity        any           `json:"entity"`
	Hash          string        `json:"hash"`
	ID            int           `json:"id"`
	Source        string        `json:"source"`
	SourceType    string        `json:"sourceType"`
	Type          string        `json:"type"`
	User          EventUser     `json:"user"`
	UUID          string        `json:"uuid"`
	FormattedUUID string        `json:"showUuid,omitempty"`
	Applied       string        `json:"applied,omitempty"`
	DateTime      string        `json:"dateTime,omitempty"`
}

type EventUser struct {
	DisplayName string `json:"displayName"`
}

// GetCombinedEvents Gets combined events from both Sirius and LPA Store for digital LPAs
func (c *Client) GetCombinedEvents(ctx Context, uid string) (APIEvent, error) {
	var events APIEvent
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/events", uid), &events)
	return events, err
}
