package utils

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		fmt.Println("Creating new redis client")
		newRedisClient := redis.NewClient(&redis.Options{Addr: os.Getenv("AUTH_REDIS_URI")})
		redisClient = newRedisClient
	}

	return redisClient
}

func GetDb() (*gorm.DB, error) {
	if db == nil {
		sqlDB, err := sql.Open("pgx", os.Getenv("AUTH_DB_URI"))
		if err != nil {
			return nil, err
		}

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
		if err != nil {
			return nil, err
		}

		newRedisClient := GetRedisClient()
		cache, err := cache.NewGorm2Cache(&config.CacheConfig{
			CacheLevel:           config.CacheLevelAll,
			CacheStorage:         config.CacheStorageRedis,
			RedisConfig:          cache.NewRedisConfigWithClient(newRedisClient),
			InvalidateWhenUpdate: true,   // when you create/update/delete objects, invalidate cache
			CacheTTL:             100000, // 100s
			CacheMaxItemCnt:      20,     // if length of objects retrieved one single time exceeds this number, then don't cache
		})
		if err != nil {
			fmt.Println("Error creating caching layer: ", err)
			return nil, err
		}

		err = gormDB.Use(cache)
		if err != nil {
			fmt.Println("Error using caching layer: ", err)
			return nil, err
		}

		db = gormDB
	}

	return db, nil
}
