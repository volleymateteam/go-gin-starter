# database-strategy.md

## Strategy Overview

We use a **PostgreSQL** managed RDS instance. Models are designed to support:

- UUID primary keys (via `uuid_generate_v4()`)
- Timestamps (`created_at`, `updated_at`, soft deletes)
- Enums for constrained values (gender, role, status)
- JSONB for permissions & logs

## Extension Setup

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

## Environments

- `volleymate_go_dev`: local development (optional seed data)
- `volleymate_db`: production RDS database

## Migration Strategy

- Use GORM auto-migrations during dev
- Lock schemas for production
- Run migration via CLI (`go run migrate.go`)

## Tables

- `users`
- `teams`
- `matches`
- `seasons`
- `waitlist_entries`
- `admin_action_logs`
- `videos`, `scout_files`
