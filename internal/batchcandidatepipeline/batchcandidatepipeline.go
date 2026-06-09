package batchcandidatepipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"portaljuridico/internal/batchcandidatereviews"
	"portaljuridico/internal/batchdrafts"
	"portaljuridico/internal/batchfinaldrafts"
	"portaljuridico/internal/batchprepublication"
	"portaljuridico/internal/batchpublicmanifest"
	"portaljuridico/internal/batchsourcematrix"
	"portaljuridico/internal/batchsourcespecificity"
	"portaljuridico/internal/content"
	"portaljuridico/internal/paidintent"
	"portaljuridico/internal/seo"
)

type Result struct {
	Reviews            int
	Prepublication     int
	SourceSpecificity  int
	PublicManifest     int
	FinalDrafts        int
	SourceLocked       int
	SourceBlocked      int
	SelectedCandidates int
}

type Issue struct {
	Code    string
	Message string
}

type Report struct {
	Issues []Issue
}

type selectedCandidate struct {
	GateID         string
	BatchID        string
	UniqueIntentID string
	CandidatePath  string
	SourceMatrixID string
	Draft          batchdrafts.Record
}

func Refresh(root string) (Result, Report) {
	projectRoot, err := findProjectRoot(root)
	if err != nil {
		return Result{}, Report{Issues: []Issue{{Code: "project_root_not_found", Message: err.Error()}}}
	}

	candidateIndex, candidateReport := batchcandidatereviews.BuildCandidateIndex(projectRoot)
	matrixEntries, matrixReport := batchsourcematrix.LoadRecords(projectRoot)
	finalEntries, finalReport := batchfinaldrafts.LoadRecords(projectRoot)
	prepublicationEntries, prepublicationReport := batchprepublication.LoadRecords(projectRoot)
	paidEntries, paidReport := paidintent.LoadRecords(projectRoot)
	repo, repoErr := content.LoadRepository(projectRoot)

	issues := convertReviewIssues(candidateReport)
	issues = append(issues, convertMatrixIssues(matrixReport)...)
	issues = append(issues, convertFinalIssues(finalReport)...)
	issues = append(issues, convertPrepublicationIssues(prepublicationReport)...)
	issues = append(issues, convertPaidIssues(paidReport)...)
	if repoErr != nil {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_site_config_unavailable", Message: repoErr.Error()})
	}
	if len(issues) > 0 {
		return Result{}, Report{Issues: issues}
	}

	baseURL := "https://wikijuridica.com.br"
	if repo.BaseURL != "" {
		baseURL = repo.BaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	sourceURLsByMatrix := make(map[string][]string)
	for _, entry := range matrixEntries {
		sourceURLsByMatrix[entry.Record.MatrixID] = append([]string{}, entry.Record.SourceURLs...)
	}

	oldSEOByIntent := make(map[string]batchprepublication.Record)
	for _, entry := range prepublicationEntries {
		oldSEOByIntent[entry.Record.UniqueIntentID] = entry.Record
	}

	finalByIntent := make(map[string]batchfinaldrafts.Record)
	for _, entry := range finalEntries {
		finalByIntent[entry.Record.UniqueIntentID] = entry.Record
	}

	paidStatusByIntent := make(map[string]string)
	for _, entry := range paidEntries {
		paidStatusByIntent[entry.Record.UniqueIntentID] = entry.Record.PaidIntentStatus
	}

	selected := selectedCandidates(candidateIndex)
	eligibleFinal := make(map[string]bool)
	for _, candidate := range selected {
		finalDraft, ok := finalByIntent[candidate.UniqueIntentID]
		if !ok {
			continue
		}
		status := paidStatusByIntent[candidate.UniqueIntentID]
		if status == "" {
			status = paidintent.EvaluateDraft(finalDraft).PaidIntentStatus
		}
		if status == paidintent.PassedBlockedStatus {
			eligibleFinal[candidate.UniqueIntentID] = true
		}
	}

	reviews := make([]batchcandidatereviews.Record, 0, len(selected))
	prepublications := make([]batchprepublication.Record, 0, len(selected))
	sources := make([]batchsourcespecificity.Record, 0, len(selected))
	manifests := make([]batchpublicmanifest.Record, 0, len(selected))
	finalDrafts := make([]batchfinaldrafts.Record, 0, len(eligibleFinal))

	result := Result{SelectedCandidates: len(selected)}
	checkedAt := "2026-06-09"
	for _, candidate := range selected {
		if candidate.Draft.CheckedAt != "" {
			checkedAt = candidate.Draft.CheckedAt
		}
		review := buildReview(candidate, candidateIndex.BaseURLMode, candidateIndex.OfficialURLLocked, checkedAt)
		prepublication := buildPrepublication(candidate, review, baseURL, oldSEOByIntent[candidate.UniqueIntentID], checkedAt)
		sourceURLs := sourceURLsByMatrix[candidate.SourceMatrixID]
		if eligibleFinal[candidate.UniqueIntentID] {
			sourceURLs = finalByIntent[candidate.UniqueIntentID].SelectedSourceURLs
		}
		if len(sourceURLs) == 0 {
			issues = append(issues, Issue{Code: "batch_candidate_pipeline_missing_source_urls", Message: candidate.SourceMatrixID})
			continue
		}
		source := buildSourceSpecificity(candidate, prepublication, sourceURLs, eligibleFinal[candidate.UniqueIntentID], checkedAt)
		manifest := buildManifest(prepublication, source, checkedAt)

		reviews = append(reviews, review)
		prepublications = append(prepublications, prepublication)
		sources = append(sources, source)
		manifests = append(manifests, manifest)
		if source.SourceSpecificityStatus == batchsourcespecificity.LockedStatus {
			result.SourceLocked++
		}
		if source.SourceSpecificityStatus == batchsourcespecificity.BlockedStatus {
			result.SourceBlocked++
		}
	}
	if len(issues) > 0 {
		return result, Report{Issues: issues}
	}

	for _, candidate := range selected {
		if !eligibleFinal[candidate.UniqueIntentID] {
			continue
		}
		finalDraft := finalByIntent[candidate.UniqueIntentID]
		finalDrafts = append(finalDrafts, finalDraft)
	}
	sort.Slice(finalDrafts, func(left int, right int) bool {
		if finalDrafts[left].BatchID != finalDrafts[right].BatchID {
			return finalDrafts[left].BatchID < finalDrafts[right].BatchID
		}
		return finalDrafts[left].UniqueIntentID < finalDrafts[right].UniqueIntentID
	})

	if err := writeRecords(projectRoot, "data/editorial/batch_candidate_reviews.jsonl", reviews); err != nil {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_write_reviews_failed", Message: err.Error()})
	}
	if err := writeRecords(projectRoot, "data/editorial/batch_prepublication_gates.jsonl", prepublications); err != nil {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_write_prepublication_failed", Message: err.Error()})
	}
	if err := writeRecords(projectRoot, "data/editorial/batch_source_specificity_resolutions.jsonl", sources); err != nil {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_write_source_specificity_failed", Message: err.Error()})
	}
	if err := writeRecords(projectRoot, "data/editorial/batch_public_manifest_gates.jsonl", manifests); err != nil {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_write_manifest_failed", Message: err.Error()})
	}
	if err := writeRecords(projectRoot, "data/editorial/batch_final_authorial_drafts.jsonl", finalDrafts); err != nil {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_write_final_drafts_failed", Message: err.Error()})
	}
	if len(issues) > 0 {
		return result, Report{Issues: issues}
	}

	result.Reviews = len(reviews)
	result.Prepublication = len(prepublications)
	result.SourceSpecificity = len(sources)
	result.PublicManifest = len(manifests)
	result.FinalDrafts = len(finalDrafts)
	return result, Report{}
}

func selectedCandidates(index batchcandidatereviews.CandidateIndex) []selectedCandidate {
	selected := make([]selectedCandidate, 0, len(index.SelectedByIntent))
	for _, candidate := range index.SelectedByIntent {
		selected = append(selected, selectedCandidate{
			GateID:         candidate.GateID,
			BatchID:        candidate.BatchID,
			UniqueIntentID: candidate.UniqueIntentID,
			CandidatePath:  candidate.CandidatePath,
			SourceMatrixID: candidate.SourceMatrixID,
			Draft:          candidate.Draft,
		})
	}
	sort.Slice(selected, func(left int, right int) bool {
		if selected[left].BatchID != selected[right].BatchID {
			return selected[left].BatchID < selected[right].BatchID
		}
		return selected[left].UniqueIntentID < selected[right].UniqueIntentID
	})
	return selected
}

func buildReview(candidate selectedCandidate, baseURLMode string, officialURLLocked bool, checkedAt string) batchcandidatereviews.Record {
	term := strings.TrimSpace(candidate.Draft.Term)
	return batchcandidatereviews.Record{
		ReviewID:          "review-" + candidate.UniqueIntentID,
		GateID:            candidate.GateID,
		BatchID:           candidate.BatchID,
		UniqueIntentID:    candidate.UniqueIntentID,
		CandidatePath:     candidate.CandidatePath,
		BaseURLMode:       baseURLMode,
		OfficialURLLocked: officialURLLocked,
		ReviewStatus:      "batch_candidate_review_blocked",
		Language:          "pt-BR",
		LegalArea:         candidate.Draft.LegalArea,
		Term:              term,
		SourceMatrixID:    candidate.SourceMatrixID,
		ReviewerRole:      "juridico_editorial_lab",
		LegalReviewNotes: []string{
			"Conferir " + term + " a partir de documentos, datas e fonte oficial antes de qualquer orientação pública.",
			"A revisão deve separar problema do leitor, risco jurídico, prova digital e limites do atendimento remoto.",
			"O CTA fica contextual ao WhatsApp para triagem paga, sem promessa de resultado ou publicação automática.",
		},
		RequiredFixes: []string{
			"Resolver fonte específica final, revisão SEO e manifesto público antes de qualquer renderização.",
		},
		CTADraft:               "Para triagem jurídica online pelo WhatsApp, envie documentos principais, datas, protocolos, contratos, comprovantes e dúvida objetiva sobre " + term + " antes da análise.",
		CTAContextMessage:      "Origem: " + candidate.CandidatePath + " Gate: " + candidate.GateID + " Intent: " + candidate.UniqueIntentID + " Tema: " + term + ". WhatsApp recebe documentos e contexto para triagem jurídica digital paga, sem promessa de resultado.",
		CTAStatus:              "draft_contextual_not_public",
		PublicationBlockReason: "P0 bloqueado: revisão em massa exige fonte específica, SEO final, manifesto público finito e aprovação explícita antes de publicar.",
		RenderAllowed:          false,
		SitemapAllowed:         false,
		PublicationAllowed:     false,
		PublicPath:             "",
		CheckedAt:              checkedAt,
	}
}

func buildPrepublication(candidate selectedCandidate, review batchcandidatereviews.Record, baseURL string, old batchprepublication.Record, checkedAt string) batchprepublication.Record {
	title := candidateTitle(candidate.Draft.Term, old.CandidateTitle)
	meta := candidateMeta(candidate.Draft.Term, old.CandidateMetaDescription)
	return batchprepublication.Record{
		PrepublicationID:         "prepub-" + candidate.UniqueIntentID,
		ReviewID:                 review.ReviewID,
		GateID:                   candidate.GateID,
		BatchID:                  candidate.BatchID,
		UniqueIntentID:           candidate.UniqueIntentID,
		SourceMatrixID:           candidate.SourceMatrixID,
		Term:                     candidate.Draft.Term,
		CandidatePath:            candidate.CandidatePath,
		CandidateCanonicalURL:    strings.TrimRight(baseURL, "/") + candidate.CandidatePath,
		CandidateRobots:          "noindex,follow",
		CandidateTitle:           title,
		CandidateMetaDescription: meta,
		SourceSpecificityStatus:  "matrix_audited_final_source_pending",
		RemainingGates: []string{
			"fonte específica final do candidato",
			"revisão SEO final",
			"manifesto público finito",
			"aprovação para render e sitemap",
		},
		RenderAllowed:      false,
		SitemapAllowed:     false,
		PublicationAllowed: false,
		PublicPath:         "",
		CheckedAt:          checkedAt,
	}
}

func buildSourceSpecificity(candidate selectedCandidate, prepublication batchprepublication.Record, sourceURLs []string, locked bool, checkedAt string) batchsourcespecificity.Record {
	status := batchsourcespecificity.BlockedStatus
	blockingReason := "fonte oficial da matriz ainda não resolve o recorte específico de " + candidate.Draft.Term + " com precisão suficiente para render público em escala."
	neededSourceDetail := "auditar URL oficial específica, termo de uso e aplicabilidade ao subtema antes de liberar manifesto, HTML ou sitemap."
	if locked {
		status = batchsourcespecificity.LockedStatus
		blockingReason = ""
		neededSourceDetail = ""
	}
	return batchsourcespecificity.Record{
		ResolutionID:            "source-specificity-" + candidate.UniqueIntentID,
		PrepublicationID:        prepublication.PrepublicationID,
		ReviewID:                prepublication.ReviewID,
		BatchID:                 candidate.BatchID,
		UniqueIntentID:          candidate.UniqueIntentID,
		SourceMatrixID:          candidate.SourceMatrixID,
		Term:                    candidate.Draft.Term,
		CandidatePath:           candidate.CandidatePath,
		CandidateCanonicalURL:   prepublication.CandidateCanonicalURL,
		CandidateRobots:         prepublication.CandidateRobots,
		SourceSpecificityStatus: status,
		SelectedSourceURLs:      append([]string{}, sourceURLs...),
		BlockingReason:          blockingReason,
		NeededSourceDetail:      neededSourceDetail,
		UsePolicy:               batchsourcespecificity.UsePolicy,
		ScrapingAllowed:         false,
		IngestionAllowed:        false,
		RenderAllowed:           false,
		SitemapAllowed:          false,
		PublicationAllowed:      false,
		PublicPath:              "",
		CheckedAt:               checkedAt,
	}
}

func buildManifest(prepublication batchprepublication.Record, source batchsourcespecificity.Record, checkedAt string) batchpublicmanifest.Record {
	status := batchpublicmanifest.SourceBlockedStatus
	seoReviewRequired := false
	contentDraftRequired := false
	if source.SourceSpecificityStatus == batchsourcespecificity.LockedStatus {
		status = batchpublicmanifest.SEOReviewPendingStatus
		seoReviewRequired = true
		contentDraftRequired = true
	}
	return batchpublicmanifest.Record{
		ManifestGateID:           "public-manifest-" + source.UniqueIntentID,
		SourceResolutionID:       source.ResolutionID,
		PrepublicationID:         source.PrepublicationID,
		ReviewID:                 source.ReviewID,
		BatchID:                  source.BatchID,
		UniqueIntentID:           source.UniqueIntentID,
		SourceMatrixID:           source.SourceMatrixID,
		Term:                     source.Term,
		CandidatePath:            source.CandidatePath,
		CandidateCanonicalURL:    source.CandidateCanonicalURL,
		CandidateRobots:          source.CandidateRobots,
		CandidateTitle:           prepublication.CandidateTitle,
		CandidateMetaDescription: prepublication.CandidateMetaDescription,
		SourceSpecificityStatus:  source.SourceSpecificityStatus,
		ManifestGateStatus:       status,
		SelectedSourceURLs:       append([]string{}, source.SelectedSourceURLs...),
		BlockingReason:           source.BlockingReason,
		NeededSourceDetail:       source.NeededSourceDetail,
		IndexPolicy:              "noindex",
		UsePolicy:                batchpublicmanifest.UsePolicy,
		SEOReviewRequired:        seoReviewRequired,
		ContentDraftRequired:     contentDraftRequired,
		ManifestAllowed:          false,
		RenderAllowed:            false,
		SitemapAllowed:           false,
		PublicationAllowed:       false,
		PublicPath:               "",
		CheckedAt:                checkedAt,
	}
}

func candidateTitle(term string, previous string) string {
	previous = strings.TrimSpace(previous)
	if validTitle(previous) {
		return previous
	}
	title := uppercaseFirst(strings.TrimSpace(term))
	if title == "" {
		title = "Tema jurídico para triagem online"
	}
	if len([]rune(title)) < seo.TitleMinCharacters {
		title = strings.TrimSpace(title + " jurídico online")
	}
	if len([]rune(title)) > seo.TitleMaxCharacters {
		title = trimRunesAtWord(title, seo.TitleMaxCharacters)
	}
	return title
}

func candidateMeta(term string, previous string) string {
	previous = strings.TrimSpace(previous)
	if validMeta(previous) {
		return previous
	}
	meta := "Organize documentos, datas, protocolos e fonte oficial antes da triagem jurídica online sobre " + strings.TrimSpace(term) + "."
	if len([]rune(meta)) > seo.MetaDescriptionMaxCharacters {
		meta = "Organize documentos, datas e fonte oficial antes da triagem jurídica online: " + trimRunesAtWord(term, 72) + "."
	}
	if len([]rune(meta)) > seo.MetaDescriptionMaxCharacters {
		meta = "Organize documentos, datas e fonte oficial antes da triagem jurídica online do caso."
	}
	if len([]rune(meta)) < seo.MetaDescriptionMinCharacters {
		meta = strings.TrimRight(meta, ".") + ", com contexto suficiente para análise remota."
	}
	return meta
}

func validTitle(value string) bool {
	length := len([]rune(value))
	return length >= seo.TitleMinCharacters && length <= seo.TitleMaxCharacters
}

func validMeta(value string) bool {
	length := len([]rune(value))
	return length >= seo.MetaDescriptionMinCharacters && length <= seo.MetaDescriptionMaxCharacters
}

func uppercaseFirst(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	runes := []rune(value)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func trimRunesAtWord(value string, limit int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= limit {
		return string(runes)
	}
	cut := limit
	for cut > seo.TitleMinCharacters && !unicode.IsSpace(runes[cut-1]) {
		cut--
	}
	if cut <= seo.TitleMinCharacters {
		cut = limit
	}
	return strings.TrimSpace(string(runes[:cut]))
}

func writeRecords[T any](projectRoot string, relativePath string, records []T) error {
	path := filepath.Join(projectRoot, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	var builder strings.Builder
	for _, record := range records {
		encoded, err := json.Marshal(record)
		if err != nil {
			return err
		}
		builder.Write(encoded)
		builder.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(builder.String()), 0644)
}

func convertReviewIssues(report batchcandidatereviews.Report) []Issue {
	issues := make([]Issue, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_candidate_" + issue.Code, Message: issue.Message})
	}
	return issues
}

func convertMatrixIssues(report batchsourcematrix.Report) []Issue {
	issues := make([]Issue, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_matrix_" + issue.Code, Message: issue.Message})
	}
	return issues
}

func convertFinalIssues(report batchfinaldrafts.Report) []Issue {
	issues := make([]Issue, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_final_" + issue.Code, Message: issue.Message})
	}
	return issues
}

func convertPrepublicationIssues(report batchprepublication.Report) []Issue {
	issues := make([]Issue, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_prepublication_" + issue.Code, Message: issue.Message})
	}
	return issues
}

func convertPaidIssues(report paidintent.Report) []Issue {
	issues := make([]Issue, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, Issue{Code: "batch_candidate_pipeline_paid_" + issue.Code, Message: issue.Message})
	}
	return issues
}

func (r Report) Passed() bool { return len(r.Issues) == 0 }

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
			return "", fmt.Errorf("go.mod not found from %s", start)
		}
		current = parent
	}
}
