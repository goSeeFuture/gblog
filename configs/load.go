package configs

import (
	"log"

	"github.com/pelletier/go-toml"
)

const TimeLayout = "2006-01-02 15:04:05"

const (
	// 自定义404页
	CustomPage404 = "404.md"
	// 自定义页脚
	CustomPageFooter = "footer.md"
	// yaml数据头长度
	MaxMetaDataLen = 256
	// 显示未分类
	ShowUncategorized = true
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
