package api

func IsValidRange(period int32) bool {
	return period == 7 || period == 30 || period == 90
}
