package shared

import (
	"encoding/json"
)

type LpaEventSourceType int

const (
	LpaEventSourceTypeAddress LpaEventSourceType = iota
	LpaEventSourceTypeAttorney
	LpaEventSourceTypeCertificateProvider
	LpaEventSourceTypeCheckAll
	LpaEventSourceTypeClient
	LpaEventSourceTypeComplaint
	LpaEventSourceTypeCorrespondent
	LpaEventSourceTypeCrec
	LpaEventSourceTypeDeputy
	LpaEventSourceTypeDonor
	LpaEventSourceTypeEpa
	LpaEventSourceTypeHoldPeriod
	LpaEventSourceTypeIncomingDocument
	LpaEventSourceTypeInvestigation
	LpaEventSourceTypeLodgingChecklist
	LpaEventSourceTypeLpa
	LpaEventSourceTypeNote
	LpaEventSourceTypeNotifiedPerson
	LpaEventSourceTypeOrder
	LpaEventSourceTypeOutgoingDocument
	LpaEventSourceTypePayment
	LpaEventSourceTypePhoneNumber
	LpaEventSourceTypeReplacementAttorney
	LpaEventSourceTypeTask
	LpaEventSourceTypeTrustCorporation
	LpaEventSourceTypeUncheckAll
	LpaEventSourceTypeValidationCheck
	LpaEventSourceTypeWarning
	LpaEventSourceTypeEmpty
	LpaEventSourceTypeNotRecognised
)

var lpaEventSourceTypeMap = map[string]LpaEventSourceType{
	"Address":             LpaEventSourceTypeAddress,
	"Attorney":            LpaEventSourceTypeAttorney,
	"CertificateProvider": LpaEventSourceTypeCertificateProvider,
	"CheckAll":            LpaEventSourceTypeCheckAll,
	"Client":              LpaEventSourceTypeClient,
	"Complaint":           LpaEventSourceTypeComplaint,
	"Correspondent":       LpaEventSourceTypeCorrespondent,
	"Crec":                LpaEventSourceTypeCrec, //CREC??
	"Deputy":              LpaEventSourceTypeDeputy,
	"Donor":               LpaEventSourceTypeDonor,
	"Epa":                 LpaEventSourceTypeEpa, //EPA?
	"HoldPeriod":          LpaEventSourceTypeHoldPeriod,
	"IncomingDocument":    LpaEventSourceTypeIncomingDocument,
	"Investigation":       LpaEventSourceTypeInvestigation,
	"LodgingChecklist":    LpaEventSourceTypeLodgingChecklist,
	"Lpa":                 LpaEventSourceTypeLpa, //LPA?
	"Note":                LpaEventSourceTypeNote,
	"NotifiedPerson":      LpaEventSourceTypeNotifiedPerson,
	"Order":               LpaEventSourceTypeOrder,
	"OutgoingDocument":    LpaEventSourceTypeOutgoingDocument,
	"Payment":             LpaEventSourceTypePayment,
	"PhoneNumber":         LpaEventSourceTypePhoneNumber,
	"ReplacementAttorney": LpaEventSourceTypeReplacementAttorney,
	"Task":                LpaEventSourceTypeTask,
	"TrustCorporation":    LpaEventSourceTypeTrustCorporation,
	"UncheckAll":          LpaEventSourceTypeUncheckAll,
	"ValidationCheck":     LpaEventSourceTypeValidationCheck,
	"Warning":             LpaEventSourceTypeWarning,
	"":                    LpaEventSourceTypeEmpty,
	"notRecognised":       LpaEventSourceTypeNotRecognised,
}

func (l LpaEventSourceType) String() string {
	return l.Key()
}

