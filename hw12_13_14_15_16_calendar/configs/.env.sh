#!/bin/bash
export LOG_LEVEL="DEBUG"
export LOG_PATH="logs/access.log"
# export LOG_PATH="/opt/calendar/logs/access.log"
export POSTGRES_DSN="host=localhost port=5432 user=postgres password=postgres dbname=calendar sslmode=disable"
# export POSTGRES_DSN="host=host.docker.internal port=5432 user=postgres password=postgres dbname=calendar sslmode=disable"
export MIGRATION_PATH="migrations"
export HTTP_HOST="0.0.0.0"
export HTTP_PORT="8080"
export GRPC_HOST="127.0.0.1"
export GRPC_PORT="50051"
export STORAGE_TYPE="sql"
export QUEUE_URL="amqp://guest:guest@localhost:5672/"
# export QUEUE_URL="amqp://guest:guest@host.docker.internal:5672/"
export QUEUE_NAME="notifications"
export SCHEDULER_INTERVAL=10
export SCHEDULER_CLEANUP_DAYS=365
export SENDER_LOG_LEVEL="info"