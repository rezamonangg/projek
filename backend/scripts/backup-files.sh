#!/bin/bash
set -e

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/var/backups/projek"
SOURCE_DIR="/app/uploads"

mkdir -p "$BACKUP_DIR"

echo "Starting file storage backup..."

tar -czf "$BACKUP_DIR/files_backup_$TIMESTAMP.tar.gz" -C "$(dirname "$SOURCE_DIR")" "$(basename "$SOURCE_DIR")"

find "$BACKUP_DIR" -name "files_backup_*.tar.gz" -mtime +7 -delete

echo "File storage backup completed: files_backup_$TIMESTAMP.tar.gz"

ls -lh "$BACKUP_DIR/files_backup_$TIMESTAMP.tar.gz"
