package main

import (
    "net/http"
    "html/template"
    "regexp"
)

var templates = template.Must(template.ParseFiles("home.html", "about.html", "contact.html"))
//var validPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(home|about|contact)$")

func renderTemplate(w http.ResponseWriter, tmpl string) {
    err := templates.ExecuteTemplate(w, tmpl + ".html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func handler(w http.ResponseWriter, r *http.Request, path string) {
    renderTemplate(w, path)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
          http.NotFound(w, r)
          return
        }
        fn(w, r, m[1])
    }
}

func main() {
    http.HandleFunc("/home", makeHandler(handler))
    http.HandleFunc("/about", makeHandler(handler))
    http.HandleFunc("/contact", makeHandler(handler))

    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
    http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

    http.ListenAndServe(":8080", nil)
}
