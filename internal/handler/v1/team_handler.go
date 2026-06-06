package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
	"github.com/jendrix/worldcup-stats-service/internal/service"
)

const (
	defaultPage = 1
	defaultSize = 20
	maxSize     = 100
)

// TeamHandler handles HTTP requests for teams.
type TeamHandler struct {
	service service.TeamService
}

// NewTeamHandler creates a new handler with the given service.
func NewTeamHandler(svc service.TeamService) *TeamHandler {
	return &TeamHandler{service: svc}
}

// RegisterRoutes registers all team routes on the given router group.
func (h *TeamHandler) RegisterRoutes(rg *gin.RouterGroup) {
	teams := rg.Group("/teams")
	{
		teams.GET("", h.List)
		teams.GET("/:code", h.GetByCode)
	}
}

// List godoc
// @Summary List teams with filters and pagination
// @Produce json
// @Router /api/teams [get]
func (h *TeamHandler) List(c *gin.Context) {
	filter, err := parseTeamFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teams, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve teams"})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// GetByCode godoc
// @Summary Get a team by code
// @Produce json
// @Param code path string true "Team code"
// @Router /api/teams/{code} [get]
func (h *TeamHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	team, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		if isTeamNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve team"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func parseTeamFilter(c *gin.Context) (domain.TeamFilter, error) {
	filter := domain.TeamFilter{
		Name:           c.Query("name"),
		FederationName: c.Query("federationName"),
		FederationCode: c.Query("federationCode"),
		Page:           defaultPage,
		Size:           defaultSize,
	}

	if rawPage := c.Query("page"); rawPage != "" {
		page, err := strconv.Atoi(rawPage)
		if err != nil {
			return domain.TeamFilter{}, errors.New("invalid page parameter")
		}
		filter.Page = page
	}

	if rawSize := c.Query("size"); rawSize != "" {
		size, err := strconv.Atoi(rawSize)
		if err != nil {
			return domain.TeamFilter{}, errors.New("invalid size parameter")
		}
		filter.Size = size
	}

	if filter.Page < 1 {
		return domain.TeamFilter{}, errors.New("invalid page parameter")
	}
	if filter.Size < 1 || filter.Size > maxSize {
		return domain.TeamFilter{}, errors.New("invalid size parameter")
	}

	if rawConfederationCode := c.Query("confederationCode"); rawConfederationCode != "" {
		filter.ConfederationCode = &rawConfederationCode
	}

	if rawIncludeDissolved := c.Query("includeDissolved"); rawIncludeDissolved != "" {
		includeDissolved, err := strconv.ParseBool(rawIncludeDissolved)
		if err != nil {
			return domain.TeamFilter{}, errors.New("invalid includeDissolved parameter")
		}
		filter.IncludeDissolved = includeDissolved
	}

	return filter, nil
}

func isTeamNotFoundError(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}
