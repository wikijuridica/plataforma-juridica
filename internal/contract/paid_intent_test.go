package contract_test

import (
	"sort"
	"testing"

	"portaljuridico/internal/batchcandidateexpansion"
	"portaljuridico/internal/paidintent"
)

func TestPaidIntentRejectsFreeServiceSignalsAndAcceptsHiringSignals(t *testing.T) {
	free := paidintent.ValidateText("Quero advogado gratuito, defensoria pública, consulta grátis e serviço jurídico sem pagar.", paidintent.RequirePaidSignal)
	if free.Passed() {
		t.Fatal("free-service text passed, want paid intent blocker")
	}
	if !free.HasIssue("paid_intent_free_service_signal") {
		t.Fatalf("missing free-service issue in %v", free.Codes())
	}

	paid := paidintent.ValidateText("Busco contratação particular de advogado online, com orçamento de honorários e triagem paga dos documentos.", paidintent.RequirePaidSignal)
	if !paid.Passed() {
		t.Fatalf("paid-service text failed: %v", paid.Messages())
	}
}

func TestPaidIntentRejectsDraftsWithoutHiringSignal(t *testing.T) {
	report := paidintent.Validate(".")
	if !report.Passed() {
		t.Fatalf("paid intent gate failed repository contract: %v", report.Messages())
	}

	neutral := paidintent.ValidateText("Origem: /temas/exemplo/ | Intent: exemplo | Documentos: contrato e protocolo.", paidintent.RequirePaidSignal)
	if neutral.Passed() {
		t.Fatal("neutral CTA passed, want explicit paid hiring signal")
	}
	if !neutral.HasIssue("paid_intent_missing_paid_signal") {
		t.Fatalf("missing paid-signal issue in %v", neutral.Codes())
	}
}

func TestPaidIntentRejectsPublicAssistanceAndSelfServiceRisk(t *testing.T) {
	publicAssistance := paidintent.ValidateText("BPC LOAS negado por renda familiar, CadÚnico, vulnerabilidade e justiça gratuita. Contratação: triagem particular com honorários.", paidintent.RequirePaidSignal)
	if publicAssistance.Passed() {
		t.Fatal("public-assistance text passed, want commercial blocker")
	}
	if !publicAssistance.HasIssue("paid_intent_public_assistance_free_risk") {
		t.Fatalf("missing public-assistance issue in %v", publicAssistance.Codes())
	}

	selfService := paidintent.ValidateText("Cumprimento de exigência do INSS sem resposta, anexar documentos pelo Meu INSS e reenviar arquivos. Contratação: análise particular com honorários.", paidintent.RequirePaidSignal)
	if selfService.Passed() {
		t.Fatal("administrative self-service text passed, want commercial blocker")
	}
	if !selfService.HasIssue("paid_intent_admin_self_service_risk") {
		t.Fatalf("missing self-service issue in %v", selfService.Codes())
	}
}

func TestPaidIntentRejectsCTAOnlyPaidSignal(t *testing.T) {
	ctaOnly := paidintent.ValidateTextParts(
		"Plano de saude negou procedimento urgente com negativa formal, contrato, protocolo, documentos, risco clinico e valor envolvido.",
		"WhatsApp contextual: contratar advogado particular online com orcamento de honorarios.",
		paidintent.RequirePaidSignal,
	)
	if ctaOnly.Passed() {
		t.Fatal("CTA-only paid signal passed, want blocker before mass generation")
	}
	if !ctaOnly.HasIssue("paid_intent_cta_only_paid_signal") {
		t.Fatalf("missing CTA-only issue in %v", ctaOnly.Codes())
	}

	strong := paidintent.ValidateTextParts(
		"Cliente quer contratar advogado particular online para analisar negativa formal do plano de saude, com honorarios, contrato, protocolo, documentos, risco clinico e valor envolvido.",
		"WhatsApp contextual: origem da pagina, tema de saude suplementar e documentos principais.",
		paidintent.RequirePaidSignal,
	)
	if !strong.Passed() {
		t.Fatalf("strong paid-intent core failed: %v", strong.Messages())
	}
}

func TestPaidIntentCoversExpansionReadinessTargetsWithoutPublishing(t *testing.T) {
	report := paidintent.Validate(".")
	if !report.Passed() {
		t.Fatalf("paid intent gate failed repository contract: %v", report.Messages())
	}

	readinessRecords, readinessReport := batchcandidateexpansion.LoadRecords(".")
	if !readinessReport.Passed() {
		t.Fatalf("could not load expansion readiness records: %v", readinessReport.Messages())
	}
	targets := make(map[string]bool)
	for _, entry := range readinessRecords {
		for _, intentID := range entry.Record.ExpansionCandidateIntentIDs {
			targets[intentID] = true
		}
	}
	if len(targets) != 180 {
		t.Fatalf("readiness targets=%d, want 180 unique expansion intents", len(targets))
	}

	paidRecords, loadReport := paidintent.LoadRecords(".")
	if !loadReport.Passed() {
		t.Fatalf("could not load paid intent gates: %v", loadReport.Messages())
	}
	paidByIntent := make(map[string]paidintent.Record)
	for _, entry := range paidRecords {
		record := entry.Record
		if record.GateScope == "" {
			t.Fatalf("line=%d paid intent gate without scope", entry.Line)
		}
		if record.IndexPolicy != "noindex" || record.RenderAllowed || record.SitemapAllowed || record.PublicationAllowed || record.PublicPath != "" {
			t.Fatalf("line=%d paid gate escaped blocked publication contract", entry.Line)
		}
		paidByIntent[record.UniqueIntentID] = record
	}

	missing := make([]string, 0)
	for intentID := range targets {
		if _, ok := paidByIntent[intentID]; !ok {
			missing = append(missing, intentID)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf("missing paid intent gates for %d expansion targets, first missing=%s", len(missing), missing[0])
	}
}

func TestPaidIntentRepositoryGatesBlockWeakCommercialDrafts(t *testing.T) {
	records, loadReport := paidintent.LoadRecords(".")
	if !loadReport.Passed() {
		t.Fatalf("could not load paid intent gates: %v", loadReport.Messages())
	}
	statusByIntent := make(map[string]string)
	for _, entry := range records {
		statusByIntent[entry.Record.UniqueIntentID] = entry.Record.PaidIntentStatus
		if entry.Record.RenderAllowed || entry.Record.SitemapAllowed || entry.Record.PublicationAllowed || entry.Record.PublicPath != "" {
			t.Fatalf("line=%d paid gate escaped blocked contract", entry.Line)
		}
	}
	for intent, status := range map[string]string{
		"previdenciario-bpc-loas-cadunico-renda":      paidintent.PublicAssistanceBlockedStatus,
		"previdenciario-cumprimento-exigencia-parado": paidintent.AdminSelfServiceBlockedStatus,
	} {
		if statusByIntent[intent] != status {
			t.Fatalf("intent=%s status=%q, want %q", intent, statusByIntent[intent], status)
		}
	}
	if statusByIntent["saude-suplementar-procedimento-urgente-negado"] != paidintent.PassedBlockedStatus {
		t.Fatalf("health urgent draft should pass paid-intent after body refinement, got %q", statusByIntent["saude-suplementar-procedimento-urgente-negado"])
	}
}
