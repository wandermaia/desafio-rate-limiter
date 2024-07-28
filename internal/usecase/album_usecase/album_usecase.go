package album_usecase

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/wandermaia/desafio-rate-limiter/internal/entity/album_entity"
	"github.com/wandermaia/desafio-rate-limiter/internal/infra/cache"
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

const (
	AlbumPrefixKeyCache = "album_"
	CacheTTL            = "ALBUM_CACHE_TTL_SECONDS" // Nome da chave da variável de ambiente
)

// Password: os.Getenv(ALBUM_CACHE_TTL_SECONDS),
// UseCase para o Album. Utiliza uma interface do repository para manter o desacoplamento
type AlbumUseCase struct {
	albumRepositoryInterface album_entity.AlbumRepositoryInterface
	albumCacheInterface      cache.CacheInterface
}

// Função "Construtora"
func NewAlbumUseCase(albumRepoInterface album_entity.AlbumRepositoryInterface, cacheInterface cache.CacheInterface) AlbumUseCaseInterface {
	return &AlbumUseCase{
		albumRepositoryInterface: albumRepoInterface,
		albumCacheInterface:      cacheInterface,
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

	//Salvando o album no cache
	err = auc.saveAlbumCache(album)
	if err != nil {
		log.Printf("Error writing Album to cache: %s", err)
	}

	return nil
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

	// Verificando se o album existe no cache
	// se não existir, pesquisar no banco e salva no cache
	albumEntity, err := auc.searchAlbumCache(id)
	if err != nil {

		//Registrando a busca no cache
		log.Printf("Error retrieving album from cache: %s", err)

		// Buscando os dados no repositorio
		albumEntity, err = auc.albumRepositoryInterface.FindAlbumById(ctx, id)
		if err != nil {
			log.Printf("Error retrieving album from repository: %s", err)
			return nil, err
		}

		//Salvando o album no cache
		err = auc.saveAlbumCache(albumEntity)
		if err != nil {
			log.Printf("Error writing Album to cache: %s", err)
		}
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

	//Salvando o album no cache
	err := auc.deleteAlbumCache(id)
	if err != nil {
		log.Printf("Error when deleting album from cache: %s", err)
	}

	// Deletando o álbum do repositório
	return auc.albumRepositoryInterface.DeleteAlbumByID(ctx, id)
}

// Pesquisa do album no cache
func (auc *AlbumUseCase) searchAlbumCache(id string) (*album_entity.Album, error) {

	// Entidade para receber os dados, caso esteja no cache
	var albumEntity album_entity.Album

	// Recuperação dos dados no cache
	key := AlbumPrefixKeyCache + id
	albumJson, err := auc.albumCacheInterface.Get(key)
	if err != nil {
		return &album_entity.Album{}, err
	}

	// Realizando o unmarshall e retornando a entidade
	err = json.Unmarshal([]byte(albumJson), &albumEntity)
	return &albumEntity, err
}

// Inclusão de album no cache
func (auc *AlbumUseCase) saveAlbumCache(album *album_entity.Album) error {

	// Criando a chave
	key := AlbumPrefixKeyCache + album.Id
	AlbumJson, err := json.Marshal(album)
	if err != nil {
		return err
	}

	//Gravando no cache com tempo de vida com base na variável de ambiente definda em CacheTTL
	// Caso haja algum erro, o valor 60 será definido como padrão.
	ttl, err := strconv.Atoi(os.Getenv(CacheTTL))
	log.Printf("ttl cache value: %d", ttl)
	if err != nil {
		log.Printf("Error convert ttl string to int - %s", err)
		ttl = 60
	}

	//return auc.albumCacheInterface.Set(key, string(AlbumJson), time.Second*60)
	return auc.albumCacheInterface.Set(key, string(AlbumJson), time.Second*time.Duration(ttl))
}

// Deletando o objeto do cache com base no ID
func (auc *AlbumUseCase) deleteAlbumCache(id string) error {

	// Criando a chave de buscas pelo id e deletando do cache
	key := AlbumPrefixKeyCache + id
	return auc.albumCacheInterface.Delete(key)
}
