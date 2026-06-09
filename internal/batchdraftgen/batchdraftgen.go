package batchdraftgen

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"portaljuridico/internal/batchdrafts"
	"portaljuridico/internal/batchsourcematrix"
	"portaljuridico/internal/humanscore"
	"portaljuridico/internal/scalablebatches"
	"portaljuridico/internal/storage"
)

type Options struct {
	SamplesPerBatch int
	CheckedAt       string
}

type Metric struct {
	BatchID                    string  `json:"batch_id"`
	GenerationStatus           string  `json:"generation_status"`
	GeneratedSamples           int     `json:"generated_samples"`
	PassedSamples              int     `json:"passed_samples"`
	RewrittenSamples           int     `json:"rewritten_samples"`
	MinimumHumanScore          int     `json:"minimum_human_score"`
	MaximumAILikeScore         int     `json:"maximum_ai_like_score"`
	MaximumSimilarity          float64 `json:"maximum_similarity"`
	SourceMatrixCoveredSamples int     `json:"source_matrix_covered_samples"`
	StructuralPatternRisk      float64 `json:"structural_pattern_risk"`
	LabEstimatedCPUUnits       int     `json:"lab_estimated_cpu_units"`
	RenderAllowed              bool    `json:"render_allowed"`
	SitemapAllowed             bool    `json:"sitemap_allowed"`
	PublicationAllowed         bool    `json:"publication_allowed"`
	PublicPath                 string  `json:"public_path"`
	CheckedAt                  string  `json:"checked_at"`
	NextValidationAction       string  `json:"next_validation_action"`
}

type MetricEntry struct {
	Line   int
	Metric Metric
}

type Result struct {
	Drafts  []batchdrafts.Record
	Metrics []Metric
}

type Issue struct {
	Code    string
	Message string
}

type Report struct {
	Issues []Issue
}

type scenario struct {
	IntentID        string
	Term            string
	ReaderProblem   string
	SourceHook      string
	DocumentContext string
	RiskContext     string
	DigitalAction   string
}

func DefaultOptions() Options {
	return Options{SamplesPerBatch: 10, CheckedAt: "2026-06-09"}
}

func Generate(root string, options Options) (Result, Report) {
	if options.SamplesPerBatch < 3 {
		return Result{}, Report{Issues: []Issue{{Code: "generation_samples_too_low", Message: fmt.Sprintf("samples=%d", options.SamplesPerBatch)}}}
	}
	if options.CheckedAt == "" {
		return Result{}, Report{Issues: []Issue{{Code: "generation_checked_at_missing", Message: "checked_at vazio"}}}
	}

	entries, loadReport := scalablebatches.LoadRecords(root)
	if !loadReport.Passed() {
		return Result{}, convertBatchReport(loadReport)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Record.BatchID < entries[j].Record.BatchID
	})

	issues := make([]Issue, 0)
	drafts := make([]batchdrafts.Record, 0, len(entries)*options.SamplesPerBatch)
	for _, entry := range entries {
		if validation := scalablebatches.ValidateRecord(entry.Record); !validation.Passed() {
			for _, issue := range validation.Issues {
				issues = append(issues, Issue{Code: "generation_invalid_batch", Message: issue.Code + ":" + issue.Message})
			}
			continue
		}
		scenarios := expandScenarios(entry.Record.LegalArea, scenariosForArea(entry.Record.LegalArea), options.SamplesPerBatch)
		if len(scenarios) < options.SamplesPerBatch {
			issues = append(issues, Issue{Code: "generation_scenario_bank_too_small", Message: entry.Record.LegalArea})
			continue
		}
		for i := 0; i < options.SamplesPerBatch; i++ {
			drafts = append(drafts, buildDraft(entry.Record, scenarios[i], options.CheckedAt))
		}
	}
	if len(issues) > 0 {
		return Result{Drafts: drafts}, Report{Issues: issues}
	}

	result := Result{Drafts: drafts, Metrics: buildMetrics(drafts, options.CheckedAt)}
	return result, ValidateResultWithRoot(root, result)
}

func ValidateResult(result Result) Report {
	return validateResult(result, nil)
}

func ValidateResultWithRoot(root string, result Result) Report {
	matrixEntries, matrixReport := batchsourcematrix.LoadRecords(root)
	if !matrixReport.Passed() {
		return validateResult(result, nil)
	}
	return validateResult(result, matrixEntries)
}

func validateResult(result Result, matrixEntries []batchsourcematrix.Entry) Report {
	issues := make([]Issue, 0)
	if len(result.Drafts) < 30 {
		issues = append(issues, Issue{Code: "generation_too_few_drafts", Message: fmt.Sprintf("drafts=%d", len(result.Drafts))})
	}
	seen := make(map[string]bool)
	for _, draft := range result.Drafts {
		if seen[draft.UniqueIntentID] {
			issues = append(issues, Issue{Code: "generation_duplicate_intent", Message: draft.UniqueIntentID})
		}
		seen[draft.UniqueIntentID] = true
		validation := batchdrafts.ValidateRecord(draft)
		for _, issue := range validation.Issues {
			issues = append(issues, Issue{Code: "generation_invalid_draft", Message: draft.UniqueIntentID + ":" + issue.Code + ":" + issue.Message})
		}
		if !ContainsOrigin(draft.CTAContext, draft.UniqueIntentID) {
			issues = append(issues, Issue{Code: "generation_cta_without_origin", Message: draft.UniqueIntentID})
		}
	}
	if result.RewrittenCount() < 6 {
		issues = append(issues, Issue{Code: "generation_too_few_rewrites", Message: fmt.Sprintf("rewritten=%d", result.RewrittenCount())})
	}
	if result.MaximumPairSimilarity() > 0.64 {
		issues = append(issues, Issue{Code: "generation_similarity_too_high", Message: fmt.Sprintf("max=%.2f", result.MaximumPairSimilarity())})
	}
	if result.StructuralPatternRisk() > 0.35 {
		issues = append(issues, Issue{Code: "generation_structural_pattern_risk_too_high", Message: fmt.Sprintf("risk=%.2f", result.StructuralPatternRisk())})
	}
	if len(matrixEntries) > 0 {
		coverage := batchsourcematrix.ValidateDraftCoverage(matrixEntries, result.Drafts)
		for _, issue := range coverage.Issues {
			issues = append(issues, Issue{Code: "generation_source_matrix_gap", Message: issue.Code + ":" + issue.Message})
		}
	}
	metricReport := ValidateMetrics(result.Metrics)
	issues = append(issues, metricReport.Issues...)
	return Report{Issues: issues}
}

