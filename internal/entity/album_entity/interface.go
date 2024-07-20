package album_entity

import "context"

// Interface com as definições de criação do repositório
type AlbumRepositoryInterface interface {
	CreateAlbum(ctx context.Context, albumEntity *Album) error
	FindAllAlbums(ctx context.Context) ([]Album, error)
	FindAlbumById(ctx context.Context, id string) (*Album, error)
	DeleteAlbumByID(ctx context.Context, id string) error
}
