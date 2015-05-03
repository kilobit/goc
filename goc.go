package main

import "net/http"
import "log"
import "io"
import "flag"
import "encoding/csv"
import "os"
import "encoding/json"

var (
	root = flag.String("root", "", "Document Root")
	path = flag.String("path", "/", "Root URL Path")
)

func fooHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Bar!")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello World!"))
}

type CSVHandler struct {
	reader *csv.Reader
}

func NewCSVHandler (r io.Reader) *CSVHandler {
	return &CSVHandler{csv.NewReader(r)}
}

func (csvh *CSVHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jsone := json.NewEncoder(w)

	row, err := csvh.reader.Read()
	if err != nil {
		return
	}
	
	jsone.Encode(row)
	
	jsone.Encode([...]string{"Hello", "World!", "CSV", "Style"})
}

func main() {

	mux := http.NewServeMux()

	flag.Parse()
	
	if *root != "" {
		mux.Handle(*path, http.StripPrefix(*path, http.FileServer(http.Dir(*root))))		
	}
	
	mux.HandleFunc("/foo", fooHandler)
	
	mux.HandleFunc("/chew", rootHandler)

	mux.Handle("/data", NewCSVHandler(os.Stdin))
	
	log.Fatal(http.ListenAndServe(":8080", mux))
}
