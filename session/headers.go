package session

import (
	"math/rand"
)

// 常见 User-Agent 列表
var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:134.0) Gecko/20100101 Firefox/134.0",
}

// 随机生成 User-Agent
func getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

// SetGeneralHeaders 设置常用的浏览器 Headers
func (sc *SessionClient) SetGeneralHeaders(overwriteHeaders map[string]string) {
	sc.SetHeader("User-Agent", getRandomUserAgent())
	sc.SetHeader(
		"Accept",
		"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	)
	sc.SetHeader("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	// sc.SetHeader("Accept-Encoding", "gzip, deflate, br, zstd")
	sc.SetHeader("Connection", "keep-alive")
	sc.SetHeader("Cache-Control", "max-age=0")

	// 覆盖默认 Headers
	for key, value := range overwriteHeaders {
		sc.SetHeader(key, value)
	}
}
