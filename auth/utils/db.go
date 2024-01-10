package utils

import (
	"internal/helpers"
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var db *gorm.DB
var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		redisClient = helpers.GetRedisClient(os.Getenv("AUTH_REDIS_URI"))
	}

	return redisClient
}

func GetDb() (*gorm.DB, error) {
	if db == nil {
		gormDB, err := helpers.GetDb(helpers.GetDbProps{
			DbUri:       os.Getenv("AUTH_DB_URI"),
			RedisUri:    os.Getenv("AUTH_REDIS_URI"),
			RedisClient: redisClient,
		})
		if err != nil {
			return nil, err
		}
		db = gormDB
	}
	return db, nil
}
