package repository

import (
	"context"
	"time"

	"github.com/openuniland/good-guy/configs"
	article "github.com/openuniland/good-guy/internal/articles"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/pkg/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionName = "test_articles"

type articleRepo struct {
	mongodb *mongodb.MongoDB
	cfg     *configs.Configs
}

func NewArticleRepository(cfg *configs.Configs, mongodb *mongodb.MongoDB) article.Repository {
	return &articleRepo{cfg: cfg, mongodb: mongodb}
}

func (a *articleRepo) Create(ctx context.Context, article *models.Article) (*mongo.InsertOneResult, error) {
	dbName := a.cfg.MongoDB.MongoDBName

	coll := a.mongodb.Client.Database(dbName).Collection(collectionName)

	article.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	res, err := coll.InsertOne(ctx, article)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *articleRepo) Upsert(ctx context.Context, article *models.Article) (*mongo.UpdateResult, error) {
	dbName := a.cfg.MongoDB.MongoDBName

	coll := a.mongodb.Client.Database(dbName).Collection(collectionName)

	filter := bson.M{"aid": article.Aid}
	update := bson.M{
		"$set": bson.M{
			"aid":            article.Aid,
			"link":           article.Link,
			"title":          article.Title,
			"subscribed_ids": article.SubscribedIDs,
			"updated_at":     primitive.NewDateTimeFromTime(time.Now()),
		},
		"$setOnInsert": bson.M{
			"created_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	opts := options.Update().SetUpsert(true)

	res, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *articleRepo) FindByAid(ctx context.Context, aid int) (*models.Article, error) {
	dbName := a.cfg.MongoDB.MongoDBName

	coll := a.mongodb.Client.Database(dbName).Collection(collectionName)

	filter := bson.M{"aid": aid}
	var article *models.Article

	err := coll.FindOne(ctx, filter).Decode(article)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (a *articleRepo) Find(ctx context.Context, filter interface{}) ([]*models.Article, error) {
	dbName := a.cfg.MongoDB.MongoDBName

	coll := a.mongodb.Client.Database(dbName).Collection(collectionName)

	var articles []*models.Article

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var article *models.Article
		err := cur.Decode(&article)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (a *articleRepo) FindOne(ctx context.Context, filter interface{}) (*models.Article, error) {

	dbName := a.cfg.MongoDB.MongoDBName

	coll := a.mongodb.Client.Database(dbName).Collection(collectionName)

	var article *models.Article

	err := coll.FindOne(ctx, filter).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	// TODO: remove this
	article.SubscribedIDs = append(article.SubscribedIDs, "5611397915547076")

	return article, nil
}

func (a *articleRepo) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	dbName := a.cfg.MongoDB.MongoDBName

	coll := a.mongodb.Client.Database(dbName).Collection(collectionName)

	res, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *articleRepo) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) (*models.Article, error) {
	dbName := a.cfg.MongoDB.MongoDBName

	coll := a.mongodb.Client.Database(dbName).Collection(collectionName)

	var article *models.Article

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	updateDoc := bson.M{
		"$set": update,
	}

	err := coll.FindOneAndUpdate(ctx, filter, updateDoc, opts).Decode(&article)
	if err != nil {
		return nil, err
	}

	return article, nil
}
