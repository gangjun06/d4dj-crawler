package conf

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

type Conf struct {
	AssetPath       string
	ToolPath        string
	CrawlerPool     int
	AssetServerPath string
	ExtractorPath   string
	Aws             struct {
		BucketName string
		Region     string
		AccessKey  string
		SecretKey  string
	}
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
