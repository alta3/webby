package main

import ()



// ---------------------------------------------------------


type Testimonials struct {
	Testimonials []struct {
		Date  string `yaml:"date"  json:"date"`
		Quote string `yaml:"quote" json:"quote"`
		Stars string `yaml:"stars" json:"stars"`
		Name  string `yaml:"name"  json:"name"`
	} `yaml:"testimonials" json:"testimonials"`
}
