package bookmark_handler

import (
	bookmark_cache "github.com/bookmark-project-learn/bookmark-service/internal/cache/bookmark"
	"github.com/bookmark-project-learn/bookmark-service/internal/service/queue"
	"github.com/gin-gonic/gin"
)

type BookmarkHandler interface {
	CreateBookmark(c *gin.Context)
	GetBookmarks(c *gin.Context)
	UpdateBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
	ImportBookmarks(c *gin.Context)
}

type bookmarkHandler struct {
	svc   bookmark_cache.BookmarkCache
	queue queue.QueueService
}

func NewBookmarkHandler(svc bookmark_cache.BookmarkCache, queue queue.QueueService) BookmarkHandler {
	return &bookmarkHandler{svc, queue}
}
