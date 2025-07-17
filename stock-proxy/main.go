package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	pb "cleanup/proto"

	"google.golang.org/grpc"
)

// StockServer implements the StockService gRPC service
type StockServer struct {
	pb.UnimplementedStockServiceServer
}

// GetStock implements the GetStock RPC method
func (s *StockServer) GetStock(ctx context.Context, req *pb.StockRequest) (*pb.StockResponse, error) {
	ticker := req.GetTicker()
	log.Printf("Received gRPC request for ticker: %s", ticker)

	// Validate ticker (basic validation)
	if ticker == "" {
		return &pb.StockResponse{
			Ticker: ticker,
			Price:  0,
			Status: "error: empty ticker symbol",
		}, nil
	}

	// TODO: This is where you would implement actual stock data fetching
	// For now, we'll simulate different responses based on ticker

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	var price float64
	var status string

	// Mock different responses for different tickers
	switch strings.ToUpper(ticker) {
	case "AAPL":
		price = 150.25 + rand.Float64()*10 // Random price around 150
		status = "success"
	case "GOOGL":
		price = 2800.50 + rand.Float64()*50 // Random price around 2800
		status = "success"
	case "TSLA":
		price = 200.75 + rand.Float64()*20 // Random price around 200
		status = "success"
	case "MSFT":
		price = 300.15 + rand.Float64()*15 // Random price around 300
		status = "success"
	case "ERROR-TESTING":
		// Simulate an error case for testing
		ticker = ""
		price = 0
		status = "error: ticker not found"
	default:
		// Default mock price for unknown tickers
		price = 50.0 + rand.Float64()*100 // Random price between 50-150
		status = "success"
	}

	response := &pb.StockResponse{
		Ticker: ticker,
		Price:  price,
		Status: status,
	}

	log.Printf("Returning gRPC response: ticker=%s, price=%.2f, status=%s",
		response.Ticker, response.Price, response.Status)

	return response, nil
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	port := os.Getenv("PROXY_PORT")
	if port == "" {
		port = "50051"
	}

	// Create TCP listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Create gRPC server
	s := grpc.NewServer()

	// Register our service implementation
	stockServer := &StockServer{}
	pb.RegisterStockServiceServer(s, stockServer)

	log.Printf("Stock Proxy gRPC Server starting on port %s...", port)
	log.Printf("Server registered and ready to accept connections")

	// Start serving
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
