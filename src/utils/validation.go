package utils

import "regexp"

// Pre-compiled regex for email validation - only compiled once at package initialization
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

// ValidateEmail checks if a string is a valid email address
func ValidateEmail(email string) bool {
	// Quick length check (shortest valid email is at least "a@b.xx")
	if len(email) < 6 {
		return false
	}

	// Fast check for @ symbol before applying regex
	hasAt := false
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			hasAt = true
			break
		}
	}

	if !hasAt {
		return false
	}

	return emailRegex.MatchString(email)
}

// Pre-compiled regex for URL validation - only compiled once at package initialization
var urlRegex = regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(?:/[a-zA-Z0-9\-\._~:/?#[\]@!$&'()*+,;=]*)?$`)

// ValidateURL checks if a string is a valid URL
func ValidateURL(urlStr string) bool {
	// Quick length check (shortest valid URL is at least "http://a.xx")
	if len(urlStr) < 10 {
		return false
	}

	// Fast prefix check before applying expensive regex
	if !(urlStr[0:7] == "http://" || urlStr[0:8] == "https://") {
		return false
	}

	return urlRegex.MatchString(urlStr)
}

// IsNumeric checks if a string contains only numeric characters
func IsNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Using direct byte access for maximum performance
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
