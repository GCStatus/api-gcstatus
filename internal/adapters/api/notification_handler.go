package api

import (
	"gcstatus/internal/errors"
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *usecases.NotificationService
	userService         *usecases.UserService
}

func NewNotificationHandler(
	notificationService *usecases.NotificationService,
	userService *usecases.UserService,
) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		userService:         userService,
	}
}

func (h *NotificationHandler) GetAllForUser(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	notifications, err := h.notificationService.GetAllForUser(user.ID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch you notifications! "+err.Error())
		log.Fatalf("failed to fetch user notifications: %+v", err)
		return
	}

	var transformedNotifications any

	if len(notifications) > 0 {
		transformedNotifications = resources.TransformNotifications(notifications)
	} else {
		transformedNotifications = []resources.NotificationResource{}
	}

	response := resources.Response{
		Data: transformedNotifications,
	}

	c.JSON(http.StatusOK, response)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	h.markNotification(c, true)
}

func (h *NotificationHandler) MarkAsUnread(c *gin.Context) {
	h.markNotification(c, false)
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	notificationIDStr := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid notification ID: "+err.Error())
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	if err := h.notificationService.DeleteNotification(user.ID, uint(notificationID)); err != nil {
		if httpErr, ok := err.(*errors.HttpError); ok {
			RespondWithError(c, httpErr.Code, httpErr.Message)
			return
		}
		RespondWithError(c, http.StatusInternalServerError, "Failed to delete notification: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The notification was successfully deleted."})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	if err = h.notificationService.MarkAllAsRead(user.ID); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to mark notifications as read: "+err.Error())
		log.Fatalf("failed to mark user %+v notifications as read: %+v", user, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your notifications was successfully marked as read!"})
}

func (h *NotificationHandler) MarkAllAsUnread(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	if err = h.notificationService.MarkAllAsUnread(user.ID); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to mark notifications as unread: "+err.Error())
		log.Fatalf("failed to mark user %+v notifications as unread: %+v", user, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your notifications was successfully marked as unread!"})
}

func (h *NotificationHandler) DeleteAllNotifications(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	if err = h.notificationService.DeleteAllNotifications(user.ID); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to delete notifications: "+err.Error())
		log.Fatalf("failed to delete user %+v notifications: %+v", user, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your notifications was successfully deleted!"})
}

func (h *NotificationHandler) markNotification(c *gin.Context, asRead bool) {
	notificationIDStr := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid notification ID: "+err.Error())
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	var responseMessage string
	if asRead {
		if err := h.notificationService.MarkAsRead(user.ID, uint(notificationID)); err != nil {
			if httpErr, ok := err.(*errors.HttpError); ok {
				RespondWithError(c, httpErr.Code, httpErr.Message)
				return
			}
			RespondWithError(c, http.StatusInternalServerError, "Failed to mark notification as read: "+err.Error())
			return
		}
		responseMessage = "The notification was marked as read successfully."
	} else {
		if err := h.notificationService.MarkAsUnread(user.ID, uint(notificationID)); err != nil {
			if httpErr, ok := err.(*errors.HttpError); ok {
				RespondWithError(c, httpErr.Code, httpErr.Message)
				return
			}
			RespondWithError(c, http.StatusInternalServerError, "Failed to mark notification as unread: "+err.Error())
			return
		}
		responseMessage = "The notification was marked as unread successfully."
	}

	c.JSON(http.StatusOK, gin.H{"message": responseMessage})
}
