package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/jendrix/worldcup-stats-service/internal/middleware"
)

func TestConfigureCORS_AllowsConfiguredOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	configureCORS(router, []string{"http://localhost:5173"})
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodOptions, "/health", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	req.Header.Set("Access-Control-Request-Headers", middleware.VersionHeader)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), http.MethodGet)
	assert.True(t, containsHeaderName(w.Header().Get("Access-Control-Allow-Headers"), middleware.VersionHeader))

	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	assert.True(t, containsHeaderName(w.Header().Get("Access-Control-Expose-Headers"), middleware.VersionUsedHeader))
}

func TestConfigureCORS_SkipsMiddlewareWithoutOrigins(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	configureCORS(router, nil)
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func containsHeaderName(headers string, expected string) bool {
	for _, header := range strings.Split(headers, ",") {
		if strings.EqualFold(strings.TrimSpace(header), expected) {
			return true
		}
	}

	return false
}
