package contract_test

import (
	"testing"

	"portaljuridico/internal/batchcandidateexpansion"
	"portaljuridico/internal/batchcandidategates"
	"portaljuridico/internal/paidintent"
)

func TestBatchCandidateGatesSelectArchiveDraftsWithoutPublishing(t *testing.T) {
	report := batchcandidategates.Validate(".")
	if !report.Passed() {
		t.Fatalf("batch candidate gates failed contract: %v", report.Messages())
	}

	records, loadReport := batchcandidategates.LoadRecords(".")
	if !loadReport.Passed() {
		t.Fatalf("could not load batch candidate gates: %v", loadReport.Messages())
	}
	if len(records) != 6 {
		t.Fatalf("candidate gates=%d, want 6 batch families", len(records))
	}
	readinessRecords, readinessReport := batchcandidateexpansion.LoadRecords(".")
	if !readinessReport.Passed() {
		t.Fatalf("could not load expansion readiness: %v", readinessReport.Messages())
	}
	readinessByBatch := make(map[string]batchcandidateexpansion.Record)
	for _, entry := range readinessRecords {
		readinessByBatch[entry.Record.BatchID] = entry.Record
	}
	paidRecords, paidReport := paidintent.LoadRecords(".")
	if !paidReport.Passed() {
		t.Fatalf("could not load paid intent gates: %v", paidReport.Messages())
	}
	paidStatusByIntent := make(map[string]string)
	for _, entry := range paidRecords {
		paidStatusByIntent[entry.Record.UniqueIntentID] = entry.Record.PaidIntentStatus
	}

	totalSelected := 0
	for _, entry := range records {
		record := entry.Record
		readiness, ok := readinessByBatch[record.BatchID]
		if !ok {
			t.Fatalf("%s missing expansion readiness", record.BatchID)
		}
		expected := make([]string, 0)
		for _, intentID := range readiness.ExpansionCandidateIntentIDs {
			if paidStatusByIntent[intentID] == paidintent.PassedBlockedStatus {
				expected = append(expected, intentID)
			}
		}
		if !sameStringSet(record.SelectedUniqueIntentIDs, expected) {
			t.Fatalf("%s selected intents do not match paid-passed readiness candidates: selected=%d expected=%d", record.BatchID, len(record.SelectedUniqueIntentIDs), len(expected))
		}
		if len(record.SelectedUniqueIntentIDs) < 18 {
			t.Fatalf("%s selected intents=%d, want at least 18 paid-passed expansion candidates", record.BatchID, len(record.SelectedUniqueIntentIDs))
		}
		totalSelected += len(record.SelectedUniqueIntentIDs)
		if record.RenderAllowed || record.SitemapAllowed || record.PublicationAllowed || record.PublicPath != "" {
			t.Fatalf("%s escaped blocked gate: render=%t sitemap=%t publication=%t public_path=%q", record.BatchID, record.RenderAllowed, record.SitemapAllowed, record.PublicationAllowed, record.PublicPath)
		}
		if record.BaseURLMode != "official_configured" || !record.OfficialURLLocked {
			t.Fatalf("%s must track locked official URL while preserving blocked publication, mode=%q locked=%t", record.BatchID, record.BaseURLMode, record.OfficialURLLocked)
		}
	}
	if totalSelected != 168 {
		t.Fatalf("selected intents=%d, want 168 paid-passed expansion candidates", totalSelected)
	}
}

func TestBatchCandidateGateRejectsPublicOrUnknownArchiveIntent(t *testing.T) {
	record := batchcandidategates.Record{
		GateID:                  "candidate-saude",
		BatchID:                 "batch-saude-suplementar-digital",
		GateStatus:              "blocked_batch_candidate",
		SourceArchivePath:       "data/editorial/batch_draft_expansion_archive.jsonl",
		ArchiveMinimumRecords:   100,
		SelectedUniqueIntentIDs: []string{"intent-inexistente"},
		CandidatePathPrefix:     "/temas/",
		BaseURLMode:             "lab_placeholder",
		OfficialURLLocked:       false,
		MaxSimilarityAllowed:    0.64,
		MaxSimilarityObserved:   0.64,
		MinimumHumanScore:       88,
		CTAContextRequired:      true,
		SourceURLAuditRequired:  true,
		PublicationBlockReason:  "",
		RenderAllowed:           true,
		SitemapAllowed:          true,
		PublicationAllowed:      true,
		PublicPath:              "/temas/intent-inexistente/",
		CheckedAt:               "2026-06-09",
	}

	report := batchcandidategates.ValidateRecordAgainstArchive(record, batchcandidategates.ArchiveIndex{
		BaseURLMode:       "official_configured",
		OfficialURLLocked: true,
	})
	for _, code := range []string{
		"batch_candidate_selected_intent_missing",
		"batch_candidate_render_allowed",
		"batch_candidate_sitemap_allowed",
		"batch_candidate_publication_allowed",
		"batch_candidate_has_public_path",
		"batch_candidate_missing_block_reason",
		"batch_candidate_url_config_mismatch",
	} {
		if !report.HasIssue(code) {
			t.Fatalf("missing issue %q in %v", code, report.Codes())
		}
	}
}

func sameStringSet(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	seen := make(map[string]int)
	for _, value := range left {
		seen[value]++
	}
	for _, value := range right {
		seen[value]--
	}
	for _, count := range seen {
		if count != 0 {
			return false
		}
	}
	return true
}
