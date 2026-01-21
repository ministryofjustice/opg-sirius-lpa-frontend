package sirius

import (
	"fmt"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type LpaEvents []LpaEvent

type LpaEventsResponse struct {
	Events   LpaEvents `json:"events"`
	Limit    int       `json:"limit"`
	Total    int       `json:"total"`
	Pages    any       `json:"pages"`
	Metadata any       `json:"metadata"`
}

type LpaEvent struct {
	Changes    any                       `json:"changeSet,omitempty"`
	CreatedOn  string                    `json:"createdOn"`
	Entity     any                       `json:"entity,omitempty"`
	Assignee   LpaEventUser              `json:"assignee,omitempty"`
	User       LpaEventUser              `json:"user,omitempty"`
	Hash       string                    `json:"hash"`
	OwningCase OwningCase                `json:"owningCase,omitempty"`
	ID         int                       `json:"id,omitempty"`
	UID        string                    `json:"uid,omitempty"`
	SourceType shared.LpaEventSourceType `json:"sourceType"`
	Type       string                    `json:"type,omitempty"`
}

type OwningCase struct {
	ID          int    `json:"id,omitempty"`
	UID         string `json:"uId,omitempty"`
	CaseSubtype string `json:"caseSubtype,omitempty"`
	CaseType    string `json:"caseType,omitempty"`
}

type LpaEventUser struct {
	DisplayName string    `json:"displayName,omitempty"`
	Email       string    `json:"email,omitempty"`
	Deleted     bool      `json:"deleted,omitempty"`
	PhoneNumber string    `json:"phoneNumber,omitempty"`
	Teams       []LpaTeam `json:"teams,omitempty"`
}

type LpaTeam struct {
	DisplayName string `json:"displayName,omitempty"`
}

// will refactor for getting events across multiple cases
func (c *Client) GetEvents(ctx Context, donorId int) (LpaEventsResponse, error) {
	var resp LpaEventsResponse
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/persons/%d/events?sort=id:desc", donorId), &resp)
	return resp, err
}
