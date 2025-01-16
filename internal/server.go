package internal

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

//go:embed index.html
var indexHTML string

func StartServer(ctx context.Context, wg *sync.WaitGroup, dag DAG[*TaskNode], events chan *TaskNode) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// if internal/index.html exists, serve that
		_, err := os.Stat("internal/index.html")
		if err == nil {
			http.ServeFile(w, r, "internal/index.html")
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(indexHTML))
	})
	mux.HandleFunc("/dag", func(w http.ResponseWriter, r *http.Request) {
		// return the dag
		marshal, err := json.Marshal(dag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(marshal)
	})
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		// return an event stream
		w.Header().Set("Content-Type", "text/event-stream")
		for event := range events {
			marshal, err := json.Marshal(event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("data: "))
			w.Write(marshal)
			w.Write([]byte("\n\n"))
			w.(http.Flusher).Flush()
		}
	})

	server := &http.Server{
		// only allow local connections
		Addr:    "127.0.0.1:3000",
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	log.Println("starting Kit server on http://localhost:3000")

	wg.Add(1)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
