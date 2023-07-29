package storage

import (
	"context"
	"encoding/json"
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func InitClient(redisUrl string) {
	client = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Test() error {
	return client.Ping(context.Background()).Err()
}

func AddElement(key models.MessageMatch, value string) error {
	ctx := context.Background()

	marshal, err := json.Marshal(key)
	if err != nil {
		return err
	}

	return client.SAdd(ctx, string(marshal), value).Err()
}

func RemoveElement(key models.MessageMatch, value string) error {
	ctx := context.Background()

	marshal, err := json.Marshal(key)
	if err != nil {
		return err
	}

	return client.SRem(ctx, string(marshal), value).Err()
}

func GetAll(key models.MessageMatch) ([]string, error) {
	ctx := context.Background()

	marshal, err := json.Marshal(key)
	if err != nil {
		return nil, err
	}

	return client.SMembers(ctx, string(marshal)).Result()
}

func GetAllKeys() ([]models.MessageMatch, error) {
	ctx := context.Background()

	result, err := client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	var allMessageMatchers []models.MessageMatch
	for _, s := range result {
		var messageMatch models.MessageMatch
		err := json.Unmarshal([]byte(s), &messageMatch)
		if err != nil {
			return nil, err
		}
		allMessageMatchers = append(allMessageMatchers, messageMatch)
	}

	return allMessageMatchers, err
}
