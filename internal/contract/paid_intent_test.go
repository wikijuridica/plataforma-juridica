package contract_test

import (
	"testing"

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
