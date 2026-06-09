package paidintent

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"portaljuridico/internal/batchfinaldrafts"
)

type Mode string

const RequirePaidSignal Mode = "require_paid_signal"

const (
	minimumPaidBusinessScore      = 4
	PassedBlockedStatus           = "paid_intent_passed_blocked_publication"
	PublicAssistanceBlockedStatus = "paid_intent_blocked_public_assistance_free_risk"
	AdminSelfServiceBlockedStatus = "paid_intent_blocked_admin_self_service_risk"
)

type Record struct {
	PaidIntentGateID   string   `json:"paid_intent_gate_id"`
	DraftID            string   `json:"draft_id"`
	UniqueIntentID     string   `json:"unique_intent_id"`
	BatchID            string   `json:"batch_id"`
	SourceMatrixID     string   `json:"source_matrix_id"`
	Term               string   `json:"term"`
	CandidatePath      string   `json:"candidate_path"`
	PaidIntentStatus   string   `json:"paid_intent_status"`
	PaidBusinessScore  int      `json:"paid_business_score"`
	PaidSignals        []string `json:"paid_signals"`
	BusinessSignals    []string `json:"business_signals"`
	RiskSignals        []string `json:"risk_signals"`
	RoutingDecision    string   `json:"routing_decision"`
	IndexPolicy        string   `json:"index_policy"`
	RenderAllowed      bool     `json:"render_allowed"`
	SitemapAllowed     bool     `json:"sitemap_allowed"`
	PublicationAllowed bool     `json:"publication_allowed"`
	PublicPath         string   `json:"public_path"`
	CheckedAt          string   `json:"checked_at"`
}

type Entry struct {
	Line   int
	Record Record
}

type Issue struct {
	Code    string
	Message string
}

type Report struct {
	Issues []Issue
}

type textScore struct {
	paidSignals             []string
	businessSignals         []string
	freeSignals             []string
	researchSignals         []string
	publicAssistanceSignals []string
	adminSelfServiceSignals []string
	paidBusinessScore       int
}

var paidServiceSignals = []string{
	"contratacao particular",
	"contratar advogado",
	"contratar advogada",
	"advogado particular",
	"advogada particular",
	"atendimento particular",
	"consulta paga",
	"triagem paga",
	"analise paga",
	"servico juridico pago",
	"honorarios",
	"orcamento de honorarios",
	"cliente particular",
}

var businessReadinessSignals = []string{
	"valor envolvido",
	"valor da causa",
	"valor pago",
	"prejuizo",
	"patrimonio",
	"empresa",
	"contrato",
	"plano de saude",
	"banco",
	"inss",
	"beneficio negado",
	"negativa formal",
	"protocolo",
	"documentos",
	"comprovante",
	"nota fiscal",
	"recibo",
	"orcamento",
	"risco clinico",
	"urgente",
	"liminar",
}

var publicAssistanceRiskSignals = []string{
	"bpc loas",
	"cadunico",
	"beneficio assistencial",
	"renda familiar",
	"vulnerabilidade",
	"baixa renda",
	"miserabilidade",
	"defensoria publica",
	"justica gratuita",
}

var adminSelfServiceRiskSignals = []string{
	"cumprimento de exigencia",
	"anexar documentos",
	"reenviar arquivos",
	"autoatendimento",
}

var freeServiceSignals = []string{
	"advogado gratuito",
	"advogada gratuita",
	"consulta gratis",
	"consulta gratuita",
	"servico gratuito",
	"atendimento gratuito",
	"juridico gratuito",
	"defensoria publica",
	"justica gratuita",
	"sem pagar",
	"nao posso pagar",
	"de graca",
	"pro bono",
	"gratuidade",
}

var researchOnlySignals = []string{
	"apenas curiosidade",
	"so pesquisando",
	"somente pesquisando",
	"para estudar",
	"trabalho de faculdade",
	"tcc",
	"resumo para prova",
	"modelo pronto",
	"modelo de peticao",
	"peticao pronta",
	"pdf gratis",
	"jurisprudencia para estudo",
	"significado juridico",
}

