package sirius

import (
	"fmt"
)

type UpdateDecisions struct {
	WhenTheLpaCanBeUsed                         string `json:"whenTheLpaCanBeUsed,omitempty"`
	LifeSustainingTreatmentOption               string `json:"lifeSustainingTreatmentOption,omitempty"`
	HowAttorneysMakeDecisions                   string `json:"howAttorneysMakeDecisions,omitempty"`
	HowAttorneysMakeDecisionsDetails            string `json:"howAttorneysMakeDecisionsDetails,omitempty"`
	HowReplacementAttorneysStepIn               string `json:"howReplacementAttorneysStepIn,omitempty"`
	HowReplacementAttorneysStepInDetails        string `json:"howReplacementAttorneysStepInDetails,omitempty"`
	HowReplacementAttorneysMakeDecisions        string `json:"howReplacementAttorneysMakeDecisions,omitempty"`
	HowReplacementAttorneysMakeDecisionsDetails string `json:"howReplacementAttorneysMakeDecisionsDetails,omitempty"`
}

func (c *Client) UpdateDecisions(ctx Context, caseUID string, decisionDetails UpdateDecisions) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/decisions", caseUID), decisionDetails, nil)
}
