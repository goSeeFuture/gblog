package content

import (
	"path/filepath"
	"strings"

	configs "github.com/goSeeFuture/gblogconfigs"
)

func ArticlesByCategoryPage(categoryId string, pageSize, pageNumber int) (total int, heads []MetaData) {
	s := pageSize * (pageNumber - 1)
	e := s + pageSize

	isUncategorized := categoryId == UncategorizedId

	heads = []MetaData{}
	for _, item := range articles {
		var chioce bool
		if isUncategorized {
			if filepath.Dir(item.Filename) == "articles" {
				chioce = true
			}
		} else {
			chioce = item.CategoryID == categoryId
		}
		if chioce {
			if total >= s && total < e {
				heads = append(heads, item)
			}
			total++
		}
	}

	return
}

func articleCategory(articles []MetaData) []configs.Category {
	var unique = make(map[string]bool)
	var categories []configs.Category
	var articleDir = configs.Setting.ArticleDir + "/"
	var count = make(map[string]int)
	for _, article := range articles {
		dir := filepath.Dir(article.Filename)
		if configs.Setting.ArticleDir == dir {
			count[UncategorizedName]++
			continue
		}

		dir = strings.TrimPrefix(dir, articleDir)
		cate := configs.Category{
			Path: dir,
			Name: dir,
			ID:   strings.Replace(dir, "/", "-", -1),
		}

		count[cate.Name]++
		if unique[cate.Name] {
			continue
		}

		unique[cate.Name] = true
		categories = append(categories, cate)
	}

	// 添加未分类条目
	if configs.Setting.ShowUncategorized {
		categories = append(categories, configs.Category{
			Name: UncategorizedName,
			Path: UncategorizedId,
			ID:   UncategorizedId,
		})
	}

	for i, e := range categories {
		e.Count = count[e.Name]
		categories[i] = e
	}
	return categories
}

// 合并b到a，忽略不在a中分类，以Path为key
func mergeCategory(a, b []configs.Category) {
	for _, eb := range b {
		for i, ea := range a {
			if eb.Path == ea.Path {
				ea.Name = eb.Name
				a[i] = ea
				break
			}
		}
	}
}
