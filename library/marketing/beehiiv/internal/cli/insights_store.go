// Copyright 2026 Kevin Magnan and contributors. Licensed under Apache-2.0. See LICENSE.
// Shared local-store plumbing for the insights command family. All commands in
// this family are computed from the synced SQLite mirror; none call the API.

package cli

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/mvanhorn/printing-press-library/library/marketing/beehiiv/internal/store"
)

func init() {
	registerNovelCommand(func(root *cobra.Command, flags *rootFlags) {
		addNovelCommandIfAbsent(root, newNovelInsightsCmd(flags))
	})
}

// insightsStore opens the read-only mirror, or reports a missing mirror with a
// sync hint and an empty JSON envelope. ok=false means the caller should
// return nil (output already written).
func insightsStore(cmd *cobra.Command, flags *rootFlags, dbPath string) (*store.Store, func(), bool) {
	if dbPath == "" {
		dbPath = defaultDBPath("beehiiv-pp-cli")
	}
	if _, statErr := os.Stat(dbPath); os.IsNotExist(statErr) {
		fmt.Fprintf(cmd.ErrOrStderr(), "no local mirror at %s\nrun: beehiiv-pp-cli sync --resources publications,subscriptions,posts --db %s\n", dbPath, dbPath)
		if !wantsHumanTable(cmd.OutOrStdout(), flags) {
			_ = printJSONFiltered(cmd.OutOrStdout(), map[string]any{"rows": []any{}, "note": "no local mirror; run sync first"}, flags)
		}
		return nil, func() {}, false
	}
	db, err := store.OpenReadOnly(dbPath)
	if err != nil {
		return nil, func() {}, false
	}
	hintIfUnsynced(cmd, db, "")
	return db, func() { _ = db.Close() }, true
}

// scanRows drains a result set into (id, data) pairs before any follow-up
// queries. SQLite single-connection semantics forbid nested queries.
func scanRows(ctx context.Context, db *store.Store, query string, args ...any) ([]insightsRow, error) {
	rows, err := db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []insightsRow{}
	for rows.Next() {
		var r insightsRow
		if err := rows.Scan(&r.ID, &r.Data); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

type insightsRow struct {
	ID   string
	Data string
}

func (r insightsRow) Map() map[string]any {
	m := map[string]any{}
	_ = json.Unmarshal([]byte(r.Data), &m)
	return m
}

// subscriptionsRows scans the typed subscriptions table (fast path).
type subRow struct {
	ID            string
	Email         string
	Status        string
	Tier          string
	ReferralCode  string
	ReferringSite string
	UTMSource     string
	UTMMedium     string
	UTMChannel    string
	UTMCampaign   string
	CreatedUnix   sql.NullInt64
}

func scanSubscriptions(ctx context.Context, db *store.Store) ([]subRow, error) {
	rows, err := db.DB().QueryContext(ctx, `SELECT id,
			COALESCE(email,''), COALESCE(status,''), COALESCE(subscription_tier,''),
			COALESCE(referral_code,''), COALESCE(referring_site,''),
			COALESCE(utm_source,''), COALESCE(utm_medium,''), COALESCE(utm_channel,''), COALESCE(utm_campaign,''),
			created
		FROM subscriptions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []subRow{}
	for rows.Next() {
		var r subRow
		if err := rows.Scan(&r.ID, &r.Email, &r.Status, &r.Tier, &r.ReferralCode, &r.ReferringSite, &r.UTMSource, &r.UTMMedium, &r.UTMChannel, &r.UTMCampaign, &r.CreatedUnix); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

type countPair struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func topCounts(counts map[string]int, limit int) []countPair {
	pairs := make([]countPair, 0, len(counts))
	for k, v := range counts {
		pairs = append(pairs, countPair{Name: k, Count: v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Name < pairs[j].Name
	})
	if limit > 0 && len(pairs) > limit {
		pairs = pairs[:limit]
	}
	return pairs
}

func countNonEmpty(values []string) map[string]int {
	m := map[string]int{}
	for _, v := range values {
		if v == "" {
			m["(unknown)"]++
			continue
		}
		m[v]++
	}
	return m
}

func syncedPublications(ctx context.Context, db *store.Store) []map[string]any {
	pubTypes := []string{"publications"}
	placeholders := make([]string, 0, len(pubTypes))
	args := make([]any, 0, len(pubTypes))
	for _, t := range pubTypes {
		placeholders = append(placeholders, "?")
		args = append(args, t)
	}
	rows, err := scanRows(ctx, db, fmt.Sprintf(`SELECT id, data FROM resources WHERE resource_type IN (%s)`, strings.Join(placeholders, ",")), args...)
	if err != nil {
		return []map[string]any{}
	}
	out := []map[string]any{}
	for _, r := range rows {
		m := r.Map()
		m["id"] = r.ID
		out = append(out, m)
	}
	return out
}

func unixSlot(sec int64) (weekday string, hour int, ok bool) {
	if sec <= 0 {
		return "", 0, false
	}
	t := time.Unix(sec, 0).UTC()
	return t.Weekday().String(), t.Hour(), true
}

// optionalArg returns the first positional (publicationId) or "" when absent.
func optionalArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}
