package bookmark_cache

import (
	"context"

	bookmark_model "github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api/bookmark"
)

// UpdateBookmark
func (b *BookmarkCacheInstance) UpdateBookmark(ctx context.Context, request *bookmark_model.UpdateBookmarkRequest, userId, bookmarkId string) error {
	err := b.deleteGroupKey(ctx, userId)
	if err != nil {
		return err
	}
	return b.service.UpdateBookmark(ctx, request, userId, bookmarkId)
}