func ValidateMetric(metric Metric) Report {
	issues := make([]Issue, 0)
	if metric.BatchID == "" {
		issues = append(issues, Issue{Code: "generation_metric_missing_batch", Message: "batch_id vazio"})
	}
	if metric.GenerationStatus != "batch_generation_scored_blocked" {
		issues = append(issues, Issue{Code: "generation_metric_invalid_status", Message: metric.GenerationStatus})
	}
	if metric.GeneratedSamples < 5 {
		issues = append(issues, Issue{Code: "generation_metric_too_few_samples", Message: fmt.Sprintf("generated=%d", metric.GeneratedSamples)})
	}
	if metric.PassedSamples != metric.GeneratedSamples || metric.PassedSamples == 0 {
		issues = append(issues, Issue{Code: "generation_metric_unpassed_samples", Message: fmt.Sprintf("passed=%d generated=%d", metric.PassedSamples, metric.GeneratedSamples)})
	}
	if metric.MinimumHumanScore < 85 {
		issues = append(issues, Issue{Code: "generation_metric_human_score_too_low", Message: fmt.Sprintf("min=%d", metric.MinimumHumanScore)})
	}
	if metric.MaximumAILikeScore > 20 {
		issues = append(issues, Issue{Code: "generation_metric_ai_score_too_high", Message: fmt.Sprintf("max=%d", metric.MaximumAILikeScore)})
	}
	if metric.MaximumSimilarity <= 0 || metric.MaximumSimilarity > 0.64 {
		issues = append(issues, Issue{Code: "generation_metric_similarity_too_high", Message: fmt.Sprintf("max=%.2f", metric.MaximumSimilarity)})
	}
	if metric.SourceMatrixCoveredSamples != metric.GeneratedSamples {
		issues = append(issues, Issue{Code: "generation_metric_source_matrix_gap", Message: fmt.Sprintf("covered=%d generated=%d", metric.SourceMatrixCoveredSamples, metric.GeneratedSamples)})
	}
	if metric.StructuralPatternRisk > 0.35 {
		issues = append(issues, Issue{Code: "generation_metric_structural_risk_too_high", Message: fmt.Sprintf("risk=%.2f", metric.StructuralPatternRisk)})
	}
	if metric.LabEstimatedCPUUnits <= 0 {
		issues = append(issues, Issue{Code: "generation_metric_cpu_estimate_missing", Message: metric.BatchID})
	}
	if metric.RenderAllowed {
		issues = append(issues, Issue{Code: "generation_metric_render_allowed", Message: metric.BatchID})
	}
	if metric.SitemapAllowed {
		issues = append(issues, Issue{Code: "generation_metric_sitemap_allowed", Message: metric.BatchID})
	}
	if metric.PublicationAllowed {
		issues = append(issues, Issue{Code: "generation_metric_publication_allowed", Message: metric.BatchID})
	}
	if metric.PublicPath != "" {
		issues = append(issues, Issue{Code: "generation_metric_has_public_path", Message: metric.PublicPath})
	}
	if metric.CheckedAt == "" {
		issues = append(issues, Issue{Code: "generation_metric_without_checked_at", Message: metric.BatchID})
	}
	return Report{Issues: issues}
}

func ValidateMetrics(metrics []Metric) Report {
	issues := make([]Issue, 0)
	if len(metrics) < 6 {
		issues = append(issues, Issue{Code: "generation_metrics_too_few_batches", Message: fmt.Sprintf("metrics=%d", len(metrics))})
	}
	seen := make(map[string]bool)
	for _, metric := range metrics {
		if seen[metric.BatchID] {
			issues = append(issues, Issue{Code: "generation_metric_duplicate_batch", Message: metric.BatchID})
		}
		seen[metric.BatchID] = true
		report := ValidateMetric(metric)
		issues = append(issues, report.Issues...)
	}
	return Report{Issues: issues}
}

func ValidateStoredMetrics(root string) Report {
	entries, report := LoadMetrics(root)
	if !report.Passed() {
		return report
	}
	metrics := make([]Metric, 0, len(entries))
	for _, entry := range entries {
		metrics = append(metrics, entry.Metric)
	}
	return ValidateMetrics(metrics)
}

func LoadMetrics(root string) ([]MetricEntry, Report) {
	projectRoot, err := findProjectRoot(root)
	if err != nil {
		return nil, Report{Issues: []Issue{{Code: "project_root_not_found", Message: err.Error()}}}
	}
	path := filepath.Join(projectRoot, "data", "editorial", "batch_generation_metrics.jsonl")
	file, err := os.Open(path)
	if err != nil {
		return nil, Report{Issues: []Issue{{Code: "generation_metrics_missing", Message: err.Error()}}}
	}
	defer file.Close()

	entries := make([]MetricEntry, 0)
	issues := make([]Issue, 0)
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var metric Metric
		if err := json.Unmarshal([]byte(line), &metric); err != nil {
			issues = append(issues, Issue{Code: "generation_metric_invalid_json", Message: fmt.Sprintf("line=%d %s", lineNumber, err.Error())})
			continue
		}
		entries = append(entries, MetricEntry{Line: lineNumber, Metric: metric})
	}
	if err := scanner.Err(); err != nil {
		issues = append(issues, Issue{Code: "generation_metric_scan_failed", Message: err.Error()})
	}
	return entries, Report{Issues: issues}
}

