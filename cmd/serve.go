package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/0987363/2table-backend/handlers"
	"github.com/0987363/2table-backend/middleware"
)

const defaultAddress = ":8080"

var serveCmd = &cobra.Command{
	Use:    "serve",
	Short:  "Start server",
	PreRun: LoadConfiguration,
	Run:    serve,
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP(
		"address",
		"a",
		defaultAddress,
		"Address the server binds to",
	)
	viper.BindPFlag("address", serveCmd.Flags().Lookup("address"))
}

func serve(cmd *cobra.Command, args []string) {
	if err := middleware.ConnectDB(viper.GetString("database.address")); err != nil {
		log.Fatalf("connect to db: %s failed: %v", viper.GetString("database.address"), err)
	}

	if err := middleware.ConnectStorageManager(); err != nil {
		log.Fatalf("connect to storage failed: %v", err)
	}

	//handlers.Init(cors)

	address := viper.GetString("address")
	log.Infof("Starting server on %s.", address)
	handlers.RootMux.Run(address)
}
