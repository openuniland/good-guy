package repository

import (
	"context"
	"time"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/cookies"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/pkg/db/mongodb"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionName = "cookies"

type cookieRepo struct {
	cfg     *configs.Configs
	mongodb *mongodb.MongoDB
}

func NewCookieRepository(cfg *configs.Configs, mongodb *mongodb.MongoDB) cookies.Repository {
	if err := mongodb.DB.CreateCollection(context.Background(), collectionName); err != nil {
		log.Warn().Err(err).Msg("collection already exists: " + collectionName)
	}
	coll := mongodb.DB.Collection(collectionName)
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.M{
				"username": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"cookies": 1,
			},
			Options: options.Index().SetUnique(false),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := coll.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create indexes: " + collectionName)
	}

	return &cookieRepo{cfg: cfg, mongodb: mongodb}
}

func (c *cookieRepo) Create(ctx context.Context, cookie *models.Cookie) (*mongo.InsertOneResult, error) {
	dbName := c.cfg.MongoDB.MongoDBName
	coll := c.mongodb.Client.Database(dbName).Collection(collectionName)

	cookie.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	return coll.InsertOne(ctx, cookie)
}

func (c *cookieRepo) FindOne(ctx context.Context, filter bson.M) (*models.Cookie, error) {
	dbName := c.cfg.MongoDB.MongoDBName
	coll := c.mongodb.Client.Database(dbName).Collection(collectionName)

	var cookie *models.Cookie
	res := coll.FindOne(ctx, filter)
	if err := res.Decode(&cookie); err != nil {
		return nil, err
	}

	return cookie, nil
}

func (c *cookieRepo) UpdateOne(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	dbName := c.cfg.MongoDB.MongoDBName
	coll := c.mongodb.Client.Database(dbName).Collection(collectionName)

	return coll.UpdateOne(ctx, filter, update)
}
