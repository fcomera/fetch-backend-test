package utils

import (
	"log"
	"time"
)


// Parses a date and time string to a time.Time type for time operations
func ParseDateAndTime(purchaseDate string, purchaseTime string) time.Time {
	dateTransformed, error := time.Parse("2006-01-02 15:04", purchaseDate+" "+purchaseTime)

	if error != nil {
		log.Println(error)
	} else {
		log.Println(dateTransformed)
	}

	return dateTransformed
}
