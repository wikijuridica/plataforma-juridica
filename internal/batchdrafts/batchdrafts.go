package batchdrafts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"portaljuridico/internal/humanscore"
	"portaljuridico/internal/scalablebatches"
)

type Record struct {
	BatchID            string   `json:"batch_id"`
	UniqueIntentID     string   `json:"unique_intent_id"`
	LegalArea          string   `json:"legal_area"`
	DraftStatus        string   `json:"draft_status"`
	Language           string   `json:"language"`
	Term               string   `json:"term"`
	ReaderProblem      string   `json:"reader_problem"`
	SourceHook         string   `json:"source_hook"`
	DocumentContext    string   `json:"document_context"`
	RiskContext        string   `json:"risk_context"`
	DigitalAction      string   `json:"digital_action"`
	CTAContext         string   `json:"cta_context"`
	SourceFamilies     []string `json:"source_families"`
	RewriteStatus      string   `json:"rewrite_status"`
	InitialIssueCodes  []string `json:"initial_issue_codes"`
	RewriteAttempts    int      `json:"rewrite_attempts"`
	HumanScore         int      `json:"human_score"`
	AILikeScore        int      `json:"ai_like_score"`
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

func Validate(root string) Report {
	entries, report := LoadRecords(root)
	if !report.Passed() {
		return report
	}
	issues := make([]Issue, 0)
	if len(entries) == 0 {
		issues = append(issues, Issue{Code: "batch_drafts_empty", Message: "data/editorial/batch_drafts.jsonl"})
	}
	validBatches := validBatchIDs(root)
	seen := make(map[string]int)
	perBatch := make(map[string]int)
	for _, entry := range entries {
		recordReport := ValidateRecord(entry.Record)
		for _, issue := range recordReport.Issues {
			issue.Message = fmt.Sprintf("line=%d %s", entry.Line, issue.Message)
			issues = append(issues, issue)
		}
		if !validBatches[entry.Record.BatchID] {
			issues = append(issues, Issue{Code: "batch_draft_unknown_batch", Message: fmt.Sprintf("line=%d batch=%s", entry.Line, entry.Record.BatchID)})
		}
		if previous := seen[entry.Record.UniqueIntentID]; previous > 0 {
			issues = append(issues, Issue{Code: "batch_draft_duplicate_intent", Message: fmt.Sprintf("line=%d previous_line=%d id=%s", entry.Line, previous, entry.Record.UniqueIntentID)})
		}
		seen[entry.Record.UniqueIntentID] = entry.Line
		perBatch[entry.Record.BatchID]++
	}
	for batchID, count := range perBatch {
		if count < 3 {
			issues = append(issues, Issue{Code: "batch_draft_too_few_samples", Message: fmt.Sprintf("%s=%d", batchID, count)})
		}
	}
	if RewrittenCount(entries) < 3 {
		issues = append(issues, Issue{Code: "batch_draft_too_few_rewrites", Message: fmt.Sprintf("rewritten=%d", RewrittenCount(entries))})
	}
	if max := MaximumPairSimilarity(entries); max > 0.64 {
		issues = append(issues, Issue{Code: "batch_draft_similarity_too_high", Message: fmt.Sprintf("max=%.2f", max)})
	}
	return Report{Issues: issues}
}

func LoadRecords(root string) ([]Entry, Report) {
	projectRoot, err := findProjectRoot(root)
	if err != nil {
		return nil, Report{Issues: []Issue{{Code: "project_root_not_found", Message: err.Error()}}}
	}
	path := filepath.Join(projectRoot, "data", "editorial", "batch_drafts.jsonl")
	file, err := os.Open(path)
	if err != nil {
		return nil, Report{Issues: []Issue{{Code: "batch_drafts_missing", Message: err.Error()}}}
	}
	defer file.Close()

	issues := make([]Issue, 0)
	entries := make([]Entry, 0)
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
			issues = append(issues, Issue{Code: "batch_draft_invalid_json", Message: fmt.Sprintf("line=%d %s", lineNumber, err.Error())})
			continue
		}
		entries = append(entries, Entry{Line: lineNumber, Record: record})
	}
	if err := scanner.Err(); err != nil {
		issues = append(issues, Issue{Code: "batch_draft_scan_failed", Message: err.Error()})
	}
	return entries, Report{Issues: issues}
}

