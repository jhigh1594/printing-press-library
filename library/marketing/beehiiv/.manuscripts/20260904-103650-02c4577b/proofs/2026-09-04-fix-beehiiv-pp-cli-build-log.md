Manifest transcendence rows: 7 planned, 7 built. Phase 3 complete: all 7 ship. (Plus absorbed row 7 growth-summary rebuilt as an 8th hand-code unit.)

## What was built
- 7 novel insights commands (store-computed, read-only, pp:data-source computed):
  subscriber-sources, churn-sources, subscriber-lookup (typed exit 0,3), post-performance,
  send-times, referral-health, compare-publications.
- Absorbed: insights growth-summary rebuilt on the enlarged store (podcasts/exports schema now synced).
- Registration hook file internal/cli/insights_store.go registers the insights group via
  registerNovelCommand; shared drain-first store helpers with NULL-safe COALESCE scans and
  missing-mirror guards (empty JSON + sync hint).
- Long-description scope redirects carried from the manifest into Cobra Long fields.

## Deferred / notes
- Per-publication attribution: subscriptions sync without publication tags. compare-publications
  reports store totals plus an attribution note; per-pub exactness needs per-pub mirrors.
- send-times rates slots only when the mirror carries expand=stats; note emitted otherwise.
- No API-key available: verify/dogfood run in no-auth mode (401 expected paths skip live checks).

## Skipped body fields
- None; the 16 new spec endpoints emit typed commands from the merged spec.
