package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Repositório para acesso ao MongoDB
type CacheRedis struct {
	client *redis.Client
}

// Função "construtora" do repositório
func NewCacheRedis(redisDB *redis.Client) *CacheRedis {

	// Cria o repositorio do mongo utilizando a collection "albums"
	return &CacheRedis{
		client: redisDB,
	}
}

// Insere um objeto no cache
func (cr *CacheRedis) Set(key string, value string, ttl time.Duration) error {
	return cr.client.Set(context.Background(), key, value, ttl).Err()
}

// Pesquisa um objeto pela chave informada
func (cr *CacheRedis) Get(key string) (string, error) {
	return cr.client.Get(context.Background(), key).Result()
}

// Delete um objeto com base na chave informada.
func (cr *CacheRedis) Delete(key string) error {
	return cr.client.Del(context.Background(), key).Err()
}
