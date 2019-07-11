package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "time"
)


//API GetEvents - Returns all events 
func (e Events) getevents() http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       now := time.Now()
       fmt.Printf("Time right now: %s\n", now.Format("2006-01-08"))
       var js []byte
       var err error
       var ce  Events  //ce means CurrentEvents
   // Iterate over all events, skipping past events.
       layout := "2006-01-02"
       for _, ThisEvent := range e.Events {
          fmt.Printf(" Event: %s, %s\n",ThisEvent.StartDate, ThisEvent.Title)
          t, _ := time.Parse(layout, ThisEvent.StartDate)
          if t.After(now) {
              ce.Events = append(ce.Events, ThisEvent)
          } else {
            fmt.Println("----------------OLD EVENT---------------")
            fmt.Printf("| OLD!!!: %s %s\n", ThisEvent.Title, ThisEvent.StartDate)
            fmt.Println("----------------OLD EVENT---------------")
            }
       }
       js, err = json.Marshal(ce)
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


