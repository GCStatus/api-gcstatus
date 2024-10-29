package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
	"gcstatus/pkg/sqs"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService *usecases.ProfileService
	userService    *usecases.UserService
}

func NewProfileHandler(
	profileService *usecases.ProfileService,
	userService *usecases.UserService,
) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		userService:    userService,
	}
}

func (h *ProfileHandler) UpdatePicture(c *gin.Context) {
	var request struct {
		File *multipart.FileHeader `json:"file" form:"file" binding:"required"`
	}

	if err := c.ShouldBind(&request); err != nil {
		errorMessages := utils.FormatValidationError(err)
		RespondWithError(c, http.StatusUnprocessableEntity, "Invalid request data: "+strings.Join(errorMessages, " "))
		return
	}

	contentType := request.File.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		RespondWithError(c, http.StatusUnprocessableEntity, "Only image files are allowed.")
		return
	}

	file, err := request.File.Open()
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "We could not open the uploaded file. Please, try again.")
		return
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("failed to close file: %s", err.Error())
		}
	}()

	fileSize := request.File.Size
	if fileSize > 5*1024*1024 {
		RespondWithError(c, http.StatusUnprocessableEntity, "The file size is too high. The max file size is up to 5MB.")
		return
	}

	fileContent := make([]byte, fileSize)
	if _, err := file.Read(fileContent); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "We could not open the uploaded file. Please, try again.")
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	folder := "profiles"
	fileName := fmt.Sprintf("%s's-profile.png", user.Nickname)

	ctx := context.TODO()

	if len(user.Profile.Photo) > 0 {
		err := s3.GlobalS3Client.RemoveFile(ctx, user.Profile.Photo)

		if err != nil {
			log.Fatalf("failed to remove profile picture from s3 server: %s", err.Error())
		}
	}

	filePath, err := s3.GlobalS3Client.UploadFile(ctx, folder, fileName, fileContent)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to upload file to S3. Please, try again later."+err.Error())
		return
	}

	err = h.profileService.UpdatePicture(user.Profile.ID, filePath)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to update your profile picture: "+err.Error())
		return
	}

	h.createTrackProfilePictureSQS(c, user)

	c.JSON(http.StatusOK, gin.H{"message": "Your profile picture was successfully updated!"})
}

func (h *ProfileHandler) UpdateSocials(c *gin.Context) {
	var request ports.UpdateSocialsRequest

	if err := c.ShouldBind(&request); err != nil {
		errorMessages := utils.FormatValidationError(err)
		RespondWithError(c, http.StatusUnprocessableEntity, "Invalid request data:"+strings.Join(errorMessages, " "))
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	err = h.profileService.UpdateSocials(user.Profile.ID, request)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not update socials: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your profile socials was successfully updated!"})
}

func (h *ProfileHandler) createTrackProfilePictureSQS(c *gin.Context, user *domain.User) {
	trackProgressMessage := map[string]any{
		"type": "TrackProgressProfilePicture",
		"body": map[string]any{
			"user_id":   user.ID,
			"increment": 1,
		},
	}

	trackProgressMessageBody, err := json.Marshal(trackProgressMessage)
	if err != nil {
		log.Fatalf("failed to serialize track progress profile picture message to JSON: %+v", err)
	}

	if err := sqs.GlobalSQSClient.SendMessage(c.Request.Context(), sqs.GetAwsQueue(), string(trackProgressMessageBody)); err != nil {
		log.Fatalf("failed to enqueue track progress profile picture message to SQS: %+v", err)
	}
}
