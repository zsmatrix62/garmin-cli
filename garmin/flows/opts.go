package flows

import "fmt"

type BasicOption struct {
	GarminHost   string
	Username     string
	Password     string
	StateFileDir string
	PersistState bool
}

func NewBasicOption(
	region, username, password, state_file_dir string,
	persist_state bool,
) *BasicOption {
	switch region {
	case "cn":
	default:
		region = "com"
	}

	return &BasicOption{
		GarminHost:   fmt.Sprintf("garmin.%s", region),
		Username:     username,
		Password:     password,
		StateFileDir: state_file_dir,
		PersistState: persist_state,
	}
}
