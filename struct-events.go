package main

import ()



// ---------------------------------------------------------

type Events struct {
  Events         []Event
}

type Event struct {
  Id             string         `yaml:"id"`
  Category       string         `yaml:"category"`
  Title          string         `yaml:"title"`
  StartDate      string         `yaml:"startdate"`
  EndDate        string         `yaml:"enddate"`
  CourseId       string         `yaml:"courseid"`
  Image          string         `yaml:"image"`
  Location       string         `yaml:"location"`
	Price          int            `yaml:"price"`
}


type EventsMenu struct {
    EventsMenu struct {
        PastEvents           []Event   `json:"past-events"`
        Webinars             []Event   `json:"webinars"`
        UpcomingClasses      []Event   `json:"upcoming-classes"`
    } `json:"eventsmenu"`
}

