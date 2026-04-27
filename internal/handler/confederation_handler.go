package handler

import (
	"net/http"
	"strconv"
	"strings"

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
		confederations.GET("/:id", h.GetByID)
		confederations.POST("", h.Create)
		confederations.PUT("/:id", h.Update)
		confederations.DELETE("/:id", h.Delete)
	}
}

// List godoc
// @Summary List all confederations
// @Produce json
// @Success 200 {array} domain.Confederation
// @Router /api/v1/confederations [get]
func (h *ConfederationHandler) List(c *gin.Context) {
	confederations, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve confederations"})
		return
	}
	c.JSON(http.StatusOK, confederations)
}

// GetByID godoc
// @Summary Get a confederation by ID
// @Produce json
// @Param id path int true "Confederation ID"
// @Success 200 {object} domain.Confederation
// @Router /api/v1/confederations/{id} [get]
func (h *ConfederationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	confederation, err := h.service.GetByID(c.Request.Context(), id)
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

// Create godoc
// @Summary Create a new confederation
// @Accept json
// @Produce json
// @Param body body domain.CreateConfederationRequest true "Confederation data"
// @Success 201 {object} domain.Confederation
// @Router /api/v1/confederations [post]
func (h *ConfederationHandler) Create(c *gin.Context) {
	var req domain.CreateConfederationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	confederation, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if isDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "a confederation with this code already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create confederation"})
		return
	}

	c.JSON(http.StatusCreated, confederation)
}

// Update godoc
// @Summary Update a confederation
// @Accept json
// @Produce json
// @Param id path int true "Confederation ID"
// @Param body body domain.UpdateConfederationRequest true "Updated confederation data"
// @Success 200 {object} domain.Confederation
// @Router /api/v1/confederations/{id} [put]
func (h *ConfederationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	var req domain.UpdateConfederationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	confederation, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if isDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "a confederation with this code already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update confederation"})
		return
	}

	c.JSON(http.StatusOK, confederation)
}

// Delete godoc
// @Summary Delete a confederation
// @Param id path int true "Confederation ID"
// @Success 204
// @Router /api/v1/confederations/{id} [delete]
func (h *ConfederationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if isNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete confederation"})
		return
	}

	c.Status(http.StatusNoContent)
}

// isNotFoundError checks if the error is a not-found error from the service layer.
func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found")
}

// isDuplicateKeyError checks if the error is a unique constraint violation from PostgreSQL.
func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "23505")
}
