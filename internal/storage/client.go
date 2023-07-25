package storage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "redis:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func Test() {
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)
}

func AddElement(key, value string) {
	ctx := context.Background()

	err := client.SAdd(ctx, key, value).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func RemoveElement(key, value string) {
	ctx := context.Background()

	err := client.SRem(ctx, key, value).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func GetLength(key string) int64 {
	ctx := context.Background()

	length, err := client.SCard(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}

	return length
}

func GetAll(key string) []string {
	ctx := context.Background()

	val, err := client.SMembers(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}

	return val
}

func GetRandom(key string) string {
	ctx := context.Background()

	val, err := client.SRandMember(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}

	return val
}
