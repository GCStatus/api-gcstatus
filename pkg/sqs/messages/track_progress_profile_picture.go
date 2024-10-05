package messages

import (
	"context"
	"encoding/json"
	"gcstatus/internal/domain"
	"gcstatus/internal/usecases"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type TrackProgressProfilePictureHandler struct {
	taskService *usecases.TaskService
}

func NewTrackProgressProfilePictureHandler(
	taskService *usecases.TaskService,
) *TrackProgressProfilePictureHandler {
	return &TrackProgressProfilePictureHandler{
		taskService: taskService,
	}
}

func (h *TrackProgressProfilePictureHandler) HandleTrackProgressProfilePictureMessage(ctx context.Context, message types.Message) {
	var messageWrapper struct {
		Type string          `json:"type"`
		Body json.RawMessage `json:"body"`
	}

	if err := json.Unmarshal([]byte(*message.Body), &messageWrapper); err != nil {
		log.Printf("Error unmarshalling main message wrapper: %v", err)
		return
	}

	var trackMsg struct {
		UserID    uint `json:"user_id"`
		Increment uint `json:"increment"`
	}

	if err := json.Unmarshal(messageWrapper.Body, &trackMsg); err != nil {
		log.Printf("Error unmarshalling track profile picture body: %v", err)
		return
	}

	if err := h.taskService.TrackTitleProgress(trackMsg.UserID, domain.ProfilePictureTitleRequirementKey, 1); err != nil {
		log.Fatalf("failed to track progress for user on title %+v. Error: %+s", trackMsg.UserID, err.Error())
	}

	if err := h.taskService.TrackMissionProgress(trackMsg.UserID, domain.ProfilePictureTitleRequirementKey, 1); err != nil {
		log.Fatalf("failed to track progress for user on mission %+v. Error: %+s", trackMsg.UserID, err.Error())
	}
}
