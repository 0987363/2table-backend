package main

import (
	"github.com/0987363/2table-backend/cmd"
	"github.com/0987363/2table-backend/models"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"time"
)

// Exported onstants for storing build information
var (
	BuildVersion string
	BuildDate    string
	BuildCommit  string
	BuildOwner   string
)

func init() {
	time.Local = time.UTC
	runtime.GOMAXPROCS(runtime.NumCPU())

	models.BuildInfo.Version = BuildVersion
	models.BuildInfo.Date = BuildDate
	models.BuildInfo.Commit = BuildCommit
	models.BuildInfo.Owner = BuildOwner
	models.BuildInfo.Env = os.Getenv("PLATFORM_ENV")
}

// @title           Ecotopia Manager API
// @version         1.0
// @description     This is Druid Ecotopia Manager Server.
// @termsOfService  https://druid.tech/terms-of-use/

// @contact.name   API Support
// @contact.url    https://druid.tech/request-form/
// @contact.email  sales@druid.tech

// @host      localhost:9080
// @BasePath  /manager/api

// @securityDefinitions.basic  BasicAuth
func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
