package handler

import (
	"context"
	"net/http"

	"booking-system/internal/detector"
	"booking-system/internal/extractor"
	"booking-system/internal/repository"

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

	// ⭐ MAIN INGEST PIPELINE
	r.POST("/ingest/email", func(c *gin.Context) {

		var req EmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// ⭐ Step 1 — Detection
		detectRes := detector.DetectEmail(req.Subject, req.Sender, req.Body)

		if !detectRes.IsBooking {
			c.JSON(200, gin.H{
				"stage":      "detection",
				"message":    "not a booking email",
				"confidence": detectRes.Confidence,
				"signals":    detectRes.Signals,
			})
			return
		}

		// ⭐ Step 2 — Extraction
		extractRes := extractor.ExtractFlight(req.Body, req.UserID)

		if extractRes.Confidence < 0.4 {
			c.JSON(200, gin.H{
				"stage":      "extraction",
				"message":    "booking detected but extraction weak",
				"confidence": extractRes.Confidence,
				"fields":     extractRes.FieldsFound,
			})
			return
		}

		booking := extractRes.Booking

		// ⭐ Step 3 — Duplicate check
		exists, err := repo.Exists(context.TODO(), booking.UserID, booking.BookingRef)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if exists {
			c.JSON(200, gin.H{
				"stage":   "storage",
				"message": "duplicate booking",
			})
			return
		}

		// ⭐ Step 4 — Save
		err = repo.Save(context.TODO(), booking)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"stage":      "completed",
			"message":    "booking stored",
			"confidence": extractRes.Confidence,
			"booking":    booking,
		})
	})
}