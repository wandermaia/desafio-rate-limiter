package album_usecase

import (
	"context"
	"log"

	"github.com/wandermaia/desafio-rate-limiter/internal/entity/album_entity"
)

// Dados de entrada
type AlbumInputDTO struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Dados de saída
type AlbumOutputDTO struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Função "Construtora"
func NewAlbumUseCase(albumReporyInterface album_entity.AlbumRepositoryInterface) AlbumUseCaseInterface {
	return &AlbumUseCase{
		albumRepositoryInterface: albumReporyInterface,
	}
}

// Cadastrar um novo álbum
func (auc *AlbumUseCase) CreateNewAlbum(ctx context.Context, albumInput *AlbumInputDTO) error {

	//Entidade que receberá os dados do DTO
	album, err := album_entity.CreateAlbum(
		albumInput.Title,
		albumInput.Artist,
		albumInput.Price,
	)
	if err != nil {
		log.Printf("Erro ao executar o método CreateAlbum da entidade: %s", err)
		return err
	}

	//Inserindo os dados no banco de dados
	if err := auc.albumRepositoryInterface.CreateAlbum(ctx, album); err != nil {
		log.Printf("Erro ao executar CreateAlbum repositorio: %s", err)
		return err
	}

	return nil
}

// UseCase para o Album. Utiliza uma interface do repository para manter o desacoplamento
type AlbumUseCase struct {
	albumRepositoryInterface album_entity.AlbumRepositoryInterface
}

// Listar todos os álbuns
func (auc *AlbumUseCase) GetAllAlbums(ctx context.Context) ([]AlbumOutputDTO, error) {
	albumsEntities, err := auc.albumRepositoryInterface.FindAllAlbums(ctx)
	if err != nil {
		log.Printf("Erro ao listar todos os álbuns: %s", err)
		return nil, err
	}

	// Variável (array) para receber os dados de todos os álbuns
	var albumsOutputs []AlbumOutputDTO

	for _, albumEntity := range albumsEntities {
		albumsOutputs = append(albumsOutputs, AlbumOutputDTO{
			Id:     albumEntity.Id,
			Title:  albumEntity.Title,
			Artist: albumEntity.Artist,
			Price:  albumEntity.Price,
		})
	}

	return albumsOutputs, nil

}

// Recuperar os dados de um álbum pelo ID
func (auc *AlbumUseCase) GetAlbumByID(ctx context.Context, id string) (*AlbumOutputDTO, error) {

	// Recuperando a entidade pelo id informado
	albumEntity, err := auc.albumRepositoryInterface.FindAlbumById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Retornando os dados encontrados
	return &AlbumOutputDTO{
		Id:     albumEntity.Id,
		Title:  albumEntity.Title,
		Artist: albumEntity.Artist,
		Price:  albumEntity.Price,
	}, nil
}

// Deleta um album do ID informado
func (auc *AlbumUseCase) DeleteAlbumByID(ctx context.Context, id string) error {
	return auc.albumRepositoryInterface.DeleteAlbumByID(ctx, id)
}
