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
		c := getLpaStoreCategoryFromChangeType(ct)

		return c
	}

	return UnknownCategory
}

func getLpaStoreCategoryFromChangeType(ct LpaStoreChangeType) LpaStoreCategory {
	switch ct {
	case DonorFirstNamesChange,
		DonorLastNameChange,
		DonorOtherNamesKnowByChange,
		DonorDateOfBirthChange,
		DonorEmailChange,
		DonorAddressLine1Change,
		DonorAddressLine2Change,
		DonorAddressLine3Change,
		DonorAddressTownChange,
		DonorAddressPostCodeChange,
		DonorAddressCountryChange:
		return DonorCategory
	case CertificateProviderFirstNamesChange,
		CertificateProviderLastNameChange,
		CertificateProviderEmailChange,
		CertificateProviderPhoneChange,
		CertificateProviderAddressLine1Change,
		CertificateProviderAddressLine2Change,
		CertificateProviderAddressLine3Change,
		CertificateProviderAddressTownChange,
		CertificateProviderAddressPostCodeChange,
		CertificateProviderAddressCountryChange,
		CertificateProviderSignedAtChange,
		CertificateProviderIdentityCheckedAtChange,
		CertificateProviderIdentityCheckTypeChange:
		return CertificateProvidersCategory
	}

	return UnknownCategory
}
