package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/cache"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type MissionCompleteMessageHandler struct {
	walletService       *usecases.WalletService
	userService         *usecases.UserService
	taskService         *usecases.TaskService
	missionService      *usecases.MissionService
	transactionService  *usecases.TransactionService
	notificationService *usecases.NotificationService
}

func NewMissionCompleteMessageHandler(
	walletService *usecases.WalletService,
	userService *usecases.UserService,
	taskService *usecases.TaskService,
	missionService *usecases.MissionService,
	transactionService *usecases.TransactionService,
	notificationService *usecases.NotificationService,
) *MissionCompleteMessageHandler {
	return &MissionCompleteMessageHandler{
		walletService:       walletService,
		userService:         userService,
		taskService:         taskService,
		missionService:      missionService,
		transactionService:  transactionService,
		notificationService: notificationService,
	}
}

func (h *MissionCompleteMessageHandler) HandleCompleteMissionMessage(ctx context.Context, message types.Message) {
	var messageWrapper struct {
		Type string          `json:"type"`
		Body json.RawMessage `json:"body"`
	}

	if err := json.Unmarshal([]byte(*message.Body), &messageWrapper); err != nil {
		log.Printf("Error unmarshalling main message wrapper: %v", err)
		return
	}

	var completeMissionMsg struct {
		UserID    uint `json:"user_id"`
		MissionID uint `json:"mission_id"`
	}

	if err := json.Unmarshal(messageWrapper.Body, &completeMissionMsg); err != nil {
		log.Printf("Error unmarshalling complete message body: %v", err)
		return
	}

	mission, err := h.missionService.FindByID(completeMissionMsg.MissionID)
	if err != nil {
		log.Fatalf("mission not found: %+v", err)
	}

	for _, reward := range mission.Rewards {
		if reward.RewardableType == "titles" {
			if err := h.taskService.AwardTitleToUser(completeMissionMsg.UserID, reward.RewardableID); err != nil {
				log.Fatalf("error awarding title: %+v", err)
			}

			h.createRewardNotification(*mission, completeMissionMsg.UserID)
		}
	}

	if err := h.walletService.Add(completeMissionMsg.UserID, mission.Coins); err != nil {
		log.Fatalf("failed to add coins to user wallet %+v. Error: %+s", completeMissionMsg.UserID, err.Error())
	}

	h.createTransactionForMissionCoins(*mission, completeMissionMsg.UserID)

	h.createMissionCompleteNotification(*mission, completeMissionMsg.UserID)

	if err := h.userService.AddExperience(completeMissionMsg.UserID, mission.Experience, h.taskService.AwardTitleToUser); err != nil {
		log.Fatalf("failed to experience to user %+v. Error: %+s", completeMissionMsg.UserID, err.Error())
	}

	cache.GlobalCache.RemoveUserFromCache(completeMissionMsg.UserID)
}

func (h *MissionCompleteMessageHandler) createTransactionForMissionCoins(mission domain.Mission, userID uint) {
	transaction := &domain.Transaction{
		Amount:            mission.Coins,
		Description:       fmt.Sprintf("Received coins from mission %s.", mission.Mission),
		UserID:            userID,
		TransactionTypeID: domain.AdditionTransactionTypeID,
	}

	if err := h.transactionService.CreateTransaction(transaction); err != nil {
		log.Printf("Failed to create a transaction for user title purchase: %+v", err)
	}
}

func (h *MissionCompleteMessageHandler) createRewardNotification(mission domain.Mission, userID uint) {
	notificationContent := &domain.NotificationData{
		Title:     fmt.Sprintf("You have achieved a new title from mission: %s", mission.Mission),
		ActionUrl: "/profile/?section=missions",
		Icon:      "FaAward",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Printf("Failed to marshal notification content: %+v", err)
	}

	notification := &domain.Notification{
		Type:   "NewTitleAchieved",
		Data:   string(dataJson),
		UserID: userID,
	}

	if err := h.notificationService.CreateNotification(notification); err != nil {
		log.Printf("Failed to save the reward title notification: %+v", err)
	}
}

func (h *MissionCompleteMessageHandler) createMissionCompleteNotification(mission domain.Mission, userID uint) {

	notificationContent := &domain.NotificationData{
		Title:     fmt.Sprintf("You have completed the mission: %s", mission.Mission),
		ActionUrl: "/profile/?section=missions",
		Icon:      "FaAward",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Printf("Failed to marshal notification content: %+v", err)
	}

	notification := &domain.Notification{
		Type:   "NewCompleteMission",
		Data:   string(dataJson),
		UserID: userID,
	}

	if err := h.notificationService.CreateNotification(notification); err != nil {
		log.Printf("Failed to save the complete mission notification: %+v", err)
	}
}
