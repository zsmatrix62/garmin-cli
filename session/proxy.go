package session

import (
	"net/http"
	"net/url"
)

// SetProxy 设置代理
func SetProxy(sc *SessionClient, proxyURL string) {
	proxy, _ := url.Parse(proxyURL)
	sc.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
}
