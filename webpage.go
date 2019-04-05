package main

import (
//  "regexp"


//	"github.com/tdewolff/minify"
//	"github.com/tdewolff/minify/html"

"encoding/json"
	"fmt"
	"html/template"
        "path/filepath"
	"log"
        "unicode"
        "bytes"
        "errors"
	"net/http"
        "strings"
	"os"
	"path"
        "io/ioutil"
        "strconv"
        "github.com/ghodss/yaml"
        "time"
  "github.com/gomarkdown/markdown"
  "github.com/gomarkdown/markdown/parser"
//  "gopkg.in/russross/blackfriday.v2"

//        "gopkg.in/yaml.v2"
//        "github.com/gorilla/mux"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)



type CheckoutToken struct {
	Token      stripe.Token `json:"token"`
	CourseId   string       `json:"course_id"`
	CourseName string       `json:"course_name"`
}

func Checkout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ct           CheckoutToken
			chargeParams *stripe.ChargeParams
			ch           *stripe.Charge
			price        int64
			ok           bool
		)

		coursePrices := map[string]int64{
			"openstack-selfpaced":     14995,
			"avaya-selfpaced":         14995,
			"volte-selfpaced":         44500,
		}

		d := json.NewDecoder(r.Body)
		err := d.Decode(&ct)
		if err != nil {
			err = fmt.Errorf("failed to decode checkoutToken: %s", err)
			goto respond
		}
		log.Printf("%+v", ct)

		price, ok = coursePrices[ct.CourseId]
		if !ok {
			err = fmt.Errorf("failed to lookup price: %s", err)
			goto respond
		}

		stripe.Key = "xxxxxxxxxxxxxxxxxxxxxxxx"
		//stripe.Key = "xxxxxxxxxxxxxxxxxxxxxxxx"
		chargeParams = &stripe.ChargeParams{
			Amount:          stripe.Int64(price),
			Currency:        stripe.String(string(stripe.CurrencyUSD)),
			Description:     stripe.String("Alta3 Research - " + ct.CourseName + " for " + ct.Token.Email),
			ReceiptEmail:    stripe.String(ct.Token.Email),
		}
		err = chargeParams.SetSource(ct.Token.ID)
		if err != nil {
			err = fmt.Errorf("failed to charge, bad token: %s", err)
			goto respond
		}
		ch, err = charge.New(chargeParams)

		if err != nil {
			err = fmt.Errorf("failed to charge, bad params: %s", err)
			goto respond
		}
		log.Printf("%+v", ch)

	respond:
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)

		if err != nil {
			log.Printf("failed, %v", err)
			enc.Encode(map[string]string{
				"error":   "checkout failed",
				"success": "true",
			})
			return
		}
		enc.Encode(map[string]string{"success": "true"})
		return
	})
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func errorHandler(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if status == http.StatusNotFound {
			r.URL.Path = "404.html"
			Template().ServeHTTP(w, r)
		}
	})
}

func Template() http.Handler {



	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
                        r.URL.Path = "index.html" 
		} else if r.URL.Path == "deploy/robots.txt" {
			http.ServeFile(w, r, "robots.txt")
			return
		} else if r.URL.Path == "deploy/sitemap.xml" {
			http.ServeFile(w, r, "sitemap.xml")
			return
		}

                // create the full pathnames of the files to be merged
		lp := path.Join("deploy/templates", "layout.html")
		fp := path.Join("deploy/templates", r.URL.Path)

                // check if file fp points to actually exits, 404 if file not there 
                // os.Stat returns TWO values, Info and an error code, hence <info, err := > below
		info, err := os.Stat(fp)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("404 - This file does NOT exist: %s", r.URL.Path)
				errorHandler(http.StatusNotFound).ServeHTTP(w, r)
				return
			}
		}
                //if fp is a directory, send 404
		if info.IsDir() {
                        fmt.Println("404 - The file request is pointing to a directory! ")
			errorHandler(http.StatusNotFound).ServeHTTP(w, r)
			return
		}

                // ParseFiles gathers all the defines ad ignores the rest  reads is all the {{ jinja-like stuff }}
                // template.ParseFiles returns TWO values, the combined jinja stuff, and error code
		tmpl, err := template.ParseFiles(lp, fp)
		if err != nil {
			log.Printf("Template Error: %s", err)
		}
                // merge tmpl with {{ define layout }} TO  {{ end }} and write response 
		err = tmpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			log.Printf("Template Error: %s", err)
		}
		return
	})
}

