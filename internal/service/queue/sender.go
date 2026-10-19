package queue

import (
	"context"
	"encoding/json"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/array"
)

const BatchSize = 10

func (q *queueService) SendImportBookmarksJob(ctx context.Context, message *ImportMessage) error {
	batchs := array.SplitIntoBatch(message.Bookmarks, BatchSize)

	for _, bookmarkBatch := range batchs {
		err := q.sendJob(ctx, message.UID, bookmarkBatch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (q *queueService) sendJob(ctx context.Context, id string, bookmarks []*ImportBookmarkInput) error {
	batch_message := &ImportMessage{UID: id, Bookmarks: bookmarks}
	data, err := json.Marshal(batch_message)
	if err != nil {
		return err
	}
	// push to rabbitmq
	err = q.rabbitMqSender.Publish(ctx, data)
	return err
}
