package contract_test

import (
	"testing"

	"portaljuridico/internal/batchfinaldrafts"
)

func TestBatchFinalAuthorialDraftsCoverEligibleManifestWithoutPublishing(t *testing.T) {
	report := batchfinaldrafts.Validate(".")
	if !report.Passed() {
		t.Fatalf("batch final authorial drafts failed contract: %v", report.Messages())
	}

	records, loadReport := batchfinaldrafts.LoadRecords(".")
	if !loadReport.Passed() {
		t.Fatalf("could not load batch final authorial drafts: %v", loadReport.Messages())
	}
	if len(records) != 7 {
		t.Fatalf("final authorial drafts=%d, want one blocked draft for each source-locked manifest gate", len(records))
	}

	for _, entry := range records {
		record := entry.Record
		if record.Language != "pt-BR" {
			t.Fatalf("line=%d language=%q", entry.Line, record.Language)
		}
		if record.HumanScore < 85 || record.AILikeScore > 20 {
			t.Fatalf("line=%d invalid score human=%d ai=%d", entry.Line, record.HumanScore, record.AILikeScore)
		}
		if record.RenderAllowed || record.SitemapAllowed || record.PublicationAllowed || record.PublicPath != "" {
			t.Fatalf("line=%d draft escaped blocked contract: render=%t sitemap=%t publication=%t path=%q", entry.Line, record.RenderAllowed, record.SitemapAllowed, record.PublicationAllowed, record.PublicPath)
		}
		if record.CTAContextMessage == "" || !record.HasContextualCTA() {
			t.Fatalf("line=%d CTA must carry origin, intent and documents", entry.Line)
		}
		if !record.HasInformationalNotice() {
			t.Fatalf("line=%d missing informational legal notice", entry.Line)
		}
	}
}

func TestBatchFinalAuthorialDraftRejectsSourceBlockedOrPublicDraft(t *testing.T) {
	index := batchfinaldrafts.ManifestIndex{
		BaseURL: "https://wikijuridica.com.br",
		EligibleByIntent: map[string]batchfinaldrafts.EligibleManifest{
			"consumidor-financeiro-pix-fraude-resposta-banco": {
				ManifestGateID:           "public-manifest-consumidor-financeiro-pix-fraude-resposta-banco",
				UniqueIntentID:           "consumidor-financeiro-pix-fraude-resposta-banco",
				BatchID:                  "batch-consumidor-financeiro-digital",
				SourceMatrixID:           "consumidor-financeiro-pix-fraude-resposta-banco",
				Term:                     "fraude via Pix e resposta insuficiente do banco",
				CandidatePath:            "/temas/consumidor-financeiro-pix-fraude-resposta-banco/",
				CandidateCanonicalURL:    "https://wikijuridica.com.br/temas/consumidor-financeiro-pix-fraude-resposta-banco/",
				CandidateTitle:           "Fraude via Pix e resposta insuficiente do banco",
				CandidateMetaDescription: "Entenda como cronologia, comprovante Pix, boletim e protocolos ajudam a avaliar resposta bancária insuficiente.",
				SelectedSourceURLs:       []string{"https://www.bcb.gov.br/estabilidadefinanceira/pix"},
			},
		},
		BlockedByIntent: map[string]bool{
			"familia-divorcio-consensual-filhos-bens": true,
		},
	}

	record := batchfinaldrafts.Record{
		DraftID:                  "bad-final-draft",
		ManifestGateID:           "public-manifest-familia-divorcio-consensual-filhos-bens",
		UniqueIntentID:           "familia-divorcio-consensual-filhos-bens",
		BatchID:                  "batch-familia-digital",
		SourceMatrixID:           "familia-divorcio-consensual-filhos-bens",
		Term:                     "divórcio consensual online com filhos e bens",
		CandidatePath:            "/temas/familia-divorcio-consensual-filhos-bens/?utm=1",
		CandidateCanonicalURL:    "https://portal-juridico.example/temas/familia-divorcio-consensual-filhos-bens/",
		CandidateTitle:           "Divórcio",
		CandidateMetaDescription: "Meta curta.",
		SelectedSourceURLs:       []string{"https://www.cnj.jus.br/"},
		DraftStatus:              "published",
		Language:                 "pt-BR",
		Opening:                  "Texto curto e genérico.",
		SourceUse:                "Sem fonte específica.",
		DocumentGuidance:         "Sem documentos.",
		DigitalTriage:            "Sem triagem.",
		CTAContextMessage:        "Mensagem sem origem.",
		InformationalNotice:      "Aviso fraco.",
		HumanScore:               100,
		AILikeScore:              0,
		IndexPolicy:              "index",
		RenderAllowed:            true,
		SitemapAllowed:           true,
		PublicationAllowed:       true,
		PublicPath:               "/temas/familia-divorcio-consensual-filhos-bens/",
		CheckedAt:                "2026-06-09",
	}

	report := batchfinaldrafts.ValidateRecordAgainstManifestIndex(record, index)
	if report.Passed() {
		t.Fatal("ValidateRecordAgainstManifestIndex passed, want blocked draft failures")
	}
	for _, code := range []string{
		"batch_final_draft_source_blocked",
		"batch_final_draft_candidate_path_not_clean",
		"batch_final_draft_canonical_mismatch",
		"batch_final_draft_title_too_short",
		"batch_final_draft_meta_too_short",
		"batch_final_draft_invalid_status",
		"batch_final_draft_invalid_index_policy",
		"batch_final_draft_text_failed_human_score",
		"batch_final_draft_text_failed_quality",
		"batch_final_draft_cta_not_contextual",
		"batch_final_draft_missing_notice",
		"batch_final_draft_render_allowed",
		"batch_final_draft_sitemap_allowed",
		"batch_final_draft_publication_allowed",
		"batch_final_draft_has_public_path",
	} {
		if !report.HasIssue(code) {
			t.Fatalf("missing issue %q in %v", code, report.Codes())
		}
	}
}