func BlogTemplate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			r.URL.Path = "blog"
			Template().ServeHTTP(w, r)
			return
		}
		lp := path.Join("templates", "layout.html")
		cp := path.Join("blog", "blog_layout.html")
		fp := path.Join("blog", r.URL.Path+".html")
		info, err := os.Stat(fp)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("404: %s", r.URL.Path)
				errorHandler(http.StatusNotFound).ServeHTTP(w, r)
				return
			}
		}
		if info.IsDir() {
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles(lp, cp, fp)
		if err != nil {
			log.Printf("Template Error: %s", err)
		}
		err = tmpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			log.Printf("Template Error: %s", err)
		}
		return
	})
}


// type Include struct {
//   Item          string          `json:"item"`
//   Description   string          `json:"description"`
//}

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


type Event struct {
  Id             string         `yaml:"id"`
  Title          string         `yaml:"title"`
  StartDate      string         `yaml:"startdate"`
  EndDate        string         `yaml:"enddate"`
  CourseId       string         `yaml:"courseid"`
  Image          string         `yaml:"image"`
  Location       string         `yaml:"location"`
}


type Events struct {
  Events         []Event
}


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



//Load EVENTS
//------------------------------------------------------------
func LoadEvents() Events {
      // Create a OS compliant path: microsoft "\" or linux "/"
      dirname := path.Join("deploy", "event")
      d, err := os.Open(dirname)
      if err != nil {
          log.Printf("No events directory! %s, %s" , d, err)
          os.Exit(1)
      }
      var    ev   Events
      fmt.Println("---------------LOADING EVENTS---------------------")
      fmt.Printf(" Reading events files in directory: %s\n", dirname)
      thisfile := path.Join(dirname, "events.yaml")
      _ , err = os.Stat(thisfile)
      if err != nil {
          if os.IsNotExist(err) {
              log.Printf(" file is missing!: %s\n ", thisfile)
          }
      } 
      yammy, err := ioutil.ReadFile(thisfile)
      if err != nil {
          log.Printf("yammy.Get err: %s\n", err)
          }
     // unmarshal byteArray using the JSON tags 
	    err = yaml.Unmarshal(yammy, &ev)
      if err != nil {
				 log.Printf("Unmarshal: %v", err)
				  }

			fmt.Printf(" Successfully read: %s\n", thisfile) 
      fmt.Printf(" Events: %+v\n", ev )
      fmt.Println("--------------------------------------------------")
      return ev 
}

//Load COURSES
//------------------------------------------------------------
func Load() Courses {
      // Create a OS compliant path: microsoft "\" or linux "/"
      dirname := path.Join("deploy", "courses")
      d, err := os.Open(dirname)
      if err != nil {
          log.Printf("No courses directory! %s" , err)
          os.Exit(1)
      }
      // If n > 0, Readdirnames(n) returns at most n names
      // If n < 0, Readdirnames(n) returns ALL names
      n := -1 
      // reads < n > files in directory < d >
      filenames, err := d.Readdirnames(n)
      if err != nil {
          log.Printf("No files in course directory! %s\n" , err)
          os.Exit(1)
      }
      c := make([]Course,50) 
      var jsonCatalogFile Courses
      fmt.Println("--------------------------------------------------")
      fmt.Printf(" Reading files in this directory: %s\n", dirname)
      i := 0
      for _, filename := range filenames {
          thisfile := path.Join(dirname, filename)
          _ , err := os.Stat(thisfile)
          if err != nil {
              if os.IsNotExist(err) {
                  log.Printf("file is missing!: %s\n ", filename)
              }
          } 
          if filepath.Ext(thisfile) == ".yaml" {
              yammy, err := ioutil.ReadFile(thisfile)
              if err != nil {
                 log.Printf("yammy.Get err: %s\n", err)
              }
              fmt.Printf("%d Sucessfully read: %s\n" , i,thisfile) 
           // unmarshal byteArray using the JSON tags 
              jsonFile, err := ToJSON(yammy)
              json.Unmarshal(jsonFile, &c[i])
              jsonCatalogFile.Cc = append(jsonCatalogFile.Cc, c[i])
                fmt.Println("\nAny zero output is bad and indicates a YAML error.")        
                fmt.Println("--------------------------------------------------")
                fmt.Println("              Course: "       + jsonCatalogFile.Cc[i].Id)
                fmt.Println("            Duration: " + strconv.Itoa(jsonCatalogFile.Cc[i].Duration.Hours))
                fmt.Printf("      Book Price Tags %d\n", len(jsonCatalogFile.Cc[i].Price.Book.PriceTags))
                fmt.Printf("    Public Price Tags %d\n", len(jsonCatalogFile.Cc[i].Price.Public.PriceTags))
                fmt.Printf("   Private Price Tags %d\n", len(jsonCatalogFile.Cc[i].Price.Private.PriceTags))
                fmt.Printf("Self Paced Price Tags %d\n", len(jsonCatalogFile.Cc[i].Price.Selfpaced.PriceTags))
                fmt.Printf("Extend LMS Price Tags %d\n", len(jsonCatalogFile.Cc[i].Price.ExtendLmsAccess.PriceTags))
                fmt.Printf("         Testimonials %d\n", len(jsonCatalogFile.Cc[i].Testimonials))
                fmt.Printf("                 Tags %d\n", len(jsonCatalogFile.Cc[i].Tags))
                fmt.Printf("             Chapters %d\n", len(jsonCatalogFile.Cc[i].Chapters))
                fmt.Printf("                 Labs %d\n", len(jsonCatalogFile.Cc[i].Labs))
              i++
              yammy = nil
              jsonFile = nil
          }
      }
      d.Close()
      return jsonCatalogFile 
}


