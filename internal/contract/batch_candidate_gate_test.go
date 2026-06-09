package contract_test

import (
	"testing"

	"portaljuridico/internal/batchcandidategates"
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
	for _, entry := range records {
		record := entry.Record
		if len(record.SelectedUniqueIntentIDs) < 3 {
			t.Fatalf("%s selected intents=%d, want at least 3", record.BatchID, len(record.SelectedUniqueIntentIDs))
		}
		if record.RenderAllowed || record.SitemapAllowed || record.PublicationAllowed || record.PublicPath != "" {
			t.Fatalf("%s escaped blocked gate: render=%t sitemap=%t publication=%t public_path=%q", record.BatchID, record.RenderAllowed, record.SitemapAllowed, record.PublicationAllowed, record.PublicPath)
		}
		if record.BaseURLMode != "official_configured" || !record.OfficialURLLocked {
			t.Fatalf("%s must track locked official URL while preserving blocked publication, mode=%q locked=%t", record.BatchID, record.BaseURLMode, record.OfficialURLLocked)
		}
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
