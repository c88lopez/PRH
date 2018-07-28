package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
)

// Start Start
func Start() {
	http.HandleFunc("/", convertHandler)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

// convertHandler / endpoint handler
func convertHandler(w http.ResponseWriter, req *http.Request) {
	var v interface{}

	if req.Method == "POST" {
		file, _, err := req.FormFile("archivito")
		if err != nil {
			log.Println("Error uploading file: ", err)
		}
		defer file.Close()

		w.Header().Set("Content-Type", "text/csv")
		io.Copy(w, Convert(file))
	} else {
		w.Header().Set("Content-type", "text/html")
		t, _ := template.ParseFiles("templates/forms/archivito.html")
		t.Execute(w, v)
	}
}
