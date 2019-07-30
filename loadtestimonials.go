package main

import (
  "fmt"
  "log"
  "os"
  "path"
  "io/ioutil"
  "github.com/ghodss/yaml"
)




//Load TESTIMONIALS
//------------------------------------------------------------
func LoadTestimonials() Testimonials {
      // Create a OS compliant path: microsoft "\" or linux "/"
      dirname := path.Join("deploy", "testimonials")
      d, err := os.Open(dirname)
      if err != nil {
          log.Printf("No testimonials directory! %s, %s" , d, err)
          os.Exit(1)
      }
      var    t   Testimonials
      fmt.Println("-------------LOADING TESTIMONIALS------------------")
      fmt.Printf(" Reading testimonials YAML in directory: %s\n", dirname)
      thisfile := path.Join(dirname, "testimonials.yaml")
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
     // unmarshal byteArray using the YAML tags 
	    err = yaml.Unmarshal(yammy, &t)
      if err != nil {
				 log.Printf("Unmarshal: %v", err)
				  }

			fmt.Printf(" Successfully read: %s\n", thisfile) 
      fmt.Printf(" Events: %+v\n", t )
      fmt.Println("--------------------------------------------------")
      return t
}



