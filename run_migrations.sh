#!/bin/bash

#docker compose exec db psql -U postgres -d app -f /app/internal/repository/pg/migrations/000001_create_examples.up.sql

export DATABASE_HOST=localhost
export DATABASE_USER=postgres
export DATABASE_PASSWORD=postgres
export DATABASE_NAME=app

make migration/up