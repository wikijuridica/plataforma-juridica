package contract_test

import (
	"testing"

	"portaljuridico/internal/batchcandidateexpansion"
)

func TestBatchCandidateExpansionReadinessUsesArchiveWithoutPublishing(t *testing.T) {
	report := batchcandidateexpansion.Validate(".")
	if !report.Passed() {
		t.Fatalf("batch candidate expansion readiness failed contract: %v", report.Messages())
	}

	records, loadReport := batchcandidateexpansion.LoadRecords(".")
	if !loadReport.Passed() {
		t.Fatalf("could not load batch candidate expansion readiness: %v", loadReport.Messages())
	}
	if len(records) != 6 {
		t.Fatalf("readiness records=%d, want one per batch family", len(records))
	}

	for _, entry := range records {
		record := entry.Record
		if record.ReadinessStatus != batchcandidateexpansion.PaidGateMissingStatus {
			t.Fatalf("line=%d status=%q", entry.Line, record.ReadinessStatus)
		}
		if record.PaidIntentGatePath != batchcandidateexpansion.PaidIntentGatePath {
			t.Fatalf("line=%d paid gate path=%q", entry.Line, record.PaidIntentGatePath)
		}
		if record.PaidIntentMissingCount <= 0 {
			t.Fatalf("line=%d expected missing paid-intent gates before expansion", entry.Line)
		}
		totalPaidClassified := record.PaidIntentPassedCount + record.PaidIntentBlockedCount + record.PaidIntentMissingCount
		if totalPaidClassified != len(record.ExpansionCandidateIntentIDs) {
			t.Fatalf("line=%d paid classified=%d intents=%d", entry.Line, totalPaidClassified, len(record.ExpansionCandidateIntentIDs))
		}
		if record.TargetCandidateCount < 30 {
			t.Fatalf("line=%d target=%d, want at least 30", entry.Line, record.TargetCandidateCount)
		}
		if len(record.ExpansionCandidateIntentIDs) < record.TargetCandidateCount {
			t.Fatalf("line=%d selected=%d target=%d", entry.Line, len(record.ExpansionCandidateIntentIDs), record.TargetCandidateCount)
		}
		if !record.PaidIntentRequired || !record.CTAContextRequired || !record.SourceURLAuditRequired || !record.SourceSpecificityRequired {
			t.Fatalf("line=%d missing required expansion guards", entry.Line)
		}
		if record.IndexPolicy != "noindex" || record.ManifestAllowed || record.RenderAllowed || record.SitemapAllowed || record.PublicationAllowed || record.PublicPath != "" {
			t.Fatalf("line=%d readiness escaped blocked contract: index=%q manifest=%t render=%t sitemap=%t publication=%t public_path=%q", entry.Line, record.IndexPolicy, record.ManifestAllowed, record.RenderAllowed, record.SitemapAllowed, record.PublicationAllowed, record.PublicPath)
		}
		if record.LegalArea == "" || record.TargetCandidateTier == "" || len(record.ActionableBlockers) == 0 {
			t.Fatalf("line=%d missing area, tier or actionable blockers", entry.Line)
		}
	}
}

func TestBatchCandidateExpansionReadinessRejectsWeakOrPublicExpansion(t *testing.T) {
	index := batchcandidateexpansion.ExpansionIndex{
		ArchiveByIntent: map[string]batchcandidateexpansion.ArchiveDraft{
			"familia-divorcio-consensual-filhos-bens-prova-documental": {
				BatchID:        "batch-familia-digital",
				UniqueIntentID: "familia-divorcio-consensual-filhos-bens-prova-documental",
				HumanScore:     100,
				SourceMatrixID: "familia-divorcio-consensual-filhos-bens",
				CTAContext:     "Origem: familia-divorcio-consensual-filhos-bens-prova-documental; Fonte: https://www.cnj.jus.br/; documentos esperados para triagem por WhatsApp.",
			},
		},
		ArchiveCountByBatch: map[string]int{"batch-familia-digital": 100},
		CurrentCandidateCountByBatch: map[string]int{
			"batch-familia-digital": 3,
		},
		MaxSimilarity: 0.52,
	}
	record := batchcandidateexpansion.Record{
		ReadinessID:                 "bad-expansion",
		BatchID:                     "batch-familia-digital",
		ReadinessStatus:             "ready_to_publish",
		SourceArchivePath:           "data/editorial/outro.jsonl",
		ArchiveRecordsRequired:      100,
		ArchiveRecordsObserved:      10,
		CurrentCandidateCount:       1,
		TargetCandidateCount:        5,
		TargetCandidateTier:         "target_1000",
		ExpansionCandidateIntentIDs: []string{"intent-inexistente"},
		PaidIntentGatePath:          "data/editorial/outro.jsonl",
		PaidIntentPassedCount:       0,
		PaidIntentBlockedCount:      0,
		PaidIntentMissingCount:      0,
		MinimumHumanScore:           70,
		MaxSimilarityAllowed:        0.90,
		MaxSimilarityObserved:       0.52,
		CTAContextRequired:          false,
		SourceURLAuditRequired:      false,
		SourceSpecificityRequired:   false,
		PaidIntentRequired:          false,
		PublicationBlockReason:      "",
		IndexPolicy:                 "index",
		ManifestAllowed:             true,
		RenderAllowed:               true,
		SitemapAllowed:              true,
		PublicationAllowed:          true,
		PublicPath:                  "/temas/familia/",
		CheckedAt:                   "2026-06-09",
	}

	report := batchcandidateexpansion.ValidateRecordAgainstIndex(record, index)
	if report.Passed() {
		t.Fatal("ValidateRecordAgainstIndex passed, want blocked expansion failures")
	}
	for _, code := range []string{
		"batch_candidate_expansion_status_not_blocked",
		"batch_candidate_expansion_wrong_archive_path",
		"batch_candidate_expansion_observed_count_mismatch",
		"batch_candidate_expansion_current_count_mismatch",
		"batch_candidate_expansion_target_too_low",
		"batch_candidate_expansion_too_few_intents",
		"batch_candidate_expansion_intent_missing",
		"batch_candidate_expansion_paid_gate_path_invalid",
		"batch_candidate_expansion_paid_gate_missing",
		"batch_candidate_expansion_paid_counts_mismatch",
		"batch_candidate_expansion_missing_legal_area",
		"batch_candidate_expansion_target_tier_invalid",
		"batch_candidate_expansion_missing_actionable_blocker",
		"batch_candidate_expansion_invalid_index_policy",
		"batch_candidate_expansion_manifest_allowed",
		"batch_candidate_expansion_similarity_limit_invalid",
		"batch_candidate_expansion_minimum_score_too_low",
		"batch_candidate_expansion_without_cta_requirement",
		"batch_candidate_expansion_without_source_audit",
		"batch_candidate_expansion_without_source_specificity",
		"batch_candidate_expansion_without_paid_intent",
		"batch_candidate_expansion_missing_block_reason",
		"batch_candidate_expansion_render_allowed",
		"batch_candidate_expansion_sitemap_allowed",
		"batch_candidate_expansion_publication_allowed",
		"batch_candidate_expansion_has_public_path",
	} {
		if !report.HasIssue(code) {
			t.Fatalf("missing issue %q in %v", code, report.Codes())
		}
	}
}
