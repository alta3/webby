package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strings"
)


//   _____ ______          _____   _____ _    _    _____ ____  _    _ _____   _____ ______  _____ 
//  / ____|  ____|   /\   |  __ \ / ____| |  | |  / ____/ __ \| |  | |  __ \ / ____|  ____|/ ____|
// | (___ | |__     /  \  | |__) | |    | |__| | | |   | |  | | |  | | |__) | (___ | |__  | (___  
//  \___ \|  __|   / /\ \ |  _  /| |    |  __  | | |   | |  | | |  | |  _  / \___ \|  __|  \___ \ 
//  ____) | |____ / ____ \| | \ \| |____| |  | | | |___| |__| | |__| | | \ \ ____) | |____ ____) |
// |_____/|______/_/    \_\_|  \_\\_____|_|  |_|  \_____\____/ \____/|_|  \_\_____/|______|_____/ 

// API Search Given a search string, returns ONLY course IDs for all matching courses.
func (cs Courses ) search() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //Retreive variable from GET URL
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    ss := r.URL.Path
    fmt.Printf("func searchstring searching for: %s\n",ss)
    // Create a new composite Course type. Interestingly, by adding existing subordinate
    // types to the cloned struct, items will OMIT them from the marshalling.
    // see: https://mycodesmells.com/post/working-with-embedded-structs
    // c := []PublicCourse{} 
    var id  []string
    var allid []string
    hits := 0
    totalhits := 0
    var js []byte
    var err error
    //Iterate over all courses, looking for the search string (ss)
    //If a match is found, add the course to id[i]
    for i, ThisCourse := range cs.Cc {
       hits = strings.Count( strings.ToLower(fmt.Sprintf("%v", cs.Cc[i])), ss )
       totalhits = totalhits + hits
       allid = append(id,ThisCourse.Id)
       if  hits > 0 {
           // Here is how you graft an existing type into a new "composite" type.
           id = append(id,ThisCourse.Id)
           fmt.Printf("%s Course has %d hits\n", ThisCourse.Id, hits )
       }
    }
    //If no courses match, SEND THEM ALL! 
    if totalhits == 0 {
       js, err = json.Marshal(allid)
    } else {
       js, err = json.Marshal(id)
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


