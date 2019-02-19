package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"fmt"
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

	if(isSecure) {
		log.Print(http.ListenAndServeTLS(":443", "server.crt", "server.key", http.FileServer(http.Dir(pwd))))
	} else {
		log.Print(http.ListenAndServe(fmt.Sprint(":",port), http.FileServer(http.Dir(pwd))))
	}
}