func (cs Courses)  Select(id string) (Courses, error) {
     log.Printf("WORKING: Looking for %s\n", id)
     var c Courses
		 for _, ThisCourse := range cs.Cc  {
          if ThisCourse.Id == id  {
              c.Cc = append(c.Cc, ThisCourse)
              fmt.Printf("FOUND %d Record, returning: %s\n" , len(c.Cc), c.Cc[0].Id)
              return c, nil
          }
      }
     return c, errors.New(fmt.Sprintf("Course ID \"%s\" does NOT exist\n", id ))
}


func (cs Courses)  Search(ls string) (Courses, error) {
     ls = strings.ToLower(ls)
     fmt.Println("--------------------------------------------------")
     log.Printf("SEARCH FUNC REPORTING: Searching for %s\n", ls)
     var c Courses
     i := 0
     hits := 0
     totalhits := 0
     for _, ThisCourse := range cs.Cc {
        hits = strings.Count( strings.ToLower(fmt.Sprintf("%v", cs.Cc[i])), ls )
        totalhits = totalhits + hits
        if  hits > 0 {
            c.Cc = append(c.Cc, ThisCourse)
            fmt.Printf("%s Course has %d hits\n", ThisCourse.Id, hits )
        }
        i++
      }
     if  totalhits == 0 {
       return c, errors.New(fmt.Sprintf("No course contains any information regarding \"%s\"" , ls ))
     }
    return c, nil
}




//----------------------------------------------------------------
//Allow painless Ingesting of YAML
//----------------------------------------------------------------
func ToJSON(data []byte) ([]byte, error) {
    if hasJSONPrefix(data) {
        return data, nil
    }
    return yaml.YAMLToJSON(data)
}

var jsonPrefix = []byte("{")

// hasJSONPrefix returns true if the provided buffer starts with "{".
func hasJSONPrefix(buf []byte) bool {
    return hasPrefix(buf, jsonPrefix)
}

// Return true if the first non-whitespace bytes in buf is prefix.
func hasPrefix(buf []byte, prefix []byte) bool {
    trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
    return bytes.HasPrefix(trim, prefix)
}



//----------------------------------------------------------------

func  (cs Courses) CourseTemplate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			r.URL.Path = "courses"
			Template().ServeHTTP(w, r)
			return
		}
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
                var err error
		lp := path.Join("deploy/templates", "layout.html")
		fp := path.Join("deploy/courses", r.URL.Path)
		info, err := os.Stat(fp)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("404: %s", r.URL.Path)
				errorHandler(http.StatusNotFound).ServeHTTP(w, r)
				return
			}
		}
		if info.IsDir() {
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles(lp,  fp)
		if err != nil {
			log.Printf("Template Error: %s", err)
		}
		err = tmpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			log.Printf("Template Error: %s", err)
		}
		return
	})
}



