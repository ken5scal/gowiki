package main

import "io/ioutil"

/*
  Data Structures
*/

// Wikiページ.ページデータがメモリに保存される
type Page struct {
	Title string
	Body  []byte
}

// This is a method that takes as its receiver p, a pointer to Page.
// It takes no parameters, and returns a value of type error
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}
