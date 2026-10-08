package bookmark_handler

import (
	"net/http"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/request_ultils"
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/response"
	"github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api"
	"github.com/gin-gonic/gin"
)

// DeleteBookmark
// @Summary Delete bookmark
// @Tags bookmark
// @Accept json
// @Produce json
// @Param id path string true "Bookmark id"
// @Success 200 {object} api.MessageResponse
// @Failure      400  {object}  api.MessageResponse
// @Failure      500  {object}  api.MessageResponse
// @Router /v1/bookmarks/{id} [delete]
// @Security BearerAuth
func (handler *bookmarkHandler) DeleteBookmark(c *gin.Context) {
	id := c.Params.ByName("id")
	// authorization
	userId, err := request_ultils.GetSubjectFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ToMessageReposnse(request_ultils.ClaimsNotFound))
		return
	}

	// call service
	err = handler.svc.DeleteBookmark(c, userId, id)
	res := &api.MessageResponse{
		Message: "Delete bookmark successfully",
	}
	if err != nil {
		res.Message = response.ErrorHandling(err).Error()
		c.JSON(response.MapErrorToHttpCode[err], res)
		return
	}
	c.JSON(http.StatusOK, res)
}
