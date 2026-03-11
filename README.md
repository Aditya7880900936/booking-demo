# 📩 Booking Email Extraction Service

> A backend service that detects travel booking emails, extracts structured booking information, and stores it securely — built in **Go**.

---

## 🚀 Problem Statement

Users receive booking confirmations via email (flights, hotels, trains, etc.) in dozens of different formats from dozens of different providers.

This system:

- **Detects** whether an email is a booking confirmation
- **Extracts** key booking details using regex-based parsing
- **Normalizes** data into a fixed canonical schema
- **Stores** bookings in MongoDB with duplicate prevention
- **Exposes APIs** to query and retrieve bookings

Designed with a focus on:

- 🔒 **Privacy preservation** — no third-party AI services required
- 💸 **Low processing cost** — rule-based detection before AI fallback
- 🧱 **Modular architecture** — each stage is independently replaceable
- ⚡ **Minimal dependencies** — lightweight and fast

---

## 🧠 Architecture Overview

Pipeline-based design — each email flows through stages sequentially:

```
POST /ingest/email
        │
        ▼
┌──────────────────┐
│  Detection Engine │  ← keyword scoring, sender heuristics, regex signals
└──────────────────┘
        │ (confidence ≥ threshold)
        ▼
┌──────────────────┐
│ Extraction Engine │  ← regex-based field parsing, confidence scoring
└──────────────────┘
        │
        ▼
┌──────────────────┐
│  Normalization   │  ← canonical schema mapping
└──────────────────┘
        │
        ▼
┌──────────────────┐
│    MongoDB       │  ← persistence, duplicate prevention
└──────────────────┘
        │
        ▼
   Query APIs
```

### Key Components

| Component | Responsibility |
|-----------|---------------|
| **Detector** | Keyword scoring, sender heuristics, regex signal detection |
| **Extractor** | Regex-based flight field parsing, confidence-driven extraction |
| **Pipeline Service** | Orchestrates detection → extraction → normalization → storage |
| **Repository** | MongoDB persistence with duplicate prevention |
| **Handler** | HTTP request/response handling via Gin |

---

## ⚙️ Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go |
| Web Framework | Gin |
| Database | MongoDB |
| Parsing | Regex-based |

---

## 📦 Project Structure

```
booking-system/
├── cmd/
│   └── server/          # Entry point
├── internal/
│   ├── handler/         # HTTP handlers
│   ├── pipeline/        # Pipeline orchestration
│   ├── detector/        # Booking detection engine
│   ├── extractor/       # Booking field extractor
│   └── repository/      # MongoDB persistence layer
├── models/              # Shared data models / schemas
└── pkg/                 # Shared utilities
```

---

## 🛠️ Setup & Running

### Prerequisites

- Go 1.21+
- MongoDB running locally

### 1. Clone the Repository

```bash
git clone <repo-url>
cd booking-system
```

### 2. Start MongoDB

Ensure MongoDB is running at the default address:

```
mongodb://localhost:27017
```

### 3. Run the Server

```bash
go run ./cmd/server
```

Server starts at: `http://localhost:8080`

---

## 📡 API Reference

### `GET /health` — Health Check

```json
{
  "status": "ok"
}
```

---

### `POST /ingest/email` — Ingest an Email

Simulates receiving an email and runs it through the full pipeline.

**Request Body:**

```json
{
  "user_id": "u1",
  "sender": "noreply@indigo.in",
  "subject": "Flight Booking Confirmed",
  "body": "Passenger: Aditya Srivastav. Your PNR AB12CD. Flight 6E203 from Delhi to Bangalore on 20 Mar 2026"
}
```

**Success Response** (booking detected and stored):

```json
{
  "success": true,
  "stage": "completed",
  "message": "booking stored",
  "confidence": 0.9,
  "booking": {
    "user_id": "u1",
    "type": "flight",
    "provider": "INDIGO",
    "booking_ref": "AB12CD",
    "passenger_name": "Aditya Srivastav",
    "flight_number": "6E203",
    "departure": "Delhi",
    "arrival": "Bangalore",
    "start_date": "2026-03-20T00:00:00Z"
  }
}
```

**Non-Booking Email Response** (rejected at detection stage):

```json
{
  "user_id": "u1",
  "sender": "newsletter@amazon.com",
  "subject": "Great Deals",
  "body": "Big sale is live"
}
```

```json
{
  "success": false,
  "stage": "detection",
  "message": "not a booking email",
  "confidence": 0.2
}
```

---

### `GET /bookings` — Get All Bookings

Returns a list of all stored bookings for all users.

---

### `GET /bookings/upcoming` — Get Upcoming Bookings

Returns bookings where the travel date is **≥ current date**.

---

## 🧩 Design Decisions

**Rule-based detection before AI**
Reduces cost and latency. AI is reserved as a future fallback for unrecognized templates.

**Confidence scoring**
Both detection and extraction produce confidence scores, enabling the pipeline to make informed decisions at each stage rather than failing silently.

**Canonical schema**
All booking types (flight, hotel, train) normalize into a single fixed schema, simplifying downstream queries and storage.

**MongoDB**
Chosen for its flexible document model, which accommodates varying booking structures without rigid migrations.

**Pipeline orchestration**
A clean pipeline service separates concerns — each stage can be developed, tested, and replaced independently.

---

## 🔮 Future Improvements

- [ ] Gmail / Outlook API integration for real email ingestion
- [ ] AI fallback extractor for unrecognized email templates
- [ ] Hotel and train booking extraction
- [ ] OCR support for scanned or image-based tickets
- [ ] Message queue (e.g. Kafka/RabbitMQ) for large-scale ingestion
- [ ] Multi-user authentication and per-user booking access

---

## 🎯 Project Goal

This prototype demonstrates how a scalable, **privacy-first** backend system can be designed to extract booking intelligence from emails — using minimal AI, rule-based parsing, and a clean pipeline architecture.

---

## 📝 License

MIT
