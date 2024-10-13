package api

import (
	"encoding/json"
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"gcstatus/pkg/cache"
	"gcstatus/pkg/sqs"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TitleHandler struct {
	titleService        *usecases.TitleService
	userService         *usecases.UserService
	walletService       *usecases.WalletService
	taskService         *usecases.TaskService
	transactionService  *usecases.TransactionService
	notificationService *usecases.NotificationService
}

func NewTitleHandler(
	titleService *usecases.TitleService,
	userService *usecases.UserService,
	walletService *usecases.WalletService,
	taskService *usecases.TaskService,
	transactionService *usecases.TransactionService,
	notificationService *usecases.NotificationService,
) *TitleHandler {
	return &TitleHandler{
		titleService:        titleService,
		userService:         userService,
		walletService:       walletService,
		taskService:         taskService,
		transactionService:  transactionService,
		notificationService: notificationService,
	}
}

func (h *TitleHandler) GetAllForUser(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	titles, err := h.titleService.GetAllForUser(user.ID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch title: "+err.Error())
		return
	}

	var transformedTitles []resources.TitleResource

	if len(titles) > 0 {
		transformedTitles = resources.TransformTitles(titles)
	} else {
		transformedTitles = []resources.TitleResource{}
	}

	response := resources.Response{
		Data: transformedTitles,
	}

	c.JSON(http.StatusOK, response)
}

func (h *TitleHandler) ToggleEnableTitle(c *gin.Context) {
	titleIDStr := c.Param("id")

	titleID, err := strconv.ParseUint(titleIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid title ID: "+err.Error())
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	err = h.titleService.ToggleEnableTitle(user.ID, uint(titleID))
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to enable/disable title: "+err.Error())
		return
	}

	cache.GlobalCache.RemoveUserFromCache(user.ID)

	c.JSON(http.StatusOK, gin.H{"message": "The selected title was successfully enabled/disabled!"})
}

func (h *TitleHandler) BuyTitle(c *gin.Context) {
	titleIDStr := c.Param("id")

	titleID, err := strconv.ParseUint(titleIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid title ID: "+err.Error())
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	title, err := h.titleService.FindById(uint(titleID))
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Something went wrong on trying to find the requested title. "+err.Error())
		log.Fatalf("failed to fetch title by id: %+v", err)
		return
	}

	if !title.Purchasable {
		RespondWithError(c, http.StatusBadRequest, "This title is not available for purchase!")
		return
	}

	if title.Cost == nil || user.Wallet.ID == 0 {
		RespondWithError(c, http.StatusBadRequest, "There is a problem with the title cost or your wallet. Please, contact support!")
		return
	}

	if *title.Cost > user.Wallet.Amount {
		RespondWithError(c, http.StatusBadRequest, "Insufficient funds to purchase the title!")
		return
	}

	if err = h.walletService.Subtract(user.ID, uint(*title.Cost)); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to subtract the amount from wallet. "+err.Error())
		log.Fatalf("failed to subtract the amount from user wallet: %+v", err)
		return
	}

	if err = h.taskService.AwardTitleToUser(user.ID, title.ID); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to process the title to user. "+err.Error())
		err := h.walletService.Add(user.ID, uint(*title.Cost))
		if err != nil {
			log.Fatalf("failed to chargeback user wallet amount: %+v", err)
			return
		}

		log.Fatalf("failed to process the title to user: %+v", err)
		return
	}

	purchaseMessage := map[string]any{
		"type": "PurchaseTitle",
		"body": map[string]any{
			"user_id":  user.ID,
			"title_id": title.ID,
			"cost":     uint(*title.Cost),
			"title":    title.Title,
		},
	}

	messageBody, err := json.Marshal(purchaseMessage)
	if err != nil {
		log.Fatalf("failed to serialize purchase message to JSON: %+v", err)
	}

	err = sqs.GlobalSQSClient.SendMessage(c.Request.Context(), sqs.GetAwsQueue(), string(messageBody))
	if err != nil {
		log.Fatalf("failed to enqueue purchase message to SQS: %+v", err)
	}

	cache.GlobalCache.RemoveUserFromCache(user.ID)

	c.JSON(http.StatusOK, gin.H{"message": "You have successfully purchased the selected title!"})
}
