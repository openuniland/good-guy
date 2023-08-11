package repository

import (
	"context"
	"time"

	"github.com/openuniland/good-guy/configs"
	examschedules "github.com/openuniland/good-guy/internal/exam_schedules"
	"github.com/openuniland/good-guy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionName = "exam_schedules"

type examSchedulesRepo struct {
	mongoClient *mongo.Client
	cfg         *configs.Configs
}

func NewExamSchedulesRepository(cfg *configs.Configs, mongoClient *mongo.Client) examschedules.Repository {
	return &examSchedulesRepo{cfg: cfg, mongoClient: mongoClient}
}

func (e *examSchedulesRepo) Create(ctx context.Context, examSchedules *models.ExamSchedules) (*mongo.InsertOneResult, error) {
	dbName := e.cfg.MongoDB.MongoDBName

	coll := e.mongoClient.Database(dbName).Collection(collectionName)

	examSchedules.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	res, err := coll.InsertOne(ctx, examSchedules)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *examSchedulesRepo) Find(ctx context.Context, filter interface{}) ([]*models.ExamSchedules, error) {
	dbName := e.cfg.MongoDB.MongoDBName

	coll := e.mongoClient.Database(dbName).Collection(collectionName)

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

	coll := e.mongoClient.Database(dbName).Collection(collectionName)

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

	coll := e.mongoClient.Database(dbName).Collection(collectionName)

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