func ValidateRecord(record Record) Report {
	issues := make([]Issue, 0)
	if record.BatchID == "" || record.UniqueIntentID == "" || record.LegalArea == "" || record.Term == "" {
		issues = append(issues, Issue{Code: "batch_draft_missing_identity", Message: record.UniqueIntentID})
	}
	if record.DraftStatus != "batch_draft_scored_blocked" {
		issues = append(issues, Issue{Code: "batch_draft_invalid_status", Message: record.DraftStatus})
	}
	if record.Language != "pt-BR" {
		issues = append(issues, Issue{Code: "batch_draft_not_ptbr", Message: record.UniqueIntentID})
	}
	if len(record.SourceFamilies) < 2 {
		issues = append(issues, Issue{Code: "batch_draft_too_few_sources", Message: record.UniqueIntentID})
	}
	text := record.FullText()
	score := humanscore.ScoreText(text)
	if record.HumanScore < 85 || score.HumanScore < 85 {
		issues = append(issues, Issue{Code: "batch_draft_human_score_too_low", Message: fmt.Sprintf("%s record=%d computed=%d", record.UniqueIntentID, record.HumanScore, score.HumanScore)})
	}
	if record.AILikeScore > 20 || score.AILikeScore > 20 {
		issues = append(issues, Issue{Code: "batch_draft_ai_score_too_high", Message: fmt.Sprintf("%s record=%d computed=%d", record.UniqueIntentID, record.AILikeScore, score.AILikeScore)})
	}
	if len(score.BlockingIssues) > 0 {
		issues = append(issues, Issue{Code: "batch_draft_text_failed_human_score", Message: record.UniqueIntentID + ":" + strings.Join(score.Codes(), ",")})
	}
	if record.RewriteStatus == "rewritten_after_score_failure" {
		if len(record.InitialIssueCodes) == 0 || record.RewriteAttempts < 1 {
			issues = append(issues, Issue{Code: "batch_draft_rewrite_evidence_missing", Message: record.UniqueIntentID})
		}
	} else if record.RewriteStatus != "clean_first_pass" {
		issues = append(issues, Issue{Code: "batch_draft_rewrite_missing", Message: record.UniqueIntentID + ":" + record.RewriteStatus})
	}
	if record.RenderAllowed {
		issues = append(issues, Issue{Code: "batch_draft_render_allowed", Message: record.UniqueIntentID})
	}
	if record.SitemapAllowed {
		issues = append(issues, Issue{Code: "batch_draft_sitemap_allowed", Message: record.UniqueIntentID})
	}
	if record.PublicationAllowed {
		issues = append(issues, Issue{Code: "batch_draft_publication_allowed", Message: record.UniqueIntentID})
	}
	if record.PublicPath != "" {
		issues = append(issues, Issue{Code: "batch_draft_has_public_path", Message: record.PublicPath})
	}
	if record.CheckedAt == "" {
		issues = append(issues, Issue{Code: "batch_draft_without_checked_at", Message: record.UniqueIntentID})
	}
	return Report{Issues: issues}
}

func (r Record) FullText() string {
	parts := []string{
		r.Term,
		r.ReaderProblem,
		r.SourceHook,
		r.DocumentContext,
		r.RiskContext,
		r.DigitalAction,
		r.CTAContext,
	}
	return strings.Join(parts, " ")
}

func RewrittenCount(entries []Entry) int {
	count := 0
	for _, entry := range entries {
		if entry.Record.RewriteStatus == "rewritten_after_score_failure" {
			count++
		}
	}
	return count
}

func MaximumPairSimilarity(entries []Entry) float64 {
	max := 0.0
	for i := 0; i < len(entries); i++ {
		left := signalSet(entries[i].Record.FullText())
		for j := i + 1; j < len(entries); j++ {
			right := signalSet(entries[j].Record.FullText())
			score := jaccard(left, right)
			if score > max {
				max = score
			}
		}
	}
	return max
}

func validBatchIDs(root string) map[string]bool {
	entries, report := scalablebatches.LoadRecords(root)
	if !report.Passed() {
		return map[string]bool{}
	}
	valid := make(map[string]bool)
	for _, entry := range entries {
		if scalablebatches.ValidateRecord(entry.Record).Passed() {
			valid[entry.Record.BatchID] = true
		}
	}
	return valid
}

func signalSet(value string) map[string]bool {
	words := strings.Fields(strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return ' '
	}, value))
	stop := map[string]bool{
		"para": true, "com": true, "que": true, "uma": true, "por": true, "dos": true, "das": true, "pelo": true, "pela": true,
		"origem": true, "fonte": true, "documentos": true, "whatsapp": true, "triagem": true, "digital": true, "juridica": true, "jurídica": true,
	}
	set := make(map[string]bool)
	for _, word := range words {
		if len(word) > 3 && !stop[word] {
			set[word] = true
		}
	}
	return set
}

func jaccard(left map[string]bool, right map[string]bool) float64 {
	if len(left) == 0 && len(right) == 0 {
		return 0
	}
	intersection := 0
	union := make(map[string]bool)
	for word := range left {
		union[word] = true
		if right[word] {
			intersection++
		}
	}
	for word := range right {
		union[word] = true
	}
	return float64(intersection) / float64(len(union))
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
