package sirius

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPaymentTypes(t *testing.T) {
	testcases := map[string]PaymentType{
		"telephone": PaymentType("PHONE"),
		"online":    PaymentType("ONLINE"),
	}

	for str, expectedPaymentType := range testcases {
		t.Run(str, func(t *testing.T) {
			p, err := ParsePaymentType(str)

			assert.Nil(t, err)
			assert.Equal(t, expectedPaymentType, p)
		})
	}
}

func TestInvalidPaymentTypes(t *testing.T) {
	p, err := ParsePaymentType("invalid")

	assert.Equal(t, "could not parse payment type", err.Error())
	assert.Equal(t, PaymentType(""), p)
}
