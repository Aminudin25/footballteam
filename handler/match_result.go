package handler

import (
	"footballteam/helper"
	"footballteam/match_result"
	"footballteam/player"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type matchResultHandler struct {
	service       match_result.Service
	playerService player.Service
}

func NewMatchResultHandler(s match_result.Service, p player.Service) *matchResultHandler {
	return &matchResultHandler{
		service:       s,
		playerService: p,
	}
}

// POST /match_results
func (h *matchResultHandler) CreateMatchResult(c *gin.Context) {
	var input match_result.CreateMatchResultInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response := helper.APIResponse("Failed to create match result", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := h.service.Create(input)
	if err != nil {
		response := helper.APIResponse("Failed to create match result", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Ambil nama player untuk goals
	playerMap := make(map[int]string)
	for _, g := range result.Goals {
		if _, exists := playerMap[g.PlayerID]; !exists {
			p, err := h.playerService.GetPlayerByID(g.PlayerID)
			if err == nil {
				playerMap[g.PlayerID] = p.Name
			} else {
				playerMap[g.PlayerID] = "Unknown"
			}
		}
	}

	formatted := match_result.FormatMatchResult(result, playerMap)
	response := helper.APIResponse("Match result created", http.StatusOK, "success", formatted)
	c.JSON(http.StatusOK, response)
}

// GET /match_results
func (h *matchResultHandler) GetMatchResults(c *gin.Context) {
	results, err := h.service.FindAll()
	if err != nil {
		response := helper.APIResponse("Failed to get match results", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var formattedList []match_result.MatchResultFormatter
	playerMap := make(map[int]string)

	for _, r := range results {
		for _, g := range r.Goals {
			if _, exists := playerMap[g.PlayerID]; !exists {
				p, err := h.playerService.GetPlayerByID(g.PlayerID)
				if err == nil {
					playerMap[g.PlayerID] = p.Name
				} else {
					playerMap[g.PlayerID] = "Unknown"
				}
			}
		}
		formattedList = append(formattedList, match_result.FormatMatchResult(r, playerMap))
	}

	response := helper.APIResponse("List of match results", http.StatusOK, "success", formattedList)
	c.JSON(http.StatusOK, response)
}


// GET /match_results/:id
func (h *matchResultHandler) GetMatchResultByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam) // convert string ke int
	if err != nil {
		response := helper.APIResponse("Invalid match result ID", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := h.service.FindByID(id)
	if err != nil {
		response := helper.APIResponse("Match result not found", http.StatusNotFound, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	playerMap := make(map[int]string)
	for _, g := range result.Goals {
		if _, exists := playerMap[g.PlayerID]; !exists {
			p, err := h.playerService.GetPlayerByID(g.PlayerID)
			if err == nil {
				playerMap[g.PlayerID] = p.Name
			} else {
				playerMap[g.PlayerID] = "Unknown"
			}
		}
	}

	formatted := match_result.FormatMatchResult(result, playerMap)
	response := helper.APIResponse("Match result detail", http.StatusOK, "success", formatted)
	c.JSON(http.StatusOK, response)
}

func (h *matchResultHandler) GetMatchResultsReport(c *gin.Context) {
	report, err := h.service.GetMatchResultsReport()
	if err != nil {
		response := helper.APIResponse("Failed to get match results report", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Match results report", http.StatusOK, "success", report)
	c.JSON(http.StatusOK, response)
}


