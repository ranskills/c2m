package util

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"strings"
	"time"
)

func WatchDirectory(dir string, processFile func(filePath string)) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					time.Sleep(1 * time.Second)
					processFile(event.Name)
				}
				//fmt.Printf("EVENT! %#v\n", event)

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	if err := watcher.Add(dir); err != nil {
		fmt.Println("ERROR", err)
	} else {
		fmt.Printf("Watching the directory %s for new files\n", dir)
		fmt.Println(strings.Repeat("-", 80))
	}

	<-done
}
