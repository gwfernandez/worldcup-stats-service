package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

// ScorerHandler handles HTTP requests for historical scorers.
type ScorerHandler struct {
	service service.ScorerService
}

// NewScorerHandler creates a new handler with the given service.
func NewScorerHandler(svc service.ScorerService) *ScorerHandler {
	return &ScorerHandler{service: svc}
}

// RegisterRoutes registers all scorer routes on the given router group.
func (h *ScorerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	scorers := rg.Group("/scorers")
	{
		scorers.GET("", h.List)
	}
}

// List godoc
// @Summary List historical scorers with filters and pagination
// @Produce json
// @Router /api/scorers [get]
func (h *ScorerHandler) List(c *gin.Context) {
	filter, err := parseScorerFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scorers, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve scorers"})
		return
	}

	c.JSON(http.StatusOK, scorers)
}

func parseScorerFilter(c *gin.Context) (domain.ScorerFilter, error) {
	filter := domain.ScorerFilter{
		Name:              c.Query("name"),
		TeamCode:          c.Query("teamCode"),
		ConfederationCode: c.Query("confederationCode"),
		Page:              defaultPage,
		Size:              defaultSize,
	}

	if rawPage := c.Query("page"); rawPage != "" {
		page, err := strconv.Atoi(rawPage)
		if err != nil {
			return domain.ScorerFilter{}, errors.New("invalid page parameter")
		}
		filter.Page = page
	}

	if rawSize := c.Query("size"); rawSize != "" {
		size, err := strconv.Atoi(rawSize)
		if err != nil {
			return domain.ScorerFilter{}, errors.New("invalid size parameter")
		}
		filter.Size = size
	}

	if filter.Page < 1 {
		return domain.ScorerFilter{}, errors.New("invalid page parameter")
	}
	if filter.Size < 1 || filter.Size > maxSize {
		return domain.ScorerFilter{}, errors.New("invalid size parameter")
	}

	return filter, nil
}
