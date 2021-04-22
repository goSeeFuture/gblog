package content

import (
	"gblog/configs"
	"log"
	"sort"
)

const UncategorizedId = "uncategorized"
const UncategorizedName = "未分类"

var (
	articles []MetaData
	tagIndex map[string][]int
	Tags     []Tag
)

func Load() {
	err := Layout()
	if err != nil {
		log.Fatalln("load layout failed:", err)
	}

	articles = List()
	// 按照时间倒序
	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].UpdateAt.After(articles[j].UpdateAt)
	})

	tagIndex, Tags = makeTagIndex(articles)

	categories := articleCategory(articles)
	mergeCategory(categories, configs.Setting.Categories)

	configs.Setting.Categories = categories
}
