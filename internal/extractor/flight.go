package extractor

import (
	"booking-system/models"
	"regexp"
	"strings"
	"time"
)

type ExtractionResult struct {
	Booking     *models.Booking
	Confidence  float64
	FieldsFound []string
}

var (
	pnrRegex       = regexp.MustCompile(`[A-Z0-9]{5,6}`)
	flightRegex    = regexp.MustCompile(`[A-Z]{2}[0-9]{3,4}`)
	dateRegex      = regexp.MustCompile(`\d{1,2}\s[A-Za-z]{3}\s\d{4}`)
	departureRegex = regexp.MustCompile(`from\s([A-Za-z]+)`)
	arrivalRegex   = regexp.MustCompile(`to\s([A-Za-z]+)`)
	passengerRegex = regexp.MustCompile(`passenger[:\s]+([A-Za-z\s]+)`)
)

func normalize(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	return strings.ToLower(text)
}

func detectProvider(text string) string {

	t := strings.ToLower(text)

	switch {
	case strings.Contains(t, "indigo"):
		return "INDIGO"
	case strings.Contains(t, "air india"):
		return "AIR_INDIA"
	case strings.Contains(t, "spicejet"):
		return "SPICEJET"
	default:
		return "UNKNOWN"
	}
}

func ExtractFlight(body string, userID string) ExtractionResult {

	fields := []string{}
	score := 0.0

	raw := body
	text := normalize(body)

	booking := &models.Booking{
		UserID: userID,
		Type:   "flight",
		Status: "extracted",
	}

	// ⭐ PNR
	if pnr := pnrRegex.FindString(raw); pnr != "" {
		booking.BookingRef = pnr
		fields = append(fields, "pnr")
		score += 0.3
	}

	// ⭐ Flight number
	if f := flightRegex.FindString(raw); f != "" {
		booking.FlightNumber = f
		fields = append(fields, "flight_number")
		score += 0.2
	}

	// ⭐ Date
	if d := dateRegex.FindString(raw); d != "" {
		if parsed, err := time.Parse("2 Jan 2006", d); err == nil {
			booking.StartDate = parsed
			fields = append(fields, "date")
			score += 0.2
		}
	}

	// ⭐ Departure
	if m := departureRegex.FindStringSubmatch(text); len(m) > 1 {
		booking.Departure = strings.Title(m[1])
		fields = append(fields, "departure")
		score += 0.1
	}

	// ⭐ Arrival
	if m := arrivalRegex.FindStringSubmatch(text); len(m) > 1 {
		booking.Arrival = strings.Title(m[1])
		fields = append(fields, "arrival")
		score += 0.1
	}

	// ⭐ Passenger
	if m := passengerRegex.FindStringSubmatch(text); len(m) > 1 {
		booking.PassengerName = strings.Title(strings.TrimSpace(m[1]))
		fields = append(fields, "passenger")
		score += 0.1
	}

	booking.Confidence = score
	booking.Provider = detectProvider(raw)

	return ExtractionResult{
		Booking:     booking,
		Confidence:  score,
		FieldsFound: fields,
	}
}
