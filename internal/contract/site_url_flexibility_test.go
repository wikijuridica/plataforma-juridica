package contract_test

import (
	"strings"
	"testing"

	"portaljuridico/internal/content"
	"portaljuridico/internal/seo"
)

func TestProjectBaseURLUsesOfficialWikiJuridicaDomain(t *testing.T) {
	repo, err := content.LoadRepository(".")
	if err != nil {
		t.Fatalf("could not load repository: %v", err)
	}
	if !seo.IsAbsoluteHTTPSURL(repo.BaseURL) {
		t.Fatalf("base_url must be absolute HTTPS, got %q", repo.BaseURL)
	}
	if repo.BaseURL != "https://wikijuridica.com.br" {
		t.Fatalf("base_url=%q, want official project domain https://wikijuridica.com.br", repo.BaseURL)
	}
	if repo.BaseURLMode != "official_configured" {
		t.Fatalf("base_url_mode=%q, want official_configured", repo.BaseURLMode)
	}
	if repo.OfficialURLStatus != "locked" {
		t.Fatalf("official_url_status=%q, want locked", repo.OfficialURLStatus)
	}
	if !repo.OfficialURLLocked {
		t.Fatal("official project URL must be locked after wikijuridica.com.br was defined")
	}
	if strings.Contains(repo.BaseURL, ".example") && repo.BaseURLMode != "lab_placeholder" {
		t.Fatalf("example domain can only be used as lab placeholder, got mode=%q", repo.BaseURLMode)
	}
}
