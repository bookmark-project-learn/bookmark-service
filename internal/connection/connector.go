package connection

import (
	"github.com/bookmark-project-learn/bookmark-common-libs/pkg/rabbitmq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DBConnector interface {
	GetRedisClient() *redis.Client
	GetSqlDB() *gorm.DB
	GetRabbitMqSender() rabbitmq.RabbitMqSender
}

type dbConnector struct {
	redisClient    *redis.Client
	db             *gorm.DB
	rabbitMqSender rabbitmq.RabbitMqSender
}

func NewDBConnector(redisClient *redis.Client, db *gorm.DB, rabbitMqSender rabbitmq.RabbitMqSender) DBConnector {
	return &dbConnector{
		redisClient:    redisClient,
		db:             db,
		rabbitMqSender: rabbitMqSender,
	}
}

func (db *dbConnector) GetRabbitMqSender() rabbitmq.RabbitMqSender {
	return db.rabbitMqSender
}
func (db *dbConnector) GetRedisClient() *redis.Client {
	return db.redisClient
}

func (db *dbConnector) GetSqlDB() *gorm.DB {
	return db.db
}
