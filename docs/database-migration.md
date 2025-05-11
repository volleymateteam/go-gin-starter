# 🦪 Volleymate Go Backend — Local to Production DB Migration

## 📋 Goal

Migrate important data from your **local development PostgreSQL** database (`volleymate_go_dev`) into your **AWS RDS production database** (`volleymate_db`), and ensure proper access/ownership by the `db_admin_vm` role.

---

### 🧱 Prerequisites

- Local PostgreSQL with populated `volleymate_go_dev`
- AWS RDS instance (PostgreSQL) running and accessible
- SSH access to production EC2
- `pg_dump`, `psql`, and `scp` installed

---

### ✅ Step-by-Step Instructions

#### 1️⃣ Dump local database to SQL file

```bash
pg_dump -U rezabarzegar -d volleymate_go_dev > volleymate_local_dump.sql
```

#### 2️⃣ Transfer SQL file to EC2

```bash
scp volleymate_local_dump.sql ubuntu@16.170.150.129:/home/ubuntu/
```

#### 3️⃣ SSH into EC2

```bash
ssh ubuntu@16.170.150.129
```

#### 4️⃣ Connect to `postgres` DB to drop & recreate `volleymate_db`

```bash
psql -U db_admin_vm -h <your-rds-endpoint> -d postgres
```

Inside psql:

```sql
DROP DATABASE volleymate_db;
CREATE DATABASE volleymate_db;
\q
```

#### 5️⃣ Import the local dump into RDS

```bash
psql -U db_admin_vm -h <your-rds-endpoint> -d volleymate_db -f /home/ubuntu/volleymate_local_dump.sql
```

> ⚠️ You may see errors like `role "rezabarzegar" does not exist`. This is expected due to role mismatch. See next step.

---

### 👑 Fix Ownership of Tables

Because the SQL dump contains tables owned by `rezabarzegar`, and that user doesn’t exist on RDS, reassign ownership:

```sql
DO $$ DECLARE
    r RECORD;
BEGIN
    FOR r IN (
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'public'
    ) LOOP
        EXECUTE format('ALTER TABLE public.%I OWNER TO db_admin_vm', r.tablename);
    END LOOP;
END $$;
```

---

### 🔄 Restart Backend

```bash
sudo systemctl restart volleymate-backend
```

---

### 🧲 Test with Postman

Verify endpoints like:

- `POST /register`
- `GET /matches`
- `GET /health`

---

### 🛠️ Future Recommendations

- Use migration tools like [golang-migrate](https://github.com/golang-migrate/migrate) or `goose`
- Automate dumps with timestamps: `volleymate_$(date +%F).sql`
- Schedule backups using AWS RDS automated snapshots

---

**Author**: Reza Barzegar
**Last updated**: `2025-05-11`
