package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strings"
)


// blogsearch returns the blog ID of all blogs that match the search string.
func (b Blogs) blogsearch() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //Retreive search string value from GET URL Path
		//Asterisk will return all blogs
   (w).Header().Set("Access-Control-Allow-Methods", "*") 
   (w).Header().Set("Access-Control-Allow-Origin", "*") 
   (w).Header().Set("Access-Control-Allow-Headers","*")
		ss := r.URL.Path
		if ss == "" {
			ss = "*"
		}
    fmt.Printf("Searching blogs for: '%s'\n",ss)
    var id  []string
    hits := 0
    totalhits := 0
    var js []byte
    var err error
    //Iterate over all courses, looking for the search string (ss)
    //If a match is found, add the blog Id to id[i]
    for i, Thisblog := range b {
			 if ss == "*" {
					 id = append(id,Thisblog.Id)
					 fmt.Printf("Appending %s Blog\n", Thisblog.Id)
	     }  else {
           hits = strings.Count( strings.ToLower(fmt.Sprintf("%v", b[i])), ss )
           totalhits = totalhits + hits
			 }
       if  hits > 0 {
           id = append(id,Thisblog.Id)
           fmt.Printf("%s Blog has %d hits\n", Thisblog.Id, hits )
       }
    }
    //If no blog match, SEND Null 
    if totalhits == 0 {
       js, err = json.Marshal(id)
			 fmt.Printf("No hits - sending NULL\n")
    } else {
       js, err = json.Marshal(id)
    }
    if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
       fmt.Printf("Error %s:\n", err)
       return
    }
    w.Write(js)
    return
    })
}

