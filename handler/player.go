package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"footballteam/helper"
	"footballteam/player"
)

type playerHandler struct {
	playerService player.Service
}

func NewPlayerHandler(playerService player.Service) *playerHandler {
	return &playerHandler{playerService}
}

func (h *playerHandler) GetPlayers(c *gin.Context) {
	players, err := h.playerService.GetAllPlayers()
	if err != nil {
		response := helper.APIResponse("Failed to get players", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("List of players", http.StatusOK, "success", player.FormatPlayers(players))
	c.JSON(http.StatusOK, response)
}

func (h *playerHandler) GetPlayerByID(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	p, err := h.playerService.GetPlayerByID(id)
	if err != nil {
		response := helper.APIResponse("Player not found", http.StatusNotFound, "error", err.Error())
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Player detail", http.StatusOK, "success", player.FormatPlayer(p))
	c.JSON(http.StatusOK, response)
}

func (h *playerHandler) GetPlayersByTeam(c *gin.Context) {
	teamIDParam := c.Param("team_id")
	teamID, _ := strconv.Atoi(teamIDParam)

	players, err := h.playerService.GetPlayersByTeam(teamID)
	if err != nil {
		response := helper.APIResponse("Failed to get players by team", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Players by team", http.StatusOK, "success", player.FormatPlayers(players))
	c.JSON(http.StatusOK, response)
}

func (h *playerHandler) CreatePlayer(c *gin.Context) {
	var input player.CreatePlayerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response := helper.APIResponse("Invalid input", http.StatusBadRequest, "error", helper.FormatValidationError(err))
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newPlayer, err := h.playerService.CreatePlayer(input)
	if err != nil {
		response := helper.APIResponse("Failed to create player", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Player created successfully", http.StatusCreated, "success", player.FormatPlayer(newPlayer))
	c.JSON(http.StatusCreated, response)
}

func (h *playerHandler) UpdatePlayer(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var input player.UpdatePlayerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response := helper.APIResponse("Invalid input", http.StatusBadRequest, "error", helper.FormatValidationError(err))
		c.JSON(http.StatusBadRequest, response)
		return
	}

	updatedPlayer, err := h.playerService.UpdatePlayer(id, input)
	if err != nil {
		response := helper.APIResponse("Failed to update player", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Player updated successfully", http.StatusOK, "success", player.FormatPlayer(updatedPlayer))
	c.JSON(http.StatusOK, response)
}

func (h *playerHandler) DeletePlayer(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	err := h.playerService.DeletePlayer(id)
	if err != nil {
		response := helper.APIResponse("Failed to delete player", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Player deleted successfully", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
