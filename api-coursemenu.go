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

// Course Menu
// Returns a course menu that will be easy for the front end js to implement

type CourseMenu  []MiniMenu    // `json:"coursemenu"`

type MiniMenu struct {
  MiniMenuTitle  string          `json:"mini-menu"`
  MmItems        []MmItem        `json:"mm-items"`
}

type MmItem struct {
  Id              string         `json:"id"`
  CourseTitle     string         `json:"course-title"`
}

func (cc Courses)  menumaker() CourseMenu  {
    var minimenutitles  []string
    var mmitem          MmItem
    var cm              CourseMenu
    var existing        bool
    var match           bool
    var mm              MiniMenu
    var mmitems         []MmItem
    var id              string
    var coursetitle     string
    // Iterate over all courses.tags[], and derive a list of unique tags called MiniMenuTitles
    for _, thiscourse := range cc.Cc {
        // Iterate over all tags within a course
        for _, thistag := range thiscourse.Tags {
             existing = false
             //Iterate over all minimenutitle, add new ones to the list
             for _, thismmt := range minimenutitles {
                 if thistag == thismmt {
                    existing = true
                 }
             }
             //No matches? Then add thistag to the list
             if existing == false {
             minimenutitles = append (minimenutitles, thistag)
             }
        }
    }
    fmt.Println("Items: %+v\n\n", minimenutitles)
    // Interate over each MiniMenuTitle 
    for _, mmt := range  minimenutitles {
          //Iternate over every course in the catalog
          for _, thiscourse := range cc.Cc {
                match = false
                //Iterate over every course's Tags
                for _, thistag := range thiscourse.Tags {
                   //If this course has a matching tag entry, grab data
                   if mmt == thistag {
                      match = true
                      id = thiscourse.Id
                      coursetitle = thiscourse.CourseTitle
                      break
                   }
                }//End of tag iteration so add item if match is true
                if match == true {
                   mmitem.Id = id
                   mmitem.CourseTitle = coursetitle
                   mmitems = append(mmitems, mmitem)
                }
          } //End of Course Itermation
          mm.MiniMenuTitle = mmt
          mm.MmItems = mmitems
          fmt.Printf("mm: %v\n", mm)
          cm = append (cm, mm)
          mmitems = nil
   }//End of minimenu iteration
    return cm
}


//API - New MegaMenu - Returns a canned megamenu pre-sorted by course type. 
func (cm CourseMenu ) coursemenu() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    var js []byte
    var err error
    js, err = json.Marshal(cm)
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


