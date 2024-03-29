package repository

import (
	"context"
	"time"

	"github.com/openuniland/good-guy/configs"
	examschedules "github.com/openuniland/good-guy/internal/exam_schedules"
	"github.com/openuniland/good-guy/internal/models"
	"github.com/openuniland/good-guy/pkg/db/mongodb"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionName = "exam_schedules"

type examSchedulesRepo struct {
	cfg     *configs.Configs
	mongodb *mongodb.MongoDB
}

func NewExamSchedulesRepository(cfg *configs.Configs, mongodb *mongodb.MongoDB) examschedules.Repository {

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
				"subjects": 1,
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
	return &examSchedulesRepo{cfg: cfg, mongodb: mongodb}
}

func (e *examSchedulesRepo) Create(ctx context.Context, examSchedules *models.ExamSchedules) (*mongo.InsertOneResult, error) {
	dbName := e.cfg.MongoDB.MongoDBName

	coll := e.mongodb.Client.Database(dbName).Collection(collectionName)

	examSchedules.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	res, err := coll.InsertOne(ctx, examSchedules)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *examSchedulesRepo) Find(ctx context.Context, filter interface{}) ([]*models.ExamSchedules, error) {
	dbName := e.cfg.MongoDB.MongoDBName

	coll := e.mongodb.Client.Database(dbName).Collection(collectionName)

	var examSchedules []*models.ExamSchedules

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var res *models.ExamSchedules
		err := cur.Decode(&res)
		if err != nil {
			return nil, err
		}

		examSchedules = append(examSchedules, res)
	}

	return examSchedules, nil
}

func (e *examSchedulesRepo) FindOne(ctx context.Context, filter interface{}) (*models.ExamSchedules, error) {

	dbName := e.cfg.MongoDB.MongoDBName

	coll := e.mongodb.Client.Database(dbName).Collection(collectionName)

	var examSchedules *models.ExamSchedules

	err := coll.FindOne(ctx, filter).Decode(&examSchedules)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return examSchedules, nil
}

func (e *examSchedulesRepo) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) (*models.ExamSchedules, error) {
	dbName := e.cfg.MongoDB.MongoDBName

	coll := e.mongodb.Client.Database(dbName).Collection(collectionName)

	var examSchedules *models.ExamSchedules

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	updateDoc := bson.M{
		"$set": update,
	}

	err := coll.FindOneAndUpdate(ctx, filter, updateDoc, opts).Decode(&examSchedules)
	if err != nil {
		return nil, err
	}

	return examSchedules, nil
}
