package redisDB

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// Constantes contendo as chaves que serão buscadas nas variáveis de ambiente
const (
	redisDsn      = "REDIS_DSN"
	redisPassword = "REDIS_PASSWORD"
)

// Cria uma nova conexão com o Redis
func NewRedisConnection(ctx context.Context) (*redis.Client, error) {

	//Inicializando o client
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv(redisDsn),
		Password: os.Getenv(redisPassword),
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Error when trying to connect to redis: %s", err)
		return nil, err
	}

	return client, nil
}
