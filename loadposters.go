package main

import (
  "fmt"
  "log"
  "os"
  "path"
  "io/ioutil"
  "github.com/ghodss/yaml"
)




//Load POSTERS
//------------------------------------------------------------
func LoadPosters() Posters {
      // Create a OS compliant path: microsoft "\" or linux "/"
      dirname := path.Join("deploy", "posters")
      d, err := os.Open(dirname)
      if err != nil {
          log.Printf("No posterss directory! %s, %s" , d, err)
          os.Exit(1)
      }
      var    p   Posters
      fmt.Println("---------------LOADING EVENTS---------------------")
      fmt.Printf(" Reading poster YAML in directory: %s\n", dirname)
      thisfile := path.Join(dirname, "posters.yaml")
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
	    err = yaml.Unmarshal(yammy, &p)
      if err != nil {
				 log.Printf("Unmarshal: %v", err)
				  }

			fmt.Printf(" Successfully read: %s\n", thisfile) 
      fmt.Printf(" Events: %+v\n", p )
      fmt.Println("--------------------------------------------------")
      return p
}



