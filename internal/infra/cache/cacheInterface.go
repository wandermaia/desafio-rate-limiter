package cache

import "time"

// Interface de definição do Cache
type CacheInterface interface {
	Set(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}
