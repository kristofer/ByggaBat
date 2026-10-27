# Use cases and roadmap

This is a working list for product discovery, not a commitment to implement every candidate. Use instructor mode and **Act as student** to compare both roles while refining these workflows.

## Current use cases

- **Student practice:** A learner works through seeded Python examples and codebases, choosing diagnose, tutor, or code-review mode and receiving AI feedback.
- **Instructor review:** An instructor browses the bug bank, reads hints and expected diagnoses, and reviews the 50 most recent submitted attempts.
- **Exercise generation:** An instructor provides a reference solution, curriculum week, and skill to generate a flawed exercise for student practice.
- **Local or single-host trial:** A facilitator runs the app for a trusted group, with attempts and exercises stored in SQLite. Docker Compose supports host access and a persistent database volume.

## Future use-case candidates

- **Classroom administration:** Manage classes, enrollments, and student identities; separate activity between classes.
- **Progress review:** Summarize learner practice over time and help instructors identify concepts that need follow-up.
- **Exercise authoring and sharing:** Create, edit, import, export, and organize exercises without relying only on seeded content or AI generation.
- **Curriculum expansion:** Support additional languages, curricula, and skill progressions.
- **Configurable AI feedback:** Let instructors select or configure models and feedback behavior for a course.
- **More deployment contexts:** Support authenticated, multi-instructor deployments beyond a trusted single-host prototype.

## Product and operational gaps

- There is no authentication or class-level access control; instructor notes and student activity are accessible to anyone who can reach the app.
- AI features depend on a configured Anthropic API key and currently use the configured model for all users.
- Exercise management is limited to seeded exercises and AI-generated variants; generated content should be reviewed by an instructor.
- SQLite and one named Docker volume support a single-host trial, but backup, restore, and multi-instance operation are not provided.

Use feedback from the current workflows to prioritize which future use cases to validate before implementation.
