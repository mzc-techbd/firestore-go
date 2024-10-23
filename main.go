package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	firestoreReads = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "firestore_reads_total",
			Help: "Total number of reads from Firestore",
		},
		[]string{"collection", "status"},
	)
	firestoreReadLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "firestore_read_latency_seconds",
			Help:    "Latency of Firestore read operations",
			Buckets: prometheus.DefBuckets, // Use default buckets
		},
		[]string{"collection", "status"},
	)
)

func init() {
	prometheus.MustRegister(firestoreReads)
	prometheus.MustRegister(firestoreReadLatency)
}

func main() {
	projectId := "kktae-demo"
	collectionPath := "test-collection"
	docId := "test-document"

	// Create a Firestore client.
	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer firestoreClient.Close()

	// Cloud Run health check endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Firestore Reader is running!")
	})

	// Endpoint to read from Firestore and return the result
	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {

		// Start timer for latency measurement
		startTime := time.Now()

		status := "success"
		isFailed := false

		// Perform a simple read operation.
		doc, err := firestoreClient.Collection(collectionPath).Doc(docId).Get(ctx)
		if err != nil {
			log.Printf("Failed to read document: %v", err)
			http.Error(w, "Failed to read document", http.StatusInternalServerError)

			isFailed = true
			status = "error"
		}

		// Calculate latency and record the metric
		latency := time.Since(startTime).Seconds()

		// Increment read counter and observe latency with status label
		firestoreReads.WithLabelValues(collectionPath, status).Inc()
		firestoreReadLatency.WithLabelValues(collectionPath, status).Observe(latency)

		if !isFailed {
			// Return the document data as a JSON response.
			data := doc.Data()
			fmt.Fprintf(w, "%v", data)
		}
	})

	// Expose Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
