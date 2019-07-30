package main

import ()



// ---------------------------------------------------------


type Posters struct {
	Posters []struct {
		Name        string `yaml:"name"        json: "name"`
		Description string `yaml:"description" json:"description"`
		Downloadurl string `yaml:"downloadurl" json::"downloadurl"`
	} `yaml:"posters" json::"posters"`
}



