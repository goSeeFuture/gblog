package configs

var (
	Setting Config
)

// 服务配置
type Config struct {
	// 监听地址
	Listen string
	// 文章目录
	ArticleDir string
	// 文章目录绝对路径
	AbsArticleDir string `toml:"-"`
	// 文章中引用资源路径，比如图片
	ArticleReferenceDir string
	// 文章中引用资源路径，比如图片
	ArticleReferenceDirs []string
	// 头部yaml长度
	MaxMetaDataLen int
	// 文章分类别名
	Categories []Category `toml:"Category"`
	// 站点名称
	WebsiteName string
	// 显示未分类
	ShowUncategorized bool
	//每页显示文章数
	PageSize int
	// 分页前缀，用于避免与ArticleDir中的目录命名冲突
	// 不能使用url保留字符
	PageNumberPrefix string
	// 文章内代码语法高亮风格
	// 取值范围：https://xyproto.github.io/splash/docs/all.html
	ArticleCodeStyle string
	// 文章内代码显示行号
	ArticleCodeShowLineNumber bool
	// 文章支持MathJax，详情参见：https://github.com/litao91/goldmark-mathjax
	ArticleMathJax bool
	// 网站页脚内容
	// 需要实现`ArticleDir`/footer.md
	CustomWebsiteFooter bool
	// 网站无法找到页面
	// 需要实现`ArticleDir`/404.md
	CustomWebsite404 bool
}

// 文章分类别名
type Category struct {
	// 分类名称，默认取Path
	Name string
	// 所在路径，基于ArticleDir的相对路径
	Path string
	// Path的变型
	ID string `toml:"-"`
	// 文章数量
	Count int `toml:"-"`
}
