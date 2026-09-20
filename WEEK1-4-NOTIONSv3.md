# WEEK1-4-NOTIONS.md — Notional Machine × Video Spiral, Section 1

Companion to `DESIGN.md` §3. This document maps the weeks 1–4 notional
machine (propositions 1–5) onto Section 1 of the level-100 Python course
("Variables, Control Flow, Containers, Functions" — Videos 1–4 plus the
matching textbook chapter), to show that these aren't two separate design
ideas sitting next to each other — the video sequence is already, in
practice, a working delivery mechanism for the notional machine, built
around one fixed running example, MiniYelp.

## 1. The two ideas, briefly restated

**The notional machine (DESIGN.md §3, weeks 1–4):**

1. A single instruction pointer moves through the code.
2. Names point at objects; they don't contain them.
3. Objects are coarsely one of two kinds: mutable, and immutable.
4. A function call opens a fresh, temporary workspace (frame).
5. An error immediately halts and unwinds.

**The Spirals depth-spiral idea, applied here:** object stays fixed, one
system, resolution increases pass by pass, each pass opens a black box
without contradicting what the student already believes. Section 1 uses
MiniYelp itself as that fixed object. At this stage it's barely an "app" —
a handful of ratings and businesses — but it's the same MiniYelp later
sections extend with classes, persistence, and a real server. The videos
don't introduce five unrelated topics; they add resolution to one running
example, one proposition at a time.

## 2. Where they already integrate, video by video

