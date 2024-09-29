package api

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/cache"
	"gcstatus/pkg/email"
	"gcstatus/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TitleHandler struct {
	titleService       *usecases.TitleService
	userService        *usecases.UserService
	walletService      *usecases.WalletService
	taskService        *usecases.TaskService
	transactionService *usecases.TransactionService
}

func NewTitleHandler(
	titleService *usecases.TitleService,
	userService *usecases.UserService,
	walletService *usecases.WalletService,
	taskService *usecases.TaskService,
	transactionService *usecases.TransactionService,
) *TitleHandler {
	return &TitleHandler{
		titleService:       titleService,
		userService:        userService,
		walletService:      walletService,
		taskService:        taskService,
		transactionService: transactionService,
	}
}

func (h *TitleHandler) GetAllForUser(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch titles: "+err.Error())
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
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user: "+err.Error())
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
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch titles: "+err.Error())
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

	transaction := &domain.Transaction{
		Amount:            uint(*title.Cost),
		Description:       fmt.Sprintf("Purchase of title %s by %v coins.", title.Title, uint(*title.Cost)),
		UserID:            user.ID,
		TransactionTypeID: domain.SubtractionTransactionTypeID,
	}

	if err = h.transactionService.CreateTransaction(transaction); err != nil {
		log.Fatalf("failed to create a transaction for user title purchase: %+v", err)
	}

	if err = email.SendTransactionEmail(user, transaction, email.Send); err != nil {
		log.Fatalf("failed to send transaction email: %+v", err)
	}

	cache.GlobalCache.RemoveUserFromCache(user.ID)

	c.JSON(http.StatusOK, gin.H{"message": "You have successfully purchased the selected title!"})
}
