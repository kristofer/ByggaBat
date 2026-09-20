# RATIONALE.md — Why This Curriculum Is Built This Way

## Preface

`DESIGN.md` and `WEEK1-4-NOTIONS.md` describe *what* was decided. This
document gathers the *why* — the reasoning, theory, and worked examples
that produced those decisions — in one place, so a design choice three
months from now can be checked against its original argument instead of
re-derived from memory.

It also reads, deliberately, like the outline of something longer: a case
for teaching programming differently now that an agent can write correct
syntax faster than any student can. Part VIII names that possibility
directly rather than leaving it implicit.

---

## Part I — What agents actually change

### 1.1 The one-minute problem

Handed a typical week-3 assignment — "read a CSV of ratings, compute the
average, print it formatted to two decimal places" — a frontier LLM
produces a correct answer in under a minute, unprompted for hints, with no
struggle at all. This isn't a hypothetical to plan around; it's already
true today, for every assignment in this curriculum's first eight weeks.

The natural first reaction is to treat this as a threat to route around —
lock down tools, forbid AI use, make assignments AI-proof. That's a losing
strategy against something this capable, and it also misdiagnoses the
actual problem.

### 1.2 What was being tested, and what wasn't

Writing a loop from a blank page was never the actual goal of a beginner
assignment — it was a *proxy* for something else: has this student
internalized enough of the language to reason about it correctly? The
proxy and the target used to move together, because there was no way to
produce correct syntax without that internal understanding. An agent
breaks that correlation completely — correct output no longer implies
understanding. The target skill hasn't disappeared. The proxy has just
stopped measuring it, and needs replacing.

### 1.3 Where real differentiation always lived

What separates a senior engineer from a junior one was never typing speed
or syntax recall — it's reading unfamiliar code fast, forming a hypothesis
about a bug and testing it, knowing when something *looks* right but
isn't, decomposing a vague problem into independently verifiable pieces,
and reviewing someone else's work critically. These skills were always the
real bottleneck to seniority; agents just moved the point where they're
needed from "years three through five of a career" to "the first time you
ask an assistant to write something for you." A curriculum that still
treats authorship-from-scratch as the finish line is optimizing for a
skill that stopped being scarce.

---

## Part II — Why reading can't be taught first: the pedagogy argument

### 2.1 Decoding before comprehension

Reading-instruction research draws a sharp line between *decoding* (turning
symbols into words) and *comprehension* (understanding what those words
mean). Comprehension-focused instruction only works once decoding is
automatic — a reader who has to sound out every word has no working memory
left over for meaning. Skip straight to "read critically" before decoding
is fluent, and the result isn't critical reading; it's guessing.

### 2.2 The code analogue: the notional machine

Code has the same split. The **notional machine** — CS-education
terminology (du Boulay, Sorva) for the internal model a programmer holds of
what the language is actually doing — plays the role decoding plays in
reading. Without one, "review this AI-generated function critically" isn't
a judgment call a student can make; there's no internal simulator to run
the code against and compare. The student can only pattern-match against
code that looks similar to something they've seen before, which is
exactly the failure mode this curriculum is trying to avoid.

### 2.3 The synthesis: front-load, then shrink — don't eliminate

This is why writing practice isn't cut, just reweighted and moved earlier:

| Weeks | Writing : Reading+Debugging+Review |
|---|---|
| 1–2 | ~70:30 |
| 3–4 | ~50:50 |
| 5–6 | ~30:70 |
| 7–8 | ~20:80 |
| 9–12 | reading/review-dominant |

Concretely, a week-1 exercise still asks a student to trace this by hand,
value by value, before ever running it:

```python
total = 0
for r in ratings:
    total += r
avg = total / len(ratings)
```

A week-9 exercise, by contrast, hands the student a working
AI-generated version of the same idea with a deliberately planted flaw and
asks for a written explanation of what's wrong and why — the decoding is
assumed complete; what's being tested is comprehension under scrutiny.

