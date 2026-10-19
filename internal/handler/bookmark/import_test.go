package bookmark_handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bookmark-project-learn/bookmark-service/internal/service/queue"
	queue_mocks "github.com/bookmark-project-learn/bookmark-service/internal/service/queue/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestImportBookmarks(t *testing.T) {
	queueError := errors.New("queue unavailable")
	validCSV := "description,url\ndemo,https://example.com\n"

	testCases := []struct {
		name         string
		withFile     bool
		fileContent  string
		fileType     string
		setClaims    bool
		claims       any
		queueError   error
		hasQueueCall bool
		expectedCode int
		expectedText string
	}{
		{
			name:         "no file uploaded",
			expectedCode: http.StatusBadRequest,
			expectedText: "No file uploaded",
		},
		{
			name:         "file exceeds size limit",
			withFile:     true,
			fileContent:  strings.Repeat("x", LimitCSVFileSize+1),
			fileType:     "text/csv",
			expectedCode: http.StatusBadRequest,
			expectedText: "File size exceeds the limit.",
		},
		{
			name:         "invalid file type",
			withFile:     true,
			fileContent:  validCSV,
			fileType:     "application/json",
			expectedCode: http.StatusBadRequest,
			expectedText: "Invalid file type. Please upload a CSV file.",
		},
		{
			name:         "invalid CSV",
			withFile:     true,
			claims:       jwt.MapClaims{"sub": "user-1"},
			fileContent:  "description,url\ndemo,\"https://example.com\n",
			fileType:     "text/csv",
			expectedCode: http.StatusBadRequest,
			expectedText: "Failed to parse CSV file",
		},
		{
			name:         "invalid bookmark URL",
			withFile:     true,
			fileContent:  "description,url\ndemo,not-a-url\n",
			fileType:     "text/csv",
			expectedCode: http.StatusBadRequest,
			expectedText: "URL",
		},
		{
			name:         "missing user claims",
			withFile:     true,
			fileContent:  validCSV,
			fileType:     "text/csv",
			expectedCode: http.StatusUnauthorized,
			expectedText: "Not found user id in claims",
		},
		{
			name:         "queue error",
			withFile:     true,
			fileContent:  validCSV,
			fileType:     "text/csv",
			setClaims:    true,
			claims:       jwt.MapClaims{"sub": "user-1"},
			queueError:   queueError,
			hasQueueCall: true,
			expectedCode: http.StatusInternalServerError,
			expectedText: "Failed to send bookmarks imports to queue",
		},
		{
			name:         "success",
			withFile:     true,
			fileContent:  validCSV,
			fileType:     "text/csv",
			setClaims:    true,
			claims:       jwt.MapClaims{"sub": "user-1"},
			hasQueueCall: true,
			expectedCode: http.StatusOK,
			expectedText: "Successfully sent bookmarks imports to queue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			requestBody, writer, err := CreateMultipartWriterWithContentString(tc.withFile, tc.fileType, tc.fileContent)
			if err != nil {
				t.Fatal(err.Error())
			}

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks/import", requestBody)
			ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			if tc.setClaims {
				ctx.Set("claims", tc.claims)
			}

			queueMock := queue_mocks.NewQueueService(t)
			if tc.hasQueueCall {
				queueMock.On("SendImportBookmarksJob", ctx, mock.MatchedBy(func(message *queue.ImportMessage) bool {
					return message != nil &&
						message.UID == "user-1" &&
						len(message.Bookmarks) == 1 &&
						message.Bookmarks[0].Description == "demo" &&
						message.Bookmarks[0].URL == "https://example.com"
				})).Return(tc.queueError)
			}

			handler := NewBookmarkHandler(nil, queueMock)
			handler.ImportBookmarks(ctx)

			assert.Equal(t, tc.expectedCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedText)
		})
	}
}
