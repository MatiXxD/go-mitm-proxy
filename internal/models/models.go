package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ParsedRequest struct {
	Method        string         `bson:"method"`
	URL           string         `bson:"url"`
	Form          url.Values     `bson:"queryParams"`
	Header        http.Header    `bson:"headers"`
	Cookies       []*http.Cookie `bson:"cookies"`
	Body          string         `bson:"body"`
	ContentLength int64          `bson:"contentLength"`
	PostForm      url.Values     `bson:"postForm"`
}

func NewParsedRequest(r *http.Request) (*ParsedRequest, error) {
	url := r.URL.String()
	if r.URL.Hostname() == "" {
		url = "https://" + r.Host + r.URL.Path
	}

	body := ""
	if r.Body != nil {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		body = string(bytes)
		r.Body = io.NopCloser(strings.NewReader(body))
	}

	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	return &ParsedRequest{
		Method:        r.Method,
		URL:           url,
		Form:          r.Form,
		Header:        r.Header,
		Cookies:       r.Cookies(),
		Body:          body,
		ContentLength: r.ContentLength,
		PostForm:      r.PostForm,
	}, nil
}

type ParsedResponse struct {
	Status        string         `bson:"status"`
	StatusCode    int            `bson:"statusCode"`
	Header        http.Header    `bson:"headers"`
	Cookies       []*http.Cookie `bson:"cookies"`
	Body          string         `bson:"body"`
	ContentLength int64          `bson:"contentLength"`
}

func NewParsedResponse(r *http.Response) (*ParsedResponse, error) {
	body := ""
	if r.Body != nil {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		body = string(bytes)
		r.Body = io.NopCloser(strings.NewReader(body))
	}

	return &ParsedResponse{
		Status:        r.Status,
		StatusCode:    r.StatusCode,
		Header:        r.Header,
		Cookies:       r.Cookies(),
		Body:          body,
		ContentLength: r.ContentLength,
	}, nil
}

type RequestInfo struct {
	Request   *ParsedRequest  `bson:"request"`
	Response  *ParsedResponse `bson:"response"`
	CreatedAt time.Time       `bson:"createdAt"`
}

type RequestInfoWithID struct {
	ID        primitive.ObjectID `bson:"_id"`
	Request   *ParsedRequest     `bson:"request"`
	Response  *ParsedResponse    `bson:"response"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func NewRequestInfo(req *ParsedRequest, resp *ParsedResponse) *RequestInfo {
	return &RequestInfo{
		Request:   req,
		Response:  resp,
		CreatedAt: time.Now(),
	}
}
