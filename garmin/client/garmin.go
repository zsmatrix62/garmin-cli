package client

import (
	"fmt"
	"log"
	"net/http/cookiejar"

	"github.com/zsmatrix62/garmin-cli/garmin/types/state"
	"github.com/zsmatrix62/garmin-cli/session"
)

func WithGarminCookie(
	state *state.GarminState,
) (*cookiejar.Jar, error) {
	jar, err := session.DeserializeJar(state.CookieDataBytes())
	if err != nil {
		log.Printf("Failed to deserialize cookie data: %v", err)
		return jar, fmt.Errorf("Failed to deserialize cookie data: %v", err)
	} else {
		log.Printf("Deserialize cookie data success: %v", len(state.CookieData))
	}

	return jar, nil
}

func WithGarminStateFile(
	stateFile string,
) (*cookiejar.Jar, error) {
	state, err := state.LoadStateFromFile(stateFile)
	if err != nil {
		log.Printf("Failed to load state file: %v", err)
		return nil, fmt.Errorf("Failed to load state file: %v", err)
	} else {
		log.Printf("Load state file success: %v", state.Ticket.ServiceTicketID)
	}

	return WithGarminCookie(state)
}
