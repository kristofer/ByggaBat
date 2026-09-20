# WEEK9-12-NOTIONS.md — Notional Machine × Video Spiral, Sections 4–5, and the Capstone

Companion to `DESIGN.md` §3, `WEEK1-4-NOTIONS.md`, and `WEEK5-8-NOTIONS.md`.

**Note on status:** prescriptive, as with `WEEK5-8-NOTIONS.md` — no scripts
for these sections exist yet to map against. Entry state assumed: MiniYelp
is a multi-file package with `Business`/`Review` classes, custom
exceptions, logging, and a small pytest suite (`WEEK5-8-NOTIONS.md` §5).

## 1. Propositions in scope (from DESIGN.md §3, weeks 9–12)

15. A decorator is proposition 4 + proposition 2, composed.
16. A generator's frame can be suspended, not just destroyed.
17. `for` calls a fixed protocol (`iter`/`next`/`StopIteration`).
18. Multiple inheritance extends attribute lookup to a resolved order.
19. Dunder methods are attribute lookup triggered by syntax.
20. `nonlocal`/LEGB is one chain-walk, same shape as attribute lookup.
21. `async`/`await` — read-only, a suspension point like `yield`.
22. Package structure is proposition 14 applied across files.
23. Environments/`import` resolution are a genuinely separate box.
24. Type hints are metadata, not enforcement.
25. Custom exception hierarchies/logging add no new mechanics.

## 2. Sections 4–5 (weeks 9–10) — recombination, made visible in MiniYelp

The pedagogical point of this level, restated from `DESIGN.md` §3:
several of these aren't new machinery, just prior propositions composed.
The integration goal is to make students *say that themselves*, using
MiniYelp as the shared reference, not to present each as a fresh topic.

**Decorators (proposition 15)** get a real, motivated use — timing or
logging a business lookup — introduced by explicitly desugaring it first:

```python
def logged(f):
    def wrapper(*args, **kwargs):
        print(f"calling {f.__name__}")
        return f(*args, **kwargs)
    return wrapper

@logged
def get_business(name):
    ...
```

Delivered as: "`@logged` above `get_business` is exactly
`get_business = logged(get_business)` — proposition 4, a function call,
followed by proposition 2, a name repointed. You already know both
halves."

**Generators (proposition 16)** get a genuine motivating problem rather
than an abstract `yield` demo — MiniYelp's review data has grown across
ten weeks; loading it all into a list at once is exactly the moment a
generator earns its place instead of being an arbitrary alternative
syntax:

```python
def stream_reviews(filepath):
    with open(filepath) as f:
        for line in f:
            yield parse_review(line)
```

Worth calling out explicitly: this is the same `with open(...)` /
`for line in f` shape introduced back in weeks 1–2 (`WEEK1-4-NOTIONS.md`
§5), now wrapped in a function that pauses instead of running straight
through — proposition 1 (a single instruction pointer) gets its first
real complication: the pointer can now be *parked* mid-function and
resumed later, exactly what proposition 16 describes.

**`for`'s protocol (proposition 17)** is worth demonstrating by having
`Business` (or a `ReviewCollection` wrapper) implement `__iter__`, so
`for review in business` works — tying dunder methods (19) and the
iteration protocol (17) together in one example instead of two.

**Dunder methods (proposition 19)** get sorting and display for free if
introduced together:

```python
class Business:
    def __lt__(self, other):
        return (self.average_rating() or 0) < (other.average_rating() or 0)

    def __len__(self):
        return len(self.reviews)
```

Delivered as: "`sorted(businesses)` and `len(business)` now work — not
because Python special-cased `Business`, but because `<` and `len()` are
just attribute lookups (proposition 11) triggered by syntax instead of a
dot, exactly like `+` calling `__add__`."

**Multiple inheritance / mixins (proposition 18)** get a real, small use —
a `Timestamped` mixin shared by `Business` and `Review` — rather than an
abstract diamond-inheritance example, keeping the MRO discussion attached
to something students already trust.

**Package structure (proposition 22) and virtual environments
(proposition 23)** land together, motivated by necessity rather than
introduced pre-emptively: MiniYelp needs its first *external* dependency
around here (a web framework, ahead of the capstone in §4 below), which is
the natural, honest moment to introduce `requirements.txt` and a venv —
proposition 23 is explicitly flagged in `DESIGN.md` §3 as a genuinely
separate box from execution semantics, and it should be taught as one:
"none of the frame/object model you've learned explains *which* code
`import` finds — that's a different job the interpreter does first."

## 3. Section 5 (weeks 11–12) — type hints, then the capstone

**Type hints (proposition 24)** work well as a retroactive exercise on
code students already trust, rather than a new-code exercise: annotate
`Business`'s existing methods, then deliberately show that Python doesn't
enforce them at runtime — pass a string where `average_rating() -> float`
is declared and watch it run anyway. The negative framing from `DESIGN.md`
§3 matters here specifically because some students will have prior
exposure to statically-typed languages and will assume otherwise.

**`async`/`await` (proposition 21)** stays read-only and small — cover it
only if the capstone's web framework choice (§4) is async-capable, framed
purely as "a suspension point like `yield` (proposition 16), resumed by
something other than the direct caller." Not authored.

## 4. The capstone: UI, Python middle end, database backend

This is the point the whole 12-week arc has been building toward, and it
should be framed to students exactly that way: not a new project, but
MiniYelp — the same running example since Video 1 — finally getting the
three pieces a real application needs. Three tiers:

- **UI** — the presentation layer the user actually touches
- **Python middle end** — the application/business-logic layer; this is
  where nearly everything built in weeks 1–10 already lives (the
  `Business`/`Review` classes, lookup functions, custom exceptions)
