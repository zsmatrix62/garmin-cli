package client

import (
	"fmt"

	"github.com/zsmatrix62/garmin-cli/garmin/pkg/helpers"
	"github.com/zsmatrix62/garmin-cli/session"
)

func NewGarminClient(
	garminHost, username, password, state_file_dir string,
) (c *session.SessionClient, isWithStateOk bool) {
	implicitHeaders := map[string]string{
		"Referer": fmt.Sprintf("https://connect.%s/", garminHost),
		"Host":    fmt.Sprintf("sso.%s", garminHost),
	}

	jar, _ := WithGarminStateFile(
		helpers.StateFileName(state_file_dir, username),
	)
	c = session.NewSessionClient(implicitHeaders, jar)
	return c, jar != nil
}
