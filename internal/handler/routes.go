package handler

import (
	"context"
	"net/http"

	"booking-system/internal/repository"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, repo *repository.BookingRepository) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/bookings", func(c *gin.Context) {

		data, err := repo.GetAll(context.TODO())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, data)
	})
}