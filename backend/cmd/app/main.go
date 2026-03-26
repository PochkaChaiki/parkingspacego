package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/handler"
	"github.com/pochkachaiki/parkingspace/internal/repository"
	"github.com/pochkachaiki/parkingspace/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// === MongoDB Connection ===
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDBName := os.Getenv("MONGO_DB")
	if mongoDBName == "" {
		mongoDBName = "parking"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Проверка подключения к MongoDB
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}
	log.Printf("Connected to MongoDB: %s/%s", mongoURI, mongoDBName)

	// === Initialize Layers ===
	db := client.Database(mongoDBName)
	coll := db.Collection("parking_records")

	// Repository layer
	repo := repository.NewMongoRepository(coll)
	log.Println("Repository initialized")

	// Service layer
	svc := service.NewService(repo)
	log.Println("Service initialized")

	// Handler layer

	h := handler.NewHandler(svc, log.New(os.Stdout, "Backend: ", log.LstdFlags))
	log.Println("Handler initialized")

	// === Setup HTTP Routes ===
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", h.Health)

	// POST /api/sessions - create session
	mux.HandleFunc("POST /api/sessions", h.StartSession)

	// GET /api/sessions/{phone} - get session info
	mux.HandleFunc("GET /api/sessions/{phone}", h.GetSession)

	// PATCH /api/records/{phone} - prolong session
	mux.HandleFunc("PATCH /api/sessions/{phone}", h.ProlongSession)

	// DELETE /api/records/{phone} - stop session
	mux.HandleFunc("DELETE /api/sessions/{phone}", h.StopSession)

	// === CORS Middleware ===
	corsHandler := corsMiddleware(mux)

	// === Start Server ===
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// corsMiddleware добавляет CORS headers для всех ответов
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
