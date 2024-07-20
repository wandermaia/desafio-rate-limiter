package album_usecase

import (
	"context"
)

// Interface de definição do AlbumUseCase
type AlbumUseCaseInterface interface {
	GetAllAlbums(ctx context.Context) ([]AlbumOutputDTO, error)
	GetAlbumByID(ctx context.Context, id string) (*AlbumOutputDTO, error)
	CreateNewAlbum(ctx context.Context, albumInput *AlbumInputDTO) error
	DeleteAlbumByID(ctx context.Context, id string) error
}
