package request

import (
	"fmt"
	"github.com/MatiXxD/go-mitm-proxy/internal/models"
	"github.com/MatiXxD/go-mitm-proxy/internal/repository/request"
	"go.uber.org/zap"
	"net/http"
)

type RequestUsecase struct {
	repo   *request.RequestRepository
	logger *zap.Logger
}

func NewRequestUsecase(repo *request.RequestRepository, logger *zap.Logger) *RequestUsecase {
	return &RequestUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (ru *RequestUsecase) AddRequest(req *http.Request, resp *http.Response) error {
	parsedReq, err := models.NewParsedRequest(req)
	if err != nil {
		ru.logger.Error("failed to parse request", zap.Error(err))
		return fmt.Errorf("can't add request to db")
	}

	parsedResp, err := models.NewParsedResponse(resp)
	if err != nil {
		ru.logger.Error("failed to parse response", zap.Error(err))
		return fmt.Errorf("can't add request to db")
	}

	if err := ru.repo.AddRequest(parsedReq, parsedResp); err != nil {
		ru.logger.Error("can't add request", zap.Error(err))
		return fmt.Errorf("can't add request to db")
	}

	return nil
}

func (ru *RequestUsecase) GetRequestsInfo() ([]*models.RequestInfoWithID, error) {
	reqs, err := ru.repo.GetRequests()
	if err != nil {
		ru.logger.Error("failed to get requests", zap.Error(err))
		return nil, fmt.Errorf("failed to get requests from db")
	}
	return reqs, nil
}

func (ru *RequestUsecase) GetRequestById(id string) (*models.RequestInfo, error) {
	req, err := ru.repo.GetRequestById(id)
	if err != nil {
		ru.logger.Error("failed to get request", zap.Error(err))
		return nil, fmt.Errorf("failed to get request from db")
	}
	return req, nil
}
