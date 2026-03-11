package main

import (
	"context"
	"log"
	"time"

	"booking-system/internal/handler"
	"booking-system/internal/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("bookingdb").Collection("bookings")

	repo := &repository.BookingRepository{
		Collection: collection,
	}

	router := gin.Default()

	handler.RegisterRoutes(router, repo)

	log.Println("Server running on :8080")
	router.Run(":8080")
}