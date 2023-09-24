package service

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/jpmchia/ip2location-pfsense/backend/config"
)

var ssl_cert string
var ssl_key string
var UseSSL bool

func initTls() {

	config.Configure()
}

func loadCert(certFile string, keyFile string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return cert, err
	}
	return cert, nil
}

func loadCACert(caFile string) (*x509.CertPool, error) {
	caCert, err := tls.LoadX509KeyPair(caFile, caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert.Certificate[0])

	return caCertPool, nil
}
