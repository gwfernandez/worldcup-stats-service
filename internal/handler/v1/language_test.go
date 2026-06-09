package v1

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveLanguage(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{name: "empty header defaults to spanish", expected: "es"},
		{name: "spanish header", header: "es", expected: "es"},
		{name: "english header", header: "en", expected: "en"},
		{name: "english regional header", header: "en-US", expected: "en"},
		{name: "unsupported header defaults to spanish", header: "pt-BR", expected: "es"},
		{name: "uses first supported language in priority list", header: "fr-CA, en;q=0.8, es;q=0.6", expected: "en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/confederations", nil)
			assert.NoError(t, err)
			if tt.header != "" {
				req.Header.Set("Accept-Language", tt.header)
			}

			assert.Equal(t, tt.expected, resolveLanguage(req))
		})
	}
}
