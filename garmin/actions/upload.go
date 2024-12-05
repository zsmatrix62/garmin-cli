package actions

import (
	"fmt"
	"log"

	"github.com/zsmatrix62/garmin-cli/garmin/types"
	"github.com/zsmatrix62/garmin-cli/garmin/types/state"
	"github.com/zsmatrix62/garmin-cli/session"
)

func ActionUploadFitActivity(
	client *session.SessionClient,
	host string,
	state *state.GarminState,
	filePath string,
) (res *types.ActionUploadResponse, err error) {
	//  visit import-data url in case of 403
	{
		resp, gerr := client.Get(fmt.Sprintf("https://connect.%s/modern/import-data", host))
		if gerr != nil {
			return
		}
		if resp.StatusCode != 200 {
			log.Fatalf("Error: %v", resp.Status)
			return
		}
	}

	fmt.Println("uploading fit file...")

	uploadUrl := fmt.Sprintf("https://connect.%s/upload-service/upload/.fit", host)
	resp, err := client.PostFile(
		uploadUrl,
		"userfile",
		filePath,
		map[string]string{
			"Accept":        "application/json, text/plain, */*",
			"Authorization": fmt.Sprintf("Bearer %s", state.Token.AccessToken),
			"DI-Backend":    fmt.Sprintf("connectapi.%s", host),
		},
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		ferr := fmt.Errorf("Error: %v", resp.Status)
		return nil, ferr
	}
	uploadResponse := new(types.ActionUploadResponse)
	if err = client.MarshallBodyToStruct(resp.Body, uploadResponse); err != nil {
		return nil, err
	}

	return uploadResponse, nil
}
