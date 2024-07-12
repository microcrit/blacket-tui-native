package config

import (
	configTypes "crit.rip/blacket-tui/config/types"
	utilTypes "crit.rip/blacket-tui/util"
	"github.com/pelletier/go-toml"
)

type ProxyScraperConfig struct {
	File    *string
	Max     *int
	Threads *int
}

type Config struct {
	Accounts     utilTypes.Either[*[]configTypes.Account, *[]interface{}] `toml:"Accounts"`
	ProxyScraper *ProxyScraperConfig                                      `toml:"ProxyScraper"`
}

func ParseConfig(path string) map[string]interface{} {
	config, err := toml.LoadFile(path)
	if err != nil {
		panic(err)
	}
	mapx := config.ToMap()
	return mapx
}
