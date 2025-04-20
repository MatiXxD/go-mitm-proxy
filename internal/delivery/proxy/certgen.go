package proxy

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/MatiXxD/go-mitm-proxy/pkg/env"
)

const (
	pemCert = "certPEM"
	pemKey  = "keyPEM"
)

func (pd *ProxyDelivery) getTLSConfig(host string) (*tls.Config, error) {
	cert, err := pd.getTLSCert(host)
	if err != nil {
		return nil, fmt.Errorf("can't get tls cert: %v", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{*cert},
		ServerName:   host,
	}, nil
}

func (pd *ProxyDelivery) getTLSCert(host string) (*tls.Certificate, error) {
	certData, ok := pd.proxyRepo.Get(host)
	if !ok {
		genCertData, err := pd.generateCert(host)
		if err != nil {
			return nil, fmt.Errorf("can't generate cert for host %s: %v", host, err)
		}
		pd.proxyRepo.Store(host, genCertData)
		certData = genCertData
	}
	tlsCert, err := tls.X509KeyPair(certData[pemCert], certData[pemKey])
	if err != nil {
		return nil, fmt.Errorf("can't get tls certificate for host %s: %v", host, err)
	}
	return &tlsCert, nil
}

func (pd *ProxyDelivery) generateCert(host string) (map[string][]byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("can't generate private key for host %s: %v", host, err)
	}

	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1e9))
	if err != nil {
		return nil, fmt.Errorf("can't generate serial number for host %s: %v", host, err)
	}

	certTmpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   host,
			Organization: []string{"Solist"},
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(host); ip != nil {
		certTmpl.IPAddresses = append(certTmpl.IPAddresses, ip)
	} else {
		certTmpl.DNSNames = append(certTmpl.DNSNames, host)
	}

	privateCert, err := x509.ParseCertificate(pd.cert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("can't parse private cert: %v", err)
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, privateCert, &key.PublicKey, pd.cert.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("can't create private cert for host %s: %v", host, err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("can't marshal private key")
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	})

	certData := map[string][]byte{
		pemCert: certPEM,
		pemKey:  keyPEM,
	}
	return certData, nil
}

func getPrivateCert(cfg *env.Config) (*tls.Certificate, error) {
	keyPEM, err := os.ReadFile(cfg.ProxyConfig.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("can't read ca key: %v", err)
	}
	keyDecoded, _ := pem.Decode(keyPEM)
	key, err := x509.ParsePKCS8PrivateKey(keyDecoded.Bytes)
	if err != nil {
		return nil, fmt.Errorf("can't parse ca key: %v", err)
	}

	certPEM, err := os.ReadFile(cfg.ProxyConfig.CertPath)
	if err != nil {
		return nil, fmt.Errorf("can't read ca cert: %v", err)
	}
	certDecoded, _ := pem.Decode(certPEM)
	cert, err := x509.ParseCertificate(certDecoded.Bytes)
	if err != nil {
		return nil, fmt.Errorf("can't parse ca cert: %v", err)
	}

	return &tls.Certificate{Certificate: [][]byte{cert.Raw}, PrivateKey: key}, nil
}
