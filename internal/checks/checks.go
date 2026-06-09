package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"portaljuridico/internal/architecture"
	"portaljuridico/internal/build"
	"portaljuridico/internal/content"
	"portaljuridico/internal/crawl"
	"portaljuridico/internal/editorial"
	"portaljuridico/internal/quality"
	"portaljuridico/internal/sources"
)

var Names = []string{
	"architecture",
	"content-quality",
	"seo",
	"crawlability",
	"sources",
	"sitemaps",
	"canonicals",
	"no-duplicate-content",
	"performance-budget",
}

func Run(name string, root string) []string {
	switch name {
	case "architecture":
		return checkArchitecture(root)
	case "content-quality":
		return checkContentQuality(root)
	case "seo":
		return checkSEO(root)
	case "crawlability":
		return checkCrawlability(root)
	case "sources":
		return checkSources(root)
	case "sitemaps":
		return checkSitemaps(root)
	case "canonicals":
		return checkCanonicals(root)
	case "no-duplicate-content":
		return checkNoDuplicateContent(root)
	case "performance-budget":
		return checkPerformanceBudget(root)
	default:
		return []string{"unknown_check=" + name}
	}
}

func checkArchitecture(root string) []string {
	report := architecture.Validate(root)
	errors := append([]string{}, report.Messages...)
	requiredDocs := []string{
		"docs/PROJECT_VISION.md",
		"docs/ARCHITECTURE.md",
		"docs/CONTENT_QUALITY.md",
		"docs/SEO_CRAWL_INDEXING.md",
		"docs/DATA_SOURCES.md",
		"docs/LAB_VALIDATION.md",
		"docs/ROADMAP_P0_P5.md",
		"docs/DECISIONS.md",
		"CHECKPOINT.md",
	}
	for _, doc := range requiredDocs {
		if _, err := os.Stat(filepath.Join(root, doc)); err != nil {
			errors = append(errors, "missing_required_doc="+doc)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "adr")); err != nil {
		errors = append(errors, "missing_required_doc=docs/adr/")
	}
	if _, err := os.Stat(filepath.Join(root, "tools", "lab-cycle")); err != nil {
		errors = append(errors, "missing_required_tool=tools/lab-cycle")
	}
	if externalDependencies(root) {
		errors = append(errors, "external_dependency_present_without_adr")
	}
	if pythonFiles(root) {
		errors = append(errors, "python_file_present_after_go_decision")
	}
	return errors
}

func checkContentQuality(root string) []string {
	repo, err := content.LoadRepository(root)
	if err != nil {
		return []string{err.Error()}
	}
	registry, err := sources.LoadRegistry(root)
	if err != nil {
		return []string{err.Error()}
	}
	sourceReport := registry.ValidateForP0()
	if !sourceReport.Passed() {
		return sourceReport.Messages()
	}
	return quality.ValidatePagesWithSources(repo.Pages, registry).Messages()
}

func checkSEO(root string) []string {
	out, cleanup, err := buildTemp(root)
	if err != nil {
		return []string{err.Error()}
	}
	defer cleanup()

	errors := make([]string, 0)
	filepath.Walk(out, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			errors = append(errors, err.Error())
			return nil
		}
		text := string(data)
		for _, token := range []string{"<title>", `name="description"`, `rel="canonical"`, `name="robots"`, "<main"} {
			if !strings.Contains(text, token) {
				errors = append(errors, path+":missing_"+token)
			}
		}
		if strings.Contains(strings.ToLower(text), "<script") {
			errors = append(errors, path+":script_not_allowed")
		}
		return nil
	})
	return errors
}

func checkCrawlability(root string) []string {
	repo, err := content.LoadRepository(root)
	if err != nil {
		return []string{err.Error()}
	}
	policy, err := crawl.LoadDefaultPolicy(root)
	if err != nil {
		return []string{err.Error()}
	}
	robots := crawl.RenderRobotsTXT(policy, repo.BaseURL)
	errors := make([]string, 0)
	for _, bot := range []string{"Googlebot", "OAI-SearchBot", "GPTBot"} {
		if !strings.Contains(robots, "User-agent: "+bot) {
			errors = append(errors, "missing_bot_policy="+bot)
		}
	}
	if !strings.Contains(robots, "Disallow: /buscar/") {
		errors = append(errors, "missing_search_disallow")
	}
	if !strings.Contains(robots, "Sitemap: ") {
		errors = append(errors, "missing_sitemap_directive")
	}
	return errors
}

