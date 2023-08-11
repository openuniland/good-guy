package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/openuniland/good-guy/configs"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var mongoDB *MongoDB
var mongoDBOnce sync.Once

func NewMongoDB(cfg *configs.Configs) (*MongoDB, error) {
	mongoDBOnce.Do(func() {
		mongoDB = &MongoDB{}
		err := mongoDB.connect(cfg)
		if err != nil {
			log.Error().Err(err).Msg("failed to connect to mongodb")
			return
		}

		mongoDB.DB = mongoDB.GetDB(cfg.MongoDB.MongoDBName)

	})
	return mongoDB, nil
}

func (m *MongoDB) connect(cfg *configs.Configs) error {
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
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	m.Client = client

	log.Info().Msg("connected to mongodb")
	return nil
}

func (m *MongoDB) GetClient() *mongo.Client {
	return m.Client
}

func (m *MongoDB) GetDB(dbname string) *mongo.Database {
	return m.Client.Database(dbname)
}

func (m *MongoDB) Close() {
	if m.Client != nil {
		m.Client.Disconnect(context.Background())
	}
}
