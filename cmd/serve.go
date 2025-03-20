package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/0987363/2table-backend/handlers"
	"github.com/0987363/2table-backend/middleware"
)

const defaultAddress = ":8080"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:    "serve",
	Short:  "Start manager server",
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
	log.Infof("level:%v, dst=%v, store=%v", viper.GetString("log.level"), viper.GetString("log.dst"),
		viper.GetString("log.store"))
	if err := middleware.ConnectLogger(viper.GetString("log.level"),
		viper.GetString("log.dst"),
		viper.GetString("log.store")); err != nil {
		log.Errorf("Connect logger error:%v", err)
	}

	// Try to connect to the database
	if err := middleware.ConnectDB(viper.GetString("database.mongodb")); err != nil {
		log.Fatalf("connect to db: %s failed: %v", viper.GetString("database.mongodb"), err)
	}

	address := viper.GetString("address")
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")

	cors := viper.GetStringSlice("cors")
	log.Infof("Server cors status: %v.", cors)
	handlers.Init(cors)
	//	handlers.Init(BuildInfo.Version + "-" + BuildInfo.Date + "-" + BuildInfo.Commit)
	if cert != "" && key != "" {
		log.Infof("Starting black manager tls server on %s.", address)
		handlers.RootMux.RunTLS(address, cert, key)
	} else {
		log.Infof("Starting black manager server on %s.", address)
		handlers.RootMux.Run(address)
	}
}