func Validate(root string) Report {
	gates, loadReport := LoadRecords(root)
	issues := append([]Issue{}, loadReport.Issues...)
	if !loadReport.Passed() {
		return Report{Issues: issues}
	}
	drafts, draftReport := batchfinaldrafts.LoadRecords(root)
	issues = append(issues, convertDraftIssues(draftReport)...)
	if len(gates) == 0 && loadReport.Passed() {
		issues = append(issues, Issue{Code: "paid_intent_gates_empty", Message: "data/editorial/batch_paid_intent_gates.jsonl"})
	}
	if len(gates) != len(drafts) {
		issues = append(issues, Issue{Code: "paid_intent_gate_count_mismatch", Message: fmt.Sprintf("gates=%d drafts=%d", len(gates), len(drafts))})
	}
	draftByID := make(map[string]batchfinaldrafts.Record)
	for _, entry := range drafts {
		draftByID[entry.Record.DraftID] = entry.Record
	}
	seen := make(map[string]int)
	for _, entry := range gates {
		report := ValidateRecord(entry.Record, draftByID)
		for _, issue := range report.Issues {
			issue.Message = fmt.Sprintf("line=%d %s", entry.Line, issue.Message)
			issues = append(issues, issue)
		}
		if previous := seen[entry.Record.UniqueIntentID]; previous > 0 {
			issues = append(issues, Issue{Code: "paid_intent_duplicate_intent", Message: fmt.Sprintf("line=%d previous_line=%d id=%s", entry.Line, previous, entry.Record.UniqueIntentID)})
		}
		seen[entry.Record.UniqueIntentID] = entry.Line
	}
	return Report{Issues: issues}
}

func LoadRecords(root string) ([]Entry, Report) {
	projectRoot, err := findProjectRoot(root)
	if err != nil {
		return nil, Report{Issues: []Issue{{Code: "project_root_not_found", Message: err.Error()}}}
	}
	path := filepath.Join(projectRoot, "data", "editorial", "batch_paid_intent_gates.jsonl")
	file, err := os.Open(path)
	if err != nil {
		return nil, Report{Issues: []Issue{{Code: "paid_intent_gates_missing", Message: err.Error()}}}
	}
	defer file.Close()

	entries := make([]Entry, 0)
	issues := make([]Issue, 0)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 4096), 65536)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var record Record
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			issues = append(issues, Issue{Code: "paid_intent_gate_invalid_json", Message: fmt.Sprintf("line=%d %s", lineNumber, err.Error())})
			continue
		}
		entries = append(entries, Entry{Line: lineNumber, Record: record})
	}
	if err := scanner.Err(); err != nil {
		issues = append(issues, Issue{Code: "paid_intent_gate_scan_failed", Message: err.Error()})
	}
	return entries, Report{Issues: issues}
}

func ValidateRecord(record Record, draftsByID map[string]batchfinaldrafts.Record) Report {
	issues := make([]Issue, 0)
	if record.PaidIntentGateID == "" || record.DraftID == "" || record.UniqueIntentID == "" || record.BatchID == "" || record.SourceMatrixID == "" {
		issues = append(issues, Issue{Code: "paid_intent_gate_missing_identity", Message: record.PaidIntentGateID})
	}
	draft, ok := draftsByID[record.DraftID]
	if !ok {
		issues = append(issues, Issue{Code: "paid_intent_gate_missing_draft", Message: record.DraftID})
	} else {
		if record.UniqueIntentID != draft.UniqueIntentID || record.BatchID != draft.BatchID || record.SourceMatrixID != draft.SourceMatrixID || record.Term != draft.Term || record.CandidatePath != draft.CandidatePath {
			issues = append(issues, Issue{Code: "paid_intent_gate_draft_mismatch", Message: record.UniqueIntentID})
		}
		expected := EvaluateDraft(draft)
		issues = append(issues, compareGateToEvaluation(record, expected)...)
	}
	if !validStatus(record.PaidIntentStatus) {
		issues = append(issues, Issue{Code: "paid_intent_gate_invalid_status", Message: record.PaidIntentStatus})
	}
	if record.IndexPolicy != "noindex" {
		issues = append(issues, Issue{Code: "paid_intent_gate_invalid_index_policy", Message: record.IndexPolicy})
	}
	if record.RenderAllowed {
		issues = append(issues, Issue{Code: "paid_intent_gate_render_allowed", Message: record.UniqueIntentID})
	}
	if record.SitemapAllowed {
		issues = append(issues, Issue{Code: "paid_intent_gate_sitemap_allowed", Message: record.UniqueIntentID})
	}
	if record.PublicationAllowed {
		issues = append(issues, Issue{Code: "paid_intent_gate_publication_allowed", Message: record.UniqueIntentID})
	}
	if record.PublicPath != "" {
		issues = append(issues, Issue{Code: "paid_intent_gate_has_public_path", Message: record.PublicPath})
	}
	if record.CheckedAt == "" {
		issues = append(issues, Issue{Code: "paid_intent_gate_without_checked_at", Message: record.UniqueIntentID})
	}
	return Report{Issues: issues}
}

