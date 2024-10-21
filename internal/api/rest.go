package rest

import (
	"context"
	"errors"
	v1 "irule-api/internal/api/v1"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(gtcx context.Context, wg *sync.WaitGroup, dbPool *pgxpool.Pool) {
	defer wg.Done()

	r := chi.NewRouter()
	r.Mount("/api/v1", v1.New(dbPool))

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-gtcx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Error closing server: %v", err)
	}

	log.Println("Server stopped")
}