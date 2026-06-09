package cta

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Policy struct {
	Enabled                      bool     `json:"enabled"`
	Channel                      string   `json:"channel"`
	PhonePlaceholder             string   `json:"phone_placeholder"`
	Intent                       string   `json:"intent"`
	RequiresApprovedLegalContent bool     `json:"requires_approved_legal_content"`
	RequiresSourceProvenance     bool     `json:"requires_source_provenance"`
	RequiresEditorialReview      bool     `json:"requires_editorial_review"`
	AllowedPageTypes             []string `json:"allowed_page_types"`
	BlockedStatuses              []string `json:"blocked_statuses"`
}

func LoadPolicy(root string) (Policy, error) {
	projectRoot, err := findProjectRoot(root)
	if err != nil {
		return Policy{}, err
	}
	data, err := os.ReadFile(filepath.Join(projectRoot, "content", "cta_policy.json"))
	if err != nil {
		return Policy{}, err
	}
	var policy Policy
	if err := json.Unmarshal(data, &policy); err != nil {
		return Policy{}, err
	}
	return policy, nil
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
