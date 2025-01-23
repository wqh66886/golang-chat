package templates

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type TemplateHandler struct {
	once     sync.Once
	FileName string
	template *template.Template
}

func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.FileName)))
	})
	err := t.template.Execute(w, r)
	if err != nil {
		log.Println("template load error: ", err)
		return
	}
}
