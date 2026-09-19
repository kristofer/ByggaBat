# ByggaBat — Design Document

> Formerly named **BildBaut**. Renamed to **ByggaBat**.

## 1. Purpose

ByggaBat is a small local LLM + retrieval system built to support ZipCode
Wilmington's 12-week Python bootcamp curriculum. It exists to help students
develop the skills that actually matter once LLMs routinely write small code
snippets for them: reading code fluently, debugging by hypothesis, reviewing
and critiquing (including AI-generated) code, and decomposing vague problems
into verifiable pieces.

The core premise: writing small loops and functions from scratch is no longer
the differentiating skill for a working programmer. What differentiates a
senior engineer — reading unfamiliar code fast, spotting subtle bugs,
verifying correctness, decomposing ambiguous problems — is where beginners
should be trained early, not after years of experience. Writing-by-hand
practice is retained, but reframed as a *means* to build the "notional
machine" (an internal execution model of the language), not the end goal
of the curriculum.

## 2. Skill progression (redrafted)

The ratio of hands-on writing to reading/debugging/review shifts steadily
across the 12 weeks:

| Weeks | Focus | Writing : Reading+Debugging+Review |
|---|---|---|
| 1–2 | Build the notional machine — core syntax, trace-before-run, shell/git basics | ~70:30 |
| 3–4 | Transition — functions/OOP/file I/O, debugging formalized, JSON as a reading skill first | ~50:50 |
| 5–6 | Reading/debugging dominant — SQL and HTTP introduced as reading skills, debug-by-hypothesis formalized, first code review exercises | ~30:70 |
| 7–8 | Decomposition and verification — vague-problem decomposition, tests used to check AI-generated code, git PR/review workflow | ~20:80 |
| 9–10 | Unfamiliar codebases — modify code they didn't write, subtler bug-injection challenges, peer review as core activity | reading/review-dominant |
| 11–12 | Capstone — built with LLM assistance throughout, but every stage requires artifacts proving comprehension | reading/review-dominant |

Capstone requirements specifically include: a written decomposition plan
before coding, tests that verify LLM-generated code, a code-review writeup
of their own AI-assisted output, and one documented case of catching and
fixing a subtly wrong AI suggestion.

Beginner skills (SQL, HTTP, JSON, git, shell) are treated as reading/
interpretation skills first, writing skills second — consistent with how
they're actually used day to day (nobody expects a junior to author an HTTP
protocol; they need to read a response and reason about a status code).

## 3. System architecture

**Base model:** small code-capable open model in the 3B–8B range —
Qwen2.5-Coder-7B-Instruct as the leading candidate (strong Python
performance, permissive license, runs on modest hardware). Qwen2.5-Coder-3B
or DeepSeek-Coder-6.7B as fallback options if a lighter footprint is needed.

**Retrieval layer** indexes:
- The 12-week project specs, rubrics, and skill checklists
- Existing code-tracing exercise library (40 JSON files) and Expert Jigsaw
  materials
- The calibrated bug bank (see §5)
- Reference solutions, marked "never surface directly" except as controlled
  bug-injection source material

**Mode selection** is done by the calling context, not left to the model to
infer.

## 4. Modes

### Critique / bug-injection mode (primary)
Takes a reference solution for the current week plus a difficulty target,
generates a plausible variant with a flaw calibrated to be findable-but-not-
obvious at that point in the curriculum. The student explains the bug in
prose; ByggaBat assesses whether the explanation identifies the real issue.
**This is the highest-risk, highest-value component — see §7.**

### Debug-companion mode
Given a failing test or error, does not supply the fix. Helps the student
form a hypothesis and design a check that would confirm or rule it out.
Socratic, but oriented around hypothesis-testing rather than withheld
syntax knowledge.

### Reviewer mode
Checks submitted code against the project rubric, weighted toward "would
this survive review" — style, edge cases, robustness — rather than just
"does this run." Fits Horizon LMS's no-grades, completion-based philosophy.

### Assistant / Q&A mode
Open concept questions. Secondary rather than central, since most student
time is spent reading/critiquing rather than authoring from scratch.

### Triage
Routes an incoming student message to the right mode. Candidate use case
for a typed-decision (Jev-style) classification pattern — see §6.

## 5. Data pipeline

1. Digitize the 12-week project sequence: objectives, required skills,
   common pitfalls, rubric items
2. Build a **calibrated bug bank** — seeded from real historical student
   mistakes where available, extended with LLM-generated synthetic variants,
   hand-graded for difficulty per week (only the instructor can calibrate
   this initially)
3. Build a **rubric for grading free-text bug explanations** — partial
   credit for identifying the symptom, full credit for the correct causal
   explanation
4. Write worked tutoring/debug-companion transcripts as few-shot exemplars
5. Tag all materials by week/skill so retrieval filters to what a student
   at that point should already know

## 6. Notes on LocalJev / Jev pattern

[LocalJev](https://github.com/githubnext/localjev) is a githubnext
experimental bridge implementing the Jev typed-decision protocol
(typed `choice` / `score` / `noul` questions over a context, returning
structured answers with probabilities) via DiffusionGemma on oMLX.

Assessment: not adopted as infrastructure — it's an early-stage prototype
tied to a Mac/oMLX/DiffusionGemma-specific stack, and its probabilities are
self-reported by the model rather than read from real logits (flagged by
the project itself as needing independent calibration). The **pattern** is
worth reusing — typed, scored questions instead of freeform grading prompts
— for reviewer-mode rubric checks and triage classification, implemented as
plain structured-output prompts against whichever model ByggaBat already
uses, rather than as a dependency on LocalJev itself.

## 7. Primary risk

Grading a student's free-text bug explanation reliably is a harder problem
than generating the bug itself. This should be prototyped and stress-tested
early, against real or newly collected student explanations, before
committing to critique/bug-injection as ByggaBat's centerpiece mode.

## 8. Deployment

Undecided — build and validate standalone first (faster iteration, no
coupling to the Horizon LMS release cycle), then integrate into Horizon LMS
once prompts/retrieval/grading are proven. Horizon LMS integration adds
real value to reviewer mode specifically, since it can pull the actual
rubric for the actual project and track completion in place.

## 9. Build phases

1. Digitize project specs/rubrics; build the calibrated bug bank for the
   first 2–3 projects, hand-graded for difficulty
2. Stand up base model + RAG; build bug-injection mode first; small eval
   set testing whether it produces bugs a target-week student should catch
   but not trivially
3. Build the free-text critique-grading rubric; test against real or newly
   collected student explanations
4. Add debug-companion and reviewer modes; pilot standalone with a small
   cohort
5. Horizon LMS integration
