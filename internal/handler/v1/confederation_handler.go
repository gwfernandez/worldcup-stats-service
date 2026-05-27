package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

// ConfederationHandler handles HTTP requests for confederations.
type ConfederationHandler struct {
	service service.ConfederationService
}

// NewConfederationHandler creates a new handler with the given service.
func NewConfederationHandler(svc service.ConfederationService) *ConfederationHandler {
	return &ConfederationHandler{service: svc}
}

// RegisterRoutes registers all confederation routes on the given router group.
func (h *ConfederationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	confederations := rg.Group("/confederations")
	{
		confederations.GET("", h.List)
		confederations.GET("/:code", h.GetByCode)
	}
}

// List godoc
// @Summary List all confederations
// @Produce json
// @Success 200 {array} domain.Confederation
// @Router /api/confederations [get]
func (h *ConfederationHandler) List(c *gin.Context) {
	confederations, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve confederations"})
		return
	}
	c.JSON(http.StatusOK, confederations)
}

// GetByCode godoc
// @Summary Get a confederation by code
// @Produce json
// @Param code path string true "Confederation code"
// @Success 200 {object} domain.Confederation
// @Router /api/confederations/{code} [get]
func (h *ConfederationHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	confederation, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve confederation"})
		return
	}

	c.JSON(http.StatusOK, confederation)
}

// isNotFoundError checks if the error is a not-found error from the service layer.
func isNotFoundError(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}
