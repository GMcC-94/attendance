package helpers

import "time"

func CalculateAge(dob time.Time) int {
	loc, _ := time.LoadLocation("Europe/Dublin")
	now := time.Now().In(loc)

	years := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		years--
	}

	return years
}
