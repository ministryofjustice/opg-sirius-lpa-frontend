package sirius

import (
	"fmt"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type APIEvents []Event
type LpaEvents []LpaEvent

type LpaEventsResponse struct {
	Events   LpaEvents `json:"events"`
	Limit    int       `json:"limit"`
	Total    int       `json:"total"`
	Pages    any       `json:"pages"`
	Metadata any       `json:"metadata"`
}

type Event struct {
	Changes             []shared.LpaStoreChange `json:"changes"`
	CreatedOn           string                  `json:"createdOn"`
	Entity              any                     `json:"entity"`
	Hash                string                  `json:"hash"`
	ID                  string                  `json:"id"`
	Source              string                  `json:"source"`
	SourceType          string                  `json:"sourceType"`
	Type                string                  `json:"type"`
	User                EventUser               `json:"user"`
	UUID                string                  `json:"uuid"`
	Applied             string                  `json:"applied,omitempty"`
	DateTime            string                  `json:"dateTime,omitempty"`
	FormattedLpaStoreId string
	Category            string
}

type LpaEvent struct {
	Changes       any         `json:"changeSet,omitempty"`
	CreatedOn     string      `json:"createdOn"`
	Entity        any         `json:"entity,omitempty"`
	Assignee      EventUser   `json:"assignee,omitempty"`
	User          EventUser   `json:"user,omitempty"`
	Hash          string      `json:"hash"`
	OwningCase    OwningCase  `json:"owningCase,omitempty"`
	ID            int         `json:"id,omitempty"`
	UID           string      `json:"uid,omitempty"`
	SourceType    string      `json:"sourceType,omitempty"`
	Type          string      `json:"type,omitempty"`
	SourceTask    any         `json:"sourceTask,omitempty"`
	SourceAddress any         `json:"sourceAddress,omitempty"`
	SourceCase    any         `json:"sourceCase,omitempty"`
	SourcePerson  any         `json:"sourcePerson,omitempty"`
	SourcePhone   SourcePhone `json:"sourcePhoneNumber,omitempty"`
}

type SourcePhone struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Type        string `json:"type,omitempty"`
}

type OwningCase struct {
	ID          int    `json:"id,omitempty"`
	UID         string `json:"uId,omitempty"`
	CaseSubtype string `json:"caseSubtype,omitempty"`
}

//Lots of source item options, most likely missing

func (e Event) IsLpaStore() bool {
	return e.Source == "lpa_store"
}

type EventUser struct {
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	Deleted     bool   `json:"deleted,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

// GetCombinedEvents Gets combined events from both Sirius and LPA Store for digital LPAs
func (c *Client) GetCombinedEvents(ctx Context, uid string) (APIEvents, error) {
	var events APIEvents
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/events", uid), &events)
	return events, err
}

// will probably have to refactor for getting events across multiple cases
func (c *Client) GetEvents(ctx Context, donorId int, caseId int, sourceTypes []string, sortBy string) (LpaEvents, error) {
	var resp LpaEventsResponse
	filters := ""
	for _, t := range sourceTypes {
		filters += fmt.Sprintf(",sourceType:%s", t)
	}
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/persons/%d/events?filter=case:%d%s&sort=id:%s", donorId, caseId, filters, sortBy), &resp)
	return resp.Events, err
}
