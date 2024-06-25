package sirius

type AnomalyDisplaySection string

const (
	RootSection                AnomalyDisplaySection = "root"
	DonorSection               AnomalyDisplaySection = "donor"
	CertificateProviderSection AnomalyDisplaySection = "certificateProvider"
	AttorneysSection           AnomalyDisplaySection = "attorneys"
	PeopleToNotifySection      AnomalyDisplaySection = "peopleToNotify"
)

// AnomalyDisplay - Anomalies for the whole LPA details page
type AnomalyDisplay struct {
	AnomaliesBySection map[AnomalyDisplaySection]AnomaliesForSection
}

func (ad *AnomalyDisplay) AddAnomalyToSection(s AnomalyDisplaySection, a Anomaly) {
	anomaliesForSection, ok := ad.AnomaliesBySection[s]
	if !ok {
		ad.AnomaliesBySection = map[AnomalyDisplaySection]AnomaliesForSection{}
		anomaliesForSection = AnomaliesForSection{Section: s}
	}
	anomaliesForSection.AddAnomalyToObject(a)
	ad.AnomaliesBySection[s] = anomaliesForSection
}

func (ad *AnomalyDisplay) GetAnomaliesForSection(s AnomalyDisplaySection) AnomaliesForSection {
	return ad.AnomaliesBySection[s]
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

// AnomaliesForSection - Anomalies for a section of the LPA details page
type AnomaliesForSection struct {
	Section AnomalyDisplaySection

	// key is the UID of the object (fieldOwnerUid) which has the anomalies;
	// if an object has no anomalies, it will have no key in this map;
	// if no object has anomalies, the map will be empty
	Objects map[ObjectUid]AnomaliesForObject
}

func (afs *AnomaliesForSection) AddAnomalyToObject(a Anomaly) {
	anomaliesForObject, ok := afs.Objects[a.FieldOwnerUid]
	if !ok {
		afs.Objects = map[ObjectUid]AnomaliesForObject{}
		anomaliesForObject = AnomaliesForObject{Uid: a.FieldOwnerUid}
	}
	anomaliesForObject.AddAnomaly(a)
	afs.Objects[a.FieldOwnerUid] = anomaliesForObject
}

func (afs *AnomaliesForSection) GetAnomaliesForObject(uid ObjectUid) AnomaliesForObject {
	anomaliesForObject, ok := afs.Objects[uid]
	if !ok {
		return AnomaliesForObject{}
	}
	return anomaliesForObject
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
	anomalies, ok := afo.Anomalies[a.FieldName]
	if !ok {
		afo.Anomalies = map[ObjectFieldName][]Anomaly{}
		anomalies = []Anomaly{}
	}
	afo.Anomalies[a.FieldName] = append(anomalies, a)
}

func (afo *AnomaliesForObject) GetAnomaliesForField(fieldName ObjectFieldName) []Anomaly {
	anomalies, ok := afo.Anomalies[fieldName]
	if !ok {
		return []Anomaly{}
	}
	return anomalies
}

// getSectionForUid - Map a UID to an object inside an LPA and return which section it's in
func getSectionForUid(lpa *LpaStoreData, uid ObjectUid) AnomalyDisplaySection {
	if ObjectUid(lpa.Donor.Uid) == uid {
		return DonorSection
	} else if ObjectUid(lpa.CertificateProvider.Uid) == uid {
		return CertificateProviderSection
	} else {
		for _, attorney := range lpa.Attorneys {
			if ObjectUid(attorney.Uid) == uid {
				return AttorneysSection
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
