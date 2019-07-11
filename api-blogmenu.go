package main

import (
  "encoding/json"
  "fmt"
  "net/http"
)


//API - BlogMenu - Returns a canned Blogmenu pre-sorted by course type.
// The cannned blogmenu was generated at boot time, or server reload.

func (bmenu BlogMenus ) blogmenu() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    var js []byte
    var err error
    js, err = json.Marshal(bmenu)
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



