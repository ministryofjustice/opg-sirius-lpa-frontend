package server

import (
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestReadable(t *testing.T) {
	testCases := map[string]struct {
		category LpaStoreCategory
		readable string
	}{
		"Donor category": {
			category: DonorCategory,
			readable: "Donor",
		},
		"Trust Corporations category": {
			category: TrustCorporationsCategory,
			readable: "Trust Corporations",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := tc.category.Readable()

			assert.Equal(t, tc.readable, r)
		})
	}
}

func TestGetLpaStoreCategoryFromChanges(t *testing.T) {
	testCases := map[string]struct {
		changes  []shared.LpaStoreChange
		category LpaStoreCategory
	}{
		"Donor first names": {
			changes: []shared.LpaStoreChange{{
				Key:      "/donor/firstNames",
				Old:      "Arthur",
				New:      "Benjamin",
				Template: "test-template",
				Readable: "First names",
			}},
			category: DonorCategory,
		},
		"CP address line 1": {
			changes: []shared.LpaStoreChange{{
				Key:      "/certificateProvider/address/line1",
				Old:      "Arthur",
				New:      "Benjamin",
				Template: "test-template",
				Readable: "First names",
			}},
			category: CertificateProvidersCategory,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c := getLpaStoreCategoryFromChanges(tc.changes)

			assert.Equal(t, tc.category, c)
		})
	}
}

func TestGetLpaStoreCategoryFromChangeType(t *testing.T) {
	testCases := map[string]struct {
		key              LpaStoreChangeType
		lpaStoreCategory LpaStoreCategory
	}{
		"Donor last name": {
			key:              "/donor/lastName",
			lpaStoreCategory: DonorCategory,
		},
		"Certificate provider": {
			key:              "/certificateProvider/email",
			lpaStoreCategory: CertificateProvidersCategory,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c := getLpaStoreCategoryFromChangeType(tc.key)

			assert.Equal(t, tc.lpaStoreCategory, c)
		})
	}
}
