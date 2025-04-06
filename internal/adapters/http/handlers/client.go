package handlers

import (
	"net/http"
	"strconv"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
)

// GetProjectsByClient godoc
// @Summary Retrieves all projects for a specific client
// @Description Fetches a list of projects associated with a given client
// @Tags Clients
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {array} dto.ProjectResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/clients/{id}/projects [get]
func GetProjectsByClient(clientService inbound.ClientHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid client ID"},
			)
			return
		}

		projects, err := clientService.GetProjects(
			c.Request.Context(), clientID,
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
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
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/clients [post]
func CreateClient(clientService inbound.ClientHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ClientCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		client, err := clientService.Create(c.Request.Context(), &req)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
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
// @Param query query dto.ClientFilter false "Query parameters"
// @Success 200 {array} dto.ClientResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/clients [get]
func GetAllClients(clientService inbound.ClientHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, paramErr := enums.ParseClientType(c.Query("type"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusUnprocessableEntity,
				utils.ErrorResponse{Error: paramErr},
			)
			return
		}

		req := dto.ClientFilter{Type: t}
		clients, err := clientService.GetAll(c.Request.Context(), &req)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, clients)
	}
}

// UpdateClient godoc
// @Summary Updates an existing client
// @Description Modifies client data based on the provided ID
// @Tags Clients
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Client ID"
// @Param request body dto.ClientUpdate true "Updated client data"
// @Success 204
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/clients/{id} [patch]
func UpdateClient(clientService inbound.ClientHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid client ID"},
			)
			return
		}

		var req dto.ClientUpdate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		err := clientService.Update(c.Request.Context(), clientID, &req)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// GetClient godoc
// @Summary Retrieves a client by ID
// @Description Fetches the details of a specific client using its ID
// @Tags Clients
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} dto.ClientResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/clients/{id} [get]
func GetClient(clientService inbound.ClientHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid client ID"},
			)
			return
		}

		clients, err := clientService.GetByID(c.Request.Context(), clientID)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, clients)
	}
}