func (l LpaEventSourceType) Translation() string {
	switch l {
	case LpaEventSourceTypeAddress:
		return "Address"
	case LpaEventSourceTypeAttorney:
		return "Attorney"
	case LpaEventSourceTypeCertificateProvider:
		return "Certificate Provider"
	case LpaEventSourceTypeCheckAll:
		return "Check all"
	case LpaEventSourceTypeClient:
		return "Client"
	case LpaEventSourceTypeComplaint:
		return "Complaint"
	case LpaEventSourceTypeCorrespondent:
		return "Correspondent"
	case LpaEventSourceTypeCrec:
		return "CREC"
	case LpaEventSourceTypeDeputy:
		return "Deputy"
	case LpaEventSourceTypeDonor:
		return "Person (Create / Edit)"
	case LpaEventSourceTypeEpa:
		return "EPA (Create / Edit)"
	case LpaEventSourceTypeHoldPeriod:
		return "Hold Period"
	case LpaEventSourceTypeIncomingDocument:
		return "Incoming Document"
	case LpaEventSourceTypeInvestigation:
		return "Investigation"
	case LpaEventSourceTypeLodgingChecklist:
		return "Lodging Checklist"
	case LpaEventSourceTypeLpa:
		return "LPA (Create / Edit)"
	case LpaEventSourceTypeNote:
		return "Event"
	case LpaEventSourceTypeNotifiedPerson:
		return "Notified Person"
	case LpaEventSourceTypeOrder:
		return "Order"
	case LpaEventSourceTypeOutgoingDocument:
		return "Outbound document"
	case LpaEventSourceTypePayment:
		return "Payment"
	case LpaEventSourceTypePhoneNumber:
		return "Phone Number"
	case LpaEventSourceTypeReplacementAttorney:
		return "Replacement Attorney"
	case LpaEventSourceTypeTask:
		return "Task"
	case LpaEventSourceTypeTrustCorporation:
		return "Trust Corporation"
	case LpaEventSourceTypeUncheckAll:
		return "Uncheck all"
	case LpaEventSourceTypeValidationCheck:
		return "Validation Check"
	case LpaEventSourceTypeWarning:
		return "Warning"
	default:
		return "lpa event source type NOT RECOGNISED: " + l.String()
	}
}

func (l LpaEventSourceType) Key() string {
	switch l {
	case LpaEventSourceTypeAddress:
		return "Address"
	case LpaEventSourceTypeAttorney:
		return "Attorney"
	case LpaEventSourceTypeCertificateProvider:
		return "CertificateProvider"
	case LpaEventSourceTypeCheckAll:
		return "CheckAll"
	case LpaEventSourceTypeClient:
		return "Client"
	case LpaEventSourceTypeComplaint:
		return "Complaint"
	case LpaEventSourceTypeCorrespondent:
		return "Correspondent"
	case LpaEventSourceTypeCrec:
		return "Crec"
	case LpaEventSourceTypeDeputy:
		return "Deputy"
	case LpaEventSourceTypeDonor:
		return "Donor"
	case LpaEventSourceTypeEpa:
		return "Epa"
	case LpaEventSourceTypeHoldPeriod:
		return "HoldPeriod"
	case LpaEventSourceTypeIncomingDocument:
		return "IncomingDocument"
	case LpaEventSourceTypeInvestigation:
		return "Investigation"
	case LpaEventSourceTypeLodgingChecklist:
		return "LodgingChecklist"
	case LpaEventSourceTypeLpa:
		return "Lpa"
	case LpaEventSourceTypeNote:
		return "Event"
	case LpaEventSourceTypeNotifiedPerson:
		return "NotifiedPerson"
	case LpaEventSourceTypeOrder:
		return "Order"
	case LpaEventSourceTypeOutgoingDocument:
		return "OutboundDocument"
	case LpaEventSourceTypePayment:
		return "Payment"
	case LpaEventSourceTypePhoneNumber:
		return "PhoneNumber"
	case LpaEventSourceTypeReplacementAttorney:
		return "ReplacementAttorney"
	case LpaEventSourceTypeTask:
		return "Task"
	case LpaEventSourceTypeTrustCorporation:
		return "TrustCorporation"
	case LpaEventSourceTypeUncheckAll:
		return "UncheckAll"
	case LpaEventSourceTypeValidationCheck:
		return "ValidationCheck"
	case LpaEventSourceTypeWarning:
		return "Warning"
	default:
		return ""
	}
}

func ParseLpaEventSourceType(s string) LpaEventSourceType {
	value, ok := lpaEventSourceTypeMap[s]
	if !ok {
		return LpaEventSourceTypeNotRecognised
	}
	return value
}

func (l LpaEventSourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Key())
}

func (l *LpaEventSourceType) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*l = ParseLpaEventSourceType(s)
	return nil
}
