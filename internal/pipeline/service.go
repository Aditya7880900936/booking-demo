package pipeline

import (
	"context"
	"log"

	"booking-system/internal/detector"
	"booking-system/internal/extractor"
	"booking-system/internal/repository"
	"booking-system/models"
	"booking-system/pkg/utils"
)

type Service struct {
	Repo *repository.BookingRepository
}

type PipelineResult struct {
	Success    bool            `json:"success"`
	Stage      string          `json:"stage"`
	Message    string          `json:"message"`
	Confidence float64         `json:"confidence,omitempty"`
	Booking    *models.Booking `json:"booking,omitempty"`
	Signals    []string        `json:"signals,omitempty"`
	Fields     []string        `json:"fields,omitempty"`
}

func (s *Service) ProcessEmail(
	ctx context.Context,
	userID string,
	sender string,
	subject string,
	body string,
) (*PipelineResult, error) {

	log.Println("📩 Processing new email ingestion")

	// ⭐ Step 1 — Detection
	log.Println("🔎 Running detection engine")

	detectRes := detector.DetectEmail(subject, sender, body)

	log.Printf("Detection confidence: %.2f Signals: %v\n",
		detectRes.Confidence, detectRes.Signals)

	if !detectRes.IsBooking || detectRes.Confidence < utils.DetectionThreshold {
		return &PipelineResult{
			Success:    false,
			Stage:      "detection",
			Message:    "not a booking email",
			Confidence: detectRes.Confidence,
			Signals:    detectRes.Signals,
		}, nil
	}

	// ⭐ Step 2 — Extraction
	log.Println("🧠 Running extraction engine")

	extractRes := extractor.ExtractFlight(body, userID)

	log.Printf("Extraction confidence: %.2f Fields: %v\n",
		extractRes.Confidence, extractRes.FieldsFound)

	if extractRes.Confidence < utils.ExtractionThreshold {
		return &PipelineResult{
			Success:    false,
			Stage:      "extraction",
			Message:    "booking detected but extraction weak",
			Confidence: extractRes.Confidence,
			Fields:     extractRes.FieldsFound,
		}, nil
	}

	booking := extractRes.Booking

	// ⭐ Step 3 — Duplicate check
	log.Println("📦 Checking duplicate booking")

	exists, err := s.Repo.Exists(ctx, booking.UserID, booking.BookingRef)
	if err != nil {
		return nil, err
	}

	if exists {
		return &PipelineResult{
			Success: false,
			Stage:   "storage",
			Message: "duplicate booking",
		}, nil
	}

	// ⭐ Step 4 — Save booking
	log.Println("💾 Saving booking to database")

	err = s.Repo.Save(ctx, booking)
	if err != nil {
		return nil, err
	}

	log.Println("✅ Booking stored successfully")

	return &PipelineResult{
		Success:    true,
		Stage:      "completed",
		Message:    "booking stored",
		Confidence: extractRes.Confidence,
		Booking:    booking,
	}, nil
}