- **Database, in the backend** — persistent storage, replacing the
  in-memory dict-of-lists and the `reviews.txt` text file from weeks 1–2

**This is also where every "beginner" auxiliary skill (`DESIGN.md` §2)
converges — but given the Developer's Workbench prerequisite (`DESIGN.md`
§1), this needs a more precise claim than "first time write-list" for most
of them.** Four of the five skills below were already brought to working
fluency in the Workbench, before Python started; the capstone is where
Python syntax gets attached to concepts already held, not where those
concepts are learned. SQL is the one genuine exception, and its capstone
appearance is real first-time application of recently-learned material.
Worth stating to students precisely, since it changes what's actually
hard about the capstone — four of five converging skills are syntax-
transfer, not new semantics, which is a meaningfully lower load than "five
new things at once" would be:

| Skill | Concept status entering the capstone | What the capstone adds |
|---|---|---|
| SQL | Learned fresh, weeks 5–6 (`DESIGN.md` §3) | First real write-list use against a live database — genuinely new territory |
| HTTP | Fully known — Workbench Videos 10–12 (request/response, verbs, status codes, idempotency, via real `curl` traffic) | Python syntax for concepts already fluent: the middle end *is* an HTTP server, authored in Flask/FastAPI rather than just read about |
| JSON | Fully known — Workbench Videos 13–14; `json` module syntax already introduced weeks 3–4 (`WEEK1-4-NOTIONS.md` §6) | Using the already-familiar module as the live wire format between UI and middle end |
| git | Fully known — Workbench Videos 5–9 (branch, merge conflict, PR); applied to a real project since weeks 7–8 | Sustained use across a genuinely multi-week, higher-stakes project — scale, not new mechanics |
| shell | Fully known — Workbench Videos 1–4 | Running the server, managing the venv, database migrations — new commands, familiar shape |

**Recommended stack, given the "intermediate" ceiling `DESIGN.md` §3
already draws** (metaclasses and the full descriptor protocol explicitly
excluded there): SQLite for the database — no separate server to run or
administer, and it's the natural continuation of the `reviews.txt` file
already familiar since weeks 1–2. Flask or FastAPI for the Python middle
end — either keeps the "middle end" honestly small and Python-centered
rather than introducing a second full ecosystem. For the UI, server-rendered
HTML with HTMX (consistent with the stack `Horizon LMS` already uses) over
a full JavaScript SPA framework — keeps the UI layer inside the
"intermediate Python programmer" scope this curriculum targets, rather
than quietly requiring a second language track. This should be confirmed,
not assumed — flagged as an open decision in §6.

### Two genuinely new black boxes, not recombinations

Unlike most of weeks 9–12, these two pieces of the capstone are **not**
just prior propositions composed — worth being honest about that, the
same way `DESIGN.md` §3 was honest about proposition 23 being a genuinely
separate box rather than a refinement.

**Persistence outlives the process.** Nothing in propositions 1–25
explains what happens to state after the Python interpreter exits.
Every object, frame, and name discussed so far disappears when the
program ends — that was true and unremarked-on since proposition 4. A
database row is the first thing in this entire curriculum that survives a
restart. Worth naming directly: "everything you've learned about names and
objects describes what happens *while the program is running*. The
database is where MiniYelp keeps something that has to still be true
*after* it stops."

**A web server is many short executions, not one long one.** Proposition
1 ("a single instruction pointer moves through the code") was never
amended after Video 2 — it didn't need to be, until now. Each incoming
HTTP request starts a fresh pass through the middle end's code, with no
memory of the previous request except whatever was deliberately saved to
the database. This is worth an explicit callback: "back in Video 1, one
program ran start to finish. Here, the same function might run a thousand
times, for a thousand different visitors, sharing nothing with each other
except the database."

### Recommended framing device

Open the capstone module with the exact device that worked in Video 3
(`WEEK1-4-NOTIONS.md` §2) and is recommended again at the Section 2
opening (`WEEK5-8-NOTIONS.md` §2): an explicit "back in Video 1..."
callback to the original `reviews.txt` file, then to the in-memory
dict-of-lists, showing the same MiniYelp data taking a third, final shape.

## 5. Misconception catalogue additions this section motivates

- Assuming a decorator changes what a function *is* rather than what name
  points at it (confusing proposition 15 with something magical)
- Consuming a generator twice, expecting it to restart (direct
  proposition-16 consequence)
- Forgetting a database write needs to be committed — code that looks
  correct but doesn't persist, a new category of bug this curriculum
  hasn't had a vehicle for until now
- Treating each HTTP request as if it shared Python-level state with the
  last one (the "many short executions" black box, above)
- JSON serialization mismatches between what the middle end sends and
  what the UI expects — a realistic, common capstone-stage bug worth
  seeding into the bug bank

## 6. Open questions

- **Stack confirmation:** SQLite/Flask-or-FastAPI/HTMX is a recommendation
  here, not a decision — needs sign-off before it shapes exercise content.
- **Scope of the UI tier:** how much frontend work is "intermediate
  Python programmer" scope versus a distinct skill this curriculum isn't
  targeting — worth an explicit boundary statement, the same way
  metaclasses were explicitly excluded in `DESIGN.md` §3.
- **Whether persistence deserves its own numbered proposition** (a
  proposition 26, addressing state outliving the process) formally added
  to `DESIGN.md` §3, given it's flagged here as a genuinely new black box
  rather than a recombination — worth deciding before ByggaBat's
  retrieval tagging (`DESIGN.md` §6–7) needs to reference it.
