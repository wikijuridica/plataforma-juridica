package contract_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"portaljuridico/internal/cta"
	"portaljuridico/internal/scale"
)

func TestAgentsFirstLineForGoalModeDoesNotAllowStopping(t *testing.T) {
	text := readRootFile(t, "AGENTS.md")
	firstLine := strings.SplitN(text, "\n", 2)[0]
	want := "SE o Codex estiver em modo /goal, não é para parar enquanto tem trabalho."
	if !strings.Contains(firstLine, want) {
		t.Fatalf("AGENTS.md first line = %q, want it to contain %q", firstLine, want)
	}
}

func TestContractsRequireContinuationBeyondCheckpoints(t *testing.T) {
	for _, path := range []string{"AGENTS.md", "GOAL.md"} {
		text := readRootFile(t, path)
		requireContains(t, text, "Checkpoint não é ordem de parada")
		requireContains(t, text, "continuar o próximo ciclo")
		requireContains(t, text, "só pode encerrar quando o escopo completo estiver comprovado")
		requireContains(t, text, "mínimo 10 mil páginas")
		requireContains(t, text, "WhatsApp")
		requireContains(t, text, "Antes de publicar conteúdo jurídico em escala")
		requireContains(t, text, "laboratório")
		requireContains(t, text, "validar, refinar, testar novamente")
		requireContains(t, text, "não criar 10 mil páginas como spam")
		requireContains(t, text, "escrita natural")
		requireContains(t, text, "não considerar fontes oficiais como alvo de scraping")
		requireContains(t, text, "não criar clone")
		requireContains(t, text, "conteúdo próprio, natural e único")
		requireContains(t, text, "Codex deve ser autônomo")
		requireContains(t, text, "engenheiro sênior, arquiteto e criador de conteúdo jurídico")
		requireContains(t, text, "deve continuar até terminar")
		requireContains(t, text, "corrigir sem pedir aprovação")
		requireContains(t, text, "algoritmos inteligentes, auditáveis e explicáveis")
		requireContains(t, text, "Se um algoritmo estiver burro")
		requireContains(t, text, "o código deve explicar suas próprias decisões")
		requireContains(t, text, "Conteúdo visível ao público deve ser escrito em PT-BR")
		requireContains(t, text, "grafia correta")
		requireContains(t, text, "pesquisar a Central da Pesquisa Google no dia da sessão")
		requireContains(t, text, "não há limite fixo oficial de caracteres")
		requireContains(t, text, "orçamento conservador do projeto")
		requireContains(t, text, "`title`: 20 a 65 caracteres Unicode")
		requireContains(t, text, "metadescrição: 70 a 160 caracteres Unicode")
		requireContains(t, text, "CPU deve ser reservado para tráfego legítimo")
		requireContains(t, text, "Nada deve ser deixado para o futuro")
		requireContains(t, text, "Tudo que estiver no escopo é para esta sessão")
		requireContains(t, text, "alternativa segura que resolva o requisito")
		requireContains(t, text, "É proibido mascarar pendência como entrega")
		requireContains(t, text, "não usar falta de acesso, incerteza ou pendência como desculpa para parar")
		requireContains(t, text, "É para usar rede quando a rede for necessária")
		requireContains(t, text, "Se a rede do sandbox falhar")
		requireContains(t, text, "sandbox_permissions")
		requireContains(t, text, "require_escalated")
		requireContains(t, text, "sem perguntar no chat")
		requireContains(t, text, "banco de dados leve e organizado")
		requireContains(t, text, "term_seeds")
		requireContains(t, text, "separar ingestão de termos")
		requireContains(t, text, "não misturar fonte bruta, auditoria de fonte, rascunho editorial e conteúdo publicado")
		requireContains(t, text, "draft_only")
	}
}

func TestCheckpointCarriesNextExecutionPlanInsteadOfStopping(t *testing.T) {
	text := readRootFile(t, "CHECKPOINT.md")

	requireContains(t, text, "nao e ordem de parada")
	requireContains(t, text, "plano de continuidade")
	requireContains(t, text, "continuar P0")
	requireContains(t, text, "nao publicar 10 mil paginas")
	requireContains(t, text, "laboratorio")
	requireContains(t, text, "validar novamente")
}

func TestP0ScalePlanTargetsAtLeastTenThousandPagesWithoutPublishingThem(t *testing.T) {
	plan, err := scale.LoadPlan(".")
	if err != nil {
		t.Fatal(err)
	}

	if plan.TotalPlannedPages() < 10000 {
		t.Fatalf("TotalPlannedPages = %d, want at least 10000", plan.TotalPlannedPages())
	}
	if plan.PublishedPagesDuringP0() != 0 {
		t.Fatalf("PublishedPagesDuringP0 = %d, want 0 while P0 architecture is incomplete", plan.PublishedPagesDuringP0())
	}
	if !plan.RequiresArchitectureBeforeContent {
		t.Fatal("scale plan must require architecture before content")
	}
	if !plan.RequiresSourceProvenance {
		t.Fatal("scale plan must require source provenance before legal content")
	}
	if !plan.RequiresEditorialReview {
		t.Fatal("scale plan must require editorial review before legal content")
	}
}

func TestWhatsAppCTAContractExistsButDoesNotBypassLegalQuality(t *testing.T) {
	policy, err := cta.LoadPolicy(".")
	if err != nil {
		t.Fatal(err)
	}

	if !policy.Enabled {
		t.Fatal("CTA policy should be enabled as an architecture contract")
	}
	if policy.Channel != "whatsapp" {
		t.Fatalf("Channel = %q, want whatsapp", policy.Channel)
	}
	if !policy.RequiresApprovedLegalContent || !policy.RequiresSourceProvenance || !policy.RequiresEditorialReview {
		t.Fatalf("CTA policy cannot bypass legal quality: %+v", policy)
	}
	if strings.TrimSpace(policy.PhonePlaceholder) == "" {
		t.Fatal("CTA policy must carry a placeholder for owned WhatsApp configuration")
	}
}

func readRootFile(t *testing.T, path string) string {
	t.Helper()
	root := findRoot(t)
	data, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func findRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("project root not found")
		}
		dir = parent
	}
}
