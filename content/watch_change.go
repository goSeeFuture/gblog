package content

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/goSeeFuture/gblog/configs"
)

func watchArticleChange() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("cannot start article watcher:", err)
		return
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("article change event:", event)
				log.Println("start reload content")
				err := loadContent()
				if err != nil {
					log.Println("reload content failed:", err)
				} else {
					log.Println("reload content ok")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("article watcher error:", err)
			}
		}
	}()

	err = watcher.Add(configs.Setting.ArticleDir)
	if err != nil {
		log.Println("cannot watch article dir:", err)
		return
	}
	<-done
}
