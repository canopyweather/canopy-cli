package utils

import "time"

func GetDatesInRange(startDate, endDate string) ([]string, error) {
	// Date format constant
	const dateFormat = "2006-01-02"

	// Parse the start and end dates
	start, err := time.Parse(dateFormat, startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(dateFormat, endDate)
	if err != nil {
		return nil, err
	}

	// Loop through each date and append to the array as string
	var dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format(dateFormat))
	}

	return dates, nil
}

func ValidateDate(date string) bool {
	const dateFormat = "2006-01-02" // Reference time in Go
	_, err := time.Parse(dateFormat, date)
	return err == nil
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func ParseDate(date, format string) (time.Time, error) {
	return time.Parse(format, date)
}
