package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/(.+)$")

type Database struct {
    DbName string
    Data []byte
}

func loadDatabase(name string) (*Database, error) {
    filename := name + ".sql"
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Database{DbName: name, Data: data}, nil
}


func viewHandler(w http.ResponseWriter, r *http.Request, name string) {
    db, err := loadDatabase(name)
    if err != nil {
        fmt.Fprintf(w, "File not found")
        return
    }
    fmt.Fprintf(w, "<p>%s</p>", db.Data)
}


func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

func main() {
   http.HandleFunc("/view/", makeHandler(viewHandler))
   http.ListenAndServe(":8080", nil)
}
