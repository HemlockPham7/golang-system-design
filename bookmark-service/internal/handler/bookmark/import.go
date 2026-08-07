package bookmark

import (
	"net/http"
	"slices"

	"github.com/HemlockPham7/golang-system-design/internal/service/queue"
	"github.com/HemlockPham7/golang-system-design/pkg/csv"
	"github.com/HemlockPham7/golang-system-design/pkg/requestutils"
	"github.com/HemlockPham7/golang-system-design/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	MaxFileSize = 10 << 20
)

var allowedFileTypes = []string{"text/csv", "aplication/csv", "text/plain", "application/octet-stream"}

// ImportBookmarks handles file uploads
// @Summary Upload and parse a csv file of bookmarks
// @Description Accepts a csv file and import bookmarks
// @Tags bookmark
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file"
// @Success 200 {object} object{message=string} "Success"
// @Router /v1/bookmarks/import [post]
func (h *bookmarkHandler) ImportBookmarks(c *gin.Context) {
	// Get uid
	uid, err := requestutils.GetUserIDFromRequest(c)
	if err != nil {
		return
	}

	// Get .csv file in request
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &response.Message{
			Message: "Invalid file",
		})
		return
	}

	// check file size
	if file.Size > MaxFileSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, &response.Message{
			Message: "File size is too large",
		})
		return
	}

	// validate file
	fileType := file.Header.Get("Content-Type")
	if !slices.Contains(allowedFileTypes, fileType) {
		c.AbortWithStatusJSON(http.StatusBadRequest, &response.Message{
			Message: "Invalid file type",
		})
		return
	}

	// parse file to struct (import message)
	var importInput []*queue.ImportBookmarkInput
	err = csv.ParseFromMultipartFile(file, &importInput)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &response.Message{
			Message: "Unable to parse file",
		})
		return
	}

	err = requestutils.InputValidator.Var(importInput, "dive")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	// sent message to queue
	err = h.messageQueue.SendImportBookmarkJob(c, uid, importInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send import bookmark job")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InputErrResponse)
		return
	}

	c.JSON(http.StatusOK, &response.Message{
		Message: "Import bookmark job sent successfully",
	})
}
