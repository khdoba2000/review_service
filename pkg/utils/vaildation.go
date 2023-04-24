package utils

func IsValidRating(rating uint8) bool {
	return rating >= 1 && rating <= 5
}