func EvaluateDraft(record batchfinaldrafts.Record) Record {
	text := strings.Join([]string{
		record.Term,
		record.CandidateTitle,
		record.CandidateMetaDescription,
		record.Opening,
		record.DocumentGuidance,
		record.DigitalTriage,
		record.CTAContextMessage,
	}, " ")
	score := scoreText(text)
	status := PassedBlockedStatus
	routing := "paid_intent_candidate_blocked_publication"
	risks := append([]string{}, score.freeSignals...)
	risks = append(risks, score.researchSignals...)
	risks = append(risks, score.publicAssistanceSignals...)
	risks = append(risks, score.adminSelfServiceSignals...)
	if len(score.publicAssistanceSignals) > 0 {
		status = PublicAssistanceBlockedStatus
		routing = "commercial_publication_blocked_public_assistance"
	} else if len(score.adminSelfServiceSignals) > 0 {
		status = AdminSelfServiceBlockedStatus
		routing = "commercial_publication_blocked_admin_self_service"
	}
	return Record{
		PaidIntentGateID:   "paid-intent-" + record.UniqueIntentID,
		DraftID:            record.DraftID,
		UniqueIntentID:     record.UniqueIntentID,
		BatchID:            record.BatchID,
		SourceMatrixID:     record.SourceMatrixID,
		Term:               record.Term,
		CandidatePath:      record.CandidatePath,
		PaidIntentStatus:   status,
		PaidBusinessScore:  score.paidBusinessScore,
		PaidSignals:        score.paidSignals,
		BusinessSignals:    score.businessSignals,
		RiskSignals:        uniqueStrings(risks),
		RoutingDecision:    routing,
		IndexPolicy:        "noindex",
		RenderAllowed:      false,
		SitemapAllowed:     false,
		PublicationAllowed: false,
		PublicPath:         "",
		CheckedAt:          record.CheckedAt,
	}
}

func ValidateText(text string, mode Mode) Report {
	score := scoreText(text)
	issues := make([]Issue, 0)
	if len(score.freeSignals) > 0 {
		issues = append(issues, Issue{
			Code:    "paid_intent_free_service_signal",
			Message: "texto induz gratuidade/nao-pagamento: " + strings.Join(score.freeSignals, ","),
		})
	}
	if len(score.researchSignals) > 0 {
		issues = append(issues, Issue{
			Code:    "paid_intent_research_only_signal",
			Message: "texto parece pesquisa academica/curiosidade, nao contratacao: " + strings.Join(score.researchSignals, ","),
		})
	}
	if len(score.publicAssistanceSignals) > 0 {
		issues = append(issues, Issue{
			Code:    "paid_intent_public_assistance_free_risk",
			Message: "texto indica assistencia publica/beneficio de baixa renda: " + strings.Join(score.publicAssistanceSignals, ","),
		})
	}
	if len(score.adminSelfServiceSignals) > 0 {
		issues = append(issues, Issue{
			Code:    "paid_intent_admin_self_service_risk",
			Message: "texto indica fluxo administrativo de autoatendimento: " + strings.Join(score.adminSelfServiceSignals, ","),
		})
	}
	if len(issues) > 0 {
		return Report{Issues: issues}
	}
	if mode == RequirePaidSignal && len(score.paidSignals) == 0 {
		issues = append(issues, Issue{
			Code:    "paid_intent_missing_paid_signal",
			Message: "faltou sinal explicito de contratacao paga, honorarios, consulta paga ou atendimento particular",
		})
	}
	if mode == RequirePaidSignal && score.paidBusinessScore < minimumPaidBusinessScore {
		issues = append(issues, Issue{
			Code:    "paid_intent_low_business_score",
			Message: fmt.Sprintf("score=%d minimo=%d paid=%s business=%s", score.paidBusinessScore, minimumPaidBusinessScore, strings.Join(score.paidSignals, ","), strings.Join(score.businessSignals, ",")),
		})
	}
	return Report{Issues: issues}
}

