package content

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/goSeeFuture/gblog/configs"

	"gopkg.in/yaml.v2"
)

type MetaData struct {
	// 文章标签
	Tags []string `yaml:"Tags"`
	// 文章更新时间
	UpdateAt time.Time `yaml:"UpdateAt"`
	// 文章副标题
	Subtitle string `yaml:"Subtitle"`
	// 文章标题
	Title string `yaml:"Title"`
	// 草稿标记
	// 通过控制台qrcode得到预览草稿链接
	IsDraft bool `yaml:"Draft"`

	// 文件路径
	Filename string `yaml:"-"`
	// 除去yaml头的偏移量
	Offset int `yaml:"-"`
	// 概要，取文章前256个字符
	Summary template.HTML `yaml:"-"`
	// 内联标题
	InlineTitle bool `yaml:"-"`
	// 存在meta头
	HasMetaHead bool `yaml:"-"`
	// 分类ID
	CategoryID string `yaml:"-"`
	// 文件修改时间，用于检查修改变动
	ModifyTime time.Time `yaml:"-"`

	// 分类主题介绍
	isCategoryTopic bool
}

func (md MetaData) UpdateFromNow() string {
	dur := time.Since(md.UpdateAt)
	if dur.Seconds() < 10 {
		return "刚刚"
	}

	if dur.Minutes() < 1 {
		return strconv.FormatInt(int64(dur.Seconds()), 10) + " 秒前"
	}

	if dur.Hours() < 1 {
		return strconv.FormatInt(int64(dur.Minutes()), 10) + " 分钟前"
	}

	if dur.Hours() < 24 {
		return strconv.FormatInt(int64(dur.Hours()), 10) + " 小时前"
	}

	if dur.Hours() < 720 {
		return strconv.FormatInt(int64(dur.Hours()), 10) + " 天前"
	}

	if dur.Hours() < 8760 {
		return strconv.FormatInt(int64(dur.Hours()/720), 10) + " 个月前"
	}

	return strconv.FormatInt(int64(dur.Hours()/8760), 10) + " 年前"
}

// 罗列所有文章的Meta头
func List() (articles []MetaData) {
	link, _ := os.Readlink(configs.Setting.ArticleDir)
	if link != "" {
		configs.Setting.AbsArticleDir = link
	} else {
		configs.Setting.AbsArticleDir, _ = filepath.Abs(configs.Setting.ArticleDir)
	}

	err := filepath.Walk(configs.Setting.AbsArticleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isSpecialArticle(path) {
			return nil
		}

		if isArticleRefDirectory(path, link) {
			// 跳过引用资源目录
			return nil
		}

		if strings.ToLower(filepath.Ext(info.Name())) != ".md" {
			return nil
		}

		md, err := loadMetaData(path)
		if err != nil {
			return nil
		}

		articles = append(articles, md)
		return nil
	})
	if err != nil {
		log.Println("foreach article failed:", err)
	}
	return
}

func reloadArticleMetaData(filename string) (MetaData, bool) {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()

	md, err := loadMetaData(filename)
	if err != nil {
		return MetaData{}, false
	}

	articles := allarticles.Load().([]MetaData)
	t := make([]MetaData, 0, len(articles))
	for _, e := range articles {
		if isSameFile(filename, e.Filename) {
			t = append(t, md)
		} else {
			t = append(t, e)
		}
	}

	allarticles.Store(t)
	return md, true
}

func loadMetaData(filename string) (MetaData, error) {
	dir := filepath.Dir(filename)
	dir = strings.TrimPrefix(strings.TrimPrefix(dir, configs.Setting.AbsArticleDir), "/")
	if dir != "" && dir[0] == filepath.Separator {
		dir = dir[1:]
	}

	file := strings.TrimPrefix(strings.TrimPrefix(filename, configs.Setting.AbsArticleDir), "/")
	if file != "" && file[0] == filepath.Separator {
		file = file[1:]
	}
	file = filepath.Join(configs.Setting.ArticleDir, file)
	file = filepath.ToSlash(file)

	var err error
	md := MetaData{
		Title:      filepath.Base(filename),
		Filename:   file,
		CategoryID: strings.ReplaceAll(dir, "/", "-"),
	}

	meta, offset := getMetaData(filename)
	if meta != nil {
		err = yaml.Unmarshal(meta, &md)
		if err != nil {
			return MetaData{}, err
		}
		md.HasMetaHead = true
	}

	head, h1, modtm, istopic := getHeadContent(filename, offset)
	if len(head) != 0 {
		md.Summary = template.HTML(head)
	}
	md.ModifyTime = modtm
	md.isCategoryTopic = istopic
	if h1 != "" {
		md.InlineTitle = true
		md.Title = h1
	}
	if md.UpdateAt.IsZero() && !modtm.IsZero() {
		md.UpdateAt = modtm
	}

	md.Offset = offset
	return md, nil
}

