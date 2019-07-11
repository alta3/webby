package main

import ()




//  ____  _      ____   _____  _____ 
// |  _ \| |    / __ \ / ____|/ ____|
// | |_) | |   | |  | | |  __| (___  
// |  _ <| |   | |  | | | |_ |\___ \ 
// | |_) | |___| |__| | |__| |____) |
// |____/|______\____/ \_____|_____/ 

//------------------BLOGS------------------------
type Blog struct {
	Id             string         `yaml:"id"            json:"id"`
	Author         string         `yaml:"author"        json:"author"`
	Category       string         `yaml:"category"      json:"category"`
	Date           string         `yaml:"date"          json:"date"`
	Title          string         `yaml:"title"         json:"title"`
	Weight         string         `yaml:"weight"        json:"weight"`
	Intro          string         `yaml:"intro"         json:"intro"`
	VideoLink      string         `yaml:"video-link"    json:"video-link"`
	HtmlContent    string         `yaml:"html-content"  json:"html-content"`
}

type Blogs []Blog

// Blog Menu
// Returns a blog menu that does all the work for the front end developers
// Best to use this when the server boots since this menu does not often change

type BlogsByCategory struct {
	BlogCategory     string        `json:"blog-category"`
	Blogs            []Blog        `json:"blogs"`
}

type BlogMenus     []BlogsByCategory

type BlogCategory  []string


