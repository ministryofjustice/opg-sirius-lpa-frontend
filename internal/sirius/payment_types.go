package sirius

import "errors"

type PaymentType string

const (
	PaymentTypePhone    = PaymentType("PHONE")
	PaymentTypeMake     = PaymentType("MAKE")
	PaymentTypeOnline   = PaymentType("ONLINE")
	PaymentTypeOther    = PaymentType("OTHER")
	PaymentTypeMigrated = PaymentType("MIGRATED")
)

func ParsePaymentType(s string) (PaymentType, error) {
	switch s {
	case "telephone":
		return PaymentTypePhone, nil
	case "make":
		return PaymentTypeMake, nil
	case "online":
		return PaymentTypeOnline, nil
	case "other":
		return PaymentTypeOther, nil
	case "migrated":
		return PaymentTypeMigrated, nil
	}

	return PaymentType(""), errors.New("could not parse payment type")
}
