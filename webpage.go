package main

import (
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
        "time"
	"path"
        "io/ioutil"
        "strconv"
        "github.com/ghodss/yaml"
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
			"sip-selfpaced":           14995,
			"sdnnfv-virtual":          259500,
                        "sdnnfv-selfpaced":        44500,
			"sip-virtual":             239500,
			"volte-virtual":           239500,
			"avaya-virtual":           249500,
			"openstack-virtual":       259500,
			"tip-jar":                 100,
			"sip-book":                5000,
			"dd-sip-selfpaced":        794325,
			"godfrey-sip-selfpaced":   22250,
			"godfrey-avaya-selfpaced": 34750,
                        "python1-selfpaced":       14995,
                        "python1-virtual":         259500,
                        "python2-selfpaced":       59500,
                        "python2-virtual":         259500,
                        "python3-selfpaced":       59500,
                        "python3-virtual":         259500,
                        "ansible-virtual":         259500,
                        "ansible-selfpaced":       14995,
                        "rhcsa-virtual":           259500,
                        "rhcsa-selfpaced":         595000,
                        "k8s-virtual":             199500,
                        "k8s-selfpaced":           39500,
                        "ipsec-virtual":           199500,
                        "ipsec-selfpaced":         29500,
                        "network-virtual":         199500,
                        "network-selfpaced":       29500,
                        "ceph-virtual":            259500,
                        "ceph-selfpaced":          59500,
                        "5g-virtual":              179500,
                        "5g-selfpaced":            14995,
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
   GUID  string                  `rethinkdb:"guid" json:"guid"`
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
  Name          string          `rethinkdb:"name" json:"name"`
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
  VideoLink     string          `json:"videolink"`
  Overview      string          `json:"overview"`
  Tags          []string        `json:"tags"`
  Courseicon    string          `json:"courseicon"`                  // TODO make sure courseicons are under images/courseicons and all available to put in yaml or json file as coursicon: xxxx.png
}

type Courses struct {
  Cc           []Course         `json:"courses"`
}



type PublicCourse struct {
  Course
  Chapters      []Chapter       `json:"chapters,omitempty"` // TODO update to single-slide-mode
  Labs          []Lab           `json:"labs,omitempty"`    // TODO Write
} 



type CourseCatalog interface {
  Select(id string)       []Course
  Load()                  []Course
  Search(ss string)       []Course
}




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
      c := make([]Course,10) 
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
      fmt.Println("YAMMY  Course: " + jsonCatalogFile.Cc[0].Id)
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

// ========JUST WORKING ON THIS RIGHT NOW===========================================
//                          V


// I just started working on this, making the /courses/courses-list.html page render search results
// NEXT, make a page show course-detail.  Hopefully I will have both HTML pages set up before returning
// to this go code.

func  (cs Courses) CourseTemplate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			r.URL.Path = "courses"
			Template().ServeHTTP(w, r)
			return
		}

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



//Not complete yet, still in testing phase. 
//Going to return the course Id, Title, Stars, Duration, Description, selfpaced price, and live price, courseicon.
//func (cs Courses ) getpopup() http.Handler {
//   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//    type PopupItems struct {
//      Id           string          `json:"id"`
//      Name         string          `json:"name"`
//      Stars        int             `json:"stars"`
//      Duration     int             `json:"duration"`
//      Overview     string          `json:"overview"`
//      Price        int             `json:"price"`
//      Courseicon   string          `json:"courseicon"`
//    }
//    popi := PopupItems{}
//    popis := []PopupItems{}
//    var js []byte
//    var err error
//    //Iterate over all courses, Copy Id, Name, Stars, Duration, Overview, Price, and Courseicon
//    for _, ThisCourse := range cs.Cc {
//       fmt.Println("--------------------------------------------------")
//       fmt.Printf("_Course PopUp_  = %s, %s, %s, %s, %s, %s, %s\n", ThisCourse.Id, ThisCourse.Name, ThisCourse.Stars, ThisCourse.Duration, ThisCourse.Overview, ThisCourse.Price, ThisCourse.Courseicon)
//       popi.Id=ThisCourse.Id
//       popi.Name=ThisCourse.Name
//       popi.Stars=ThisCourse.Stars
//       popi.Duration=ThisCourse.Duration
//       popi.Overview=ThisCourse.Overview
//       popi.Price=ThisCourse.Price
//       popi.Courseicon=ThisCourse.Courseicon
//       popis = append(popis,popi)
//    }
//    //If no courses match, SEND THEM ALL! 
//       js, err = json.Marshal(popis)
//    if err != nil {
//       http.Error(w, err.Error(), http.StatusInternalServerError)
//       fmt.Printf("Error %s:\n", err)
//       return
//    }
//    w.Header().Set("Content-Type", "application/json")
//    w.Write(js)
//    return
//    })
//}






