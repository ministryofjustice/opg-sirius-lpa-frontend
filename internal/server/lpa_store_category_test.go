package server

import (
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestCategoryReadable(t *testing.T) {
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
		"Attorneys mobile": {
			changes: []shared.LpaStoreChange{{
				Key:      "/attorneys/0/mobile",
				Old:      "07697 780428",
				New:      "07123 425784",
				Template: "test-template",
				Readable: "Mobile",
			}},
			category: AttorneysCategory,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c := getLpaStoreCategoryFromChanges(tc.changes)

			assert.Equal(t, tc.category, c)
		})
	}
}
