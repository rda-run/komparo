# Komparo

Komparo (Esperanto for "comparison") is a standalone, fast, and reliable CLI
tool written in Go to compare PostgreSQL database schemas.

Relying solely on migration version numbers (e.g., `schema_migrations`) can lead
to silent failures, partial migrations, or out-of-band schema changes. Komparo
solves this by taking a structural snapshot of your expected database schema
during your CI/CD pipeline and allowing you to validate it against any live
database environment (like production) at runtime.

## Features

- **Snapshot Generation**: Extract the exact structure of a reference PostgreSQL
  database into a standard JSON file.
- **Live Validation**: Compare a JSON snapshot against a live database to
  instantly detect schema drifts.
- **Comprehensive Detection**: Identifies missing/extra tables, column type
  mismatches, missing indices, differing constraints, nullable rule violations,
  and more.
- **Zero Dependencies**: An independent statically-linked Go binary.
- **CI/CD Ready**: Easily integrate into GitHub Actions, GitLab CI, etc.,
  without requiring production credentials during the build phase.

## Installation

```bash
# Linux/MacOS/Windows
curl -fsSL https://install.rda.run/rda-run/komparo@latest! | bash

# RPM-Based distros
sudo dnf localinstall -y https://rpm.rda.run/repo.rpm
sudo dnf install komparo
```

## Configuration

Komparo supports loading database connection parameters from a `.env` file located in the current directory. When these variables are present, the `--db` flag becomes optional for the `snapshot` and `validate` commands.

Supported environment variables:

- `PGSQL_HOST` (e.g., `192.168.1.2`)
- `PGSQL_PORT` (e.g., `5432`)
- `PGSQL_USER` (e.g., `postgres`)
- `PGSQL_DB` (e.g., `txlog_dev`)
- `PGSQL_PASSWORD` (optional)
- `PGSQL_SSLMODE` (optional, e.g., `disable` or `require`)

If you provide the `--db` flag explicitly, it will override the `.env` configurations.

## Quick Start

### 1. Capture the Expected Schema

Run your migrations against an ephemeral database in your CI pipeline, then
capture the snapshot:

```bash
komparo snapshot --db "postgres://user:pass@localhost:5432/ephemeral_db?sslmode=disable" --out expected_schema.json
```

### 2. Validate a Target Database

Point Komparo to your staging or production database using the previously
generated snapshot:

```bash
komparo validate --db "postgres://user:pass@prod-db.internal:5432/production_db?sslmode=require" --file expected_schema.json
```

The schema can be hosted online, if you want.

```bash
komparo validate --db "postgres://user:pass@prod-db.internal:5432/production_db?sslmode=require" --file https://example.com/expected_schema.json
```

### Output Example

If differences are found, Komparo will report them:

```text
⚠️ Schema Divergence Detected!

STATUS     TYPE      OBJECT                DETAILS
Missing    Column    users.last_login_at   Expected: timestamp 
Mismatch   Column    orders.total_amount   Expected: numeric(10,2) | Actual: numeric(12,4)
Mismatch   Index     idx_users_email       Expected: Match | Actual: Differ
Extra      View      vw_legacy_reports     Actual: Exists
```

## Why Komparo?

Tools like `pgdiff` or `migra` are excellent, but they require **two
simultaneous live database connections**. This is often impossible or insecure
in CI/CD pipelines where the build server doesn't have access to the production
database.

Komparo splits the process into two phases (`snapshot` and `validate`), allowing
you to embed or ship the lightweight JSON schema definition alongside your
application, validating the database strictly at deployment time or on-demand
without connecting two environments.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.
