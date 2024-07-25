package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/wandermaia/desafio-rate-limiter/configuration/database/mongodb"
	"github.com/wandermaia/desafio-rate-limiter/configuration/database/redisDB"
	"github.com/wandermaia/desafio-rate-limiter/internal/infra/api/handler"
	"github.com/wandermaia/desafio-rate-limiter/internal/infra/cache"
	database "github.com/wandermaia/desafio-rate-limiter/internal/infra/database/album"
	"github.com/wandermaia/desafio-rate-limiter/internal/usecase/album_usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

//desafio-rate-limiter/internal/infra/database/album

func main() {

	// Pegando o contexto
	ctx := context.Background()

	// Carregando as variáveis de ambiente
	if err := godotenv.Load("cmd/album/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	// Criando a conexão com o mongo
	mongoConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Criando a conexão com o redis
	redisConnection, err := redisDB.NewRedisConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Inicializando router
	router := gin.Default()

	//userController, bidController, auctionsController := initDependencies(databaseConnection)
	albumHandler := initDependencies(mongoConnection, redisConnection)

	// Definindo as rotas e inicializando o server
	router.POST("/login", handler.Login)
	router.POST("/album", albumHandler.CreateAlbum)
	router.GET("/album/:albumId", albumHandler.FindAlbumById)
	router.GET("/album", albumHandler.FindAllAlbums)
	router.DELETE("/album/:albumId", albumHandler.DeleteAlbumByID)
	router.Run(":8080")
}

// Função para inicializar as dependências e retornar o handler nomeado albumHandler
func initDependencies(mongoDB *mongo.Database, redisDB *redis.Client) (albumHandler *handler.AlbumHandler) {

	// Criando o repositorio para o mongo
	albumRepository := database.NewAlbumRepository(mongoDB)

	// Criando o cache utilizando o redis
	cacheRedis := cache.NewCacheRedis(redisDB)

	// Criando o Usecase
	albumHandler = handler.NewAlbumHandler(album_usecase.NewAlbumUseCase(albumRepository, cacheRedis))

	// Retorno da função. Como a variável já está nomeada na definição da função, não é necessário passar o nome aqui.
	return
}
