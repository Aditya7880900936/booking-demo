package pipeline

import (
	"context"

	"booking-system/internal/detector"
	"booking-system/internal/extractor"
	"booking-system/internal/repository"
	"booking-system/models"
)

type Service struct {
	Repo *repository.BookingRepository
}

type PipelineResult struct {
	Stage      string
	Message    string
	Confidence float64
	Booking    *models.Booking
	Signals    []string
	Fields     []string
}

func (s *Service) ProcessEmail(
	ctx context.Context,
	userID string,
	sender string,
	subject string,
	body string,
) (*PipelineResult, error) {

	// ⭐ Step 1 — Detection
	detectRes := detector.DetectEmail(subject, sender, body)

	if !detectRes.IsBooking {
		return &PipelineResult{
			Stage:      "detection",
			Message:    "not a booking email",
			Confidence: detectRes.Confidence,
			Signals:    detectRes.Signals,
		}, nil
	}

	// ⭐ Step 2 — Extraction
	extractRes := extractor.ExtractFlight(body, userID)

	if extractRes.Confidence < 0.4 {
		return &PipelineResult{
			Stage:      "extraction",
			Message:    "booking detected but extraction weak",
			Confidence: extractRes.Confidence,
			Fields:     extractRes.FieldsFound,
		}, nil
	}

	booking := extractRes.Booking

	// ⭐ Step 3 — Duplicate check
	exists, err := s.Repo.Exists(ctx, booking.UserID, booking.BookingRef)
	if err != nil {
		return nil, err
	}

	if exists {
		return &PipelineResult{
			Stage:   "storage",
			Message: "duplicate booking",
		}, nil
	}

	// ⭐ Step 4 — Save
	err = s.Repo.Save(ctx, booking)
	if err != nil {
		return nil, err
	}

	return &PipelineResult{
		Stage:      "completed",
		Message:    "booking stored",
		Confidence: extractRes.Confidence,
		Booking:    booking,
	}, nil
}