func WriteMetrics(root string, metrics []Metric) error {
	contract, err := storage.LoadContract(root)
	if err != nil {
		return err
	}
	layer, ok := contract.LayerByName("batch_generation_metrics")
	if !ok {
		return fmt.Errorf("unknown_layer=batch_generation_metrics")
	}
	path := filepath.Join(contract.Root, layer.Path)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, metric := range metrics {
		if report := ValidateMetric(metric); !report.Passed() {
			return fmt.Errorf("invalid_generation_metric=%s", strings.Join(report.Messages(), " | "))
		}
		data, err := json.Marshal(metric)
		if err != nil {
			return err
		}
		if len(data)+1 > layer.RecordMaxBytes {
			return fmt.Errorf("record_too_large=batch_generation_metrics bytes=%d max=%d", len(data)+1, layer.RecordMaxBytes)
		}
		if _, err := file.Write(append(data, '\n')); err != nil {
			return err
		}
	}
	return nil
}

func (r Result) DraftIDs() []string {
	ids := make([]string, 0, len(r.Drafts))
	for _, draft := range r.Drafts {
		ids = append(ids, draft.UniqueIntentID)
	}
	return ids
}

func (r Result) RewrittenCount() int {
	count := 0
	for _, draft := range r.Drafts {
		if draft.RewriteStatus == "rewritten_after_score_failure" {
			count++
		}
	}
	return count
}

func (r Result) MaximumPairSimilarity() float64 {
	entries := make([]batchdrafts.Entry, 0, len(r.Drafts))
	for i, draft := range r.Drafts {
		entries = append(entries, batchdrafts.Entry{Line: i + 1, Record: draft})
	}
	return batchdrafts.MaximumPairSimilarity(entries)
}

func (r Result) StructuralPatternRisk() float64 {
	return structuralPatternRisk(r.Drafts)
}

func ContainsOrigin(value string, uniqueIntentID string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(uniqueIntentID))
}

func buildDraft(batch scalablebatches.Record, scenario scenario, checkedAt string) batchdrafts.Record {
	record := batchdrafts.Record{
		BatchID:            batch.BatchID,
		UniqueIntentID:     batch.LegalArea + "-" + scenario.IntentID,
		SourceMatrixID:     batch.LegalArea + "-" + sourceMatrixID(scenario),
		LegalArea:          batch.LegalArea,
		DraftStatus:        "batch_draft_scored_blocked",
		Language:           "pt-BR",
		Term:               scenario.Term,
		ReaderProblem:      scenario.ReaderProblem,
		SourceHook:         scenario.SourceHook,
		DocumentContext:    scenario.DocumentContext,
		RiskContext:        scenario.RiskContext,
		DigitalAction:      scenario.DigitalAction,
		SourceFamilies:     append([]string{}, batch.SourceFamilies...),
		RewriteStatus:      "rewritten_after_score_failure",
		InitialIssueCodes:  initialIssueCodes(scenario.Term),
		RewriteAttempts:    1,
		RenderAllowed:      false,
		SitemapAllowed:     false,
		PublicationAllowed: false,
		PublicPath:         "",
		CheckedAt:          checkedAt,
	}
	sourceFamily := ""
	if len(batch.SourceFamilies) > 0 {
		sourceFamily = batch.SourceFamilies[0]
	}
	record.CTAContext = strings.ReplaceAll(batch.CTAContextTemplate, "{unique_intent_id}", record.UniqueIntentID)
	record.CTAContext = strings.ReplaceAll(record.CTAContext, "{source_family}", sourceFamily)
	score := humanscore.ScoreText(record.FullText())
	record.HumanScore = score.HumanScore
	record.AILikeScore = score.AILikeScore
	return record
}

func initialIssueCodes(term string) []string {
	text := "Este conteúdo explica " + term + " de forma geral. Documentos necessários. Quando procurar advogado. " + strings.Repeat(term+" ", 8)
	score := humanscore.ScoreText(text)
	codes := score.Codes()
	if len(codes) == 0 {
		return []string{"ai_like_generic_markers"}
	}
	return codes
}

func buildMetrics(drafts []batchdrafts.Record, checkedAt string) []Metric {
	byBatch := make(map[string][]batchdrafts.Record)
	order := make([]string, 0)
	for _, draft := range drafts {
		if _, ok := byBatch[draft.BatchID]; !ok {
			order = append(order, draft.BatchID)
		}
		byBatch[draft.BatchID] = append(byBatch[draft.BatchID], draft)
	}
	sort.Strings(order)
	metrics := make([]Metric, 0, len(order))
	for _, batchID := range order {
		items := byBatch[batchID]
		minHuman := 100
		maxAI := 0
		passed := 0
		rewritten := 0
		entries := make([]batchdrafts.Entry, 0, len(items))
		for i, draft := range items {
			if draft.HumanScore < minHuman {
				minHuman = draft.HumanScore
			}
			if draft.AILikeScore > maxAI {
				maxAI = draft.AILikeScore
			}
			if batchdrafts.ValidateRecord(draft).Passed() {
				passed++
			}
			if draft.RewriteStatus == "rewritten_after_score_failure" {
				rewritten++
			}
			entries = append(entries, batchdrafts.Entry{Line: i + 1, Record: draft})
		}
		metrics = append(metrics, Metric{
			BatchID:                    batchID,
			GenerationStatus:           "batch_generation_scored_blocked",
			GeneratedSamples:           len(items),
			PassedSamples:              passed,
			RewrittenSamples:           rewritten,
			MinimumHumanScore:          minHuman,
			MaximumAILikeScore:         maxAI,
			MaximumSimilarity:          batchdrafts.MaximumPairSimilarity(entries),
			SourceMatrixCoveredSamples: sourceMatrixCoveredSamples(items),
			StructuralPatternRisk:      structuralPatternRisk(items),
			LabEstimatedCPUUnits:       len(items) * 12,
			RenderAllowed:              false,
			SitemapAllowed:             false,
			PublicationAllowed:         false,
			PublicPath:                 "",
			CheckedAt:                  checkedAt,
			NextValidationAction:       "aumentar lote com semantica por subtema, matriz de fonte especifica, similaridade por familia e publicacao bloqueada ate revisao e SEO passarem",
		})
	}
	return metrics
}

