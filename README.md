# ByggaBat

A small Go, SQLite, and HTMX prototype for Python code-reading practice. Students can work through seeded flawed examples in diagnose, Socratic tutor, and code-review modes. Instructors can inspect the bug bank and recent student work, view answer notes, and generate a calibrated variant from a reference solution.

## Run locally

Requirements: Go 1.24+ and a C compiler (the SQLite driver uses CGO).

```sh
export ANTHROPIC_API_KEY=your-api-key
go run .
```

Open [http://localhost:8080](http://localhost:8080). The database is created as `byggabat.db`; set `DATABASE_PATH`, `PORT`, or `LISTEN_ADDR` to override the defaults. The server binds to localhost by default. `ANTHROPIC_MODEL` can override the default `claude-haiku-4-5` model.

The app seeds three beginner examples and five medium-sized Python codebases on first startup. Diagnose, tutor, review, and exercise generation use the Anthropic Messages API. Attempts are saved for the instructor view at `/instructor`. The prototype has no authentication; run it only in a trusted local environment and configure a valid Anthropic API key for AI features.

## Tests

```sh
go test ./...
```
