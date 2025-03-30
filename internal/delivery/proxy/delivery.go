package proxy

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/MatiXxD/go-mitm-proxy/internal/repository/proxy"
	"github.com/MatiXxD/go-mitm-proxy/pkg/env"
)

type ProxyDelivery struct {
	proxyRepo *proxy.MemProxyRepository
	cert      *tls.Certificate
	cfg       *env.Config
}

func NewProxyDelivery(pr *proxy.MemProxyRepository, cfg *env.Config) (*ProxyDelivery, error) {
	cert, err := getPrivateCert(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't get private tls certificate")
	}

	return &ProxyDelivery{
		proxyRepo: pr,
		cert:      cert,
		cfg:       cfg,
	}, nil
}

func (pd *ProxyDelivery) Handle(conn net.Conn) error {
	defer conn.Close()
	r := bufio.NewReader(conn)

	req, err := http.ReadRequest(r)
	if err != nil {
		return fmt.Errorf("error handle connection: %v", err)
	}

	if req.Method == http.MethodConnect {
		return pd.handleHTTPS(conn, req)
	}
	return pd.handleHTTP(conn, req, nil)
}

func (pd *ProxyDelivery) handleHTTPS(conn net.Conn, req *http.Request) error {
	if _, err := conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
		return fmt.Errorf("can't send CONNECT to connection: %v", err)
	}

	tlsCfg, err := pd.getTLSConfig(req.URL.Hostname())
	if err != nil {
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
			return fmt.Errorf("can't read HTTPS requset: %v", err)
		}
		err = pd.handleHTTP(tlsConn, req, tlsCfg)
		if err != nil {
			return fmt.Errorf("error handaling HTTPS request: %v", err)
		}
	}

	return nil
}

func (pd *ProxyDelivery) handleHTTP(conn net.Conn, req *http.Request, tlsCfg *tls.Config) error {
	log.Println("request info: ", req.Method, req.Host, req.RequestURI)
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
		return fmt.Errorf("can't connect to host: %v", err)
	}

	resp, err := pd.sendRequest(dial, req)
	if err != nil {
		return fmt.Errorf("can't send request: %v", err)
	}

	if err := resp.Write(conn); err != nil {
		return fmt.Errorf("can't send response to client: %v", err)
	}

	return nil
}

func (pd *ProxyDelivery) sendRequest(dial net.Conn, req *http.Request) (*http.Response, error) {
	if err := req.Write(dial); err != nil {
		return nil, fmt.Errorf("can't send request: %v", err)
	}

	resp, err := http.ReadResponse(bufio.NewReader(dial), req)
	if err != nil {
		return nil, fmt.Errorf("can't read response: %v", err)
	}

	return resp, nil
}
