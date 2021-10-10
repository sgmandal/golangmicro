/*didnt clearly understand the concept of middleware*/

package main

import (
	"context"
	"log"
	world "micro/video7/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
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

	//we create a new servvemux
	//here we're not going to use the defaultservemux

	//create a new serve mux using gorilla mux
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter() //creating a get subrouter
	getRouter.HandleFunc("/", ph.GetProducts)

	//creating a put subrouter
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation) //this will get executed before the actual handler function

	//creating post subrouter
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

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
