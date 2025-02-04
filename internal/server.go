package internal

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

//go:embed index.html
var indexHTML string

func StartServer(ctx context.Context, port int, wg *sync.WaitGroup, dag DAG[*TaskNode], events chan *TaskNode) {

	streams := &sync.Map{}

	go func() {
		for event := range events {
			streams.Range(func(key, value any) bool {
				value.(chan *TaskNode) <- event
				return true
			})
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// if internal/index.html exists, serve that
		_, err := os.Stat("internal/index.html")
		if err == nil {
			http.ServeFile(w, r, "internal/index.html")
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, err = w.Write([]byte(indexHTML))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/dag", func(w http.ResponseWriter, r *http.Request) {
		// return the dag
		marshal, err := json.Marshal(dag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(marshal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {

		id := rand.Int()

		// create a stream for this connection
		stream := make(chan *TaskNode, 100)

		// load the stream with the current state
		for _, node := range dag.Nodes {
			stream <- node
		}
		streams.Store(id, stream)
		defer func() {
			streams.Delete(id)
		}()

		// return an event stream
		w.Header().Set("Content-Type", "text/event-stream")
		for event := range stream {
			marshal, err := json.Marshal(event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = fmt.Fprintf(w, "data: %s\n\n", marshal)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.(http.Flusher).Flush()
		}
	})
	mux.HandleFunc("/logs/{task}", func(w http.ResponseWriter, r *http.Request) {
		//ctx := r.Context()
		task := r.PathValue("task")
		node, ok := dag.Nodes[task]
		if !ok {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		file, err := os.Open(node.logFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "text/event-stream")

		for {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				_, err := fmt.Fprintf(w, "data: %s\n\n", line)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.(http.Flusher).Flush()
			}

			if err := scanner.Err(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Sleep for a short duration before checking for new lines
			time.Sleep(1 * time.Second)

			// Reset the scanner to continue reading new lines
			_, err := file.Seek(0, io.SeekCurrent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	server := &http.Server{
		// only allow local connections
		Addr:    fmt.Sprintf("localhost:%d", port),
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

	log.Printf("UI available on http://%s", server.Addr)

	wg.Add(1)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
