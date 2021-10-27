package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Fudoshin2596/curly-telegram/backend/data"
	"github.com/Fudoshin2596/curly-telegram/backend/handlers"
	"github.com/go-openapi/runtime/middleware"
	gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	var wait time.Duration
	var bindAddress string

	flag.DurationVar(&wait, "graceful-timeout", 30*time.Second, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.StringVar(&bindAddress, "BIND_ADDRESS", ":9090", "Bind address for the server")
	flag.Parse()

	l := log.New(os.Stdout, "products-api-", log.LstdFlags)
	v := data.NewValidation()

	// create the handlers
	productHandler := handlers.NewProducts(l, v)

	// create a new serve Mux and register the handlers
	r := mux.NewRouter()
	getRouter := r.Methods(http.MethodGet).Subrouter()
	postRouter := r.Methods(http.MethodPost).Subrouter()
	putRouter := r.Methods(http.MethodPut).Subrouter()
	deleteRouter := r.Methods(http.MethodDelete).Subrouter()

	// CRUD
	getRouter.HandleFunc("/products", productHandler.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.ListSingle)

	postRouter.HandleFunc("/products", productHandler.Create)
	postRouter.Use(productHandler.MiddlewareValidateProduct)

	putRouter.HandleFunc("/products/", productHandler.Update)
	putRouter.Use(productHandler.MiddlewareValidateProduct)

	deleteRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	redoc := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", redoc)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// Apply the CORS middleware to our top-level router, with the defaults.
	gCors := gorilla.CORS(gorilla.AllowedOrigins([]string{"*"}))

	// create a new server
	srv := &http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      gCors(r),                 // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			l.Printf("Error starting server: %s\n", err)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// or SIGKILL (Ctrl+/)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// Pass a context with a timeout to tell a blocking function that it
	// should abandon its work after the timeout elapses.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	// Even though ctx will be expired, it is good practice to call its
	// cancellation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		return
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}
