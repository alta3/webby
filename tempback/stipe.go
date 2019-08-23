package main

import (
//  "regexp"


//	"github.com/tdewolff/minify"
//	"github.com/tdewolff/minify/html"

  "encoding/json"
	"fmt"
	"log"
	"net/http"
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


