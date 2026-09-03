# web-audits

API service for running website performance audits using Google PageSpeed Insights.

## Quick Start

```bash
# Development
make dev

# Docker
make up

# View logs
make logs
```

## Features

- Website performance auditing via PageSpeed Insights
- Audit history and trend tracking
- Detailed performance reports
- RESTful API with OpenAPI documentation
- SQLite storage for audit results

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/audits` | List all audits |
| GET | `/audits/{id}` | Get audit by ID |
| POST | `/audits` | Create new audit |
| DELETE | `/audits/{id}` | Delete audit |
| POST | `/audit/url` | Audit a specific URL |
| GET | `/audit/{id}/report` | Get detailed audit report |
| GET | `/scores` | List performance scores |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PAGESPEED_API_KEY` | Google PageSpeed API key | - |
| `DATABASE_PATH` | SQLite database file path | `./data/audits.db` |
| `OUTPUT_DIR` | Report output directory | `./output` |
| `API_KEY` | Service API authentication key | - |

## Deployment

```bash
# Build and deploy with Docker Compose
make deploy

# Or manually
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

## API Documentation

Once running, visit:
- Swagger UI: http://localhost:8095/docs
- ReDoc: http://localhost:8095/redoc