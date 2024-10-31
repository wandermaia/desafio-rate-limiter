package repository

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

// Função de inicizalização do repositório
func NewRedisRepository(address, password string) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
	})

	return &RedisRepository{client: client}
}

// Função para validar o rate limite
func (r *RedisRepository) Allow(ip string, token string, maxRequests int, duration time.Duration) bool {
	ctx := context.Background()
	key := "rate_limiter:" + ip
	if token != "" {
		key = "rate_limiter:token:" + token
	}

	// Coletando o valor da chave
	count, err := r.client.Get(ctx, key).Int()

	// Caso a chave não exista, será retornoado o erro redis.Nil
	if err == redis.Nil {

		//Criando a chave nova no redis
		err = r.client.Set(ctx, key, 1, duration).Err()
		if err != nil {
			log.Printf("Erro ao criar a chave no redis: %s", err)
			return false
		}

		// Inicializando  contador de acessos com o valor 1. Será utilizado na chamada de comparação do limite mais abaixo
		count = 1

	} else if err != nil {
		log.Printf("Erro acessar o redis - %s", err)
		return false
	}

	log.Printf("Contador redis (valor): %v, máximo de requests; %v", count, maxRequests)
	// Se o limite foi alcançado
	if count > maxRequests {
		return false
	}

	// Incrementando o contador do redis, que na verdade é o valor
	_, err = r.client.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("Erro acessar o redis para incrementar o contador - %s", err)
		return false
	}

	// Retorna true se ainda não houver alcançado o limite
	return true

}