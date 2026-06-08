package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

// StandingHandler handles HTTP requests for historical standings.
type StandingHandler struct {
	service service.StandingService
}

// NewStandingHandler creates a new handler with the given service.
func NewStandingHandler(svc service.StandingService) *StandingHandler {
	return &StandingHandler{service: svc}
}

// RegisterRoutes registers all standing routes on the given router group.
func (h *StandingHandler) RegisterRoutes(rg *gin.RouterGroup) {
	standings := rg.Group("/standings")
	{
		standings.GET("", h.List)
	}
}

// List godoc
// @Summary List historical standings with filters and pagination
// @Produce json
// @Router /api/standings [get]
func (h *StandingHandler) List(c *gin.Context) {
	filter, err := parseStandingFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	standings, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve standings"})
		return
	}

	c.JSON(http.StatusOK, standings)
}

func parseStandingFilter(c *gin.Context) (domain.StandingFilter, error) {
	filter := domain.StandingFilter{
		Name:              c.Query("name"),
		ConfederationCode: c.Query("confederationCode"),
		Page:              defaultPage,
		Size:              defaultSize,
	}

	if rawPage := c.Query("page"); rawPage != "" {
		page, err := strconv.Atoi(rawPage)
		if err != nil {
			return domain.StandingFilter{}, errors.New("invalid page parameter")
		}
		filter.Page = page
	}

	if rawSize := c.Query("size"); rawSize != "" {
		size, err := strconv.Atoi(rawSize)
		if err != nil {
			return domain.StandingFilter{}, errors.New("invalid size parameter")
		}
		filter.Size = size
	}

	if filter.Page < 1 {
		return domain.StandingFilter{}, errors.New("invalid page parameter")
	}
	if filter.Size < 1 || filter.Size > maxSize {
		return domain.StandingFilter{}, errors.New("invalid size parameter")
	}

	return filter, nil
}
