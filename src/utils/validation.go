package utils

import "regexp"

// ValidateEmail checks if a string is a valid email address
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}

// ValidateURL checks if a string is a valid URL
func ValidateURL(url string) bool {
	re := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:/[a-zA-Z0-9\-\._~:/?#[\]@!$&'()*+,;=]*)?$`)
	return re.MatchString(url)
}

// IsNumeric checks if a string contains only numeric characters
func IsNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}
