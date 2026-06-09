package contract_test

import (
	"reflect"
	"testing"

	"portaljuridico/internal/batchdraftgen"
	"portaljuridico/internal/batchdrafts"
)

func TestBatchDraftGeneratorProducesDeterministicBlockedScoredDrafts(t *testing.T) {
	options := batchdraftgen.Options{SamplesPerBatch: 5, CheckedAt: "2026-06-09"}
	result, report := batchdraftgen.Generate(".", options)
	if !report.Passed() {
		t.Fatalf("batch draft generation failed: %v", report.Messages())
	}

	again, secondReport := batchdraftgen.Generate(".", options)
	if !secondReport.Passed() {
		t.Fatalf("second generation failed: %v", secondReport.Messages())
	}
	if !reflect.DeepEqual(result.DraftIDs(), again.DraftIDs()) {
		t.Fatalf("generation is not deterministic\nfirst=%v\nsecond=%v", result.DraftIDs(), again.DraftIDs())
	}

	if len(result.Drafts) < 30 {
		t.Fatalf("generated drafts=%d, want at least 30", len(result.Drafts))
	}
	if len(result.Metrics) < 6 {
		t.Fatalf("generation metrics=%d, want at least 6", len(result.Metrics))
	}
	if result.RewrittenCount() < 6 {
		t.Fatalf("rewritten drafts=%d, want at least one automatic rewrite per batch", result.RewrittenCount())
	}
	if result.MaximumPairSimilarity() > 0.64 {
		t.Fatalf("max generated similarity=%.2f, want <=0.64", result.MaximumPairSimilarity())
	}

	for _, draft := range result.Drafts {
		if validation := batchdrafts.ValidateRecord(draft); !validation.Passed() {
			t.Fatalf("generated draft %s failed batch draft contract: %v", draft.UniqueIntentID, validation.Messages())
		}
		if draft.RenderAllowed || draft.SitemapAllowed || draft.PublicationAllowed || draft.PublicPath != "" {
			t.Fatalf("generated draft %s escaped blocked-publication lab contract", draft.UniqueIntentID)
		}
		if draft.CTAContext == "" || draft.UniqueIntentID == "" || !batchdraftgen.ContainsOrigin(draft.CTAContext, draft.UniqueIntentID) {
			t.Fatalf("generated draft %s lacks contextual WhatsApp origin: %q", draft.UniqueIntentID, draft.CTAContext)
		}
	}

	metricReport := batchdraftgen.ValidateMetrics(result.Metrics)
	if !metricReport.Passed() {
		t.Fatalf("generated metrics failed contract: %v", metricReport.Messages())
	}

	storedMetricReport := batchdraftgen.ValidateStoredMetrics(".")
	if !storedMetricReport.Passed() {
		t.Fatalf("stored generation metrics failed contract: %v", storedMetricReport.Messages())
	}
}

func TestBatchDraftGenerationRejectsUnsafeOptionsAndPublicMetrics(t *testing.T) {
	_, report := batchdraftgen.Generate(".", batchdraftgen.Options{SamplesPerBatch: 1, CheckedAt: "2026-06-09"})
	if !report.HasIssue("generation_samples_too_low") {
		t.Fatalf("missing generation_samples_too_low in %v", report.Codes())
	}

	metric := batchdraftgen.Metric{
		BatchID:            "batch-saude-suplementar-digital",
		GenerationStatus:   "published",
		GeneratedSamples:   1,
		PassedSamples:      0,
		MinimumHumanScore:  40,
		MaximumAILikeScore: 90,
		MaximumSimilarity:  0.91,
		RenderAllowed:      true,
		SitemapAllowed:     true,
		PublicationAllowed: true,
		PublicPath:         "/temas/unsafe/",
		CheckedAt:          "2026-06-09",
	}
	report = batchdraftgen.ValidateMetric(metric)
	for _, code := range []string{
		"generation_metric_invalid_status",
		"generation_metric_too_few_samples",
		"generation_metric_unpassed_samples",
		"generation_metric_human_score_too_low",
		"generation_metric_ai_score_too_high",
		"generation_metric_similarity_too_high",
		"generation_metric_render_allowed",
		"generation_metric_sitemap_allowed",
		"generation_metric_publication_allowed",
		"generation_metric_has_public_path",
	} {
		if !report.HasIssue(code) {
			t.Fatalf("missing issue %q in %v", code, report.Codes())
		}
	}
}
