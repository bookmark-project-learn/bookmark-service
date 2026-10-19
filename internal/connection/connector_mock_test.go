package connection

import (
	"errors"
	"testing"

	rabbitmq "github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq/mocks"
	"github.com/bookmark-project-learn/bookmark-service/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
)

func setupRunTestInit(t *testing.T, rabbitMqSenderMock *rabbitmq.RabbitMqSender) {
	fix := fixture.NewBookmarkTestCase(t)
	connectorMock, errConnector := InitDBConnectorMock(t, fix, rabbitMqSenderMock)
	if errConnector != nil {
		t.Fatal(errConnector)
	}
	assert.NoError(t, errConnector)
	assert.NotNil(t, connectorMock)
}

func TestInitDBConnectorMock(t *testing.T) {
	t.Run("run rabbit sender no error", func(t *testing.T) {
		rabbitMqSenderMock := rabbitmq.NewRabbitMqSender(t, nil)
		setupRunTestInit(t, rabbitMqSenderMock)
	})

	t.Run("run rabbit sender no error", func(t *testing.T) {
		rabbitMqSenderMock := rabbitmq.NewRabbitMqSender(t, errors.New("test error"))
		setupRunTestInit(t, rabbitMqSenderMock)
	})

}
