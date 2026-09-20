# ByggaBat — Design Document

v6

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

### Prerequisite: the Developer's Workbench

Python Course 1 assumes a shorter prerequisite course, the Developer's
Workbench (20 videos), is already complete. It's the first concepts ledger
in the whole program, and every later course — Python, Java, and Go alike
— assumes everything in it is already in hand and never re-logs it. It
covers, in order: shell (filesystem, redirection, piping, scripts, exit
codes, `$PATH`), git (repository through branching, merge conflicts read
as "git asking a person to decide," `diff`, `stash`, pull requests), HTTP
and JSON (request/response, verbs, status codes, idempotency, all five
JSON types, nesting, exact grammar) via real `curl` traffic rather than
abstraction, then regex, tool installation (Homebrew, with an explicit
"machine-wide tool vs. project dependency" distinction), and Markdown.

This matters directly to how four of the five "beginner" auxiliary skills
below are scoped: shell, git, HTTP, and JSON are **not introduced by the
Python curriculum at all** — students already hold the concepts. Python's
job for those four is syntax-transfer onto knowledge already in place
("you already know what a `POST` does; here's how you send one from
Python"), not concept-teaching. **SQL is the one exception** — nothing in
the Workbench covers it, so Python Course 1 is where SQL is genuinely
taught from first principles, and it should be scoped with the same care
as any other new topic, unlike the other four.

## 2. Skill progression (redrafted)

The ratio of hands-on writing to reading/debugging/review shifts steadily
across the 12 weeks:

| Weeks | Focus | Writing : Reading+Debugging+Review |
|---|---|---|
| 1–2 | Build the notional machine — core syntax, trace-before-run | ~70:30 |
| 3–4 | Transition — functions/OOP/file I/O, debugging formalized, `json` module as the payoff for a Workbench hook (§3) | ~50:50 |
| 5–6 | Reading/debugging dominant — SQL introduced fresh (no prerequisite coverage), HTTP applied via Python syntax onto concepts already known, debug-by-hypothesis formalized, first code review exercises | ~30:70 |
| 7–8 | Decomposition and verification — vague-problem decomposition, tests used to check AI-generated code, git workflow applied to a real project (not taught fresh — see §3) | ~20:80 |
| 9–10 | Unfamiliar codebases — modify code they didn't write, subtler bug-injection challenges, peer review as core activity | reading/review-dominant |
| 11–12 | Capstone — built with LLM assistance throughout, but every stage requires artifacts proving comprehension | reading/review-dominant |

Capstone requirements specifically include: a written decomposition plan
before coding, tests that verify LLM-generated code, a code-review writeup
of their own AI-assisted output, and one documented case of catching and
fixing a subtly wrong AI suggestion.

Beginner skills (SQL, HTTP, JSON, git, shell): SQL is treated as reading
skill first, writing second, exactly as originally planned, since it has
no prerequisite coverage. HTTP, JSON, git, and shell are different — the
Workbench already brought all four to working fluency before Python
starts, so Python's role is applying that existing fluency in Python
syntax, not building it from scratch.

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

### Python topic outline

Language topics sequenced against the week bands from §2, tagged **W**
(write-list — must author from scratch) or **R** (read-list — must
recognize and not be confused by; authorship not required). Given the
Developer's Workbench prerequisite (§1), shell and git appear in this
table only where Python-specific syntax needs teaching (module imports,
`subprocess`, packaging) — the underlying concepts aren't introduced here.
SQL is the one auxiliary skill genuinely new to this curriculum.

| Weeks | Core language | Data structures | OOP / structure | Errors & I/O |
|---|---|---|---|---|
| 1–2 | Variables & names, types (int/float/bool/str/None), operators, `if`/`elif`/`else`, `while`, `for` **[W]** | Lists, indexing/slicing basics **[W]** | — | `print`/`input`, reading a traceback **[R→W]**; reading a text file (`with`, `for line in f`) **[R]** |
| 3–4 | Functions (def, params, return, defaults), scope basics (local vs. global) **[W]** | Dicts, tuples, sets; nested structures **[W]** | — | `try`/`except` basics **[W]**; `open()`/file read-write **[W]**; `json.loads`/`json.dumps` **[W]** — syntax only, concepts already known from the Workbench (§1) |
| 5–6 | f-strings, comprehensions (read first, then write), ternary expressions **[R→W]** | List/dict comprehensions **[W]**; `collections` module (`Counter`, `defaultdict`) **[R]** | Classes: `__init__`, attributes, methods **[W]** | Custom exceptions, `finally` **[W]**; context managers (`with`) **[R→W]**; SQL — genuinely new, read first **[R]** |
| 7–8 | `*args`/`**kwargs` **[R]**; lambda **[R]**; closures **[R]**; `re` module **[R]** — pattern concepts already known from the Workbench | Instance vs. class attributes, basic inheritance **[W]** | Modules & packages, imports, `if __name__ == "__main__"` **[W]** | Logging basics **[R]**; `assert` / test-writing **[W]**; `requests` library **[R]** — HTTP concepts already known |
| 9–10 | Decorators **[R]**; generators & `yield` **[R]**; magic/dunder methods beyond `__init__`/`__str__` **[R]** | Multiple inheritance / mixins **[R]** | Package structure across multiple files **[W]** | `pytest` basics **[W]**; virtual environments, `requirements.txt` **[W]** — pip's version of the Homebrew tool-vs-dependency distinction already known |
| 11–12 | Type hints **[R]**; `async`/`await` (recognize only) **[R]** | — | Reading unfamiliar module structure **[R]** | Reading production-style logging & error handling **[R]** |

### The notional machine, by week band

Three resolution levels, each strictly additive — nothing at a later level
contradicts an earlier one, per the depth-spiral rule. A real tool worth
using or embedding for the visualization capability noted below:
[Python Tutor](https://pythontutor.com) already renders almost exactly this
model (frames stacked on a call stack, objects in a separate heap, arrows
instead of copied values) at every level described here.

#### Weeks 1–4 — the coarse, load-bearing model

1. **A single instruction pointer moves through the code** — one line at a
   time, top to bottom, except where `if`/`elif`/`else` branches it or
   `while`/`for` sends it backward. No concurrency, no lookahead.
2. **Names point at objects; they don't contain them.** `x = [1, 2, 3]`
   creates a list somewhere; `x` points at it. `y = x` points a second name
   at the *same* object, not a copy. This is the single proposition that
   makes list aliasing, "why did my function change something it wasn't
   supposed to," and `is` vs. `==` all comprehensible later without
   contradiction.
3. **Objects are coarsely one of two kinds: mutable, and immutable.**
   Lists/dicts can be changed in place; numbers/strings/tuples can't — any
   "change" to a string produces a new object and repoints the name. Full
   mutability semantics (shallow vs. deep copy, sharing across calls) wait
   until weeks 5–8.
4. **A function call opens a fresh, temporary workspace.** A call creates a
   new frame — a fresh set of name-to-object bindings — stacked on top of
   the caller. `return` hands a value back and the frame disappears; its
   names stop existing, but any object they pointed at survives if
   something outside still points at it too.
5. **An error immediately halts and unwinds.** Execution stops at that line
   and control jumps outward until something catches it, or nothing does.
   A traceback at this resolution reads as "where was the pointer when
   everything stopped," not yet a full call-stack diagnostic.

**Read-only text I/O, introduced early, as a second block-syntax instance.**
`with open("reviews.txt") as f: for line in f: ...` fits as soon as
proposition 1 exists — it needs `for` and strings, nothing else. Its value
here isn't the file access itself; it's that `with` gives students a fourth
block (alongside `if`/`while`/`for`) that shares the identical `:` +
indented-suite shape but, unlike the other three, does not redirect the
instruction pointer — a useful early counter-example that "block syntax"
and "control flow" aren't the same idea, seen cheaply now rather than
sprung as a surprise once `with`'s real mechanics (proposition 13) open in
weeks 5–8. `for line in f` is also a second iteration domain (files, not
just lists) at negligible cost — a small early rep toward the general
iteration protocol (proposition 17, weeks 9–12).

*Deliberately closed at this level:* generators/laziness, decorators, full
LEGB scope beyond local-vs-global, closures, multiple inheritance, deep
copy, comprehension internals, the context-manager protocol behind `with`
(proposition 13) — `with` is used here, not yet explained.

#### Weeks 5–8 — opening specific boxes

6. **A frame survives as long as something still points at it.** Sharpens
   proposition 4: a closure is the case where a nested function keeps a
   reference back to its enclosing frame's names even after that frame's
   call has returned — why a closure can still read a variable from a
   function that already "finished." (Read-list.)
7. **Mutable default arguments are created once, at definition time.**
   Sharpens 4 without contradicting it: the call still gets fresh
   *bindings*, but a mutable default object is one single object created
   when `def` ran, silently shared across every call that doesn't pass its
   own — the canonical mutable-default bug, now derivable from 2 + 4
   rather than memorized.
8. **A comprehension opens its own small frame too** — same mechanics as
   4, just unnamed and one-shot.
9. **`*args`/`**kwargs` are parameter binding, not new syntax** — extra
   positionals collected into a tuple, extra keywords into a dict, both
   just names bound in the frame. **Lambda** is an unnamed function frame,
   same call mechanics as `def`.
10. **An object is its own private dict of names plus a pointer to its
    class.** Instantiation creates a fresh object with an empty dict, then
    calls `__init__` as an ordinary call (4) whose first parameter, `self`,
    is a name pointing at that new object (2). `self.x = 5` is proposition
    2 again, in a new location.
11. **`obj.attr` checks the object's own dict first, then the class's
    dict.** The coarse version of attribute lookup — what instance-vs-class
    attributes *are*: two dicts checked in order, not a special rule.
12. **Raising an error creates an exception object; unwinding passes that
    object outward looking for a type match.** Each `except` is a type
    check performed as the object passes; `finally` runs unconditionally.
13. **`with` is a special case of the unwind path.** Objects with
    `__enter__`/`__exit__` hooks: `with` calls `__enter__` going in and
    guarantees `__exit__` runs going out — including mid-unwind — why
    `with open(...)` reliably closes on failure too.
14. **A module runs top-to-bottom once on import (1) and becomes an object
    whose dict holds everything it defined.** `import` is proposition 2
    again — a new name pointing at an object that already exists.

*Deliberately closed at this level:* decorators, generators/`yield`,
`nonlocal`/multi-level enclosing scope, multiple inheritance/MRO beyond the
two-dict lookup, metaclasses, async.

#### Weeks 9–12 — recombination, not new mechanism, in most cases

The shift here: several of these aren't new machinery at all, just prior
propositions composed — recognizing that is itself evidence of a strong
mental model, versus memorizing each as its own topic.

15. **A decorator is proposition 4 + proposition 2, composed.**
    `@decorator` above `def f(): ...` desugars to `f = decorator(f)`: a
    call (4), whose return value repoints the name `f` (2).
16. **A generator's frame can be suspended, not just destroyed** —
    refines 6: calling a generator function creates a generator object
    holding a paused frame; each `next()` resumes exactly where `yield`
    last paused, until the frame runs off the end and raises
    `StopIteration`. This is mechanically why a generator can only be
    consumed once.
17. **`for` calls a fixed protocol, not built-in magic.** It calls
    `iter()`, then `next()` each pass, catching `StopIteration`. Lists,
    generators, and custom `__iter__`/`__next__` classes all go through
    this one mechanism (11, applied consistently).
18. **Multiple inheritance extends proposition 11's chain from one link to
    a resolved list (MRO)** — same walk, just longer and order-computed
    rather than assumed.
19. **Dunder methods are attribute lookup (11) triggered by syntax instead
    of a dot.** `a + b` looks up `__add__`; `len(x)` looks up `__len__`;
    operators and built-ins are sugar for "look up a name, then call it"
    (11 + 4).
20. **`nonlocal` makes the enclosing-scope walk explicit — LEGB is one
    chain-walk, same shape as attribute lookup (11).** Check the local
    frame, then each enclosing frame still alive (6), then the module dict
    (14), then built-ins.
21. **`async`/`await` — read-only: a suspension point like `yield` (16),
    resumed by something other than the direct caller.** Students need to
    recognize "many paused frames, something else schedules their turns,"
    not author a concurrent system.
22. **Package structure is proposition 14 applied across files** — each
    file its own module-dict; reading unfamiliar multi-file code is mostly
    tracing which name in which file's dict points at which object
    elsewhere.
23. **Environments and `import` resolution are a genuinely separate black
    box** — not part of the execution model at all. `sys.path` and
    installed packages decide what a module name resolves to *before* any
    frame/object machinery runs. This has a direct pre-existing scaffold:
    the Developer's Workbench already teaches the Homebrew distinction —
    "what does this machine have" vs. "what does this project depend on"
    — so `pip`/`venv` should be introduced explicitly as that same
    distinction's Python-specific instance, not a new concept. Worth
    naming the box regardless, since students otherwise fold it into
    "more Python syntax" and get confused when identical code behaves
    differently across environments.
24. **Type hints are metadata, not enforcement** — stated as an explicit
    negative, since the natural wrong assumption (especially from prior
    exposure to a statically-typed language) is that they behave like one.
25. **Custom exception hierarchies and structured logging add no new
    mechanics** — propositions 12 and 18 at production scale, plus
    ordinary calls (4) that happen to write to a log.

*Deliberately closed entirely:* metaclasses, the full descriptor protocol —
flagged as real and worth naming, past the defined "intermediate" ceiling,
not something this curriculum owes a pass on.

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
