package sirius

import (
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestAnomaliesForObject(t *testing.T) {
	afo := AnomaliesForObject{}

	anomaly := Anomaly{
		Status:        AnomalyDetected,
		FieldName:     "firstNames",
		FieldOwnerUid: "1",
	}

	afo.AddAnomaly(anomaly)

	assert.Equal(t, []Anomaly{anomaly}, afo.GetAnomaliesForFieldWithStatus("firstNames", "detected"))
}

func TestAnomaliesForSection(t *testing.T) {
	afs := AnomaliesForSection{}

	anomaly := Anomaly{
		Status:        AnomalyDetected,
		FieldName:     "firstNames",
		FieldOwnerUid: "1",
	}

	assert.False(t, afs.HasAnomalies())

	afs.AddAnomalyToObject(anomaly)

	expected := AnomaliesForObject{
		Uid: "1",
		Anomalies: map[ObjectFieldName][]Anomaly{
			"firstNames": {anomaly},
		},
	}

	assert.Equal(t, expected, *afs.GetAnomaliesForObject("1"))
	assert.True(t, afs.HasAnomalies())
}

func TestAnomalyDisplay(t *testing.T) {
	ad := AnomalyDisplay{}

	anomaly := Anomaly{
		Status:        AnomalyDetected,
		FieldName:     "firstNames",
		FieldOwnerUid: "1",
	}

	assert.False(t, ad.HasAnomalies())
	ad.AddAnomalyToSection(DonorSection, anomaly)
	assert.True(t, ad.HasAnomalies())

	expected := AnomaliesForSection{
		Section: DonorSection,
		Objects: map[ObjectUid]AnomaliesForObject{
			"1": {
				Uid: "1",
				Anomalies: map[ObjectFieldName][]Anomaly{
					"firstNames": {anomaly},
				},
			},
		},
	}
	assert.Equal(t, expected, *ad.GetAnomaliesForSection("donor"))
}

func TestAnomalyDisplay_SectionHasAnomalies(t *testing.T) {
	ad := AnomalyDisplay{
		AnomaliesBySection: map[AnomalyDisplaySection]AnomaliesForSection{
			DonorSection: {
				Section: DonorSection,
				Objects: map[ObjectUid]AnomaliesForObject{
					ObjectUid("1"): {},
				},
			},
		},
	}

	assert.True(t, ad.SectionHasAnomalies("donor"))
}

func TestAnomaliesForSection_GetAnomaliesForObject(t *testing.T) {
	afs := AnomaliesForSection{
		Section: AnomalyDisplaySection("donor"),
		Objects: nil,
	}

	assert.Equal(t, &AnomaliesForObject{}, afs.GetAnomaliesForObject("donor"))
}

func TestAnomaliesForObject_GetAnomaliesForFieldWithStatus(t *testing.T) {
	anomaly := Anomaly{
		Id:     1,
		Status: AnomalyDetected,
	}

	afo := AnomaliesForObject{
		Uid: ObjectUid("1"),
		Anomalies: map[ObjectFieldName][]Anomaly{
			"firstNames": {anomaly},
		},
	}

	assert.Equal(t, []Anomaly{anomaly}, afo.GetAnomaliesForFieldWithStatus("firstNames", "detected"))
	assert.Equal(t, []Anomaly(nil), afo.GetAnomaliesForFieldWithStatus("lastName", "detected"))
}

func TestAnomalyDisplay_Group(t *testing.T) {
	ad := AnomalyDisplay{}

	lpa := LpaStoreData{
		Donor: LpaStoreDonor{
			LpaStorePerson: LpaStorePerson{
				Uid: "1",
			},
		},
		Attorneys: []LpaStoreAttorney{
			{
				LpaStorePerson: LpaStorePerson{
					Uid: "2",
				},
				Status:          shared.ActiveAttorneyStatus.String(),
				AppointmentType: shared.OriginalAppointmentType.String(),
			},
		},
		// to test that no anomalies are returned for this section
		CertificateProvider: LpaStoreCertificateProvider{
			LpaStorePerson: LpaStorePerson{
				Uid: "4",
			},
		},
	}

	anomalies := []Anomaly{
		// donor (2 anomalies on one field)
		{
			Status:        AnomalyDetected,
			FieldName:     ObjectFieldName("firstNames"),
			FieldOwnerUid: ObjectUid("1"),
		},
		{
			Status:        AnomalyDetected,
			FieldName:     ObjectFieldName("firstNames"),
			FieldOwnerUid: ObjectUid("1"),
		},
		// attorneys
		{
			Status:        AnomalyDetected,
			FieldName:     ObjectFieldName("lastName"),
			FieldOwnerUid: ObjectUid("2"),
		},
		// root
		{
			Status:        AnomalyDetected,
			FieldName:     ObjectFieldName("howAttorneysMakeDecisions"),
			FieldOwnerUid: ObjectUid(""),
		},
		{
			Status:        AnomalyDetected,
			FieldName:     ObjectFieldName("whenTheLpaCanBeUsed"),
			FieldOwnerUid: ObjectUid(""),
		},
		{
			Status:        AnomalyDetected,
			FieldName:     ObjectFieldName("lifeSustainingTreatmentOption"),
			FieldOwnerUid: ObjectUid(""),
		},
	}

	ad.Group(&lpa, anomalies)

	expectedDonorAnomalies := AnomaliesForSection{
		Section: DonorSection,
		Objects: map[ObjectUid]AnomaliesForObject{
			ObjectUid("1"): {
				Uid: ObjectUid("1"),
				Anomalies: map[ObjectFieldName][]Anomaly{
					ObjectFieldName("firstNames"): {
						{
							Status:        AnomalyDetected,
							FieldName:     ObjectFieldName("firstNames"),
							FieldOwnerUid: ObjectUid("1"),
						},
						{
							Status:        AnomalyDetected,
							FieldName:     ObjectFieldName("firstNames"),
							FieldOwnerUid: ObjectUid("1"),
						},
					},
				},
			},
		},
	}

	expectedAttorneyAnomalies := AnomaliesForSection{
		Section: AttorneysSection,
		Objects: map[ObjectUid]AnomaliesForObject{
			ObjectUid("2"): {
				Uid: ObjectUid("2"),
				Anomalies: map[ObjectFieldName][]Anomaly{
					ObjectFieldName("lastName"): {
						{
							Status:        AnomalyDetected,
							FieldName:     ObjectFieldName("lastName"),
							FieldOwnerUid: ObjectUid("2"),
						},
					},
				},
			},
		},
	}

	expectedRootAnomalies := AnomaliesForSection{
		Section: RootSection,
		Objects: map[ObjectUid]AnomaliesForObject{
			ObjectUid(""): {
				Uid: ObjectUid(""),
				Anomalies: map[ObjectFieldName][]Anomaly{
					ObjectFieldName("howAttorneysMakeDecisions"): {
						{
							Status:    AnomalyDetected,
							FieldName: ObjectFieldName("howAttorneysMakeDecisions"),
						},
					},
					ObjectFieldName("whenTheLpaCanBeUsed"): {
						{
							Status:        AnomalyDetected,
							FieldName:     ObjectFieldName("whenTheLpaCanBeUsed"),
							FieldOwnerUid: ObjectUid(""),
						},
					},
					ObjectFieldName("lifeSustainingTreatmentOption"): {
						{
							Status:        AnomalyDetected,
							FieldName:     ObjectFieldName("lifeSustainingTreatmentOption"),
							FieldOwnerUid: ObjectUid(""),
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expectedDonorAnomalies, *ad.GetAnomaliesForSection("donor"))
	assert.Equal(t, expectedAttorneyAnomalies, *ad.GetAnomaliesForSection("attorneys"))
	assert.Equal(t, expectedRootAnomalies, *ad.GetAnomaliesForSection("root"))
	assert.Equal(t, AnomaliesForSection{}, *ad.GetAnomaliesForSection("certificateProvider"))
}

func TestGetSectionForUid(t *testing.T) {
	lpa := LpaStoreData{
		Donor: LpaStoreDonor{
			LpaStorePerson: LpaStorePerson{
				Uid: "1",
			},
		},
		Attorneys: []LpaStoreAttorney{
			{
				LpaStorePerson: LpaStorePerson{
					Uid: "2",
				},
				Status:          shared.ActiveAttorneyStatus.String(),
				AppointmentType: shared.OriginalAppointmentType.String(),
			},
			{
				LpaStorePerson: LpaStorePerson{
					Uid: "3",
				},
				Status:          shared.InactiveAttorneyStatus.String(),
				AppointmentType: shared.ReplacementAppointmentType.String(),
			},
			{
				LpaStorePerson: LpaStorePerson{
					Uid: "4",
				},
				Status:          shared.ActiveAttorneyStatus.String(),
				AppointmentType: shared.ReplacementAppointmentType.String(),
			},
		},
		CertificateProvider: LpaStoreCertificateProvider{
			LpaStorePerson: LpaStorePerson{
				Uid: "5",
			},
		},
		PeopleToNotify: []LpaStorePersonToNotify{
			{
				LpaStorePerson{
					Uid: "6",
				},
			},
		},
	}

	assert.Equal(t, DonorSection, getSectionForUid(&lpa, "1"))
	assert.Equal(t, AttorneysSection, getSectionForUid(&lpa, "2"))
	assert.Equal(t, ReplacementAttorneysSection, getSectionForUid(&lpa, "3"))
	assert.Equal(t, AttorneysSection, getSectionForUid(&lpa, "4"))
	assert.Equal(t, CertificateProviderSection, getSectionForUid(&lpa, "5"))
	assert.Equal(t, PeopleToNotifySection, getSectionForUid(&lpa, "6"))
	assert.Equal(t, RootSection, getSectionForUid(&lpa, ""))
}

func TestGetHintTextForAnomalyField(t *testing.T) {
	tests := []struct {
		name             string
		anomalies        []Anomaly
		want             string
		whoHasTheAnomaly string
	}{
		{
			name: "LastNameMatchesDonor and LastNameMatchesAttorney",
			anomalies: []Anomaly{
				{RuleType: LastNameMatchesDonor},
				{RuleType: LastNameMatchesAttorney},
			},
			want: "Review last name - this matches the donor and at least one of the attorneys. Check certificate provider's eligibility",
		},
		{
			name: "Empty anomaly",
			anomalies: []Anomaly{
				{RuleType: Empty},
			},
			want: "Review certificate provider's last name",
		},
		{
			name: "LastNameMatchesAttorney only",
			anomalies: []Anomaly{
				{RuleType: LastNameMatchesAttorney},
			},
			want: "Review last name - this matches at least one of the attorney's. Check certificate provider's eligibility",
		},
		{
			name: "LastNameMatchesDonor only",
			anomalies: []Anomaly{
				{RuleType: LastNameMatchesDonor},
			},
			want: "Review last name - this matches the donor's. Check certificate provider's eligibility",
		},
		{
			name:      "No anomalies",
			anomalies: []Anomaly{},
			want:      "Review certificate provider's last name",
		},
		{
			name: "NoCountry for CP",
			anomalies: []Anomaly{
				{RuleType: NoCountry},
			},
			want: "Review address as there is no country",
		},
		{
			name: "Bad address for CP",
			anomalies: []Anomaly{
				{RuleType: InvalidAddress},
			},
			whoHasTheAnomaly: "certificate provider's",
			want:             "Review certificate provider's address",
		},
		{
			name: "SignedAt and IdCheck dates too dates too far apart",
			anomalies: []Anomaly{
				{RuleType: IdAndSignedDateFarApart},
			},
			want: "Review signature date - check this is within 6 months either side of the donor’s ID check",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			afo := &AnomaliesForObject{}
			got := afo.GetHintTextForAnomalyField(tt.anomalies, tt.whoHasTheAnomaly)
			if got != tt.want {
				t.Errorf("GetHintTextForAnomalyField() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestContainsAnomalyType(t *testing.T) {
	tests := []struct {
		name       string
		anomalies  []Anomaly
		queryType  AnomalyRuleType
		wantResult bool
	}{
		{
			name:       "empty anomalies list",
			anomalies:  []Anomaly{},
			queryType:  "TypeA",
			wantResult: false,
		},
		{
			name: "anomaly type exists",
			anomalies: []Anomaly{
				{RuleType: "TypeA"},
				{RuleType: "TypeB"},
			},
			queryType:  "TypeA",
			wantResult: true,
		},
		{
			name: "anomaly type does not exist",
			anomalies: []Anomaly{
				{RuleType: "TypeB"},
				{RuleType: "TypeC"},
			},
			queryType:  "TypeA",
			wantResult: false,
		},
		{
			name: "multiple anomalies, match later",
			anomalies: []Anomaly{
				{RuleType: "TypeB"},
				{RuleType: "TypeC"},
				{RuleType: "TypeA"},
			},
			queryType:  "TypeA",
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsAnomalyType(tt.anomalies, tt.queryType)
			if got != tt.wantResult {
				t.Errorf("containsAnomalyType() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}
