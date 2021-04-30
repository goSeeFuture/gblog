package content

import (
	"html/template"
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

	// 区分出主题介绍
	var t1 []MetaData
	var t2 = make(map[string]MetaData)
	for _, article := range articles {
		if article.isCategoryTopic {
			t2[article.CategoryID] = article
		} else {
			t1 = append(t1, article)
		}
	}
	articles = t1

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

	for i, category := range categories {
		md, exist := t2[category.ID]
		if !exist {
			continue
		}
		category.Topic = template.HTML(md.Summary)
		category.TopicTitle = md.Title
		if category.TopicTitle == "" {
			category.TopicTitle = category.Name
		}
		categories[i] = category
	}

	allcategories.Store(categories)

	return
}
