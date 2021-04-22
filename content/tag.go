package content

type Tag struct {
	Name  string
	Count int
}

func makeTagIndex(md []MetaData) (map[string][]int, []Tag) {
	var tags []Tag
	index := make(map[string][]int)
	for i, e := range md {
		for _, tag := range e.Tags {
			array := index[tag]
			array = append(array, i)
			_, exist := index[tag]
			index[tag] = array
			if !exist {
				tags = append(tags, Tag{Name: tag})
			}
		}
	}

	for tag, e := range index {
		for i, et := range tags {
			if et.Name == tag {
				et.Count = len(e)
				tags[i] = et
				break
			}
		}
	}
	return index, tags
}

func GetTags() []Tag {
	var _tags []Tag
	cotentmutex.RLock()
	_tags = tags
	cotentmutex.RUnlock()
	return _tags
}

func ArticlesByTagPage(tag string, pageSize, pageNumber int) (total int, heads []MetaData) {
	s := pageSize * (pageNumber - 1)
	e := s + pageSize

	var slice []int
	var _articles []MetaData
	cotentmutex.RLock()
	slice = tagIndex[tag]
	_articles = articles
	cotentmutex.RUnlock()

	total = len(slice)

	if s >= total {
		return total, []MetaData{}
	}

	for i, index := range slice {
		if i >= s && i < e {
			heads = append(heads, _articles[index])
		}
	}
	return
}
