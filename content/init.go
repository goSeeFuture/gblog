package content

import (
	"log"
	"sort"
	"sync"

	"github.com/goSeeFuture/gblog/configs"
)

const UncategorizedId = "uncategorized"
const UncategorizedName = "未分类"

var (
	articles []MetaData
	tagIndex map[string][]int
	tags     []Tag

	cotentmutex sync.RWMutex
)

func Load() {
	err := loadContent()
	if err != nil {
		log.Fatalln(err)
	}

	go watchArticleChange()
}

func loadContent() (err error) {
	cotentmutex.Lock()
	defer cotentmutex.Unlock()

	err = Layout()
	if err != nil {
		log.Println("load layout failed:", err)
		return
	}

	articles = List()
	// 按照时间倒序
	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].UpdateAt.After(articles[j].UpdateAt)
	})

	tagIndex, tags = makeTagIndex(articles)

	categories := articleCategory(articles)
	mergeCategory(categories, configs.Setting.Categories)

	configs.Setting.Categories = categories

	return
}
