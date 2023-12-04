package cli

import (
	"log"
	"os"

	"github.com/arindampal-0/gdown-cli-golang/internal/gdown"
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

func init() {
	/* Init -l flag for Download Subcommand */
	downlaodCmd.Flags().BoolVarP(&alsoListDetails, "list", "l", false, "also list details of the item")

	/* Register Subcommands */
	rootCmd.AddCommand(listCmd, downlaodCmd)
}
