package service

import "time"

func CalculateAge(birthDate, refDate time.Time) int {
	age := refDate.Year() - birthDate.Year()
	if refDate.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}
