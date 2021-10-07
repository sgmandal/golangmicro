package main

import (
	"log"
	"micro/video2/handlers"
	"net/http"
	"os"
)

func main() {

	//this func creates a new logger
	//os.Stdout signifies where to write the log
	//The prefix appears at the beginning of each generated log line
	//flags as last argument, defining logging properties
	l := log.New(os.Stdout, "product-api", log.LstdFlags) //creating a new log

	//we got the instance of struct in hh variable
	hh := handlers.NewHello(l)
	gh := handlers.NewBye(l)

	//we create a new servvemux
	//here we're not going to use the defaultservemux
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	//we pass our custom mux here with handle
	http.ListenAndServe(":9090", sm)

}
