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
)

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
}

func (cs CaseStatus) ReadableString() string {
	switch cs {
	case CaseStatusTypeDraft:
		return "Draft"
	case CaseStatusTypeInProgress:
		return "In progress"
	case CaseStatusTypeStatutoryWaitingPeriod:
		return "Statutory waiting period"
	case CaseStatusTypeDoNotRegister:
		return "Do not register"
	case CaseStatusTypeExpired:
		return "Expired"
	case CaseStatusTypeRegistered:
		return "Registered"
	case CaseStatusTypeCannotRegister:
		return "Cannot register"
	case CaseStatusTypeCancelled:
		return "Cancelled"
	case CaseStatusTypeDeRegistered:
		return "De-registered"
	case CaseStatusTypeSuspended:
		return "Suspended"
	case CaseStatusTypePerfect:
		return "Perfect"
	case CaseStatusTypePending:
		return "Pending"
	case CaseStatusTypePaymentPending:
		return "Payment Pending"
	case CaseStatusTypeReducedFeesPending:
		return "Reduced Fees Pending"
	case CaseStatusTypeRejected:
		return "Rejected"
	case CaseStatusTypeWithdrawn:
		return "Withdrawn"
	case CaseStatusTypeReturnUnpaid:
		return "Return - unpaid"
	case CaseStatusTypeDeleted:
		return "Deleted"
	case CaseStatusTypeRevoked:
		return "Revoked"
	case CaseStatusTypeImperfect:
		return "Imperfect"
	case CaseStatusTypeInvalid:
		return "Invalid"
	default:
		return ""
	}
}

func (cs CaseStatus) StringForApi() string {
	switch cs {
	case CaseStatusTypeDraft:
		return "draft"
	case CaseStatusTypeInProgress:
		return "in-progress"
	case CaseStatusTypeStatutoryWaitingPeriod:
		return "statutory-waiting-period"
	case CaseStatusTypeDoNotRegister:
		return "do-not-register"
	case CaseStatusTypeExpired:
		return "expired"
	case CaseStatusTypeRegistered:
		return "registered"
	case CaseStatusTypeCannotRegister:
		return "cannot-register"
	case CaseStatusTypeCancelled:
		return "cancelled"
	case CaseStatusTypeDeRegistered:
		return "de-registered"
	case CaseStatusTypeSuspended:
		return "suspended"
	default:
		return ""
	}
}

func (cs CaseStatus) Colour() string {
	switch cs {
	case CaseStatusTypeRegistered:
		return "green"
	case CaseStatusTypePerfect:
		return "turquoise"
	case CaseStatusTypeStatutoryWaitingPeriod:
		return "yellow"
	case CaseStatusTypeInProgress:
		return "light-blue"
	case CaseStatusTypePending, CaseStatusTypePaymentPending, CaseStatusTypeReducedFeesPending:
		return "blue"
	case CaseStatusTypeDraft:
		return "purple"
	case CaseStatusTypeCancelled, CaseStatusTypeRejected, CaseStatusTypeRevoked, CaseStatusTypeWithdrawn, CaseStatusTypeReturnUnpaid,
		CaseStatusTypeDeleted, CaseStatusTypeDoNotRegister, CaseStatusTypeExpired, CaseStatusTypeCannotRegister, CaseStatusTypeDeRegistered:
		return "red"
	default:
		return "grey"
	}
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
	switch cs {
	case CaseStatusTypeDraft:
		return true
	default:
		return false
	}
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

func (cs *CaseStatus) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*cs = ParseCaseStatusType(s)
	return nil
}
