package infrastructure

import (
	"context"
	"time"

	"github.com/bookmark-project-learn/bookmark-service/constant"
	"github.com/bookmark-project-learn/bookmark-service/internal/api"
	"github.com/bookmark-project-learn/bookmark-service/internal/config"
	"github.com/bookmark-project-learn/bookmark-service/internal/connection"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	jwt_pkg "github.com/bookmark-project-learn/bookmark-common-libs/pkg/jwt"
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq"
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/sqldb"
)

func CreateApi() api.Engine {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// create redis client
	rdClient := CreateRedisClient()

	// create SQL client
	sqlClient := CreateSqlClient()

	migrator := sqldb.BuildMigrate(sqlClient, constant.MigrationPath)
	migrator.SetLogging()
	err = migrator.MigrateUp()
	if err != nil {
		panic("start migration failed " + err.Error())
	}

	// connector
	ctx_sender, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := rabbitmq.NewConnection(ctx_sender, cfg.QueueName)
	if err != nil {
		panic("failed to connect to rabbitmq: " + err.Error())
	}
	rabbitMqSender, err := rabbitmq.NewRabbitSender(ctx_sender, cfg.QueueName, conn)
	if err != nil {
		panic("failed to connect to rabbitmq: " + err.Error())
	}
	jwtGenerator, jwtValidator := CreateJwtProvider()

	return CoreNewApi(rdClient, sqlClient, rabbitMqSender, cfg, jwtGenerator, jwtValidator)
}

func CoreNewApi(rdClient *redis.Client, sqlClient *gorm.DB, rabbitMqSender rabbitmq.RabbitMqSender, cfg *config.Config, jwtGenerator jwt_pkg.JwtGenerator, jwtValidator jwt_pkg.JwtValidator) api.Engine {
	connector := connection.NewDBConnector(rdClient, sqlClient, rabbitMqSender)
	apiEngine := api.NewEngine(&api.EnginOpt{
		App:          gin.New(),
		Cfg:          cfg,
		Connector:    connector,
		JwtGenerator: jwtGenerator,
		JwtValidator: jwtValidator,
	})
	return apiEngine
}
