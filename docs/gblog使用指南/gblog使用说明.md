---
Tags:
    - gblog
    - 使用指南
UpdateAt: 2021-04-12 00:44:00
---
# gblog使用说明

[gblog](https://github.com/goSeeFuture/gblog) 是一款基于Markdown文件的轻量博客。它的使用和维护都比较简单，无需编程知识，只需要在`articles`目录中放入Markdown文件，然后启动服务即可。

## 快速开始

### 启动博客

```bash
# gblog 博客服务目录结构
.
├── articles                                    # 将你的Markdown文件放在此目录中
│   └── gblog使用说明.md           # 本教程文章
├── config.toml                            # 默认配置
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

![快速开始结果](docs/assets/image/quickstart.png)

### 文章中引用图片

让我们丰富以下文章，加入一张图片，Markdown引用图片如下：

```md
![快速开始结果](/image/quickstart.png)
```

重新启动服务器，按`Ctrl+C`终止`./gblog`运行，并再次运行它。然后打开我们的博客首页

> 访问 http://127.0.0.1:3000 即博客首页

![没有显示图片](/image/nopicture.png)

结果如上图，图片并没有显示出来，这是因为我们还没有指定文章中引用图片的位置，接下来我们指定它。

1. 在`articles`目录下建立`assets/image`目录
2. 将`quickstart.png`文件放入到`articles/assets/image`中
3. 再次重启服务器

再次访问博客，可以看到图片成功的显示出来了

![引入图片成功](/image/pictureok.png)

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

比如，将我刚刚写的文章放到 `articles/gblog使用指南` 目录下，然后重启服务

![新建分类](/image/newcategory.png)

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

将文章放入`articles`目录，重启服务后，访问[该文章](http://localhost:3000/articles/%e7%b2%be%e7%a5%9e%e7%9a%84%e5%81%a5%e5%ba%b7%e5%8f%96%e5%86%b3%e4%ba%8e%e6%83%85%e6%84%9f.md)可以看到

![新增标签](/image/newtag.png)

在文章右上角，出现我们新增的3个标签，点击某个标签，就能找到拥有同样标签的文章。另外，也可以通过导航栏中的“标签”菜单，查看所有拥有此标签的文章。

![当行栏中的标签菜单](/image/navtags.png)

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
