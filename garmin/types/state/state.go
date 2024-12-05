package state

import (
	"encoding/json"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/zsmatrix62/garmin-cli/garmin/pkg/helpers"
	"github.com/zsmatrix62/garmin-cli/garmin/types"
	"github.com/zsmatrix62/garmin-cli/session"
)

var garminUrls = []string{
	"https://connect.garmin.cn",
	"https://sso.garmin.cn",
	"https://connectus.garmin.cn",
	"https://connect.garmin.com",
	"https://sso.garmin.com",
	".garmin.cn",
}

func DeleteStateFile(dir, host string, username string) error {
	filePath := helpers.StateFileName(dir, host, username)
	return os.Remove(filePath)
}

func LoadStateFromFile(stateFile string) (*GarminState, error) {
	f, err := os.Open(stateFile)
	if err != nil {
		return nil, err
	}

	state := &GarminState{}

	if err := json.NewDecoder(f).Decode(state); err != nil {
		return nil, err
	}

	return state, nil
}

func LoadState(dir, host string, username string) (*GarminState, error) {
	filePath := helpers.StateFileName(dir, host, username)
	return LoadStateFromFile(filePath)
}

func SaveState(
	client *session.SessionClient,
	baseDir, host string,
	username string,
	ticket *types.ActionTicketResponse,
	token *types.ActionTokenResponse,
) error {
	urls := make([]*url.URL, 0)
	for _, u := range garminUrls {
		url, _ := url.Parse(u)
		urls = append(urls, url)
	}
	filePath := helpers.StateFileName(baseDir, host, username)

	state := &GarminState{
		Ticket:     ticket,
		CookieData: session.CookieDataFromJar(client.Jar, urls),
	}

	state.SetToken(token)

	// serialize state into json bytes and save to file
	stateBytes, err := json.Marshal(state)
	if err != nil {
		return err
	}

	fileDir := path.Dir(filePath)
	if _, err2 := os.Stat(fileDir); os.IsNotExist(err2) {
		if err3 := os.MkdirAll(fileDir, os.ModePerm); err3 != nil {
			return err3
		}
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.Write(stateBytes)
	if err != nil {
		return err
	}

	return nil
}

type GarminState struct {
	Ticket           *types.ActionTicketResponse
	Token            *types.ActionTokenResponse
	CookieData       []session.CookieData
	TokenRefreshedAt int64
}

func (g *GarminState) CookieDataBytes() []byte {
	data, _ := json.Marshal(g.CookieData)
	return data
}

func (g *GarminState) SetToken(token *types.ActionTokenResponse) {
	g.Token = token
	g.TokenRefreshedAt = time.Now().Unix()
}

func (g *GarminState) TokenExpired() bool {
	tokenExpiresIn := g.Token.ExpiresIn
	refreshedAt := g.TokenRefreshedAt
	if refreshedAt == 0 {
		return true
	}
	return time.Now().Unix() > refreshedAt+tokenExpiresIn
}

func (g *GarminState) CanRefreshToken() bool {
	tokenRefreshExpiresIn := g.Token.RefreshTokenExpiresIn
	if tokenRefreshExpiresIn == 0 {
		return false
	}
	return time.Now().Unix() < g.TokenRefreshedAt+tokenRefreshExpiresIn
}
