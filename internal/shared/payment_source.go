package shared

var paymentSourceActionMap = map[string]string{
	"PHONE":         "paid over the phone",
	"ONLINE":        "paid via GOV.UK Pay (card payment)",
	"MAKE":          "paid through Make an LPA",
	"MIGRATED":      "payment was migrated",
	"FEE_REDUCTION": "fee reduction",
	"CHEQUE":        "paid by cheque",
	"OTHER":         "paid through other method",
}

func PaymentSourceToAction(paymentSource string) string {
	if val, ok := paymentSourceActionMap[paymentSource]; ok {
		return val
	}
	return "paid through other method"
}
