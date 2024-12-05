package flows

import (
	"log"
	"os"

	"github.com/zsmatrix62/garmin-cli/garmin/actions"
	"github.com/zsmatrix62/garmin-cli/garmin/types"
	"github.com/zsmatrix62/garmin-cli/garmin/types/state"
	"github.com/zsmatrix62/garmin-cli/session"
)

func Login(
	c *session.SessionClient,
	opt *BasicOption,
) (gres *types.FlowGenericResp[types.ActionTicketResponse]) {
	garminHost := opt.GarminHost
	username := opt.Username
	password := opt.Password
	state_file_dir := opt.StateFileDir

	gres = new(types.FlowGenericResp[types.ActionTicketResponse])

	var err error
	defer func() {
		if err != nil {
			gres.Err = err
		}
	}()

	var ticket *types.ActionTicketResponse
	log.Println("Step 1: Accessing Garmin sign-in page...")
	if err = actions.ActionAccessSignInPage(c, garminHost); err != nil {
		return
	}
	log.Println("Step 2: Getting ticket in...")
	t, err := actions.ActionGetTicket(c, garminHost, username, password)
	if err != nil {
		return
	}
	ticket = t
	log.Println("Step 3: Create Session...")
	if err = actions.ActionCreateSession(c, garminHost, ticket); err != nil {
		return
	}
	var token *types.ActionTokenResponse

	{
		log.Println("Step 4: Exchange token...")
		t, aErr := actions.ActionExchange(c, garminHost)
		if aErr != nil {
			err = aErr
			return
		}
		token = t
	}

	if opt.PersistState {
		// check if state_file_dir exists, if not, create it
		if _, sErr := os.Stat(state_file_dir); os.IsNotExist(sErr) {
			_ = os.MkdirAll(state_file_dir, os.ModePerm)
		}

		if err = state.SaveState(c, state_file_dir, username, ticket, token); err != nil {
			return
		}
	}

	return
}
