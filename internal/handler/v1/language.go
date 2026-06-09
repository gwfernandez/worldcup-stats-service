package v1

import (
	"net/http"
	"strings"
)

const (
	defaultLanguage = "es"
	englishLanguage = "en"
)

// resolveLanguage returns the supported API language from the Accept-Language header.
func resolveLanguage(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Accept-Language"))
	if header == "" {
		return defaultLanguage
	}

	for _, rawPart := range strings.Split(header, ",") {
		part := strings.TrimSpace(rawPart)
		if part == "" {
			continue
		}

		tag := strings.ToLower(strings.TrimSpace(strings.Split(part, ";")[0]))
		base := strings.Split(tag, "-")[0]
		switch base {
		case defaultLanguage, englishLanguage:
			return base
		}
	}

	return defaultLanguage
}
