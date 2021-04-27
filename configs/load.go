package configs

import (
	"log"

	"github.com/pelletier/go-toml"
)

const TimeLayout = "2006-01-02 15:04:05"
const (
	CustomPage404    = "404.md"
	CustomPageFooter = "footer.md"
)

func Load() {
	tree, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatalln("load config.toml file failed:", err)
	}

	err = tree.Unmarshal(&Setting)
	if err != nil {
		log.Fatalln("config.toml file syntax error:", err)
	}
}
