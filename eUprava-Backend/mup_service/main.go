package main

import (
	"context"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"log"
	"mup_service/data"
	"mup_service/handlers"
	"mup_service/middlewares"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	JaegerAddress = os.Getenv("JAEGER_ADDRESS")
)

func main() {

	port := os.Getenv("MUP_SERVICE_PORT")

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
	tracer := tp.Tracer("mup_service")
	otel.SetTextMapPropagator(propagation.TraceContext{})

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

	mupHandler := handlers.NewMupHandler(logger, store, tracer)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	router.Use(middlewares.MiddlewareContentTypeSet)

	casbinMiddleware, err := middlewares.InitializeCasbinMiddleware("./rbac_model.conf", "./policy.csv")
	if err != nil {
		log.Fatal(err)
	}
	router.Use(casbinMiddleware)

	kreirajLicnuKartu := router.Methods(http.MethodPut).Subrouter()
	kreirajLicnuKartu.HandleFunc("/kreirajLicnuKartu/{id}", mupHandler.KreirajLicnuKartu)

	dobaviKorisnike := router.Methods(http.MethodGet).Subrouter()
	dobaviKorisnike.HandleFunc("/dobaviKorisnike", mupHandler.DobaviKorisnike)

	kreirajVozackuDozvolu := router.Methods(http.MethodPut).Subrouter()
	kreirajVozackuDozvolu.HandleFunc("/kreirajVozackuDozvolu/{id}", mupHandler.KreirajVozackuDozvolu)

	kreirajSaobracajnuDozvolu := router.Methods(http.MethodPut).Subrouter()
	kreirajSaobracajnuDozvolu.HandleFunc("/kreirajSaobracajnuDozvolu/{id}", mupHandler.KreirajSaobracajnuDozvolu)

	kreirajNalogZaPracenje := router.Methods(http.MethodPost).Subrouter()
	kreirajNalogZaPracenje.HandleFunc("/kreirajNalogZaPracenje", mupHandler.KreirajNalogZaPracenje)

	dobaviNalogeZaPracenje := router.Methods(http.MethodGet).Subrouter()
	dobaviNalogeZaPracenje.HandleFunc("/dobaviNalogeZaPracenje", mupHandler.DobaviNalogeZaPracenje)

	kreirajPasos := router.Methods(http.MethodPut).Subrouter()
	kreirajPasos.HandleFunc("/kreirajPasos/{id}", mupHandler.KreirajPasos)

	validirajDokumente := router.Methods(http.MethodPost).Subrouter()
	validirajDokumente.HandleFunc("/validirajDokumente", mupHandler.ValidirajDokumente)

	dobaviJmbgKorisnika := router.Methods(http.MethodGet).Subrouter()
	dobaviJmbgKorisnika.HandleFunc("/dobaviJmbgKorisnika/{id}", mupHandler.DobaviJmbgKorisnika)

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
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("reservations_service"),
		),
	)

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
