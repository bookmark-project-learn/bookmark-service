package bookmark_handler

import (
	"net/http"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/request_ultils"
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/response"

	"github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api"
	bookmark_model "github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api/bookmark"
	"github.com/gin-gonic/gin"
)

// GetBookmarks
// @Summary Get bookmarks
// @Tags bookmark
// @Accept json
// @Produce json
// @Param query query bookmark_model.GetBookmarksQuery true "Query required"
// @Success 200 {object} api.PaginatedResponse[bookmark_model.BookmarkInfo]
// @Failure      400  {object}  api.MessageResponse
// @Failure      500  {object}  api.MessageResponse
// @Router /v1/bookmarks [get]
// @Security BearerAuth
func (h *bookmarkHandler) GetBookmarks(c *gin.Context) {
	var query bookmark_model.GetBookmarksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := request_ultils.GetSubjectFromClaims(c)
	error_response := &api.MessageResponse{}
	if err != nil {
		error_response.Message = request_ultils.ClaimsNotFound.Error()
		c.JSON(http.StatusUnauthorized, error_response)
		return
	}

	// call service
	res, err := h.svc.GetBookmarks(c, userId, &query)
	if err != nil {
		error_response.Message = response.ErrorHandling(err).Error()
		c.JSON(http.StatusInternalServerError, error_response)
		return
	}

	c.JSON(http.StatusOK, res)
}
