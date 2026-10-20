package bookmark_endpoint

import (
	"fmt"
	"net/http"
	"testing"

	jwt_pkg "github.com/bookmark-project-learn/bookmark-common-libs/pkg/jwt"
	rabbitmq "github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq"
	rabbitmq_mocks "github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq/mocks"
	"github.com/bookmark-project-learn/bookmark-service/internal/api"
	"github.com/bookmark-project-learn/bookmark-service/internal/config"
	"github.com/bookmark-project-learn/bookmark-service/internal/connection"
	"github.com/bookmark-project-learn/bookmark-service/internal/infrastructure"
	"github.com/bookmark-project-learn/bookmark-service/internal/test/data/fixture"
	"github.com/golang-jwt/jwt/v5"
)

const testBookmarkUserID = "4e90220a-51f6-49e4-bc0e-44e2f321475a"

func buildBookmarkIntegrationConfig() *config.Config {
	return &config.Config{
		AppPort:     "8080",
		ServiceName: "app_service",
		InstanceID:  "instance_01",
	}
}

func BuildBookmarkApiEngine(t *testing.T, fix fixture.Fixture, jwtMock *jwt_pkg.MockJwt, rabbitMqSenderMock rabbitmq.RabbitMqSender) api.Engine {
	t.Helper()
	connectorMock, err := connection.InitDBConnectorMock(t, fix, rabbitMqSenderMock)
	if err != nil {
		t.Fatal(err)
	}

	return infrastructure.CoreNewApi(connectorMock.GetRedisClient(), connectorMock.GetSqlDB(), connectorMock.GetRabbitMqSender(), buildBookmarkIntegrationConfig(), jwtMock.JwtGenarate, jwtMock.JwtValidate)
}

func GetRabbitMqMock(t *testing.T, errorPublishMessageExpected error) *rabbitmq_mocks.RabbitMqSender {
	rabbitMqSenderMock := rabbitmq_mocks.NewRabbitMqSender(t, errorPublishMessageExpected)
	return rabbitMqSenderMock

}

var dummyUserId = "4e90220a-51f6-49e4-bc0e-44e2f321475a"

func RouterSetupAuthorization(jwtMock *jwt_pkg.MockJwt, w http.ResponseWriter, req *http.Request) {
	claims := jwt.MapClaims{
		"sub": dummyUserId,
	}
	token, err := jwtMock.JwtGenarate.GenerateToken(claims)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

}

func buildBookmarkURL(path string) string {
	return fmt.Sprintf("/v1/bookmarks%s", path)
}
