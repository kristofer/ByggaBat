# WEEK5-8-NOTIONS.md — Notional Machine × Video Spiral, Sections 2–3

Companion to `DESIGN.md` §3 and `WEEK1-4-NOTIONS.md`.

**Note on status:** `WEEK1-4-NOTIONS.md` mapped propositions 1–5 onto
video/textbook material that already exists (Videos 1–4, uploaded and
reviewed directly). No equivalent scripts for weeks 5–8 have been supplied
yet. This document is therefore **prescriptive** — a design brief for how
that material should integrate the notional machine and continue the
MiniYelp spiral — not a retrospective mapping. It should be checked
against the actual scripts once they exist, the way Section 1's mapping
was checked against Videos 1–4.

## 1. Propositions in scope (from DESIGN.md §3, weeks 5–8)

6. A frame survives as long as something still points at it (closures).
7. Mutable default arguments are created once, at definition time.
8. A comprehension opens its own small frame.
9. `*args`/`**kwargs` are parameter binding; lambda is an unnamed frame.
10. An object is its own private dict of names plus a class pointer.
11. `obj.attr` checks the object's own dict first, then the class's.
12. Raising an error creates an exception object; unwinding looks for a
    type match.
13. `with` is a special case of the unwind path (`__enter__`/`__exit__`).
14. A module runs once on import and becomes an object with its own dict.

## 2. Continuing MiniYelp as the fixed object

Weeks 1–4 left MiniYelp as loose variables and a nested dict-of-lists —
correct, but with no behavior attached. Section 2 should not introduce
classes as a new topic; it should **re-describe data MiniYelp already
has**, exactly as `WEEK1-4-NOTIONS.md` §4 recommended: open with a "back
in Video 4..." callback, reopening proposition 4 (fresh frame per call)
into proposition 10 (object = dict + class pointer, `__init__` is an
ordinary call whose first parameter, `self`, is a name like any other).

```python
class Business:
    def __init__(self, name):
        self.name = name
        self.reviews = []

    def add_review(self, rating):
        self.reviews.append(rating)

    def average_rating(self):
        if not self.reviews:
            return None
        return sum(self.reviews) / len(self.reviews)
```

This is the same accumulator pattern from Video 1 (`total`, built up one
item at a time) and the same empty-list guard from Video 3 (`not
ratings`), now living inside a method instead of a bare script. Nothing
about the *logic* is new — only where it lives. That's the point: students
should recognize the pattern before they're told it's "the same as
before."

## 2a. What's genuinely new here vs. applied here

Given the Developer's Workbench prerequisite (DESIGN.md §1), the auxiliary
skills scoped for weeks 5–8 are not on equal footing:

- **SQL (weeks 5–6) is genuinely new** — no Workbench coverage exists.
  It should be paced and scaffolded like any first-encounter topic,
  read-first as already planned, with the same care given to proposition
  1–5-level content in weeks 1–4.
- **HTTP (via `requests`, weeks 7–8) is not new** — the Workbench already
  covers request/response, verbs, status codes, and idempotency via real
  `curl` traffic. This section's job is showing what that already-known
  behavior looks like from inside Python, not teaching what a `POST` does.
- **Regex (`re` module, weeks 7–8) is not new either** — the Workbench
  teaches pattern-matching concepts directly, including the explicit
  discipline "a deliberately simple pattern beats a technically complete
  one." Introduce `re` with a direct callback to that framing rather than
  re-deriving it.
- **Git (weeks 7–8 PR workflow) is not new** — full branch/merge/PR
  workflow is already Workbench material. What's new is doing it inside a
  real, growing Python project rather than a standalone exercise.

## 3. Proposed integration, section by section

### Section 2 (weeks 5–6) — classes, comprehensions, exceptions, `with`

**Classes open proposition 10–11.** `Business` and `Review` as two
classes, MiniYelp's dict-of-businesses now holding `Business` objects
instead of raw dicts. Instance vs. class attributes (proposition 11) get a
natural, slightly dangerous example worth using deliberately:

