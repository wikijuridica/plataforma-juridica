package contract_test

import (
	"testing"

	"portaljuridico/internal/batchdraftgen"
	"portaljuridico/internal/batchsourcematrix"
)

func TestBatchDraftGeneratorScalesWithSourceMatrixCoverage(t *testing.T) {
	matrixReport := batchsourcematrix.Validate(".")
	if !matrixReport.Passed() {
		t.Fatalf("source matrix failed contract: %v", matrixReport.Messages())
	}

	records, loadReport := batchsourcematrix.LoadRecords(".")
	if !loadReport.Passed() {
		t.Fatalf("could not load source matrix: %v", loadReport.Messages())
	}
	if len(records) < 30 {
		t.Fatalf("source matrix records=%d, want at least 30 subthemes", len(records))
	}

	result, report := batchdraftgen.Generate(".", batchdraftgen.Options{SamplesPerBatch: 10, CheckedAt: "2026-06-09"})
	if !report.Passed() {
		t.Fatalf("scaled batch draft generation failed: %v", report.Messages())
	}
	if len(result.Drafts) < 60 {
		t.Fatalf("generated drafts=%d, want at least 60", len(result.Drafts))
	}
	if result.StructuralPatternRisk() > 0.35 {
		t.Fatalf("structural risk=%.2f, want <=0.35", result.StructuralPatternRisk())
	}
	if result.MaximumPairSimilarity() > 0.64 {
		t.Fatalf("max generated similarity=%.2f, want <=0.64", result.MaximumPairSimilarity())
	}

	coverage := batchsourcematrix.ValidateDraftCoverage(records, result.Drafts)
	if !coverage.Passed() {
		t.Fatalf("generated drafts lack source matrix coverage: %v", coverage.Messages())
	}
	for _, metric := range result.Metrics {
		if metric.GeneratedSamples < 10 {
			t.Fatalf("%s generated_samples=%d, want >=10", metric.BatchID, metric.GeneratedSamples)
		}
		if metric.SourceMatrixCoveredSamples != metric.GeneratedSamples {
			t.Fatalf("%s source_matrix_covered=%d generated=%d", metric.BatchID, metric.SourceMatrixCoveredSamples, metric.GeneratedSamples)
		}
		if metric.StructuralPatternRisk > 0.35 {
			t.Fatalf("%s structural risk=%.2f, want <=0.35", metric.BatchID, metric.StructuralPatternRisk)
		}
		if metric.LabEstimatedCPUUnits <= 0 {
			t.Fatalf("%s missing lab cpu estimate", metric.BatchID)
		}
	}
}

func TestBatchSourceMatrixRejectsWeakOrPublicSourceRecord(t *testing.T) {
	record := batchsourcematrix.Record{
		MatrixID:                "unsafe",
		BatchID:                 "batch-saude-suplementar-digital",
		LegalArea:               "saude-suplementar",
		SubthemeID:              "unsafe",
		SourceStatus:            "published",
		SourceURLs:              []string{"https://example.com/blog"},
		SourceTypes:             []string{"blog"},
		SourceSpecificityScore:  30,
		OfficialSourcesVerified: false,
		RobotsReviewRequired:    false,
		UsePolicy:               "copy_text",
		RenderAllowed:           true,
		SitemapAllowed:          true,
		PublicationAllowed:      true,
		PublicPath:              "/temas/unsafe/",
		CheckedAt:               "2026-06-09",
	}
	report := batchsourcematrix.ValidateRecord(record)
	for _, code := range []string{
		"source_matrix_invalid_status",
		"source_matrix_too_few_official_urls",
		"source_matrix_too_few_source_types",
		"source_matrix_specificity_too_low",
		"source_matrix_not_verified",
		"source_matrix_robots_review_missing",
		"source_matrix_use_policy_invalid",
		"source_matrix_render_allowed",
		"source_matrix_sitemap_allowed",
		"source_matrix_publication_allowed",
		"source_matrix_has_public_path",
	} {
		if !report.HasIssue(code) {
			t.Fatalf("missing issue %q in %v", code, report.Codes())
		}
	}
}
