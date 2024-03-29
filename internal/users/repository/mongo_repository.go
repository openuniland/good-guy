package repository

import (
	"context"
	"time"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/internal/users"
	"github.com/openuniland/good-guy/pkg/db/mongodb"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionName = "users"

type userRepo struct {
	cfg     *configs.Configs
	mongodb *mongodb.MongoDB
}

func NewUserRepository(cfg *configs.Configs, mongodb *mongodb.MongoDB) users.Repository {
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
				"subscribed_id": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := coll.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create indexes: " + collectionName)
	}

	return &userRepo{cfg: cfg, mongodb: mongodb}
}

func (u *userRepo) Create(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	dbName := u.cfg.MongoDB.MongoDBName

	coll := u.mongodb.Client.Database(dbName).Collection(collectionName)

	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	res, err := coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *userRepo) GetUsers(ctx context.Context) ([]*models.User, error) {
	dbName := u.cfg.MongoDB.MongoDBName

	coll := u.mongodb.Client.Database(dbName).Collection(collectionName)

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

func (u *userRepo) Find(ctx context.Context, filter interface{}) ([]*models.User, error) {
	dbName := u.cfg.MongoDB.MongoDBName

	coll := u.mongodb.Client.Database(dbName).Collection(collectionName)

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

func (u *userRepo) FindOneUserByCondition(ctx context.Context, filter interface{}) (*models.User, error) {
	dbName := u.cfg.MongoDB.MongoDBName

	coll := u.mongodb.Client.Database(dbName).Collection(collectionName)

	singleResult := coll.FindOne(ctx, filter)

	var user *models.User
	singleResult.Decode(&user)

	return user, nil
}

func (u *userRepo) FindOneAndUpdate(ctx context.Context, filter, update bson.M) (*models.User, error) {
	dbName := u.cfg.MongoDB.MongoDBName

	coll := u.mongodb.Client.Database(dbName).Collection(collectionName)

	var user *models.User

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	update["updated_at"] = primitive.NewDateTimeFromTime(time.Now())
	updateDoc := bson.M{
		"$set": update,
	}

	err := coll.FindOneAndUpdate(ctx, filter, updateDoc, opts).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepo) FindOneAndDelete(ctx context.Context, filter interface{}) *mongo.SingleResult {
	dbName := u.cfg.MongoDB.MongoDBName

	coll := u.mongodb.Client.Database(dbName).Collection(collectionName)

	res := coll.FindOneAndDelete(ctx, filter)

	return res
}
