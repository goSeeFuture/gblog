package service

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/goSeeFuture/gblog/configs"
	"github.com/goSeeFuture/gblog/content"

	"github.com/gofiber/fiber/v2"
)

func homePage(c *fiber.Ctx) error {
	return c.Redirect("/list", http.StatusFound)
}

// 页头，页脚等公共数据
// 以覆盖的方式修改
func frame(c *fiber.Ctx) map[string]interface{} {
	return map[string]interface{}{
		"Title":         configs.Setting.WebsiteName,
		"Categories":    content.Categories(),
		"Tags":          content.Tags(),
		"Footer":        content.Footer(),
		"EnableMathJax": configs.Setting.ArticleMathJax,
	}
}

// 合并b到a
func mergeBind(a, b map[string]interface{}) map[string]interface{} {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func listPage(c *fiber.Ctx) error {
	var curPage int
	page := c.Params("page")
	if page == "" {
		curPage = pageNumber(c, "listpn")
	} else {
		if curPage = getPageNumber(page); curPage == 0 {
			return fiber.NewError(http.StatusBadRequest, "page must be number")
		}
	}

	total, heads := content.ArticlesByPage(configs.Setting.PageSize, curPage)
	maxPage := maxPage(total, configs.Setting.PageSize)
	listtype := "list"

	bind := map[string]interface{}{
		"Total":    total,
		"Pages":    pages(maxPage, curPage),
		"CurPage":  strconv.Itoa(curPage),
		"Articles": heads,
		"ListType": listtype,
		"PrevPage": prevPage(maxPage, curPage, listtype),
		"NextPage": nextPage(maxPage, curPage, listtype),
	}
	return content.Render(c, "list", mergeBind(frame(c), bind))
}

func articlePage(c *fiber.Ctx) error {
	article := c.Params("*")
	article, _ = url.QueryUnescape(article)
	filename := filepath.Join(configs.Setting.AbsArticleDir, article)
	log.Println("获取文章:", filename)

	bind := make(map[string]interface{})

	md, exist := content.FindMetaData(article)
	if exist {
		bind["Title"] = md.Title
		bind["ArticleTags"] = md.Tags
		bind["InlineTitle"] = md.InlineTitle
		if md.UpdateAt.IsZero() {
			bind["UpdateAt"] = ""
		} else {
			bind["UpdateAt"] = md.UpdateAt.Format(configs.TimeLayout)
		}
	}
	bind["HasMetaHead"] = md.HasMetaHead

	data, err := content.MarkdownPage(filename, md.Offset)
	if err != nil {
		log.Println("get article failed:", err)

		return content.Render(c, "custom", mergeBind(frame(c), map[string]interface{}{
			"Content": content.Page404(),
		}))
	}

	if md.IsDraft && !isAuthor(c) {
		return content.Render(c, "custom", mergeBind(frame(c), map[string]interface{}{
			"Content": content.PageNotAuthor(),
		}))
	}

	bind["Article"] = template.HTML(data.Content)
	bind["TOC"] = data.Heads
	return content.Render(c, "article", mergeBind(frame(c), bind))
}

func categoryPage(c *fiber.Ctx) error {
	var curPage int
	page := c.Params("page")
	if page == "" {
		curPage = pageNumber(c, "categorypn")
	} else {
		if curPage = getPageNumber(page); curPage == 0 {
			return fiber.NewError(http.StatusBadRequest, "page must be number")
		}
	}

	categoryId := c.Params("categoryid")
	categoryId, _ = url.QueryUnescape(categoryId)

	total, heads := content.ArticlesByCategoryPage(categoryId, configs.Setting.PageSize, curPage)
	maxPage := maxPage(total, configs.Setting.PageSize)

	category, _ := content.GetCategory(categoryId)

	listtype := "category/" + categoryId
	bind := map[string]interface{}{
		"Total":    total,
		"Pages":    pages(maxPage, curPage),
		"CurPage":  strconv.Itoa(curPage),
		"Articles": heads,
		"Category": category,
		"ListType": listtype,
		"PrevPage": prevPage(maxPage, curPage, listtype),
		"NextPage": nextPage(maxPage, curPage, listtype),
	}
	return content.Render(c, "list", mergeBind(frame(c), bind))
}

func tagPage(c *fiber.Ctx) error {
	var curPage int
	page := c.Params("page")
	if page == "" {
		curPage = pageNumber(c, "tagpn")
	} else {
		if curPage = getPageNumber(page); curPage == 0 {
			return fiber.NewError(http.StatusBadRequest, "page must be number")
		}
	}

	tag := c.Params("tag")
	tag, _ = url.QueryUnescape(tag)
	total, heads := content.ArticlesByTagPage(tag, configs.Setting.PageSize, curPage)
	maxPage := maxPage(total, configs.Setting.PageSize)

	listtype := "tag/" + tag
	bind := map[string]interface{}{
		"Total":    total,
		"Pages":    pages(maxPage, curPage),
		"CurPage":  strconv.Itoa(curPage),
		"Articles": heads,
		"ListType": listtype,
		"PrevPage": prevPage(maxPage, curPage, listtype),
		"NextPage": nextPage(maxPage, curPage, listtype),
	}
	return content.Render(c, "list", mergeBind(frame(c), bind))
}

func maxPage(total int, pageSize int) int {
	maxPageNumber := total / pageSize
	if total%pageSize != 0 {
		maxPageNumber++
	}
	return maxPageNumber
}

func pages(maxPage, curPage int) []string {
	if maxPage < 1 {
		return []string{}
	}

	if curPage > maxPage {
		curPage = maxPage
	}

	const hellip = "..."

	pages := make([]int, 0, 7)
	pages = append(pages, 1)

	prev := curPage - 1
	if prev > 1 && prev != maxPage {
		pages = append(pages, prev)
	}

	if curPage != 1 && curPage != maxPage {
		pages = append(pages, curPage)
	}

	next := curPage + 1
	if next != 1 && next < maxPage {
		pages = append(pages, next)
	}

	if maxPage != 1 {
		pages = append(pages, maxPage)
	}

	pg := make([]string, 0, len(pages)+2)
	var pp = 0
	for _, page := range pages {
		if pp+1 < page {
			pg = append(pg, hellip)
		}
		pg = append(pg, strconv.Itoa(page))
		pp = page
	}

	return pg
}

func nextPage(maxPage, curPage int, listType string) string {
	if curPage >= maxPage {
		return "#"
	}

	return "/" + listType + "/" + strconv.Itoa(curPage+1)
}

func prevPage(maxPage, curPage int, listType string) string {
	if curPage <= 1 {
		return "#"
	}

	return "/" + listType + "/" + strconv.Itoa(curPage-1)
}

func getPageNumber(s string) int {
	if len(s) < 1 {
		return 0
	}

	var err error
	page, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return page
}

func author(c *fiber.Ctx) error {
	value := c.Params("value")

	apvalue, err := getAuthorPreviewValue()
	if err != nil {
		panic(err)
	}

	if value != apvalue {
		return fiber.NewError(fiber.StatusUnauthorized, "认证失败！")
	}

	sess, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	sess.Set(authorField, true)
	log.Println("set author")
	if err := sess.Save(); err != nil {
		panic(err)
	}

	return homePage(c)
}
