package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func Start() {
	http.HandleFunc("/", convertHandler)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

// convertHandler / endpoint handler
func convertHandler(w http.ResponseWriter, req *http.Request) {
	var v interface{}

	if req.Method == "POST" {
		log.Println("post")

		file, handler, err := req.FormFile("archivito")
		if err != nil {
			log.Println("Error uploading file: ", err)
		}
		defer file.Close()

		f, err := os.OpenFile(fmt.Sprintf("%s%c%s", csvsFolderPath, os.PathSeparator, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("Error saving uploaded file: ", err)
		}
		defer f.Close()

		io.Copy(f, file)
		processFiles()

		fConv, err := os.Open(strings.Replace(f.Name(), ".csv", " converted.csv", 1))
		if err != nil {
			log.Println("Error getting converted file: ", err)
		}

		fConvStat, _ := fConv.Stat()

		w.Header().Set("Content-Disposition", "attachment; filename="+fConv.Name())
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Length", strconv.FormatInt(fConvStat.Size(), 10))

		io.Copy(w, fConv)
	} else {
		w.Header().Set("Content-type", "text/html")
		t, _ := template.ParseFiles("templates/forms/archivito.html")
		t.Execute(w, v)
	}
}

func processFiles() {
	csvsFolder, err := os.Open(csvsFolderPath)
	if err != nil {
		log.Fatal(err)
	}

	csvFiles, err := csvsFolder.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, csvFile := range csvFiles {
		file, err := os.Open(fmt.Sprintf("%s%c%s", csvsFolderPath, os.PathSeparator, csvFile.Name()))

		if skipFile(file) {
			continue
		}

		if err != nil {
			log.Fatal("Error opening csv file", fmt.Sprintf("%s%c%s", csvsFolderPath, os.PathSeparator, csvFile.Name()), err)
		}
		defer file.Close()

		waitGroup.Add(1)
		go CreateConvertedFile(file)
	}

	waitGroup.Wait()
}
