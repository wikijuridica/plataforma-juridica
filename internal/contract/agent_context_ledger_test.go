package contract_test

import (
	"testing"

	"portaljuridico/internal/agentcontext"
)

func TestAgentContextLedgerPersistsSubagentSummaries(t *testing.T) {
	report := agentcontext.Validate(".")
	if !report.Passed() {
		t.Fatalf("agent context ledger failed: %v", report.Messages())
	}

	records, loadReport := agentcontext.LoadRecords(".")
	if !loadReport.Passed() {
		t.Fatalf("could not load agent context ledger: %v", loadReport.Messages())
	}
	if len(records) < 3 {
		t.Fatalf("agent context records=%d, want at least cycle 40 subagent summaries", len(records))
	}

	var hasSourceResearch, hasReadinessAudit, hasPaidIntentReview bool
	for _, entry := range records {
		record := entry.Record
		if record.Cycle < 40 {
			continue
		}
		if record.UsagePolicy != agentcontext.ReferenceOnlyPolicy {
			t.Fatalf("line=%d policy=%q", entry.Line, record.UsagePolicy)
		}
		if record.RepoWriteAllowed || !record.CodexValidationRequired || !record.ClosedBeforeCheckpoint {
			t.Fatalf("line=%d unsafe context flags write=%t validation=%t closed=%t", entry.Line, record.RepoWriteAllowed, record.CodexValidationRequired, record.ClosedBeforeCheckpoint)
		}
		switch record.TaskKind {
		case "source_research":
			hasSourceResearch = true
		case "readiness_audit":
			hasReadinessAudit = true
		case "paid_intent_review":
			hasPaidIntentReview = true
		}
	}
	if !hasSourceResearch || !hasReadinessAudit || !hasPaidIntentReview {
		t.Fatalf("missing expected cycle 40 agent summaries: source=%t readiness=%t paid=%t", hasSourceResearch, hasReadinessAudit, hasPaidIntentReview)
	}
}

func TestAgentContextLedgerRejectsUnsafeDelegationState(t *testing.T) {
	record := agentcontext.Record{
		Cycle:                   41,
		AgentID:                 "019eadc3-example",
		Nickname:                "Unsafe",
		TaskKind:                "source_research",
		Scope:                   "pesquisa de fontes oficiais",
		Status:                  "completed",
		UsagePolicy:             "direct_write",
		Summary:                 "saida curta",
		Evidence:                []string{},
		Risks:                   []string{},
		IntegrationDecision:     "",
		RepoWriteAllowed:        true,
		CodexValidationRequired: false,
		ClosedBeforeCheckpoint:  false,
		RecordedAt:              "2026-06-09T21:00:00-03:00",
	}

	report := agentcontext.ValidateRecord(record)
	if report.Passed() {
		t.Fatal("unsafe agent context record passed, want failures")
	}
	for _, code := range []string{
		"agent_context_usage_policy_invalid",
		"agent_context_repo_write_allowed",
		"agent_context_validation_not_required",
		"agent_context_not_closed_before_checkpoint",
		"agent_context_missing_integration_decision",
		"agent_context_missing_evidence",
		"agent_context_missing_risk",
	} {
		if !report.HasIssue(code) {
			t.Fatalf("missing issue %q in %v", code, report.Codes())
		}
	}
}
