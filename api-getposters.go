package main

import (
  "encoding/json"
  "fmt"
  "net/http"
)


//API - posterMenu - Returns a list of posters..
// The canned poster structure was loaded at boot time, or server reload.

func (posters Posters ) getposters() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    var js []byte
    var err error
    js, err = json.Marshal(posters)
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



