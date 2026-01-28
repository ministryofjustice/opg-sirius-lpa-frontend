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
	SourceType shared.LpaEventSourceType `json:"sourceType"`
	Total      int                       `json:"total"`
}

func (c *Client) GetEvents(ctx Context, donorId string, caseIds []string, sourceTypes []string, sortBy string) (LpaEventsResponse, error) {
	var resp LpaEventsResponse

	query := ""
	for i, caseId := range caseIds {
		if i == 0 {
			query = "filter=case:" + caseId
		} else {
			query += ",case:" + caseId
		}
	}

	for _, sourceType := range sourceTypes {
		if query == "" {
			query = "filter=sourceType:" + sourceType
		} else {
			query += ",sourceType:" + sourceType
		}
	}

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/persons/%s/events?%s&sort=id:%s&limit=999", donorId, query, sortBy), &resp)
	return resp, err
}