---

## Part III — The notional machine, formally

### 3.1 Why it must never contradict itself

Two pieces of theory justify the shape this takes, not just the existence
of "teach a mental model" as an idea. Bruner's spiral curriculum: revisit
the same object at increasing resolution rather than teaching topics once
and moving on. Sweller's cognitive load theory: a novice lacks the schema
to absorb a whole system at once, so each pass must add detail without
requiring the student to unlearn anything from the previous pass — a
*correction* costs more working memory than a *first encounter* would
have.

### 3.2 Weeks 1–4 — the load-bearing model

The proposition that carries the most weight, by far:

```python
>>> ratings = [4.5, 3.0, 5.0, 4.0]
>>> backup = ratings
>>> backup.append(100)
>>> ratings
[4.5, 3.0, 5.0, 4.0, 100]
```

**Names point at objects; they don't contain them.** `backup = ratings`
doesn't copy the list — it points a second name at the same one. Nearly
every beginner surprise in weeks 5–8 (mutable default arguments, "why did
my function change something it wasn't supposed to," `is` vs. `==`) is
this same idea, revisited — not new content to memorize each time.

### 3.3 Weeks 5–8 — refining, not correcting

The flagship illustration of what "refine without contradiction" looks
like in practice:

```python
def add_review(review, reviews=[]):
    reviews.append(review)
    return reviews
```

A default argument is created *once*, when `def` runs — not fresh per
call. Combined with "names point at objects" (§3.2), this bug is no longer
a gotcha to memorize; it's a direct, derivable consequence of a rule the
student already has. That's the test of whether a depth-spiral pass
actually worked: the harder case should feel *inevitable* once stated, not
like a new rule bolted onto the old one.

### 3.4 Weeks 9–12 — recombination as evidence of mastery

By this level, several "new" topics turn out not to be new mechanisms at
all:

```python
def decorator(f):
    def wrapper(*args, **kwargs):
        print(f"calling {f.__name__}")
        return f(*args, **kwargs)
    return wrapper

@decorator
def greet(name):
    return f"hello {name}"
```

`@decorator` above `def greet` desugars to `greet = decorator(greet)` — a
function call (proposition 4), whose return value repoints a name
(proposition 2). A student who can state that decomposition unprompted has
demonstrated a *stronger* mental model than one who memorized decorator
syntax as its own topic — recognizing recombination is itself the skill
being assessed at this level, not a bonus insight.

---

## Part IV — The Spirals mechanism, and where it's already working

### 4.1 Depth, breadth, and macro spirals

Depth spiral: one object, increasing resolution pass by pass. Breadth
spiral: the same shape, applied to a different domain, at roughly constant
resolution (a rep, not a lesson). Macro spiral: a big idea revisited much
later, at a harder case, spaced out across the curriculum rather than
taught once and left behind.

### 4.2 MiniYelp as the fixed object

Section 1 of the level-100 Python course already runs this mechanism, not
just in theory but in the actual narration. Video 1 teaches "a name, not a
box" and explicitly plants a hook: *"Hold onto that; it's going to matter
again, in a more surprising way, once we get to lists you can actually
change in place."* Video 3 cashes that hook in directly: *"Back in Video 1,
two names pointing at the same number moved independently... Lists behave
differently, and this is the promised payoff for holding onto that earlier
lesson."* That's a depth-spiral pass, verbatim, happening in production
course material — not an abstraction DESIGN.md imposed after the fact.

### 4.3 Block syntax as its own pattern

The same mechanism extends past mutability. Adding read-only text I/O
right after Video 2 —

```python
with open("reviews.txt") as f:
    for line in f:
        print(line.strip())
```

— gives students a fourth block sharing `if`/`while`/`for`'s exact `:` +
indented-suite shape, but the first one that *doesn't* redirect the
instruction pointer. Seen cheaply now, "block syntax ≠ control flow" is
half-absorbed well before it needs full explanation as the context-manager
protocol in weeks 5–8.

---

## Part V — Designing exercises an agent can't shortcut

### 5.1 The one-shot litmus test

For any new exercise: if a frontier model can solve it in one shot with
zero explanation, and that constitutes a complete, gradeable answer, it
isn't an exercise for this curriculum anymore. It's either raw material
for reading practice, or it needs a critique/explain layer added before
it's usable.

### 5.2 Mine, don't author

Two kinds of material scale better sourced from the real world than
invented: exercises pulled from permissively-licensed standard-library
source and small open-source projects (exposing real quality variation —
good idioms next to mediocre ones, which is exactly the discernment a
working programmer needs), and bugs mined from real bug-fix commit diffs
rather than synthetic flaws, which gives a realistic bug distribution and
free gold-standard explanations (the commit message often *is* the correct
answer).

### 5.3 Requiring justification, not just output

An exercise grounded in a real file, with a required trace or test
alongside any code, can't be discharged by pasting a correct answer —
correctness alone earns no credit if it isn't paired with evidence of
understanding.

---

## Part VI — Assessing understanding, not output

### 6.1 Two catalogues, not one

The bug bank targets surface-level code defects. A separate **misconception
catalogue** targets conceptual errors — mutable default arguments,
aliasing, late-binding closures — each entry tagged by which notional-machine
proposition it violates, not just by week or difficulty.

### 6.2 Calibration with a principled basis

A bug is "findable but not obvious" once the depth-spiral pass that
explains it has actually been taught — not an arbitrary per-week guess.
Tagging bugs by proposition is what makes that calibration checkable
rather than a matter of instructor intuition each time.

### 6.3 "Explain the model" checkpoints

A qualitative, never percentage-scored artifact, matching the Spirals
checkpoint definition: given `backup = ratings; backup.append(100)`, sketch
or describe what `ratings` and `backup` point at afterward. This tests the
model directly, rather than testing whether the student can match expected
output — which an agent can already do without understanding anything.

### 6.4 The open risk

Grading a free-text explanation reliably is a harder problem than
generating the bug that provoked it. This is the central unresolved risk
in the whole design and should be prototyped early, against real student
explanations, before anything here is treated as settled.

---

## Part VII — What ByggaBat is, given all of this

Each mode operationalizes one part of the argument above, rather than
existing as an independent feature:

| Mode | Answers to |
|---|---|
| Critique / bug-injection | Part V (agent-proof exercises) + Part VI (misconception-tagged assessment) |
| Debug-companion | Part II (build the simulator before reviewing) + Part IV (mental-model diagrams) |
| Reviewer | Part I.3 (review as the real differentiator) |
| Assistant / Q&A | Secondary by design — Part I.2, since answering isn't the scarce skill anymore |
| Triage | Infrastructure only; not an argument of its own |

---

## Part VIII — Where this could become a book

Read end to end, Parts I–VII sketch an argument that isn't really specific
to ZCW or to Python: agentic AI has made authorship cheap and made
judgment — reading, debugging, reviewing, deciding what to build — the
actual scarce skill, and CS education hasn't caught up to that shift.
A working title along the lines of *Thinking in Python After Agents* isn't
far from what's already written here.

What a book-length treatment would need beyond this curriculum-specific
document: empirical validation that the reweighted writing:reading ratio
actually produces stronger outcomes than the traditional sequence, not
just a theoretically motivated one; comparison against how this plays out
in other languages and other domains, since the notional-machine content
here is Python-specific but the underlying argument (Part I–II) isn't;
some account of what changes for the instructor, not just the curriculum —
grading free-text explanations at scale (Part VI.4) is a new skill for
teaching staff too; and a harder look at where this argument might be
wrong — cases where authorship-first still matters, or where an agent's
presence undermines the notional-machine-building phase itself rather than
just changing what comes after it.

None of that needs resolving to keep building ByggaBat and this
curriculum. But it's worth naming honestly: what's collected here reads
like a first draft of the argument, not just a design document.
