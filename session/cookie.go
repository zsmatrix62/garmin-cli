package session

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

type CookieData struct {
	URL     string         `json:"url"`
	Cookies []*http.Cookie `json:"cookies"`
}

func CookieDataFromJar(jar *cookiejar.Jar, urls []*url.URL) []CookieData {
	var cookieData []CookieData

	// 遍历所有 URL，提取对应的 cookies
	for _, u := range urls {
		cookies := jar.Cookies(u)
		if len(cookies) > 0 {
			cookieData = append(cookieData, CookieData{
				URL:     u.String(),
				Cookies: cookies,
			})
		}
	}
	return cookieData
}

func SerializeJar(jar *cookiejar.Jar, urls []*url.URL) ([]byte, error) {
	return json.Marshal(CookieDataFromJar(jar, urls))
}

func DeserializeJar(data []byte) (*cookiejar.Jar, error) {
	var cookieData []CookieData
	err := json.Unmarshal(data, &cookieData)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// 恢复 cookies 到 jar
	for _, cd := range cookieData {
		u, err := url.Parse(cd.URL)
		if err != nil {
			return nil, err
		}
		jar.SetCookies(u, cd.Cookies)
	}

	return jar, nil
}

func SaveJar(jar *cookiejar.Jar, path string, urls []*url.URL) error {
	data, err := SerializeJar(jar, urls)
	if err != nil {
		return err
	}

	// save to file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func LoadJar(path string) (*cookiejar.Jar, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return DeserializeJar(data)
}
