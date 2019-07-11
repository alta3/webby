package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "time"
  "strings"
)


//API geteventsmenu - Returns all events 
func (e Events) geteventsmenu() http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       now := time.Now()
       fmt.Printf("Time right now: %s\n", now.Format("2006-01-08"))
       var js  []byte
       var err error
       var em  EventsMenu
       var ce  Events  //ce means CurrentEvents
   // Iterate over all events, skipping past events.
       layout := "2006-01-02"
       for _, ThisEvent := range e.Events {
          fmt.Printf(" Event: %s, %s\n",ThisEvent.StartDate, ThisEvent.Title)
          t, _ := time.Parse(layout, ThisEvent.StartDate)
          if t.After(now) {
              ce.Events = append(ce.Events, ThisEvent)
                // Iterate over every event and place if correnct sub-struct
              if strings.Contains(strings.ToUpper(ThisEvent.Category),"WEBINAR") {
                   em.EventsMenu.Webinars = append(em.EventsMenu.Webinars, ThisEvent)
                   fmt.Printf("  - %s = WENINAR\n", ThisEvent.Title)
              } else if strings.Contains(strings.ToUpper(ThisEvent.Category),"CLASS") {
                   em.EventsMenu.UpcomingClasses = append(em.EventsMenu.UpcomingClasses, ThisEvent)
                   fmt.Printf("  - %s = CLASS\n", ThisEvent.Title)
              }
          } else {
                fmt.Println("----------------OLD EVENT---------------")
                fmt.Printf("| OLD!!!: %s %s\n", ThisEvent.Title, ThisEvent.StartDate)
                fmt.Println("----------------OLD EVENT---------------")
                em.EventsMenu.PastEvents = append(em.EventsMenu.PastEvents, ThisEvent)
                fmt.Printf("  - %s\n", ThisEvent.Title)
            }
       }
       js, err = json.Marshal(em)
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

