package queue

import (
	"context"
	"errors"

	"testing"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq"
	mocks_rabbitmq "github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq/mocks"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test error")

func TestRabbitMqSender_Publish(t *testing.T) {
	testCase := []struct {
		name        string
		mockFn      func(ctx context.Context) rabbitmq.RabbitMqSender
		message     *ImportMessage
		errExpected error
	}{
		{
			name: "success",
			mockFn: func(ctx context.Context) rabbitmq.RabbitMqSender {
				mockInternalPublisher := mocks_rabbitmq.NewRabbitMqSender(t, nil)
				return mockInternalPublisher
			},
			message: &ImportMessage{
				UID: "1",
				Bookmarks: []*ImportBookmarkInput{
					{
						Description: "test",
						URL:         "https://google.com",
					},
				},
			},
			errExpected: nil,
		},
		{
			name: "failed",
			mockFn: func(ctx context.Context) rabbitmq.RabbitMqSender {
				mockInternalPublisher := mocks_rabbitmq.NewRabbitMqSender(t, errTest)
				return mockInternalPublisher
			},
			message: &ImportMessage{
				UID: "1",
				Bookmarks: []*ImportBookmarkInput{
					{
						Description: "test",
						URL:         "https://google.com",
					},
				},
			},
			errExpected: errTest,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockInternalPublisher := tc.mockFn(ctx)
			queueSender := NewQueueService(mockInternalPublisher)
			err := queueSender.SendImportBookmarksJob(ctx, tc.message)
			assert.Equal(t, tc.errExpected, err)
		})
	}
}
