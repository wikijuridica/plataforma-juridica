package quality

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"portaljuridico/internal/content"
	"portaljuridico/internal/editorial"
	"portaljuridico/internal/router"
	"portaljuridico/internal/seo"
	"portaljuridico/internal/sources"
)

const minIndexableWords = 90
const minInternalLinks = 2
const nearDuplicateThreshold = 0.82

var nonText = regexp.MustCompile(`[^a-z0-9 ]+`)

type Issue struct {
	PagePath string
	Code     string
	Message  string
}

type Report struct {
	Issues []Issue
}

func (r Report) Passed() bool {
	return len(r.Issues) == 0
}

func (r Report) Codes() []string {
	codes := make([]string, 0, len(r.Issues))
	for _, issue := range r.Issues {
		codes = append(codes, issue.Code)
	}
	return codes
}

func (r Report) HasIssue(code string) bool {
	for _, issue := range r.Issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

func (r Report) Messages() []string {
	messages := make([]string, 0, len(r.Issues))
	for _, issue := range r.Issues {
		messages = append(messages, fmt.Sprintf("%s: %s: %s", issue.PagePath, issue.Code, issue.Message))
	}
	return messages
}

func ValidatePages(pages []content.Page) Report {
	issues := make([]Issue, 0)
	indexable := make([]content.Page, 0)

	for _, page := range pages {
		if !editorial.ValidStatuses[page.Status] {
			addIssue(&issues, page, "invalid_status", "status editorial desconhecido")
		}
		if !editorial.ValidIndexPolicies[page.IndexPolicy] {
			addIssue(&issues, page, "invalid_index_policy", "politica index/noindex invalida")
		}
		if !router.IsCleanPublicPath(page.Path) {
			addIssue(&issues, page, "unclean_url", "URL deve ser limpa, minuscula e sem parametros")
		}
		if !seo.IsAbsoluteHTTPSURL(page.CanonicalURL) {
			addIssue(&issues, page, "invalid_canonical", "canonical deve ser HTTPS absoluto")
		}
		if router.CanonicalPath(page.CanonicalURL) != page.Path {
			addIssue(&issues, page, "canonical_path_mismatch", "canonical deve apontar para a propria rota")
		}
		if editorial.IsIndexable(page) {
			indexable = append(indexable, page)
		}
	}

	for _, page := range indexable {
		if page.UniqueIntentID == "" {
			addIssue(&issues, page, "missing_unique_intent", "pagina indexavel sem unique_intent_id")
		}
		if page.Title == "" {
			addIssue(&issues, page, "missing_title", "pagina indexavel sem titulo")
		}
		if page.MetaDescription == "" {
			addIssue(&issues, page, "missing_meta_description", "pagina indexavel sem meta description")
		}
		if wordCount(page.PlainText()) < minIndexableWords {
			addIssue(&issues, page, "thin_content", "conteudo indexavel abaixo do minimo textual")
		}
		if len(page.InternalLinks) < minInternalLinks {
			addIssue(&issues, page, "missing_useful_internal_links", "pagina indexavel exige links internos uteis")
		}
		if page.PublicationPurpose == "" {
			addIssue(&issues, page, "missing_publication_purpose", "pagina indexavel exige motivo de publicacao")
		}
		if page.IsLegalContent() && len(page.SourceProvenance) == 0 {
			addIssue(&issues, page, "legal_content_without_source", "conteudo juridico indexavel exige fonte")
		}
		if page.IsLegalContent() && (page.ReviewedAt == "" || page.Reviewer == "") {
			addIssue(&issues, page, "legal_content_without_review", "conteudo juridico indexavel exige revisao")
		}
		if page.IsLegalContent() && page.LegalNotice == "" {
			addIssue(&issues, page, "legal_content_without_notice", "conteudo juridico exige aviso informativo")
		}
	}

	checkDuplicates(indexable, &issues)
	return Report{Issues: issues}
}

func ValidatePagesWithSources(pages []content.Page, registry sources.Registry) Report {
	report := ValidatePages(pages)
	issues := append([]Issue{}, report.Issues...)
	for _, page := range pages {
		if !editorial.IsIndexable(page) || !page.IsLegalContent() {
			continue
		}
		for _, provenance := range page.SourceProvenance {
			source, ok := registry.ByID(provenance.SourceID)
			if !ok {
				addIssue(&issues, page, "source_not_registered", "fonte juridica nao registrada")
				continue
			}
			if source.AuditStatus != "approved" || source.IngestionEnabled {
				addIssue(&issues, page, "source_not_approved_for_indexable_legal_content", "fonte ainda nao aprovada para conteudo juridico indexavel")
			}
		}
	}
	return Report{Issues: issues}
}

func addIssue(issues *[]Issue, page content.Page, code string, message string) {
	*issues = append(*issues, Issue{PagePath: page.Path, Code: code, Message: message})
}

func checkDuplicates(indexable []content.Page, issues *[]Issue) {
	checkField(indexable, issues, "duplicate_intent", "intencao unica duplicada", func(p content.Page) string {
		return p.UniqueIntentID
	})
	checkField(indexable, issues, "duplicate_title", "titulo duplicado", func(p content.Page) string {
		return p.Title
	})
	checkField(indexable, issues, "duplicate_meta_description", "meta description duplicada", func(p content.Page) string {
		return p.MetaDescription
	})
	checkField(indexable, issues, "duplicate_canonical", "canonical duplicado", func(p content.Page) string {
		return p.CanonicalURL
	})

	seenHashes := make(map[string]content.Page)
	for _, page := range indexable {
		digest := normalizedHash(page.PlainText())
		if previous, exists := seenHashes[digest]; exists {
			addIssue(issues, page, "duplicate_content_hash", "hash normalizado duplicado")
			addIssue(issues, previous, "duplicate_content_hash", "hash normalizado duplicado")
			continue
		}
		seenHashes[digest] = page
	}

	for i := 0; i < len(indexable); i++ {
		for j := i + 1; j < len(indexable); j++ {
			if shingleSimilarity(indexable[i].PlainText(), indexable[j].PlainText()) >= nearDuplicateThreshold {
				addIssue(issues, indexable[i], "near_duplicate_content", "conteudo parecido acima do limite")
				addIssue(issues, indexable[j], "near_duplicate_content", "conteudo parecido acima do limite")
			}
		}
	}
}

func checkField(indexable []content.Page, issues *[]Issue, code string, message string, value func(content.Page) string) {
	seen := make(map[string]content.Page)
	for _, page := range indexable {
		normalized := normalizeText(value(page))
		if previous, exists := seen[normalized]; exists {
			addIssue(issues, page, code, message)
			addIssue(issues, previous, code, message)
			continue
		}
		seen[normalized] = page
	}
}

func normalizeText(value string) string {
	lower := strings.ToLower(value)
	clean := nonText.ReplaceAllString(lower, " ")
	return strings.Join(strings.Fields(clean), " ")
}

func normalizedHash(value string) string {
	sum := sha256.Sum256([]byte(normalizeText(value)))
	return hex.EncodeToString(sum[:])
}

func wordCount(value string) int {
	normalized := normalizeText(value)
	if normalized == "" {
		return 0
	}
	return len(strings.Fields(normalized))
}

func shingles(value string, size int) map[string]bool {
	words := strings.Fields(normalizeText(value))
	result := make(map[string]bool)
	if len(words) < size {
		return result
	}
	for i := 0; i <= len(words)-size; i++ {
		result[strings.Join(words[i:i+size], " ")] = true
	}
	return result
}

func shingleSimilarity(left string, right string) float64 {
	leftSet := shingles(left, 5)
	rightSet := shingles(right, 5)
	if len(leftSet) == 0 || len(rightSet) == 0 {
		return 0
	}
	intersection := 0
	union := make(map[string]bool)
	for value := range leftSet {
		union[value] = true
		if rightSet[value] {
			intersection++
		}
	}
	for value := range rightSet {
		union[value] = true
	}
	return float64(intersection) / float64(len(union))
}
