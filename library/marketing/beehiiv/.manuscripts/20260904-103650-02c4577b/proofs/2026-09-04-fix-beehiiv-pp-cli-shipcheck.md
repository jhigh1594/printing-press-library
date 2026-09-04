# Shipcheck Report — beehiiv-pp-cli (reprint 2026-09-04)

## Final leg results
| Leg | Result |
|---|---|
| verify | PASS |
| validate-narrative | PASS |
| dogfood | PASS |
| workflow-verify | PASS |
| apify-audit | PASS |
| verify-skill | PASS |
| scorecard | HOLD (only: live_api_verification unverified — no API key on machine) |

Scorecard: 96/100 (Grade A). Sample output probe: 7/7 novel commands pass (publication_id echoed).

## Fix loops (2 used)
1. Loop 1: verify-skill positional-args findings + live-probe token misses shared one root cause:
   novel insights commands ignored the optional publicationId positional that rendered examples pass.
   Fixed signatures to `… [publicationId]`, echo publication_id in output.
2. Loop 2: README troubleshoot referenced a wrong exports path. Fixed at source (research.json),
   regenerated; implementations survived --force reconciliation. Cleared stale preserve snapshot.

## Enrichment outcomes (vs prior print scorecard 89/100)
- mcp_remote_transport 5 → 10 (stdio+http), mcp_token_efficiency 4 → omitted (thin orchestration),
  mcp_tool_design 5 → 10, mcp_surface_strategy 2 → 10 (Cloudflare pattern, 98 endpoints).
- cache_freshness 5 → 10 (cache block, 24h stale_after). auth_protocol → 10 (canonical BEEHIIV_API_KEY).
- Type fidelity restored on 16 newly documented endpoints (podcasts, exports, complimentary access,
  workspace permissions, post preview/test-send).

## Verdict: ship-pending-live-dogfood
All functional thresholds met. Only live_api_verification is unverified; needs BEEHIIV_API_KEY
for Phase 5 live dogfood. If the user provides a key, run live dogfood; else write
phase5-skip.json (auth_required_no_credential) and promote.
