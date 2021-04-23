package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goSeeFuture/gblog/configs"
	"github.com/goSeeFuture/gblog/content"

	"github.com/gofiber/fiber/v2"
)

func homePage(c *fiber.Ctx) error {
	return c.Redirect("/list", http.StatusFound)
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
	return content.Render(c, "list", map[string]interface{}{
		"Categories":       content.Categories(),
		"Title":            configs.Setting.WebsiteName,
		"Footer":           content.Footer(),
		"Total":            total,
		"Pages":            pages(maxPage, curPage),
		"CurPage":          strconv.Itoa(curPage),
		"Articles":         heads,
		"ListType":         "list",
		"PrevPage":         prevPage(maxPage, curPage),
		"NextPage":         nextPage(maxPage, curPage),
		"PageNumberPrefix": configs.Setting.PageNumberPrefix,
		"EnableMathJax":    configs.Setting.ArticleMathJax,
		"Tags":             content.Tags(),
	})
}

func articlePage(c *fiber.Ctx) error {
	filename := filepath.Join(configs.Setting.ArticleDir, c.Params("*"))
	filename, _ = url.QueryUnescape(filename)
	fmt.Println("获取文章:", filename)

	bind := map[string]interface{}{
		"Categories": content.Categories(),
		"Title":      configs.Setting.WebsiteName,
		"Footer":     content.Footer(),
	}

	md, exist := content.FindMetaData(filename)
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
		return content.Render(c, "404", map[string]interface{}{
			"Categories": content.Categories(),
			"Title":      configs.Setting.WebsiteName,
		})
	}
	bind["Article"] = template.HTML(data)
	bind["EnableMathJax"] = configs.Setting.ArticleMathJax

	return content.Render(c, "article", bind)
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

	return content.Render(c, "list", map[string]interface{}{
		"Categories":       content.Categories(),
		"Title":            configs.Setting.WebsiteName,
		"Footer":           content.Footer(),
		"Total":            total,
		"Pages":            pages(maxPage, curPage),
		"CurPage":          strconv.Itoa(curPage),
		"Articles":         heads,
		"ListType":         "category/" + categoryId,
		"PrevPage":         prevPage(maxPage, curPage),
		"NextPage":         nextPage(maxPage, curPage),
		"PageNumberPrefix": configs.Setting.PageNumberPrefix,
		"EnableMathJax":    configs.Setting.ArticleMathJax,
		"Tags":             content.Tags(),
	})
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

	return content.Render(c, "list", map[string]interface{}{
		"Categories":       content.Categories(),
		"Title":            configs.Setting.WebsiteName,
		"Footer":           content.Footer(),
		"Total":            total,
		"Pages":            pages(maxPage, curPage),
		"CurPage":          strconv.Itoa(curPage),
		"Articles":         heads,
		"ListType":         "tag/" + tag,
		"PrevPage":         prevPage(maxPage, curPage),
		"NextPage":         nextPage(maxPage, curPage),
		"PageNumberPrefix": configs.Setting.PageNumberPrefix,
		"EnableMathJax":    configs.Setting.ArticleMathJax,
		"Tags":             content.Tags(),
	})
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

func nextPage(maxPage, curPage int) string {
	if curPage >= maxPage {
		return "#"
	}

	return strconv.Itoa(curPage + 1)
}

func prevPage(maxPage, curPage int) string {
	if curPage <= 1 {
		return "#"
	}

	return strconv.Itoa(curPage - 1)
}

func getPageNumber(s string) int {
	if len(s) < 1+len(configs.Setting.PageNumberPrefix) {
		return 0
	}

	if configs.Setting.PageNumberPrefix != "" && !strings.HasPrefix(s, configs.Setting.PageNumberPrefix) {
		return 0
	}

	s = strings.TrimPrefix(s, configs.Setting.PageNumberPrefix)
	var err error
	page, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return page
}
