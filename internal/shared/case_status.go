package shared

import (
	"encoding/json"
)

type CaseStatus int

const (
	CaseStatusTypeUnknown CaseStatus = iota
	CaseStatusTypeRegistered
	CaseStatusTypePerfect
	CaseStatusTypeStatutoryWaitingPeriod
	CaseStatusTypeInProgress
	CaseStatusTypePending
	CaseStatusTypePaymentPending
	CaseStatusTypeReducedFeesPending
	CaseStatusTypeDraft
	CaseStatusTypeCancelled
	CaseStatusTypeRejected
	CaseStatusTypeRevoked
	CaseStatusTypeWithdrawn
	CaseStatusTypeReturnUnpaid
	CaseStatusTypeDeleted
	CaseStatusTypeDoNotRegister
	CaseStatusTypeExpired
	CaseStatusTypeCannotRegister
	CaseStatusTypeDeRegistered
	CaseStatusTypeSuspended
	CaseStatusTypeImperfect
	CaseStatusTypeInvalid
	CaseStatusTypeWithCop
	CaseStatusTypeProcessing
)

type caseStatusMeta struct {
	Readable string
	API      string
	Colour   string
}

var caseStatusMetadata = map[CaseStatus]caseStatusMeta{
	CaseStatusTypeDraft:                  {"Draft", "draft", "purple"},
	CaseStatusTypeInProgress:             {"In progress", "in-progress", "light-blue"},
	CaseStatusTypeStatutoryWaitingPeriod: {"Statutory waiting period", "statutory-waiting-period", "yellow"},
	CaseStatusTypeDoNotRegister:          {"Do not register", "do-not-register", "red"},
	CaseStatusTypeExpired:                {"Expired", "expired", "red"},
	CaseStatusTypeRegistered:             {"Registered", "registered", "green"},
	CaseStatusTypeCannotRegister:         {"Cannot register", "cannot-register", "red"},
	CaseStatusTypeCancelled:              {"Cancelled", "cancelled", "red"},
	CaseStatusTypeDeRegistered:           {"De-registered", "de-registered", "red"},
	CaseStatusTypeSuspended:              {"Suspended", "suspended", "red"},
	CaseStatusTypePerfect:                {"Perfect", "", "turquoise"},
	CaseStatusTypePending:                {"Pending", "", "blue"},
	CaseStatusTypePaymentPending:         {"Payment Pending", "", "blue"},
	CaseStatusTypeReducedFeesPending:     {"Reduced Fees Pending", "", "blue"},
	CaseStatusTypeRejected:               {"Rejected", "", "red"},
	CaseStatusTypeWithdrawn:              {"Withdrawn", "", "red"},
	CaseStatusTypeReturnUnpaid:           {"Return - unpaid", "", "red"},
	CaseStatusTypeDeleted:                {"Deleted", "", "red"},
	CaseStatusTypeRevoked:                {"Revoked", "", "red"},
	CaseStatusTypeImperfect:              {"Imperfect", "", "grey"},
	CaseStatusTypeInvalid:                {"Invalid", "", "grey"},
	CaseStatusTypeWithCop:                {"With Cop", "", "grey"},
	CaseStatusTypeProcessing:             {"Processing", "processing", "grey"},
}

var caseStatusTypeMap = map[string]CaseStatus{
	"draft":                    CaseStatusTypeDraft,
	"Draft":                    CaseStatusTypeDraft,
	"In progress":              CaseStatusTypeInProgress,
	"in-progress":              CaseStatusTypeInProgress,
	"Statutory waiting period": CaseStatusTypeStatutoryWaitingPeriod,
	"statutory-waiting-period": CaseStatusTypeStatutoryWaitingPeriod,
	"Do not register":          CaseStatusTypeDoNotRegister,
	"do-not-register":          CaseStatusTypeDoNotRegister,
	"Expired":                  CaseStatusTypeExpired,
	"expired":                  CaseStatusTypeExpired,
	"Registered":               CaseStatusTypeRegistered,
	"registered":               CaseStatusTypeRegistered,
	"Cannot register":          CaseStatusTypeCannotRegister,
	"cannot-register":          CaseStatusTypeCannotRegister,
	"Cancelled":                CaseStatusTypeCancelled,
	"cancelled":                CaseStatusTypeCancelled,
	"De-registered":            CaseStatusTypeDeRegistered,
	"de-registered":            CaseStatusTypeDeRegistered,
	"Suspended":                CaseStatusTypeSuspended,
	"suspended":                CaseStatusTypeSuspended,
	"Perfect":                  CaseStatusTypePerfect,
	"Pending":                  CaseStatusTypePending,
	"Payment Pending":          CaseStatusTypePaymentPending,
	"Reduced Fees Pending":     CaseStatusTypeReducedFeesPending,
	"Rejected":                 CaseStatusTypeRejected,
	"Withdrawn":                CaseStatusTypeWithdrawn,
	"Return - unpaid":          CaseStatusTypeReturnUnpaid,
	"Deleted":                  CaseStatusTypeDeleted,
	"Revoked":                  CaseStatusTypeRevoked,
	"Imperfect":                CaseStatusTypeImperfect,
	"imperfect":                CaseStatusTypeImperfect,
	"Invalid":                  CaseStatusTypeInvalid,
	"invalid":                  CaseStatusTypeInvalid,
	"With Cop":                 CaseStatusTypeWithCop,
	"processing":               CaseStatusTypeProcessing,
}

func (cs CaseStatus) ReadableString() string {
	if meta, ok := caseStatusMetadata[cs]; ok {
		return meta.Readable
	}
	return ""
}

func (cs CaseStatus) StringForApi() string {
	if meta, ok := caseStatusMetadata[cs]; ok {
		return meta.API
	}
	return ""
}

func (cs CaseStatus) Colour() string {
	if meta, ok := caseStatusMetadata[cs]; ok && meta.Colour != "" {
		return meta.Colour
	}
	return "grey"
}

func (cs CaseStatus) IsValidStatusForObjection() bool {
	switch cs {
	case CaseStatusTypeInProgress, CaseStatusTypeDraft, CaseStatusTypeStatutoryWaitingPeriod:
		return true
	default:
		return false
	}
}

func (cs CaseStatus) IsDraft() bool {
	return cs == CaseStatusTypeDraft
}

func ParseCaseStatusType(s string) CaseStatus {
	value, ok := caseStatusTypeMap[s]
	if !ok {
		return CaseStatus(0)
	}
	return value
}

func (cs CaseStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(cs.ReadableString())
}

func (cs *CaseStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*cs = ParseCaseStatusType(s)
	return nil
}
