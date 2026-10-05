package bookmark_handler

import (
	bookmark_cache "github.com/bookmark-project-learn/bookmark-service/internal/cache/bookmark"
	"github.com/gin-gonic/gin"
)

type BookmarkHandler interface {
	CreateBookmark(c *gin.Context)
	GetBookmarks(c *gin.Context)
	UpdateBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
}

type bookmarkHandler struct {
	svc bookmark_cache.BookmarkCache
}

func NewBookmarkHandler(svc bookmark_cache.BookmarkCache) BookmarkHandler {
	return &bookmarkHandler{svc}
}
