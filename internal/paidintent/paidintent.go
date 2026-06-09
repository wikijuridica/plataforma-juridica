package paidintent

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"portaljuridico/internal/batchfinaldrafts"
)

type Mode string

const RequirePaidSignal Mode = "require_paid_signal"

const minimumPaidBusinessScore = 4

type Issue struct {
	Code    string
	Message string
}

type Report struct {
	Issues []Issue
}

type textScore struct {
	paidSignals       []string
	businessSignals   []string
	freeSignals       []string
	researchSignals   []string
	paidBusinessScore int
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
	entries, loadReport := batchfinaldrafts.LoadRecords(root)
	issues := convertDraftIssues(loadReport)
	if len(entries) == 0 && loadReport.Passed() {
		issues = append(issues, Issue{Code: "paid_intent_no_final_drafts", Message: "data/editorial/batch_final_authorial_drafts.jsonl"})
	}
	for _, entry := range entries {
		record := entry.Record
		text := strings.Join([]string{
			record.Term,
			record.CandidateTitle,
			record.CandidateMetaDescription,
			record.DigitalTriage,
			record.CTAContextMessage,
		}, " ")
		report := ValidateText(text, RequirePaidSignal)
		for _, issue := range report.Issues {
			issue.Message = fmt.Sprintf("line=%d intent=%s %s", entry.Line, record.UniqueIntentID, issue.Message)
			issues = append(issues, issue)
		}
	}
	return Report{Issues: issues}
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
	score := len(paid)*3 + min(len(business), 3)
	return textScore{
		paidSignals:       paid,
		businessSignals:   business,
		freeSignals:       free,
		researchSignals:   research,
		paidBusinessScore: score,
	}
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

func convertDraftIssues(report batchfinaldrafts.Report) []Issue {
	issues := make([]Issue, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, Issue{Code: "paid_intent_draft_" + issue.Code, Message: issue.Message})
	}
	return issues
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
