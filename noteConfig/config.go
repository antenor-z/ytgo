package noteConfig

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Server struct {
	Domain    string `toml:"domain"`
	DebugMode bool   `toml:"debugmode"`
}
type Config struct {
	Server Server
}

var config Config

func ConfigInit() {
	dat, err := os.ReadFile("config.toml")
	if err != nil {
		panic("ConfigInit(): " + err.Error())
	}
	_, err2 := toml.Decode(string(dat), &config)
	if err2 != nil {
		panic("ConfigInit(): " + err2.Error())
	}
}

func GetDomain() string {
	return config.Server.Domain
}

func IsDebug() bool {
	return config.Server.DebugMode
}

func GetVersion() string {
	dat, err := os.ReadFile("VERSION")
	if err != nil {
		dat, err := os.ReadFile("../VERSION")
		if err != nil {
			panic("GetVersion(): " + err.Error())
		}
		return string(dat)
	}
	return string(dat)
}
