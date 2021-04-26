package content

import (
	"log"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/goSeeFuture/gblog/configs"
)

const UncategorizedId = "uncategorized"
const UncategorizedName = "未分类"

var (
	// allarticles []MetaData
	allarticles atomic.Value
	// articleVersion map[string]int64
	articleVersion atomic.Value
	// tagIndex    map[string][]int
	tagIndex atomic.Value
	// alltags  []Tag
	alltags atomic.Value
	// []configs.Category
	allcategories atomic.Value

	reloadMutex sync.Mutex

	// linkDir string
	linkDir atomic.Value
)

func Load() {
	err := loadContent()
	if err != nil {
		log.Fatalln(err)
	}

	go watchArticleChange()
}

func loadContent() (err error) {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()

	err = initLayoutTemplate()
	if err != nil {
		log.Println("init layout template failed:", err)
		return
	}

	articles := List()
	if len(articles) == 0 {
		articles = []MetaData{}
	}

	// 按照时间倒序
	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].UpdateAt.After(articles[j].UpdateAt)
	})
	allarticles.Store(articles)

	version := make(map[string]int64)
	for _, e := range articles {
		version[e.Filename] = e.ModifyTime.Unix()
	}
	articleVersion.Store(version)

	tagindex, tags := makeTagIndex(articles)
	if len(tagindex) == 0 {
		tagindex = make(map[string][]int)
	}
	if len(tags) == 0 {
		tags = []Tag{}
	}

	alltags.Store(tags)
	tagIndex.Store(tagindex)

	categories := articleCategory(articles)
	mergeCategory(categories, configs.Setting.Categories)
	if len(categories) == 0 {
		categories = []configs.Category{}
	}
	allcategories.Store(categories)

	return
}
