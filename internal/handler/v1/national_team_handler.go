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

// NationalTeamHandler handles HTTP requests for national teams.
type NationalTeamHandler struct {
	service service.NationalTeamService
}

// NewNationalTeamHandler creates a new handler with the given service.
func NewNationalTeamHandler(svc service.NationalTeamService) *NationalTeamHandler {
	return &NationalTeamHandler{service: svc}
}

// RegisterRoutes registers all national team routes on the given router group.
func (h *NationalTeamHandler) RegisterRoutes(rg *gin.RouterGroup) {
	nationalTeams := rg.Group("/national-teams")
	{
		nationalTeams.GET("", h.List)
		nationalTeams.GET("/:code", h.GetByCode)
	}
}

// List godoc
// @Summary List national teams with filters and pagination
// @Produce json
// @Router /api/national-teams [get]
func (h *NationalTeamHandler) List(c *gin.Context) {
	filter, err := parseNationalTeamFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teams, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve national teams"})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// GetByCode godoc
// @Summary Get a national team by code
// @Produce json
// @Param code path string true "National Team code"
// @Router /api/national-teams/{code} [get]
func (h *NationalTeamHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	team, err := h.service.GetByCode(c.Request.Context(), code)	
	if err != nil {
		if isNationalTeamNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve national team"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func parseNationalTeamFilter(c *gin.Context) (domain.NationalTeamFilter, error) {
	filter := domain.NationalTeamFilter{
		Name:           c.Query("name"),
		FederationName: c.Query("federation_name"),
		FederationCode: c.Query("federation_code"),
		Page:           defaultPage,
		Size:           defaultSize,
	}

	if rawPage := c.Query("page"); rawPage != "" {
		page, err := strconv.Atoi(rawPage)
		if err != nil {
			return domain.NationalTeamFilter{}, errors.New("invalid page parameter")
		}
		filter.Page = page
	}

	if rawSize := c.Query("size"); rawSize != "" {
		size, err := strconv.Atoi(rawSize)
		if err != nil {
			return domain.NationalTeamFilter{}, errors.New("invalid size parameter")
		}
		filter.Size = size
	}

	if filter.Page < 1 {
		return domain.NationalTeamFilter{}, errors.New("invalid page parameter")
	}
	if filter.Size < 1 || filter.Size > maxSize {
		return domain.NationalTeamFilter{}, errors.New("invalid size parameter")
	}

	if rawConfederationCode := c.Query("confederation_code"); rawConfederationCode != "" {
		filter.ConfederationCode = &rawConfederationCode
	}

	if rawIncludeDissolved := c.Query("include_dissolved"); rawIncludeDissolved != "" {
		includeDissolved, err := strconv.ParseBool(rawIncludeDissolved)
		if err != nil {
			return domain.NationalTeamFilter{}, errors.New("invalid include_dissolved parameter")
		}
		filter.IncludeDissolved = includeDissolved
	}

	return filter, nil
}

func isNationalTeamNotFoundError(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}
