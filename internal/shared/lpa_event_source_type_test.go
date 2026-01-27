package shared

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLpaEventSourceTypeTranslation(t *testing.T) {
	tests := []struct {
		input    LpaEventSourceType
		expected string
	}{
		{input: LpaEventSourceTypeAddress, expected: "Address"},
		{input: LpaEventSourceTypeAttorney, expected: "Attorney"},
		{input: LpaEventSourceTypeCertificateProvider, expected: "Certificate Provider"},
		{input: LpaEventSourceTypeCheckAll, expected: "Check all"},
		{input: LpaEventSourceTypeClient, expected: "Client"},
		{input: LpaEventSourceTypeComplaint, expected: "Complaint"},
		{input: LpaEventSourceTypeCorrespondent, expected: "Correspondent"},
		{input: LpaEventSourceTypeCrec, expected: "CREC"},
		{input: LpaEventSourceTypeDeputy, expected: "Deputy"},
		{input: LpaEventSourceTypeDonor, expected: "Person (Create / Edit)"},
		{input: LpaEventSourceTypeEpa, expected: "EPA (Create / Edit)"},
		{input: LpaEventSourceTypeHoldPeriod, expected: "Hold Period"},
		{input: LpaEventSourceTypeIncomingDocument, expected: "Incoming Document"},
		{input: LpaEventSourceTypeInvestigation, expected: "Investigation"},
		{input: LpaEventSourceTypeLodgingChecklist, expected: "Lodging Checklist"},
		{input: LpaEventSourceTypeLpa, expected: "LPA (Create / Edit)"},
		{input: LpaEventSourceTypeNote, expected: "Event"},
		{input: LpaEventSourceTypeNotifiedPerson, expected: "Notified Person"},
		{input: LpaEventSourceTypeOrder, expected: "Order"},
		{input: LpaEventSourceTypeOutgoingDocument, expected: "Outbound document"},
		{input: LpaEventSourceTypePayment, expected: "Payment"},
		{input: LpaEventSourceTypePhoneNumber, expected: "Phone Number"},
		{input: LpaEventSourceTypeReplacementAttorney, expected: "Replacement Attorney"},
		{input: LpaEventSourceTypeTask, expected: "Task"},
		{input: LpaEventSourceTypeTrustCorporation, expected: "Trust Corporation"},
		{input: LpaEventSourceTypeUncheckAll, expected: "Uncheck all"},
		{input: LpaEventSourceTypeValidationCheck, expected: "Validation Check"},
		{input: LpaEventSourceTypeWarning, expected: "Warning"},
	}

	for _, tc := range tests {
		got := tc.input.Translation()
		if got != tc.expected {
			t.Errorf("LpaEventSourceType(%v) = %v, expected %v", tc.input, got, tc.expected)
		}
	}
}

func TestLpaEventSourceTypeTranslationDefault(t *testing.T) {
	assert.Equal(t, "lpa event source type NOT RECOGNISED: ", LpaEventSourceTypeNotRecognised.Translation())
}

func TestLpaEventSourceTypeKey(t *testing.T) {
	tests := []struct {
		input    LpaEventSourceType
		expected string
	}{
		{input: LpaEventSourceTypeAddress, expected: "Address"},
		{input: LpaEventSourceTypeAttorney, expected: "Attorney"},
		{input: LpaEventSourceTypeCertificateProvider, expected: "CertificateProvider"},
		{input: LpaEventSourceTypeCheckAll, expected: "CheckAll"},
		{input: LpaEventSourceTypeClient, expected: "Client"},
		{input: LpaEventSourceTypeComplaint, expected: "Complaint"},
		{input: LpaEventSourceTypeCorrespondent, expected: "Correspondent"},
		{input: LpaEventSourceTypeCrec, expected: "Crec"},
		{input: LpaEventSourceTypeDeputy, expected: "Deputy"},
		{input: LpaEventSourceTypeDonor, expected: "Donor"},
		{input: LpaEventSourceTypeEpa, expected: "Epa"},
		{input: LpaEventSourceTypeHoldPeriod, expected: "HoldPeriod"},
		{input: LpaEventSourceTypeIncomingDocument, expected: "IncomingDocument"},
		{input: LpaEventSourceTypeInvestigation, expected: "Investigation"},
		{input: LpaEventSourceTypeLodgingChecklist, expected: "LodgingChecklist"},
		{input: LpaEventSourceTypeLpa, expected: "Lpa"},
		{input: LpaEventSourceTypeNote, expected: "Event"},
		{input: LpaEventSourceTypeNotifiedPerson, expected: "NotifiedPerson"},
		{input: LpaEventSourceTypeOrder, expected: "Order"},
		{input: LpaEventSourceTypeOutgoingDocument, expected: "OutboundDocument"},
		{input: LpaEventSourceTypePayment, expected: "Payment"},
		{input: LpaEventSourceTypePhoneNumber, expected: "PhoneNumber"},
		{input: LpaEventSourceTypeReplacementAttorney, expected: "ReplacementAttorney"},
		{input: LpaEventSourceTypeTask, expected: "Task"},
		{input: LpaEventSourceTypeTrustCorporation, expected: "TrustCorporation"},
		{input: LpaEventSourceTypeUncheckAll, expected: "UncheckAll"},
		{input: LpaEventSourceTypeValidationCheck, expected: "ValidationCheck"},
		{input: LpaEventSourceTypeWarning, expected: "Warning"},
	}

	for _, tc := range tests {
		got := tc.input.Key()
		if got != tc.expected {
			t.Errorf("LpaEventSourceType(%v) = %v, expected %v", tc.input, got, tc.expected)
		}
	}
}

