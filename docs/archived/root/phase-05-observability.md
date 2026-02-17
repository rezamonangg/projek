# Phase 5: Observability & Admin Dashboard

Metrics endpoint, admin dashboard, email alerts, and community settings.

## Task 5.1: Prometheus Metrics Setup

**Commits:**
- `feat: add prometheus metrics collector`
- `feat: add HTTP request metrics middleware`
- `feat: add database query metrics`
- `feat: add business metrics (projects, tasks created)`

**Metrics:**
- `projek_http_requests_total` - Counter with method, path, status
- `projek_http_duration_seconds` - Histogram
- `projek_db_queries_total` - Counter with table, operation
- `projek_db_query_duration_seconds` - Histogram
- `projek_active_sessions` - Gauge
- `projek_communities_total` - Gauge
- `projek_projects_total` - Gauge by community
- `projek_tasks_total` - Gauge by status

## Task 5.2: Metrics Configuration

**Commits:**
- `feat: add metrics endpoint with configurable enable`
- `feat: add metrics settings to community settings`
- `feat: add admin toggle for metrics endpoint`

## Task 5.3: Admin Settings API

**Commits:**
- `feat: add admin settings repository`
- `feat: add admin settings service`
- `feat: add get community settings handler`
- `feat: add update community settings handler`

**Settings updatable by admin:**
- metrics_enabled
- alert_email
- error_threshold
- storage_type (if not using local)
- email_provider settings

## Task 5.4: Error Tracking & Alerts

**Commits:**
- `feat: add error tracking middleware`
- `feat: add error rate monitoring`
- `feat: add email alert service for errors`
- `feat: add alert throttling`

**Alert conditions:**
- Error count exceeds threshold in 5-minute window
- Send email to admin
- Include error details, stack trace

## Task 5.5: Admin Dashboard API - Stats

**Commits:**
- `feat: add admin dashboard stats service`
- `feat: add stats aggregation queries`
- `feat: add admin dashboard data endpoint`

**Stats to provide:**
```go
type AdminDashboardStats struct {
    TotalMembers      int
    ActiveMembers     int // logged in last 30 days
    TotalProjects     int
    TotalTasks        int
    TasksByStatus     map[TaskStatus]int
    TotalWikiPages    int
    StorageUsed       int64
    RecentActivity    []ActivityLog
}
```

## Task 5.6: Activity Logging

**Commits:**
- `feat: add activity log model`
- `feat: add activity logging middleware`
- `feat: add get recent activity endpoint`

**Log key actions:**
- Project created
- Task created/moved
- Member invited
- Wiki page updated
- File uploaded

## Task 5.7: Prometheus Development Setup

**Commits:**
- `chore: add prometheus configuration`
- `chore: add prometheus to docker-compose`
- `chore: add grafana to docker-compose (optional)`

## Task 5.8: Email Notification Settings

**Commits:**
- `feat: add email notification preferences`
- `feat: add task assignment notification`
- `feat: add task status change notification`

---

## Total Commits: 20
**Estimated Time:** 2-3 days
