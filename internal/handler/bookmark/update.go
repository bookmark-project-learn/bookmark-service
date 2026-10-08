package bookmark_handler

import (
	"net/http"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/request_ultils"
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/response"
	"github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api"
	bookmark_model "github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api/bookmark"
	"github.com/gin-gonic/gin"
)

// UpdateBookmark
// @Summary Update bookmark
// @Tags bookmark
// @Accept json
// @Produce json
// @Param id path string true "Bookmark id"
// @Param input body bookmark_model.UpdateBookmarkRequest true "Input required"
// @Success 200 {object} api.MessageResponse
// @Failure      400  {object}  api.MessageResponse
// @Failure      500  {object}  api.MessageResponse
// @Router /v1/bookmarks/{id} [put]
// @Security BearerAuth
func (handler *bookmarkHandler) UpdateBookmark(c *gin.Context) {
	id := c.Params.ByName("id")
	res := &api.MessageResponse{
		Message: "Update bookmark successfully",
	}
	// model validation
	request, err := request_ultils.ModelBindValidation[bookmark_model.UpdateBookmarkRequest](c)
	if err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}
	// authorizaton
	userId, err := request_ultils.GetSubjectFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ToDataResponse[*bookmark_model.BookmarkInfo](request_ultils.ClaimsNotFound))
		return
	}
	// call service
	err = handler.svc.UpdateBookmark(c, request, userId, id)
	if err != nil {
		res.Message = response.ErrorHandling(err).Error()
		c.JSON(response.MapErrorToHttpCode[err], res)
		return
	}
	c.JSON(http.StatusOK, res)
}