func TestLpaEventSourceTypeStringDefault(t *testing.T) {
	assert.Equal(t, "", LpaEventSourceTypeNotRecognised.Key())
}

func TestParseLpaEventSourceType(t *testing.T) {
	tests := []struct {
		input    string
		expected LpaEventSourceType
	}{
		{input: "Address", expected: LpaEventSourceTypeAddress},
		{input: "Attorney", expected: LpaEventSourceTypeAttorney},
		{input: "CertificateProvider", expected: LpaEventSourceTypeCertificateProvider},
		{input: "CheckAll", expected: LpaEventSourceTypeCheckAll},
		{input: "Client", expected: LpaEventSourceTypeClient},
		{input: "Complaint", expected: LpaEventSourceTypeComplaint},
		{input: "Correspondent", expected: LpaEventSourceTypeCorrespondent},
		{input: "Crec", expected: LpaEventSourceTypeCrec}, //CREC?
		{input: "Deputy", expected: LpaEventSourceTypeDeputy},
		{input: "Donor", expected: LpaEventSourceTypeDonor},
		{input: "Epa", expected: LpaEventSourceTypeEpa}, //EPA?
		{input: "HoldPeriod", expected: LpaEventSourceTypeHoldPeriod},
		{input: "IncomingDocument", expected: LpaEventSourceTypeIncomingDocument},
		{input: "Investigation", expected: LpaEventSourceTypeInvestigation},
		{input: "LodgingChecklist", expected: LpaEventSourceTypeLodgingChecklist},
		{input: "Lpa", expected: LpaEventSourceTypeLpa}, //LPA?
		{input: "Note", expected: LpaEventSourceTypeNote},
		{input: "NotifiedPerson", expected: LpaEventSourceTypeNotifiedPerson},
		{input: "Order", expected: LpaEventSourceTypeOrder},
		{input: "OutgoingDocument", expected: LpaEventSourceTypeOutgoingDocument},
		{input: "Payment", expected: LpaEventSourceTypePayment},
		{input: "PhoneNumber", expected: LpaEventSourceTypePhoneNumber},
		{input: "ReplacementAttorney", expected: LpaEventSourceTypeReplacementAttorney},
		{input: "Task", expected: LpaEventSourceTypeTask},
		{input: "TrustCorporation", expected: LpaEventSourceTypeTrustCorporation},
		{input: "UncheckAll", expected: LpaEventSourceTypeUncheckAll},
		{input: "ValidationCheck", expected: LpaEventSourceTypeValidationCheck},
		{input: "Warning", expected: LpaEventSourceTypeWarning},
		{input: "", expected: LpaEventSourceTypeEmpty},
		{input: "notRecognised", expected: LpaEventSourceTypeNotRecognised},
	}

	for _, tc := range tests {
		got := ParseLpaEventSourceType(tc.input)
		if got != tc.expected {
			t.Errorf("string(%v) = %v, expected %v", tc.input, got, tc.expected)
		}
	}
}

func TestParseLpaEventSourceTypeErrors(t *testing.T) {
	got := ParseLpaEventSourceType("---")
	assert.Equal(t, LpaEventSourceTypeNotRecognised, got)
}

