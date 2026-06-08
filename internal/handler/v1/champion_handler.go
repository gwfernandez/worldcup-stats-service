package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

// ChampionHandler handles HTTP requests for champions.
type ChampionHandler struct {
	service service.ChampionService
}

// NewChampionHandler creates a new handler with the given service.
func NewChampionHandler(svc service.ChampionService) *ChampionHandler {
	return &ChampionHandler{service: svc}
}

// RegisterRoutes registers all champion routes on the given router group.
func (h *ChampionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	champions := rg.Group("/champions")
	{
		champions.GET("", h.List)
	}
}

// List godoc
// @Summary List champions with pagination
// @Produce json
// @Router /api/champions [get]
func (h *ChampionHandler) List(c *gin.Context) {
	filter, err := parseChampionFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	champions, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve champions"})
		return
	}

	c.JSON(http.StatusOK, champions)
}

func parseChampionFilter(c *gin.Context) (domain.ChampionFilter, error) {
	filter := domain.ChampionFilter{
		Page: defaultPage,
		Size: defaultSize,
	}

	if rawPage := c.Query("page"); rawPage != "" {
		page, err := strconv.Atoi(rawPage)
		if err != nil {
			return domain.ChampionFilter{}, errors.New("invalid page parameter")
		}
		filter.Page = page
	}

	if rawSize := c.Query("size"); rawSize != "" {
		size, err := strconv.Atoi(rawSize)
		if err != nil {
			return domain.ChampionFilter{}, errors.New("invalid size parameter")
		}
		filter.Size = size
	}

	if filter.Page < 1 {
		return domain.ChampionFilter{}, errors.New("invalid page parameter")
	}
	if filter.Size < 1 || filter.Size > maxSize {
		return domain.ChampionFilter{}, errors.New("invalid size parameter")
	}

	return filter, nil
}
