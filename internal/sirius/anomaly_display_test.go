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

	assert.Equal(t, []Anomaly{anomaly}, afo.GetAnomaliesForField("firstNames"))
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
	}

	anomalies := []Anomaly{
		// donor (2 anomalies on one field)
		{
			Status:        AnomalyDetected,
			FieldName:     "firstNames",
			FieldOwnerUid: "1",
		},
		{
			Status:        AnomalyDetected,
			FieldName:     "firstNames",
			FieldOwnerUid: "1",
		},
	}

	ad.Group(&lpa, anomalies)

	expectedDonorAnomalies := AnomaliesForSection{
		Section: DonorSection,
		Objects: map[ObjectUid]AnomaliesForObject{
			"1": {
				Uid: "1",
				Anomalies: map[ObjectFieldName][]Anomaly{
					"firstNames": {
						{
							Status:        AnomalyDetected,
							FieldName:     "firstNames",
							FieldOwnerUid: "1",
						},
						{
							Status:        AnomalyDetected,
							FieldName:     "firstNames",
							FieldOwnerUid: "1",
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expectedDonorAnomalies, *ad.GetAnomaliesForSection("donor"))
}
