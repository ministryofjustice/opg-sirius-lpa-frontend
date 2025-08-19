package shared

import "strings"

type CaseStatus int

const (
	Unknown CaseStatus = iota
	Registered
	Perfect
	StatutoryWaitingPeriod
	InProgress
	Pending
	PaymentPending
	ReducedFeesPending
	Draft
	Cancelled
	Rejected
	Revoked
	Withdrawn
	ReturnUnpaid
	Deleted
	DoNotRegister
	Expired
	CannotRegister
	DeRegistered
	Suspended
)

func CaseStatusParseStatus(s string) CaseStatus {
	switch strings.ToLower(s) {
	case "registered":
		return Registered
	case "perfect":
		return Perfect
	case "statutory waiting period":
		return StatutoryWaitingPeriod
	case "in progress":
		return InProgress
	case "pending":
		return Pending
	case "payment pending":
		return PaymentPending
	case "reduced fees pending":
		return ReducedFeesPending
	case "draft":
		return Draft
	case "cancelled":
		return Cancelled
	case "rejected":
		return Rejected
	case "revoked":
		return Revoked
	case "withdrawn":
		return Withdrawn
	case "return - unpaid":
		return ReturnUnpaid
	case "deleted":
		return Deleted
	case "do not register":
		return DoNotRegister
	case "expired":
		return Expired
	case "cannot register":
		return CannotRegister
	case "de-registered":
		return DeRegistered
	default:
		return Unknown
	}
}

func (cs CaseStatus) CaseStatusColour() string {
	switch cs {
	case Registered:
		return "green"
	case Perfect:
		return "turquoise"
	case StatutoryWaitingPeriod:
		return "yellow"
	case InProgress:
		return "light-blue"
	case Pending, PaymentPending, ReducedFeesPending:
		return "blue"
	case Draft:
		return "purple"
	case Cancelled, Rejected, Revoked, Withdrawn, ReturnUnpaid,
		Deleted, DoNotRegister, Expired, CannotRegister, DeRegistered:
		return "red"
	default:
		return "grey"
	}
}

func (cs CaseStatus) String() string {
	switch cs {
	case Draft:
		return "Draft"
	case InProgress:
		return "In progress"
	case StatutoryWaitingPeriod:
		return "Statutory waiting period"
	case Registered:
		return "Registered"
	case Suspended:
		return "Suspended"
	case DoNotRegister:
		return "Do not register"
	case Expired:
		return "Expired"
	case CannotRegister:
		return "Cannot register"
	case Cancelled:
		return "Cancelled"
	case DeRegistered:
		return "De-registered"
	default:
		return "draft"
	}
}
