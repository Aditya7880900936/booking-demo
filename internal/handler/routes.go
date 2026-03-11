package handler

import (
	"context"
	"net/http"

	"booking-system/internal/pipeline"
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

	// ⭐ create pipeline service
	svc := &pipeline.Service{
		Repo: repo,
	}

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

	// ⭐ CLEAN INGEST ENDPOINT
	r.POST("/ingest/email", func(c *gin.Context) {

		var req EmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		res, err := svc.ProcessEmail(
			context.TODO(),
			req.UserID,
			req.Sender,
			req.Subject,
			req.Body,
		)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, res)
	})
}