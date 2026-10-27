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

## Import problems

In the instructor workspace, upload a `.json` file containing an array of problems. Each problem must have a unique `title`, `level`, positive integer `week`, `skill`, `prompt`, `code`, `hint`, and `solution`, for example:

```json
[
  {
    "title": "An empty ratings list",
    "level": "Beginner",
    "week": 1,
    "skill": "Lists and edge cases",
    "prompt": "Explain why this function fails for an empty list.",
    "code": "def average(ratings):\n    return sum(ratings) / len(ratings)",
    "hint": "What is the length of an empty list?",
    "solution": "Check for an empty list before dividing."
  }
]
```

Imports are limited to 1 MB. Titles identify problems: importing a title that already exists skips it without changing the existing problem, so re-uploading a file is safe.

## Run with Docker

Requirements: Docker with the Compose plugin.

```sh
export ANTHROPIC_API_KEY=your-api-key
docker compose up --build
```

Open [http://localhost:8080](http://localhost:8080). Compose stores the SQLite database in the persistent `byggabat-data` volume. To stop the app, run `docker compose down`; the database remains. The app listens on all container interfaces, while Compose publishes it only on the Docker host's loopback interface. Without an API key, the app still starts, but AI features cannot respond.

## Switching roles

Use **Instructor mode** in the header to open the instructor workspace. From there, choose **Act as student** to return to the student practice view; the same links are available from the exercise pages. This is a navigation convenience, not authentication: instructor pages and answer notes are available to anyone who can access the app.

## Use cases and roadmap

See [ROADMAP.md](ROADMAP.md) for current use cases, future use-case candidates, and deployment/product gaps to guide refinement.

## Tests

```sh
go test ./...
```