func scoreText(text string) textScore {
	normalized := normalize(text)
	paid := matchedSignals(normalized, paidServiceSignals)
	business := matchedSignals(normalized, businessReadinessSignals)
	free := matchedSignals(normalized, freeServiceSignals)
	research := matchedSignals(normalized, researchOnlySignals)
	publicAssistance := matchedSignals(normalized, publicAssistanceRiskSignals)
	adminSelfService := matchedSignals(normalized, adminSelfServiceRiskSignals)
	score := len(paid)*3 + min(len(business), 3)
	return textScore{
		paidSignals:             paid,
		businessSignals:         business,
		freeSignals:             free,
		researchSignals:         research,
		publicAssistanceSignals: publicAssistance,
		adminSelfServiceSignals: adminSelfService,
		paidBusinessScore:       score,
	}
}

func compareGateToEvaluation(record Record, expected Record) []Issue {
	issues := make([]Issue, 0)
	if record.PaidIntentGateID != expected.PaidIntentGateID {
		issues = append(issues, Issue{Code: "paid_intent_gate_id_mismatch", Message: record.UniqueIntentID})
	}
	if record.PaidIntentStatus != expected.PaidIntentStatus {
		issues = append(issues, Issue{Code: "paid_intent_gate_status_mismatch", Message: fmt.Sprintf("%s record=%s expected=%s", record.UniqueIntentID, record.PaidIntentStatus, expected.PaidIntentStatus)})
	}
	if record.RoutingDecision != expected.RoutingDecision {
		issues = append(issues, Issue{Code: "paid_intent_gate_routing_mismatch", Message: record.UniqueIntentID})
	}
	if record.PaidBusinessScore != expected.PaidBusinessScore {
		issues = append(issues, Issue{Code: "paid_intent_gate_score_mismatch", Message: fmt.Sprintf("%s record=%d expected=%d", record.UniqueIntentID, record.PaidBusinessScore, expected.PaidBusinessScore)})
	}
	if !sameStringSet(record.PaidSignals, expected.PaidSignals) {
		issues = append(issues, Issue{Code: "paid_intent_gate_paid_signals_mismatch", Message: fmt.Sprintf("%s record=%s expected=%s", record.UniqueIntentID, strings.Join(record.PaidSignals, ","), strings.Join(expected.PaidSignals, ","))})
	}
	if !sameStringSet(record.RiskSignals, expected.RiskSignals) {
		issues = append(issues, Issue{Code: "paid_intent_gate_risk_signals_mismatch", Message: fmt.Sprintf("%s record=%s expected=%s", record.UniqueIntentID, strings.Join(record.RiskSignals, ","), strings.Join(expected.RiskSignals, ","))})
	}
	return issues
}

func validStatus(value string) bool {
	return value == PassedBlockedStatus || value == PublicAssistanceBlockedStatus || value == AdminSelfServiceBlockedStatus
}

func matchedSignals(text string, signals []string) []string {
	matches := make([]string, 0)
	for _, signal := range signals {
		if strings.Contains(text, signal) {
			matches = append(matches, signal)
		}
	}
	sort.Strings(matches)
	return matches
}

func normalize(text string) string {
	mapped := strings.Map(func(r rune) rune {
		switch r {
		case 'á', 'à', 'â', 'ã', 'ä', 'Á', 'À', 'Â', 'Ã', 'Ä':
			return 'a'
		case 'é', 'è', 'ê', 'ë', 'É', 'È', 'Ê', 'Ë':
			return 'e'
		case 'í', 'ì', 'î', 'ï', 'Í', 'Ì', 'Î', 'Ï':
			return 'i'
		case 'ó', 'ò', 'ô', 'õ', 'ö', 'Ó', 'Ò', 'Ô', 'Õ', 'Ö':
			return 'o'
		case 'ú', 'ù', 'û', 'ü', 'Ú', 'Ù', 'Û', 'Ü':
			return 'u'
		case 'ç', 'Ç':
			return 'c'
		default:
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				return unicode.ToLower(r)
			}
			return ' '
		}
	}, text)
	return strings.Join(strings.Fields(mapped), " ")
}

func min(left int, right int) int {
	if left < right {
		return left
	}
	return right
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func sameStringSet(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	leftCopy := append([]string{}, left...)
	rightCopy := append([]string{}, right...)
	sort.Strings(leftCopy)
	sort.Strings(rightCopy)
	for index := range leftCopy {
		if leftCopy[index] != rightCopy[index] {
			return false
		}
	}
	return true
}

func convertDraftIssues(report batchfinaldrafts.Report) []Issue {
	issues := make([]Issue, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, Issue{Code: "paid_intent_draft_" + issue.Code, Message: issue.Message})
	}
	return issues
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
