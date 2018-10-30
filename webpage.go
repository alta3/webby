package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

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
			r.URL.Path = "404"
			Template().ServeHTTP(w, r)
		}
	})
}

func Template() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			http.ServeFile(w, r, "deploy/html_menu_1/index.html")
			return
		} else if r.URL.Path == "deploy/html_menu_1/robots.txt" {
			http.ServeFile(w, r, "robots.txt")
			return
		} else if r.URL.Path == "deploy/html_menu_1/sitemap.xml" {
			http.ServeFile(w, r, "sitemap.xml")
			return
		}

		lp := path.Join("deploy/html_menu_1/templates", "layout.html")
		fp := path.Join("deploy/html_menu_1/templates", r.URL.Path)

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
		tmpl, err := template.ParseFiles(lp, fp)
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

func main() {
	// All the static folders
	http.Handle("/downloads/", http.StripPrefix("/downloads/", http.FileServer(http.Dir("downloads"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("deploy/html_menu_1/img"))))
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
