package session

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/antchfx/htmlquery"
)

type SessionClient struct {
	client  *http.Client
	headers map[string]string
	Jar     *cookiejar.Jar
}

// NewSessionClient 初始化SessionClient
func NewSessionClient(headers map[string]string, jar *cookiejar.Jar) *SessionClient {
	if jar == nil {
		jar, _ = cookiejar.New(nil)
	}

	sc := &SessionClient{
		client: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
					MaxVersion: tls.VersionTLS13,
				},
			},
		},
		headers: make(map[string]string),
		Jar:     jar,
	}

	// 设置默认 Headers 和其他配置
	sc.SetGeneralHeaders(headers)
	sc.EnableRedirect()
	return sc
}

// SetHeader 设置单个 Header
func (sc *SessionClient) SetHeader(key, value string) {
	sc.headers[key] = value
}

// RemoveHeader 删除单个 Header
func (sc *SessionClient) RemoveHeader(key string) {
	delete(sc.headers, key)
}

// applyHeaders 应用 Header 到请求
func (sc *SessionClient) applyHeaders(req *http.Request) {
	for key, value := range sc.headers {
		req.Header.Set(key, value)
	}
}

// Get 发送 GET 请求
func (sc *SessionClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	sc.applyHeaders(req)
	return sc.client.Do(req)
}

// PostJSON 发送 POST 请求
func (sc *SessionClient) PostJSON(
	url string,
	jsonBody []byte,
	implicitHeaders map[string]string,
) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	sc.applyHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	sc.applyImplicitHeaders(req, implicitHeaders)
	return sc.client.Do(req)
}

// PostForm 发送 POST Form 请求
func (sc *SessionClient) PostForm(url string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	sc.applyHeaders(req)
	return sc.client.Do(req)
}

// Post File
func (sc *SessionClient) PostFile(
	url string,
	fileKey, filePath string,
	implicitHeaders map[string]string,
) (*http.Response, error) {
	file, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fileKey, filepath.Base(file.Name()))
	_, _ = io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	sc.applyHeaders(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	sc.applyImplicitHeaders(req, implicitHeaders)
	return sc.client.Do(req)
}

func (sc *SessionClient) applyImplicitHeaders(
	req *http.Request,
	implicitHeaders map[string]string,
) {
	if len(implicitHeaders) > 0 {
		for key, value := range implicitHeaders {
			req.Header.Set(key, value)
		}
	}
}

// defaultRedirectHandler 返回一个默认的重定向处理函数
func (sc *SessionClient) defaultRedirectHandler() func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		// 限制最大重定向次数
		if len(via) >= 10 {
			return http.ErrUseLastResponse // 停止重定向
		}
		// 默认跟随重定向
		return nil
	}
}

// DisableRedirect 禁用重定向
func (sc *SessionClient) DisableRedirect() {
	sc.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // 返回重定向响应而不继续跟随
	}
}

// EnableRedirect 启用默认的重定向处理
func (sc *SessionClient) EnableRedirect() {
	sc.client.CheckRedirect = sc.defaultRedirectHandler()
}

// WaitForElementXPath 使用 XPath 查询 HTML 响应中的元素
func (sc *SessionClient) WaitForElementXPath(
	respBody io.Reader,
	xpath string,
	timeout time.Duration,
) error {
	start := time.Now()
	for {
		// 解析 HTML
		doc, err := htmlquery.Parse(respBody)
		if err != nil {
			return err
		}

		// 查询元素
		node := htmlquery.FindOne(doc, xpath)
		if node != nil {
			return nil
		}

		// 检查超时
		if time.Since(start) > timeout {
			return errors.New("timeout waiting for element matching xpath: " + xpath)
		}

		// 等待一段时间再尝试
		time.Sleep(500 * time.Millisecond)
	}
}

// ready body and reset cursor to 0
func (sc *SessionClient) ReadBody(body io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(body)
	if err != nil {
		return "", err
	}
	body.Close()
	return buf.String(), nil
}

// get body from response, and marshall it to pointer of struct passed in
func (sc *SessionClient) MarshallBodyToStruct(body io.ReadCloser, v interface{}) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(body)
	if err != nil {
		return err
	}
	body.Close()
	return json.Unmarshal(buf.Bytes(), v)
}

func (sc *SessionClient) SaveJar(path string, urls []*url.URL) error {
	return SaveJar(sc.Jar, path, urls)
}
