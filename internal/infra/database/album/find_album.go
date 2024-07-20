package database

import (
	"context"
	"log"

	"github.com/wandermaia/desafio-rate-limiter/internal/entity/album_entity"
	"go.mongodb.org/mongo-driver/bson"
)

// Retorna todos os álbuns existentes
func (ar AlbumRepositoryMongo) FindAllAlbums(ctx context.Context) ([]album_entity.Album, error) {

	// Filtro vazio para realizar a pesquisa no MongoDB de todos os albums
	filter := bson.M{}

	// Cursor para recuperar todos os albuns
	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Erro ao realizar a busca no MongoDB: %s", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	// Variável para receber os dados coletados
	var albumsMongo []AlbumEntityMongo
	if err := cursor.All(ctx, &albumsMongo); err != nil {
		log.Printf("Erro ao realizar o decode dos dados do mongo %s", err)
		return nil, err
	}

	// Transformar os dados nas entidades
	var albumsEntity []album_entity.Album

	for _, album := range albumsMongo {
		albumsEntity = append(albumsEntity, album_entity.Album{
			Id:     album.Id,
			Title:  album.Title,
			Artist: album.Artist,
			Price:  album.Price,
		})
	}

	return albumsEntity, nil
}

// Realiza a busca no MongoDB pelo ID informado
func (ar AlbumRepositoryMongo) FindAlbumById(ctx context.Context, id string) (*album_entity.Album, error) {

	// Filtro pelo ID para realizar a pesquisa no Mongo
	filter := bson.M{"_id": id}

	// Realizando a pesquisa do dados
	var albumEntityMongo AlbumEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&albumEntityMongo); err != nil {
		log.Printf("Erro ao realizar a busca no mongo pelo ID: %s", err)
		return nil, err
	}

	return &album_entity.Album{
		Id:     albumEntityMongo.Id,
		Title:  albumEntityMongo.Title,
		Artist: albumEntityMongo.Artist,
		Price:  albumEntityMongo.Price,
	}, nil

}
