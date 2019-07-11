package main

import (
  "time"
)



// ----------------COURSE STRUCT------------------------------- 

type Include struct {
  Item          string          `json:"item"`
  Description   string          `json:"description"`
}

type PriceTag struct {
  Id            string          `json:"id"`
  price         int             `json:"price"`
  Available     bool            `json:"available"`
  Description   string          `json:"description"`
  Includes      []Include       `json:"includes"`
}

type Book struct {
  PriceTags   []PriceTag           `json:"price-tags"`
}

type Selfpaced struct {
  PriceTags   []PriceTag           `json:"price-tags"`
}

type Public struct {
  PriceTags   []PriceTag           `json:"price-tags"`
}

type Private struct {
  PriceTags   []PriceTag           `json:"price-tags"`
}

type ExtendLmsAccess struct {
  PriceTags   []PriceTag           `json:"price-tags"`
}

type Price struct {
  Book             Book            `json:"book"`
  Selfpaced        Selfpaced       `json:"self-paced"`
  Public           Public          `json:"public"`
  Private          Private         `json:"private"`
  ExtendLmsAccess  ExtendLmsAccess `json:"extend-lms-access"`
}

type Slide struct {
//   GUID  string                  `rethinkdb:"guid" json:"guid"`
   Title string                  `rethinkdb:"title" json:"title"`
}

type Subchapter struct {
  Title  string                  `rethinkdb:"title" json:"title"`
  Slides []Slide                 `rethinkdb:"slides" json:"slides"`
}

type Chapter struct {
  Title       string             `rethinkdb:"title" json:"title"`
  SubChapters []Subchapter       `rethinkdb:"subchapters" json:"subchapters"` // TODO sync with codepen json
}

type Duration struct {
  Hours       int               `json:"hours"`
  Days        int               `json:"days"`
}

type Testimonial struct {
  Quote         string          `json:"quote"`
  Stars         int             `json:"stars"`
}

type Lab struct {
  Title string                  `rethinkdb:"title" json:"title"`
  File  string                  `rethinkdb:"file" json:"file"`
}


type Course struct {
  Id            string          `rethinkdb:"id" json:"id"`
  Filename      string          `rethinkdb:"filename" json:"filename"`
  WebURL        string          `rethinkdb:"weburl" json:"weburl"`
  CourseTitle   string          `rethinkdb:"name" json:"course-title"`
  HasSlides     bool            `rethinkdb:"has-slides" json:"has-slides"`
  HasLabs       bool            `rethinkdb:"has-labs" json:"has-labs"`
  HasVideos     bool            `rethinkdb:"has-videos" json:"has-videos"`
  Private       bool            `rethinkdb:"private" json:"private"`
  Chapters      []Chapter       `rethinkdb:"chapters" json:"chapters"` // TODO update to single-slide-mode
  Labs          []Lab           `rethinkdb:"labs" json:"labs"`         // TODO Write
  Expires       time.Time       `rethinkdb:"-" json:"-"`
  Purchased     bool            `rethinkdb:"-" json:"-"`
  Price         Price           `rethinkdb:"-" json:"price"`
  Duration      Duration        `json:"duration"`
  Testimonials  []Testimonial   `json:"testimonials"`
  VideoLink     string          `json:"video-link"`
  Overview      string          `json:"overview"`
  Tags          []string        `json:"tags"`
  Courseicon    string          `json:"courseicon"`      // TODO courseicons will be under images/courseicons
  Stars         int             `json:"stars"`
  Audience      string          `json:"audience"`
  Prereqs       []string        `json:"prereqs"`
  Postreqs      []string        `json:"postreqs"`
}

type Courses struct {
  Cc           []Course         `json:"courses"`
}


type PublicCourse struct {
  Course
  Chapters      []Chapter       `json:"chapters,omitempty"` // TODO update to single-slide-mode
  Labs          []Lab           `json:"labs,omitempty"`    // TODO Write
}


