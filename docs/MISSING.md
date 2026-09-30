# MISSING — deferred & future frontend features

Anything raised but intentionally **not** in V1 lands here instead of expanding scope.
When a new idea comes up mid-build that isn't already planned, add a row rather than
growing the current phase. Promote to a real plan when it's time.

## Deferred from V1 (decided during planning)

| Feature                                                            | Why deferred / context                                                                                              |
|--------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------|
| Phone-grade responsive (<768px) on dense surfaces                  | V1 is desktop-first, mobile-tolerable to ~768px. Diffs/tables/dashboard-grid reflow for phones is a separate pass.  |
| SSE endpoint for AI suggestions                                    | V1 AI is a batched single-request review-diff (no streaming). SSE only matters if/when streaming is reintroduced.   |
| AI-suggestion token *streaming* (inline live paint)                | Dropped deliberately for the cheaper batched review-diff pattern. Revisit if real-time feel is wanted.              |

> **Shipped since:** lab-mutating manager operations (service start/stop/restart
> and config push) were delivered under #236, governed by
> [ADR 0001](adr/0001-lab-mutating-operation-boundaries.md),
> [ADR 0002](adr/0002-start-stop-lab-mutating-operations.md) and
> [ADR 0003](adr/0003-config-push-lab-mutating-operation.md).

## Suggested-later (raised in build, not yet planned)

| Feature                                                    | Why deferred / context                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
|------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `http.Flusher` on `middleware.responseWriter`              | Logging middleware wraps `http.ResponseWriter` but only implements `Hijacker` (added for WebSocket). `Flusher` is not needed yet — no endpoints stream responses. Add `Flush()` method (type-assert underlying writer, delegate) when SSE, file export streaming, or AI token streaming is added.                                                                                                                                                                                                                            |
