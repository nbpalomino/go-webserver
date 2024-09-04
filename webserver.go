package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var port int
var isSecure bool

func main() {
	flag.IntVar(&port, "p", 8080, "Puerto del WebServer")
	flag.BoolVar(&isSecure, "s", false, "Servidor HTTPS")
	flag.Parse()

	log.Print("-> Iniciando WebServer en puerto ", port)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("-> Directorio ", pwd)

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			log.Printf("-> Showing upload.html")
			tmpl := template.Must(template.ParseFiles("./static/upload.html"))
			tmpl.Execute(w, nil)
			return
		}

		// Parse the multipart form
		r.ParseMultipartForm(10 << 20 << 20) // 10GB maximum file size

		// Retrieve the file from form data
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println("Error retrieving the file:", err)
			http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
			return
		}
		log.Printf("-> Retrieving file from form.")
		defer file.Close()

		// Create a temporary file on the server
		tempFile, err := os.Create(filepath.Join("uploads", handler.Filename))
		if err != nil {
			fmt.Println("Error creating a temporary file:", err)
			http.Error(w, "Error creating a temporary file", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		// Copy the uploaded file to the temporary file
		_, err = io.Copy(tempFile, file)
		if err != nil {
			fmt.Println("Error saving the file:", err)
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}
		log.Printf("-> File copied to destiny.")

		// Render a response back to the user
		fmt.Fprintf(w, "Successfully uploaded file: %s", handler.Filename)
	})

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	http.Handle("/", http.FileServer(http.Dir(pwd)))

	if isSecure {
		log.Print(http.ListenAndServeTLS(":443", "server.crt", "server.key", http.FileServer(http.Dir(pwd))))
	} else {
		log.Print(http.ListenAndServe(fmt.Sprint(":", port), nil))
	}
}
