package request

import (
	"bytes"
	"crypto/tls"
	"github.com/MatiXxD/go-mitm-proxy/internal/models"
	"go.uber.org/zap"
	"net/http"
)

func (rd *RequestDelivery) sendRequest(reqInfo *models.RequestInfo) error {
	req, err := newRequest(reqInfo.Request)
	if err != nil {
		rd.logger.Error("can't create request", zap.Error(err))
		return err
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		rd.logger.Error("can't send request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if err := rd.usecase.AddRequest(req, resp); err != nil {
		rd.logger.Error("can't add request", zap.Error(err))
		return err
	}

	return nil
}

func newRequest(parsedReq *models.ParsedRequest) (*http.Request, error) {
	req, err := http.NewRequest(parsedReq.Method, parsedReq.URL, bytes.NewBuffer([]byte(parsedReq.Body)))
	if err != nil {
		return nil, err
	}

	req.Host = parsedReq.URL
	req.ContentLength = parsedReq.ContentLength

	req.Header = parsedReq.Header
	for _, cookie := range parsedReq.Cookies {
		req.AddCookie(cookie)
	}

	req.PostForm = parsedReq.PostForm
	req.Form = parsedReq.Form

	return req, nil
}
