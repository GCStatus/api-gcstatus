package api

import (
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *usecases.TransactionService
	userService        *usecases.UserService
}

func NewTransactionHandler(
	transactionService *usecases.TransactionService,
	userService *usecases.UserService,
) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		userService:        userService,
	}
}

func (h *TransactionHandler) GetAllForUser(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	transactions, err := h.transactionService.GetAllForUser(user.ID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	var transformedTransactions any

	if len(transactions) > 0 {
		transformedTransactions = resources.TransformTransactions(transactions)
	} else {
		transformedTransactions = []resources.TransactionResource{}
	}

	response := resources.Response{
		Data: transformedTransactions,
	}

	c.JSON(http.StatusOK, response)
}
