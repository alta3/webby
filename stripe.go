package main

import (
//  "regexp"


//	"github.com/tdewolff/minify"
//	"github.com/tdewolff/minify/html"

  "encoding/json"
  "fmt"
  "html/template"
  "log"
  "net/http"
  "os"
  "path"
// "gopkg.in/russross/blackfriday.v2"

// "gopkg.in/yaml.v2"
// "github.com/gorilla/mux"

  stripe "github.com/stripe/stripe-go"
  "github.com/stripe/stripe-go/charge"
)


// MOVE TO stripe.go 


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



