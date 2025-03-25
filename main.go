package main

import (
	"runtime"
	"time"

	"github.com/0987363/2table-backend/cmd"
	"github.com/0987363/2table-backend/models"
	log "github.com/sirupsen/logrus"
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
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
