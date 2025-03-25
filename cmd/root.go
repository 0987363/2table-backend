package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "2table",
	Short: "2table backend server",
}

var configFilePath string

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&configFilePath, "config", "c", "", "Path to the config file",
	)
}

func LoadConfiguration(cmd *cobra.Command, args []string) {
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/2table")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		if viper.ConfigFileUsed() == "" {
			log.Fatalf("Unable to find configuration file.")
		}

		log.Fatalf("Failed to load %s: %v", viper.ConfigFileUsed(), err)
	} else {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}
}
