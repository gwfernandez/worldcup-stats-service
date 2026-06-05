package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestVersioning(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Default version when header is missing", func(t *testing.T) {
		router := gin.New()
		router.Use(Versioning())
		router.GET("/test", func(c *gin.Context) {
			v, _ := c.Get(ContextVersionKey)
			c.JSON(http.StatusOK, gin.H{"version": v})
		})

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "1", resp.Header().Get("API-Version-Used"))
	})

	t.Run("Specific version from header", func(t *testing.T) {
		router := gin.New()
		router.Use(Versioning())
		router.GET("/test", func(c *gin.Context) {
			v, _ := c.Get(ContextVersionKey)
			c.JSON(http.StatusOK, gin.H{"version": v})
		})

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("API-Version", "2")
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "2", resp.Header().Get("API-Version-Used"))
	})

	t.Run("Invalid version format", func(t *testing.T) {
		router := gin.New()
		router.Use(Versioning())
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("API-Version", "abc")
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Negative version", func(t *testing.T) {
		router := gin.New()
		router.Use(Versioning())
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("API-Version", "-1")
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
