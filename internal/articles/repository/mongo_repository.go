package repository

import (
	"context"
	"time"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/articles"
	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionName = "articles"

type articleRepo struct {
	mongoClient *mongo.Client
	cfg         *configs.Configs
}

func NewUserRepository(mongoClient *mongo.Client) articles.Repository {
	return &articleRepo{mongoClient: mongoClient}
}

func (a *articleRepo) Create(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	dbName := a.cfg.MongoDB.MongoDBName

	coll := a.mongoClient.Database(dbName).Collection(collectionName)

	user.CreatedAt = time.Now().Format(time.RFC3339)
	res, err := coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *articleRepo) GetUsers(ctx context.Context) ([]*models.User, error) {
	dbName := a.cfg.MongoDB.MongoDBName

	coll := a.mongoClient.Database(dbName).Collection(collectionName)

	filter := bson.D{{
		Key: "is_deleted", Value: false,
	}}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []*models.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
