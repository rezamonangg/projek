#!/bin/bash
set -e

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/var/backups/projek"
DB_NAME="projek"
DB_USER="projek"

mkdir -p "$BACKUP_DIR"

echo "Starting database backup..."

pg_dump -U "$DB_USER" -h localhost "$DB_NAME" | gzip > "$BACKUP_DIR/db_backup_$TIMESTAMP.sql.gz"

find "$BACKUP_DIR" -name "db_backup_*.sql.gz" -mtime +7 -delete

echo "Database backup completed: db_backup_$TIMESTAMP.sql.gz"

ls -lh "$BACKUP_DIR/db_backup_$TIMESTAMP.sql.gz"
