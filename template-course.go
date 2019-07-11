package main

import (
//  "regexp"
//	"github.com/tdewolff/minify"
//	"github.com/tdewolff/minify/html"
//  "encoding/json"
  "html/template"
  "log"
  "net/http"
  "os"
  "path"
// "gopkg.in/russross/blackfriday.v2"
// "gopkg.in/yaml.v2"
// "github.com/gorilla/mux"

)



//------------CourseTemplate-------------------------------------------

func  (cs Courses) CourseTemplate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			r.URL.Path = "courses"
			Template().ServeHTTP(w, r)
			return
		}
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
                var err error
		lp := path.Join("deploy/templates", "layout.html")
		fp := path.Join("deploy/courses", r.URL.Path)
		info, err := os.Stat(fp)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("404: %s", r.URL.Path)
				errorHandler(http.StatusNotFound).ServeHTTP(w, r)
				return
			}
		}
		if info.IsDir() {
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles(lp,  fp)
		if err != nil {
			log.Printf("Template Error: %s", err)
		}
		err = tmpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			log.Printf("Template Error: %s", err)
		}
		return
	})
}


