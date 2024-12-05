package flows

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/zsmatrix62/garmin-cli/garmin/actions"
	"github.com/zsmatrix62/garmin-cli/garmin/client"
	"github.com/zsmatrix62/garmin-cli/garmin/types"
	"github.com/zsmatrix62/garmin-cli/garmin/types/state"
)

// FlowUploadActivity uploads a .fit file to Garmin Connect
//
// # Always returns OK with the URL of the uploaded activity despite it being successful or not
//
// Example: {"Err":{},"Ok":"https://connect.garmin.cn/modern/activity/436478761"}
func FlowUploadActivity(opt *BasicOption, filePath string) (gres *types.FlowGenericResp[string]) {
	var err error
	gres = new(types.FlowGenericResp[string])
	defer func() {
		if err != nil {
			gres.Err = err
		}
	}()

	if !strings.EqualFold(path.Ext(filePath), ".fit") {
		err = errors.New("File must be a .fit file")
		return
	}
	if _, fErr := os.Stat(filePath); os.IsNotExist(fErr) {
		err = errors.New("File does not exist")
		return
	}

	garminHost := opt.GarminHost
	username := opt.Username
	password := opt.Password
	state_file_dir := opt.StateFileDir
	gClient, stateLoaded := client.NewGarminClient(
		garminHost,
		username,
		password,
		state_file_dir,
	)

	var ticket *types.ActionTicketResponse
	if !stateLoaded || !opt.PersistState {
		loginRes := Login(gClient, opt)
		if loginRes.Err != nil {
			err = loginRes.Err
			return
		} else {
			ticket = loginRes.Ok
		}
	}

	s, _ := state.LoadState(state_file_dir, garminHost, username)
	if s.TokenExpired() && s.CanRefreshToken() {
		log.Println("Token expired, refreshing...")
		// if token, err := actions.ActionRefreshToken(gClient, garminHost, s); err != nil {
		if token, gerr := actions.ActionExchange(gClient, garminHost); err != nil {
			err = gerr
			return
		} else {
			s.SetToken(token)

			if ticket != nil && opt.PersistState {
				// Step 5: save sate
				if err = state.SaveState(gClient, state_file_dir, garminHost, username, ticket, token); err != nil {
					return
				}
			} else {
				if err = state.DeleteStateFile(state_file_dir, garminHost, username); err != nil {
					return
				}
			}
		}
	}

	uRes, err := actions.ActionUploadFitActivity(gClient, garminHost, s, filePath)
	if err != nil {
		_ = state.DeleteStateFile(state_file_dir, garminHost, username)
		return
	}

	activityUrlFmt := fmt.Sprintf("https://connect.%s/modern/activity/%s", garminHost, "%d")
	if uRes.Fails() {
		if len(uRes.Failures()) > 0 {
			gres.Err = errors.New(uRes.Failures()[0].Messages[0].Content)
			u := fmt.Sprintf(activityUrlFmt, uRes.Failures()[0].InternalID)
			gres.Ok = &u
		}
	} else {
		_s := uRes.DetailedImportResult.UploadUUID.Uuid
		if uRes2, err := actions.ActionCheckActivityStatus(gClient, garminHost, s, _s); err != nil {
			gres.Err = err
			return
		} else {
			if uRes2.Fails() {
				if len(uRes2.Failures()) > 0 {
					gres.Err = errors.New(uRes.Failures()[0].Messages[0].Content)
					u := fmt.Sprintf(activityUrlFmt, uRes.Failures()[0].InternalID)
					gres.Ok = &u
				}
			} else {
				if uRes2.Success() {
					if len(uRes2.Successes()) > 0 {
						u := fmt.Sprintf(activityUrlFmt, uRes2.Successes()[0].InternalID)
						gres.Ok = &u
					}
				}
			}
		}
	}
	return
}
