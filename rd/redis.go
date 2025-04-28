package rd

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	ctx := context.Background()
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Printf("Connected to Redis: %s", pong)
}

// Publish message to Redis
func Publish(channel string, payload interface{}) error {
	return RedisClient.Publish(ctx, channel, payload).Err()
}

func Subscribe(channel string) <-chan *redis.Message {
	pubsub := RedisClient.Subscribe(ctx, channel)
	return pubsub.Channel()
}
