package batchcandidategates

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"portaljuridico/internal/batchdraftarchive"
	"portaljuridico/internal/batchdrafts"
	"portaljuridico/internal/content"
	"portaljuridico/internal/router"
)

type Record struct {
	GateID                  string   `json:"gate_id"`
	BatchID                 string   `json:"batch_id"`
	GateStatus              string   `json:"gate_status"`
	SourceArchivePath       string   `json:"source_archive_path"`
	ArchiveMinimumRecords   int      `json:"archive_minimum_records"`
	SelectedUniqueIntentIDs []string `json:"selected_unique_intent_ids"`
	CandidatePathPrefix     string   `json:"candidate_path_prefix"`
	BaseURLMode             string   `json:"base_url_mode"`
	OfficialURLLocked       bool     `json:"official_url_locked"`
	MaxSimilarityAllowed    float64  `json:"max_similarity_allowed"`
	MaxSimilarityObserved   float64  `json:"max_similarity_observed"`
	MinimumHumanScore       int      `json:"minimum_human_score"`
	CTAContextRequired      bool     `json:"cta_context_required"`
	SourceURLAuditRequired  bool     `json:"source_url_audit_required"`
	PublicationBlockReason  string   `json:"publication_block_reason"`
	RenderAllowed           bool     `json:"render_allowed"`
	SitemapAllowed          bool     `json:"sitemap_allowed"`
	PublicationAllowed      bool     `json:"publication_allowed"`
	PublicPath              string   `json:"public_path"`
	CheckedAt               string   `json:"checked_at"`
}

type Entry struct {
	Line   int
	Record Record
}

type ArchiveIndex struct {
	ByIntent          map[string]batchdrafts.Record
	CountByBatch      map[string]int
	MaxSimilarity     float64
	BaseURLMode       string
	OfficialURLLocked bool
}

type Issue struct {
	Code    string
	Message string
}

type Report struct {
	Issues []Issue
}

var idPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func Validate(root string) Report {
	entries, report := LoadRecords(root)
	if !report.Passed() {
		return report
	}
	issues := make([]Issue, 0)
	if len(entries) != 6 {
		issues = append(issues, Issue{Code: "batch_candidate_gate_count_invalid", Message: fmt.Sprintf("gates=%d", len(entries))})
	}
	index, indexReport := BuildArchiveIndex(root)
	issues = append(issues, indexReport.Issues...)
	seen := make(map[string]int)
	for _, entry := range entries {
		recordReport := ValidateRecordAgainstArchive(entry.Record, index)
		for _, issue := range recordReport.Issues {
			issue.Message = fmt.Sprintf("line=%d %s", entry.Line, issue.Message)
			issues = append(issues, issue)
		}
		if previous := seen[entry.Record.BatchID]; previous > 0 {
			issues = append(issues, Issue{Code: "batch_candidate_duplicate_batch", Message: fmt.Sprintf("line=%d previous_line=%d batch=%s", entry.Line, previous, entry.Record.BatchID)})
		}
		seen[entry.Record.BatchID] = entry.Line
	}
	return Report{Issues: issues}
}

func BuildArchiveIndex(root string) (ArchiveIndex, Report) {
	archiveEntries, archiveReport := batchdraftarchive.LoadRecords(root)
	if !archiveReport.Passed() {
		return ArchiveIndex{}, Report{Issues: convertArchiveIssues(archiveReport)}
	}
	repo, err := content.LoadRepository(root)
	if err != nil {
		return ArchiveIndex{}, Report{Issues: []Issue{{Code: "batch_candidate_site_config_unavailable", Message: err.Error()}}}
	}
	index := ArchiveIndex{
		ByIntent:          make(map[string]batchdrafts.Record),
		CountByBatch:      make(map[string]int),
		MaxSimilarity:     batchdrafts.MaximumPairSimilarity(archiveEntries),
		BaseURLMode:       repo.BaseURLMode,
		OfficialURLLocked: repo.OfficialURLLocked,
	}
	for _, entry := range archiveEntries {
		index.ByIntent[entry.Record.UniqueIntentID] = entry.Record
		index.CountByBatch[entry.Record.BatchID]++
	}
	return index, Report{}
}

