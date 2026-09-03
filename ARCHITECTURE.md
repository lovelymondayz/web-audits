# web-audits

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        web-audits                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │   Client     │───▶│   FastAPI   │───▶│  Services   │     │
│  │  (HTTP/REST) │    │   (8095)    │    │  (Business) │     │
│  └─────────────┘    └─────────────┘    └─────────────┘     │
│                              │                    │         │
│                              ▼                    ▼         │
│                       ┌─────────────┐    ┌─────────────┐   │
│                       │   SQLite    │    │  PageSpeed  │   │
│                       │   (db)      │    │   Insights  │   │
│                       └─────────────┘    └─────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Python 3.11+ |
| Framework | FastAPI |
| Server | Uvicorn |
| Database | SQLite |
| ORM | SQLAlchemy |
| Validation | Pydantic |
| Testing | pytest |
| Containerization | Docker |

## Design Decisions

- **FastAPI**: Async support, automatic OpenAPI docs, Pydantic validation
- **SQLite**: Lightweight storage for audit results and history
- **PageSpeed Insights API**: Google's web performance metrics
- **Service Layer**: Separates audit logic from API routes
- **Configurable API Key**: Support for multiple API keys via environment

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