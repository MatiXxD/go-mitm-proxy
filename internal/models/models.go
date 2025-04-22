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
	Host          string         `bson:"host"`
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
		Host:          r.Host,
		Form:          r.Form,
		Header:        r.Header,
		Cookies:       r.Cookies(),
		Body:          body,
		ContentLength: r.ContentLength,
		PostForm:      r.PostForm,
	}, nil
}

func CloneParsedRequest(original *ParsedRequest) *ParsedRequest {
	if original == nil {
		return nil
	}

	cloneValues := func(values url.Values) url.Values {
		copy := make(url.Values)
		for k, v := range values {
			copy[k] = append([]string(nil), v...)
		}
		return copy
	}

	cloneHeaders := func(h http.Header) http.Header {
		copy := make(http.Header)
		for k, v := range h {
			copy[k] = append([]string(nil), v...)
		}
		return copy
	}

	cloneCookies := func(cookies []*http.Cookie) []*http.Cookie {
		copy := make([]*http.Cookie, len(cookies))
		for i, c := range cookies {
			if c != nil {
				cCopy := *c
				copy[i] = &cCopy
			}
		}
		return copy
	}

	return &ParsedRequest{
		Method:        original.Method,
		URL:           original.URL,
		Host:          original.Host,
		Form:          cloneValues(original.Form),
		Header:        cloneHeaders(original.Header),
		Cookies:       cloneCookies(original.Cookies),
		Body:          original.Body,
		ContentLength: original.ContentLength,
		PostForm:      cloneValues(original.PostForm),
	}
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

func NewRequestInfo(req *ParsedRequest, resp *ParsedResponse) *RequestInfo {
	return &RequestInfo{
		Request:   req,
		Response:  resp,
		CreatedAt: time.Now(),
	}
}

type RequestInfoWithID struct {
	ID        primitive.ObjectID `bson:"_id"`
	Request   *ParsedRequest     `bson:"request"`
	Response  *ParsedResponse    `bson:"response"`
	CreatedAt time.Time          `bson:"createdAt"`
}
