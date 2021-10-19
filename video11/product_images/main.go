package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"product_img/product_images/files"
	"product_img/product_images/handlers"
	"time"

	gohandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9091", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "./imagestore", "Base path to save images")

func main() {
	env.Parse()

	//hclogger main initialization is done here
	// and then that data is used in handlers/files.go
	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product_images",
			Level: hclog.LevelFromString(*logLevel), //returned for nicholasjackson/env
		},
	)

	// create a logger for the server from the default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// create the storage class, use local storage
	// max filesize 5MB
	stor, err := files.NewLocal(*basePath, 1024*1000*5)
	if err != nil {
		l.Error("unable to create storage", "error", err)
		os.Exit(1)
	}

	// create the handlers
	fh := handlers.NewFiles(stor, l) //logger passed

	// create a new serve mux and register the handlers
	sm := mux.NewRouter() // done in video

	ch := gohandler.CORS(gohandler.AllowedOrigins([]string{"*"})) //CORS bypass, used for security reasons to connect to frontend

	// filename regex: {filename:[a-zA-Z]+\\.[a-z]{3}}
	// problem with FileServer is it lacks some features
	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.RetardFunction)
	ph.HandleFunc("/", fh.UploadMultipart)

	// get files
	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))),
	)

	/*as we've implemented the ServeHTTP interface, the package http calls such interface
	and do all the required operations, hence interfaces are useful(serveHTTP is used internally in http package which is used further)*/
	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      ch(sm),            // set the default handler, wraps main router
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request form the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max timeout for connections using TCP Keep-Alive
	}

	//start the server
	go func() {
		l.Info("starting the server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("unable to start the server", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1) // buffered channel of type os.Signal
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received
	sig := <-c

	l.Info("shutting the server", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