func sourceMatrixID(scenario scenario) string {
	for _, suffix := range []string{"-documentos-prazo", "-revisao-fonte", "-triagem-risco"} {
		if strings.HasSuffix(scenario.IntentID, suffix) {
			return strings.TrimSuffix(scenario.IntentID, suffix)
		}
	}
	return scenario.IntentID
}

func sourceMatrixCoveredSamples(items []batchdrafts.Record) int {
	count := 0
	for _, item := range items {
		if item.SourceMatrixID != "" {
			count++
		}
	}
	return count
}

func structuralPatternRisk(items []batchdrafts.Record) float64 {
	if len(items) == 0 {
		return 0
	}
	counts := make(map[string]int)
	max := 0
	for _, item := range items {
		signature := patternSignature(item.ReaderProblem)
		counts[signature]++
		if counts[signature] > max {
			max = counts[signature]
		}
	}
	return float64(max) / float64(len(items))
}

func patternSignature(value string) string {
	words := strings.Fields(strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return ' '
	}, value))
	signals := make([]string, 0, 3)
	stop := map[string]bool{"o": true, "a": true, "os": true, "as": true, "um": true, "uma": true, "de": true, "do": true, "da": true, "e": true, "mas": true, "com": true, "para": true, "precisa": true}
	for _, word := range words {
		if !stop[word] {
			signals = append(signals, word)
		}
		if len(signals) == 3 {
			break
		}
	}
	if len(signals) == 0 {
		return "empty"
	}
	return strings.Join(signals, "-")
}

func expandScenarios(area string, base []scenario, count int) []scenario {
	if count <= len(base) {
		return base[:count]
	}
	expanded := append([]scenario{}, base...)
	variant := 1
	for len(expanded) < count {
		for _, item := range base {
			if len(expanded) >= count {
				break
			}
			expanded = append(expanded, scenarioVariant(area, item, variant))
		}
		variant++
	}
	return expanded
}

func scenarioVariant(area string, base scenario, variant int) scenario {
	switch variant % 3 {
	case 1:
		return scenario{
			IntentID:        base.IntentID + "-documentos-prazo",
			Term:            base.Term + " com prova digital organizada",
			ReaderProblem:   "Para " + strings.ReplaceAll(base.IntentID, "-", " e ") + ", o leitor precisa transformar a dúvida em linha do tempo, com datas, protocolo, resposta formal e documento que prove o ponto sensível do caso.",
			SourceHook:      "A fonte de apoio continua específica ao subtema: " + base.SourceHook,
			DocumentContext: "A segunda leitura usa os mesmos documentos centrais do caso, mas cobra ordem cronológica, comprovante digital preservado e identificação clara de quem respondeu: " + base.DocumentContext,
			RiskContext:     "O risco muda quando a prova está incompleta: " + base.RiskContext,
			DigitalAction:   "O atendimento online recebe origem, subtema e arquivos essenciais para decidir se falta documento, fonte ou revisão antes de qualquer página pública.",
		}
	case 2:
		return scenario{
			IntentID:        base.IntentID + "-revisao-fonte",
			Term:            base.Term + " com revisão de fonte oficial",
			ReaderProblem:   "O usuário quer entender se a fonte oficial realmente cobre o seu cenário, sem depender de notícia, promessa comercial ou modelo genérico de página.",
			SourceHook:      base.SourceHook + " A variação exige comparação entre fonte primária, orientação administrativa e documentos do caso.",
			DocumentContext: "A leitura pede decisão, protocolo, contrato, comprovantes, mensagens, identificação das partes e documento que conecte o fato ao subtema.",
			RiskContext:     "A cautela principal é não transformar dúvida comum em página pública antes de provar especificidade, utilidade e diferença real de intenção.",
			DigitalAction:   "O fluxo digital registra origem, subtema e documentos esperados para que o WhatsApp receba contexto útil de contratação remota.",
		}
	default:
		return scenario{
			IntentID:        base.IntentID + "-triagem-risco",
			Term:            base.Term + " para triagem jurídica online",
			ReaderProblem:   "A dúvida já tem urgência prática, porém precisa de triagem cuidadosa para separar documento essencial, risco jurídico, prazo e fonte oficial.",
			SourceHook:      base.SourceHook + " O foco é confirmar se o subtema possui fonte suficiente para avançar no laboratório.",
			DocumentContext: "São separados documentos obrigatórios, registros auxiliares, resposta do órgão ou empresa, comprovantes financeiros e comunicações preservadas.",
			RiskContext:     "O risco é publicar texto parecido com outros do lote ou prometer solução sem validar a prova mínima do caso.",
			DigitalAction:   "A etapa online mede score, similaridade, fonte e CTA contextual antes de qualquer render público.",
		}
	}
}

