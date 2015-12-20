package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

// Template Caching(Panic when passed a non-nil error value)
var templates = template.Must(
	template.ParseFiles("edit.html", "view.html"))

// Validation(Panic if expression compilation fails)
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

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
	//title := r.URL.Path[len("/view/"):]
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	//p, _ := loadPage(title) //Shows a page containing HTML as it tries to fill template with no data
	p, err := loadPage(title)
	if err != nil {
		// Redirecting the client to the edit Page
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

// loads the page (if it the page doesn't exist, create an empty Page struct
func editWikiHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/edit/"):]
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
	//	fmt.Fprintf(w, "<h1>Editing %s</h1>"+
	//		"<form action=\"/save/%s\" method=\"POST\">"+
	//		"<textarea name=\"body\">%s</textarea><br>"+
	//		"<input type=\"submit\" value=\"Save\">"+
	//		"</form>",
	//		p.Title, p.Title, p.Body)
}

// Handle the submission of forms located on the edit pages.
func saveWikiHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/save/"):]
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("body") // Returns type string
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//	t, err := template.ParseFiles(tmpl + ".html") // retun a *template.Template
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	err = t.Execute(w, p) // Execute the template ,writing generated HTML to w
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // the title is the 2nd subexpression:w
}

/**********************************
***********************************
  Data Structures
***********************************
************************************/
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
