package helpers

import "strings"

func CleanPrice(price string) string {
	price = ToEnglishDigits(price)
	price = cleanTooman(price)
	price = CleanAllCommas(price)
	return strings.TrimSpace(price)

}

func cleanTooman(price string) string {
	return strings.ReplaceAll(price, "تومان", "")
}
