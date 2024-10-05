package api

import (
	"encoding/json"
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/cache"
	"gcstatus/pkg/sqs"
	"gcstatus/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MissionHandler struct {
	missionService *usecases.MissionService
	userService    *usecases.UserService
}

func NewMissionHandler(
	missionService *usecases.MissionService,
	userService *usecases.UserService,
) *MissionHandler {
	return &MissionHandler{
		missionService: missionService,
		userService:    userService,
	}
}

func (h *MissionHandler) GetAllForUser(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	missions, err := h.missionService.GetAllForUser(user.ID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user missions: "+err.Error())
		return
	}

	var transformedMissions []resources.MissionResource

	if len(missions) > 0 {
		transformedMissions = resources.TransformMissions(missions)
	} else {
		transformedMissions = []resources.MissionResource{}
	}

	response := resources.Response{
		Data: transformedMissions,
	}

	c.JSON(http.StatusOK, response)
}

func (h *MissionHandler) CompleteMission(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	missionIDStr := c.Param("id")

	missionID, err := strconv.ParseUint(missionIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid mission ID: "+err.Error())
		return
	}

	mission, err := h.missionService.FindByID(uint(missionID))
	if err != nil {
		RespondWithError(c, http.StatusNotFound, "Mission not found!")
		return
	}

	trackProgressMessage := map[string]any{
		"type": "CompleteMission",
		"body": map[string]any{
			"user_id":    user.ID,
			"mission_id": mission.ID,
		},
	}

	missionCompleteMessage, err := json.Marshal(trackProgressMessage)
	if err != nil {
		log.Fatalf("failed to serialize mission complete message to JSON: %+v", err)
	}

	if err := sqs.GlobalSQSClient.SendMessage(c.Request.Context(), sqs.GetAwsQueue(), string(missionCompleteMessage)); err != nil {
		log.Fatalf("failed to enqueue complete mission message to SQS: %+v", err)
	}

	if err := h.missionService.CompleteMission(user.ID, uint(missionID)); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to complete the mission: "+err.Error())
		return
	}

	cache.GlobalCache.RemoveUserFromCache(user.ID)

	c.JSON(http.StatusOK, gin.H{"message": "You have successfully completed the mission!"})
}
