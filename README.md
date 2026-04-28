# Resume Screener

An AI-powered resume screening platform that analyzes resumes against job descriptions, scores candidates, and ranks applicants — built with Go (backend) and Next.js (frontend).

## Features

- Upload and parse resumes (PDF)
- Create job postings with descriptions
- AI-driven analysis: match score, strengths, missing skills, and recommendation
- Candidate ranking per job
- JWT authentication with rate limiting (Redis)
- File storage via AWS S3 or Cloudflare R2
- Supports OpenAI or local Ollama as AI provider

## Tech Stack

| Layer     | Technology                          |
|-----------|-------------------------------------|
| Backend   | Go, Gin, PostgreSQL, Redis          |
| Frontend  | Next.js 15, TypeScript, Tailwind CSS |
| AI        | OpenAI API / Ollama (local)         |
| Storage   | AWS S3 / Cloudflare R2 / MinIO      |
| Infra     | Docker Compose, Render              |

## Project Structure

```
resume/
├── backend/
│   ├── cmd/              # Entry point
│   ├── configs/          # Config loading from env
│   ├── internal/
│   │   ├── api/          # Gin router, handlers, middleware
│   │   ├── model/        # Domain models
│   │   ├── repository/   # PostgreSQL queries
│   │   └── service/      # Business logic
│   ├── pkg/
│   │   ├── ai/           # AI provider interface (OpenAI / Ollama)
│   │   ├── parser/       # PDF text extraction
│   │   └── storage/      # S3-compatible file storage
│   └── migrations/       # SQL migration files
└── frontend/
    ├── app/              # Next.js App Router pages
    └── lib/              # API client, auth helpers
```

## Getting Started

### Prerequisites

- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- Go 1.21+ (for local backend development)
- Node.js 18+ (for local frontend development)

### Run with Docker Compose

```bash
cp .env.example .env
# Edit .env and set JWT_SECRET and optionally OPENAI_API_KEY
docker compose up --build
```

Services started:

| Service   | URL                        |
|-----------|----------------------------|
| Frontend  | http://localhost:3000      |
| Backend   | http://localhost:8080      |
| MinIO UI  | http://localhost:9001      |
| PostgreSQL| localhost:5433             |
| Redis     | localhost:6379             |

### Local Development

**Backend:**

```bash
cd backend
cp ../.env.example .env   # fill in values
go run ./cmd/main.go
```

**Frontend:**

```bash
cd frontend
npm install
npm run dev
```

## Environment Variables

Copy `.env.example` to `.env` and configure:

```env
PORT=8081
DATABASE_URL=postgres://postgres:postgres@localhost:5432/resume_screener?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=<at-least-32-char-random-string>

# AI — leave blank to use local Ollama
OPENAI_API_KEY=sk-...

# Storage — Cloudflare R2 or AWS S3
S3_BUCKET=resume-screener
S3_REGION=auto
S3_ENDPOINT=https://<account-id>.r2.cloudflarestorage.com
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8081/api
```

## API Endpoints

All protected routes require `Authorization: Bearer <token>`.

| Method | Path                    | Auth | Description                  |
|--------|-------------------------|------|------------------------------|
| POST   | /api/auth/register      | No   | Register a new user          |
| POST   | /api/auth/login         | No   | Login, returns JWT           |
| POST   | /api/resume/upload      | Yes  | Upload a resume (PDF)        |
| GET    | /api/resume             | Yes  | List uploaded resumes        |
| GET    | /api/resume/:id         | Yes  | Get a single resume          |
| POST   | /api/job                | Yes  | Create a job posting         |
| GET    | /api/job                | Yes  | List job postings            |
| POST   | /api/analyze            | Yes  | Analyze a resume vs a job    |
| GET    | /api/results/:id        | Yes  | Get analysis result          |
| GET    | /api/ranking/:jobId     | Yes  | Get ranked candidates for job|

## AI Providers

The backend auto-selects the AI provider at startup:

- **OpenAI** — set `OPENAI_API_KEY` in `.env`
- **Ollama** (local, free) — leave `OPENAI_API_KEY` blank; install [Ollama](https://ollama.com) and pull a model

## Deployment

The project includes a `render.yaml` for one-click deploy to [Render](https://render.com). Set the required environment variables in the Render dashboard before deploying.
