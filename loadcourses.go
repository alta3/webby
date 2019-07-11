package main

import (
//  "regexp"


//	"github.com/tdewolff/minify"
//	"github.com/tdewolff/minify/html"

  "encoding/json"
  "fmt"
  "path/filepath"
  "log"
  "os"
  "path"
  "io/ioutil"
  "strconv"
// "gopkg.in/russross/blackfriday.v2"

// "gopkg.in/yaml.v2"
// "github.com/gorilla/mux"

)

//Load COURSES
//------------------------------------------------------------
func Load() Courses {
      // Create a OS compliant path: microsoft "\" or linux "/"
      dirname := path.Join("deploy", "courses")
      d, err := os.Open(dirname)
      if err != nil {
          log.Printf("No courses directory! %s" , err)
          os.Exit(1)
      }
      // If n > 0, Readdirnames(n) returns at most n names
      // If n < 0, Readdirnames(n) returns ALL names
      n := -1
      // reads < n > files in directory < d >
      filenames, err := d.Readdirnames(n)
      if err != nil {
          log.Printf("No files in course directory! %s\n" , err)
          os.Exit(1)
      }
      c := make([]Course,50)
      var jsonCatalogFile Courses
      fmt.Println("--------------------------------------------------")
      fmt.Printf(" Reading files in this directory: %s\n", dirname)
      i := 0
      for _, filename := range filenames {
          thisfile := path.Join(dirname, filename)
          _ , err := os.Stat(thisfile)
          if err != nil {
              if os.IsNotExist(err) {
                  log.Printf("file is missing!: %s\n ", filename)
              }
          } 
          if filepath.Ext(thisfile) == ".yaml" {
              yammy, err := ioutil.ReadFile(thisfile)
              if err != nil {
                 log.Printf("yammy.Get err: %s\n", err)
              }
              fmt.Printf("%d Sucessfully read: %s\n" , i,thisfile) 
           // unmarshal byteArray using the JSON tags 
              jsonFile, err := ToJSON(yammy)
              json.Unmarshal(jsonFile, &c[i])
              jsonCatalogFile.Cc = append(jsonCatalogFile.Cc, c[i])
                fmt.Println("\nAny zero output is bad and indicates a YAML error.")
                fmt.Println("--------------------------------------------------")
                fmt.Println("              Course: "       + jsonCatalogFile.Cc[i].Id)
                fmt.Println("            Duration: " + strconv.Itoa(jsonCatalogFile.Cc[i].Duration.Hours))
                fmt.Printf("      Book Price Tags %d\n", len(jsonCatalogFile.Cc[i].Price.Book.PriceTags))
                fmt.Printf("    Public Price Tags %d\n", len(jsonCatalogFile.Cc[i].Price.Public.PriceTags))
                fmt.Printf("   Private Price Tags %d\n", len(jsonCatalogFile.Cc[i].Price.Private.PriceTags))
                fmt.Printf("Self Paced Price Tags %d\n", len(jsonCatalogFile.Cc[i].Price.Selfpaced.PriceTags))
                fmt.Printf("Extend LMS Price Tags %d\n", len(jsonCatalogFile.Cc[i].Price.ExtendLmsAccess.PriceTags))
                fmt.Printf("         Testimonials %d\n", len(jsonCatalogFile.Cc[i].Testimonials))
                fmt.Printf("                 Tags %d\n", len(jsonCatalogFile.Cc[i].Tags))
                fmt.Printf("             Chapters %d\n", len(jsonCatalogFile.Cc[i].Chapters))
                fmt.Printf("                 Labs %d\n", len(jsonCatalogFile.Cc[i].Labs))
              i++
              yammy = nil
              jsonFile = nil
          }
      }
      d.Close()
      return jsonCatalogFile
}



