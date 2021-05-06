---
Tags:
    - gblog
    - 使用指南
UpdateAt: 2021-04-12 00:44:00
---
# gblog使用说明

[gblog](https://github.com/goSeeFuture/gblog/) 是一款基于Markdown文件的轻量博客。它的使用和维护都比较简单，无需编程知识，只需要在`articles`目录中放入Markdown文件，然后启动服务即可。

## 快速开始

### 启动博客

执行 ./gblog 文件，启动博客服务

```bash
# gblog 博客服务目录结构
.
├── articles                                    # 将你的Markdown文件放在此目录中
│   └── gblog使用说明.md           # 本教程文章
├── config.toml                            # 默认配置
├── static                                      # 博客网站使用的静态资源
│     ├── css                                   # 博客用到的层叠样式表
│     │     └── bulma.min.css        # 默认采用bulma库
│     └── image                             # 网站用的图片
│            ├── favicon.ico               # 网站收藏图标
│            └── logo.png                   # 网站logo
└── gblog                                      # gblog服务执行文件

# 启动服务命令
./gblog

# 命令执行结果：
 # ┌───────────────────────────────────────────────────┐ 
 # │                    Fiber v2.8.0                   │ 
 # │               http://127.0.0.1:3000               │ 
 # │                                                   │ 
 # │ Handlers ............ 16  Processes ........... 1 │ 
 # │ Prefork ....... Disabled  PID ........... 1473109 │ 
 # └───────────────────────────────────────────────────┘ 
```

可以看到 `http://127.0.0.1:3000` 就是博客网页地址，访问它就能打开你的博客，就这么简单。

![快速开始结果](/docs/assets/image/quickstart.png)

在`./gblog`

### 文章中引用图片

让我们丰富以下文章，加入一张图片，Markdown引用图片如下：

```md
![快速开始结果](/docs/assets/image/quickstart.png)
```

打开我们的博客首页，按F5键刷新

> 访问 http://127.0.0.1:3000 即博客首页

![没有显示图片](/docs/assets/image/nopicture.png)

结果如上图，图片并没有显示出来，这是因为我们还没有指定文章中引用图片的位置，接下来我们指定它。

1. 在`articles`目录下建立`assets/image`目录
2. 将`quickstart.png`文件放入到`articles/assets/image`中

再次访问博客，可以看到图片成功的显示出来了

![引入图片成功](/docs/assets/image/pictureok.png)

此时，articles的目录结构如下：

```bash
articles
├── assets
│   └── image
│       └── quickstart.png
└── gblog使用说明.md
```

### 增加分类

当你的博客有许多不同类别的文章时，希望能进行分类整理，便于博客访问者找到自己感兴趣的那类文章。
在gblog中，你只需要建立以类别命名的目录，同时把相关文章放进去，博客顶部导航中就会显示出分类。

比如，将我刚刚写的文章放到 `articles/gblog使用指南` 目录下，然后刷新博客页面

![新建分类](/docs/assets/image/newcategory.png)

如图，点击某个分类，gblog博客就会罗列出分类中的所有文章。大家发现，除了我们新建的分类外，还有一个
名叫`未分类`的项，这个项显示的是`articles`目录下的所有文章，当我们博客文章数量较少时，用不到分类，
直接把文章放进`articles`目录，此时访问首页和`未分类`都可以列出所有文章。

```bash
# 新建`gblog使用指南`分类后的目录结构
articles
├── assets
│   └── image
│       ├── nopicture.png
│       ├── pictureok.png
│       └── quickstart.png
└── gblog使用指南
    └── gblog使用说明.md
```

## 增加标签

博客中的一篇文章，可能存在多个小类，
比如我们有一篇文章《精神的健康取决于情感》，它跟心理、情感、健康都有关系，
用上面说的分类法就很难将它归类，这时我们可以通过给文章增加标签来达到目的。
具体做法如下：

在文章顶部中增加“头数据”

```yaml
---
Tags:           # 为此文章设定标签，可以是1个或者多个
    - 心理      # 设定“心理”标签，行首有1个tab缩进，1个“-”字符，以及1个空格
    - 情感      # 设定“情感”标签
    - 健康       # 设定“健康”标签
---

# 精神的健康取决于情感

这是一篇有关心理、情感、健康的文章。
```

**注意**：“头数据”前后，必须有一行“---”gblog才能正确识别。

将文章放入`articles`目录，刷新页面访问[该文章](http://localhost:3000/articles/%e7%b2%be%e7%a5%9e%e7%9a%84%e5%81%a5%e5%ba%b7%e5%8f%96%e5%86%b3%e4%ba%8e%e6%83%85%e6%84%9f.md)可以看到

![新增标签](/docs/assets/image/newtag.png)

在文章右上角，出现我们新增的3个标签，点击某个标签，就能找到拥有同样标签的文章。另外，也可以通过导航栏中的“标签”菜单，查看所有拥有此标签的文章。

![当行栏中的标签菜单](/docs/assets/image/navtags.png)

此时目录结构

```md
articles
├── 精神的健康取决于情感.md
├── assets
│   └── image
│       ├── newcategory.png
│       ├── newtag.png
│       ├── nopicture.png
│       ├── pictureok.png
│       └── quickstart.png
└── gblog使用指南
    └── gblog使用说明.md
```

可以看到这篇文章放在`articles`目录下，点击导航栏中【分类】->【未分类】可以看出它属于“未分类”。

## 修改配置

通过修改配置`config.toml`可以定制博客外观，修改需要将在重启gblog服务后生效。

### 设置站点名称

使用文本编辑器打开 `config.toml` 文件，修改`WebsiteName`后面的值，比如：

```toml
# 站点名称
WebsiteName = "牛莫旺的窝窝"
```

**注意**：该值需要用英文双引号包裹起来。

嗯，你发现了这个文件里面还有许多其他的配置，如果你能读明白注释，可以尝试自行配置其他项。

> 在配置文件config.toml中，#号以及后面的文字都被视为注释，gblog会忽略它们。

### 自定义页脚

设置`CustomWebsiteFooter`项的的值为`true`，然后重启gblog服务。

```toml
# 网站页脚内容
# 需要实现`ArticleDir`/footer.md 
CustomWebsiteFooter = true
```

新建文件`articles/footer.md`，并写入页脚显示的内容，比如：

```md
欢迎您访问 **[牛莫旺的窝窝](http://youdomain.com)** ，想寻求合作请联系 <email>niumowang@qq.com</email>
```

### 设置404页

设置`CustomWebsite404`项的的值为`true`，然后重启gblog服务。

```toml
# 网站无法找到页面
# 需要实现`ArticleDir`/404.md
CustomWebsite404 = true
```

新建文件`articles/404.md`，并写入内容，比如：

```md
## 抱歉，无法找到该页面！
```

### 修改Logo与收藏图标

- 修改Logo文件，用你的Logo文件`logo.png`替换掉`static/image/logo.png`。
- 修改收藏图标，用你的图标收藏文件`favicon.ico`替换掉`static/image/favicon.ico`。
- 完成后重启gblog服务，修改将会生效。

### 修改文章内代码样式

`ArticleCodeStyle`设置代码高亮风格，取值请参见[此网页](https://xyproto.github.io/splash/docs/all.html)

`ArticleCodeShowLineNumber`设置是否显示行号，true显示，false隐藏

```toml
# 文章内代码语法高亮风格
# 取值范围：https://xyproto.github.io/splash/docs/all.html
ArticleCodeStyle = "perldoc"
# 文章内代码显示行号
ArticleCodeShowLineNumber = true
```

需重启服务后修改才能设置生效。

### 文章支持MathJax

文章支持显示漂亮的数学公式，通过开启MathJax

```toml
# 文章支持MathJax，详情参见：https://github.com/litao91/goldmark-mathjax
ArticleMathJax = true # 取消此行注释，重启服务后生效
```

MathJax支持LaTeX语法，在Markdown使用格式如下：

```md
$$
/<你的LaTeX内容/>
$$
```

可以看出仅需要在LaTeX内容首尾各增加一行“$$”即可。

行内的公式，可以使用内联语法，即在内容前后各增加“$”即可。如下：

```md
仅剩$\frac{1}{2}$个梦想
```

重启gblog服务，生效后的效果

![内联LaTeX内容](/docs/assets/image/inlinelatex.png)

### 设置分类别名

gblog使用目录相对路径作为分类名，当目录层级太深时，你也许想要给长长名字取个短点的别名。
比如我们有个`articles/我/的/名字/很长/啊.md`文件，要修改它所在的分类`我/的/名字/很长`，
为`短类名`，目录结构如下

```bash
articles/我
└── 的
    └── 名字
        └── 很长
            └── 啊.md
```

![长类别名](/docs/assets/image/longcategory.png)

只需要在`config.toml`末尾加上如下几行配置

```toml
[[Category]]
Name = "短类名"
Path = "我/的/名字/很长"
```

重启gblog服务，生效后的效果

![修改类别名称](/docs/assets/image/changecategory.png)

### 草稿

gblog 可以将未完成的文章发布未草稿，被设定为草稿的文章，需要作者权限才能阅读全文。
你可以先将草稿发布到博客中，作为预告，待完成后，再对所有人开放阅读全文权限。

接下来我将演示如何创建一篇草稿文章。

#### 配置对外域名

打开配置文件`config.toml`，并设定好`PublicDomain`项的值。

```toml
# 公开域名，不使用草稿可不填写
# 比如：PublicDomain = "http://www.youdomain.com"
PublicDomain = "http://localhost:3000"
```

接下来，重启服务，让配置生效。

#### 标记草稿文章

例如，我们创建一篇文章 `articles/这是一篇预发布的文章.md`文章最前面，写入以下内容

```md
---
Draft: true
---
```

这样，gblog将会把此文章识别为草稿文章。我们刷新博客，将看到如下文章列表

![列表中的草稿文章](/docs/assets/image/listdraft.png)

点击文章标题，将提示如下

![无法浏览草稿文章](/docs/assets/image/cannotviewdraft.png)

这是因为gblog不允许读者访问草稿文章的全部内容，作为作者的你需要阅读这篇草稿文章的全部内容，只需要进行如下操作

查看`gblog`服务的启动日志，访问预览草稿链接

![预览草稿链接](/docs/assets/image/authorlink.png)

```
二维码与链接是同样内容，你可以使用手机扫描二维码获得作者权限。
```

然后再次访问草稿文章即可显示文章全部内容

![草稿全文](/docs/assets/image/viewdraft.png)

## 为分类添加简介

gblog支持为分类增加简短的文字介绍，要使用该功能需要在配置文件中设置如下项

```toml
# 为分类书写主题介绍
# 需要实现`ArticleDir`/分类/topic.md
CategoryTopic = true
```

    重启服务之后生效

比如，要给`gblog使用指南`分类添加简介，需要建立`articles/gblog使用指南/topic.md`文件，并写入内容

```toml
# gblog

简绍gblog博客系统的基础用法和使用技巧。
```

在导航栏上选择`gblog使用指南`分类，可以看到我们刚写的简介

![分类简介](/docs/assets/image/categorytopic.png)
