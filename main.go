package main

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zsmatrix62/garmin-cli/garmin/flows"
)

var (
	username string
	password string
	flow     string

	region       string // optional
	persistState bool   // optional
	stateFileDir string // optional
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "app",
		Short: "A brief description of your application",
		Run: func(cmd *cobra.Command, args []string) {
			if username == "" || password == "" {
				log.Fatal("Error: username and password are required.")
			}

			if strings.EqualFold(flow, "") {
				log.Fatal("flow is required")
			}

			if persistState && stateFileDir == "" {
				log.Fatal("Error: state file dir is required")
			}

			opts := flows.NewBasicOption(region, username, password, stateFileDir, persistState)

			switch flow {
			case "upload-fit":
				if len(args) != 1 {
					log.Fatalf("Error: upload-fit requires a single argument")
				}
				gres := flows.FlowUploadActivity(opts, args[0])
				if gres.Err != nil {
					log.Fatalf("Error: %v", gres.Err)
				} else {
					log.Printf("Success: %v", gres.Ok)
				}
			default:
				log.Fatalf("Error: unknown flow %s", flow)
			}
		},
	}

	rootCmd.Flags().
		StringVarP(&region, "region", "r", "cn", "cn or global, cn - uses garmin.cn, and global - uses garmin.com")
	rootCmd.Flags().StringVarP(&username, "username", "u", "", "Username for Garmin account")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "Password for Garmin account")
	rootCmd.Flags().
		StringVarP(&stateFileDir, "save_state", "", "./garmin-cli-states", "state base dir path")
	rootCmd.Flags().StringVarP(&flow, "flow", "w", "", "Flow to run (e.g., upload)")
	rootCmd.Flags().
		BoolVarP(&persistState, "persist_state", "", true, "Persist ticket, cookie and tokens to file")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(1)
	}
}
