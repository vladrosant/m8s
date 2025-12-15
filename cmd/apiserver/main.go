package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vladrosant/m8s/pkg/api"
	"github.com/vladrosant/m8s/pkg/store"
)

func main() {
	fmt.Println("starting m8s API server...")

	stateFile := "/var/lib/m8s/state.json"
	st, err := store.NewStore(stateFile)
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}
	fmt.Println("state store initialized: %s\n", stateFile)

	server := api.NewServer(st)

	http.HandleFunc("/api/v1/pods", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			server.HandleCreatePod(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("name") != "" {
				server.HandleGetPod(w, r)
			} else {
				server.HandleListPods(w, r)
			}
		case http.MethodDelete:
			server.HandleDeletePod(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("m8s API server - your container otchestrator matey!\n"))
	})

	port := ":8080"
	go func() {
		fmt.Printf("API server listening on http://localhost%s\n", port)
		fmt.Println("\nAvailable endpoints:")
		fmt.Println("	GET	/health")
		fmt.Println("	POST	/api/v1/pods")
		fmt.Println("	GET	/api/v1/pods")
		fmt.Println("	GET	/api/v1/pods?name=<name>&namespace=<namespace>")
		fmt.Println("	DELETE	/api/v1/pods?name=<name>&namespace=<namespace>")
		fmt.Println("\npress ctrl+c to stop")

		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nshutting down API server...")
}
