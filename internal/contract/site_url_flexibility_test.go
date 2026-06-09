package contract_test

import (
	"strings"
	"testing"

	"portaljuridico/internal/content"
	"portaljuridico/internal/seo"
)

func TestProjectBaseURLIsConfigurableUntilOfficialURLLocks(t *testing.T) {
	repo, err := content.LoadRepository(".")
	if err != nil {
		t.Fatalf("could not load repository: %v", err)
	}
	if !seo.IsAbsoluteHTTPSURL(repo.BaseURL) {
		t.Fatalf("base_url must be absolute HTTPS, got %q", repo.BaseURL)
	}
	if repo.BaseURLMode == "" {
		t.Fatal("base_url_mode must state whether URL is lab_placeholder or official_configured")
	}
	switch repo.BaseURLMode {
	case "lab_placeholder":
		if repo.OfficialURLLocked {
			t.Fatal("lab placeholder base_url cannot be marked as official locked")
		}
		if repo.OfficialURLStatus != "not_locked" {
			t.Fatalf("lab placeholder official_url_status=%q, want not_locked", repo.OfficialURLStatus)
		}
	case "official_configured":
		if !repo.OfficialURLLocked {
			t.Fatal("official_configured base_url must be locked intentionally")
		}
	default:
		t.Fatalf("unsupported base_url_mode=%q", repo.BaseURLMode)
	}
	if strings.Contains(repo.BaseURL, ".example") && repo.BaseURLMode != "lab_placeholder" {
		t.Fatalf("example domain can only be used as lab placeholder, got mode=%q", repo.BaseURLMode)
	}
}
