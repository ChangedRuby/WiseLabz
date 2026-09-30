# Product

## Register

product

## Users

Technical homelabbers — people who run their own infrastructure at home or in a
small lab: Proxmox VE, Docker/Portainer, pfSense/OPNsense, self-hosted apps.
They are comfortable with APIs, YAML config, Docker Compose, and the terminal.
Their context: managing a fleet of services that drift over time, where the
documentation is always stale and nobody enjoys writing it. They reach for
WiseLabz to stop hand-maintaining docs and to get an honest, live picture of
what changed since the last sync.

Primary jobs to be done:
- See live service status and what drifted at a glance.
- Read accurate, auto-generated documentation for each service.
- Review and accept/reject diffs the sync engine surfaces.
- Configure connectors, templates, auth, and AI providers without fighting the UI.

## Product Purpose

WiseLabz automatically generates and maintains documentation for a homelab. It
connects to each service's native API, fetches live data, renders structured docs
from configurable templates, and runs a diff engine that detects change since the
last sync — then either drafts an AI-assisted section update or raises an alert so
nothing drifts undocumented. The web dashboard shows live service data alongside
the generated docs.

Documentation is the heart of the product, but not the whole of it. WiseLabz is a
homelab *manager*, not only a documentation manager — the name is the thesis: it
makes your lab *wise*. The dashboard is a command surface, not a passive status
wall: you see live state and act on it, and everything you click routes you into
the docs that explain it. Docs stay the center of gravity; the manager layer is how
you reach and operate on them.

Success: the user trusts the docs without re-checking them, spends near-zero effort
keeping them current, notices drift through WiseLabz before it bites them, and runs
the lab from one place instead of jumping between service consoles.

## Product decisions (pre-planning, v1)

Locked before the planning session. These are product/UX scope calls; technical
enforcement and build details live in `docs/ARCHITECTURE.md`.

1. **IA spine — Docs is home; the dashboard is the entry/command surface.** The
   dashboard is a work-and-control surface (what drifted, what's stale, what needs
   review, what you can act on) that routes into docs. Never a big-number hero wall.
2. **Manager scope is deliberately narrow, and grows deliberately.** In scope:
   trigger sync (global + per-service), enable/disable a connector, add/remove
   connectors, and lab-mutating operations — service start/stop/restart and
   config push — each gated by step-up elevation and a confirmation/dry-run
   preview per ADRs 0001–0003 (`docs/adr/`).
3. **AI drafts are hybrid.** Minor section updates surface inline in the doc with a
   provenance marker (drafted / human-confirmed / synced-raw); structural changes
   route through the Changes accept/reject loop. Nothing implies more confidence
   than the sync actually has.
4. **Rack Panel is *the* identity, not one theme among many.** Industrial
   instrument panel — mono-forward, hairline grid, rounded-but-tight radii, a
   two-accent system (caution-orange for primary action/selection, cable-blue
   for secondary emphasis), status as shape + word. Dark-only; theme
   playground exists behind Settings → Theme for exploration but the default
   is the brand, not a menu of equal options.
5. **Motion is a meaning-bound differentiator.** Four moments earn animation: diff
   reveal, sync pulse (live heartbeat), drift surfacing, and state transitions
   (a row resolving on accept/reject). Everything else is plain state feedback; no
   page-load choreography.

