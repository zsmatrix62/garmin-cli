package actions

import (
	"fmt"
	"log"
	"strings"

	"github.com/k0kubun/pp/v3"
	"github.com/zsmatrix62/garmin-cli/garmin/types"
	"github.com/zsmatrix62/garmin-cli/garmin/types/state"
	"github.com/zsmatrix62/garmin-cli/session"
)

// FIXME: not working as of 2024-12-05
func ActionCheckActivityStatus(
	client *session.SessionClient,
	host string,
	s *state.GarminState,
	uuid string,
) (res *types.ActionUploadResponse, err error) {
	// FIXME: 1733399269811 below is a dynamic value, it should be fetched from somewhere

	// https://connect.garmin.cn/activity-service/activity/status/1733399269811/9f40a0e5791a4b308c4f81531b846c07
	url := fmt.Sprintf(
		"https://connect.garmin.cn/activity-service/activity/status/1733399269811/%s",
		strings.ReplaceAll(uuid, "-", ""),
	)

	resp, err := client.Get(url,
		map[string]string{
			"Accept":        "application/json, text/plain, */*",
			"Authorization": fmt.Sprintf("Bearer %s", s.Token.AccessToken),
			"DI-Backend":    fmt.Sprintf("connectapi.%s", host),
		},
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		ferr := fmt.Errorf("Error: %v", resp.Status)
		pp.Println(3, ferr)
		return nil, ferr
	}

	uploadResponse := new(types.ActionUploadResponse)
	if err = client.MarshallBodyToStruct(resp.Body, uploadResponse); err != nil {
		return nil, err
	}

	return uploadResponse, nil
}

func ActionUploadFitActivity(
	client *session.SessionClient,
	host string,
	state *state.GarminState,
	filePath string,
) (res *types.ActionUploadResponse, err error) {
	//  visit import-data url in case of 403
	{
		resp, gerr := client.Get(fmt.Sprintf("https://connect.%s/modern/import-data", host), nil)
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
