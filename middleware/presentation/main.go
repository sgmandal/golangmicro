package main

import (
	"context"
	"hello/presentation/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	getRouter := sm.Methods(http.MethodGet).Subrouter() // request
	getRouter.HandleFunc("/", ph.GetRaw)

	post := sm.Methods(http.MethodPost).Subrouter()
	post.HandleFunc("/", ph.AddRaw)
	post.Use(ph.MiddlewareProductValidation)

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
}
