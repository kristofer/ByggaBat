# ByggaBat — Design Document
### V3

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

## 3. Mental models (the notional machine)

Writing fluency was already reframed in §1 as a means to build the
**notional machine** — the internal model a programmer holds of what the
language is actually doing — rather than an end in itself. This section
makes that thread explicit and gives it its own supporting structure,
drawing on the depth-spiral mechanism described in ZCW's broader Spirals
curriculum research (object stays fixed, resolution increases pass over
pass, each pass opens a black box without contradicting what the student
already believes).

Applied to Python specifically, the "object" held fixed across a depth
spiral doesn't have to be an external system (a search engine, a banking
app) — it can be one of Python's own semantic black boxes, each revisited
at increasing resolution as the 12 weeks progress rather than taught once
and left behind:

- **What a name refers to** — binding vs. copying, `b = a` vs.
  `b = a.copy()`, list aliasing
- **Mutability** — mutable default arguments, shared state via aliasing,
  shallow vs. deep copy
- **Scope and closures** — LEGB resolution, late-binding closures over
  loop variables
- **Execution / call stack** — what a traceback actually represents,
  frame by frame
- **Iteration and laziness** — generators, exhaustion, why a generator can
  only be consumed once
- **Identity vs. equality** — `is` vs. `==`, small-int/string interning
  artifacts

Each is its own small depth-spiral thread: a coarse-but-correct model
introduced early, then reopened at higher resolution once the surrounding
concepts exist to make the next layer make sense (e.g. late-binding
closures only click once closures themselves are already a working idea).
Three or four of these run as an explicit macro thread reappearing at
increasing resolution around weeks ~2, ~5, and ~8, rather than each being
covered once in isolation.

### Misconception catalogue

A sibling structure to the calibrated bug bank (§6), but targeting
conceptual errors rather than surface-level code bugs: mutable default
arguments, aliasing, late-binding closures, shallow/deep copy confusion,
generator exhaustion, and similar well-documented Python misconceptions.
Sourced the same way as the rest of the corpus (§7) — mined from real,
well-documented "Python gotchas" collections rather than invented from
scratch.

### Mental-model-tagged bug injection

Bugs in the critique/bug-injection mode (§5) get tagged by *which
misconception they're designed to surface*, not just by difficulty. This
gives the "findable but not obvious" calibration a principled basis:
findable once the relevant depth-spiral pass has been taught for that
concept, rather than an arbitrary per-week guess.

### "Explain the model" checkpoints

A qualitative, never percentage-scored artifact: the student describes or
sketches what's happening in memory or execution for a given snippet —
what do `a`, `b`, and the list they point to look like after this code
runs? This tests whether the underlying mental model is correct, not just
whether the student can match expected output.

### Visualization in debug-companion mode

Python's reference and object model is invisible by default, which is a
large part of why it's misunderstood. Debug-companion mode (§5) should be
able to render a simple box-and-pointer or stack-frame diagram alongside
its Socratic hints, not just text — a real capability addition worth
scoping separately from the base text-generation pipeline.

## 4. System architecture

**Base model:** small code-capable open model in the 3B–8B range —
Qwen2.5-Coder-7B-Instruct as the leading candidate (strong Python
performance, permissive license, runs on modest hardware). Qwen2.5-Coder-3B
or DeepSeek-Coder-6.7B as fallback options if a lighter footprint is needed.

**Retrieval layer** indexes:
- The 12-week project specs, rubrics, and skill checklists
- Existing code-tracing exercise library (40 JSON files) and Expert Jigsaw
  materials
- The calibrated bug bank and misconception catalogue (see §3, §6)
- Reference solutions, marked "never surface directly" except as controlled
  bug-injection source material

**Mode selection** is done by the calling context, not left to the model to
infer.

## 5. Modes

### Critique / bug-injection mode (primary)
Takes a reference solution for the current week plus a difficulty target,
generates a plausible variant with a flaw calibrated to be findable-but-not-
obvious at that point in the curriculum. The student explains the bug in
prose; ByggaBat assesses whether the explanation identifies the real issue.
**This is the highest-risk, highest-value component — see §9.**

### Debug-companion mode
Given a failing test or error, does not supply the fix. Helps the student
form a hypothesis and design a check that would confirm or rule it out.
Socratic, but oriented around hypothesis-testing rather than withheld
syntax knowledge. Can render mental-model diagrams alongside hints (§3).

### Reviewer mode
Checks submitted code against the project rubric, weighted toward "would
this survive review" — style, edge cases, robustness — rather than just
"does this run." Fits Horizon LMS's no-grades, completion-based philosophy.

### Assistant / Q&A mode
Open concept questions. Secondary rather than central, since most student
time is spent reading/critiquing rather than authoring from scratch.

### Triage
Routes an incoming student message to the right mode. Candidate use case
for a typed-decision (Jev-style) classification pattern — see §8.

## 6. Data pipeline

