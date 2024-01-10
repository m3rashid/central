package helpers

import (
	"database/sql"
	"fmt"

	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetRedisClient(redisUri string) *redis.Client {
	fmt.Println("Creating new redis client")
	return redis.NewClient(&redis.Options{Addr: redisUri})
}

type GetDbProps struct {
	DbUri       string
	RedisUri    string
	RedisClient *redis.Client
}

func GetDb(props GetDbProps) (*gorm.DB, error) {
	sqlDB, err := sql.Open("pgx", props.DbUri)
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	newRedisClient := props.RedisClient
	if newRedisClient == nil {
		newRedisClient = GetRedisClient(props.RedisUri)
	}

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

	return gormDB, nil
}
