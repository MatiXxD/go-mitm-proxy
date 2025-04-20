package proxy

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/MatiXxD/go-mitm-proxy/internal/usecase/request"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/MatiXxD/go-mitm-proxy/internal/repository/proxy"
	"github.com/MatiXxD/go-mitm-proxy/pkg/env"
)

type ProxyDelivery struct {
	proxyRepo      *proxy.MemProxyRepository
	requestUsecase *request.RequestUsecase
	cert           *tls.Certificate
	cfg            *env.Config
	logger         *zap.Logger
}

func NewProxyDelivery(pr *proxy.MemProxyRepository, ru *request.RequestUsecase, cfg *env.Config, logger *zap.Logger) (*ProxyDelivery, error) {
	cert, err := getPrivateCert(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't get private tls certificate")
	}

	return &ProxyDelivery{
		proxyRepo:      pr,
		requestUsecase: ru,
		cert:           cert,
		cfg:            cfg,
		logger:         logger,
	}, nil
}

func (pd *ProxyDelivery) Handle(conn net.Conn) error {
	defer conn.Close()
	r := bufio.NewReader(conn)

	req, err := http.ReadRequest(r)
	if err != nil {
		pd.logger.Error("can't read request", zap.Error(err))
		return fmt.Errorf("error handle connection: %v", err)
	}

	if req.Method == http.MethodConnect {
		return pd.handleHTTPS(conn, req)
	}
	return pd.handleHTTP(conn, req, nil)
}

func (pd *ProxyDelivery) handleHTTPS(conn net.Conn, req *http.Request) error {
	if _, err := conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
		pd.logger.Error("can't send HTTP/1.1 response", zap.Error(err))
		return fmt.Errorf("can't send CONNECT to connection: %v", err)
	}

	tlsCfg, err := pd.getTLSConfig(req.URL.Hostname())
	if err != nil {
		pd.logger.Error("can't get tls config", zap.Error(err))
		return fmt.Errorf("can't get TLS config: %v", err)
	}
	tlsConn := tls.Server(conn, tlsCfg)
	defer tlsConn.Close()

	// read request in loop, because https conn
	for {
		req, err := http.ReadRequest(bufio.NewReader(tlsConn))
		if err == io.EOF {
			break
		} else if err != nil {
			pd.logger.Error("can't read request", zap.Error(err))
			return fmt.Errorf("can't read HTTPS requset: %v", err)
		}
		err = pd.handleHTTP(tlsConn, req, tlsCfg)
		if err != nil {
			pd.logger.Error("can't handle HTTP", zap.Error(err))
			return fmt.Errorf("error handaling HTTPS request: %v", err)
		}
	}

	return nil
}

func (pd *ProxyDelivery) handleHTTP(conn net.Conn, req *http.Request, tlsCfg *tls.Config) error {
	pd.logger.Info(fmt.Sprintln("request info: ", req.Method, req.Host, req.RequestURI))
	pd.deleteHeaders(req)

	var dial net.Conn
	var err error
	if tlsCfg == nil {
		port := "80"
		if req.URL.Port() != "" {
			port = req.URL.Port()
		}
		dial, err = net.Dial("tcp", net.JoinHostPort(req.Host, port))
	} else {
		dial, err = tls.Dial("tcp", net.JoinHostPort(req.Host, "443"), tlsCfg)
	}
	if err != nil {
		pd.logger.Error("can't dial", zap.Error(err))
		return fmt.Errorf("can't connect to host: %v", err)
	}

	resp, err := pd.sendRequest(dial, req)
	if err != nil {
		pd.logger.Error("can't send request", zap.Error(err))
		return fmt.Errorf("can't send request: %v", err)
	}

	// need to save before resp.Write
	if err := pd.requestUsecase.AddRequest(req, resp); err != nil {
		pd.logger.Error("can't add request", zap.Error(err))
	}

	if err := resp.Write(conn); err != nil {
		pd.logger.Error("can't write response", zap.Error(err))
		return fmt.Errorf("can't send response to client: %v", err)
	}

	return nil
}

func (pd *ProxyDelivery) sendRequest(dial net.Conn, req *http.Request) (*http.Response, error) {
	// make body readable more than one time
	if req.Body != nil {
		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			pd.logger.Error("error reading body", zap.Error(err))
			return nil, err
		}
		req.Body = io.NopCloser(strings.NewReader(string(bytes)))
	}

	if err := req.Write(dial); err != nil {
		pd.logger.Error("can't send request", zap.Error(err))
		return nil, fmt.Errorf("can't send request: %v", err)
	}

	resp, err := http.ReadResponse(bufio.NewReader(dial), req)
	if err != nil {
		pd.logger.Error("can't read response", zap.Error(err))
		return nil, fmt.Errorf("can't read response: %v", err)
	}

	return resp, nil
}