### Video 1 — Variables, arithmetic, string formatting
**Delta:** variables, arithmetic operators, f-strings.
**Opens:** proposition 2, explicitly — "the better picture: `ratings` is a
name, and right now that name points at this list of four numbers" (Slide
2), reinforced through reassignment (Slide 4: "the name can point at
something new at any time; the values themselves never change") and two
independent names (Slide 5).
**Plants a hook:** Slide 5 ends with "Hold onto that; it's going to matter
again, in a more surprising way, once we get to lists you can actually
change in place" — a deliberate macro-spiral seed for Video 3, not
incidental repetition.
**Touches:** proposition 5, lightly — the `ZeroDivisionError` crash at
Slide 14 is shown and deliberately left unfixed, explicitly deferred to
Video 2.

### Video 2 — Control flow: if/elif/else, for, while
**Opens:** proposition 1 — the whole video's subject is the instruction
pointer branching (`if`) and looping (`for`/`while`).
**Closes a loop:** resolves the Video 1 `ZeroDivisionError` hook by
guarding the division with `if`.

### Video 3 — Lists and dictionaries as containers
**Opens:** proposition 3 — this is the payoff video, and it's the clearest
example of a genuine depth-spiral pass in the whole section. Slide 6
explicitly reopens Video 1's material at higher resolution: "Back in Video
1, two names pointing at the same *number* moved independently... Lists
behave differently, and this is the promised payoff for holding onto that
earlier lesson." Proposition 2 isn't retaught — it's the premise the new
distinction (mutable vs. immutable) is built on top of.
**Reinforces:** proposition 5 across three containers at once —
`IndexError`, `KeyError`, and `ZeroDivisionError` are named together in the
recap as "the same underlying story," an explicit macro-thread on "Python
refuses to guess" rather than three unrelated vocabulary words.
**Structural note:** slicing (Slide 5) is introduced immediately before
aliasing (Slide 6) specifically as a contrast case — a copy, not a shared
reference — so the aliasing distinction lands sharper two slides later.

### Video 4 — Functions, input/output
**Opens:** proposition 4 — taught operationally rather than named
outright. "Same function, two completely different lists" and the
`return`-vs-`print` distinction are exactly proposition 4 in practice:
each call gets fresh bindings, and only `return` lets a value survive past
the frame's end. The composing-functions example
(`format_rating(highest(ratings))`) is the coarse call-stack case — one
call nested inside another — without needing the word "frame" at all.
**Names the actual goal:** the textbook chapter's summary states the
target explicitly — "recognizing these patterns on sight... is the actual
skill" — which is the same thing DESIGN.md §3 identifies as the point of
notional-machine teaching generally, not incidental phrasing.

## 3. What makes this a genuine integration, not parallel content

- MiniYelp is the one fixed object; each video adds resolution to its data
  model (a bare list → aliased vs. copied lists → a nested dict-of-lists →
  functions operating on it) without contradicting anything shown earlier.
  That's a depth spiral, not a topic list with a shared theme.
- Propositions aren't introduced as abstract rules — the narration teaches
  by revisiting a *named moment* ("back in Video 1..."), which is exactly
  the mechanism the misconception catalogue and "explain the model"
  checkpoints (DESIGN.md §3) depend on: a student already has a shared
  reference point to anchor an explanation to (e.g. "this is the
  `backup = ratings` situation") instead of restating the rule from
  scratch each time.
- Errors are taught as one recurring idea (proposition 5) surfacing in
  three unrelated-looking contexts, not three separate vocabulary words —
  the macro-spiral mechanism (same big idea, harder case, spaced out)
  operating at a much smaller scale than the six-big-idea table in the
  Spirals research, but the identical shape.

## 4. Gaps and recommendations for extending the pattern

- Propositions 1 and 5 are both touched but deliberately left open at
  week-4 resolution — correct, since neither is meant to close yet. Worth
  making that *open* status explicit so weeks 5–8 material (which reopens
  proposition 5 into typed exception objects, proposition 12 in DESIGN.md
  §3) can cite Video 3's recap ("three errors now firmly in your
  vocabulary") as its entry state, rather than re-teaching from zero.
- No video yet stages a deliberate "predict, then run" moment purely for
  proposition 4 (frames) the way Video 3's Slide 6 does for proposition 3.
  Video 4 teaches `return` vs. `print` operationally but never has the
  student trace two overlapping calls by hand before seeing the answer.
  Worth a check-yourself question in that shape for Section 1, matching
  the pattern already used for aliasing.
- Recommend Section 2 (classes) open with the same rhetorical device Video
  3 used — an explicit "back in Video 4..." callback — reopening
  proposition 4 into proposition 10 (an object is a dict plus a class
  pointer, and `__init__` is an ordinary call) rather than introducing
  objects as new vocabulary.

## 5. Adding simple text I/O

Read-only text I/O fits right after Video 2 — it needs only `for` and
strings, both already in place, and nothing from Videos 3–4:

```python
with open("reviews.txt") as f:
    for line in f:
        print(line.strip())
```

**Why this reinforces `if`/`while`/`for`, not just adds a topic:**

- **Block syntax as its own recognizable pattern.** `if`, `while`, `for`,
  and now `with` all share the identical shape — a line ending in `:`,
  followed by an indented suite. Three instances make that shape feel
  arbitrary; four make it a rule a student can predict forward onto a
  construct they haven't met yet.
- **A useful counter-example, seen early and cheaply.** Every block met so
  far redirects the instruction pointer (proposition 1) — `if` branches
  it, `while`/`for` send it backward. `with` is the first block that
  *doesn't* — it wraps the pointer's path without branching or looping it.
  Meeting that distinction now, while the other three are still fresh,
  means "block syntax ≠ control flow" is half-absorbed before it needs
  full explanation as proposition 13 in weeks 5–8, rather than arriving as
  a surprise then.
- **A second `for`-shape at negligible cost.** `for line in f` iterates
  over something that isn't a list — a small breadth-spiral rep (same
  shape, different domain), in the same spirit as Video 1's
  donations-vs-ratings device — and quietly starts building toward the
  general iteration protocol (proposition 17, weeks 9–12) without needing
  to say anything about it yet.

**What stays deferred:** full write-mode I/O stays in weeks 3–4 as already
scoped in DESIGN.md's topic table. The actual mechanics of `with`
(the context-manager protocol, proposition 13) stay closed until weeks
5–8 — `with` is *used* here, on schedule with the rest of the weeks 1–2
resolution level, not *explained* yet.

**One error deliberately not surfaced here:** `open()` on a missing file
raises `FileNotFoundError`, which would be a clean fourth instance of the
proposition-5 thread Video 3 already built (`IndexError`/`KeyError`/
`ZeroDivisionError`, recapped as "the same underlying story"). Better left
out for now — Video 3's recap closes cleanly as "three errors," and a
fourth so soon dilutes that. It's available whenever weeks 5–8 wants an
opener for reopening proposition 5 into typed exception objects.

**Checkpoint wording needs updating.** Section 1's current checkpoint
language ("No class, no file I/O yet — this section is deliberately
syntax-only," from `_index.md`) should narrow to "no file *writing* yet,"
since reading is now in scope right after Video 2.

## 6. The prerequisite Workbench, and the `json` hook it sets up

Entry state for this whole section, from `_index.md`: "Developer's
Workbench complete." That course's own capstone (Videos 19–20) deliberately
extracts JSON data with `grep -o`, framed explicitly as "the hard way, on
purpose — and the first honest look at why a real parser is worth wanting
later." That's a direct, already-authored setup this section should cash
in rather than let go stale.

Since JSON's concepts (five types, nesting, exact grammar) are already
fully known going into Python, `json.loads`/`json.dumps` don't need to
wait for a natural "we need this" moment the way a fresh concept would —
they're pure syntax-transfer onto something already understood, and the
motivating hook is already sitting there. Right after Video 3 (dicts) is
the natural placement: a JSON object maps directly onto a dict, so showing
`json.loads('{"name": "Joe'\''s Pizza", "rating": 4.5}')` immediately after
dicts are introduced lets the mapping do the explaining, and can open with
a direct callback: "Remember the Workbench capstone, doing this with
`grep`? Here's the tool that was worth wanting."

This also means the earlier placement suggestion in DESIGN.md's topic
table (§3) — `json.loads`/`json.dumps` at weeks 3–4 — should be read as
*this* specific moment, not a loose week-level suggestion.

## 7. Implication for ByggaBat

- Retrieval and the bug bank / misconception catalogue (DESIGN.md §6–7)
  should tag material by *which video/slide first taught the relevant
  proposition*, not just by week number — so a hint in debug-companion
  mode can say "like the `backup = ratings` case from Video 3" instead of
  re-explaining aliasing generically. This is a stronger, cheaper version
  of the "ground exercises in something real" principle from §7 — grounded
  in material the student has already built a mental model around, instead
  of a fresh unfamiliar example each time.
- MiniYelp should be the default source for any weeks 1–4 bug-injection or
  "explain the model" checkpoint content ByggaBat generates, rather than a
  fresh example — reusing the fixed object the depth spiral already
  established, instead of resetting the resolution to zero every time a
  new exercise is generated.
- Where the misconception catalogue lists "list aliasing" as an entry, its
  worked exemplar should be, literally, the `backup = ratings` /
  `backup.append(100)` example from Video 3 — already the canonical
  instance in this student's mental model, not a new one to build fresh.
- The same video-tagging principle should extend one course earlier: a
  hint touching HTTP, JSON, git, or shell can cite the Developer's
  Workbench video that already taught it ("like the merge conflict from
  Workbench Video 8") instead of re-explaining the concept — since, per
  §6, those four are genuinely prior knowledge here, not new material.
