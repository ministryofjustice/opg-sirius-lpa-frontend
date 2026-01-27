package sirius

import (
	"fmt"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type LpaEventsResponse struct {
	Events   []LpaEvent    `json:"events"`
	Limit    int           `json:"limit"`
	Total    int           `json:"total"`
	Pages    Pages         `json:"pages"`
	Metadata EventMetaData `json:"metadata"`
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

type EventMetaData struct {
	CaseIds     any          `json:"caseIds"`
	SourceTypes []SourceType `json:"sourceTypes"`
}

type SourceType struct {
	SourceType string `json:"sourceType"`
	Total      int    `json:"total"`
}

func (c *Client) GetEvents(ctx Context, donorId string, caseIds []string, sourceTypes []string, sortBy string) (LpaEventsResponse, error) {
	var resp LpaEventsResponse

	//TODO: Simplify this logic into a query builder
	selectedCaseIds := ""
	for i, caseId := range caseIds {
		if i == 0 {
			selectedCaseIds = selectedCaseIds + "filter=case:" + caseId
		} else {
			selectedCaseIds = selectedCaseIds + ",case:" + caseId
		}
	}

	filters := ""
	for i, t := range sourceTypes {
		if len(caseIds) == 0 && i == 0 {
			filters = "filter=" + fmt.Sprintf("sourceType:%s", t)
		} else {
			filters += fmt.Sprintf(",sourceType:%s", t)
		}
	}

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/persons/%s/events?%s%s&sort=id:%s&limit=999", donorId, selectedCaseIds, filters, sortBy), &resp)
	return resp, err
}
