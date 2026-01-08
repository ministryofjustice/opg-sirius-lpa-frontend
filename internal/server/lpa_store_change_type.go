package server

import (
	"regexp"
	"strings"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type LpaStoreChangeType string

const (
	UnknownChange LpaStoreChangeType = ""

	DonorFirstNamesChange       LpaStoreChangeType = "/donor/firstNames"
	DonorOtherNamesKnowByChange LpaStoreChangeType = "/donor/otherNamesKnownBy"
	DonorLastNameChange         LpaStoreChangeType = "/donor/lastName"
	DonorDateOfBirthChange      LpaStoreChangeType = "/donor/dateOfBirth"
	DonorEmailChange            LpaStoreChangeType = "/donor/email"
	DonorAddressLine1Change     LpaStoreChangeType = "/donor/address/line1"
	DonorAddressLine2Change     LpaStoreChangeType = "/donor/address/line2"
	DonorAddressLine3Change     LpaStoreChangeType = "/donor/address/line3"
	DonorAddressTownChange      LpaStoreChangeType = "/donor/address/town"
	DonorAddressPostCodeChange  LpaStoreChangeType = "/donor/address/postcode"
	DonorAddressCountryChange   LpaStoreChangeType = "/donor/address/country"
	SignedAtChange              LpaStoreChangeType = "/signedAt"

	DonorIdentityCheckedAtChange LpaStoreChangeType = "/donor/identityCheck/checkedAt"
	DonorIdentityCheckTypeChange LpaStoreChangeType = "/donor/identityCheck/type"

	AuthorisedSignatoryFirstNamesChange LpaStoreChangeType = "/authorisedSignatory/firstNames"
	AuthorisedSignatoryLastNameChange   LpaStoreChangeType = "/authorisedSignatory/lastName"

	IndependentWitnessFirstNamesChange      LpaStoreChangeType = "/independentWitness/firstNames"
	IndependentWitnessLastNameChange        LpaStoreChangeType = "/independentWitness/lastName"
	IndependentWitnessAddressLine1Change    LpaStoreChangeType = "/independentWitness/address/line1"
	IndependentWitnessAddressLine2Change    LpaStoreChangeType = "/independentWitness/address/line2"
	IndependentWitnessAddressLine3Change    LpaStoreChangeType = "/independentWitness/address/line3"
	IndependentWitnessAddressTownChange     LpaStoreChangeType = "/independentWitness/address/town"
	IndependentWitnessAddressPostCodeChange LpaStoreChangeType = "/independentWitness/address/postcode"
	IndependentWitnessAddressCountryChange  LpaStoreChangeType = "/independentWitness/address/country"

	CertificateProviderFirstNamesChange        LpaStoreChangeType = "/certificateProvider/firstNames"
	CertificateProviderLastNameChange          LpaStoreChangeType = "/certificateProvider/lastName"
	CertificateProviderEmailChange             LpaStoreChangeType = "/certificateProvider/email"
	CertificateProviderPhoneChange             LpaStoreChangeType = "/certificateProvider/phone"
	CertificateProviderAddressLine1Change      LpaStoreChangeType = "/certificateProvider/address/line1"
	CertificateProviderAddressLine2Change      LpaStoreChangeType = "/certificateProvider/address/line2"
	CertificateProviderAddressLine3Change      LpaStoreChangeType = "/certificateProvider/address/line3"
	CertificateProviderAddressTownChange       LpaStoreChangeType = "/certificateProvider/address/town"
	CertificateProviderAddressPostCodeChange   LpaStoreChangeType = "/certificateProvider/address/postcode"
	CertificateProviderAddressCountryChange    LpaStoreChangeType = "/certificateProvider/address/country"
	CertificateProviderSignedAtChange          LpaStoreChangeType = "/certificateProvider/signedAt"
	CertificateProviderIdentityCheckedAtChange LpaStoreChangeType = "/certificateProvider/identityCheck/checkedAt"
	CertificateProviderIdentityCheckTypeChange LpaStoreChangeType = "/certificateProvider/identityCheck/type"

	AttorneysFirstNamesChange      LpaStoreChangeType = "/attorneys/firstNames"
	AttorneysLastNameChange        LpaStoreChangeType = "/attorneys/lastName"
	AttorneysDateOfBirthChange     LpaStoreChangeType = "/attorneys/dateOfBirth"
	AttorneysEmailChange           LpaStoreChangeType = "/attorneys/email"
	AttorneysMobileChange          LpaStoreChangeType = "/attorneys/mobile"
	AttorneysAddressLine1Change    LpaStoreChangeType = "/attorneys/address/line1"
	AttorneysAddressLine2Change    LpaStoreChangeType = "/attorneys/address/line2"
	AttorneysAddressLine3Change    LpaStoreChangeType = "/attorneys/address/line3"
	AttorneysAddressTownChange     LpaStoreChangeType = "/attorneys/address/town"
	AttorneysAddressPostCodeChange LpaStoreChangeType = "/attorneys/address/postcode"
	AttorneysAddressCountryChange  LpaStoreChangeType = "/attorneys/address/country"
	AttorneysSignedAtChange        LpaStoreChangeType = "/attorneys/signedAt"

	TrustCorporationNameChange               LpaStoreChangeType = "/trustCorporation/name"
	TrustCorporationEmailChange              LpaStoreChangeType = "/trustCorporation/email"
	TrustCorporationMobileChange             LpaStoreChangeType = "/trustCorporation/mobile"
	TrustCorporationCompanyNumberChange      LpaStoreChangeType = "/trustCorporation/companyNumber"
	TrustCorporationAddressLine1ChangeChange LpaStoreChangeType = "/trustCorporation/address/line1"
	TrustCorporationAddressLine2ChangeChange LpaStoreChangeType = "/trustCorporation/address/line2"
	TrustCorporationAddressLine3ChangeChange LpaStoreChangeType = "/trustCorporation/address/line3"
	TrustCorporationAddressTownChange        LpaStoreChangeType = "/trustCorporation/address/town"
	TrustCorporationAddressPostcodeChange    LpaStoreChangeType = "/trustCorporation/address/postcode"
	TrustCorporationAddressCountryChange     LpaStoreChangeType = "/trustCorporation/address/country"

	HowAttorneysMakeDecisionsChange                   LpaStoreChangeType = "/howAttorneysMakeDecisions"
	HowAttorneysMakeDecisionsDetailsChange            LpaStoreChangeType = "/howAttorneysMakeDecisionsDetails"
	HowReplacementAttorneysStepInChange               LpaStoreChangeType = "/howReplacementAttorneysStepIn"
	HowReplacementAttorneysStepInDetailsChange        LpaStoreChangeType = "/howReplacementAttorneysStepInDetails"
	HowReplacementAttorneysMakeDecisionsChange        LpaStoreChangeType = "/howReplacementAttorneysMakeDecisions"
	HowReplacementAttorneysMakeDecisionsDetailsChange LpaStoreChangeType = "/howReplacementAttorneysMakeDecisionsDetails"
	LifeSustainingTreatmentOptionChange               LpaStoreChangeType = "/lifeSustainingTreatmentOption"
	WhenTheLpaCanBeUsedChange                         LpaStoreChangeType = "/whenTheLpaCanBeUsed"
)

func (l LpaStoreChangeType) Readable() string {
	switch l {
	case DonorFirstNamesChange,
		AttorneysFirstNamesChange,
		CertificateProviderFirstNamesChange:
		return "First names"

	case DonorOtherNamesKnowByChange:
		return "Other names know by"

	case DonorLastNameChange,
		AttorneysLastNameChange,
		CertificateProviderLastNameChange:
		return "Last name"

	case DonorDateOfBirthChange,
		AttorneysDateOfBirthChange:
		return "Date of birth"

	case DonorEmailChange,
		AttorneysEmailChange,
		CertificateProviderEmailChange:
		return "Email"

	case AttorneysMobileChange:
		return "Mobile"

	case CertificateProviderPhoneChange:
		return "Phone"

	case DonorAddressLine1Change,
		AttorneysAddressLine1Change,
		CertificateProviderAddressLine1Change:
		return "Address line 1"

	case DonorAddressLine2Change,
		AttorneysAddressLine2Change,
		CertificateProviderAddressLine2Change:
		return "Address line 2"

	case DonorAddressLine3Change,
		AttorneysAddressLine3Change,
		CertificateProviderAddressLine3Change:
		return "Address line 3"

	case DonorAddressTownChange,
		AttorneysAddressTownChange,
		CertificateProviderAddressTownChange:
		return "Town"

	case DonorAddressPostCodeChange,
		AttorneysAddressPostCodeChange,
		CertificateProviderAddressPostCodeChange:
		return "Post code"

	case DonorAddressCountryChange,
		AttorneysAddressCountryChange,
		CertificateProviderAddressCountryChange,
		IndependentWitnessAddressCountryChange:
		return "Country"

	case AttorneysSignedAtChange,
		CertificateProviderSignedAtChange:
		return "Signed at"

	case CertificateProviderIdentityCheckedAtChange:
		return "CP checked at"

	case DonorIdentityCheckedAtChange:
		return "Identity checked at"

	case DonorIdentityCheckTypeChange:
		return "Identity checked type"

	default:
		return l.guessReadable()
	}
}

func (l LpaStoreChangeType) GetTemplate() string {
	if l == DonorDateOfBirthChange || l == AttorneysDateOfBirthChange {
		return "history-date-updated-from-to"
	}

	return "history-updated-from-to"
}

func (l LpaStoreChangeType) guessReadable() string {
	// Remove preceding slash
	m := regexp.MustCompile("^/")
	r := m.ReplaceAllString(string(l), "")

	// Replace slash with space
	m = regexp.MustCompile("/")
	r = m.ReplaceAllString(r, " ")

	// Separate by upper
	m = regexp.MustCompile("([A-Z])")
	r = m.ReplaceAllString(r, " $1")

	// Lower
	r = strings.ToLower(r)

	// Split and upper first letter
	s := strings.Split(r, "")
	s[0] = strings.ToUpper(s[0])
	r = strings.Join(s, "")

	return r
}

func (l LpaStoreChangeType) getCategory() LpaStoreCategory {
	var categoryChangeTypeMap = map[LpaStoreCategory][]LpaStoreChangeType{
		DonorCategory: {
			DonorFirstNamesChange,
			DonorLastNameChange,
			DonorOtherNamesKnowByChange,
			DonorDateOfBirthChange,
			DonorEmailChange,
			DonorAddressLine1Change,
			DonorAddressLine2Change,
			DonorAddressLine3Change,
			DonorAddressTownChange,
			DonorAddressPostCodeChange,
			DonorAddressCountryChange,
			SignedAtChange,
			DonorIdentityCheckedAtChange,
			DonorIdentityCheckTypeChange,
			AuthorisedSignatoryFirstNamesChange,
			AuthorisedSignatoryLastNameChange,
			IndependentWitnessFirstNamesChange,
			IndependentWitnessLastNameChange,
			IndependentWitnessAddressLine1Change,
			IndependentWitnessAddressLine2Change,
			IndependentWitnessAddressLine3Change,
			IndependentWitnessAddressTownChange,
			IndependentWitnessAddressPostCodeChange,
			IndependentWitnessAddressCountryChange,
		},
		CertificateProvidersCategory: {
			CertificateProviderFirstNamesChange,
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
			CertificateProviderIdentityCheckTypeChange,
		},
		AttorneysCategory: {
			AttorneysFirstNamesChange,
			AttorneysLastNameChange,
			AttorneysDateOfBirthChange,
			AttorneysEmailChange,
			AttorneysMobileChange,
			AttorneysAddressLine1Change,
			AttorneysAddressLine2Change,
			AttorneysAddressLine3Change,
			AttorneysAddressTownChange,
			AttorneysAddressPostCodeChange,
			AttorneysAddressCountryChange,
			AttorneysSignedAtChange,
		},
		TrustCorporationsCategory: {
			TrustCorporationNameChange,
			TrustCorporationEmailChange,
			TrustCorporationMobileChange,
			TrustCorporationCompanyNumberChange,
			TrustCorporationAddressLine1ChangeChange,
			TrustCorporationAddressLine2ChangeChange,
			TrustCorporationAddressLine3ChangeChange,
			TrustCorporationAddressTownChange,
			TrustCorporationAddressPostcodeChange,
			TrustCorporationAddressCountryChange,
		},
		DecisionsCategory: {
			HowAttorneysMakeDecisionsChange,
			HowAttorneysMakeDecisionsDetailsChange,
			HowReplacementAttorneysStepInChange,
			HowReplacementAttorneysStepInDetailsChange,
			HowReplacementAttorneysMakeDecisionsChange,
			HowReplacementAttorneysMakeDecisionsDetailsChange,
			LifeSustainingTreatmentOptionChange,
			WhenTheLpaCanBeUsedChange,
		},
	}

	for cat, cts := range categoryChangeTypeMap {
		for _, ct := range cts {
			if l == ct {
				return cat
			}
		}
	}

	return UnknownCategory
}

func getLpaStoreChangeTypeFromChange(lse shared.LpaStoreChange) LpaStoreChangeType {
	m := regexp.MustCompile("/[0-9]+")
	k := m.ReplaceAllString(lse.Key, "")

	return LpaStoreChangeType(k)
}
