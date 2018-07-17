package helper

import (
	"time"
)

func GetNowTime() time.Time {
	loc, err := GetLoc()
	if err != nil {
		return time.Now().Add(time.Hour * time.Duration(7)) // Add 7 hours from local time
	}

	t := time.Now().In(loc)
	return t
}

func GetLoc() (loc *time.Location, err error) {
	// 1st Try Indonesia
	loc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// 2nd Try Thailand
		loc, err = time.LoadLocation("Asia/Bangkok")
		if err != nil {
			// 3rd Try Vietnam
			loc, err = time.LoadLocation("Asia/Saigon")
		}
	}

	return loc, err
}