```python
class Business:
    all_businesses = []   # class attribute — one shared list

    def __init__(self, name):
        self.name = name          # instance attribute — one per object
        self.reviews = []
        Business.all_businesses.append(self)
```

This is a genuine parallel to proposition 7's mutable-default-argument
bug — a mutable class attribute is a single shared object, exactly the way
a mutable default argument is a single object created once at `def` time.
Teaching them back-to-back, explicitly named as the same shape, is a
strong candidate for a depth-spiral pass: proposition 7 reopened, not a
new gotcha.

**Comprehensions (proposition 8) get a real MiniYelp use**, not a toy list:

```python
high_rated = [b.name for b in Business.all_businesses if b.average_rating() and b.average_rating() >= 4.0]
```

**Custom exceptions (proposition 12) cash in two deferred hooks at once.**
`WEEK1-4-NOTIONS.md` §5 deliberately held back `FileNotFoundError` from
Video 3's error recap, flagging it as "available whenever weeks 5–8 wants
an opener." This is that opener — reopen proposition 5 ("an error
immediately halts and unwinds") into proposition 12 (errors are typed
objects unwinding through a type-matching search) using
`FileNotFoundError` as the bridge, then introduce a MiniYelp-specific
custom exception built on the same mechanism:

```python
class BusinessNotFoundError(Exception):
    pass

def get_business(name):
    for b in Business.all_businesses:
        if b.name == name:
            return b
    raise BusinessNotFoundError(f"no business named {name!r}")
```

**`with` gets its real mechanics (proposition 13) — the other deferred
hook.** `with open(...)` has been *used* since weeks 1–2 without
explanation, flagged explicitly for this moment. Now show what makes it
work:

```python
class LoggedLookup:
    def __enter__(self):
        print("starting lookup")
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        print("lookup finished")
```

The payoff line, delivered directly: "`with open(...)` has worked since
week 1 because file objects implement exactly this pair of methods — now
you know why."

### Section 3 (weeks 7–8) — modules, closures, logging, testing

**Modules and packages open proposition 14.** Split MiniYelp into
`business.py`, `review.py`, `main.py` — each file a module object with its
own dict, `import business` a plain proposition-2 example (a new name
pointing at something that already exists) applied to a module instead of
a list.

**Closures and lambda (propositions 6, 9) earn their keep via sorting**,
rather than being demonstrated in isolation:

```python
top_businesses = sorted(Business.all_businesses, key=lambda b: b.average_rating() or 0, reverse=True)
```

A closure example worth pairing with it — a rating threshold "baked into"
a filter function, closing over a variable from the enclosing scope — is
the natural moment to reopen proposition 6 (a frame survives as long as
something points at it) with a concrete payoff instead of an abstract
definition.

**Logging replaces `print` in MiniYelp's own functions**, framed
explicitly as a `return`-vs-`print` echo from Video 4: `print` was already
taught as "for humans watching," logging is the production-grade version
of the same idea, not a new concept.

**`assert` and first pytest tests** should test the exact functions
already built — `test_average_rating_empty()`, `test_business_not_found()`
— so testing arrives attached to code students already trust, not a
fresh, unfamiliar codebase.

## 4. Misconception catalogue additions this section motivates

- Mutable class attribute treated as if it were per-instance (direct
  sibling of the mutable-default-argument entry, same underlying cause)
- Catching `Exception` too broadly, silently swallowing
  `BusinessNotFoundError` along with everything else
- Forgetting `self` doesn't crash immediately but produces confusing
  `TypeError`s later — worth its own catalogue entry given how far the
  actual error can appear from the actual mistake
- Using `with` without understanding it still requires the object to
  define `__enter__`/`__exit__` — assuming any object works

## 5. Handoff to weeks 9–12

By the end of Section 3, MiniYelp should be: a multi-file package, with
classes for `Business` and `Review`, custom exceptions, logging, and a
small pytest suite. That's the exact entry state `WEEK9-12-NOTIONS.md`
assumes — package structure (proposition 22) and testing infrastructure
already in place, ready to be extended rather than introduced from
nothing.
