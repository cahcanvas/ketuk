package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ketuk.id/api/internal/app"
	"ketuk.id/api/internal/sweepjob"
)

func main() {
	ctx := context.Background()
	a, err := app.Build(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Pool.Close()

	bgCtx, stopBackground := context.WithCancel(ctx)
	defer stopBackground()
	go sweepjob.Loop(bgCtx, 15*time.Minute, a.Identity, a.Billing, a.Gift)

	srv := &http.Server{
		Addr:              a.Config.HTTPAddr,
		Handler:           a.Handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		log.Printf("ketuk-api listening on %s", a.Config.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	stopBackground()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
