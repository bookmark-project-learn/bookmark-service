package bookmark_service

import (
	"context"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/response"
	bookmark_model "github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api/bookmark"
)

// UpdateBookmark
func (s *bookmarkService) UpdateBookmark(ctx context.Context, request *bookmark_model.UpdateBookmarkRequest, userId, bookmarkId string) error {
	err := s.bookmarkRepository.UpdateBookmark(ctx, userId, bookmarkId, request.Url, request.Description)
	if err != nil {
		return response.ErrorHandling(err)
	}
	return nil
}
