package helpers

import (
	"regexp"
	"strconv"
	"strings"
)

func UnsafeAtoi(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return val
}

func CleanAllCommas(str string) string {
	// THESE TWO COMMAS ARE DIFFERENT CHARS
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, "٬", "")
	return str
}

func ArabicToPersianChars(str string) string {
	str = convertArabicCharsToPersian(str)
	str = cleanArabicErab(str)
	return str
}

func convertArabicCharsToPersian(input string) string {
	arabicChars := []string{"ي", "ك", "ى"}
	persianChars := []string{"ی", "ک", "ی"}

	for i := range arabicChars {
		input = strings.ReplaceAll(input, arabicChars[i], persianChars[i])
	}
	return input
}

func cleanArabicErab(input string) string {
	erabs := []string{"ً", "ٌ", "ٍ", "َ", "ُ", "ِ", "ّ"}
	for _, erab := range erabs {
		input = strings.ReplaceAll(input, erab, "")
	}
	return input
}

func SubStringBetweenTwoRegEx(text string, startPattern string, endPattern string) string {
	regex := regexp.MustCompile(startPattern)
	startMatch := regex.FindStringIndex(text)

	if startMatch == nil {
		return ""
	}

	startIndex := startMatch[1]

	regex = regexp.MustCompile(endPattern)
	endMatch := regex.FindStringIndex(text[startIndex:])

	if endMatch == nil {
		return ""
	}

	endIndex := endMatch[1]
	return text[startIndex : startIndex+endIndex]
}

func RemoveLastCurlyBrace(text string) string {
	trimmedText := strings.TrimSpace(text)
	if len(trimmedText) > 0 && trimmedText[len(trimmedText)-1] == '}' {
		return trimmedText[:len(trimmedText)-1]
	}
	return strings.TrimSpace(trimmedText)
}

func WordNumberToNumber(str string) int {
	numbersMap := map[string]int{
		"صفر":    0,
		"یک":     1,
		"دو":     2,
		"سه":     3,
		"چهار":   4,
		"پنج":    5,
		"شش":     6,
		"هفت":    7,
		"هشت":    8,
		"نه":     9,
		"ده":     10,
		"یازده":  11,
		"دوازده": 12,
		"سیزده":  13,
		"چهارده": 14,
		"پانزده": 15,
		"شانزده": 16,
		"هفتده":  17,
		"هجده":   18,
		"نوزده":  19,
		"بیست":   20,
	}

	number, exist := numbersMap[str]
	if !exist {
		return 0
	}

	return number
}
