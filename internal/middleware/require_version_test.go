package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequireVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success when version matches", func(t *testing.T) {
		router := gin.New()
		router.Use(Versioning()) // Required to set context
		router.GET("/test", RequireVersion(1), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(VersionHeader, "1")
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Error when version does not match", func(t *testing.T) {
		router := gin.New()
		router.Use(Versioning())
		router.GET("/test", RequireVersion(2), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(VersionHeader, "1")
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Internal server error when version context has invalid type", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set(ContextVersionKey, "not-an-int")
			c.Next()
		})
		router.GET("/test", RequireVersion(1), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("Internal server error when version context is missing", func(t *testing.T) {
		router := gin.New()
		router.GET("/test", RequireVersion(1), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
