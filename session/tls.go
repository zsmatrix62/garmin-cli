package session

import (
	"crypto/tls"
	"net/http"
)

// ConfigureTLS 配置自定义 TLS
func ConfigureTLS(sc *SessionClient) {
	sc.client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			MaxVersion:               tls.VersionTLS13,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
			PreferServerCipherSuites: false,
		},
	}
}
