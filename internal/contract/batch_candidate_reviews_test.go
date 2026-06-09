package contract_test

import (
	"testing"

	"portaljuridico/internal/batchcandidatereviews"
)

func TestBatchCandidateReviewsCoverEverySelectedIntentWithoutPublishing(t *testing.T) {
	report := batchcandidatereviews.Validate(".")
	if !report.Passed() {
		t.Fatalf("batch candidate reviews failed contract: %v", report.Messages())
	}

	records, loadReport := batchcandidatereviews.LoadRecords(".")
	if !loadReport.Passed() {
		t.Fatalf("could not load batch candidate reviews: %v", loadReport.Messages())
	}

	if len(records) != 18 {
		t.Fatalf("reviews=%d, want one blocked legal-editorial review for each of the 18 selected batch intents", len(records))
	}

	seen := make(map[string]bool)
	for _, entry := range records {
		record := entry.Record
		if seen[record.UniqueIntentID] {
			t.Fatalf("duplicate review for selected intent %q", record.UniqueIntentID)
		}
		seen[record.UniqueIntentID] = true

		if record.BaseURLMode != "official_configured" {
			t.Fatalf("line=%d base_url_mode=%q, want official_configured after wikijuridica.com.br was defined", entry.Line, record.BaseURLMode)
		}
		if !record.OfficialURLLocked {
			t.Fatalf("line=%d official_url_locked=false after wikijuridica.com.br was defined", entry.Line)
		}
		if record.RenderAllowed || record.SitemapAllowed || record.PublicationAllowed || record.PublicPath != "" {
			t.Fatalf("line=%d review escaped blocked contract: render=%t sitemap=%t publication=%t path=%q", entry.Line, record.RenderAllowed, record.SitemapAllowed, record.PublicationAllowed, record.PublicPath)
		}
		if !record.CTAContextMessageContains("Origem: " + record.CandidatePath) {
			t.Fatalf("line=%d cta context missing origin path %q: %s", entry.Line, record.CandidatePath, record.CTAContextMessage)
		}
		if !record.CTAContextMessageContains("Gate: " + record.GateID) {
			t.Fatalf("line=%d cta context missing gate id %q: %s", entry.Line, record.GateID, record.CTAContextMessage)
		}
		if !record.CTAContextMessageContains("Intent: " + record.UniqueIntentID) {
			t.Fatalf("line=%d cta context missing unique intent id %q: %s", entry.Line, record.UniqueIntentID, record.CTAContextMessage)
		}
	}
}

func TestBatchCandidateReviewRejectsUnknownIntentFixedDomainOrPublicCTA(t *testing.T) {
	index := batchcandidatereviews.CandidateIndex{
		BaseURLMode:       "official_configured",
		OfficialURLLocked: true,
		SelectedByIntent: map[string]batchcandidatereviews.SelectedCandidate{
			"familia-divorcio-consensual-filhos-bens": {
				GateID:         "gate-familia-digital",
				BatchID:        "batch-familia-digital",
				UniqueIntentID: "familia-divorcio-consensual-filhos-bens",
				CandidatePath:  "/temas/familia-divorcio-consensual-filhos-bens/",
				SourceMatrixID: "matrix-familia-digital",
			},
		},
		AuditedMatrixIDs: map[string]bool{
			"matrix-familia-digital": true,
		},
	}

	record := batchcandidatereviews.Record{
		ReviewID:               "bad-review",
		GateID:                 "gate-familia-digital",
		BatchID:                "batch-familia-digital",
		UniqueIntentID:         "familia-divorcio-consensual-inventado",
		CandidatePath:          "https://www.exemplo-oficial.com/temas/familia-divorcio-consensual-inventado/",
		BaseURLMode:            "lab_placeholder",
		OfficialURLLocked:      false,
		ReviewStatus:           "approved_public",
		Language:               "pt-BR",
		SourceMatrixID:         "matrix-nao-auditada",
		ReviewerRole:           "juridico_editorial_lab",
		LegalReviewNotes:       []string{"nota curta"},
		RequiredFixes:          []string{},
		CTADraft:               "Garantimos resultado pelo WhatsApp em 24 horas para resolver o processo digital.",
		CTAContextMessage:      "Quero atendimento.",
		CTAStatus:              "approved_public",
		PublicationBlockReason: "",
		RenderAllowed:          true,
		SitemapAllowed:         true,
		PublicationAllowed:     true,
		PublicPath:             "/temas/familia-divorcio-consensual-inventado/",
		CheckedAt:              "2026-06-09",
	}

	report := batchcandidatereviews.ValidateRecordAgainstCandidateIndex(record, index)
	if report.Passed() {
		t.Fatal("ValidateRecordAgainstCandidateIndex passed, want blocked review failures")
	}
	for _, code := range []string{
		"batch_candidate_review_unknown_selected_intent",
		"batch_candidate_review_path_contains_domain",
		"batch_candidate_review_url_config_mismatch",
		"batch_candidate_review_invalid_status",
		"batch_candidate_review_too_few_notes",
		"batch_candidate_review_missing_required_fixes",
		"batch_candidate_review_source_matrix_not_audited",
		"batch_candidate_review_promise_cta",
		"batch_candidate_review_invalid_cta_status",
		"batch_candidate_review_cta_context_without_origin",
		"batch_candidate_review_missing_block_reason",
		"batch_candidate_review_render_allowed",
		"batch_candidate_review_sitemap_allowed",
		"batch_candidate_review_publication_allowed",
		"batch_candidate_review_has_public_path",
	} {
		if !report.HasIssue(code) {
			t.Fatalf("missing issue %q in %v", code, report.Codes())
		}
	}
}
