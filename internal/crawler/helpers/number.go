package helpers

import "strings"

func ToEnglishDigits(numberStr string) string {
	persianDigits := []string{"۰", "۱", "۲", "۳", "۴", "۵", "۶", "۷", "۸", "۹"}
	arabicDigits := []string{"٠", "١", "٢", "٣", "٤", "٥", "٦", "٧", "٨", "٩"}
	englishDigits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	for i := range persianDigits {
		numberStr = strings.ReplaceAll(numberStr, persianDigits[i], englishDigits[i])
	}
	for i := range arabicDigits {
		numberStr = strings.ReplaceAll(numberStr, arabicDigits[i], englishDigits[i])
	}
	return numberStr
}
