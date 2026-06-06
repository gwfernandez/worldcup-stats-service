package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

// ChampionshipHandler handles HTTP requests for championships.
type ChampionshipHandler struct {
	service service.ChampionshipService
}

// NewChampionshipHandler creates a new ChampionshipHandler.
func NewChampionshipHandler(svc service.ChampionshipService) *ChampionshipHandler {
	return &ChampionshipHandler{service: svc}
}

// RegisterRoutes registers all championship routes on the given router group.
func (h *ChampionshipHandler) RegisterRoutes(rg *gin.RouterGroup) {
	championships := rg.Group("/championships")
	{
		championships.GET("", h.List)
		championships.GET("/:year", h.GetByYear)
		championships.GET("/:year/teams", h.ListTeamsByYear)
	}
}

// List godoc
// @Summary List championships with filters and pagination
// @Produce json
// @Router /api/championships [get]
func (h *ChampionshipHandler) List(c *gin.Context) {
	filter, err := parseChampionshipFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve championships"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByYear godoc
// @Summary Get championship details by year
// @Produce json
// @Param year path int true "Championship Year"
// @Router /api/championships/{year} [get]
func (h *ChampionshipHandler) GetByYear(c *gin.Context) {
	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year parameter"})
		return
	}

	championship, err := h.service.GetByYear(c.Request.Context(), year)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve championship details"})
		return
	}

	c.JSON(http.StatusOK, championship)
}

// ListTeamsByYear godoc
// @Summary List teams that participated in a championship year with filters and pagination
// @Produce json
// @Param year path int true "Championship Year"
// @Router /api/championships/{year}/teams [get]
func (h *ChampionshipHandler) ListTeamsByYear(c *gin.Context) {
	filter, err := parseChampionshipTeamFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.ListTeamsByYear(c.Request.Context(), filter)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve championship teams"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func parseChampionshipFilter(c *gin.Context) (domain.ChampionshipFilter, error) {
	filter := domain.ChampionshipFilter{
		Host:              c.Query("host"),
		ConfederationCode: c.Query("confederation_code"),
		Page:              defaultPage,
		Size:              defaultSize,
	}

	if rawPage := c.Query("page"); rawPage != "" {
		page, err := strconv.Atoi(rawPage)
		if err != nil {
			return domain.ChampionshipFilter{}, errors.New("invalid page parameter")
		}
		filter.Page = page
	}

	if rawSize := c.Query("size"); rawSize != "" {
		size, err := strconv.Atoi(rawSize)
		if err != nil {
			return domain.ChampionshipFilter{}, errors.New("invalid size parameter")
		}
		filter.Size = size
	}

	if filter.Page < 1 {
		return domain.ChampionshipFilter{}, errors.New("invalid page parameter")
	}
	if filter.Size < 1 || filter.Size > maxSize {
		return domain.ChampionshipFilter{}, errors.New("invalid size parameter")
	}

	if rawYear := c.Query("year"); rawYear != "" {
		year, err := strconv.Atoi(rawYear)
		if err != nil {
			return domain.ChampionshipFilter{}, errors.New("invalid year parameter")
		}
		filter.Year = year
	}

	return filter, nil
}

func parseChampionshipTeamFilter(c *gin.Context) (domain.ChampionshipTeamFilter, error) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return domain.ChampionshipTeamFilter{}, errors.New("invalid year parameter")
	}

	filter := domain.ChampionshipTeamFilter{
		Year:              year,
		Name:              c.Query("name"),
		ConfederationCode: strings.ToUpper(c.Query("confederation_code")),
		GroupCode:         strings.ToUpper(c.Query("group_code")),
		Page:              defaultPage,
		Size:              defaultSize,
	}

	if rawPage := c.Query("page"); rawPage != "" {
		page, err := strconv.Atoi(rawPage)
		if err != nil {
			return domain.ChampionshipTeamFilter{}, errors.New("invalid page parameter")
		}
		filter.Page = page
	}

	if rawSize := c.Query("size"); rawSize != "" {
		size, err := strconv.Atoi(rawSize)
		if err != nil {
			return domain.ChampionshipTeamFilter{}, errors.New("invalid size parameter")
		}
		filter.Size = size
	}

	if filter.Page < 1 {
		return domain.ChampionshipTeamFilter{}, errors.New("invalid page parameter")
	}
	if filter.Size < 1 || filter.Size > maxSize {
		return domain.ChampionshipTeamFilter{}, errors.New("invalid size parameter")
	}

	return filter, nil
}
