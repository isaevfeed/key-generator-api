package cache

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
)

type Cache struct {
	Addr     string
	Password string
	Database int
	Client   *redis.Client
}

var ctx = context.Background()

func NewCache(Addr, Password string, Database int) *Cache {
	Client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       Database,
	})

	return &Cache{Addr, Password, Database, Client}
}

func (c *Cache) Set(key, value string) error {
	ttl, _ := os.LookupEnv("REDIS_TTL")
	duration, err := strconv.Atoi(ttl)
	if err != nil {
		return err
	}

	c.Client.Set(ctx, key, value, time.Duration(duration))

	return nil
}

func (c *Cache) Get(key string) string {
	cacheValue := c.Client.Get(ctx, key)

	if cacheValue.Err() != nil {
		log.Print("Cache value get error")
		return ""
	}

	return cacheValue.Val()
}
