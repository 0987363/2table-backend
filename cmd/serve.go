package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/0987363/2table-backend/handlers"
	"github.com/0987363/2table-backend/middleware"
	"github.com/0987363/2table-backend/models"
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
	if err := middleware.ConnectDB(viper.GetString("database.address")); err != nil {
		log.Fatalf("connect to db: %s failed: %v", viper.GetString("database.address"), err)
	}

	Type := viper.GetString("storage.type")
	cfg := models.StorageConfig{Type: Type}
	switch Type {
	case models.StorageTypeLocal:
		cfg.Local = &models.LocalConfig{Path: viper.GetString("storage.local_config.path")}
	case models.StorageTypeS3:
		cfg.S3 = &models.S3Config{
			Bucket:   viper.GetString("storage.local_config.path"),
			Region:   viper.GetString("storage.local_config.region"),
			Endpoint: viper.GetString("storage.local_config.endpoint"),
		}
	default:
		log.Fatal("Unknown storage type: ", Type)
	}
	if err := middleware.ConnectStorageManager(&cfg); err != nil {
		log.Fatalf("connect to storage: %v failed: %v", cfg, err)
	}

	address := viper.GetString("address")
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")

	cors := viper.GetStringSlice("cors")
	log.Infof("Server cors status: %v.", cors)
	handlers.Init(cors)
	//	handlers.Init(BuildInfo.Version + "-" + BuildInfo.Date + "-" + BuildInfo.Commit)
	if cert != "" && key != "" {
		log.Infof("Starting tls server on %s.", address)
		handlers.RootMux.RunTLS(address, cert, key)
	} else {
		log.Infof("Starting server on %s.", address)
		handlers.RootMux.Run(address)
	}
}
