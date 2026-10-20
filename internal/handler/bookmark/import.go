package bookmark_handler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"slices"
	"strings"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/csv"
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/request_ultils"
	response_pkg "github.com/bookmark-project-learn/bookmark-common-libs/pkg/response"
	"github.com/bookmark-project-learn/bookmark-service/internal/service/queue"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var fileExtensions = []string{"text/csv", "application/csv", "text/plain"}

const LimitCSVFileSize = 10 << 20 // 10 MiB
// ImportBookmarks imports bookmarks from a CSV file
// @Summary Import bookmarks from a CSV file
// @Tags bookmark
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file to import"
// @Success 200 {object} ImportMessageResponse
// @Failure      400  {object}  ImportMessageResponse
// @Failure      500  {object}  ImportMessageResponse
// @Router /v1/bookmarks/import [post]
// @Security BearerAuth
func (h *bookmarkHandler) ImportBookmarks(c *gin.Context) {
	file, err := c.FormFile("file")
	if file == nil {
		response := &ImportMessageResponse{Message: "No file uploaded"}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	response := &ImportMessageResponse{Message: "Successfully sent bookmarks imports to queue"}
	if err != nil {
		response.Message = "Failed to get file from request"
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if file.Size > LimitCSVFileSize {
		response.Message = "File size exceeds the limit."
		c.JSON(http.StatusBadRequest, response)
		return
	}
	fileType := file.Header.Get("Content-Type")
	if !slices.Contains(fileExtensions, fileType) {
		response.Message = "Invalid file type. Please upload a CSV file."
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// parse input file and send to queue for processing
	var importMessageInputs []*queue.ImportBookmarkInput
	f, _ := file.Open()
	defer f.Close()
	byteContainer, err := ioutil.ReadAll(f)
	log.Printf("size:%d", len(byteContainer))

	if err := csv.ParseCSVFile(file, &importMessageInputs); err != nil {
		response.Message = "Failed to parse CSV file"
		log.Error().Err(err).Msg("Failed to parse CSV file")
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// validate the parsed bookmarks
	if err := request_ultils.InputValidator.Var(importMessageInputs, "dive"); err != nil {
		message_error := response_pkg.InputFieldError(err)
		response.Message = strings.Join(message_error.Details.([]string), "; ")
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// send message to rabbitmq queue for processing
	userId, err := request_ultils.GetSubjectFromClaims(c)
	if err != nil {
		response.Message = "Not found user id in claims"
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	importMessage := &queue.ImportMessage{UID: userId, Bookmarks: importMessageInputs}
	if err := h.queue.SendImportBookmarksJob(c, importMessage); err != nil {
		response.Message = "Failed to send bookmarks imports to queue"
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	// return success response
	c.JSON(http.StatusOK, response)
}

type ImportMessageResponse struct {
	Message string `json:"message"`
}

func CreateMultipartWriterWithContentString(hasfile bool, fileType string, content string) (*bytes.Buffer, *multipart.Writer, error) {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	if hasfile {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", `form-data; name="file"; filename="bookmarks.csv"`)
		header.Set("Content-Type", fileType)
		part, err := writer.CreatePart(header)
		if err != nil {
			return nil, nil, errors.New("create multipart file part: " + err.Error())
		}
		if _, err := part.Write([]byte(content)); err != nil {
			return nil, nil, errors.New("write multipart file part: " + err.Error())
		}
	}

	if err := writer.Close(); err != nil {
		return nil, nil, errors.New("close multipart writer: " + err.Error())
	}
	return requestBody, writer, nil
}
