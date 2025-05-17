package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KirillKotovsky/location_user_api/internal/handler"
	// "github.com/KirillKotovsky/location_user_api/internal/model"

	grpcclient "github.com/KirillKotovsky/location_proto/proto/pkg/grpcclient"
	"github.com/KirillKotovsky/location_user_api/internal/repository"
	"github.com/KirillKotovsky/location_user_api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/location?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	repo := repository.New(db)
	if err := repo.InitSchema(ctx); err != nil {
		log.Fatalf("failed to init schema: %v", err)
	}

	grpcAddr := os.Getenv("GRPC_LOCATION_TRACK_ADDR")
	if grpcAddr == "" {
		grpcAddr = "localhost:50051"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC: %v", err)
	}
	defer conn.Close()

	grpcClient := grpcclient.New(conn)
	svc := service.New(repo, grpcClient)
	h := handler.New(svc)

	r := gin.Default()
	r.POST("/location", h.HandleLocationUpdate)
	r.GET("/nearby", h.HandleNearbyUsers)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
