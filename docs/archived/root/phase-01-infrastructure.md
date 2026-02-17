# Phase 1: Infrastructure

Core infrastructure: configuration, database connection, logging, common utilities.

## Task 1.1: Configuration System

**Commits:**
- `feat: add configuration loader with env and yaml support`
- `feat: add configuration validation using go-playground/validator`
- `task: add configuration struct with all required fields`

**Fields needed:**
```go
type Config struct {
    AppName     string `validate:"required"`
    Environment string `validate:"oneof=development staging production"`
    
    Server struct {
        Port         int           `validate:"required"`
        ReadTimeout  time.Duration
        WriteTimeout time.Duration
    }
    
    Database struct {
        Host     string
        Port     int
        User     string
        Password string
        DBName   string
        SSLMode  string
    }
    
    Redis struct {
        Host     string
        Port     int
        Password string
        DB       int
    }
    
    Session struct {
        Secret   string
        MaxAge   time.Duration
        Secure   bool
        HttpOnly bool
    }
    
    Storage struct {
        Type      string // local or s3
        LocalPath string
        S3Bucket  string
        S3Region  string
        S3Endpoint string
        S3AccessKey string
        S3SecretKey string
    }
    
    Email struct {
        Provider   string // smtp, sendgrid
        SMTPHost   string
        SMTPPort   int
        SMTPUser   string
        SMTPPass   string
        FromEmail  string
        FromName   string
    }
    
    Metrics struct {
        Enabled bool
        Path    string
    }
}
```

## Task 1.2: Logger Implementation

**Commits:**
- `feat: implement zerolog logger with custom JSON format`
- `feat: add HTTP request logging middleware`
- `feat: add context-aware logging with fields`

**Requirements:**
- JSON format by default
- Fields in order: timestamp, level, appname, message, context, http
- Context injection support (memberId, projectId, taskId)
- HTTP logging with status, latency, path

## Task 1.3: Database Connection

**Commits:**
- `feat: add PostgreSQL connection pool with pgx`
- `feat: add database health check`
- `feat: add connection retry logic`

**Details:**
- Use pgx (native Go PostgreSQL driver)
- Connection pool configuration
- Health check endpoint
- Graceful shutdown handling

## Task 1.4: Database Migrations (dbmate)

**Commits:**
- `chore: add dbmate installation instructions`
- `feat: add initial database schema migration`
- `feat: add migration runner in Go code`
- `task: add indexes for performance`

**Initial Schema:**
```sql
-- 001_initial_schema.sql
-- migrate:up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- for text search

-- Communities (tenants)
CREATE TABLE communities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Settings per community
CREATE TABLE community_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    community_id UUID REFERENCES communities(id) ON DELETE CASCADE,
    metrics_enabled BOOLEAN DEFAULT false,
    alert_email VARCHAR(255),
    error_threshold INT DEFAULT 10,
    storage_type VARCHAR(20) DEFAULT 'local',
    email_provider VARCHAR(20) DEFAULT 'smtp',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- migrate:down
-- DROP TABLE community_settings;
-- DROP TABLE communities;
```

## Task 1.5: Redis Connection

**Commits:**
- `feat: add Redis connection with go-redis`
- `feat: add Redis health check`
- `feat: add Redis connection pooling`

## Task 1.6: Common Utilities

**Commits:**
- `feat: add standard API response helpers`
- `feat: add custom error types and handling`
- `feat: add input validation helpers`
- `feat: add password hashing utilities (bcrypt)`

## Task 1.7: HTTP Router Setup

**Commits:**
- `feat: setup Chi router with middleware`
- `feat: add CORS middleware`
- `feat: add request ID middleware`
- `feat: add panic recovery middleware`

## Task 1.8: Health Check Endpoint

**Commits:**
- `feat: add /health endpoint for liveness`
- `feat: add /ready endpoint for readiness`
- `feat: add /health/detailed with all dependencies`

---

## Total Commits: 22
**Estimated Time:** 2-3 days
