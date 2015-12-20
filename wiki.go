package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	//	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page.")}
	//	p1.save()
	//	p2, _ := loadPage("TestPage")
	//	fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", viewWikiHandler)
	http.HandleFunc("/edit/", editWikiHandler)
	http.HandleFunc("/save/", saveWikiHandler)
	http.ListenAndServe(":8080", nil)
}

// Allow users to view a wiki page.
func viewWikiHandler(w http.ResponseWriter, r *http.Request) {
	// Extracting the Page title from URL
	// Also droppoing the leading ?view?
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

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
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	// TODO Handle error
	// underscore(_) symbol is used to throw away the value
	//body, _ := ioutil.ReadFile(filename)
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}
