package cmd

import (
	"fmt"

	"github.com/0987363/2table-backend/models"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version information",
	Run:   version,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	fmt.Println(models.Version())
}
