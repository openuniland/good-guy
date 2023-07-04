package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/openuniland/good-guy/configs"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBClient(cfg *configs.Configs) (*mongo.Client, error) {
	protocol := cfg.MongoDB.MongoDBProtocol
	username := cfg.MongoDB.MongoDBUsername
	password := cfg.MongoDB.MongoDBPassword
	host := cfg.MongoDB.MongoDBHost
	database := cfg.MongoDB.MongoDBName
	replicaSet := cfg.MongoDB.MongoDBReplicaSet

	uri := fmt.Sprintf("%s://%s:%s@%s/%s?replicaSet=%s", protocol, username, password, host, database, replicaSet)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("connected to mongodb")
	return client, nil
}
