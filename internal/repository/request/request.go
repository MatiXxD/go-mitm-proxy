package request

import (
	"context"
	"github.com/MatiXxD/go-mitm-proxy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type RequestRepository struct {
	db     *mongo.Database
	logger *zap.Logger
}

func NewRequestRepository(db *mongo.Database, logger *zap.Logger) *RequestRepository {
	return &RequestRepository{
		db:     db,
		logger: logger,
	}
}

func (rr *RequestRepository) AddRequest(parsedReq *models.ParsedRequest, parsedResp *models.ParsedResponse) error {
	requestInfo := models.NewRequestInfo(parsedReq, parsedResp)
	_, err := rr.db.Collection("request").InsertOne(context.Background(), requestInfo)
	if err != nil {
		rr.logger.Error("Failed to insert request", zap.Error(err))
		return err
	}
	return nil
}

func (rr *RequestRepository) GetRequests() ([]*models.RequestInfo, error) {
	cursor, err := rr.db.Collection("request").Find(context.Background(), bson.D{})
	if err != nil {
		rr.logger.Error("Failed to get requests", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(context.Background())

	var requests []*models.RequestInfo
	for cursor.Next(context.Background()) {
		req := models.RequestInfo{}
		if err := cursor.Decode(&req); err != nil {
			rr.logger.Error("Failed to get requests", zap.Error(err))
			return nil, err
		}
		requests = append(requests, &req)
	}

	return requests, nil
}

func (rr *RequestRepository) GetRequestById(id string) (*models.RequestInfo, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rr.logger.Error("Failed to get request by id", zap.Error(err))
		return nil, err
	}

	req := models.RequestInfo{}
	if err := rr.db.Collection("request").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&req); err != nil {
		rr.logger.Error("Failed to get request by id", zap.Error(err))
		return nil, err
	}

	return &req, nil
}