//--------------------------------------------------------------
func (cs Courses ) getsummarylist() http.Handler {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    type PopupItems struct {
      Id                 string          `json:"id"`
      CourseTitle        string          `json:"course-title"`
      Stars              int             `json:"stars"`
      Duration           int             `json:"duration"`
      Overview           string          `json:"overview"`
      SelfpacedPrice     int             `json:"selfpacedprice"`
      PublicPrice        int             `json:"publicprice"`
      Courseicon         string          `json:"courseicon"`
    }
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    popi := PopupItems{}
    popis := []PopupItems{}
    var js []byte
    var err error
    //Iterate over all courses, Copy Id, Name, Stars, Duration, Overview, Price, and Courseicon
    for _, ThisCourse := range cs.Cc {
       fmt.Println("--------------------------------------------------")
       fmt.Printf("_Course PopUp_  = %s, %s, %s, %s, %s, %s, %s\n", ThisCourse.Id, ThisCourse.CourseTitle, ThisCourse.Testimonials[0].Stars, ThisCourse.Duration, ThisCourse.Overview, ThisCourse.Price, ThisCourse.Courseicon)
       popi.Id=ThisCourse.Id
       popi.CourseTitle=ThisCourse.CourseTitle
       popi.Stars=ThisCourse.Testimonials[0].Stars
       popi.Duration=ThisCourse.Duration.Hours
//       popi.Overview=ThisCourse.Overview
       popi.SelfpacedPrice=ThisCourse.Price.Selfpaced.PriceTags[0].price
       popi.SelfpacedPrice=ThisCourse.Price.Public.PriceTags[0].price
       popi.Courseicon=ThisCourse.Courseicon
       popis = append(popis,popi)
    }
    //If no courses match, SEND THEM ALL! 
       js, err = json.Marshal(popis)
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

//  __  __ ______ _____          __  __ ______ _   _ _    _ 
// |  \/  |  ____/ ____|   /\   |  \/  |  ____| \ | | |  | |
// | \  / | |__ | |  __   /  \  | \  / | |__  |  \| | |  | |
// | |\/| |  __|| | |_ | / /\ \ | |\/| |  __| | . ` | |  | |
// | |  | | |___| |__| |/ ____ \| |  | | |____| |\  | |__| |
// |_|  |_|______\_____/_/    \_\_|  |_|______|_| \_|\____/

//API - MegaMenu - Returns all courses: Id, Tag, Stars, and Name. 
func (cs Courses ) getmegamenu() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    type MenuItems struct {
      Id           string          `json:"id"`
      CourseTitle  string          `json:"course-title"`
      Tags         []string        `json:"tags"`
    }
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    mi := MenuItems{} 
    mis := []MenuItems{} 
    var js []byte
    var err error
    //Iterate over all courses, copy Id, Name, and any Tags
    for _, ThisCourse := range cs.Cc {
       fmt.Println("--------------------------------------------------")
       fmt.Printf("Menu Item = %s, %s, %s\n", ThisCourse.Id, ThisCourse.CourseTitle, ThisCourse.Tags)
       mi.Id=ThisCourse.Id
       mi.CourseTitle=ThisCourse.CourseTitle
       mi.Tags=ThisCourse.Tags
       mis = append(mis,mi)
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


//  ____  _      ____   _____  _____ 
// |  _ \| |    / __ \ / ____|/ ____|
// | |_) | |   | |  | | |  __| (___  
// |  _ <| |   | |  | | | |_ |\___ \ 
// | |_) | |___| |__| | |__| |____) |
// |____/|______\____/ \_____|_____/ 

//------------------BLOGS------------------------
type Blog struct {
  Id             string         `yaml:"id"`
  Author         string         `yaml:"author"`
	Category       string         `yaml:"category"`
  Date           string         `yaml:"date"`
  Title          string         `yaml:"title"`
  Weight         string         `yaml:"weight"`
	Intro          string         `yaml:"intro"`
	VideoLink      string         `yaml:"video-link"`
	HtmlContent    string         `yaml:"html-content"`
}

type Blogs []Blog

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


// Given a valid Blog ID, returns the blog
func (b Blogs) getblog() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //Retreive variable from GET URL
    fmt.Println("--------------------------------------------------")
		if r.URL.Path == "" {
			r.URL.Path = "noblog"
			Template().ServeHTTP(w, r)
			return
		}
	  ss := r.URL.Path
    fmt.Printf("%s blog requested\n",ss)
    var blog  Blog
    var js []byte
    var err error
    //Iterate over all courses, looking for the search string (ss)
    //If a match is found, add the blog Id to id[i]
    for _, Thisblog := range b {
			 if strings.ToLower(ss) == strings.ToLower(Thisblog.Id) {
					 blog = Thisblog
					 fmt.Printf("Found blog %s\n", Thisblog.Id)
           break
	     } 
    }
    //If no blog match, SEND Null 
    if blog.Id  == "" {
       js, err = json.Marshal("")
			 fmt.Printf("Blog %s not found\n", ss)
    } else {
       js, err = json.Marshal(blog)
    }
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


// ----------------------LOAD BLOGs------------------------------
// Load all blogs when server boots
func Loadblogs() Blogs {
    // Create a OS compliant path: microsoft "\" or linux "/"
    dirname := path.Join("deploy", "blog")
    d, err := os.Open(dirname)
    if err != nil {
        log.Printf("No blogs directory! %s" , dirname)
        os.Exit(1)
    }
    // If n > 0, Readdirnames(n) returns at most n names
    // If n < 0, Readdirnames(n) returns ALL names
    n := -1
    // reads < n > files in directory < d >
    filenames, err := d.Readdirnames(n)
    if err != nil {
        log.Printf("No files in blogs directory! %s\n" , dirname)
        os.Exit(1)
    }
    var   b          Blog
    var   allblogs   Blogs
    fmt.Println("--------------------------------------------------")
    fmt.Printf(" Reading BLOG files from directory: %s\n", dirname)
    fmt.Println(" Any zero output is bad and indicates a YAML error.")
    i := 0
    for _, filename := range filenames {
        thisfile := path.Join(dirname, filename)
        _ , err := os.Stat(thisfile)
        if err != nil {
            if os.IsNotExist(err) {
                log.Printf("file is missing!: %s\n ", filename)
            }
        }
        if filepath.Ext(thisfile) == ".yaml" {
            yammy, err := ioutil.ReadFile(thisfile)
            if err != nil {
               log.Printf("yammy.Get err: %s\n", err)
            }
            fmt.Println("--------------------------------------------------")
            fmt.Printf("#%d  Sucessfully read: %s\n" , i+1,thisfile)
            parts := string(yammy)
            //split YAML header from markdown body using "\n---" delimiter
            z := strings.Split(parts, "\n---")
            //check if there exactly two parts or skip to next file
            if len(z) != 2 {
                fmt.Printf("    ********* FILE PARSE FAIL **************\n    BROKEN FILE: %s, skipping\n",thisfile)
                fmt.Printf("    SPLIT-COUNT: %d should be 2\n",len(z))
                fmt.Printf("    Should be easy to fix, check file format\n\n\n")
            }
            if len(z) == 2 {
                fmt.Printf("FIRST SPLIT is the YAML HEADER:\n%s\n2nd SPLIT MARKDOWN: %d characters\n", z[0],len(z[1]))
                extensions := parser.CommonExtensions | parser.AutoHeadingIDs
                parser := parser.NewWithExtensions(extensions)
                md := []byte(z[1])
                //load html into b.Content
                myhtml := markdown.ToHTML(md, parser, nil)
								if err != nil {
													panic(err)
								}
                b.HtmlContent = string(myhtml)
                //der().Set("Access-Control-Allow-Headers", fmt.Printf("Content:\n--------\n %s\n", b.Content)
                // unmarshal byteArray using the JSON tags 
                jsonFile, err := ToJSON(yammy)
                if err != nil {
                   log.Printf("jsonFile error: %s\n", err)
                }
                json.Unmarshal(jsonFile, &b)
                allblogs = append(allblogs, b)
                fmt.Printf("                  ID: %s\n", allblogs[i].Id)
                fmt.Printf("               Title: %s\n", allblogs[i].Title)
                fmt.Printf("                Date: %s\n", allblogs[i].Date)
                fmt.Printf("              Weight: %s\n", allblogs[i].Weight)
                fmt.Printf("              Author: %s\n", allblogs[i].Author)
                fmt.Printf("    Content in bytes: %d\n", len(allblogs[i].HtmlContent))
                jsonFile = nil
            }
            yammy = nil
            i++
       }
    }
    d.Close()
    return allblogs
}


// Blog Menu
// Returns a blog menu that does all the work for the front end developers

type BlogsByCategory struct {
	BlogCategory     string        `json:"blog-category"`
	Blogs            []Blog        `json:"blogs"`
}

type BlogMenus     []BlogsByCategory

type BlogCategory  []string

func (b Blogs) blogmenumaker() BlogMenus  {
	  var existing         bool
	  var blogmenus        []BlogsByCategory
		var blogsbycategory  BlogsByCategory
    var blogs            Blogs
	  var categories       []string
		// Iterate over all blogs, and derive a list of unique categories
		for _, thisblog := range  b {
        //Iterate over array of categories
				existing = false
				for _,  thiscategory := range categories {
								if thiscategory == thisblog.Category {
							  existing = true
					      }
			  }
        if existing == false {
          categories = append (categories, thisblog.Category)
		    }
    }
    //At this stage, a list of unique categories has been gathered,
		//so build the blogmenu
    // Interate over each category 
		for _, thiscategory := range categories {
				fmt.Printf("\"%s\"\n",thiscategory)
				//Iternate over every blog for that category
        for _, thisblog := range b {
                if thisblog.Category == thiscategory {
							      blogs = append( blogs, thisblog)
										fmt.Printf("  - %s\n", thisblog.Title)
                }
         }
				 blogsbycategory.BlogCategory = thiscategory
				 blogsbycategory.Blogs = blogs
				 blogs = nil
				 blogmenus = append(blogmenus,blogsbycategory)
    }
		return blogmenus
}


//API - New MegaMenu - Returns a canned megamenu pre-sorted by course type. 
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





// Course Menu
// Returns a course menu that will be easy for the front end js to implement

type CourseMenu  []MiniMenu    // `json:"coursemenu"`

type MiniMenu struct {
  MiniMenuTitle  string          `json:"mini-menu"`
  MmItems        []MmItem        `json:"mm-items"`
}

type MmItem struct {
  Id              string         `json:"id"`
  CourseTitle     string         `json:"course-title"`
}

func (cc Courses)  menumaker() CourseMenu  {
    var minimenutitles  []string
    var mmitem          MmItem
    var cm              CourseMenu
    var existing        bool
    var match           bool
    var mm              MiniMenu
    var mmitems         []MmItem
    var id              string
    var coursetitle     string
    // Iterate over all courses.tags[], and derive a list of unique tags called MiniMenuTitles
    for _, thiscourse := range cc.Cc {
        // Iterate over all tags within a course
        for _, thistag := range thiscourse.Tags {
             existing = false
             //Iterate over all minimenutitle, add new ones to the list
             for _, thismmt := range minimenutitles {
                 if thistag == thismmt {
                    existing = true
                 }
             }
             //No matches? Then add thistag to the list
             if existing == false {
             minimenutitles = append (minimenutitles, thistag)
             }
        }
    }
    fmt.Println("Items: %+v\n\n", minimenutitles)
    // Interate over each MiniMenuTitle 
    for _, mmt := range  minimenutitles {
          //Iternate over every course in the catalog
          for _, thiscourse := range cc.Cc {
                match = false
                //Iterate over every course's Tags
                for _, thistag := range thiscourse.Tags {
                   //If this course has a matching tag entry, grab data
                   if mmt == thistag {
                      match = true
                      id = thiscourse.Id
                      coursetitle = thiscourse.CourseTitle
                      break
                   }
                }//End of tag iteration so add item if match is true
                if match == true {
                   mmitem.Id = id
                   mmitem.CourseTitle = coursetitle
                   mmitems = append(mmitems, mmitem)
                }
          } //End of Course Itermation
          mm.MiniMenuTitle = mmt
          mm.MmItems = mmitems
          fmt.Printf("mm: %v\n", mm)
          cm = append (cm, mm)
          mmitems = nil
   }//End of minimenu iteration
    return cm
}


//API - New MegaMenu - Returns a canned megamenu pre-sorted by course type. 
func (cm CourseMenu ) coursemenu() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    (w).Header().Set("Access-Control-Allow-Headers","*")
    (w).Header().Set("Access-Control-Allow-Origin", "*")
    var js []byte
    var err error
    js, err = json.Marshal(cm)
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





// ------------------------------------------------------------


func main() {
////      router := mux.NewRouter().StrictSlash(true)
  cs := Load()
  fmt.Println("--------------------------------------------------")
  fmt.Println("Course Loaded into MAIN: " + cs.Cc[0].Id)
  events := LoadEvents()
  blogcontent := Loadblogs()
  fmt.Println("Blogs Loaded into MAIN: " + blogcontent[0].Id)
  cmenu :=  cs.menumaker()
  fmt.Println("Menu: %+v\n", cmenu)
  fmt.Printf("BLOGMENU\n")
  fmt.Println("--------------------------------------------------")
  bmenu := blogcontent.blogmenumaker()
  fmt.Println("--------------------------------------------------")
//	fmt.Printf("blogmenu: %s\n", bmenu[0].BlogCategory)




// All the static folders
	http.Handle("/downloads/",   http.StripPrefix("/downloads/",   http.FileServer(http.Dir("downloads"))))
	http.Handle("/img/",         http.StripPrefix("/img/",         http.FileServer(http.Dir("deploy/img"))))
	http.Handle("/images/",      http.StripPrefix("/images/",      http.FileServer(http.Dir("deploy/images"))))
	http.Handle("/fonts/",       http.StripPrefix("/fonts/",       http.FileServer(http.Dir("deploy/fonts"))))
	http.Handle("/icons/",       http.StripPrefix("/icons/",       http.FileServer(http.Dir("deploy/icons"))))
	http.Handle("/css/",         http.StripPrefix("/css/",         http.FileServer(http.Dir("deploy/css"))))
	http.Handle("/js/",          http.StripPrefix("/js/",          http.FileServer(http.Dir("deploy/js"))))
  http.Handle("/assets/",      http.StripPrefix("/assets/",      http.FileServer(http.Dir("deploy/assets"))))
  http.Handle("/coming_soon/", http.StripPrefix("/coming_soon/", http.FileServer(http.Dir("deploy/coming_soon"))))
  http.Handle("/sass/",        http.StripPrefix("/sass/",        http.FileServer(http.Dir("deploy/sass"))))
  http.Handle("/video/",       http.StripPrefix("/video/",       http.FileServer(http.Dir("deploy/video"))))
  http.Handle("/layerslider/", http.StripPrefix("/layerslider/", http.FileServer(http.Dir("deploy/layerslider"))))
//http.Handle("/",             http.StripPrefix("/",             http.FileServer(http.Dir("deploy"))))

	// Templates
	http.Handle("/courses/", http.StripPrefix("/courses/", cs.CourseTemplate()))
	http.Handle("/blog/", http.StripPrefix("/blog/", BlogTemplate()))
	http.Handle("/", http.StripPrefix("/", Template()))
// JSON RESTful Interfaces
  http.Handle("/api/v1/course/search/",      http.StripPrefix("/api/v1/course/search/",     cs.search()))
  http.Handle("/api/v1/course/megamenu/",    http.StripPrefix("/api/v1/course/megamenu/",   cs.getmegamenu()))
  http.Handle("/api/v1/course/coursemenu/",  http.StripPrefix("/api/v1/course/coursemenu/", cmenu.coursemenu()))
  http.Handle("/api/v1/course/summary/id/",  http.StripPrefix("/api/v1/course/summary/id/", cs.getsummary()))
  http.Handle("/api/v1/course/detail/id/",   http.StripPrefix("/api/v1/course/detail/id/",  cs.getcoursedetail()))
  http.Handle("/api/v1/blog/search/",        http.StripPrefix("/api/v1/blog/search/",       blogcontent.blogsearch()))
  http.Handle("/api/v1/blog/id/",            http.StripPrefix("/api/v1/blog/id/",           blogcontent.getblog()))
  http.Handle("/api/v1/events/",             http.StripPrefix("/api/v1/events/",            events.getevents()))
  http.Handle("/api/v1/blog/menu/",          http.StripPrefix("/api/v1/blog/menu/",         bmenu.blogmenu()))

// Stripe Chckout
	http.Handle("/checkout", Checkout())

	log.Printf("serving...")
	http.ListenAndServe(":8888", Log(http.DefaultServeMux))
}

