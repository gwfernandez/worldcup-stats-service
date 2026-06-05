package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	// VersionHeader is the custom header used for API versioning
	VersionHeader = "API-Version"
	// VersionUsedHeader is the response header indicating the version used
	VersionUsedHeader = "API-Version-Used"
	// DefaultVersion is the fallback version if no header is provided
	DefaultVersion = 1
	// ContextVersionKey is the key used to store the version in the Gin context
	ContextVersionKey = "api_version"
)

// Versioning extracts the API version from the API-Version header.
// If the header is missing, it defaults to version 1.
// It also sets the API-Version-Used header in the response.
func Versioning() gin.HandlerFunc {
	return func(c *gin.Context) {
		versionStr := c.GetHeader(VersionHeader)
		version := DefaultVersion

		if versionStr != "" {
			v, err := strconv.Atoi(versionStr)
			if err != nil || v <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid API version format"})
				c.Abort()
				return
			}
			version = v
		}

		// Store version in context for downstream handlers/middlewares
		c.Set(ContextVersionKey, version)

		// Set the version used in the response header
		c.Header(VersionUsedHeader, strconv.Itoa(version))

		c.Next()
	}
}
