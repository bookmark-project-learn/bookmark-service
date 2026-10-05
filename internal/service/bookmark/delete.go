package bookmark_service

import (
	"context"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/response"
)

// DeleteBookmark
func (s *bookmarkService) DeleteBookmark(ctx context.Context, userId, bookmarkId string) error {
	err := s.bookmarkRepository.DeleteBookmark(ctx, userId, bookmarkId)
	if err != nil {
		return response.ErrorHandling(err)
	}
	return nil
}
