package content

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
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
}

// 罗列所有文章的Meta头
func List() (articles []MetaData) {
	err := filepath.Walk(configs.Setting.ArticleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if path == configs.Setting.WebsiteFooter {
			// log.Println("ignore website footer:", path)
			return nil
		}

		if isArticleRefDirectory(path) {
			// 跳过引用资源目录
			return nil
		}

		if strings.ToLower(filepath.Ext(info.Name())) != ".md" {
			return nil
		}

		dir := strings.TrimPrefix(filepath.Dir(path), "articles/")
		md := MetaData{
			Title:      info.Name(),
			Filename:   path,
			CategoryID: strings.Replace(strings.TrimPrefix(dir, "/articles/"), "/", "-", -1),
		}
		meta, offset := getMetaData(path)
		if meta != nil {
			err = yaml.Unmarshal(meta, &md)
			if err != nil {
				return err
			}
			md.HasMetaHead = true
		}

		head, h1, modtm := getHeadContent(path, offset)
		if len(head) != 0 {
			md.Summary = template.HTML(head)
		}
		if h1 != "" {
			md.InlineTitle = true
			md.Title = h1
		}
		if md.UpdateAt.IsZero() && !modtm.IsZero() {
			md.UpdateAt = modtm
		}

		md.Offset = offset
		articles = append(articles, md)
		return nil
	})
	if err != nil {
		log.Println("foreach article failed:", err)
	}
	return
}

func getHeadContent(filename string, offset int) ([]byte, string, time.Time) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("open article file failed:", err)
		return nil, "", time.Time{}
	}
	defer file.Close()

	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		log.Println("read seek article file failed:", err)
		return nil, "", time.Time{}
	}

	const headPartSize int64 = 256
	var readsize int64
	fs, err := file.Stat()
	if err != nil {
		log.Println("get article file head failed:", err)
		return nil, "", time.Time{}
	}

	var isPart = fs.Size() >= headPartSize
	if isPart {
		readsize = headPartSize
	} else {
		readsize = fs.Size()
	}

	part := make([]byte, readsize)
	partsize, err := file.Read(part)
	if err != nil && err != io.EOF {
		log.Println("get article file head failed:", err)
		return nil, "", time.Time{}
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

	if isPart {
		part = append(part, []byte("...")...)
	}

	header, h1 := removeH1(markdown2HTML(part))
	return header, h1, fs.ModTime()
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
	if fs.Size() < int64(configs.Setting.MaxMetaDataLen) {
		data = make([]byte, fs.Size())
	} else {
		data = make([]byte, configs.Setting.MaxMetaDataLen)
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

	var _articles []MetaData
	cotentmutex.RLock()
	_articles = articles
	cotentmutex.RUnlock()

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
	var _articles []MetaData
	cotentmutex.RLock()
	_articles = articles
	defer cotentmutex.RUnlock()
	for _, e := range _articles {
		if e.Filename == filename {
			return e, true
		}
	}
	return MetaData{}, false
}

func isArticleRefDirectory(path string) bool {
	if strings.HasPrefix(path, configs.Setting.ArticleReferenceDir) {
		// 跳过引用资源目录
		return true
	}

	for _, dir := range configs.Setting.ArticleReferenceDirs {
		if strings.HasPrefix(path, dir) {
			// 跳过引用资源目录
			return true
		}
	}

	return false
}
