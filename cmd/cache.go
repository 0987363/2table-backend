package cmd

import (
	"github.com/0987363/cache"
	"github.com/0987363/cache/persistence"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"strings"
	"time"
)

func NewCache() {
	cache.SetPageKey("bird:manager:2:" + viper.GetString("release"))

	var store persistence.CacheStore
	memAddr := viper.GetString("memcached")
	if memAddr != "" {
		addrs := strings.Split(memAddr, ", ")
		if len(addrs) > 0 {
			mem := persistence.NewMemcachedStore(addrs, time.Hour*12)
			mem.Client.Timeout = time.Second * 5
			log.Info("Init memcached store success.", addrs, store)
		}
	}
}
