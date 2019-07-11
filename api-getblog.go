package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strings"
)


// Given A valid Blog ID, returns the blog
func (b Blogs) getblog() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //Retreive variable from GET URL
    fmt.Println("--------------------------------------------------")
		if r.URL.Path == "" {
			r.URL.Path = "noblog"
			Template().ServeHTTP(w, r)
			return
		}
	  ss := r.URL.Path
    fmt.Printf("%s blog requested\n",ss)
    var blog  Blog
    var js []byte
    var err error
    //Iterate over all courses, looking for the search string (ss)
    //If a match is found, add the blog Id to id[i]
    for _, Thisblog := range b {
			 if strings.ToLower(ss) == strings.ToLower(Thisblog.Id) {
					 blog = Thisblog
					 fmt.Printf("Found blog %s\n", Thisblog.Id)
           break
	     }
    }
    //If no blog match, SEND Null 
    if blog.Id  == "" {
       js, err = json.Marshal("")
			 fmt.Printf("Blog %s not found\n", ss)
    } else {
       js, err = json.Marshal(blog)
    }
    if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
       fmt.Printf("Error %s:\n", err)
       return
    }
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
    return
    })
}


