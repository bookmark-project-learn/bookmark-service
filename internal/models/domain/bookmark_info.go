package domain_model

import (
	bookmark_model "github.com/bookmark-project-learn/bookmark-service/internal/models/dto/api/bookmark"
	"github.com/bookmark-project-learn/bookmark-service/internal/models/entity"
)

func ToBookmarkInfo(b *entity.Bookmark) *bookmark_model.BookmarkInfo {
	return &bookmark_model.BookmarkInfo{
		Id:          b.Id,
		Code:        b.Code,
		Description: b.Description,
		Url:         b.Url,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}
