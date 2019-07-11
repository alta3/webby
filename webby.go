package main

import (
//  "regexp"
//	"github.com/tdewolff/minify"
//	"github.com/tdewolff/minify/html"
//  "encoding/json"
  "fmt"
  "log"
  "net/http"
// "gopkg.in/russross/blackfriday.v2"
// "gopkg.in/yaml.v2"
// "github.com/gorilla/mux"

)



// ------------------------------------------------------------

func main() {
////      router := mux.NewRouter().StrictSlash(true)
  cs := Load()
  fmt.Println("--------------------------------------------------")
  fmt.Println("Course Loaded into MAIN: " + cs.Cc[0].Id)
  events := LoadEvents()
  blogcontent := Loadblogs()
  fmt.Println("Blogs Loaded into MAIN: " + blogcontent[0].Id)
  cmenu :=  cs.menumaker()
  fmt.Println("Menu: %+v\n", cmenu)
  fmt.Printf("BLOGMENU\n")
  fmt.Println("--------------------------------------------------")
  bmenu := blogcontent.blogmenumaker()
  fmt.Println("--------------------------------------------------")
//	fmt.Printf("blogmenu: %s\n", bmenu[0].BlogCategory)


// JSON RESTful Interfaces
  http.Handle("/api/v1/course/search/",      http.StripPrefix("/api/v1/course/search/",     cs.search()))
  http.Handle("/api/v1/course/megamenu/",    http.StripPrefix("/api/v1/course/megamenu/",   cs.getmegamenu()))
  http.Handle("/api/v1/course/coursemenu/",  http.StripPrefix("/api/v1/course/coursemenu/", cmenu.coursemenu()))
  http.Handle("/api/v1/course/summary/id/",  http.StripPrefix("/api/v1/course/summary/id/", cs.getsummary()))
  http.Handle("/api/v1/course/detail/id/",   http.StripPrefix("/api/v1/course/detail/id/",  cs.getcoursedetail()))
  http.Handle("/api/v1/blog/search/",        http.StripPrefix("/api/v1/blog/search/",       blogcontent.blogsearch()))
  http.Handle("/api/v1/blog/id/",            http.StripPrefix("/api/v1/blog/id/",           blogcontent.getblog()))
  http.Handle("/api/v1/events/",             http.StripPrefix("/api/v1/events/",            events.getevents()))
  http.Handle("/api/v1/events/menu/",        http.StripPrefix("/api/v1/events/menu",        events.geteventsmenu()))
  http.Handle("/api/v1/blog/menu/",          http.StripPrefix("/api/v1/blog/menu/",         bmenu.blogmenu()))


// All the static folders
	http.Handle("/downloads/",   http.StripPrefix("/downloads/",   http.FileServer(http.Dir("downloads"))))
	http.Handle("/img/",         http.StripPrefix("/img/",         http.FileServer(http.Dir("deploy/img"))))
	http.Handle("/images/",      http.StripPrefix("/images/",      http.FileServer(http.Dir("deploy/images"))))
	http.Handle("/fonts/",       http.StripPrefix("/fonts/",       http.FileServer(http.Dir("deploy/fonts"))))
	http.Handle("/icons/",       http.StripPrefix("/icons/",       http.FileServer(http.Dir("deploy/icons"))))
	http.Handle("/css/",         http.StripPrefix("/css/",         http.FileServer(http.Dir("deploy/css"))))
	http.Handle("/js/",          http.StripPrefix("/js/",          http.FileServer(http.Dir("deploy/js"))))
  http.Handle("/assets/",      http.StripPrefix("/assets/",      http.FileServer(http.Dir("deploy/assets"))))
  http.Handle("/coming_soon/", http.StripPrefix("/coming_soon/", http.FileServer(http.Dir("deploy/coming_soon"))))
  http.Handle("/sass/",        http.StripPrefix("/sass/",        http.FileServer(http.Dir("deploy/sass"))))
  http.Handle("/video/",       http.StripPrefix("/video/",       http.FileServer(http.Dir("deploy/video"))))
  http.Handle("/layerslider/", http.StripPrefix("/layerslider/", http.FileServer(http.Dir("deploy/layerslider"))))
//http.Handle("/",             http.StripPrefix("/",             http.FileServer(http.Dir("deploy"))))

	// Templates
	http.Handle("/courses/", http.StripPrefix("/courses/", cs.CourseTemplate()))
	http.Handle("/blog/", http.StripPrefix("/blog/", BlogTemplate()))
	http.Handle("/", http.StripPrefix("/", Template()))

// Stripe Chckout
	http.Handle("/checkout", Checkout())

	log.Printf("serving...")
	http.ListenAndServe(":8888", Log(http.DefaultServeMux))
}

