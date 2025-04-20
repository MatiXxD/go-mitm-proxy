package request

import (
	"github.com/MatiXxD/go-mitm-proxy/internal/usecase/request"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type RequestDelivery struct {
	usecase *request.RequestUsecase
	logger  *zap.Logger
}

func NewRequestDelivery(usecase *request.RequestUsecase, logger *zap.Logger) *RequestDelivery {
	return &RequestDelivery{
		usecase: usecase,
		logger:  logger,
	}
}

func (rd *RequestDelivery) GetRequestsInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		reqsInfo, err := rd.usecase.GetRequestsInfo()
		if err != nil {
			rd.logger.Error("GetRequestsInfo: ", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "could not retrieve requests",
			})
		}
		return c.JSON(http.StatusOK, reqsInfo)
	}
}

func (rd *RequestDelivery) GetRequestById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "id is required",
			})
		}

		reqInfo, err := rd.usecase.GetRequestById(id)
		if err != nil {
			rd.logger.Error("GetRequestById: ", zap.Error(err))
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "request not found",
			})
		}

		return c.JSON(http.StatusOK, reqInfo)
	}
}
