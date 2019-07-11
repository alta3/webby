package main

import (
  "encoding/json"
  "fmt"
  "net/http"
)


//   _____ ____  _    _ _____   _____ ______    _____ _    _ __  __ __  __          _______     __
//  / ____/ __ \| |  | |  __ \ / ____|  ____|  / ____| |  | |  \/  |  \/  |   /\   |  __ \ \   / /
// | |   | |  | | |  | | |__) | (___ | |__    | (___ | |  | | \  / | \  / |  /  \  | |__) \ \_/ / 
// | |   | |  | | |  | |  _  / \___ \|  __|    \___ \| |  | | |\/| | |\/| | / /\ \ |  _  / \   /  
// | |___| |__| | |__| | | \ \ ____) | |____   ____) | |__| | |  | | |  | |/ ____ \| | \ \  | |   
//  \_____\____/ \____/|_|  \_\_____/|______| |_____/ \____/|_|  |_|_|  |_/_/    \_\_|  \_\ |_|


//API - Course Summary - Given a valid courseID, returns that course's: Id, Tag, Stars, and Name. 
//Why do we need this? A search will provide a list of IDs. This API is used to get summary data PER ID.
func (cs Courses ) getsummary() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    type MenuItems struct {
      Id           string          `json:"id"`
      CourseTitle  string          `json:"course-title"`
      Tags         []string        `json:"tags"`
			Stars        int             `json:"stars"`
    }
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    //Retreive variable from GET URL
    fmt.Println("--------------------------------------------------")
		if r.URL.Path == "" {
			r.URL.Path = "nocourse"
			Template().ServeHTTP(w, r)
			return
		}
	  ss := r.URL.Path
    fmt.Printf("Menu for:   %s\n",ss)
    mi := MenuItems{}
    mis := []MenuItems{}
		var js []byte
    var err error
    //Iterate over all courses, copy Id, Name, and any Tags
    for _, ThisCourse := range cs.Cc {
       fmt.Printf("Menu Item = %s, %s, %s\n", ThisCourse.Id, ThisCourse.CourseTitle, ThisCourse.Tags)
       if  ThisCourse.Id == ss {
           mi.Id=ThisCourse.Id
           mi.CourseTitle=ThisCourse.CourseTitle
           mi.Tags=ThisCourse.Tags
					 mi.Stars=47
           mis = append(mis,mi)
           break
       }
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


