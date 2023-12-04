package cli

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/ArindamPal-0/gdown-cli-golang/internal/gdown"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gdown",
	Short: "Download files and folders from Google Drive",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var listCmd = &cobra.Command{
	Use:   "list <id>",
	Short: "List details if the item",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		err := gdown.ListItemDetails(id)
		if err != nil {
			log.Fatalf("Could not list item details.\n%v", err)
		}
	},
}

var alsoListDetails bool

var downlaodCmd = &cobra.Command{
	Use:   "download <id>",
	Short: "Download the item",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		/* id from CLI args */
		id := args[0]
		err := gdown.DownloadItem(id, alsoListDetails)
		if err != nil {
			log.Fatalf("Could not download item.\n%v", err)
		}
	},
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure gdown, creates gdown config directory and more...",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		/* Get User Config directory */
		configDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatalf("Could not get User Config Directory.\n%v", err)
		}

		/* Generate application config Path */
		configPath := path.Join(configDir, "gdown")

		/* Generate credentials path for Google Service Account */
		var saCredentialsPath string = path.Join(configPath, "service-account", "credentials.json")

		/* Get the folder path for Service Account credentials file */
		saCredentialsFolder := path.Dir(saCredentialsPath)

		/* Create the folder for Service Account Credentials file */
		os.MkdirAll(saCredentialsFolder, os.ModePerm)

		/* Ask user to put credentials.json file in the folder created */
		fmt.Printf("Put the Service Account credentials.json file in the following directory %v\n", saCredentialsFolder)
	},
}

func init() {
	/* Init -l flag for Download Subcommand */
	downlaodCmd.Flags().BoolVarP(&alsoListDetails, "list", "l", false, "also list details of the item")

	/* Register Subcommands */
	rootCmd.AddCommand(listCmd, downlaodCmd, configureCmd)
}
