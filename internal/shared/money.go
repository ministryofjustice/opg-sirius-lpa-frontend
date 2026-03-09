package shared

import "fmt"

func FormatMonetaryValue(amount int) string {
	float := float64(amount)
	return FormatMonetaryFloat(float)
}

func FormatMonetaryFloat(amount float64) string {
	return fmt.Sprintf("%.2f", amount/100)
}
