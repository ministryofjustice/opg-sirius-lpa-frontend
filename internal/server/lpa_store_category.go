package server

import "github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"

type LpaStoreCategory string

const (
	UnknownCategory              LpaStoreCategory = ""
	DonorCategory                LpaStoreCategory = "Donor"
	AttorneysCategory            LpaStoreCategory = "Attorneys"
	CertificateProvidersCategory LpaStoreCategory = "Certificate Providers"
	TrustCorporationsCategory    LpaStoreCategory = "Trust Corporations"
	DecisionsCategory            LpaStoreCategory = "Decisions"

	// Status ?
)

func (lsc LpaStoreCategory) Readable() string {
	return string(lsc)
}

func getLpaStoreCategoryFromChanges(lscs []shared.LpaStoreChange) LpaStoreCategory {
	for _, lsc := range lscs {
		ct := getLpaStoreChangeTypeFromChange(lsc)

		if c := ct.getCategory(); c != "" {
			return c
		}
	}

	return UnknownCategory
}
