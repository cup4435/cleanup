package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pb "cleanup/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const PORT = "8080"

type APIResponse struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}

func main() {
	// Get stock proxy service address
	stockProxyAddr := os.Getenv("STOCK_PROXY_ADDR")
	if stockProxyAddr == "" {
		stockProxyAddr = "localhost:50051"
	}

	// Connect to stock proxy service with retry logic
	var conn *grpc.ClientConn
	var err error

	for i := range 10 {
		conn, err = grpc.NewClient(
			stockProxyAddr,
			grpc.WithTransportCredentials(
				credentials.NewClientTLSFromCert(nil, "")))
		if err == nil {
			break
		}
		log.Printf("Failed to connect to stock proxy (attempt %d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to stock proxy after 10 attempts: %v", err)
	}
	defer conn.Close()

	// Create gRPC client
	client := pb.NewStockServiceClient(conn)

	// HTTP handlers
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	http.HandleFunc("/stock/", func(w http.ResponseWriter, r *http.Request) {
		// Extract ticker from URL path
		ticker := r.URL.Path[len("/stock/"):]
		if ticker == "" {
			http.Error(w, "Ticker symbol required", http.StatusBadRequest)
			return
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Call stock proxy service via gRPC
		req := &pb.StockRequest{
			Ticker: ticker,
		}

		resp, err := client.GetStock(ctx, req)
		if err != nil {
			log.Printf("gRPC call failed: %v", err)
			http.Error(w, fmt.Sprintf("Failed to get stock data: %v", err), http.StatusInternalServerError)
			return
		}

		// Convert gRPC response to API response
		apiResp := APIResponse{
			Ticker: resp.GetTicker(),
			Price:  resp.GetPrice(),
			Status: resp.GetStatus(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(apiResp)
	})

	log.Printf("API Gateway starting on port %s...", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
