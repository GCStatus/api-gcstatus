package api

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService *usecases.UserService
}

func NewUserHandler(userService *usecases.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transformedUsers := resources.TransformUsers(users, s3.GlobalS3Client)

	response := resources.Response{
		Data: transformedUsers,
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUserNickAndEmail(c *gin.Context) {
	var request ports.UpdateNickAndEmailRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorMessages := utils.FormatValidationError(err)
		RespondWithError(c, http.StatusUnprocessableEntity, "Invalid request data: "+strings.Join(errorMessages, " "))
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	existsNickUser, err := h.userService.FindUserByEmailOrNickname(request.Nickname)
	if err != nil && err != gorm.ErrRecordNotFound {
		RespondWithError(c, http.StatusInternalServerError, "Something went wrong. Please, try again later.")
		log.Fatalf("failed to fetch user by nickname: %s", err.Error())
		return
	}

	if existsNickUser != nil && user.ID != existsNickUser.ID {
		RespondWithError(c, http.StatusConflict, "This nickname is already in use.")
		return
	}

	existsEmailUser, err := h.userService.FindUserByEmailOrNickname(request.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		RespondWithError(c, http.StatusInternalServerError, "Something went wrong. Please, try again later.")
		log.Fatalf("failed to fetch user by email: %s", err.Error())
		return
	}

	if existsEmailUser != nil && user.ID != existsEmailUser.ID {
		RespondWithError(c, http.StatusConflict, "This email is already in use.")
		return
	}

	if equal := utils.IsHashEqualsValue(user.Password, request.Password); !equal {
		RespondWithError(c, http.StatusBadRequest, "Your password does not match. Double check it and try again!")
		return
	}

	err = h.userService.UpdateUserNickAndEmail(user.ID, request)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to update your nickname or email: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your nickname or email was successfully updated!"})
}

func (h *UserHandler) UpdateUserBasics(c *gin.Context) {
	var request ports.UpdateUserBasicsRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorMessages := utils.FormatValidationError(err)
		RespondWithError(c, http.StatusUnprocessableEntity, "Invalid request data: "+strings.Join(errorMessages, " "))
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	birthdate, err := time.Parse("2006-01-02T15:04:05", request.Birthdate)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid birthdate format.")
		return
	}

	if time.Since(birthdate).Hours() < 14*365*24 {
		RespondWithError(c, http.StatusBadRequest, "You must be at least 14 years old.")
		return
	}

	err = h.userService.UpdateUserBasics(user.ID, request)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to update your informations: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your informations was successfully updated!"})
}
