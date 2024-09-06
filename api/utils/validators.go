package api

// IsValidRange ensures that an API range input is within the valid range of
// acceptable ranges
func IsValidRange(period int) bool {
	return period == 7 || period == 30 || period == 90
}
