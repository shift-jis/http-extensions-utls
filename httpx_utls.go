package httpx

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	oohttp "github.com/ooni/oohttp"
)

var defaultDialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}

type TLSConnFactory func(conn net.Conn, config *tls.Config) oohttp.TLSConn

type SecureClientConfig struct {
	TLSFactory *SecureTLSConnFactory
	ProxyURL   *url.URL

	ForceAttemptHTTP2  bool
	MaxIdleConnections int
	RequestTimeout     time.Duration
}

func DefaultSecureClientConfig() *SecureClientConfig {
	return &SecureClientConfig{
		TLSFactory: &SecureTLSConnFactory{},
		ProxyURL:   nil,

		ForceAttemptHTTP2:  true,
		MaxIdleConnections: 100,
		RequestTimeout:     time.Second * 20,
	}
}

func NewSecureHTTPClient(config *SecureClientConfig) (*http.Client, error) {
	if config == nil {
		config = DefaultSecureClientConfig()
	}
	return newHTTPClientWithTransport(NewSecureHTTPTransport(*config), *config)
}

func newHTTPClientWithTransport(transport *oohttp.Transport, config SecureClientConfig) (*http.Client, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: &oohttp.StdlibTransport{
			Transport: transport,
		},
		Timeout:       config.RequestTimeout,
		Jar:           cookieJar,
		CheckRedirect: nil,
	}, nil
}

func NewSecureHTTPTransport(config SecureClientConfig) *oohttp.Transport {
	return &oohttp.Transport{
		Proxy: func(httpRequest *oohttp.Request) (*url.URL, error) {
			if config.ProxyURL != nil {
				return config.ProxyURL, nil
			}
			return oohttp.ProxyFromEnvironment(httpRequest)
		},
		DialContext:           defaultDialer.DialContext,
		TLSClientFactory:      config.TLSFactory.CreateTLSConnection,
		ForceAttemptHTTP2:     config.ForceAttemptHTTP2,
		MaxIdleConns:          config.MaxIdleConnections,
		TLSHandshakeTimeout:   10 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
}
