package utils

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func FormatDateTimeToGetDailySchedule() string {
	d := time.Now()
	dayOfWeek := d.Weekday()

	// Check if it's Sunday (Weekday 0)
	if dayOfWeek == time.Sunday {
		// Return the formatted date for the previous day
		prevDay := d.AddDate(0, 0, -1)
		return fmt.Sprintf("%d-%02d-%02d", prevDay.Year(), prevDay.Month(), prevDay.Day())
	}

	// Return the formatted date for the current day
	return fmt.Sprintf("%d-%02d-%02d", d.Year(), d.Month(), d.Day())
}

func TodayFormatted() string {
	today := time.Now()

	dd := today.Day()
	mm := int(today.Month())
	yyyy := today.Year()

	ddStr := fmt.Sprintf("%02d", dd)
	mmStr := fmt.Sprintf("%02d", mm)

	date := ddStr + "/" + mmStr + "/" + fmt.Sprintf("%d", yyyy)
	return date
}

func IsTomorrow(dateString string) bool {
	layout := "15:04 02/01/2006"
	t, err := time.Parse(layout, dateString)
	if err != nil {
		log.Error().Msgf("Error parsing date string: %s", dateString)
		return false
	}

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	return t.Year() == tomorrow.Year() && t.YearDay() == tomorrow.YearDay()
}