func checkSources(root string) []string {
	registry, err := sources.LoadRegistry(root)
	if err != nil {
		return []string{err.Error()}
	}
	return registry.ValidateForP0().Messages()
}

func checkSitemaps(root string) []string {
	out, cleanup, err := buildTemp(root)
	if err != nil {
		return []string{err.Error()}
	}
	defer cleanup()

	repo, err := content.LoadRepository(root)
	if err != nil {
		return []string{err.Error()}
	}
	indexData, err := os.ReadFile(filepath.Join(out, "sitemap.xml"))
	if err != nil {
		return []string{err.Error()}
	}
	pageData, err := os.ReadFile(filepath.Join(out, "sitemaps", "pages-0001.xml"))
	if err != nil {
		return []string{err.Error()}
	}
	index := string(indexData)
	pages := string(pageData)
	errors := make([]string, 0)
	if !strings.Contains(index, "<sitemapindex") {
		errors = append(errors, "sitemap_xml_not_index")
	}
	if !strings.Contains(index, "pages-0001.xml") {
		errors = append(errors, "sitemap_index_missing_partition")
	}
	for _, page := range repo.Pages {
		inSitemap := strings.Contains(pages, page.CanonicalURL)
		if editorial.IsIndexable(page) && !inSitemap {
			errors = append(errors, page.Path+":indexable_missing_from_sitemap")
		}
		if !editorial.IsIndexable(page) && inSitemap {
			errors = append(errors, page.Path+":noindex_in_sitemap")
		}
	}
	return errors
}

func checkCanonicals(root string) []string {
	repo, err := content.LoadRepository(root)
	if err != nil {
		return []string{err.Error()}
	}
	report := quality.ValidatePages(repo.Pages)
	errors := make([]string, 0)
	for _, message := range report.Messages() {
		if strings.Contains(message, "canonical") || strings.Contains(message, "unclean_url") {
			errors = append(errors, message)
		}
	}
	return errors
}

func checkNoDuplicateContent(root string) []string {
	repo, err := content.LoadRepository(root)
	if err != nil {
		return []string{err.Error()}
	}
	report := quality.ValidatePages(repo.Pages)
	errors := make([]string, 0)
	duplicateCodes := []string{
		"duplicate_intent",
		"duplicate_title",
		"duplicate_meta_description",
		"duplicate_canonical",
		"duplicate_content_hash",
		"near_duplicate_content",
	}
	for _, message := range report.Messages() {
		for _, code := range duplicateCodes {
			if strings.Contains(message, code) {
				errors = append(errors, message)
			}
		}
	}
	return errors
}

func checkPerformanceBudget(root string) []string {
	out, cleanup, err := buildTemp(root)
	if err != nil {
		return []string{err.Error()}
	}
	defer cleanup()

	errors := make([]string, 0)
	filepath.Walk(out, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		if info.Size() > 50000 {
			errors = append(errors, fmt.Sprintf("%s:html_size_over_budget=%d", path, info.Size()))
		}
		data, err := os.ReadFile(path)
		if err != nil {
			errors = append(errors, err.Error())
			return nil
		}
		if strings.Contains(strings.ToLower(string(data)), "<script") {
			errors = append(errors, path+":script_not_allowed")
		}
		return nil
	})
	return errors
}

func buildTemp(root string) (string, func(), error) {
	out, err := os.MkdirTemp("", "portaljuridico-check-")
	if err != nil {
		return "", func() {}, err
	}
	repo, err := content.LoadRepository(root)
	if err != nil {
		os.RemoveAll(out)
		return "", func() {}, err
	}
	if _, err := build.Site(repo, out); err != nil {
		os.RemoveAll(out)
		return "", func() {}, err
	}
	return out, func() { os.RemoveAll(out) }, nil
}

func externalDependencies(root string) bool {
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "require ") || trimmed == "require (" {
			return true
		}
	}
	return false
}

func pythonFiles(root string) bool {
	found := false
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, string(filepath.Separator)+".git"+string(filepath.Separator)) {
			return nil
		}
		if strings.HasSuffix(path, ".py") || strings.HasSuffix(path, ".pyc") {
			found = true
		}
		return nil
	})
	return found
}
