package queue

import (
	"context"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq"
)

//go:generate mockery --name QueueService --filename queue.go
type QueueService interface {
	SendImportBookmarksJob(ctx context.Context, message *ImportMessage) error
}

type queueService struct {
	rabbitMqSender rabbitmq.RabbitMqSender
}

func NewQueueService(rabbitMqSender rabbitmq.RabbitMqSender) QueueService {
	return &queueService{rabbitMqSender: rabbitMqSender}
}

// Model
type ImportMessage struct {
	UID       string                 `json:"uid"`
	Bookmarks []*ImportBookmarkInput `json:"data"`
}

type ImportBookmarkInput struct {
	Description string `csv:"description" json:"description" validate:"lte=255" `
	URL         string `csv:"url" json:"url" validate:"required,url"`
}
