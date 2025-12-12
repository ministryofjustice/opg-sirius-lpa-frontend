package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadble(t *testing.T) {
	assert.Equal(t, "First names", DonorFirstNamesChange.Readable())
	assert.Equal(t, "Last name", AttorneysLastNameChange.Readable())
	assert.Equal(t, "Post code", CertificateProviderAddressPostCodeChange.Readable())
}

func TestGetTemplate(t *testing.T) {
	assert.Equal(t, "history-updated-from-to", DonorFirstNamesChange.GetTemplate())
}

func TestGuessReadable(t *testing.T) {
	const Test LpaStoreChangeType = "/test/lpaStore/changeType"

	assert.Equal(t, "Test lpa store change type", Test.guessReadable())
}

func TestGetLpaStoreChangeTypeFromChange(t *testing.T) {
	testCases := map[string]struct {
		change     shared.LpaStoreChange
		changeType LpaStoreChangeType
	}{
		"Attorney DOB": {
			change: shared.LpaStoreChange{
				Key: "/attorneys/0/dateOfBirth",
				Old: "1960-01-01",
				New: "1960-01-10",
			},
			changeType: AttorneysDateOfBirthChange,
		},
		"Trust Corporations category": {
			change: shared.LpaStoreChange{
				Key: "/certificateProvider/signedAt",
				Old: "2025-01-01",
				New: "2025-01-10",
			},
			changeType: CertificateProviderSignedAtChange,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c := getLpaStoreChangeTypeFromChange(tc.change)

			assert.Equal(t, tc.changeType, c)
		})
	}
}
