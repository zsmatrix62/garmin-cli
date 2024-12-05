package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zsmatrix62/garmin-cli/garmin/types"
	"github.com/zsmatrix62/garmin-cli/session"
)

func ActionGetTicket(
	client *session.SessionClient,
	host string,
	username, password string,
) (ticketResp *types.ActionTicketResponse, err error) {
	ticketResp = &types.ActionTicketResponse{}
	rqst := types.NewTicketRequest(host)
	apiURL := rqst.Url(host)
	// 构造 JSON 请求体
	requestBody := map[string]any{
		"username":     username,
		"password":     password,
		"rememberMe":   true,
		"captchaToken": "",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return ticketResp, fmt.Errorf("failed to marshal JSON body: %w", err)
	}

	// 发送 POST 请求
	resp, err := client.PostJSON(apiURL, jsonBody, nil)
	if err != nil {
		return ticketResp, fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ticketResp, fmt.Errorf("unexpected status code during login: %d", resp.StatusCode)
	}

	err = client.MarshallBodyToStruct(resp.Body, ticketResp)
	if err != nil {
		return ticketResp, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if !strings.EqualFold(ticketResp.ResponseStatus.Type, "SUCCESSFUL") {
		return nil, fmt.Errorf("failed to get ticket: %s", ticketResp.ResponseStatus.Message)
	}
	return ticketResp, nil
}

func ActionCreateSession(
	client *session.SessionClient,
	host string,
	ticke *types.ActionTicketResponse,
) (err error) {
	url := fmt.Sprintf("https://connect.%s/modern?ticket=%s", host, ticke.ServiceTicketID)

	resp, eerr := client.Get(url, nil)
	if eerr != nil {
		return fmt.Errorf("failed to create session: %w", eerr)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code during create session: %d", resp.StatusCode)
	}

	return nil
}
