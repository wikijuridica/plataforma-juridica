package contract_test

import (
	"testing"

	"portaljuridico/internal/batchdrafts"
)

func TestBatchDraftsGenerateScoredSamplesAcrossMassBatches(t *testing.T) {
	report := batchdrafts.Validate(".")
	if !report.Passed() {
		t.Fatalf("batch drafts failed contract: %v", report.Messages())
	}

	records, loadReport := batchdrafts.LoadRecords(".")
	if !loadReport.Passed() {
		t.Fatalf("could not load batch drafts: %v", loadReport.Messages())
	}
	if len(records) < 18 {
		t.Fatalf("batch draft count=%d, want at least 18 samples", len(records))
	}
	if batchdrafts.RewrittenCount(records) < 3 {
		t.Fatalf("rewritten drafts=%d, want at least 3 automatic rewrite proofs", batchdrafts.RewrittenCount(records))
	}
	if batchdrafts.MaximumPairSimilarity(records) > 0.64 {
		t.Fatalf("max similarity=%.2f, want <=0.64", batchdrafts.MaximumPairSimilarity(records))
	}

	perBatch := make(map[string]int)
	for _, entry := range records {
		record := entry.Record
		perBatch[record.BatchID]++
		if record.DraftStatus != "batch_draft_scored_blocked" {
			t.Fatalf("%s status=%q, want batch_draft_scored_blocked", record.UniqueIntentID, record.DraftStatus)
		}
		if record.HumanScore < 85 || record.AILikeScore > 20 {
			t.Fatalf("%s weak score human=%d ai=%d", record.UniqueIntentID, record.HumanScore, record.AILikeScore)
		}
		if record.RenderAllowed || record.SitemapAllowed || record.PublicationAllowed || record.PublicPath != "" {
			t.Fatalf("%s escaped blocked publication contract", record.UniqueIntentID)
		}
	}
	for batchID, count := range perBatch {
		if count < 3 {
			t.Fatalf("batch %s has %d samples, want at least 3", batchID, count)
		}
	}
}

func TestBatchDraftRejectsMechanicalRewriteOrPublication(t *testing.T) {
	record := batchdrafts.Record{
		BatchID:            "batch-saude-suplementar-digital",
		UniqueIntentID:     "negativa-cobertura-plano-saude-spam",
		LegalArea:          "saude-suplementar",
		DraftStatus:        "published",
		Language:           "pt-BR",
		Term:               "negativa cobertura plano saúde",
		ReaderProblem:      "Este conteúdo explica negativa cobertura plano saúde de forma geral.",
		SourceHook:         "Fonte genérica.",
		DocumentContext:    "Documentos necessários.",
		RiskContext:        "Quando procurar advogado.",
		DigitalAction:      "Chamar no WhatsApp.",
		CTAContext:         "Chamar no WhatsApp.",
		RewriteStatus:      "not_rewritten",
		RewriteAttempts:    0,
		HumanScore:         50,
		AILikeScore:        70,
		RenderAllowed:      true,
		SitemapAllowed:     true,
		PublicationAllowed: true,
		PublicPath:         "/temas/spam/",
		CheckedAt:          "2026-06-09",
	}

	report := batchdrafts.ValidateRecord(record)
	if report.Passed() {
		t.Fatal("ValidateRecord passed mechanical public draft, want failures")
	}
	for _, code := range []string{
		"batch_draft_invalid_status",
		"batch_draft_human_score_too_low",
		"batch_draft_ai_score_too_high",
		"batch_draft_text_failed_human_score",
		"batch_draft_rewrite_missing",
		"batch_draft_render_allowed",
		"batch_draft_sitemap_allowed",
		"batch_draft_publication_allowed",
		"batch_draft_has_public_path",
	} {
		if !report.HasIssue(code) {
			t.Fatalf("missing issue %q in %v", code, report.Codes())
		}
	}
}
