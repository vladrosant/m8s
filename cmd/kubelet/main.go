package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vladrosant/m8s/pkg/kubelet"
)

func main() {
	nodeName := flag.String("node-name", "node-01", "Name of this node")
	apiURL := flag.String("api-url", "http://localhost:8080", "URL of the API server")
	syncInterval := flag.Int("sync-interval", 5, "Sync interval in seconds")
	flag.Parse()

	fmt.Printf(" Starting m8s Kubelet...\n")
	fmt.Printf("	Node Name: %s\n", *nodeName)
	fmt.Printf("	API Server: %s\n", *apiURL)
	fmt.Printf("	Sync Interval: %ds\n", *syncInterval)

	pm := kubelet.NewPodManager(*nodeName, *apiURL)

	ticker := time.NewTicker(time.Duration(*syncInterval) * time.Second)
	defer ticker.Stop()

	fmt.Println("\n Kubelet started! Watching for pods...")
	fmt.Println("Press ctrl+c to stop\n")

	if err := pm.Sync(); err != nil {
		log.Printf("Sync error: %v:", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			if err := pm.Sync(); err != nil {
				log.Printf("Sync error: %v", err)
			}
		case <-sigChan:
			fmt.Println("\nShutting down kubelet...")
			return
		}
	}
}
