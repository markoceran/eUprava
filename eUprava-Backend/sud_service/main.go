package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sud_service/client"
	"sud_service/data"
	"sud_service/domain"
	"sud_service/handlers"
	"sud_service/middlewares"
	"time"
)

var (
	JaegerAddress = os.Getenv("JAEGER_ADDRESS")
)

func main() {
	port := os.Getenv("SUD_SERVICE_PORT")

	timeoutContext, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	exp, err := newExporter(JaegerAddress)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)
	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(timeoutContext) }()
	otel.SetTracerProvider(tp)
	// Finally, set the tracer that can be used for this package.
	tracer := tp.Tracer("sud_service")
	otel.SetTextMapPropagator(propagation.TraceContext{})

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[sud-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[sud-store] ", log.LstdFlags)

	// NoSQL: Initialize Repository store
	store, err := data.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.DisconnectMongo(timeoutContext)
	store.Ping()

	tuzilastvoClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}
	tuzilastvoBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "tuzilastvo",
			MaxRequests: 1,
			Timeout:     10 * time.Second,
			Interval:    0,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures > 2
			},
			OnStateChange: func(name string, from, to gobreaker.State) {
				logger.Printf("CB '%s' changed from '%s' to '%s'\n", name, from, to)
			},
			IsSuccessful: func(err error) bool {
				if err == nil {
					return true
				}
				errResp, ok := err.(domain.ErrResp)
				return ok && errResp.StatusCode >= 400 && errResp.StatusCode < 500
			},
		},
	)

	tuzilastvUri := fmt.Sprintf("http://%s:%s", os.Getenv("TUZILASTVO_SERVICE_HOST"), os.Getenv("TUZILASTVO_SERVICE_PORT"))
	tuzilastvo := client.NewTuzilastvoClient(tuzilastvoClient, tuzilastvUri, tuzilastvoBreaker)

	sudHandler := handlers.NewSudHandler(logger, store, tracer, tuzilastvo)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	router.Use(middlewares.MiddlewareContentTypeSet)

	casbinMiddleware, err := middlewares.InitializeCasbinMiddleware("./rbac_model.conf", "./policy.csv")
	if err != nil {
		log.Fatal(err)
	}
	router.Use(casbinMiddleware)

	dobaviPredmete := router.Methods(http.MethodGet).Subrouter()
	dobaviPredmete.HandleFunc("/predmeti", sudHandler.DobaviPredmete)

	kreirajPredmet := router.Methods(http.MethodPost).Subrouter()
	kreirajPredmet.HandleFunc("/predmeti", sudHandler.DodajPredmet)
	kreirajPredmet.Use(sudHandler.MiddlewareDeserialization)

	dobaviPredmetPoId := router.Methods(http.MethodGet).Subrouter()
	dobaviPredmetPoId.HandleFunc("/predmeti/{id}", sudHandler.DobaviPredmetPoId)

	dodajPredmetePoZahtjevima := router.Methods(http.MethodPost).Subrouter()
	dodajPredmetePoZahtjevima.HandleFunc("/predmeti/zahtjevi", sudHandler.DodajPredmetePoZahtjevima)
	dodajPredmetePoZahtjevima.Use(sudHandler.MiddlewareDeserialization)

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

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources are set.
	r := resource.Default()

	// Set the service name.
	serviceName := "sud_service"

	// Merge additional attributes.
	mergedResource, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		panic(err)
	}

	// Merge default and additional resources.
	r, err = resource.Merge(r, mergedResource)
	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func newExporter(address string) (*jaeger.Exporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(address)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}
