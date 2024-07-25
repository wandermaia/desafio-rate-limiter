package album_entity

import (
	"context"
	"time"
)

// Interface com as definições de criação do repositório
type AlbumRepositoryInterface interface {
	CreateAlbum(ctx context.Context, albumEntity *Album) error
	FindAllAlbums(ctx context.Context) ([]Album, error)
	FindAlbumById(ctx context.Context, id string) (*Album, error)
	DeleteAlbumByID(ctx context.Context, id string) error
}

// Interface de definição do Cache
type CacheInterface interface {
	Set(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}
