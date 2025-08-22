package api

import (
	"strings"
)

var decodeMap = map[string]string{
	"%3A": ":",
	"%2F": "/",
	"%3F": "?",
	"%23": "#",
	"%5B": "[",
	"%5D": "]",
	"%40": "@",
	"%21": "!",
	"%24": "$",
	"%26": "&",
	"%27": "'",
	"%28": "(",
	"%29": ")",
	"%2A": "*",
	"%2B": "+",
	"%2C": ",",
	"%3B": ";",
	"%3D": "=",
	"%25": "%",
	// Recognizing both %20 and + for ' '. Use "%2B" if '+' in the url
	"%20": " ",
	"+":   " ",
}

// decodeURL decodes any percent-encoded characters as defined here:
// https://developer.mozilla.org/en-US/docs/Glossary/Percent-encoding
func decodeURL(url string) string {
	for k, v := range decodeMap {
		if strings.Contains(url, k) {
			url = strings.ReplaceAll(url, k, v)
		}
	}
	return url
}
