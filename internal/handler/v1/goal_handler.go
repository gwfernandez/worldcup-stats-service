package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

// GoalHandler handles HTTP requests for player goals.
type GoalHandler struct {
	service service.GoalService
}

// NewGoalHandler creates a new handler with the given service.
func NewGoalHandler(svc service.GoalService) *GoalHandler {
	return &GoalHandler{service: svc}
}

// RegisterRoutes registers all goal routes on the given router group.
func (h *GoalHandler) RegisterRoutes(rg *gin.RouterGroup) {
	players := rg.Group("/players")
	{
		players.GET("/:playerId/goals", h.List)
	}
}

// List godoc
// @Summary List goals scored by a player with optional championship filter and pagination
// @Produce json
// @Param playerId path int true "Player ID"
// @Param year query int false "World Cup year"
// @Router /api/players/{playerId}/goals [get]
func (h *GoalHandler) List(c *gin.Context) {
	filter, err := parseGoalFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goals, err := h.service.ListByPlayer(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve goals"})
		return
	}

	c.JSON(http.StatusOK, goals)
}

func parseGoalFilter(c *gin.Context) (domain.GoalFilter, error) {
	playerID, err := strconv.ParseInt(c.Param("playerId"), 10, 64)
	if err != nil || playerID < 1 {
		return domain.GoalFilter{}, errors.New("invalid playerId parameter")
	}

	filter := domain.GoalFilter{
		PlayerID: playerID,
		Language: resolveLanguage(c.Request),
		Page:     defaultPage,
		Size:     defaultSize,
	}

	if rawYear := c.Query("year"); rawYear != "" {
		year, err := strconv.Atoi(rawYear)
		if err != nil || year < 1 {
			return domain.GoalFilter{}, errors.New("invalid year parameter")
		}
		filter.Year = year
	}

	if rawPage := c.Query("page"); rawPage != "" {
		page, err := strconv.Atoi(rawPage)
		if err != nil || page < 1 {
			return domain.GoalFilter{}, errors.New("invalid page parameter")
		}
		filter.Page = page
	}

	if rawSize := c.Query("size"); rawSize != "" {
		size, err := strconv.Atoi(rawSize)
		if err != nil || size < 1 || size > maxSize {
			return domain.GoalFilter{}, errors.New("invalid size parameter")
		}
		filter.Size = size
	}

	return filter, nil
}