1. Digitize the 12-week project sequence: objectives, required skills,
   common pitfalls, rubric items
2. Build a **calibrated bug bank** — seeded from real historical student
   mistakes where available, extended with LLM-generated synthetic variants,
   hand-graded for difficulty per week (only the instructor can calibrate
   this initially)
3. Build a **misconception catalogue** (§3) alongside the bug bank, tagging
   entries by which mental-model gap they target
4. Build a **rubric for grading free-text bug explanations** — partial
   credit for identifying the symptom, full credit for the correct causal
   explanation
5. Write worked tutoring/debug-companion transcripts as few-shot exemplars
6. Tag all materials by week/skill so retrieval filters to what a student
   at that point should already know

## 7. Corpus sourcing, real bug-mining, and the one-shot litmus test

A frontier LLM can one-shot a typical week's assignment in under a minute.
Rather than fight this, it's treated as a signal about what shouldn't be
the deliverable anymore — and as a reason to source exercises from real,
existing Python at scale rather than hand-authoring synthetic toy problems.

### Widen the corpus

Pull exercise and reading material from real, permissively-licensed Python
that already exists, instead of authoring everything from scratch:

- **Standard library source** — `itertools`, `collections`, `pathlib`,
  `argparse` — well-written, readable, and exactly the idioms students
  will meet in real jobs
- **Small, well-scoped open-source libraries** (MIT/BSD/Apache licensed) —
  `requests`, `click`, contained modules from larger projects — for
  reading, tracing, and critique exercises
- **Public tutorial corpora and code datasets** (e.g. CodeSearchNet-style
  collections, freeCodeCamp's BSD-licensed curriculum) as a scalable,
  pre-tagged source of graded examples

This gives effectively unlimited, non-repeating material across cohorts,
and exposes students to real quality variation — good idioms next to
mediocre ones — which is the discernment a working programmer actually
needs and a synthetic toy problem can't teach.

### Mine real bugs instead of inventing them

Strengthens the calibrated bug bank (§6) by mining before/after diffs from
real bug-fix commits in open-source repos, rather than relying solely on
LLM-generated synthetic flaws:

- A realistic bug distribution — off-by-ones, wrong comparison operators,
  mutable-default-argument traps, scope errors — instead of whatever
  pattern a model decides to inject
- A much larger raw pool to hand-grade for difficulty than could be
  authored directly
- Free gold-standard answers for the free-text explanation rubric — the
  actual commit message and fix often *is* the correct explanation

Synthetic bug-injection from §5 still has a role, for producing fresh
variants calibrated to a specific week's target skill, but its patterns
should be modeled on what's mined from real fixes rather than invented
from nothing.

### Design exercises to resist one-shotting

Since LLMs can solve isolated, self-contained problems near-instantly, the
exercise itself is designed so an isolated solution isn't the deliverable:

- **Ground exercises in a real file or module**, not a floating problem
  statement — e.g. trace what happens when a specific real config-parsing
  function is called, so the deliverable is understanding *that* code, not
  producing a generic solution
- **Require justification alongside any working code** — tests written, an
  execution trace, a rationale for a design choice — so correct output
  alone earns no credit
- **One-shot litmus test for any new exercise:** if a frontier model can
  one-shot it with zero explanation and that constitutes a complete,
  gradeable answer, it isn't a ByggaBat exercise — it's either raw-material
  reading practice, or it needs a critique/explain layer added before use

### Pipeline addition

A corpus curation and licensing step precedes the existing data pipeline
(§6): selecting source repos, checking license compatibility for
educational reuse and redistribution to students, and building a tagging
taxonomy (comprehensions, decorators, context managers, generators, class
composition, error handling, common stdlib usage) so retrieval can pull
material at the right complexity and idiom for a given week.

## 8. Notes on LocalJev / Jev pattern

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

## 9. Primary risk

Grading a student's free-text bug explanation reliably is a harder problem
than generating the bug itself. This should be prototyped and stress-tested
early, against real or newly collected student explanations, before
committing to critique/bug-injection as ByggaBat's centerpiece mode.

## 10. Deployment

Undecided — build and validate standalone first (faster iteration, no
coupling to the Horizon LMS release cycle), then integrate into Horizon LMS
once prompts/retrieval/grading are proven. Horizon LMS integration adds
real value to reviewer mode specifically, since it can pull the actual
rubric for the actual project and track completion in place.

## 11. Build phases

1. Digitize project specs/rubrics; curate initial corpus sources and check
   licensing (§7); build the calibrated bug bank and misconception
   catalogue for the first 2–3 projects from mined real fixes plus
   hand-graded synthetic variants
2. Stand up base model + RAG; build bug-injection mode first; small eval
   set testing whether it produces bugs a target-week student should catch
   but not trivially
3. Build the free-text critique-grading rubric; test against real or newly
   collected student explanations
4. Add debug-companion and reviewer modes; pilot standalone with a small
   cohort
5. Horizon LMS integration
