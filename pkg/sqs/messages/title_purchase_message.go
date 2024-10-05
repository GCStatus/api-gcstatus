package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/ses"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type PurchaseMessageHandler struct {
	userService         *usecases.UserService
	transactionService  *usecases.TransactionService
	notificationService *usecases.NotificationService
}

func NewPurchaseMessageHandler(
	userService *usecases.UserService,
	transactionService *usecases.TransactionService,
	notificationService *usecases.NotificationService,
) *PurchaseMessageHandler {
	return &PurchaseMessageHandler{
		userService:         userService,
		transactionService:  transactionService,
		notificationService: notificationService,
	}
}

func (h *PurchaseMessageHandler) HandlePurchaseMessage(ctx context.Context, message types.Message) {
	var messageWrapper struct {
		Type string          `json:"type"`
		Body json.RawMessage `json:"body"`
	}

	if err := json.Unmarshal([]byte(*message.Body), &messageWrapper); err != nil {
		log.Printf("Error unmarshalling main message wrapper: %v", err)
		return
	}

	var purchaseMsg struct {
		UserID  uint   `json:"user_id"`
		TitleID uint   `json:"title_id"`
		Cost    uint   `json:"cost"`
		Title   string `json:"title"`
	}

	if err := json.Unmarshal(messageWrapper.Body, &purchaseMsg); err != nil {
		log.Printf("Error unmarshalling purchase body: %v", err)
		return
	}

	transaction := &domain.Transaction{
		Amount:            purchaseMsg.Cost,
		Description:       fmt.Sprintf("Purchase of title %s by %d coins.", purchaseMsg.Title, purchaseMsg.Cost),
		UserID:            purchaseMsg.UserID,
		TransactionTypeID: domain.SubtractionTransactionTypeID,
	}

	if err := h.transactionService.CreateTransaction(transaction); err != nil {
		log.Printf("Failed to create a transaction for user title purchase: %+v", err)
	}

	notificationContent := &domain.NotificationData{
		Title:     fmt.Sprintf("You bought %s title by %d coins!", purchaseMsg.Title, purchaseMsg.Cost),
		ActionUrl: "/profile/?section=transactions",
		Icon:      "CiCoinInsert",
	}

	dataJson, err := json.Marshal(notificationContent)
	if err != nil {
		log.Printf("Failed to marshal notification content: %+v", err)
	}

	notification := &domain.Notification{
		Type:   "NewTitlePurchase",
		Data:   string(dataJson),
		UserID: purchaseMsg.UserID,
	}

	if err := h.notificationService.CreateNotification(notification); err != nil {
		log.Printf("Failed to save the title purchase notification: %+v", err)
	}

	user, err := h.userService.GetUserByID(purchaseMsg.UserID)
	if err != nil {
		log.Printf("Failed to get user by id: %+v", err)
		return
	}

	if err := ses.SendTransactionEmail(user, transaction, ses.Send); err != nil {
		log.Printf("Failed to send transaction email: %+v", err)
		return
	}
}
