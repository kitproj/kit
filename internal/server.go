package internal

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
)

//go:embed index.html
var indexHTML string

func StartServer(ctx context.Context, dag DAG[*TaskNode], events chan any) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-events:
				switch event.(type) {
				case *TaskNode:
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
			}
		}
	})

	server := &http.Server{
		// only allow local connections
		Addr:    "127.0.0.1:3000",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
