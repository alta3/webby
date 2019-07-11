package main

import (
//  "regexp"


//	"github.com/tdewolff/minify"
//	"github.com/tdewolff/minify/html"

  "encoding/json"
  "fmt"
  "net/http"
// "gopkg.in/russross/blackfriday.v2"

// "gopkg.in/yaml.v2"
// "github.com/gorilla/mux"

)




//  __  __ ______ _____          __  __ ______ _   _ _    _ 
// |  \/  |  ____/ ____|   /\   |  \/  |  ____| \ | | |  | |
// | \  / | |__ | |  __   /  \  | \  / | |__  |  \| | |  | |
// | |\/| |  __|| | |_ | / /\ \ | |\/| |  __| | . ` | |  | |
// | |  | | |___| |__| |/ ____ \| |  | | |____| |\  | |__| |
// |_|  |_|______\_____/_/    \_\_|  |_|______|_| \_|\____/

//API - MegaMenu - Returns all courses: Id, Tag, Stars, and Name. 
func (cs Courses ) getmegamenu() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    type MenuItems struct {
      Id           string          `json:"id"`
      CourseTitle  string          `json:"course-title"`
      Tags         []string        `json:"tags"`
    }
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    mi := MenuItems{} 
    mis := []MenuItems{} 
    var js []byte
    var err error
    //Iterate over all courses, copy Id, Name, and any Tags
    for _, ThisCourse := range cs.Cc {
       fmt.Println("--------------------------------------------------")
       fmt.Printf("Menu Item = %s, %s, %s\n", ThisCourse.Id, ThisCourse.CourseTitle, ThisCourse.Tags)
       mi.Id=ThisCourse.Id
       mi.CourseTitle=ThisCourse.CourseTitle
       mi.Tags=ThisCourse.Tags
       mis = append(mis,mi)
    }
    // Marshal the Mega-Menu
    js, err = json.Marshal(mis)
    if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
       fmt.Printf("Error %s:\n", err)
       return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
    return
    })
}


