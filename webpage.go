package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
        "unicode"
        "bytes"
	"net/http"
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
//			http.ServeFile(w, r, "deploy/html_menu_1/index.html")
//			return
		} else if r.URL.Path == "deploy/html_menu_1/robots.txt" {
			http.ServeFile(w, r, "robots.txt")
			return
		} else if r.URL.Path == "deploy/html_menu_1/sitemap.xml" {
			http.ServeFile(w, r, "sitemap.xml")
			return
		}

                // create the full pathnames of the files to be merged
		lp := path.Join("deploy/html_menu_1/templates", "layout.html")
		fp := path.Join("deploy/html_menu_1/templates", r.URL.Path)

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

                // OK files exsist, begin the merging, remove {{ jinja-like stuff }}
                // template.ParseFiles returns TWO values, the combined template, and error code
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

func CourseTemplate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			r.URL.Path = "courses"
			Template().ServeHTTP(w, r)
			return
		}
		lp := path.Join("templates", "layout.html")
		cp := path.Join("courses", "course_layout.html")
		fp := path.Join("courses", r.URL.Path+".html")
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


// ----------------COURSE STRUCT------------------------------- 

type Include struct {
  Item          string          `json:"item"`
  Description   string          `json:"description"`
}

type PriceTag struct {
  Name          string          `json:"name"`
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

type Testimonials struct {
  Quotes        []string        `json:"quotes"`
}

type Lab struct {
  Title string `rethinkdb:"title" json:"title"`
  File  string `rethinkdb:"file" json:"file"`
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
  Testimonials  Testimonials    `json:"testimonials"`
  VideoLink     string          `json:"videolink"`
  Overview      string          `json:"overview"`
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



func main() {

//      router := mux.NewRouter().StrictSlash(true)

        yammy, err := ioutil.ReadFile("course.yaml")
        if err != nil {
           log.Printf("yammy.Get err   #%v ", err)
        }
        fmt.Println("-------------------------------------------")
        fmt.Println("Successfully Opened:  course.yaml")

        // convert yammy into JSON
        jsonFile, err := ToJSON(yammy)
        
        
        // initialize course type now that it is converted to JSON
        var course Course 
        
        // unmarshal byteArray using the JSON tags 
        json.Unmarshal(jsonFile, &course)

        // Any zero output is bad and indicates a YAML error.        
            fmt.Println("-------------------------------------------")
            fmt.Println("              Course: " + course.Id)
            fmt.Println("            Duration: " + strconv.Itoa(course.Duration.Hours))
            fmt.Printf("      Book Price Tags %d\n", len(course.Price.Book.PriceTags))
            fmt.Printf("    Public Price Tags %d\n", len(course.Price.Public.PriceTags))
            fmt.Printf("   Private Price Tags %d\n", len(course.Price.Private.PriceTags))
            fmt.Printf("Self Paced Price Tags %d\n", len(course.Price.Selfpaced.PriceTags))
            fmt.Printf("Extend LMS Price Tags %d\n", len(course.Price.ExtendLmsAccess.PriceTags))
            fmt.Printf("         Testimonials %d\n", len(course.Testimonials.Quotes))
            fmt.Printf("             Chapters %d\n", len(course.Chapters))
            fmt.Printf("                 Labs %d\n", len(course.Labs))





	// All the static folders
	http.Handle("/downloads/", http.StripPrefix("/downloads/", http.FileServer(http.Dir("downloads"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("deploy/html_menu_1/img"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("deploy/html_menu_1/images"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("deploy/html_menu_1/fonts"))))
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("deploy/html_menu_1/icons"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("deploy/html_menu_1/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("deploy/html_menu_1/js"))))
        http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("deploy/html_menu_1/assets"))))
        http.Handle("/coming_soon/", http.StripPrefix("/coming_soon/", http.FileServer(http.Dir("deploy/html_menu_1/coming_soon"))))
        http.Handle("/sass/", http.StripPrefix("/sass/", http.FileServer(http.Dir("deploy/html_menu_1/sass"))))
        http.Handle("/video/", http.StripPrefix("/video/", http.FileServer(http.Dir("deploy/html_menu_1/video"))))
        http.Handle("/layerslider/", http.StripPrefix("/layerslider/", http.FileServer(http.Dir("deploy/html_menu_1/layerslider"))))
        // http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("deploy/html_menu_1"))))

	// Templates
	http.Handle("/courses/", http.StripPrefix("/courses/", CourseTemplate()))
	http.Handle("/blog/", http.StripPrefix("/blog/", BlogTemplate()))
	http.Handle("/", http.StripPrefix("/", Template()))

	// Stripe Chckout
	http.Handle("/checkout", Checkout())

	log.Printf("serving...")
	http.ListenAndServe(":8888", Log(http.DefaultServeMux))
}
