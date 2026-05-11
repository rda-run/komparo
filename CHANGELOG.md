<!-- markdownlint-disable MD024 -->
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic
Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.1] - 2026-05-11

### Fixed

- **Migrate URL Support**: The `migrate` command now supports fetching JSON
  snapshots via HTTP/HTTPS URLs using the `-f` flag, matching the behavior
  of the `validate` command.

## [0.4.0] - 2026-05-09

### Added

- **Migration Script Generation**: New `migrate` command that generates SQL DDL
  fix scripts (`komparo_fix_*.sql`) to align a database with an expected schema
  snapshot. This command works with read-only database users and never executes
  DDL directly - the user must review and execute the generated SQL manually.
  - Analyzes differences between snapshot and live database
  - Generates ordered SQL statements respecting dependencies
  - Supports all schema object types: tables, columns, indexes, constraints,
    sequences, views, materialized views, triggers, functions, enums, extensions
  - DDL is organized in sections: DROP, ALTER, RECREATE, CREATE
  - Includes transaction wrapper with COMMIT/ROLLBACK for safe testing

### Security

- **Read-Only by Design**: The `migrate` command maintains Komparo's philosophy
  of never executing DDL directly. It only generates SQL files for manual review,
  allowing usage with restricted database users that have only SELECT privileges
  on system catalogs.

## [0.3.0] - 2026-03-24

### Added

- **Environment Variable Support**: Added support for loading database connection
  parameters (e.g., `PGSQL_HOST`, `PGSQL_USER`) from a `.env` file, making
  the `--db` flag optional for `snapshot` and `validate` commands.

### Security

- **Log Sanitization**: The `snapshot` command now obscures the password
  in the connection string output (e.g., `password=***`) to prevent
  accidental credential leaks in logs.

## [0.2.0] - 2026-03-24

### Added

- **Remote JSON Validation**: Added support for fetching JSON snapshot files
  via HTTP and HTTPS URLs in the `validate` command using the `-f` flag.

## [0.1.1] - 2026-03-23

### Added

- **Colored Terminals**: ANSI colored status indicators in the `validate`
  divergence table (Missing: Red, Mismatch: Yellow, Extra: Cyan) seamlessly
  integrated with Go's native `tabwriter`.

## [0.1.0] - 2026-03-23

### Added

- **Core CLI Foundation**: Basic Cobra CLI setup with `komparo` executable.
- **Snapshot Generation**: `snapshot` command to extract structural schema
  metadata from a PostgreSQL database into a standardized JSON file.
  - Supports extracting: Tables, Columns, Indices, Constraints, Sequences,
    Views, Materialized Views, Triggers, Functions, Enums, and Extensions.
- **Validation Engine**: `validate` command to compare a live database against a
  generated JSON snapshot.
- **Output Formats**: Tabular terminal output (default) and structured JSON
  output (`--format json`) for CI/CD integration.
- **GitHub Actions**: Automated CI pipeline running linting and tests on all PRs
  and pushes to `main`.
- **Makefile**: Pre-configured build scripts for Linux, Windows, and macOS
  (Intel & ARM64) with aggressive production optimizations (`CGO_ENABLED=0`,
  `-trimpath`, `-s -w`).
- **Version Command**: `version` subcommand printing the dynamically injected
  semver tag.
- **JSON Schema**: Formal specification for the expected
  `komparo-snapshot-v1.schema.json` output format.
