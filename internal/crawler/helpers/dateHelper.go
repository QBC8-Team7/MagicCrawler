package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func YearToAge(year string) int {
	intYear, err := strconv.Atoi(year)
	if err != nil {
		return 0
	}

	return (time.Now().Year() - 621) - intYear
}

func HumanDateToNormalDate(persianDate string) string {
	parts := strings.Fields(persianDate)

	return fmt.Sprintf("%s-%s-%s", parts[2], monthNameToMonthNumber(parts[1]), parts[0])
}

// private
func monthNameToMonthNumber(monthName string) string {
	monthName = strings.TrimSpace(monthName)

	months := map[string]int{
		"فروردین":  1,
		"اردیبهشت": 2,
		"خرداد":    3,
		"تیر":      4,
		"مرداد":    5,
		"شهریور":   6,
		"مهر":      7,
		"آبان":     8,
		"آذر":      9,
		"دی":       10,
		"بهمن":     11,
		"اسفند":    12,
	}

	number := months[monthName]
	if number < 10 {
		return fmt.Sprintf("0%d", number)
	} else {
		return fmt.Sprintf("%d", number)
	}
}