func getHeadContent(filename string, offset int) ([]byte, string, time.Time, bool) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("open article file failed:", err)
		return nil, "", time.Time{}, false
	}
	defer file.Close()

	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		log.Println("read seek article file failed:", err)
		return nil, "", time.Time{}, false
	}

	const headPartSize int64 = 256
	var readsize int64
	fs, err := file.Stat()
	if err != nil {
		log.Println("get article file head failed:", err)
		return nil, "", time.Time{}, false
	}

	var istopic = configs.Setting.CategoryTopic && filepath.Base(filename) == "topic.md"
	var isPart = fs.Size() >= headPartSize
	if !istopic && isPart {
		readsize = headPartSize
	} else {
		readsize = fs.Size()
	}

	part := make([]byte, readsize)
	partsize, err := file.Read(part)
	if err != nil && err != io.EOF {
		log.Println("get article file head failed:", err)
		return nil, "", time.Time{}, false
	}

	part = part[:partsize]
	// 去掉末尾不完整的字符
	for n := len(part) - 1; n >= 0; n-- {
		r, _ := utf8.DecodeLastRune(part[:n])
		if r != utf8.RuneError {
			part = part[:n]
			break
		}
	}

	if !istopic && isPart {
		part = append(part, []byte("...")...)
	}

	c, _ := markdown2HTML(part)
	header, h1 := removeH1(c)
	return header, h1, fs.ModTime(), istopic
}

func getMetaData(filename string) ([]byte, int) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("open article file failed:", err)
		return nil, 0
	}
	defer file.Close()

	fs, err := file.Stat()
	if err != nil {
		log.Println("open article file failed:", err)
		return nil, 0
	}

	var data []byte
	if fs.Size() < int64(configs.MaxMetaDataLen) {
		data = make([]byte, fs.Size())
	} else {
		data = make([]byte, configs.MaxMetaDataLen)
	}

	_, err = file.Read(data)
	if err != nil && err != io.EOF {
		// log.Println("get meta data from article file head failed:", err)
		return nil, 0
	}

	var yamlTag = []byte("---")
	var foundtag int
	var i int
	for n := 0; n < 2; n++ {
		i = bytes.Index(data[i:], yamlTag)
		if i == -1 {
			var errstr string
			if foundtag == 0 {
				// errstr = "not found yaml meta"
			} else {
				errstr = "not found yam deta end tag"
				log.Println("get meta data from article file head failed:", errstr)
			}
			return nil, 0
		}
		i += len(yamlTag)
	}

	return data[len(yamlTag):i], i
}

func ArticlesByPage(pageSize, pageNumber int) (total int, heads []MetaData) {
	s := pageSize * (pageNumber - 1)
	e := s + pageSize

	_articles := allarticles.Load().([]MetaData)
	heads = []MetaData{}
	total = len(_articles)
	if s > total {
		return
	}

	if e > len(_articles) {
		e = total
	}

	heads = _articles[s:e]
	return
}

func FindMetaData(filename string) (MetaData, bool) {
	for _, e := range allarticles.Load().([]MetaData) {
		dir := strings.TrimSuffix(e.Filename, filename)
		if dir != "" && dir[len(dir)-1] == filepath.Separator {
			dir = dir[:len(dir)-1]
		}

		if dir == configs.Setting.ArticleDir {
			return e, true
		}
	}
	return MetaData{}, false
}

func isArticleRefDirectory(path, link string) bool {
	path, _ = filepath.Abs(path)

	refdir := []string{configs.Setting.ArticleReferenceDir}
	refdir = append(refdir, configs.Setting.ArticleReferenceDirs...)
	for _, dir := range refdir {
		dir, _ = filepath.Abs(dir)
		if strings.HasPrefix(path, dir) {
			// 跳过引用资源目录
			log.Println("跳过引用资源目录", path)
			return true
		}
	}

	return false
}

func isSpecialArticle(path string) bool {
	path = strings.TrimPrefix(path, configs.Setting.AbsArticleDir)
	if path != "" && path[0] == filepath.Separator {
		path = path[1:]
	}

	var custompages []string
	if configs.Setting.CustomWebsiteFooter {
		custompages = append(custompages, configs.CustomPageFooter)
	}

	if configs.Setting.CustomWebsite404 {
		custompages = append(custompages, configs.CustomPage404)
	}

	for _, e := range custompages {
		if path == e {
			log.Println("ignore custom page:", path)
			return true
		}
	}

	return false
}

// filename1 是绝对路径
func isSameFile(filename1, filename2 string) bool {
	filename2 = strings.TrimPrefix(filename2, configs.Setting.ArticleDir)
	if filename2 != "" && filename2[0] == filepath.Separator {
		filename2 = filename2[1:]
	}

	dir := strings.TrimSuffix(filename1, filename2)
	if dir != "" && dir[len(dir)-1] == filepath.Separator {
		dir = dir[:len(dir)-1]
	}

	return dir == configs.Setting.AbsArticleDir
}
