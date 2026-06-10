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

func TestMaximumPairSimilarityDetailDoesNotHideDuplicateTextBehindExpansionProfiles(t *testing.T) {
	left := similarityTestRecord(
		"consumidor-financeiro-negativacao-divida-desconhecida-prova-digital-recurso-decisao-recorrida-prova-nova-rodada-02",
		"O consumidor descobriu a mesma restricao no cadastro, reuniu o mesmo protocolo, contrato contestado, extratos e resposta do credor para triagem online.",
	)
	left.SourceMatrixID = "consumidor-financeiro-negativacao-divida-desconhecida"
	left.LegalArea = "consumidor-financeiro"
	right := left
	right.UniqueIntentID = "consumidor-financeiro-negativacao-divida-desconhecida-prazo-e-urgencia-divergencia-cadastral-base-antiga-rodada-02"

	pair := MaximumPairSimilarityDetail([]Entry{
		{Line: 1, Record: left},
		{Line: 2, Record: right},
	})
	if pair.Score <= 0.64 {
		t.Fatalf("near-duplicate text hidden by expansion profiles: pair=%s/%s score=%.4f, want >0.64", pair.LeftID, pair.RightID, pair.Score)
	}
}

func TestMaximumPairSimilarityDetailCatchesDuplicateTextAcrossDifferentSourceMatrix(t *testing.T) {
	left := similarityTestRecord(
		"consumidor-financeiro-tarifa-bancaria-indevida-prova-digital",
		"O consumidor reuniu contrato, extratos, protocolo, resposta do banco e impacto financeiro para triagem juridica online.",
	)
	left.SourceMatrixID = "consumidor-financeiro-tarifa-bancaria-indevida"
	left.LegalArea = "consumidor-financeiro"
	right := left
	right.UniqueIntentID = "consumidor-financeiro-pix-fraude-resposta-banco-prazo-e-urgencia"
	right.SourceMatrixID = "consumidor-financeiro-pix-fraude-resposta-banco"

	pair := MaximumPairSimilarityDetail([]Entry{
		{Line: 1, Record: left},
		{Line: 2, Record: right},
	})
	if pair.Score <= 0.64 {
		t.Fatalf("duplicate text across source matrix passed: pair=%s/%s score=%.4f, want >0.64", pair.LeftID, pair.RightID, pair.Score)
	}
}

func TestMaximumPairSimilarityDetailIgnoresExpansionBoilerplateAcrossDifferentSourceMatrix(t *testing.T) {
	left := similarityTestRecord(
		"familia-partilha-bens-conta-digital-negociacao-previa-recurso-decisao-recorrida-prova-nova-rodada-02",
		"Histórico de tentativa prévia em família trata recurso documentado. A rodada nova muda o eixo do rascunho para fase recursal no subtema partilha de bens com conta digital e usa o caso matriz apenas como vínculo de origem editorial.",
	)
	left.SourceMatrixID = "familia-partilha-bens-conta-digital"
	left.DocumentContext = "Razões do recurso, decisão recorrida, protocolo e prova nova ficam em trilha separada. A documentação da rodada usa eixo documental partilha bens, separando prova nova, lacuna corrigida e fase atual sem repetir a lista do caso matriz."
	left.RiskContext = "Tratar recurso como simples novo pedido pode fazer a pessoa perder prazo. Risco operacional partilha bens mede se o complemento muda prazo, valor, fase ou prova, sem transformar repetição de pedido em conteúdo novo."
	right := left
	right.UniqueIntentID = "familia-pensao-revisao-desemprego-negociacao-previa-recurso-decisao-recorrida-prova-nova-rodada-02"
	right.SourceMatrixID = "familia-pensao-revisao-desemprego"
	right.ReaderProblem = "Histórico de tentativa prévia em família trata recurso documentado. A rodada nova muda o eixo do rascunho para fase recursal no subtema revisão de pensão após desemprego e usa o caso matriz apenas como vínculo de origem editorial."
	right.DocumentContext = "Razões do recurso, decisão recorrida, protocolo e prova nova ficam em trilha separada. A documentação da rodada usa eixo documental pensão desemprego, separando prova nova, lacuna corrigida e fase atual sem repetir a lista do caso matriz."
	right.RiskContext = "Tratar recurso como simples novo pedido pode fazer a pessoa perder prazo. Risco operacional pensão desemprego mede se o complemento muda prazo, valor, fase ou prova, sem transformar repetição de pedido em conteúdo novo."

	pair := MaximumPairSimilarityDetail([]Entry{
		{Line: 1, Record: left},
		{Line: 2, Record: right},
	})
	if pair.Score > 0.64 {
		t.Fatalf("expansion boilerplate dominated similarity: pair=%s/%s score=%.4f, want <=0.64", pair.LeftID, pair.RightID, pair.Score)
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
