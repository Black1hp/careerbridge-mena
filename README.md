# CareerBridge MENA

Scholarship, internship, and competition aggregation platform for the Arab world.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

CareerBridge MENA collects and organizes opportunities across the Arab world — scholarships, internships, competitions, and more. Search by type, country, and deadline with real-time countdown timers.

## Features

- Search scholarships, internships, and competitions
- Filter by country and type
- Deadline tracking with live countdown
- Elasticsearch-powered full-text search
- Automated web crawling with Playwright
- Desktop app via Electron

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 18, Electron, Vite |
| Backend | Go, Gorilla Mux |
| Crawler | Go, Playwright |
| Database | PostgreSQL 16 |
| Search | Elasticsearch 8 |
| Infra | Docker Compose |

## Project Structure

```
careerbridge-mena/
├── backend/           # Go API server
│   ├── cmd/api/       # Entry point
│   └── internal/      # Handlers, models, database, search
├── crawler/           # Go + Playwright crawlers
│   ├── cmd/crawler/   # Entry point
│   └── internal/      # Scrapers, parsers
├── frontend/          # React + Electron app
│   ├── src/           # Components, pages
│   └── electron/      # Electron shell
├── database/          # SQL migrations
├── docker-compose.yml # PostgreSQL + Elasticsearch
└── .env.example       # Environment config template
```

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose

### 1. Start databases

```bash
docker compose up -d
```

This starts PostgreSQL (port 5432) and Elasticsearch (port 9200) with the schema auto-loaded.

### 2. Run the API

```bash
cd backend
go mod tidy
go run ./cmd/api
```

API runs on `http://localhost:8080`.

### 3. Run the frontend

```bash
cd frontend
npm install
npm run dev
```

UI runs on `http://localhost:5173`.

### 4. Run the crawler

```bash
cd crawler
go mod tidy
go run ./cmd/crawler
```

Scrapes configured sources and POSTs results to the API.

### 5. Run with Electron

```bash
cd frontend
npm run electron:dev
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/search` | Search opportunities (query, type, country, page, limit) |
| GET | `/api/v1/opportunities/:id` | Get single opportunity |
| GET | `/api/v1/countries` | List all countries |
| POST | `/api/v1/ingest` | Bulk insert opportunities (used by crawler) |

## Roadmap

- [x] Phase 1 — Aggregation (search, filters, crawling)
- [ ] Phase 2 — Notifications (email, in-app, deadline alerts)
- [ ] Phase 3 — AI matching (resume review, cover letter, smart matching)

## License

MIT
