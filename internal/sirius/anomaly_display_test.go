package sirius

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestGroupAnomalies(t *testing.T) {
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
				Status: "active",
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

	assert.Equal(t, expectedDonorAnomalies, *ad.GetAnomaliesForSection("donor"))
	assert.Equal(t, expectedAttorneyAnomalies, *ad.GetAnomaliesForSection("attorneys"))
}
