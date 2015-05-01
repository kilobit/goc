package main

import "net/http"
import "log"
import "io"
import "flag"

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

func main() {

	mux := http.NewServeMux()

	flag.Parse()
	
	if *root != "" {
		mux.Handle(*path, http.StripPrefix(*path, http.FileServer(http.Dir(*root))))		
	}
	
	mux.HandleFunc("/foo", fooHandler)
	
	mux.HandleFunc("/chew", rootHandler)
	
	log.Fatal(http.ListenAndServe(":8080", mux))
}
