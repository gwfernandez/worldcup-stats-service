package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

// FixtureHandler handles HTTP requests for championship fixtures.
type FixtureHandler struct {
	service service.FixtureService
}

// NewFixtureHandler creates a new FixtureHandler.
func NewFixtureHandler(svc service.FixtureService) *FixtureHandler {
	return &FixtureHandler{service: svc}
}

// RegisterRoutes registers all fixture routes on the given router group.
func (h *FixtureHandler) RegisterRoutes(rg *gin.RouterGroup) {
	championships := rg.Group("/championships")
	{
		championships.GET("/:year/fixture", h.GetByYear)
	}
}

// GetByYear godoc
// @Summary Get championship fixture by year
// @Produce json
// @Param year path int true "Championship Year"
// @Router /api/championships/{year}/fixture [get]
func (h *FixtureHandler) GetByYear(c *gin.Context) {
	year, err := parseFixtureYear(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	fixture, err := h.service.GetByYear(c.Request.Context(), year)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "championship not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve fixture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fixture})
}

func parseFixtureYear(rawYear string) (int, error) {
	year, err := strconv.Atoi(rawYear)
	if err != nil || year <= 0 {
		return 0, errors.New("invalid year")
	}
	return year, nil
}
