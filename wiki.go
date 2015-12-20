package main

import "io/ioutil"

/*
  Data Structures
*/

// Wikiページ.ページデータがメモリに保存される
type Page struct {
	Title string
	Body  []byte // byte slice.
	// Not string bc it is the type expected by the io library
}

// This is a method that takes as its receiver p, a pointer to Page.
// It takes no parameters, and returns a value of type error
// If everything goes well, it return nil.
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
	// 0600 indicates read-write permission for the current user only
}

// Returns a pointer to a Page literal
func loadPages(title string) *Page {
	filename := title + ".txt"
	body, _ := ioutil.ReadFile(filename)
	return &Page{Title: title, Body: body}
}
