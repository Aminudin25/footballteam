package handler

import (
	"fmt"
	"footballteam/helper"
	"footballteam/team"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type teamHandler struct {
	teamService team.Service
}

func NewTeamHandler(teamService team.Service) *teamHandler {
	return &teamHandler{teamService}
}

// GET /api/teams
func (h *teamHandler) GetTeams(c *gin.Context) {
	teams, err := h.teamService.GetAllTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to get teams", http.StatusInternalServerError, "error", nil))
		return
	}

	var formattedTeams []team.TeamFormatter
	for _, t := range teams {
		formattedTeams = append(formattedTeams, team.FormatTeam(t))
	}

	response := helper.APIResponse("List of teams", http.StatusOK, "success", formattedTeams)
	c.JSON(http.StatusOK, response)
}

// GET /api/teams/:id
func (h *teamHandler) GetTeamByID(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	t, err := h.teamService.GetTeamByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.APIResponse("Team not found", http.StatusNotFound, "error", nil))
		return
	}

	response := helper.APIResponse("Team detail", http.StatusOK, "success", team.FormatTeam(t))
	c.JSON(http.StatusOK, response)
}

// POST /api/teams
func (h *teamHandler) CreateTeam(c *gin.Context) {
	var input team.CreateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, helper.APIResponse("Create team failed", http.StatusUnprocessableEntity, "error", errors))
		return
	}

	newTeam, err := h.teamService.CreateTeam(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIResponse("Create team failed", http.StatusInternalServerError, "error", nil))
		return
	}

	response := helper.APIResponse("Team created", http.StatusOK, "success", team.FormatTeam(newTeam))
	c.JSON(http.StatusOK, response)
}

// PUT /api/teams/:id
func (h *teamHandler) UpdateTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	
	if err != nil {
    	response := helper.APIResponse("Invalid team ID", http.StatusBadRequest, "error", nil)
    	c.JSON(http.StatusBadRequest, response)
    	return
	}

	var input team.UpdateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, helper.APIResponse("Update team failed", http.StatusUnprocessableEntity, "error", errors))
		return
	}

	updatedTeam, err := h.teamService.UpdateTeam(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIResponse("Update team failed", http.StatusInternalServerError, "error", nil))
		return
	}

	response := helper.APIResponse("Team updated", http.StatusOK, "success", team.FormatTeam(updatedTeam))
	c.JSON(http.StatusOK, response)
}

// DELETE /api/teams/:id
func (h *teamHandler) DeleteTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	err := h.teamService.DeleteTeam(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIResponse("Delete team failed", http.StatusInternalServerError, "error", nil))
		return
	}

	response := helper.APIResponse("Team deleted", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *teamHandler) UploadLogo(c *gin.Context) {
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIResponse("Logo upload failed", http.StatusBadRequest, "error", err.Error()))
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, helper.APIResponse("Only JPG and PNG files are allowed", http.StatusBadRequest, "error", nil))
		return
	}

	uploadDir := "uploads/logo"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to create upload directory", http.StatusInternalServerError, "error", err.Error()))
			return
		}
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	uploadPath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to save logo", http.StatusInternalServerError, "error", err.Error()))
		return
	}

	idParam := c.Param("id")
	teamID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIResponse("Invalid team ID", http.StatusBadRequest, "error", err.Error()))
		return
	}

	_, err = h.teamService.SaveLogo(teamID, uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.APIResponse("Failed to save logo to database", http.StatusInternalServerError, "error", err.Error()))
		return
	}

	response := map[string]string{
		"logo_url": "/" + uploadPath,
	}
	c.JSON(http.StatusOK, helper.APIResponse("Logo uploaded successfully", http.StatusOK, "success", response))
}
