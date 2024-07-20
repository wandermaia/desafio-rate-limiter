package database

import (
	"context"
	"log"

	"github.com/wandermaia/desafio-rate-limiter/internal/entity/album_entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// album represents data about a record album.
type AlbumEntityMongo struct {
	Id     string  `bson:"_id"`
	Title  string  `bson:"title"`
	Artist string  `bson:"artist"`
	Price  float64 `bson:"price"`
}

// Repositório para acesso ao MongoDB
type AlbumRepositoryMongo struct {
	Collection *mongo.Collection
}

// Função "construtora" do repositório
func NewAlbumRepository(database *mongo.Database) *AlbumRepositoryMongo {

	// Cria o repositorio do mongo utilizando a collection "albums"
	return &AlbumRepositoryMongo{
		Collection: database.Collection("albums"),
	}
}

// Insere um novo album no mongoDB
func (ar AlbumRepositoryMongo) CreateAlbum(ctx context.Context, albumEntity *album_entity.Album) error {

	// Criando a entidade com os dados informados
	albumEntityMongo := &AlbumEntityMongo{
		Id:     albumEntity.Id,
		Title:  albumEntity.Title,
		Artist: albumEntity.Artist,
		Price:  albumEntity.Price,
	}

	// Inserindo os dados no mongo
	_, err := ar.Collection.InsertOne(ctx, albumEntityMongo)
	if err != nil {
		log.Printf("Erro ao tentar inserir no mongodb: %s", err)
		return err
	}

	return nil
}

// Remove o album do ID informado
func (ar AlbumRepositoryMongo) DeleteAlbumByID(ctx context.Context, id string) error {

	// Filtro pelo ID para realizar a pesquisa no Mongo
	filter := bson.M{"_id": id}

	if _, err := ar.Collection.DeleteOne(ctx, filter); err != nil {
		log.Printf("Erro ao deletar o objeto ID %s no mongodb: %s", id, err)
		return err
	}

	return nil
}
