package bookmark_service

import (
	"context"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/response"
	domain_model "github.com/bookmark-project-learn/bookmark-service/internal/models/domain"
	bookmark_model "github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api/bookmark"
)

const PrefixCreate = "bookmark_url_"

// NewBookmark
func (s *bookmarkService) NewBookmark(ctx context.Context, userId string, bookmark *bookmark_model.NewBookmarkRequest) (*bookmark_model.BookmarkInfo, error) {
	bookmarkCreate, err := s.bookmarkRepository.CreateBookmark(ctx, userId, bookmark.Url, bookmark.Description, PrefixCreate)

	if err != nil {
		return nil, response.ErrorHandling(err)
	}

	return domain_model.ToBookmarkInfo(bookmarkCreate), nil
}
