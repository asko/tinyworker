package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const HttpTimeout = time.Duration(30) * time.Second

type Config struct {
	StartupDelay       time.Duration `env:"STARTUP_DELAY" envDefault:"10s"`
	LifetimeSecondsMax int           `env:"LIFETIME_SECONDS_MAX" envDefault:"600"`
	HttpPort           int           `env:"HTTP_PORT" envDefault:"8080"`
}

func fatalOn(err error) {
	if err != nil {
		fmt.Errorf("unexpected error: %s", err)
	}
}

func main() {
	c := Config{}
	err := env.Parse(&c)
	fatalOn(err)
	fatalOn(err)

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	r := rand.Intn(c.LifetimeSecondsMax)
	shutdownDelay := time.Duration(r) * time.Second
	fmt.Printf("Initializing. Startup delay %s, http lifetime %s\n", c.StartupDelay, shutdownDelay)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	srv := &http.Server{Handler: mux, Addr: fmt.Sprintf(":%d", c.HttpPort)}

	select {
	case <-stopChan:
		fmt.Println("Initial delay interrupted")
		close(stopChan)
	case <-time.After(c.StartupDelay):
		go func() {
			fmt.Println("Starting to listen for HTTP calls")
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fatalOn(err)
			}
		}()
	}

	select {
	case <-stopChan:
		fmt.Println("Shutdown delay interrupted")
	case <-time.After(shutdownDelay):
		fmt.Println("Shutting down")
	}

	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
	defer cancel()
	fatalOn(srv.Shutdown(ctx))

	fmt.Println("Stopped")
}

func handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello"))
}
