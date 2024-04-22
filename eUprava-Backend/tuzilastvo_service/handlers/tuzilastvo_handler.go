package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"tuzilastvo_service/data"
)

type TuzilastvoHandler struct {
	logger         *log.Logger
	tuzilastvoRepo *data.TuzilastvoRepo
}

var (
	tuzilastvoServiceHost = os.Getenv("TUZILASTVO_SERVICE_HOST")
	tuzilastvoServicePort = os.Getenv("TUZILASTVO_SERVICE_PORT")
)

func NewTuzilastvoHandler(l *log.Logger, r *data.TuzilastvoRepo) *TuzilastvoHandler {
	return &TuzilastvoHandler{l, r}
}

func (s *TuzilastvoHandler) MiddlewareDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		appointments := &data.Appointment{}
		err := appointments.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			s.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyProduct{}, appointments)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}
