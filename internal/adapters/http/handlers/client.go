package handlers

import (
	"net/http"
	"strconv"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
)

// GetProjectsByClient godoc
// @Summary Retrieves all projects for a specific client
// @Description Fetches a list of projects associated with a given client
// @Tags Projects
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {array} dto.ProjectResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/clients/{id}/projects [get]
func GetProjectsByClient(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
			return
		}

		projects, err := projectService.GetByClient(
			c.Request.Context(), clientID,
		)
		if err != nil {
			c.JSON(
				utils.GetDomainErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		c.JSON(http.StatusOK, projects)
	}
}

// CreateClient godoc
// @Summary Creates a new client
// @Description Adds a new client to the system
// @Tags Clients
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.ClientCreate true "Client creation data"
// @Success 201 {object} dto.ClientResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/clients [post]
func CreateClient(clientService inbound.ClientHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ClientCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(
				utils.GetBindJSONErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		client, err := clientService.Create(c.Request.Context(), &req)
		if err != nil {
			c.JSON(
				utils.GetDomainErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		c.JSON(http.StatusCreated, client)
	}
}

// GetAllClients godoc
// @Summary Retrieves all clients with optional filtering by type
// @Description Fetches a list of clients, optionally filtered by client type
// @Tags Clients
// @Security OAuth2Password
// @Produce json
// @Param query query dto.ClientQueryParams false "Query parameters"
// @Success 200 {array} dto.ClientResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/clients [get]
func GetAllClients(clientService inbound.ClientHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ClientFilter

		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(
				utils.GetBindJSONErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		clients, err := clientService.GetClients(c.Request.Context(), &req)
		if err != nil {
			c.JSON(
				utils.GetDomainErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		c.JSON(http.StatusOK, clients)
	}
}
