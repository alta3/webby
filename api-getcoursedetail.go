package main

import (
  "encoding/json"
  "fmt"
  "net/http"
)


//   ____ ___  _   _ ____  ____  _____   ____  _____ _____  _    ___ _     
//  / ___/ _ \| | | |  _ \/ ___|| ____| |  _ \| ____|_   _|/ \  |_ _| |    
// | |  | | | | | | | |_) \___ \|  _|   | | | |  _|   | | / _ \  | || |    
// | |__| |_| | |_| |  _ < ___) | |___  | |_| | |___  | |/ ___ \ | || |___ 
//  \____\___/ \___/|_| \_\____/|_____| |____/|_____| |_/_/   \_\___|_____|

// API Get Course Detail, given a valid courseID, return all course details with GUIDs blanked out.
func (cs Courses ) getcoursedetail() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    //Retreive variable from GET URL
    ss := r.URL.Path
    fmt.Printf("func searchstring searching for: %s\n",ss)
    var c Course
    hits := 0
    totalhits := 0
    var js []byte
    var err error
    //Iterate over all courses, looking for the search string (ss)
    //If a match is found, add the course to c.Cc[i]
    for _, ThisCourse := range cs.Cc {
       totalhits = totalhits + hits
       if  ThisCourse.Id == ss {
           c = ThisCourse
           hits = 1
           totalhits = 1
           fmt.Printf("Found detail for %s Course\n", ThisCourse.Id, hits )
           break
       }
    }
    //If no courses match, SEND THEM ALL! 
    if totalhits == 0 {
       js, err = json.Marshal(hits)
    } else {
       js, err = json.Marshal(c)
    }
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

