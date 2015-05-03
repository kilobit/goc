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

type CSVHandler struct {
	reader *csv.Reader
	rows [][]string
}

func NewCSVHandler (r io.Reader) *CSVHandler {
	csvh := &CSVHandler{csv.NewReader(r), [][]string{}}

	csvh.readCSVRows()

	return csvh
}

func (csvh *CSVHandler) readCSVRows() {

	for {
		row, err := csvh.reader.Read()
		if err != nil {
			return
		}

		csvh.rows = append(csvh.rows, row)
	}
}

func (csvh *CSVHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	
	jsone := json.NewEncoder(w)

	jsone.Encode(csvh.rows)
}

func main() {

	mux := http.NewServeMux()

	flag.Parse()
	
	if *root != "" {
		mux.Handle(*path, http.StripPrefix(*path, http.FileServer(http.Dir(*root))))		
	}
	
	mux.Handle("/data", NewCSVHandler(os.Stdin))
	
	log.Fatal(http.ListenAndServe(":8080", mux))
}
