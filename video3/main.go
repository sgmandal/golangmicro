package main

import (
	"context"
	"log"
	hello "micro/video2/handlers"
	world "micro/video3/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	//this func creates a new logger
	//os.Stdout signifies where to write the log
	//The prefix appears at the beginning of each generated log line
	//flags as last argument, defining logging properties
	l := log.New(os.Stdout, "product-api", log.LstdFlags) //creating a new log

	//we got the instance of struct in hh variable
	//hh := handlers.NewHello(l)
	ph := world.NewProducts(l)
	gh := hello.NewBye(l)

	//we create a new servvemux
	//here we're not going to use the defaultservemux
	sm := http.NewServeMux()
	sm.Handle("/", ph)
	sm.Handle("/goodbye", gh)

	//we pass our custom mux here with handle
	//a more feature packed listen and serve providedby http package
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	//yes the program gets terminated here
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	//graceful shutdown
	//used to close database connections so no ghost connections are present
	//when go program is terminated
	//check documentation os.Notify

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("we received termination with following reason ", sig) //keyboard interrupt is printed

	//check context package documentation
	//no leaking goroutines
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

	//http.ListenAndServe(":9090", sm) this works but we cant have 2 listen and serve so I commented it out

}
