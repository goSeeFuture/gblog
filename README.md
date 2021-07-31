# gblog

gblog是一个基于Markdown文档构建个人博客的服务，其设计意图是，让使用者付出尽可能少的理解成本，自建博客，使用个人博客。
为了达成愿景，gblog假设使用者只会编辑Markdown文件，和基本的电脑知识，包括：知道文件、目录、编辑文本文件、运行程序。

![gblog工作流程](/docs/assets/image/gblog_process.png)

## 功能介绍

- 基于Markdown文档的博文
- 基于目录路径的博文分类
- 博文多支持标签，查询包含指定标签的博文
- 监视博文目录，动态加载变更
- 采用Markdown自定义页脚、404页
- 支持为分类添加简介
- 支持草稿，只有作者可以阅读草稿全部内容
  
## 使用说明

也许你从未自己搭建过博客，有或许你只是个golang后端开发者，想找个简单的博客自己使用并维护，
我这里刚好有两份使用说明，方便不同需求的朋友阅读。

- [gblog使用指南，适合最萌最萌的萌新](/docs/gblog使用指南/gblog使用说明.md)
- [gblog设计说明，想要飙车的老鸟](?)

建议老鸟也看看`gblog使用指南`，了解一下gblog都有那些功能。
也许你就是想看看热闹，那也不妨来颗 ⭐吧。

## 特别鸣谢

特别感谢：

- [glodmark](https://github.com/yuin/goldmark)
- [fiber](https://github.com/gofiber/fiber)
- [bulma](https://github.com/jgthms/bulma)
- [go-toml](https://github.com/pelletier/go-toml)
- [etree](https://github.com/beevik/etree)
- [chroma](https://github.com/alecthomas/chroma)
- [fsnotify](https://github.com/fsnotify/fsnotify)

正是因为你们的的付出，才有了这个博客项目，谢谢！
