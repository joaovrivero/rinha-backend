package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"runtime/trace"
	"syscall"
	"time"

	"github.com/joaovrivero/rinha-backend/internal/cache"
	"github.com/joaovrivero/rinha-backend/internal/config"
	"github.com/joaovrivero/rinha-backend/internal/database"
	route "github.com/joaovrivero/rinha-backend/internal/http"
	"github.com/joaovrivero/rinha-backend/internal/pessoa"
)

func main() {
	if config.PROFILING {
		slog.Info("Running with profiling")
		// create a cpu profile
		f, err := os.Create("/pprof/profile.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// start cpu profiling
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()

		// start tracing
		tracefile, err := os.Create("/pprof/trace.out")
		if err != nil {
			panic(err)
		}
		defer tracefile.Close()
		if err := trace.Start(tracefile); err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	slog.Info("Waiting for database connection")
	time.Sleep(5 * time.Second) // wait for database on docker-compose
	if err := database.Connect(); err != nil {
		panic(err)
	}
	defer database.Close()

	if err := cache.Connect(); err != nil {
		panic(err)
	}

	chExit := make(chan struct{})
	repo := pessoa.NewRepository(database.Connection, cache.Client)

	for i := 0; i < config.NumWorkers; i++ {
		go pessoa.RunWorker(repo.ChPessoas, chExit, repo, config.NumBatch)
	}

	http.HandleFunc("/", route.Pessoas)

	slog.Info("Server running on port 80")
	go http.ListenAndServe(":80", nil)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	close(repo.ChPessoas)
	for i := 0; i < config.NumWorkers; i++ {
		<-chExit
	}
}