//Not complete yet, still in testing phase.
//Going to return blog Id and Title.
//func (cs Courses ) getblogs() http.Handler {
//   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//    type BlogItems struct {
//      Id           string          `json:"id"`
//      Name         string          `json:"name"`
//    }
//    blogi := BlogItems{}
//    blogis := []BlogItems{}
//    var js []byte
//    var err error
//    //Iterate over all courses, Copy Id, Name, Stars, Duration, Overview, Price, and Courseicon
//    for _, ThisCourse := range cs.Cc {
//       fmt.Println("--------------------------------------------------")
//       fmt.Printf("Blog Search Results  = %s, %s\n", ???.Id., ???.Name)
//       blogi.Id=???.Id
//       blogi.Name=???.Name
//       blogis = append(blogis,blogi)
//    }
//    //If no courses match, SEND THEM ALL! 
//       js, err = json.Marshal(blogis)
//    if err != nil {
//       http.Error(w, err.Error(), http.StatusInternalServerError)
//       fmt.Printf("Error %s:\n", err)
//       return
//    }
//    w.Header().Set("Content-Type", "application/json")
//    w.Write(js)
//    return
//    })
//}







// Returns a Course Id, Tag, and Name for building a Course Mega-Menu. Tag is 
func (cs Courses ) getmenu() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    type MenuItems struct {
      Id           string          `json:"id"`
      Name         string          `json:"name"`
      Tags         []string        `json:"tags"`
    }
    mi := MenuItems{} 
    mis := []MenuItems{} 
    var js []byte
    var err error
    //Iterate over all courses, copy Id, Name, and any Tags
    for _, ThisCourse := range cs.Cc {
       fmt.Println("--------------------------------------------------")
       fmt.Printf("Menu Item = %s, %s, %s\n", ThisCourse.Id, ThisCourse.Name, ThisCourse.Tags)
       mi.Id=ThisCourse.Id
       mi.Name=ThisCourse.Name
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


// Given a search string, returns course data for all matching courses, without the course outline data.
func (cs Courses ) getsearch() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    //Retreive variable from GET URL
    ss := r.FormValue("search")
    fmt.Printf("func searchstring searching for: %s\n",ss)                
    // Create a new composite Course type. Interestingly, by adding existing subordinate
    // types to the cloned struct, items will OMIT them from the marshalling.
    // see: https://mycodesmells.com/post/working-with-embedded-structs
    c := []PublicCourse{} 
    hits := 0
    totalhits := 0
    var js []byte
    var err error
    //Iterate over all courses, looking for the search string (ss)
    //If a match is found, add the course to c.Cc[i]
    for i, ThisCourse := range cs.Cc {
       hits = strings.Count( strings.ToLower(fmt.Sprintf("%v", cs.Cc[i])), ss )
       totalhits = totalhits + hits
       if  hits > 0 {
           // Here is how you graft an existing type into a new "composite" type.
           pch :=  PublicCourse{Course: ThisCourse} 
           c = append(c,pch)
           fmt.Printf("%s Course has %d hits\n", ThisCourse.Id, hits )
       }
    }
    //If no courses match, SEND THEM ALL! 
    if totalhits == 0 {
       js, err = json.Marshal(cs.Cc)
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

func main() {

////      router := mux.NewRouter().StrictSlash(true)

        cs := Load()
        fmt.Println("--------------------------------------------------")
        fmt.Println("Course Loaded into MAIN: " + cs.Cc[0].Id)

       
//        _, err := cs.Select("5g")
//           if err != nil {
//             log.Printf("SORRY: %s\n ", err)
//          }
           
//        _, err = cs.Search("has-book")
//           if err != nil {
//             log.Printf("SORRY: %s\n ", err)
//          }
 
        

	// All the static folders
	http.Handle("/downloads/", http.StripPrefix("/downloads/", http.FileServer(http.Dir("downloads"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("deploy/img"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("deploy/images"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("deploy/fonts"))))
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("deploy/icons"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("deploy/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("deploy/js"))))
        http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("deploy/assets"))))
        http.Handle("/coming_soon/", http.StripPrefix("/coming_soon/", http.FileServer(http.Dir("deploy/coming_soon"))))
        http.Handle("/sass/", http.StripPrefix("/sass/", http.FileServer(http.Dir("deploy/sass"))))
        http.Handle("/video/", http.StripPrefix("/video/", http.FileServer(http.Dir("deploy/video"))))
        http.Handle("/layerslider/", http.StripPrefix("/layerslider/", http.FileServer(http.Dir("deploy/layerslider"))))
        // http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("deploy"))))

	// Templates
	http.Handle("/courses/", http.StripPrefix("/courses/", cs.CourseTemplate()))
	http.Handle("/blog/", http.StripPrefix("/blog/", BlogTemplate()))
	http.Handle("/", http.StripPrefix("/", Template()))

        // JSON RESTful Interfaces
        http.Handle("/api/v1/search/", http.StripPrefix("/search/", cs.getsearch()))

        // Get Menu 
        http.Handle("/api/v1/menu/", http.StripPrefix("/menu/", cs.getmenu()))

        //Get Popup Course
        http.Handle("/api/v1/popup/",http.StripPrefix("/popup/", cs.getpopup()))

        //Get Searh Blog
        //http.Handle("/api/v1/searchblog/",http.StripPrefix("/searchblog/", cs.getblogs()))

        //Get Blog Details
        //http.Handle("/api/v1/blogdet/",http.StripPrefix("/blogdet/", cs.getblogd()))

	// Stripe Chckout
	http.Handle("/checkout", Checkout())

	log.Printf("serving...")
	http.ListenAndServe(":8888", Log(http.DefaultServeMux))
}
