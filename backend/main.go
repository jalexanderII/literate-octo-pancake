package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gorilla "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/literate-octo-pancake/backend/data"
	"github.com/jalexanderII/literate-octo-pancake/backend/handlers"
	"github.com/jalexanderII/literate-octo-pancake/currency/protos/currency"
	"google.golang.org/grpc"
)

const (
	server     = "localhost"
	serverPort = "9092"
)

var (
	wait        time.Duration
	bindAddress string
)

func main() {
	flag.DurationVar(&wait, "graceful-timeout", 30*time.Second, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.StringVar(&bindAddress, "BIND_ADDRESS", ":9090", "Bind address for the server")
	flag.Parse()

	l := hclog.Default()
	v := data.NewValidation()

	serverAddr := net.JoinHostPort(server, serverPort)

	// setup insecure connection
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// grpc client
	curClient := currency.NewCurrencyClient(conn)

	// create productsDB
	pdb := data.NewProductsDB(l, curClient)

	// create the handlers
	productHandler := handlers.NewProducts(l, v, pdb)

	// create a new serve Mux and register the handlers
	r := mux.NewRouter()
	getRouter := r.Methods(http.MethodGet).Subrouter()
	postRouter := r.Methods(http.MethodPost).Subrouter()
	putRouter := r.Methods(http.MethodPut).Subrouter()
	deleteRouter := r.Methods(http.MethodDelete).Subrouter()

	// CRUD
	getRouter.HandleFunc("/products", productHandler.ListAll).Queries("currency", "{[A-Z](3)}")
	getRouter.HandleFunc("/products", productHandler.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.ListSingle).Queries("currency", "{[A-Z](3)}")
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
		Addr:         bindAddress,                                      // configure the bind address
		Handler:      gCors(r),                                         // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		l.Info("starting backend service on", "bindAddress", bindAddress)
		if err := srv.ListenAndServe(); err != nil {
			l.Error("Error starting server", "error", err)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// or SIGTERM (Ctrl+/)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Info("Got signal", "signal", sig)

	// Pass a context with a timeout to tell a blocking function that it
	// should abandon its work after the timeout elapses.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	// Even though ctx will be expired, it is good practice to call its
	// cancellation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err = srv.Shutdown(ctx)
	if err != nil {
		return
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	l.Info("shutting down")
	os.Exit(0)

}
