package bookmark_handler

import (
	"net/http"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/request_ultils"
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/response"
	"github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api"
	_ "github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api"
	bookmark_model "github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api/bookmark"
	"github.com/gin-gonic/gin"
)

var SuccessCreateBookmark = "success create bookmark"

// CreateBookmark
// @Summary Create bookmark
// @Tags bookmark
// @Accept json
// @Produce json
// @Param input body bookmark_model.NewBookmarkRequest true "Input required"
// @Success 200 {object} api.Response[bookmark_model.BookmarkInfo]
// @Failure      400  {object}  api.Response[bookmark_model.BookmarkInfo]
// @Failure      500  {object}  api.Response[bookmark_model.BookmarkInfo]
// @Router /v1/bookmarks [post]
// @Security BearerAuth
func (h *bookmarkHandler) CreateBookmark(c *gin.Context) {
	request, err := request_ultils.ModelBindValidation[bookmark_model.NewBookmarkRequest](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ToDataResponse[*bookmark_model.BookmarkInfo](err))
		return
	}

	userId, err := request_ultils.GetSubjectFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ToDataResponse[*bookmark_model.BookmarkInfo](request_ultils.ClaimsNotFound))
		return
	}

	bookmark, err := h.svc.NewBookmark(c, userId, request)

	res := &api.Response[bookmark_model.BookmarkInfo]{}
	if err != nil {
		res.Message = response.ErrorHandling(err).Error()
		c.JSON(response.MapErrorToHttpCode[err], res)
		return
	}
	res.Data = bookmark
	res.Message = SuccessCreateBookmark
	c.JSON(http.StatusOK, res)
}
