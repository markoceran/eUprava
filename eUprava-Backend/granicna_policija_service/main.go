package granicna_policija_service

import (
	"context"
	"github.com/gorilla/mux"
	"granicna_policija_service/data"
	"granicna_policija_service/middlewares"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	port := os.Getenv("GRANICNA_POLICIJA_SERVICE_PORT")

	timeoutContext, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[granicna_policija-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[granicna_policija-store] ", log.LstdFlags)

	// NoSQL: Initialize Repository store
	store, err := data.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.DisconnectMongo(timeoutContext)
	store.Ping()

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	router.Use(middlewares.MiddlewareContentTypeSet)
	router.Use(middlewares.TokenValidationMiddleware)

	casbinMiddleware, err := middlewares.InitializeCasbinMiddleware("./rbac_model.conf", "./policy.csv")
	if err != nil {
		log.Fatal(err)
	}
	router.Use(casbinMiddleware)

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Println("Server listening on port", port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")

}
