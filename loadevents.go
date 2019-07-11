package main

import (
  "fmt"
  "log"
  "os"
  "path"
  "io/ioutil"
  "github.com/ghodss/yaml"
)




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



