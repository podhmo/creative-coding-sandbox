package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type App struct {
	tmpl *template.Template
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "text/html")
	fpath := strings.TrimPrefix(r.URL.Path, "/")
	switch fpath {
case "":
	err = app.List(w, r)
default:
	err = app.Draw(w, r, fpath)
}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (app *App) List(w http.ResponseWriter, r *http.Request) error {
	files, err := listFiles("*.js")
	if err != nil {
		return err
	}
	return app.tmpl.ExecuteTemplate(w, "list.tmpl", map[string]interface{}{
"Files": files,
})
}

func (app *App) Draw(w http.ResponseWriter, r *http.Request, filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
}
	return app.tmpl.ExecuteTemplate(w, "draw.tmpl", map[string]interface{}{
"Code": string(b),
})
}
func listFiles(globPatten string) ([]string, error) {
	var files []string
	targets, err := filepath.Glob(globPatten)
	if err != nil {
		return nil, err
	}
	for _, fname := range targets {
		if strings.ToUpper(fname[:1]) != fname[:1] {
			continue
		}
		files = append(files, fname)
	}
	return files, nil
}

var tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))

func main() {
	port := 4444
	if v, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		port = v
	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	app := &App{tmpl: tmpl}
	log.Println("listen ...", addr)
	http.ListenAndServe(addr, app)
}
