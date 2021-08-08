package content

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/goSeeFuture/gblog/configs"
	"github.com/goSeeFuture/gblog/pkg/toc"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/beevik/etree"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type HTMLPage struct {
	Heads   []toc.Head
	Content []byte
}

func MarkdownPage(filename string, offset int) (*HTMLPage, error) {
	fs, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	versions := articleVersion.Load().(map[string]int64)
	if fs.ModTime().Unix() != versions[filename] {
		// 重新解析文章MetaData
		md, ok := reloadArticleMetaData(filename)
		if !ok {
			return nil, errors.New("reload article meta data failed")
		}
		offset = md.Offset
		log.Println("reload article meta data", filename)
	}

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	f = f[offset:]

	f, _ = removeMarkdownH1(f)
	c, h := markdown2HTML(f)

	return &HTMLPage{Heads: h, Content: c}, nil
}

func markdown2HTML(data []byte) ([]byte, []toc.Head) {
	extensions := []goldmark.Extender{
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle(configs.Setting.ArticleCodeStyle),
			highlighting.WithFormatOptions(
				html.WithLineNumbers(configs.Setting.ArticleCodeShowLineNumber),
			),
		),
		emoji.Emoji,
		extension.Footnote,
		extension.TaskList,
	}

	if configs.Setting.ArticleMathJax {
		extensions = append(extensions, mathjax.MathJax)
	}

	md := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		))

	var buf bytes.Buffer
	slugID := toc.NewSlugID()
	context := parser.NewContext(parser.WithIDs(slugID))
	if err := md.Convert(data, &buf, parser.WithContext(context)); err != nil {
		return nil, nil
	}

	return buf.Bytes(), slugID.(*toc.SlugID).Heads()
}

func removeMarkdownH1(data []byte) ([]byte, string) {
	if len(data) == 0 {
		return nil, ""
	}

	data = bytes.TrimLeft(data, "\r\n\t ")
	i := bytes.Index(data, []byte{'\n'})
	if i == -1 || !regexp.MustCompile(`^#\s`).Match(data[:i]) {
		return data, ""
	}

	return bytes.TrimLeft(data[i+1:], "\r\n\t "), string(data[2:i])
}

func getH1(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	data = bytes.TrimLeft(data, "\r\n\t ")
	i := bytes.Index(data, []byte{'\n'})
	if i == -1 || !regexp.MustCompile(`^#\s`).Match(data[:i]) {
		return ""
	}

	return string(data[2:i])
}

func removeH1(data []byte) ([]byte, string) {
	if len(data) == 0 {
		return nil, ""
	}

	doc := etree.NewDocument()
	err := doc.ReadFromBytes(data)
	if err != nil {
		return nil, ""
	}

	var h1 string
	node := doc.FindElement("//h1")
	if node != nil {
		h1 = node.Text()
		if node.Parent() != nil {
			node.Parent().RemoveChild(node)
		} else {
			doc.RemoveChild(node)
		}
		data, err = doc.WriteToBytes()
		if err != nil {
			return nil, ""
		}
	}

	return data, h1
}
