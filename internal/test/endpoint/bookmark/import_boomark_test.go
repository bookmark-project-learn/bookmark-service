package bookmark_endpoint

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt_pkg "github.com/bookmark-project-learn/bookmark-common-libs/pkg/jwt"
	"github.com/bookmark-project-learn/bookmark-service/internal/test/data/fixture"
	general_helpers "github.com/bookmark-project-learn/bookmark-service/internal/test/general"
	"github.com/stretchr/testify/assert"
)

func TestImportBookmark(t *testing.T) {
	validCSV := "description,url\ndemo,https://example.com\n"
	invalidCSV := "description,url\ndemo,https://example.com\n\""
	errClose := errors.New("close error")
	testCases := []struct {
		name               string
		hasFile            bool
		content            string
		fileType           string
		expectedStatusCode int
		expectedContain    string
		pushMessageError   error
		disableAuth        bool
	}{
		{
			name:               "import bookmark successfully",
			content:            validCSV,
			hasFile:            true,
			fileType:           "text/csv",
			expectedStatusCode: http.StatusOK,
			expectedContain:    "Successfully sent bookmarks imports to queue",
		},
		{
			name:               "no file uploaded",
			hasFile:            false,
			expectedStatusCode: http.StatusBadRequest,
			expectedContain:    "No file uploaded",
		},
		{
			name:               "invalid CSV",
			content:            invalidCSV,
			hasFile:            true,
			fileType:           "text/csv",
			expectedStatusCode: http.StatusBadRequest,
			expectedContain:    "Failed to parse CSV file",
		},
		{
			name:               "invalid bookmark URL",
			content:            "description,url\ndemo,not-a-url\n",
			hasFile:            true,
			fileType:           "text/csv",
			expectedStatusCode: http.StatusBadRequest,
			expectedContain:    "URL is invalid",
		},
		{
			name:               "no authorization",
			content:            validCSV,
			hasFile:            true,
			fileType:           "text/csv",
			expectedStatusCode: http.StatusUnauthorized,
			expectedContain:    "Unauthorized",
			disableAuth:        true,
		},
		{
			name:               "queue error",
			content:            validCSV,
			hasFile:            true,
			fileType:           "text/csv",
			expectedStatusCode: http.StatusInternalServerError,
			expectedContain:    "Failed to send bookmarks imports to queue",
			pushMessageError:   errClose,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			jwtMock := jwt_pkg.NewMockJwt()

			router := BuildBookmarkApiEngine(testItem, fixture.NewBookmarkTestCase(testItem), jwtMock, GetRabbitMqMock(testItem, tc.pushMessageError))

			// http
			req, err := general_helpers.MakeMultipartRequest(buildBookmarkURL("/import"), tc.hasFile, tc.content, tc.fileType)
			if err != nil {
				t.Fatal("failed to setup request in test case: " + tc.name)
				return
			}
			resp := httptest.NewRecorder()
			if !tc.disableAuth {
				RouterSetupAuthorization(jwtMock, resp, req)
			} else {
				t.Log("No Authorization")
				t.Log(req.Header.Get("Authorization"))
			}
			router.ServeHTTP(resp, req)

			assert.Equal(testItem, tc.expectedStatusCode, resp.Code)
			assert.Contains(testItem, resp.Body.String(), tc.expectedContain)
		})
	}
}
