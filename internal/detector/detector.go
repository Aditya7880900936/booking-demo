package detector

import (
	"regexp"
	"strings"
)

type DetectionResult struct {
	IsBooking  bool
	Confidence float64
	Signals    []string
}

var knownDomains = []string{
	"indigo",
	"airindia",
	"makemytrip",
	"goibibo",
	"booking",
	"agoda",
	"irctc",
}

var keywords = []string{
	"pnr",
	"itinerary",
	"booking confirmed",
	"flight",
	"departure",
	"arrival",
	"check-in",
	"ticket",
}

var pnrRegex = regexp.MustCompile(`[A-Z0-9]{5,6}`)
var flightRegex = regexp.MustCompile(`[A-Z]{2}[0-9]{3,4}`)

func DetectEmail(subject, sender, body string) DetectionResult {

	score := 0.0
	signals := []string{}

	text := strings.ToLower(subject + " " + body)
	senderLower := strings.ToLower(sender)

	//  Domain heuristic
	for _, d := range knownDomains {
		if strings.Contains(senderLower, d) {
			score += 0.4
			signals = append(signals, "known_domain")
			break
		}
	}

	// Keyword scoring
	keywordHits := 0
	for _, k := range keywords {
		if strings.Contains(text, k) {
			keywordHits++
		}
	}

	if keywordHits > 0 {
		score += float64(keywordHits) * 0.1
		signals = append(signals, "keyword_match")
	}

	//  Regex signals
	if pnrRegex.MatchString(body) {
		score += 0.3
		signals = append(signals, "pnr_pattern")
	}

	if flightRegex.MatchString(body) {
		score += 0.2
		signals = append(signals, "flight_pattern")
	}

	return DetectionResult{
		IsBooking:  score >= 0.5,
		Confidence: score,
		Signals:    signals,
	}
}
