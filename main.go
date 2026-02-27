package main

import(
	"log"
	"net/http"
)





func main() {
	const port = "8080"
	const filepathRoot = "."

	
	mux := http.NewServeMux()
	file := http.Dir(filepathRoot)
	mux.Handle("/", http.FileServer(file))
	server := &http.Server{
		Addr : ":" + port,
		Handler : mux,
	}


	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}