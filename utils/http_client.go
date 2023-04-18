package utils

import (
	"crypto/tls"
	"net/http"
	"time"
)

var DefaultTimeout = 30 * time.Second

func DefaultHttpClient() *http.Client {
	return DefaultHttpClientWithTimeout(DefaultTimeout)
}

func DefaultHttpClientWithTimeout(timeout time.Duration) *http.Client {
	client := http.Client{
		Timeout: timeout,
	}
	return &client
}

func SkipVerifyClient(timeout time.Duration, verifyServerCertificate bool) *http.Client {
	client := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: verifyServerCertificate,
			},
			Proxy: http.ProxyFromEnvironment,
		},
	}
	return &client
}
