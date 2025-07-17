package timeutil

import "time"

var istLocation = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic("Failed to load IST location: " + err.Error())
	}
	return loc
}()

// NowIST returns current time in IST
func NowIST() time.Time {
	return time.Now().In(istLocation)
}

// NowISTFormatted returns current IST time as a formatted string
func NowISTFormatted() string {
	return NowIST().Format("2006-01-02T15:04:05")
}
