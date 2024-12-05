package actions

import (
	"encoding/json"
	"fmt"

	"github.com/zsmatrix62/garmin-cli/garmin/types"
	"github.com/zsmatrix62/garmin-cli/garmin/types/state"
	"github.com/zsmatrix62/garmin-cli/session"
)

// / TODO:
func ActionRefreshToken(
	client *session.SessionClient,
	host string,
	s *state.GarminState,
) (token *types.ActionTokenResponse, err error) {
	token = new(types.ActionTokenResponse)
	url := fmt.Sprintf("https://connect.%s/services/auth/token/refresh", host)

	implicitHeaders := make(map[string]string)
	implicitHeaders["Authorization"] = fmt.Sprintf("Bearer %s", s.Token.RefreshToken)
	resp, err := client.PostJSON(url, nil, implicitHeaders)
	if err != nil {
		return nil, err
	}

	jsonText, err := client.ReadBody(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	err = json.Unmarshal([]byte(jsonText), token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ActionExchange(
	client *session.SessionClient,
	host string,
) (token *types.ActionTokenResponse, err error) {
	token = new(types.ActionTokenResponse)
	url := fmt.Sprintf("https://connect.%s/modern/di-oauth/exchange", host)
	resp, err := client.PostJSON(url, nil, make(map[string]string))
	if err != nil {
		return nil, err
	}

	jsonText, err := client.ReadBody(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	err = json.Unmarshal([]byte(jsonText), token)
	if err != nil {
		return nil, err
	}
	return token, nil
}
