package sirius

import "github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"

type AnomalyDisplaySection string

const (
	RootSection                 AnomalyDisplaySection = "root"
	DonorSection                AnomalyDisplaySection = "donor"
	CertificateProviderSection  AnomalyDisplaySection = "certificateProvider"
	AttorneysSection            AnomalyDisplaySection = "attorneys"
	ReplacementAttorneysSection AnomalyDisplaySection = "replacementAttorneys"
	PeopleToNotifySection       AnomalyDisplaySection = "peopleToNotify"
)

// AnomalyDisplay - Anomalies for the whole LPA details page
type AnomalyDisplay struct {
	AnomaliesBySection map[AnomalyDisplaySection]AnomaliesForSection
}

func (ad *AnomalyDisplay) AddAnomalyToSection(s AnomalyDisplaySection, a Anomaly) {
	if ad.AnomaliesBySection == nil {
		ad.AnomaliesBySection = map[AnomalyDisplaySection]AnomaliesForSection{}
	}

	anomaliesForSection, ok := ad.AnomaliesBySection[s]
	if !ok {
		anomaliesForSection = AnomaliesForSection{Section: s}
	}

	anomaliesForSection.AddAnomalyToObject(a)
	ad.AnomaliesBySection[s] = anomaliesForSection
}

func (ad *AnomalyDisplay) GetAnomaliesForSection(s string) *AnomaliesForSection {
	anomalies := ad.AnomaliesBySection[AnomalyDisplaySection(s)]
	return &anomalies
}

// Group - Split raw anomalies across the sections of the LPA details page
func (ad *AnomalyDisplay) Group(lpa *LpaStoreData, anomalies []Anomaly) *AnomalyDisplay {
	var s AnomalyDisplaySection
	for _, a := range anomalies {
		s = getSectionForUid(lpa, a.FieldOwnerUid)
		ad.AddAnomalyToSection(s, a)
	}
	return ad
}

func (ad *AnomalyDisplay) HasAnomalies() bool {
	return len(ad.AnomaliesBySection) > 0
}

func (ad *AnomalyDisplay) SectionHasAnomalies(s string) bool {
	section := ad.GetAnomaliesForSection(s)
	return section.HasAnomalies()
}

// AnomaliesForSection - Anomalies for a section of the LPA details page
type AnomaliesForSection struct {
	Section AnomalyDisplaySection

	// key is the UID of the object (fieldOwnerUid) which has the anomalies;
	// if an object has no anomalies, it will have no key in this map;
	// if no object has anomalies, the map will be empty
	Objects map[ObjectUid]AnomaliesForObject
}

func (afs *AnomaliesForSection) AddAnomalyToObject(a Anomaly) {
	if afs.Objects == nil {
		afs.Objects = map[ObjectUid]AnomaliesForObject{}
	}

	anomaliesForObject, ok := afs.Objects[a.FieldOwnerUid]
	if !ok {
		anomaliesForObject = AnomaliesForObject{Uid: a.FieldOwnerUid}
	}

	anomaliesForObject.AddAnomaly(a)
	afs.Objects[a.FieldOwnerUid] = anomaliesForObject
}

func (afs *AnomaliesForSection) GetAnomaliesForObject(uid string) *AnomaliesForObject {
	anomaliesForObject, ok := afs.Objects[ObjectUid(uid)]
	if !ok {
		return &AnomaliesForObject{}
	}
	return &anomaliesForObject
}

func (afs *AnomaliesForSection) HasAnomalies() bool {
	return len(afs.Objects) > 0
}

// AnomaliesForObject - Anomalies for an individual object (donor, attorney etc.) within a section
type AnomaliesForObject struct {
	// fieldOwnerUid for the object which has the anomalies
	Uid ObjectUid

	// map from field names to the anomalies for that field
	Anomalies map[ObjectFieldName][]Anomaly
}

func (afo *AnomaliesForObject) AddAnomaly(a Anomaly) {
	if afo.Anomalies == nil {
		afo.Anomalies = map[ObjectFieldName][]Anomaly{}
	}
	anomalies, ok := afo.Anomalies[a.FieldName]

	if !ok {
		anomalies = []Anomaly{}
	}

	afo.Anomalies[a.FieldName] = append(anomalies, a)
}

func (afo *AnomaliesForObject) GetAnomaliesForFieldWithStatus(fieldName string, status string) []Anomaly {
	var anomaliesWithStatus []Anomaly

	anomalies, ok := afo.Anomalies[ObjectFieldName(fieldName)]
	if !ok {
		return anomaliesWithStatus
	}

	wantedStatus := AnomalyStatus(status)
	for _, anomaly := range anomalies {
		if anomaly.Status == wantedStatus {
			anomaliesWithStatus = append(anomaliesWithStatus, anomaly)
		}
	}

	return anomaliesWithStatus
}

// getSectionForUid - Map a UID to an object inside an LPA and return which section it's in
func getSectionForUid(lpa *LpaStoreData, uid ObjectUid) AnomalyDisplaySection {
	if uid == "" {
		return RootSection
	} else if ObjectUid(lpa.Donor.Uid) == uid {
		return DonorSection
	} else if ObjectUid(lpa.CertificateProvider.Uid) == uid {
		return CertificateProviderSection
	} else {
		for _, attorney := range lpa.Attorneys {
			if ObjectUid(attorney.Uid) == uid {
				if attorney.Status == shared.ActiveAttorneyStatus.String() {
					return AttorneysSection
				}

				if attorney.Status == shared.InactiveAttorneyStatus.String() &&
					attorney.AppointmentType == shared.ReplacementAppointmentType.String() {
					return ReplacementAttorneysSection
				}
			}
		}

		for _, person := range lpa.PeopleToNotify {
			if ObjectUid(person.Uid) == uid {
				return PeopleToNotifySection
			}
		}
	}

	// UID is not matched, so assume it applies to the top-level fields of the LPA
	return RootSection
}
