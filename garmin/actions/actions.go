package actions

import (
	"fmt"
	"net/http"

	"github.com/zsmatrix62/garmin-cli/session"
)

// AccessSignInPage 访问 Garmin 登录页面并等待特定元素
func ActionAccessSignInPage(client *session.SessionClient, host string) error {
	garminURL := fmt.Sprintf(
		"https://sso.%s/portal/sso/zh-CN/sign-in?clientId=GarminConnect&service=https://connect.garmin.cn/modern",
		host,
	)
	fmt.Println("Accessing Garmin sign-in page...", garminURL)
	resp, err := client.Get(garminURL)
	if err != nil {
		return fmt.Errorf("failed to access Garmin sign-in page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// // 等待页面加载并检测包含 class 为 validation-form 的元素
	// return client.WaitForElementXPath(resp.Body, `//*[@class="validation-form"]`, 30*time.Second)
	return nil
}
