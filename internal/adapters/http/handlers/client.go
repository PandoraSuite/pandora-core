package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/dto"
	"github.com/MAD-py/pandora-core/internal/adapters/http/errors"
	"github.com/MAD-py/pandora-core/internal/app/client"
)

// ClientList godoc
// @Summary Retrieves all clients with optional filtering by type
// @Description Fetches a list of clients, optionally filtered by client type
// @Tags Clients
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param query query dto.ClientFilter false "Query parameters"
// @Success 200 {array} dto.ClientResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/clients [get]
func ClientList(useCase client.ListUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.ClientFilter{Type: c.Query("type")}
		clients, err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		resp := make([]*dto.ClientResponse, len(clients))
		for i, client := range clients {
			resp[i] = dto.ClientResponseFromDomain(client)
		}
		c.JSON(http.StatusOK, resp)
	}
}

// ClientCreate godoc
// @Summary Creates a new client
// @Description Adds a new client to the system
// @Tags Clients
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.ClientCreate true "Client creation data"
// @Success 201 {object} dto.ClientResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/clients [post]
func ClientCreate(useCase client.CreateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ClientCreate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindJSONToHTTPError(req, err))
			return
		}

		client, err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, dto.ClientResponseFromDomain(client))
	}
}

// ClientGet godoc
// @Summary Retrieves a client by ID
// @Description Fetches the details of a specific client using its ID
// @Tags Clients
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} dto.ClientResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/clients/{id} [get]
func ClientGet(useCase client.GetUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid client id",
				),
			)
			return
		}

		client, err := useCase.Execute(c.Request.Context(), clientID)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.ClientResponseFromDomain(client))
	}
}

// ClientUpdate godoc
// @Summary Updates an existing client
// @Description Modifies client data based on the provided ID
// @Tags Clients
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Client ID"
// @Param request body dto.ClientUpdate true "Updated client data"
// @Success 200 {object} dto.ClientResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/clients/{id} [patch]
func ClientUpdate(useCase client.UpdateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid client id",
				),
			)
			return
		}

		var req dto.ClientUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindJSONToHTTPError(req, err))
			return
		}

		client, err := useCase.Execute(
			c.Request.Context(), clientID, req.ToDomain(),
		)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.ClientResponseFromDomain(client))
	}
}

// ClientDelete godoc
// @Summary Deletes a client by ID
// @Description Removes a specific client from the system using its ID
// @Tags Clients
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Success 204 "No Content"
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/clients/{id} [delete]
func ClientDelete(useCase client.DeleteUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid client id",
				),
			)
			return
		}

		if err := useCase.Execute(c.Request.Context(), clientID); err != nil {
			c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// ClientListProjects godoc
// @Summary Retrieves all projects for a specific client
// @Description Fetches a list of projects associated with a given client
// @Tags Clients
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {array} dto.ProjectResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/clients/{id}/projects [get]
func ClientListProjects(useCase client.ListProjectsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid client id",
				),
			)
			return
		}

		projects, err := useCase.Execute(c.Request.Context(), clientID)
		if err != nil {
			c.Error(err)
			return
		}

		resp := make([]*dto.ProjectResponse, len(projects))
		for i, project := range projects {
			resp[i] = dto.ProjectResponseFromDomain(project)
		}
		c.JSON(http.StatusOK, resp)
	}
}
