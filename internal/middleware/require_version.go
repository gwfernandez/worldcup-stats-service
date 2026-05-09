package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireVersion restricts access to routes within a group to a specific API version.
// It retrieves the version from the context (set by Versioning middleware).
func RequireVersion(requiredVersion int) gin.HandlerFunc {
	return func(c *gin.Context) {
		versionValue, exists := c.Get(ContextVersionKey)
		if !exists {
			// This should not happen if Versioning middleware is registered globally
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API version not identified"})
			c.Abort()
			return
		}

		version, ok := versionValue.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid API version type in context"})
			c.Abort()
			return
		}

		if version != requiredVersion {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("this endpoint requires API version %d, but version %d was requested", requiredVersion, version),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