func TestLpaEventSourceTypeMarshalJSON(t *testing.T) {
	tests := []struct {
		input    LpaEventSourceType
		expected string
	}{
		{input: LpaEventSourceTypeAddress, expected: `"Address"`},
		{input: LpaEventSourceTypeAttorney, expected: `"Attorney"`},
		{input: LpaEventSourceTypeCertificateProvider, expected: `"CertificateProvider"`},
		{input: LpaEventSourceTypeCheckAll, expected: `"CheckAll"`},
		{input: LpaEventSourceTypeClient, expected: `"Client"`},
		{input: LpaEventSourceTypeComplaint, expected: `"Complaint"`},
		{input: LpaEventSourceTypeCorrespondent, expected: `"Correspondent"`},
		{input: LpaEventSourceTypeCrec, expected: `"Crec"`},
		{input: LpaEventSourceTypeDeputy, expected: `"Deputy"`},
		{input: LpaEventSourceTypeDonor, expected: `"Donor"`},
		{input: LpaEventSourceTypeEpa, expected: `"Epa"`},
		{input: LpaEventSourceTypeHoldPeriod, expected: `"HoldPeriod"`},
		{input: LpaEventSourceTypeIncomingDocument, expected: `"IncomingDocument"`},
		{input: LpaEventSourceTypeInvestigation, expected: `"Investigation"`},
		{input: LpaEventSourceTypeLodgingChecklist, expected: `"LodgingChecklist"`},
		{input: LpaEventSourceTypeLpa, expected: `"Lpa"`},
		{input: LpaEventSourceTypeNote, expected: `"Event"`},
		{input: LpaEventSourceTypeNotifiedPerson, expected: `"NotifiedPerson"`},
		{input: LpaEventSourceTypeOrder, expected: `"Order"`},
		{input: LpaEventSourceTypeOutgoingDocument, expected: `"OutboundDocument"`},
		{input: LpaEventSourceTypePayment, expected: `"Payment"`},
		{input: LpaEventSourceTypePhoneNumber, expected: `"PhoneNumber"`},
		{input: LpaEventSourceTypeReplacementAttorney, expected: `"ReplacementAttorney"`},
		{input: LpaEventSourceTypeTask, expected: `"Task"`},
		{input: LpaEventSourceTypeTrustCorporation, expected: `"TrustCorporation"`},
		{input: LpaEventSourceTypeUncheckAll, expected: `"UncheckAll"`},
		{input: LpaEventSourceTypeValidationCheck, expected: `"ValidationCheck"`},
		{input: LpaEventSourceTypeWarning, expected: `"Warning"`},
	}

	for _, tc := range tests {
		got, _ := json.Marshal(tc.input)
		if string(got) != tc.expected {
			t.Errorf("LpaEventSourceType(%v) = %v, expected %v", tc.input, got, tc.expected)
		}
	}
}

func TestLpaEventSourceTypeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected LpaEventSourceType
	}{
		{input: `"Address"`, expected: LpaEventSourceTypeAddress},
		{input: `"Attorney"`, expected: LpaEventSourceTypeAttorney},
		{input: `"CertificateProvider"`, expected: LpaEventSourceTypeCertificateProvider},
		{input: `"CheckAll"`, expected: LpaEventSourceTypeCheckAll},
		{input: `"Client"`, expected: LpaEventSourceTypeClient},
		{input: `"Complaint"`, expected: LpaEventSourceTypeComplaint},
		{input: `"Correspondent"`, expected: LpaEventSourceTypeCorrespondent},
		{input: `"Crec"`, expected: LpaEventSourceTypeCrec},
		{input: `"Deputy"`, expected: LpaEventSourceTypeDeputy},
		{input: `"Donor"`, expected: LpaEventSourceTypeDonor},
		{input: `"Epa"`, expected: LpaEventSourceTypeEpa},
		{input: `"HoldPeriod"`, expected: LpaEventSourceTypeHoldPeriod},
		{input: `"IncomingDocument"`, expected: LpaEventSourceTypeIncomingDocument},
		{input: `"Investigation"`, expected: LpaEventSourceTypeInvestigation},
		{input: `"LodgingChecklist"`, expected: LpaEventSourceTypeLodgingChecklist},
		{input: `"Lpa"`, expected: LpaEventSourceTypeLpa},
		{input: `"Note"`, expected: LpaEventSourceTypeNote},
		{input: `"NotifiedPerson"`, expected: LpaEventSourceTypeNotifiedPerson},
		{input: `"Order"`, expected: LpaEventSourceTypeOrder},
		{input: `"OutgoingDocument"`, expected: LpaEventSourceTypeOutgoingDocument},
		{input: `"Payment"`, expected: LpaEventSourceTypePayment},
		{input: `"PhoneNumber"`, expected: LpaEventSourceTypePhoneNumber},
		{input: `"ReplacementAttorney"`, expected: LpaEventSourceTypeReplacementAttorney},
		{input: `"Task"`, expected: LpaEventSourceTypeTask},
		{input: `"TrustCorporation"`, expected: LpaEventSourceTypeTrustCorporation},
		{input: `"UncheckAll"`, expected: LpaEventSourceTypeUncheckAll},
		{input: `"ValidationCheck"`, expected: LpaEventSourceTypeValidationCheck},
		{input: `"Warning"`, expected: LpaEventSourceTypeWarning},
	}

	for _, tc := range tests {
		var est LpaEventSourceType
		x := []byte(tc.input)
		err := json.Unmarshal(x, &est)
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) returned error: %v", tc.input, err)
		}
		if est != tc.expected {
			t.Errorf("UnmarshalJSON(%s), expected %v", tc.input, tc.expected)
		}
	}
}

func TestLpaEventSourceTypeUnmarshalJSONErrors(t *testing.T) {
	var est LpaEventSourceType
	err := json.Unmarshal([]byte(`123`), &est)
	assert.Error(t, err)
}
