package handler

import (
	"context"
	"net/http"
	"time"

	"booking-system/internal/detector"
	"booking-system/internal/repository"
	"booking-system/models"

	"github.com/gin-gonic/gin"
)

type EmailRequest struct {
	UserID  string `json:"user_id"`
	Sender  string `json:"sender"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}


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

	r.POST("/ingest/email", func(c *gin.Context) {

		var req EmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	
		result := detector.DetectEmail(req.Subject, req.Sender, req.Body)
	
		if !result.IsBooking {
			c.JSON(200, gin.H{
				"message":    "not a booking email",
				"confidence": result.Confidence,
				"signals":    result.Signals,
			})
			return
		}
	
		// Temporary placeholder booking (real extraction next step)
		b := models.Booking{
			UserID:     req.UserID,
			Type:       "flight",
			Status:     "detected",
			Confidence: result.Confidence,
			CreatedAt:  time.Now(),
		}
	
		repo.Save(context.TODO(), &b)
	
		c.JSON(200, gin.H{
			"message":    "booking detected",
			"confidence": result.Confidence,
			"signals":    result.Signals,
		})
	})
}
