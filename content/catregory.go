package content

import (
	"fmt"
	"path/filepath"
	"strings"

	configs "github.com/goSeeFuture/gblog/configs"
)

func ArticlesByCategoryPage(categoryId string, pageSize, pageNumber int) (total int, heads []MetaData) {
	s := pageSize * (pageNumber - 1)
	e := s + pageSize

	isUncategorized := categoryId == UncategorizedId
	_articles := allarticles.Load().([]MetaData)

	heads = []MetaData{}
	for _, item := range _articles {
		id := strings.ReplaceAll(categoryId, "/", "-")
		var chioce bool
		if isUncategorized {
			if item.CategoryID == "" {
				chioce = true
			}
		} else {
			chioce = item.CategoryID == id
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
	var count = make(map[string]int)
	for _, article := range articles {
		dir := filepath.Dir(article.Filename)

		if configs.Setting.AbsArticleDir == dir {
			count[UncategorizedName]++
			continue
		}

		fmt.Println("cate1", dir)
		dir = strings.TrimPrefix(dir, configs.Setting.ArticleDir)
		if dir != "" && dir[0] == filepath.Separator {
			dir = dir[1:]
		}
		if dir == "" {
			continue
		}
		fmt.Println("cate2", dir)
		fmt.Println("AbsArticleDir", configs.Setting.AbsArticleDir)

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

func Categories() []configs.Category {
	return allcategories.Load().([]configs.Category)
}
