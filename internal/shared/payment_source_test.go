package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentSourceToAction(t *testing.T) {
	tests := map[string]string{
		"PHONE":         "paid over the phone",
		"ONLINE":        "paid via GOV.UK Pay (card payment)",
		"MAKE":          "paid through Make an LPA",
		"MIGRATED":      "payment was migrated",
		"FEE_REDUCTION": "fee reduction",
		"CHEQUE":        "paid by cheque",
		"OTHER":         "paid through other method",
	}

	for input, expected := range tests {
		result := PaymentSourceToAction(input)
		assert.Equal(t, expected, result, "paymentSource(%q) should return %q", input, expected)
	}
}
