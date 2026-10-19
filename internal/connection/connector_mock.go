package connection

import (
	"testing"

	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq"
	redisPkg "github.com/bookmark-project-learn/bookmark-common-libs/pkg/redis"
	"github.com/bookmark-project-learn/bookmark-service/internal/test/data/fixture"
)

func InitDBConnectorMock(t *testing.T, fix fixture.Fixture, rabbitMqSenderMock rabbitmq.RabbitMqSender) (DBConnector, error) {
	// fix := fixture.NewUserTestCase(t)
	gormSql := fixture.NewFixture(t, fix)
	return NewDBConnector(redisPkg.InitMockRedis(t), gormSql, rabbitMqSenderMock), nil
}
