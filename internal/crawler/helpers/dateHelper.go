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

func PersianToMiladi(dashedDate string) (time.Time, error) {
	var year, month, day int
	_, err := fmt.Sscanf(dashedDate, "%d-%d-%d", &year, &month, &day)
	if err != nil {
		return time.Time{}, err
	}

	miladi := jalali_to_gregorian(year, month, day)

	return time.Date(miladi[0], time.Month(miladi[1]), miladi[2], 0, 0, 0, 0, time.UTC), nil
}

func jalali_to_gregorian(jy int, jm int, jd int) [3]int {
	var gy, gm, gd, days int
	var sal_a = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	jy += 1595
	days = -355668 + (365 * jy) + ((jy / 33) * 8) + (((jy % 33) + 3) / 4) + jd
	if jm < 7 {
		days += (jm - 1) * 31
	} else {
		days += ((jm - 7) * 30) + 186
	}
	gy = 400 * (days / 146097)
	days %= 146097
	if days > 36524 {
		days--
		gy += 100 * (days / 36524)
		days %= 36524
		if days >= 365 {
			days++
		}
	}
	gy += 4 * (days / 1461)
	days %= 1461
	if days > 365 {
		gy += (days - 1) / 365
		days = (days - 1) % 365
	}
	gd = days + 1
	if (gy%4 == 0 && gy%100 != 0) || (gy%400 == 0) {
		sal_a[2] = 29
	}
	gm = 0
	for gm < 13 && gd > sal_a[gm] {
		gd -= sal_a[gm]
		gm++
	}
	return [3]int{gy, gm, gd}
}
