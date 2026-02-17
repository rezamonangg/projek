# Backend Phase 5: Observability & Admin

Metrics endpoint, admin dashboard API, email alerts, and community settings.

## Task 5.1: Prometheus Metrics Setup

**Commits:**
- `feat: add prometheus metrics collector`
- `feat: add HTTP request metrics middleware`
- `feat: add database query metrics`
- `feat: add business metrics`

**Package:** `internal/metrics/`

**Dependencies:**
- github.com/prometheus/client_golang

**Metrics:**
- projek_http_requests_total
- projek_http_duration_seconds
- projek_db_queries_total
- projek_active_sessions
- projek_projects_total
- projek_tasks_total

## Task 5.2: Metrics Configuration

**Commits:**
- `feat: add metrics endpoint with configurable enable`
- `feat: add metrics settings to community settings`

## Task 5.3: Admin Settings API

**Commits:**
- `feat: add admin settings repository`
- `feat: add admin settings service`
- `feat: add get community settings handler`
- `feat: add update community settings handler`

**Package:** `internal/admin/`

## Task 5.4: Error Tracking & Alerts

**Commits:**
- `feat: add error tracking middleware`
- `feat: add error rate monitoring`
- `feat: add email alert service for errors`

## Task 5.5: Admin Dashboard Stats

**Commits:**
- `feat: add admin dashboard stats service`
- `feat: add stats aggregation queries`
- `feat: add admin dashboard data endpoint`

## Task 5.6: Activity Logging

**Commits:**
- `feat: add activity log model`
- `feat: add activity logging middleware`
- `feat: add get recent activity endpoint`

## Task 5.7: Database Migrations

**Commits:**
- `task: add activity_logs table migration`

---

## Total Backend Commits (Phase 5): 18
