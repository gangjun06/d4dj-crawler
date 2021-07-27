package conf

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

type Conf struct {
	Port            int
	AssetPath       string
	ToolPath        string
	CrawlerPool     int
	AssetServerPath string
	Debug           bool
}

var conf *Conf

func init() {
	confFile, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal("Error loading config.toml file")
	}
	if _, err := toml.Decode(string(confFile), &conf); err != nil {
		log.Fatal("Error decoding. " + err.Error())
	}
}

func Get() *Conf {
	return conf
}
