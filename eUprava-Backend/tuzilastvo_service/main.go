package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"tuzilastvo_service/data"
	"tuzilastvo_service/handlers"
	"tuzilastvo_service/middlewares"
)

func main() {

	port := os.Getenv("TUZILASTVO_SERVICE_PORT")

	timeoutContext, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[acc-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[acc-store] ", log.LstdFlags)

	// NoSQL: Initialize Repository store
	store, err := data.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.DisconnectMongo(timeoutContext)
	store.Ping()

	tuzilastvoHandler := handlers.NewTuzilastvoHandler(logger, store)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	router.Use(middlewares.MiddlewareContentTypeSet)
	router.Use(middlewares.TokenValidationMiddleware)

	casbinMiddleware, err := middlewares.InitializeCasbinMiddleware("./rbac_model.conf", "./policy.csv")
	if err != nil {
		log.Fatal(err)
	}
	router.Use(casbinMiddleware)

	//getAppointmentByAccommodation := router.Methods(http.MethodGet).Subrouter()
	//getAppointmentByAccommodation.HandleFunc("/appointmentsByAccommodation/{id}", appointmentHandler.GetAppointmentsByAccommodation)
	////getAllAppointment.Use(appointmentHandler.MiddlewareAppointmentDeserialization)
	//
	//getAppointmentsByDate := router.Methods(http.MethodGet).Subrouter()
	//getAppointmentsByDate.HandleFunc("/appointmentsByDate/", appointmentHandler.GetAppointmentsByDate)
	//
	//createAppointment := router.Methods(http.MethodPost).Subrouter()
	//createAppointment.HandleFunc("/appointments", appointmentHandler.CreateAppointment)
	//createAppointment.Use(appointmentHandler.MiddlewareAppointmentDeserialization)
	//
	//getAllAppointment := router.Methods(http.MethodGet).Subrouter()
	//getAllAppointment.HandleFunc("/appointments", appointmentHandler.GetAllAppointment)
	////getAllAppointment.Use(appointmentHandler.MiddlewareAppointmentDeserialization)
	//
	//createReservation := router.Methods(http.MethodPost).Subrouter()
	//createReservation.HandleFunc("/reservations", reservationHandler.CreateReservation)
	//createReservation.Use(reservationHandler.MiddlewareReservationDeserialization)

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
