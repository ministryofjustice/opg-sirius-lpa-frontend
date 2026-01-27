package sirius

import (
	"fmt"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type LpaEventsResponse struct {
	Events   []LpaEvent `json:"events"`
	Limit    int        `json:"limit"`
	Total    int        `json:"total"`
	Pages    Pages      `json:"pages"`
	Metadata any        `json:"metadata"`
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
	UID        string                    `json:"uuid,omitempty"`
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
func (c *Client) GetEvents(ctx Context, donorId string, caseIds []string) (LpaEventsResponse, error) {
	var resp LpaEventsResponse

	selectedCaseIds := ""
	for i, caseId := range caseIds {
		if i == 0 {
			selectedCaseIds = selectedCaseIds + "filter=case:" + caseId
		} else {
			selectedCaseIds = selectedCaseIds + ",case:" + caseId
		}
	}

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/persons/%s/events?%s&sort=id:desc&limit=999", donorId, selectedCaseIds), &resp)
	return resp, err
}
