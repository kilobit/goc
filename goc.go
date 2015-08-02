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
	writer *csv.Writer
	rows [][]string
}

func NewCSVHandler (r io.Reader, w io.Writer) *CSVHandler {
	csvh := &CSVHandler{csv.NewReader(r), 
		csv.NewWriter(w),
		[][]string{}}

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

	switch {

	case r.Method == "GET":
		csvh.ReadCSV(w, r)
		return

	case r.Method == "POST":
		csvh.WriteCSV(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotImplemented)
		io.WriteString(w, "The given method has not been implemented.")
	}
}

func (csvh *CSVHandler) ReadCSV(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	
	jsone := json.NewEncoder(w)

	jsone.Encode(csvh.rows)
}

func (csvh *CSVHandler) WriteCSV(w http.ResponseWriter, r *http.Request) {

	

}

func main() {

	mux := http.NewServeMux()

	flag.Parse()
	
	if *root != "" {
		mux.Handle(*path, http.StripPrefix(*path, http.FileServer(http.Dir(*root))))		
	}
	
	mux.Handle("/data", NewCSVHandler(os.Stdin, os.Stdout))
	
	log.Fatal(http.ListenAndServe(":8080", mux))
}
