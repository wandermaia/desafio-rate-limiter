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
	auth_midlleware "github.com/wandermaia/desafio-rate-limiter/internal/middleware/auth_middleware"
	"github.com/wandermaia/desafio-rate-limiter/internal/usecase/album_usecase"
	"github.com/wandermaia/desafio-rate-limiter/internal/usecase/authentication_usecase"
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
	accessHandler, albumHandler := initDependencies(mongoConnection, redisConnection)

	// Definindo as rotas e inicializando o server
	router.POST("/token", accessHandler.CreateAccessToken)
	router.POST("/token/refresh", accessHandler.RefreshAccessToken)
	router.POST("/logout", accessHandler.Logout)
	router.GET("/health", accessHandler.Health)
	router.POST("/album", auth_midlleware.TokenAuthMiddleware(), albumHandler.CreateAlbum)
	router.GET("/album/:albumId", auth_midlleware.TokenAuthMiddleware(), albumHandler.FindAlbumById)
	router.GET("/album", auth_midlleware.TokenAuthMiddleware(), albumHandler.FindAllAlbums)
	router.DELETE("/album/:albumId", auth_midlleware.TokenAuthMiddleware(), albumHandler.DeleteAlbumByID)
	router.Run(":8080")
}

// Função para inicializar as dependências e retornar o handler nomeado albumHandler
func initDependencies(mongoDB *mongo.Database, redisDB *redis.Client) (accessHandler *handler.AccessHandler, albumHandler *handler.AlbumHandler) {

	// Criando o repositorio para o mongo
	albumRepository := database.NewAlbumRepository(mongoDB)

	// Criando o cache utilizando o redis
	cacheRedis := cache.NewCacheRedis(redisDB)

	// Criando os Usecases
	albumHandler = handler.NewAlbumHandler(album_usecase.NewAlbumUseCase(albumRepository, cacheRedis))
	accessHandler = handler.NewAccessHandler(authentication_usecase.NewAuthenticationUseCase(cacheRedis))

	// Retorno da função. Como as variáveis já estão nomeadas na definição da função, não é necessário passar o nome aqui.
	return
}
