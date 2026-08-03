package sirius

import (
	"fmt"
)

type NotifiedPerson struct {
	Person
	NoticeGivenDate DateString `json:"noticeGivenDate,omitempty"`
	CaseId          int        `json:"caseId,omitempty"`
}

func (c *Client) CreateNotifiedPerson(ctx Context, caseId int, notifiedPerson NotifiedPerson) error {
	notifiedPerson.CaseId = caseId
	notifiedPerson.PersonType = "NotifiedPerson"

	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/persons"), []NotifiedPerson{notifiedPerson}, nil)
}

func (c *Client) UpdateNotifiedPerson(ctx Context, notifiedPersonId int, notifiedPerson NotifiedPerson) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/notified-persons/%d", notifiedPersonId), notifiedPerson, nil)
}
