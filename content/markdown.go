package content

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/goSeeFuture/gblog/configs"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/beevik/etree"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func MarkdownPage(filename string, offset int) ([]byte, error) {
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

	return markdown2HTML(f), nil
}

func markdown2HTML(data []byte) []byte {
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

	md := goldmark.New(goldmark.WithExtensions(extensions...), goldmark.WithParserOptions(parser.WithAttribute()))

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(data, &buf, parser.WithContext(context)); err != nil {
		return nil
	}

	return buf.Bytes()
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
