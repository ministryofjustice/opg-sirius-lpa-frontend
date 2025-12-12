package server

import (
	"regexp"
	"strings"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type LpaStoreChangeType string

const (
	UnknownChange LpaStoreChangeType = ""

	DonorFirstNamesChange        LpaStoreChangeType = "/donor/firstNames"
	DonorOtherNamesKnowByChange  LpaStoreChangeType = "/donor/otherNamesKnownBy"
	DonorLastNameChange          LpaStoreChangeType = "/donor/lastName"
	DonorDateOfBirthChange       LpaStoreChangeType = "/donor/dateOfBirth"
	DonorEmailChange             LpaStoreChangeType = "/donor/email"
	DonorAddressLine1Change      LpaStoreChangeType = "/donor/address/line1"
	DonorAddressLine2Change      LpaStoreChangeType = "/donor/address/line2"
	DonorAddressLine3Change      LpaStoreChangeType = "/donor/address/line3"
	DonorAddressTownChange       LpaStoreChangeType = "/donor/address/town"
	DonorAddressPostCodeChange   LpaStoreChangeType = "/donor/address/postcode"
	DonorAddressCountryChange    LpaStoreChangeType = "/donor/address/country"
	DonorIdentityCheckedAtChange LpaStoreChangeType = "/donor/identityCheck/checkedAt"
	DonorIdentityCheckTypeChange LpaStoreChangeType = "/donor/identityCheck/type"

	AuthorisedSignatoryFirstNamesChange LpaStoreChangeType = "/authorisedSignatory/firstNames"
	AuthorisedSignatoryLastNameChange   LpaStoreChangeType = "/authorisedSignatory/lastName"

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
		CertificateProviderAddressCountryChange:
		return "Country"

	case AttorneysSignedAtChange,
		CertificateProviderSignedAtChange:
		return "Signed at"

	case CertificateProviderIdentityCheckedAtChange:
		return "CP checked at"
	}

	return l.guessReadable()
}

func (l LpaStoreChangeType) GetTemplate() string {
	switch l {
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
		DonorAddressCountryChange,
		DonorIdentityCheckedAtChange,
		DonorIdentityCheckTypeChange:
		return "history-updated-from-to"
	}

	return "generic-update"
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

func getLpaStoreChangeTypeFromChange(lse shared.LpaStoreChange) LpaStoreChangeType {
	m := regexp.MustCompile("/[0-9]+")
	k := m.ReplaceAllString(lse.Key, "")

	return LpaStoreChangeType(k)
}
