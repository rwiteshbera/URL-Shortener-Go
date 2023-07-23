package helpers

import (
	"strings"
)

func RemoveDomainError(url string) string {
	cleanURL := strings.ToLower(url)
	cleanURL = strings.TrimSpace(cleanURL)
	cleanURL = strings.TrimPrefix(cleanURL, "http://")
	cleanURL = strings.TrimPrefix(cleanURL, "https://")
	cleanURL = strings.TrimPrefix(cleanURL, "www.")
	cleanURL = strings.TrimSuffix(cleanURL, "/")

	return cleanURL
}

func EnforceHTTP(url string) string {
	return "http://" + url
}