func ValidateRecordAgainstArchive(record Record, index ArchiveIndex) Report {
	issues := make([]Issue, 0)
	if record.GateID == "" || !idPattern.MatchString(record.GateID) {
		issues = append(issues, Issue{Code: "batch_candidate_invalid_gate_id", Message: record.GateID})
	}
	if record.BatchID == "" || !idPattern.MatchString(record.BatchID) {
		issues = append(issues, Issue{Code: "batch_candidate_invalid_batch_id", Message: record.BatchID})
	}
	if record.GateStatus != "blocked_batch_candidate" {
		issues = append(issues, Issue{Code: "batch_candidate_status_not_blocked", Message: record.GateStatus})
	}
	if record.SourceArchivePath != "data/editorial/batch_draft_expansion_archive.jsonl" {
		issues = append(issues, Issue{Code: "batch_candidate_wrong_archive_path", Message: record.SourceArchivePath})
	}
	if record.ArchiveMinimumRecords < 100 {
		issues = append(issues, Issue{Code: "batch_candidate_archive_minimum_too_low", Message: fmt.Sprintf("%d", record.ArchiveMinimumRecords)})
	}
	if index.CountByBatch != nil && index.CountByBatch[record.BatchID] < record.ArchiveMinimumRecords {
		issues = append(issues, Issue{Code: "batch_candidate_archive_batch_too_small", Message: fmt.Sprintf("%s=%d", record.BatchID, index.CountByBatch[record.BatchID])})
	}
	if len(record.SelectedUniqueIntentIDs) < 3 {
		issues = append(issues, Issue{Code: "batch_candidate_too_few_selected_intents", Message: record.BatchID})
	}
	seenIntents := make(map[string]bool)
	for _, intentID := range record.SelectedUniqueIntentIDs {
		if !idPattern.MatchString(intentID) {
			issues = append(issues, Issue{Code: "batch_candidate_invalid_selected_intent", Message: intentID})
		}
		if seenIntents[intentID] {
			issues = append(issues, Issue{Code: "batch_candidate_duplicate_selected_intent", Message: intentID})
		}
		seenIntents[intentID] = true
		draft, ok := index.ByIntent[intentID]
		if !ok {
			issues = append(issues, Issue{Code: "batch_candidate_selected_intent_missing", Message: intentID})
			continue
		}
		if draft.BatchID != record.BatchID {
			issues = append(issues, Issue{Code: "batch_candidate_selected_intent_wrong_batch", Message: fmt.Sprintf("%s:%s", intentID, draft.BatchID)})
		}
		if draft.HumanScore < record.MinimumHumanScore {
			issues = append(issues, Issue{Code: "batch_candidate_selected_score_too_low", Message: fmt.Sprintf("%s=%d", intentID, draft.HumanScore)})
		}
		if record.CTAContextRequired && strings.TrimSpace(draft.CTAContext) == "" {
			issues = append(issues, Issue{Code: "batch_candidate_selected_without_cta_context", Message: intentID})
		}
		if draft.SourceMatrixID == "" {
			issues = append(issues, Issue{Code: "batch_candidate_selected_without_source_matrix", Message: intentID})
		}
		if draft.RenderAllowed || draft.SitemapAllowed || draft.PublicationAllowed || draft.PublicPath != "" {
			issues = append(issues, Issue{Code: "batch_candidate_selected_public_flag", Message: intentID})
		}
		candidatePath := record.CandidatePathPrefix + intentID + "/"
		if !router.IsCleanPublicPath(candidatePath) {
			issues = append(issues, Issue{Code: "batch_candidate_generated_path_not_clean", Message: candidatePath})
		}
	}
	if !router.IsCleanPublicPath(record.CandidatePathPrefix) {
		issues = append(issues, Issue{Code: "batch_candidate_prefix_not_clean", Message: record.CandidatePathPrefix})
	}
	if index.BaseURLMode != "" && record.BaseURLMode != index.BaseURLMode {
		issues = append(issues, Issue{Code: "batch_candidate_base_url_mode_mismatch", Message: fmt.Sprintf("record=%s site=%s", record.BaseURLMode, index.BaseURLMode)})
	}
	if record.OfficialURLLocked || index.OfficialURLLocked {
		issues = append(issues, Issue{Code: "batch_candidate_url_locked_before_publication", Message: record.BatchID})
	}
	if record.MaxSimilarityAllowed > 0.64 || record.MaxSimilarityAllowed <= 0 {
		issues = append(issues, Issue{Code: "batch_candidate_similarity_limit_invalid", Message: fmt.Sprintf("%.4f", record.MaxSimilarityAllowed)})
	}
	if record.MaxSimilarityObserved > record.MaxSimilarityAllowed {
		issues = append(issues, Issue{Code: "batch_candidate_similarity_observed_too_high", Message: fmt.Sprintf("%.4f", record.MaxSimilarityObserved)})
	}
	if index.MaxSimilarity > record.MaxSimilarityAllowed+0.0001 {
		issues = append(issues, Issue{Code: "batch_candidate_archive_similarity_too_high", Message: fmt.Sprintf("%.4f", index.MaxSimilarity)})
	}
	if record.MinimumHumanScore < 88 {
		issues = append(issues, Issue{Code: "batch_candidate_minimum_score_too_low", Message: fmt.Sprintf("%d", record.MinimumHumanScore)})
	}
	if !record.CTAContextRequired {
		issues = append(issues, Issue{Code: "batch_candidate_without_cta_requirement", Message: record.BatchID})
	}
	if !record.SourceURLAuditRequired {
		issues = append(issues, Issue{Code: "batch_candidate_without_source_url_audit", Message: record.BatchID})
	}
	if strings.TrimSpace(record.PublicationBlockReason) == "" {
		issues = append(issues, Issue{Code: "batch_candidate_missing_block_reason", Message: record.BatchID})
	}
	if record.RenderAllowed {
		issues = append(issues, Issue{Code: "batch_candidate_render_allowed", Message: record.BatchID})
	}
	if record.SitemapAllowed {
		issues = append(issues, Issue{Code: "batch_candidate_sitemap_allowed", Message: record.BatchID})
	}
	if record.PublicationAllowed {
		issues = append(issues, Issue{Code: "batch_candidate_publication_allowed", Message: record.BatchID})
	}
	if record.PublicPath != "" {
		issues = append(issues, Issue{Code: "batch_candidate_has_public_path", Message: record.PublicPath})
	}
	if record.CheckedAt == "" {
		issues = append(issues, Issue{Code: "batch_candidate_without_checked_at", Message: record.BatchID})
	}
	return Report{Issues: issues}
}

func LoadRecords(root string) ([]Entry, Report) {
	projectRoot, err := findProjectRoot(root)
	if err != nil {
		return nil, Report{Issues: []Issue{{Code: "project_root_not_found", Message: err.Error()}}}
	}
	path := filepath.Join(projectRoot, "data", "editorial", "batch_candidate_gates.jsonl")
	file, err := os.Open(path)
	if err != nil {
		return nil, Report{Issues: []Issue{{Code: "batch_candidate_gates_missing", Message: err.Error()}}}
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
			issues = append(issues, Issue{Code: "batch_candidate_gate_invalid_json", Message: fmt.Sprintf("line=%d %s", lineNumber, err.Error())})
			continue
		}
		entries = append(entries, Entry{Line: lineNumber, Record: record})
	}
	if err := scanner.Err(); err != nil {
		issues = append(issues, Issue{Code: "batch_candidate_gate_scan_failed", Message: err.Error()})
	}
	return entries, Report{Issues: issues}
}

func convertArchiveIssues(report batchdraftarchive.Report) []Issue {
	issues := make([]Issue, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, Issue{Code: "batch_candidate_archive_" + issue.Code, Message: issue.Message})
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
