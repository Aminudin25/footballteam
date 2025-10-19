package handler

import (
	"net/http"
	"strconv"

	"footballteam/helper"
	"footballteam/match"

	"github.com/gin-gonic/gin"
)

type matchHandler struct {
	matchService match.Service
}

func NewMatchHandler(matchService match.Service) *matchHandler {
	return &matchHandler{matchService}
}

// GET /matches
func (h *matchHandler) GetMatches(c *gin.Context) {
	matches, err := h.matchService.FindAll()
	if err != nil {
		response := helper.APIResponse("Failed to get matches", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("List of matches", http.StatusOK, "success", match.FormatMatches(matches))
	c.JSON(http.StatusOK, response)
}

// GET /matches/:id
func (h *matchHandler) GetMatchByID(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	m, err := h.matchService.FindByID(id)
	if err != nil {
		response := helper.APIResponse("Match not found", http.StatusNotFound, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Match detail", http.StatusOK, "success", match.FormatMatch(m))
	c.JSON(http.StatusOK, response)
}

// POST /matches
func (h *matchHandler) CreateMatch(c *gin.Context) {
	var input match.CreateMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response := helper.APIResponse("Invalid input", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newMatch, err := h.matchService.CreateMatch(input)
	if err != nil {
		response := helper.APIResponse("Failed to create match", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Match created successfully", http.StatusOK, "success", match.FormatMatch(newMatch))
	c.JSON(http.StatusOK, response)
}

// PUT /matches/:id
func (h *matchHandler) UpdateMatch(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var input match.UpdateMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response := helper.APIResponse("Invalid input", http.StatusBadRequest, "error", helper.FormatValidationError(err))
		c.JSON(http.StatusBadRequest, response)
		return
	}

	updatedMatch, err := h.matchService.UpdateMatch(id, input)
	if err != nil {
		response := helper.APIResponse("Failed to update match", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Match updated successfully", http.StatusOK, "success", match.FormatMatch(updatedMatch))
	c.JSON(http.StatusOK, response)
}

// DELETE /matches/:id
func (h *matchHandler) DeleteMatch(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	err := h.matchService.DeleteMatch(id)
	if err != nil {
		response := helper.APIResponse("Failed to delete match", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Match deleted successfully", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
