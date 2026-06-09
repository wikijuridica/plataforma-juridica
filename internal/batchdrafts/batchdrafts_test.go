package batchdrafts

import "testing"

func TestMaximumPairSimilarityDetailCachesByContentFingerprint(t *testing.T) {
	resetSimilarityCacheForTest()
	entries := []Entry{
		{Line: 1, Record: similarityTestRecord("consumidor-cartao-negativado-prova-digital", "Contrato do cartao e fatura mostram a cobranca contestada.")},
		{Line: 2, Record: similarityTestRecord("consumidor-cartao-negativado-linha-do-tempo", "Contrato do cartao e fatura mostram a cobranca contestada.")},
		{Line: 3, Record: similarityTestRecord("familia-divorcio-consensual-prova-digital", "Documentos familiares e minuta mostram acordo em construcao.")},
	}

	first := MaximumPairSimilarityDetail(entries)
	afterFirst := similarityCacheStatsForTest()
	if afterFirst.Misses != 1 || afterFirst.Hits != 0 || afterFirst.Entries != 1 {
		t.Fatalf("after first call stats=%+v, want misses=1 hits=0 entries=1", afterFirst)
	}

	second := MaximumPairSimilarityDetail(entries)
	afterSecond := similarityCacheStatsForTest()
	if second != first {
		t.Fatalf("cached pair changed: first=%+v second=%+v", first, second)
	}
	if afterSecond.Misses != 1 || afterSecond.Hits != 1 || afterSecond.Entries != 1 {
		t.Fatalf("after second call stats=%+v, want cache hit without recompute", afterSecond)
	}

	entries[1].Record.ReaderProblem = "Plano de saude negou procedimento com laudo, negativa formal e risco clinico urgente."
	_ = MaximumPairSimilarityDetail(entries)
	afterMutation := similarityCacheStatsForTest()
	if afterMutation.Misses != 2 || afterMutation.Hits != 1 || afterMutation.Entries != 2 {
		t.Fatalf("after mutation stats=%+v, want content fingerprint miss without masking changed text", afterMutation)
	}
}

func similarityTestRecord(intentID string, readerProblem string) Record {
	return Record{
		BatchID:         "batch-consumidor-financeiro-digital",
		UniqueIntentID:  intentID,
		SourceMatrixID:  "consumidor-cartao-negativado",
		LegalArea:       "consumidor-financeiro",
		Term:            "cartao negativado indevidamente",
		ReaderProblem:   readerProblem,
		SourceHook:      "CDC e Banco Central orientam comprovantes, contestacao e contrato.",
		DocumentContext: "Contrato, fatura, protocolo, notificacao e comprovantes organizam a triagem.",
		RiskContext:     "Negativacao indevida pode afetar credito, empresa e valor envolvido.",
		DigitalAction:   "Atendimento online organiza documentos antes de orcamento de honorarios.",
		CTAContext:      "Origem: teste; Intent: " + intentID + "; Documentos: contrato, fatura e protocolo.",
	}
}
