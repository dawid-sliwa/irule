package main

import (
	"context"
	rest "irule-api/internal/api"
	"irule-api/internal/config"
	"irule-api/internal/db"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bugsnag/panicwrap"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		zap.S().Errorw("failed to load config",

			"error", err,
		)
		os.Exit(1)
	}

	exitStatus, err := panicwrap.BasicWrap(func(s string) {
		zap.S().Errorw("panic detected",
			"panic", s,
		)
	})

	if err != nil {
		zap.S().Errorw("failed to setup panic handler",
			"error", err,
		)
		os.Exit(2)
	}

	if exitStatus >= 0 {
		os.Exit(exitStatus)
	}

	gctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	pool, err := db.NewPgPool(cfg)
	if err != nil {
		zap.S().Errorw("failed to create db pool",

			"error", err,
		)
		os.Exit(2)
	}

	done := make(chan struct{})
	wg := sync.WaitGroup{}

	go func() {
		<-quit
		cancel()

		go func() {
			select {
			case <-time.After(time.Second * 10):
			case <-quit:
			}
			zap.S().Fatal("forced shutdown")
		}()

		zap.S().Info("waiting for shutdown")

		wg.Wait()

		close(done)
	}()

	wg.Add(1)
	go rest.New(gctx, &wg, pool, cfg)

	zap.S().Info("Server started")

	<-done

	zap.S().Info("Server stopped")

	os.Exit(0)
}
