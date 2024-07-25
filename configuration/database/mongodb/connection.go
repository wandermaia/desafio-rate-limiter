package mongodb

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Constantes contendo as chaves que serão buscadas nas variáveis de ambiente
const (
	mongoDbUrl = "MONGODB_URL"
	mongoDb    = "MONGODB_DB"
)

// Cria uma nova conexão ao mongo com base nas configurações de URL
func NewMongoDBConnection(ctx context.Context) (*mongo.Database, error) {
	mongoURL := os.Getenv(mongoDbUrl)
	mongoDatabase := os.Getenv(mongoDb)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Printf("Erro ao tentar conectar no mongodb database: %s", err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("Erro ao tentar pingar o mongodb database: %s", err)
		return nil, err
	}

	return client.Database(mongoDatabase), nil
}