func scenariosForArea(area string) []scenario {
	switch area {
	case "saude-suplementar":
		return []scenario{
			{"procedimento-urgente-negado", "procedimento urgente negado pelo plano de saúde", "O leitor recebeu resposta negativa da operadora para cirurgia indicada como urgente e precisa separar contrato, relatório médico, protocolo, segmentação do plano, prazo da ANS e risco clínico antes de agir.", "A pauta usa Lei 9656, Rol da ANS, regra de cobertura assistencial, canais oficiais de reclamação e o motivo escrito pela operadora.", "A triagem pede pedido médico, relatório clínico, resposta formal, protocolo, contrato, carteirinha, exames recentes e mensagens do aplicativo.", "O risco é a demora agravar o quadro ou o usuário aceitar justificativa incompleta sem preservar prova da urgência e da negativa.", "O atendimento digital organiza arquivos, confirma datas, identifica fonte aplicável e prepara perguntas objetivas para avaliação jurídica remota."},
			{"medicamento-prescrito-recusado", "medicamento prescrito recusado pelo plano", "A família tem prescrição atual, recebeu recusa do plano e precisa entender indicação médica, cobertura, Rol da ANS, orçamento, contrato e justificativa técnica da negativa.", "A análise parte das normas da ANS, legislação dos planos, atualização do Rol, protocolo administrativo e resposta emitida pela operadora.", "São esperados prescrição, relatório do especialista, CID quando constar, orçamento, e-mail de recusa, protocolo e cópia do contrato.", "A cautela é não prometer fornecimento imediato e diferenciar urgência documentada, escolha terapêutica e exclusão contratual alegada.", "A triagem online confere a documentação, monta linha do tempo e orienta se falta reclamação administrativa, nova prova médica ou análise judicial."},
			{"reembolso-exame-fora-rede", "reembolso de exame fora da rede credenciada", "O usuário pagou exame porque não conseguiu agenda na rede e precisa verificar prazo de atendimento, nota fiscal, pedido médico, contrato e resposta da operadora.", "A pauta cruza canais da ANS, contrato assistencial, regra de rede credenciada, comprovante de indisponibilidade e registro de atendimento.", "A conferência exige nota fiscal, comprovante de pagamento, pedido médico, protocolo, justificativa da rede indisponível e data do exame.", "O risco é confundir escolha particular com indisponibilidade real da rede ou perder documento que comprove tentativa de atendimento.", "A análise remota recebe os arquivos, compara datas, identifica lacunas e registra origem do caso para WhatsApp com contexto."},
			{"home-care-reduzido", "home care reduzido ou suspenso pelo plano", "A família recebeu redução de assistência domiciliar e precisa organizar relatório clínico, contrato, prescrição de enfermagem, resposta da operadora e evolução do paciente.", "A fonte considera regras da ANS, legislação de planos, contrato de cobertura, relatório médico e histórico de autorização anterior.", "Devem ser enviados relatório clínico, prescrição, autorizações antigas, comunicado de redução, protocolos, contrato e exames que mostrem dependência.", "A cautela é diferenciar ajuste técnico de suspensão abusiva e evitar promessa de manutenção automática sem prova médica robusta.", "A triagem digital monta sequência de eventos, confere urgência e prepara leitura jurídica antes de medida administrativa ou judicial."},
			{"prazo-consulta-especialista", "prazo excessivo para consulta com especialista", "O beneficiário tenta agendar especialista, recebe data distante e precisa provar tentativas, protocolo, rede disponível, urgência e prazo máximo aplicável.", "A análise usa prazos da ANS, canais de reclamação, contrato do plano, rede credenciada e documentos médicos que indicam prioridade.", "A triagem pede protocolos, prints de agenda, pedido médico, carteirinha, contrato, resposta da central e eventual relatório clínico.", "O risco é reclamar sem demonstrar tentativa real de agendamento ou sem identificar se o caso é urgência, retorno ou consulta eletiva.", "O fluxo online organiza evidências, registra datas e prepara mensagem contextual de WhatsApp para avaliar providência proporcional."},
		}
	case "trabalhista":
		return []scenario{
			{"rescisao-indireta-assedio-salario", "rescisão indireta por assédio e atraso salarial", "O empregado relata cobrança humilhante, atraso de salário e medo de pedir demissão sem avaliar contrato, holerites, mensagens, testemunhas e datas relevantes.", "A pauta usa CLT, orientação trabalhista, documentos rescisórios, controle de jornada, prova digital e cautela sobre permanência no emprego.", "A triagem pede contrato, holerites, mensagens, advertências, comprovantes de atraso, escala, ponto e nomes de testemunhas.", "O risco é romper o vínculo sem prova suficiente ou perder verbas por uma narrativa sem datas e documentos concretos.", "O atendimento online organiza linha do tempo, preserva prints e indica perguntas para análise de reclamação trabalhista digital."},
			{"horas-extras-banco-irregular", "horas extras e compensação irregular de jornada", "O trabalhador percebe jornada maior que a registrada e precisa comparar CLT, cartões de ponto, escala, mensagens, banco de horas, holerites, contrato e períodos prescricionais.", "A análise parte da CLT, informações oficiais sobre jornada, regras de compensação, testemunhas e documentos que mostram rotina efetiva.", "São conferidos ponto, escala, mensagens, e-mails, holerites, contrato, recibos de compensação e prints preservados.", "A cautela é separar ocorrência eventual de jornada habitual e não calcular verbas sem base documental completa.", "A triagem remota recebe arquivos, agrupa períodos, revisa testemunhas e identifica se há perguntas objetivas para ação trabalhista."},
			{"justa-causa-prova-fragil", "contestação de justa causa com prova frágil", "A pessoa foi dispensada por justa causa e precisa avaliar comunicado, proporcionalidade, advertências anteriores, holerites, testemunhas e verbas retidas.", "A pauta considera CLT, documentos rescisórios, histórico disciplinar, prova empresarial e critérios de proporcionalidade.", "Devem ser enviados termo de rescisão, comunicado da falta, advertências, mensagens, controle de ponto, holerites e prova apresentada pela empresa.", "O risco é assinar documento sem ressalva, perder prazo ou ignorar prova que confirme a falta alegada.", "O fluxo digital revisa datas, arquivos e perguntas essenciais antes de orientar contestação ou negociação."},
			{"acidente-trabalho-estabilidade", "acidente de trabalho e estabilidade após afastamento", "O empregado voltou de afastamento, teme dispensa e precisa organizar CAT, laudos, atestados, benefício do INSS, contrato e comunicação da empresa.", "A análise combina CLT, documentos previdenciários, prova médica, comunicação de acidente e histórico de retorno ao trabalho.", "A triagem pede CAT, laudos, exames, atestados, comunicado de decisão do INSS, holerites, contrato e mensagens com a empresa.", "A cautela é diferenciar afastamento comum, acidente laboral e doença ocupacional sem criar promessa de reintegração automática.", "O atendimento online monta cronologia, confere benefício e registra origem para análise trabalhista e previdenciária conectada."},
			{"verbas-rescisorias-nao-pagas", "verbas rescisórias não pagas no prazo", "O ex-empregado saiu da empresa, não recebeu tudo no prazo e precisa conferir TRCT, chave do FGTS, holerites, aviso prévio e comprovantes.", "A pauta usa CLT, documentos rescisórios, guias, prazos de pagamento e prova de comunicação entre empresa e trabalhador.", "São esperados TRCT, termo de homologação quando houver, comprovante de pagamento, extrato de FGTS, holerites e mensagens.", "O risco é aceitar cálculo incompleto ou acionar sem separar atraso, diferença de valor e verba ainda controvertida.", "A triagem digital recebe documentos, organiza valores por tipo e prepara roteiro para análise jurídica remota."},
		}
	case "familia":
		return []scenario{
			{"divorcio-consensual-filhos-bens", "divórcio consensual online com filhos e bens", "O casal concorda com o divórcio, mas precisa tratar guarda, convivência, pensão, partilha de bens, renda e documentos dos filhos sem confundir cartório com processo judicial.", "A orientação considera Código Civil, regras processuais, CNJ, cartório quando cabível e diferença entre consenso formal e conflito pendente.", "A triagem pede certidão de casamento, documentos dos filhos, comprovantes de renda, acordo sobre bens, endereço das partes e despesas da criança.", "O risco é tentar via inadequada quando há filhos menores, desacordo oculto ou partilha sem documentação do patrimônio.", "O fluxo online organiza pontos de consenso, documentos e pendências para análise de advogado antes da minuta."},
			{"pensao-revisao-desemprego", "revisão de pensão alimentícia após desemprego", "Quem paga pensão perdeu renda e precisa demonstrar mudança real, sentença anterior, acordo, despesas da criança, extratos e capacidade financeira atual.", "A pauta usa legislação de família, decisões sobre alimentos, prova de renda e cautela para não estimular interrupção unilateral.", "São esperados sentença ou acordo, comprovantes de desemprego, renda atual, despesas do filho, recibos escolares e gastos médicos.", "A cautela é não parar pagamento por conta própria e diferenciar dificuldade temporária de alteração relevante para revisão.", "A triagem remota monta quadro de renda, despesas, datas e documentos para avaliar revisão, negociação ou cumprimento."},
			{"guarda-distancia-escola", "guarda e convivência quando os pais moram longe", "Os pais vivem em cidades diferentes e precisam organizar residência, escola, férias, custos de viagem, calendário de convivência e comunicação com a criança.", "A análise considera legislação de família, melhor interesse da criança, CNJ, prova da rotina familiar e acordos anteriores.", "A triagem pede certidão, comprovante de residência, calendário escolar, mensagens, despesas de deslocamento e proposta de convivência.", "O risco é transformar disputa logística em conflito maior sem demonstrar rotina, cooperação e impacto real para a criança.", "O atendimento online recebe documentos, organiza calendário e identifica pontos de acordo antes de ação ou homologação."},
			{"alimentos-avoengos", "pensão avoenga quando os pais não conseguem pagar", "A família cogita pedir alimentos aos avós e precisa avaliar renda dos pais, necessidade do filho, despesas, tentativa anterior e situação dos avós.", "A pauta parte do Código Civil, entendimento sobre obrigação complementar e prova financeira de todos os envolvidos.", "São conferidos comprovantes de renda, despesas da criança, decisão anterior, mensagens, documentos dos avós e histórico de pagamento.", "A cautela é não tratar avós como primeira opção automática nem ignorar capacidade dos pais.", "A triagem digital organiza documentos, identifica lacunas e prepara perguntas para análise de família sem promessa de valor."},
			{"partilha-bens-conta-digital", "partilha de bens com conta digital e investimentos", "No fim do casamento, uma parte suspeita de patrimônio oculto em conta digital e precisa organizar extratos, imposto de renda, bens, contrato e datas de aquisição.", "A análise considera regime de bens, Código Civil, documentos financeiros, declaração fiscal e prova de movimentação patrimonial.", "A triagem pede certidão, pacto quando houver, extratos, informes de rendimento, imposto de renda, matrícula de imóvel e comprovantes.", "O risco é alegar ocultação sem prova mínima ou deixar fora bens adquiridos no período correto.", "O atendimento online organiza ativos, datas e documentos para orientar estratégia de divórcio ou sobrepartilha."},
		}
	case "previdenciario":
		return []scenario{
			{"auxilio-incapacidade-pericia", "auxílio por incapacidade negado após perícia", "O segurado recebeu indeferimento depois da perícia e precisa comparar laudos, exames, atestados, CNIS, qualidade de segurado e atividade profissional.", "A pauta usa canais do INSS, Previdência, legislação de benefícios, comunicado de decisão e exigências documentais.", "A triagem pede decisão, laudos, exames, atestados, CNIS, carteira de trabalho, protocolos e histórico médico organizado por data.", "O risco é perder prazo de recurso ou discutir incapacidade sem demonstrar relação entre doença e trabalho exercido.", "O atendimento digital confere documentos, datas e lacunas para avaliar recurso administrativo ou ação judicial."},
			{"bpc-loas-cadunico-renda", "BPC LOAS negado por renda familiar", "A família recebeu negativa do BPC e precisa revisar CadÚnico, composição familiar, laudos, despesas de saúde, renda real e documentos da casa.", "A análise parte de regras assistenciais, INSS, Previdência, cadastro social, critérios administrativos e prova de vulnerabilidade.", "São esperados comunicado de decisão, CadÚnico, laudos, comprovantes de renda, despesas médicas, residência e documentos familiares.", "A cautela é discutir apenas renda bruta e esquecer impedimentos, gastos essenciais ou composição familiar correta.", "A triagem remota organiza documentos e perguntas para recurso, nova solicitação ou medida judicial proporcional."},
			{"cumprimento-exigencia-parado", "cumprimento de exigência do INSS sem resposta", "O segurado enviou documentos pelo Meu INSS, o pedido ficou parado e precisa comprovar protocolo, carta de exigência, envio, CNIS e histórico do benefício.", "A pauta usa canais oficiais do INSS, protocolo administrativo, prazos de análise e documento exigido no processo.", "A triagem pede print do protocolo, carta de exigência, comprovante de envio, CNIS, documentos pessoais e número do requerimento.", "O risco é tratar atraso como indeferimento ou ajuizar sem confirmar se a exigência foi realmente cumprida.", "O fluxo online revisa arquivo por arquivo, identifica pendência e orienta cobrança administrativa ou providência jurídica."},
			{"aposentadoria-cnis-incompleto", "CNIS incompleto antes de pedir aposentadoria", "O trabalhador percebe vínculos faltando no CNIS e precisa reunir carteira, contracheques, PPP, carnês, guias e documentos de empresa antiga.", "A análise parte do INSS, Previdência, cadastro de vínculos, prova de contribuição e exigências para acerto cadastral.", "São conferidos CNIS, carteira de trabalho, holerites, guias, contratos, PPP, certidões e protocolos anteriores.", "A cautela é protocolar aposentadoria com tempo incompleto e receber indeferimento evitável por falta de acerto prévio.", "A triagem digital organiza vínculos, identifica lacunas e prepara lista documental para correção administrativa."},
			{"beneficio-cessado-laudo", "benefício cessado mesmo com laudo atualizado", "O beneficiário teve auxílio cessado, ainda possui laudo recente e precisa entender perícia, decisão, exames, atestados e atividade laboral.", "A pauta considera INSS, Previdência, comunicação de cessação, prova médica e diferença entre incapacidade parcial, temporária e permanente.", "A triagem pede comunicado, laudos, exames, atestados, receitas, CNIS, função exercida e protocolos de recurso.", "O risco é perder prazo ou apresentar laudo sem conexão clara com limitações no trabalho.", "O atendimento online monta linha do tempo e avalia se cabe recurso, novo pedido ou processo."},
		}
	case "consumidor-financeiro":
		return []scenario{
			{"negativacao-divida-desconhecida", "negativação indevida por dívida desconhecida", "O consumidor descobriu restrição no cadastro e precisa guardar consulta, contrato inexistente ou contestado, protocolos, extratos e resposta do credor.", "A pauta usa defesa do consumidor, Banco Central, canais oficiais, prova documental e histórico de cobrança.", "A triagem pede consulta do cadastro, faturas, contrato se apresentado, protocolos, e-mails, boletim quando houver fraude e extratos.", "O risco é prometer indenização sem separar fraude, erro cadastral, dívida prescrita ou contratação válida.", "O atendimento digital organiza credor, datas, documentos e urgência comercial para análise jurídica remota."},
			{"consignado-nao-reconhecido", "empréstimo consignado não reconhecido", "A pessoa encontrou desconto em benefício ou salário e precisa verificar contrato, depósito, autorização, extratos, protocolos e contestação administrativa.", "A análise cruza Banco Central, INSS quando houver benefício, regras de consumidor e prova bancária do desconto.", "São esperados extrato, contrato apresentado pelo banco, comprovante de depósito, benefício, protocolos e reclamação feita.", "A cautela é diferenciar fraude, contratação eletrônica contestada, refinanciamento e valor efetivamente recebido.", "A triagem remota confere documentos, banco envolvido e datas antes de avaliar cancelamento, devolução ou demanda."},
			{"pix-fraude-resposta-banco", "fraude via Pix e resposta insuficiente do banco", "O cliente sofreu golpe, comunicou o banco e precisa organizar horários, comprovante Pix, protocolo, boletim, resposta da instituição e linha do tempo.", "A pauta considera Banco Central, regras de segurança, comunicação rápida, boletim de ocorrência e prova digital preservada.", "A triagem pede comprovante Pix, extrato, protocolo, conversa com fraudador, boletim, resposta do banco e horário dos contatos.", "O risco é prometer ressarcimento sem avaliar rapidez da comunicação, falha de segurança e conduta do usuário.", "O atendimento online monta cronologia, confere documentos e avalia reclamação administrativa ou medida jurídica cabível."},
			{"cartao-cobranca-nao-reconhecida", "cobrança não reconhecida no cartão", "O consumidor viu compras desconhecidas na fatura e precisa preservar cartão, contestação, protocolos, extratos, bloqueio e resposta da administradora.", "A análise usa Banco Central, defesa do consumidor, regras de contestação, prova de comunicação e documentos da fatura.", "São conferidos fatura, extrato, protocolo, boletim quando houver, e-mails, comprovante de bloqueio e resposta da instituição.", "A cautela é separar fraude, compra recorrente esquecida, contestação fora do prazo e falha de segurança.", "A triagem digital organiza eventos e documentos para avaliar providência administrativa ou judicial."},
			{"tarifa-bancaria-indevida", "tarifa bancária cobrada sem contratação clara", "O cliente identifica tarifa recorrente e precisa conferir contrato, cesta de serviços, extratos, protocolos, histórico de contratação e resposta do banco.", "A pauta usa Banco Central, canais oficiais, regras de relacionamento bancário e prova documental da cobrança.", "A triagem pede extratos, contrato, prints do aplicativo, protocolos, resposta do banco e histórico de alteração de pacote.", "O risco é transformar cobrança contratada em disputa sem prova ou perder repetição mensal relevante.", "O atendimento online recebe documentos, calcula período e prepara perguntas para análise de direito do consumidor."},
		}
	case "sucessorio":
		return []scenario{
			{"inventario-extrajudicial-consenso", "inventário extrajudicial com herdeiros concordes", "A família quer resolver inventário em cartório e precisa confirmar consenso, certidão de óbito, herdeiros, bens, dívidas, imposto e eventual testamento.", "A análise usa Código Civil, regras notariais, CNJ, e-Notariado quando cabível e limites do inventário extrajudicial.", "A triagem pede certidão de óbito, documentos dos herdeiros, matrícula de imóvel, extratos, dívidas, imposto e dados do cartório.", "O risco é prometer prazo curto quando há conflito, incapaz, testamento, dívida complexa ou documento faltante.", "O atendimento digital organiza bens, herdeiros e pendências para avaliar caminho cartorial ou judicial."},
			{"imovel-financiado-partilha", "partilha de imóvel financiado no inventário", "Os herdeiros precisam dividir imóvel financiado e entender saldo devedor, contrato bancário, matrícula, imposto, espólio e acordo de partilha.", "A pauta considera legislação sucessória, registro de imóveis, contrato de financiamento, regras notariais e documentos do bem.", "São esperados matrícula, contrato bancário, saldo devedor, certidão de óbito, documentos dos herdeiros e comprovantes de imposto.", "A cautela é ignorar garantia, dívida ou divergência entre herdeiros e criar partilha inviável.", "A triagem remota revisa documentos, organiza bens e dívidas e indica perguntas para inventário judicial ou extrajudicial."},
			{"alvara-valores-bancarios", "alvará para levantar valores bancários de falecido", "A família encontrou saldo pequeno e precisa saber se alvará pode substituir inventário completo, considerando herdeiros, dependentes, bens e dívidas.", "A análise parte de legislação sucessória, documentos bancários, certidão de óbito, orientações judiciais e limites do levantamento.", "A triagem pede certidão de óbito, extrato bancário, documentos dos herdeiros, declaração de dependentes e informação sobre outros bens.", "O risco é simplificar caso que tem imóvel, disputa, testamento ou dívida relevante.", "O atendimento online confere valores, bens conhecidos e documentos para avaliar alvará, inventário ou outra medida."},
			{"testamento-duvida-validade", "testamento encontrado após a morte", "Os familiares encontraram testamento e precisam entender validade formal, herdeiros necessários, bens, cartório, certidões e possível inventário judicial.", "A pauta usa Código Civil, regras sucessórias, documentos notariais, CNJ e diferença entre testamento público, cerrado ou particular.", "São conferidos testamento, certidão de óbito, documentos dos herdeiros, matrícula de bens, extratos e informações de cartório.", "A cautela é partilhar bens sem confirmar testamento ou presumir invalidade sem análise formal.", "A triagem digital organiza documentos e perguntas para orientar abertura, registro ou discussão no inventário."},
			{"divida-espolio-cobranca", "dívida do espólio cobrada durante inventário", "Os herdeiros receberam cobrança contra o falecido e precisam separar dívida legítima, contrato, extratos, bens do espólio, imposto e responsabilidade patrimonial.", "A análise considera Código Civil, documentos de cobrança, inventário, bens do espólio e limites de responsabilidade dos herdeiros.", "A triagem pede contrato, cobrança, extratos, certidão de óbito, relação de bens, dívidas conhecidas e documentos dos herdeiros.", "O risco é pagar dívida sem verificar origem ou ignorar cobrança que afeta partilha e imposto.", "O atendimento online organiza credores, valores, bens e documentos para análise sucessória remota."},
		}
	default:
		return nil
	}
}

func convertBatchReport(report scalablebatches.Report) Report {
	issues := make([]Issue, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, Issue{Code: "generation_batch_load_failed", Message: issue.Code + ":" + issue.Message})
	}
	return Report{Issues: issues}
}

func (r Report) Passed() bool { return len(r.Issues) == 0 }

func (r Report) HasIssue(code string) bool {
	for _, issue := range r.Issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

func (r Report) Codes() []string {
	codes := make([]string, 0, len(r.Issues))
	for _, issue := range r.Issues {
		codes = append(codes, issue.Code)
	}
	sort.Strings(codes)
	return codes
}

func (r Report) Messages() []string {
	messages := make([]string, 0, len(r.Issues))
	for _, issue := range r.Issues {
		messages = append(messages, issue.Code+": "+issue.Message)
	}
	return messages
}

func findProjectRoot(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", os.ErrNotExist
		}
		current = parent
	}
}
