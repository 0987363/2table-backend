package models

import "fmt"

var BuildInfo struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	Commit  string `json:"commit"`
	Owner   string `json:"owner"`
}

func Version() string {
	return fmt.Sprintf("%s %s %s %s", BuildInfo.Version, BuildInfo.Date, BuildInfo.Commit, BuildInfo.Owner)
}
