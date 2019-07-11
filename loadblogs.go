package main

import (
  "encoding/json"
  "fmt"
  "log"
  "strings"
  "os"
  "path"
  "io/ioutil"
  "path/filepath"
  "github.com/gomarkdown/markdown"
  "github.com/gomarkdown/markdown/parser"

)




// ----------------------LOAD BLOGs------------------------------
// Load all blogs when server boots
func Loadblogs() Blogs {
    // Create a OS compliant path: microsoft "\" or linux "/"
    dirname := path.Join("deploy", "blog")
    d, err := os.Open(dirname)
    if err != nil {
        log.Printf("No blogs directory! %s" , dirname)
        os.Exit(1)
    }
    // If n > 0, Readdirnames(n) returns at most n names
    // If n < 0, Readdirnames(n) returns ALL names
    n := -1
    // reads < n > files in directory < d >
    filenames, err := d.Readdirnames(n)
    if err != nil {
        log.Printf("No files in blogs directory! %s\n" , dirname)
        os.Exit(1)
    }
    var   b          Blog
    var   allblogs   Blogs
    fmt.Println("--------------------------------------------------")
    fmt.Printf(" Reading BLOG files from directory: %s\n", dirname)
    fmt.Println(" Any zero output is bad and indicates a YAML error.")
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
            fmt.Println("--------------------------------------------------")
            fmt.Printf("#%d  Sucessfully read: %s\n" , i+1,thisfile)
            parts := string(yammy)
            //split YAML header from markdown body using "\n---" delimiter
            z := strings.Split(parts, "\n---")
            //check if there exactly two parts or skip to next file
            if len(z) != 2 {
                fmt.Printf("    ********* FILE PARSE FAIL **************\n    BROKEN FILE: %s, skipping\n",thisfile)
                fmt.Printf("    SPLIT-COUNT: %d should be 2\n",len(z))
                fmt.Printf("    Should be easy to fix, check file format\n\n\n")
            }
            if len(z) == 2 {
                fmt.Printf("FIRST SPLIT is the YAML HEADER:\n%s\n2nd SPLIT MARKDOWN: %d characters\n", z[0],len(z[1]))
                extensions := parser.CommonExtensions | parser.AutoHeadingIDs
                parser := parser.NewWithExtensions(extensions)
                md := []byte(z[1])
                //load html into b.Content
                myhtml := markdown.ToHTML(md, parser, nil)
								if err != nil {
													panic(err)
								}
                b.HtmlContent = string(myhtml)
                //der().Set("Access-Control-Allow-Headers", fmt.Printf("Content:\n--------\n %s\n", b.Content)
                // unmarshal byteArray using the JSON tags 
                jsonFile, err := ToJSON(yammy)
                if err != nil {
                   log.Printf("jsonFile error: %s\n", err)
                }
                json.Unmarshal(jsonFile, &b)
                allblogs = append(allblogs, b)
                fmt.Printf("                  ID: %s\n", allblogs[i].Id)
                fmt.Printf("               Title: %s\n", allblogs[i].Title)
                fmt.Printf("                Date: %s\n", allblogs[i].Date)
                fmt.Printf("              Weight: %s\n", allblogs[i].Weight)
                fmt.Printf("              Author: %s\n", allblogs[i].Author)
                fmt.Printf("    Content in bytes: %d\n", len(allblogs[i].HtmlContent))
                jsonFile = nil
            }
            yammy = nil
            i++
       }
    }
    d.Close()
    return allblogs
}


