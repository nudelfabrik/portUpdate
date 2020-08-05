package templates

import (
	"html/template"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

var fileNames = []string{
	"header.html",
	"footer.html",
	"list.html",
}

func WatchTemplates(tch chan *template.Template) *template.Template {
	path := "../../server/templates/tmpls/"
	tmpl := template.New("")
	files := []string{}
	for _, f := range fileNames {
		files = append(files, path+f)

	}
	tmpl = template.Must(template.ParseFiles(files...))
	watcher, err := fsnotify.NewWatcher()
	watcher.Add(path)
	if err != nil {
		log.Println(err)
		return tmpl
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// In case the editor deletes and then rewrites the file, the next call panics, because the files are missing until rewritten.
				time.Sleep(100 * time.Millisecond)
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("Reloading Templates")
					tmpl = template.New("")
					tmpl = template.Must(template.ParseFiles(files...))
					tch <- tmpl
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	return tmpl
}