**Resolved by ADR: how far "manager" goes past v1.** Lab-mutating control is
shipped (#236). The permission/confirmation model is defined in
[ADR 0001](docs/adr/0001-lab-mutating-operation-boundaries.md) (restart),
[ADR 0002](docs/adr/0002-start-stop-lab-mutating-operations.md) (start/stop) and
[ADR 0003](docs/adr/0003-config-push-lab-mutating-operation.md) (config push).

## Locked frontend direction (planning session, 2026-06; revised 2026-09)

Refines the decisions above; supersedes the "sharp corners / light-dark only"
specifics of decision 4. Shell architecture and the theme flag live in
`docs/ARCHITECTURE.md`.

1. **Navigation — bottom dock.** A floating, centered, soft-dark dock (solid, not
   glass); one amber pill slides under the active item. A sidebar variant was built
   and rejected in comparison. The dock is the single shell; no shell switching.
2. **Identity — Rack Panel palette + Big Shoulders Text / Martian Mono, in
   soft-dark chrome.** Caution-orange (primary) and cable-blue (secondary) on
   warm dark steel, rounded-but-tight radii (6–16px), shadow reserved for
   genuinely raised/floating surfaces (not every resting panel). Replaces the
   OKLCH neutral-ramp-plus-single-accent generator (and the
   Geist/Inter-Tight/Space-Grotesk font pool) that, despite being consciously
   chosen to avoid genericness, had converged on the same archetype as
   Linear/Vercel/Cursor-adjacent AI-coding tools. Status colors stay distinct
   (green/lemon/red), never collapsed into either accent. Section headers use
   sentence-case labels, never uppercase-tracked eyebrows.
3. **Theme engine ships behind Settings → Theme** (presets / fonts / OKLCH sliders)
   for exploration, defaulting to the locked Blueprint + Geist combo. It is a power
   surface, not a marketing theme playground; the default is the brand.
4. **Changes read as a JetBrains-IDE git diff.** A doc Change renders the whole
   document with context — add/remove lines highlighted, word-level intra-line
   highlighting, foldable unchanged regions, unified (default) or side-by-side.
   Acknowledge = accept, Dismiss = reject; the row resolves with motion (moment #4
   above) on either.

## Brand Personality

**Precise, technical, calm.** Confident expert tooling grounded in real
homelab hardware — rack panels, patch cables, nameplates — rather than in a
generic SaaS-dashboard register (Linear/Vercel-adjacent, now itself an
oversaturated "AI tool" look). Restrained, machine-honest, never shouting.
Voice is direct and
unpadded — it states facts (status, counts, timestamps, diffs) plainly and
trusts the user to read them. No marketing gloss, no fake enthusiasm, no
hand-holding. The interface should feel like a well-built instrument: dense
where density earns its keep, quiet everywhere else.

## Anti-references

- **Generic AI dashboard slop** — identical card grids, the big-number hero-metric
  template, gradient accents, tiny uppercase tracked eyebrows over every section.
  (The v1 prototype was rejected for exactly this.)
- **The neutral-ramp-plus-single-accent-hue OKLCH generator itself**, and the
  Geist/Inter/Space-Grotesk font pool — not just card-grid/eyebrow symptoms of
  genericness but the underlying formula that produces them across most
  current AI-assisted dashboard tools.
- **Heavy enterprise consoles (Grafana / Datadog)** — chart-soup, cluttered chrome,
  gray-on-gray density with no hierarchy, the legacy ops-console feel.
- **Bootstrap / Material defaults** — off-the-shelf component-kit look with no identity.

## Design Principles

1. **Machine-honest.** Show real state — exact counts, timestamps, statuses, diffs.
   Never decorate data or imply more confidence than the sync actually has.
2. **Density that earns its keep.** Pack information where the user is working
   (rosters, diffs, tables); stay quiet and spacious everywhere else. Density is a
   tool, not a default.
3. **Status is never color alone.** Every state reads through shape + word + color
   (dot + label), so it survives color blindness, grayscale, and a glance.
4. **Restraint over reflex.** No card unless a card is the best affordance; no
   accent unless it carries meaning. Warmth comes from typography and one rare
   highlight, not from chrome.
5. **Expert-default, keyboard-friendly.** Strong defaults, fast paths (command
   palette, shortcuts), and an UI that respects that the user already knows their lab.

## Accessibility & Inclusion

Target **WCAG 2.2 AA**, plus committed extras:
- Body text ≥4.5:1 contrast; large text ≥3:1; placeholders held to the same bar.
- Visible, non-color focus indicators on every interactive element.
- Status and severity never communicated by color alone — always dot/icon + word.
- Dark-first theme designed for long sessions; light theme held to the same contrast bar.

**Motion is a default, not an afterthought.** Animation and motion are a core
differentiator for WiseLabz — the interface should be astonishingly beautiful, and
motion carries that. Build it in from the start, not as garnish. Reduced motion is
**not** a blanket suppression rule; it is a **user-controlled setting** in
Settings → Appearance (Motion: Full / Reduced / Off). The OS `prefers-reduced-motion`
signal seeds the initial value but never silently forces motion off — the user decides
and can override. Default is full motion.
