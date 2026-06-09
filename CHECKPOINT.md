# CHECKPOINT.md

## 2026-06-09 — Ciclo 1 — P0/P1

- ciclo: 1
- prioridade: P0/P1
- objetivo: assumir a base tecnica do portal juridico, escolher runtime proprio, materializar contrato Go/on-demand, criar arquitetura, gates e provas iniciais.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e encerramento do projeto e nao substitui definicao de pronto.
- continuidade: nao e ordem de parada; o proximo ciclo deve continuar P0 enquanto arquitetura, escala, fontes, qualidade e indexacao ainda nao estiverem maduras.
- laboratorio: validar, refinar e validar novamente; nao confiar em script isolado; combinar testes, scripts, build e inspecao de artefatos.
- entregas registradas no ciclo: contrato P0/P1 atualizado; base Go sem dependencias externas; gerador on demand proprio; build estatico operacional; HTML completo; robots.txt; sitemap index; sitemap particionado; busca interna noindex; gates de arquitetura, qualidade, SEO, crawl, sitemaps, canonicals, duplicidade e performance.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `.gitignore`, `go.mod`, `cmd/`, `internal/`, `content/`, `tools/`, `docs/`, `CHECKPOINT.md`.
- decisoes: Go-first; geracao on demand propria obrigatoria; zero dependencias externas; build estatico permitido apenas como apoio operacional; cache do Go em `/tmp` durante validacao no sandbox.
- comandos executados: `go version`; `go test ./...`; `GOCACHE=/tmp/opt-wiki-go-cache go test ./...`; `GOCACHE=/tmp/opt-wiki-go-cache go test -count=1 ./...`; `./tools/check-all`; `GOCACHE=/tmp/opt-wiki-go-cache go run ./cmd/build public`; `go list -m all`; `git diff --check`; `find . -name '*.py' -o -name '*.pyc'`.
- resultados: `go test -count=1 ./...` passou; `./tools/check-all` passou todos os checks; build gerou `generated_pages=3 indexable_pages=1 output_dir=public`; `go list -m all` retornou apenas `portaljuridico`; `git diff --check` sem output; busca por Python sem output.
- falhas: tentativa inicial em Python descartada por desalinhamento com a decisao tecnica do usuario; cache padrao do Go tentou escrever fora da area permitida; teste de sitemap verificava `?` no XML header; scanner de arquitetura acusava sua propria lista de tecnologias proibidas.
- correcoes: artefatos Python removidos; uso de `GOCACHE` em `/tmp`; contrato atualizado para on-demand proprio; teste de sitemap passou a verificar apenas `<loc>`; scanner ignora o proprio arquivo de definicao do validador.
- provas: `internal/contract/p0p1_test.go` cobre HTML, noindex, sitemap, robots, gates de qualidade, duplicidade, arquitetura proibitiva e geracao on demand com cache proprio; `public/` contem `index.html`, `buscar/index.html`, `fontes/planalto/index.html`, `robots.txt`, `sitemap.xml` e `sitemaps/pages-0001.xml`.
- plano de continuidade: continuar P0; fortalecer arquitetura de escala para minimo 10 mil paginas futuras; nao publicar 10 mil paginas durante P0; nao criar 10 mil paginas como spam; manter `published_pages_during_p0=0`; validar politica de CTA WhatsApp sem publicar conteudo juridico; pesquisar fontes corretas antes de qualquer conteudo; ampliar gates antes de qualquer geracao massiva; validar novamente apos cada refinamento.
- proximo ciclo: P0 de escala/arquitetura; depois P2 de contratos por tipo de pagina, auditoria documental das demais fontes oficiais e persistencia editorial/historico.
- riscos: Go local e 1.19.8; recursos modernos devem respeitar esse limite enquanto o toolchain nao for atualizado documentalmente.

## 2026-06-09 — Ciclo 2 — P0 continuidade, laboratorio e escala sem spam

- ciclo: 2
- prioridade: P0
- objetivo: corrigir contrato operacional para deixar explicito que checkpoint nao para o trabalho; adicionar laboratorio multi-validacao; preparar arquitetura para minimo 10 mil paginas futuras sem publicar spam; registrar CTA WhatsApp como politica subordinada a qualidade.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada e nao substitui definicao de pronto.
- continuidade: continuar P0; checkpoint deve deixar o proximo plano e o agente deve seguir quando nao houver bloqueio P0 real.
- entregas registradas no ciclo: `internal/scale`; `internal/cta`; `content/scale_plan.json`; `content/cta_policy.json`; `tools/lab-cycle`; `docs/LAB_VALIDATION.md`; contratos `AGENTS.md` e `GOAL.md` atualizados para laboratorio, continuidade, 10 mil paginas sem spam, fonte correta e escrita natural.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `content/scale_plan.json`, `content/cta_policy.json`, `internal/scale/`, `internal/cta/`, `internal/contract/continuity_test.go`, `internal/architecture/architecture.go`, `internal/checks/checks.go`, `tools/lab-cycle`, `docs/`.
- decisoes: 10 mil paginas e meta de produto, nao geracao de spam; `published_pages_during_p0=0`; CTA WhatsApp fica habilitado como politica arquitetural, mas bloqueado por fonte, revisao e aprovacao; laboratorio multi-prova e obrigatorio para P0/P1.
- comandos executados: `GOCACHE=/tmp/opt-wiki-go-cache go test -count=1 ./internal/contract`; `gofmt -w internal`; `chmod +x tools/lab-cycle`; `./tools/lab-cycle`.
- resultados: `./tools/lab-cycle` passou; `go test -count=1 ./...` passou; `./tools/check-all` passou; build gerou `generated_pages=3 indexable_pages=1 output_dir=public`; `go list -m all` retornou apenas `portaljuridico`; `git diff --check` sem output; busca por `.py`/`.pyc` sem output.
- falhas: teste inicial de continuidade nao encontrava a raiz do projeto ao ler `AGENTS.md`; plano de escala somava 9.600 paginas; contrato nao tinha frases literais de continuidade e anti-spam exigidas pelo novo teste.
- correcoes: teste passou a localizar `go.mod`; blueprints passaram a planejar 10.200 paginas; contratos receberam frases explicitas; `tools/lab-cycle` criado para validar, refinar e testar novamente.
- provas: `internal/contract/continuity_test.go` valida continuidade, laboratorio, meta 10 mil, anti-spam, CTA WhatsApp, bloqueio de publicacao em P0 e exigencia de fonte/revisao.
- plano de continuidade: continuar P0 com arquitetura de fontes oficiais; criar registro auditavel das fontes iniciais; impedir qualquer conteudo juridico sem fonte pesquisada e documentada; manter ausencia de geracao massiva ate P0/P1/P2 comprovados.
- proximo ciclo: P0 fontes oficiais e proveniencia antes de conteudo.
- riscos: ainda nao ha auditoria documental completa de todas as fontes oficiais; portanto conteudo juridico em escala segue bloqueado.

## 2026-06-09 — Ciclo 3 — P0 fontes oficiais e APIs publicas candidatas

- ciclo: 3
- prioridade: P0
- objetivo: continuar arquitetura antes de conteudo; criar registro auditavel de fontes oficiais e APIs publicas de facil acesso, mantendo ingestao bloqueada ate auditoria completa.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo ciclo deve fortalecer proveniencia, robots/termos e contratos de conteudo antes de qualquer redacao juridica.
- regra de continuidade: nao parar; sempre planejar o proximo passo, registrar e continuar executando enquanto nao houver bloqueio P0 real comprovado.
- entregas registradas no ciclo: `content/source_registry.json`; `internal/sources`; testes `internal/contract/sources_test.go`; documentos preliminares para Camara, Senado, LexML, CNJ Datajud, STF, STJ, CJF e APIs candidatas; `docs/DATA_SOURCES.md` atualizado.
- arquivos alterados: `content/source_registry.json`, `internal/sources/sources.go`, `internal/contract/sources_test.go`, `internal/architecture/architecture.go`, `docs/DATA_SOURCES.md`, `docs/data-sources/`.
- decisoes: fontes e APIs publicas podem ser candidatas se tiverem acesso simples, mas `ingestion_enabled=false` permanece obrigatorio em P0; fonte pesquisada nao equivale a fonte aprovada para conteudo.
- comandos executados: pesquisas web em fontes oficiais; `gofmt -w internal`; `GOCACHE=/tmp/opt-wiki-go-cache go test -count=1 ./internal/contract`; `./tools/lab-cycle`.
- resultados: `./tools/lab-cycle` passou; `go test -count=1 ./...` passou; `./tools/check-all` passou; build gerou `generated_pages=3 indexable_pages=1 output_dir=public`; modulo externo ausente (`go list -m all` retornou apenas `portaljuridico`).
- falhas: scanner de frameworks proibidos acusou falso positivo `astro` dentro de `sem_cadastro` no JSON de fontes.
- correcoes: scanner de arquitetura passou a procurar tecnologias proibidas por fronteira de token, evitando falso positivo em palavras portuguesas.
- provas: `sources_test.go` exige fontes iniciais registradas, documentadas e com ingestao bloqueada; `source_registry.json` inclui fontes oficiais e candidatas auxiliares sem cadastro individual quando aplicavel.
- plano de continuidade: continuar P0; adicionar auditoria de robots/termos por fonte; criar contrato de proveniencia por payload; criar validador que impeça qualquer pagina juridica se a fonte tiver `audit_status` preliminar.
- proximo ciclo: P0 proveniencia, robots/termos e bloqueio de conteudo por status da fonte.
- riscos: pesquisa preliminar nao substitui leitura completa de termos, robots, limites e privacidade; Datajud e dados processuais exigem cuidado LGPD reforcado.

## 2026-06-09 — Ciclo 4 — P0 proveniencia e commit por ciclo

- ciclo: 4
- prioridade: P0
- objetivo: continuar sem parada; tornar commit por ciclo regra contratual; adicionar proveniencia por payload; conectar registro de fontes ao pipeline real de qualidade/build.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada.
- continuidade: nao parar; sempre planejar o proximo passo, registrar e continuar executando enquanto nao houver bloqueio P0 real comprovado.
- entregas registradas no ciclo: regra de commit por ciclo em `AGENTS.md`, `GOAL.md`, `docs/LAB_VALIDATION.md` e `docs/DECISIONS.md`; `internal/provenance`; teste `internal/contract/provenance_test.go`; `tools/check-sources`; `cmd/check sources`; validação de fontes integrada a `content-quality` e `build`.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `docs/`, `content/source_registry.json`, `internal/provenance/`, `internal/sources/`, `internal/quality/`, `internal/build/`, `internal/checks/`, `internal/contract/`, `tools/`.
- decisoes: fonte pesquisada nao e fonte aprovada; todo payload oficial futuro exige fonte, URL oficial, data, hash SHA-256, snapshot de robots, snapshot de termos, campos e finalidade; ciclo deve ser commitado apos validacao.
- comandos executados: `GOCACHE=/tmp/opt-wiki-go-cache go test -count=1 ./internal/contract`; `gofmt -w internal`; `chmod +x tools/check-sources tools/lab-cycle`; `./tools/lab-cycle`.
- resultados: `./tools/lab-cycle` passou; inclui `go test -count=1 ./...`, `./tools/check-all`, `./tools/check-sources`, build, `go list -m all`, `git diff --check` e busca por `.py/.pyc`; build gerou `generated_pages=3 indexable_pages=1 output_dir=public`; `go list -m all` retornou apenas `portaljuridico`.
- falhas: patch composto inicial falhou por contexto textual em `AGENTS.md`; testes RED indicaram ausencia de `internal/provenance`, campos de auditoria de fontes e `ValidatePagesWithSources`.
- correcoes: patches aplicados em blocos menores; campos de robots/termos/proveniencia/privacidade adicionados ao registro; `ValidatePagesWithSources` bloqueia fonte preliminar em pagina juridica indexavel; `build.Site` e `check content-quality` usam registro de fontes.
- provas: `source_quality_test.go` bloqueia pagina juridica indexavel com fonte preliminar; `provenance_test.go` exige hash e snapshots; `sources_test.go` exige auditoria minima e ingestao bloqueada.
- commit: este ciclo deve ser persistido em Git apos a validacao final desta entrada.
- plano de continuidade: continuar P0; criar auditoria automatizada de robots/termos sem ativar ingestao; adicionar amostras de payload apenas como fixtures de laboratorio quando necessario; fortalecer bloqueio de CTA para paginas nao aprovadas.
- proximo ciclo: P0 auditoria robots/termos e CTA gating.
- riscos: termos e robots ainda estao marcados como `pendente_verificacao_final`; nenhuma ingestao ou publicacao juridica em escala autorizada.
- adendo do ciclo: fontes oficiais sao referencia/proveniencia, nao alvo de scraping, clone ou espelho; o site deve produzir conteudo proprio, natural e unico.
- adendo operacional: sempre revisar e validar; validacao isolada nao basta; revisar diff, artefatos, contratos e riscos antes de commit.
- adendo de produto: CTA WhatsApp e critico para paginas informativas de alta intencao de contratar advogado, mas continua bloqueado por fonte, revisao, aprovacao e qualidade.
- adendo de autonomia: Codex deve ser autonomo como engenheiro senior, arquiteto e criador de conteudo juridico; deve continuar ate terminar e corrigir bugs do escopo sem pedir aprovacao normal.

## 2026-06-09 — Ciclo 5 — P0 CTA WhatsApp gating

- ciclo: 5
- prioridade: P0
- objetivo: tratar CTA WhatsApp como componente critico de produto sem permitir que ele burle fonte, revisao, aprovacao e qualidade.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo passo deve fortalecer auditoria de robots/termos e depois preparar componente visual de CTA sem ativar em paginas nao aprovadas.
- entregas registradas no ciclo: contrato atualizado para CTA critico de alta intencao; `internal/cta.CanRender`; teste `internal/contract/cta_gate_test.go`; revisao explicita de diff e HTML gerado.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `docs/PROJECT_VISION.md`, `docs/ARCHITECTURE.md`, `docs/DECISIONS.md`, `docs/LAB_VALIDATION.md`, `CHECKPOINT.md`, `internal/cta/cta.go`, `internal/contract/cta_gate_test.go`.
- decisoes: CTA WhatsApp e critico, mas so renderiza quando a pagina for publicada, indexavel, de tipo permitido, revisada e com fonte aprovada; paginas atuais nao exibem CTA.
- comandos executados: `GOCACHE=/tmp/opt-wiki-go-cache go test -count=1 ./internal/contract`; `gofmt -w internal/cta/cta.go`; `./tools/lab-cycle`; `git diff --stat`; `rg` de contrato/CTA; inspecao de `public/index.html`.
- resultados: `./tools/lab-cycle` passou; `go test -count=1 ./...` passou; `./tools/check-all` passou; `./tools/check-sources` passou; build gerou `generated_pages=3 indexable_pages=1 output_dir=public`; HTML atual nao exibe CTA.
- falhas: teste RED indicou ausencia de `cta.CanRender`.
- correcoes: `CanRender` implementado com gates de politica, status, indexacao, tipo de pagina, revisao, proveniencia e fonte aprovada.
- provas: `cta_gate_test.go` bloqueia CTA em pagina `noindex`, bloqueia fonte preliminar e permite apenas fixture com fonte aprovada.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0; criar auditoria robots/termos por fonte e preparar renderizacao visual do CTA em modo bloqueado/testado, sem ativar nas paginas atuais.
- proximo ciclo: P0 robots/termos e componente visual de CTA gated.
- riscos: telefone real de WhatsApp ainda e placeholder; nao configurar numero nem ativar CTA publico antes dos gates.

## 2026-06-09 — Ciclo 6 — autonomia persistente do Codex

- ciclo: 6
- prioridade: P0
- objetivo: persistir regra operacional de autonomia do Codex como engenheiro senior, arquiteto e criador de conteudo juridico.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0 sem aguardar repeticao do usuario; se houver trabalho no escopo ou bug identificado, corrigir, validar, revisar, checkpointar e commitar.
- entregas registradas no ciclo: `AGENTS.md`, `GOAL.md`, `docs/PROJECT_VISION.md`, `docs/DECISIONS.md` e teste `continuity_test.go` atualizados para autonomia explicita.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `docs/PROJECT_VISION.md`, `docs/DECISIONS.md`, `CHECKPOINT.md`, `internal/contract/continuity_test.go`.
- decisoes: Codex tem autonomia para decisao normal de engenharia; so deve pedir intervencao em bloqueio P0 real; bugs/lacunas no escopo devem ser corrigidos sem aguardar permissao.
- comandos executados: `GOCACHE=/tmp/opt-wiki-go-cache go test -count=1 ./internal/contract`; `./tools/lab-cycle`.
- resultados: `./tools/lab-cycle` passou; `go test -count=1 ./...` passou; checks P0/P1 passaram; build gerou `generated_pages=3 indexable_pages=1 output_dir=public`.
- falhas: teste RED indicou ausencia da frase literal de autonomia nos contratos persistentes.
- correcoes: regra literal adicionada aos contratos e docs.
- provas: `continuity_test.go` exige `Codex deve ser autônomo`, `engenheiro sênior, arquiteto e criador de conteúdo jurídico`, `deve continuar até terminar` e `corrigir sem pedir aprovação`.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0; commitar ciclo 5/6; seguir para auditoria automatizada de robots/termos e componente visual de CTA gated.
- proximo ciclo: P0 robots/termos e renderizacao visual do CTA gated.
- riscos: manter autonomia sem violar P0; decisao normal e autonoma, mas bloqueio P0 real deve ser registrado.
- adendo temporal: antes de cada commit, executar `date`, registrar hora local e manter ciclo numerado em ordem.

## 2026-06-09 — Ciclo 7 — P0 HTML leve, CTA visual gated e algoritmos explicaveis

- data/hora local conferida antes do commit: `2026-06-09 08:26:47 -03`.
- ciclo: 7
- prioridade: P0/P1
- objetivo: continuar sem parada; tornar HTML publico leve um contrato reprovavel; renderizar CTA WhatsApp somente quando aprovado; registrar que algoritmos devem ser inteligentes, auditaveis e explicaveis; registrar PT-BR correto para texto visivel ao publico.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo ciclo deve criar contrato de baixo uso de CPU e melhorar inteligencia/diagnostico dos validadores, sem publicar conteudo em escala.
- entregas registradas no ciclo: gate de performance reforcado contra HTML pesado, runtime frontend e bundles; `render.PageWithCTA`; CSS de CTA condicionado a CTA real; contrato de algoritmos explicaveis; decisao de HTML leve para Googlebot/OAI-SearchBot/bots valiosos; contrato de PT-BR com grafia correta para conteudo publico; normalizacao Unicode em `internal/quality`.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `content/pages.json`, `docs/ARCHITECTURE.md`, `docs/CONTENT_QUALITY.md`, `docs/DECISIONS.md`, `docs/LAB_VALIDATION.md`, `docs/SEO_CRAWL_INDEXING.md`, `internal/build/build.go`, `internal/checks/checks.go`, `internal/checks/performance_budget_test.go`, `internal/contract/continuity_test.go`, `internal/contract/cta_render_test.go`, `internal/contract/p0p1_test.go`, `internal/quality/quality.go`, `internal/quality/quality_test.go`, `internal/render/render.go`.
- decisoes: leveza de HTML e requisito de indexacao; pagina publica nao deve carregar JavaScript, WebAssembly, import map, `modulepreload`, marcador de hidratacao ou payload de framework; CTA visual existe, mas so entra em pagina juridica aprovada; algoritmos burros devem ser melhorados com teste e diagnostico; conteudo visivel ao publico deve ser PT-BR com grafia correta.
- comandos executados: `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/checks`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/quality`; `./tools/check-performance-budget`; `./tools/check-seo`; `gofmt -w internal/build/build.go internal/checks/checks.go internal/checks/performance_budget_test.go internal/render/render.go internal/quality/quality.go internal/quality/quality_test.go internal/contract/continuity_test.go internal/contract/p0p1_test.go`; `./tools/lab-cycle`; `date '+%Y-%m-%d %H:%M:%S %Z'`; `git diff --stat`; `rg` em `public` para WhatsApp, scripts, runtime, bundles e grafia publica; `find public -name '*.html' -printf '%p %s\n'`.
- resultados: `./tools/lab-cycle` passou; `go test -count=1 ./...` passou; `./tools/check-all` passou; `./tools/check-sources` passou; build gerou `generated_pages=3 indexable_pages=1 output_dir=public`; `performance-budget: pass`; `seo: pass`; HTML publico gerado ficou em 3115, 1951 e 2324 bytes apos PT-BR correto.
- falhas: RED inicial apontou `htmlPerformanceIssues` inexistente; compilacao falhou por colisao de nomes entre politica de CTA e politica de crawl; teste de CTA detectou vazamento de `whatsapp` por CSS global em pagina nao aprovada; contrato de algoritmos explicaveis ainda nao existia em `AGENTS.md`/`GOAL.md`; contrato PT-BR ainda nao existia; normalizacao de qualidade quebrava letras acentuadas.
- correcoes: `htmlPerformanceIssues` implementado; nomes separados para `ctaPolicy` e `crawlPolicy`; CSS do CTA so e emitido quando o CTA real e renderizado; contratos receberam regra de algoritmos inteligentes, auditaveis e explicaveis; paginas e textos fixos publicos corrigidos para PT-BR; `normalizeText` passou a preservar letras Unicode.
- provas: `performance_budget_test.go` reprova `<script>`, CSS inline excessivo, `__NEXT_DATA__`, `data-reactroot`, `data-hydrate`, `modulepreload` e `.js`; `cta_render_test.go` prova que paginas atuais nao renderizam WhatsApp e fixture aprovada renderiza CTA; `quality_test.go` prova preservacao de acentos em PT-BR; varredura em `public/` nao encontrou WhatsApp indevido, script, runtime ou bundle.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 no ciclo 8 com contrato e gate de CPU baixo, reservando CPU para trafego legitimo; criar laboratorio temporario em `/tmp` para detectar conteudo raso/mecanico antes do Googlebot; melhorar diagnosticos dos validadores para explicar custo, causa e proximo ajuste; depois seguir para robots/termos por fonte.
- proximo ciclo: P0 CPU budget, detector de conteudo mecanico/raso, algoritmos mais inteligentes e auditoria robots/termos.
- riscos: orçamento de 50 KB por HTML e 8 KB de CSS inline e conservador e pode precisar ficar mais agressivo conforme as paginas crescerem; qualquer excecao deve ser ADR e prova, nao atalho.

## 2026-06-09 — Ciclo 8 — P0 Google Search appearance, CPU publico e conteudo mecanico

- data/hora local conferida antes do commit: `2026-06-09 08:38:34 -03`.
- ciclo: 8
- prioridade: P0/P1
- objetivo: continuar sem parada; pesquisar Google Search Central no dia da sessao; criar orcamento conservador de title/metadescricao/snippet; detectar conteudo raso/mecanico em laboratorio antes do Googlebot; reservar CPU baixo para runtime publico/producao sem confundir com laboratorio.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo ciclo deve seguir para auditoria robots/termos por fonte e refinamento de conteudo/CPU conforme falsos positivos ou falsos negativos.
- fontes pesquisadas no ciclo: Google Search Central `https://developers.google.com/search/docs/appearance/title-link?hl=pt-BR`; `https://developers.google.com/search/docs/appearance/snippet?hl=pt-br`; `https://developers.google.com/search/docs/essentials/technical`.
- entregas registradas no ciclo: `seo.ValidateSearchAppearance`; `max-snippet:160` para pagina indexavel; checks `google-search-appearance`, `mechanical-content` e `cpu-budget`; laboratorio temporario `lab-content-quality` em `/tmp`; comando `cmd/content-lab`; detector de texto raso/mecanico com diagnostico; contrato explicito de CPU baixo apenas para runtime publico/producao.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `cmd/content-lab/main.go`, `docs/CONTENT_QUALITY.md`, `docs/DECISIONS.md`, `docs/LAB_VALIDATION.md`, `docs/SEO_CRAWL_INDEXING.md`, `internal/checks/checks.go`, `internal/checks/cpu_budget_test.go`, `internal/contract/continuity_test.go`, `internal/contract/p0p1_test.go`, `internal/quality/quality.go`, `internal/quality/mechanical_content_test.go`, `internal/seo/seo.go`, `internal/seo/seo_test.go`, `tools/check-cpu-budget`, `tools/check-google-search-appearance`, `tools/check-mechanical-content`, `tools/lab-content-quality`, `tools/lab-cycle`.
- decisoes: Google nao tem limite fixo oficial de caracteres para `<title>` e metadescricao; o projeto usa orcamento conservador de 20-65 caracteres Unicode para title e 70-160 para metadescricao; conteudo mecanico deve ser bloqueado antes de exposicao ao Googlebot; comandos pesados podem existir em laboratorio/build/teste, mas runtime publico deve ser barato em CPU.
- comandos executados: pesquisas web oficiais; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/seo`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/quality`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/checks`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `./tools/check-google-search-appearance`; `./tools/check-mechanical-content`; `./tools/check-cpu-budget`; `./tools/lab-content-quality`; `./tools/lab-cycle`; `date '+%Y-%m-%d %H:%M:%S %Z'`; `rg` em `public` para metas, snippet, scripts e runtime.
- resultados: `./tools/lab-cycle` passou; `go test -count=1 ./...` passou; `check-all` passou todos os checks incluindo `google-search-appearance`, `mechanical-content`, `performance-budget` e `cpu-budget`; laboratorio temporario aprovou texto natural e reprovou texto mecanico com `thin_content`, `keyword_stuffing`, `repeated_phrase` e `mechanical_keyword_permutation`; build gerou `generated_pages=3 indexable_pages=1 output_dir=public`.
- falhas: RED inicial apontou ausencia de `TitleMaxCharacters`, `ValidateSearchAppearance`, `AppearanceIssue` e `AnalyzeText`; detector inicial gerou falso positivo ao contar repeticao normal de marca/titulo/heading da home como frase mecanica; contrato ainda nao diferenciava CPU de laboratorio versus CPU de producao.
- correcoes: SEO budget implementado com contagem Unicode; `AnalyzeText` e `AnalyzePage` implementados; analise de pagina passou a focar corpo editorial, evitando falso positivo de marca/titulo/heading; contrato esclareceu que CPU baixo vale para runtime publico/producao, enquanto laboratorio pode executar comandos mais pesados quando necessario para provar qualidade.
- provas: `seo_test.go` valida limites Unicode e `max-snippet:160`; `mechanical_content_test.go` aprova texto natural e reprova permutacao mecanica; `cpu_budget_test.go` reprova execucao externa e loop sem limite em runtime publico; `public/index.html` contem `index,follow,max-snippet:160` e nao contem script/runtime/WhatsApp indevido.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 com auditoria robots/termos por fonte; adicionar diagnosticos mais ricos para fonte aprovada versus preliminar; manter detector mecanico em observacao para falsos positivos/falsos negativos antes de qualquer conteudo em escala.
- proximo ciclo: P0 auditoria robots/termos e contratos de fonte por tipo de dado.
- riscos: limites SERP sao orcamento interno conservador, nao limite oficial Google; algoritmo mecanico e inicial e deve evoluir quando houver novos exemplos naturais ou mecanicos.

## 2026-06-09 — Ciclo 9 — P0 primeira linha /goal e continuidade sem parada

- data/hora local conferida antes do commit: `2026-06-09 08:43:12 -03`.
- ciclo: 9
- prioridade: P0
- objetivo: corrigir contrato persistente para deixar a primeira linha de `AGENTS.md` explícita: em modo `/goal`, Codex nao deve parar enquanto houver trabalho no escopo.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0 sem aguardar nova cobrança; o proximo ciclo planejado e auditoria robots/termos por fonte e contratos de fonte por tipo de dado.
- entregas registradas no ciclo: teste `TestAgentsFirstLineForGoalModeDoesNotAllowStopping`; primeira linha de `AGENTS.md` com regra literal de `/goal`; `GOAL.md` reforcado; `docs/DECISIONS.md` com decisao de nao parar em `/goal`.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `docs/DECISIONS.md`, `internal/contract/continuity_test.go`, `CHECKPOINT.md`.
- decisoes: resposta no thread, checkpoint ou build verde nao encerram modo `/goal`; o proximo passo planejado deve ficar em documentos persistentes e deve ser executado enquanto nao houver bloqueio P0 real comprovado.
- comandos executados: `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `date '+%Y-%m-%d %H:%M:%S %Z'`.
- resultados: teste de contrato passou e prova a primeira linha de `AGENTS.md`.
- falhas: RED inicial mostrou que a primeira linha de `AGENTS.md` ainda era o titulo do arquivo, nao a regra de `/goal`.
- correcoes: primeira linha substituida por regra literal; `GOAL.md` e `docs/DECISIONS.md` receberam regra equivalente.
- provas: `continuity_test.go` agora valida a primeira linha e os contratos de continuidade.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 imediatamente apos commit; proximo ciclo: auditoria robots/termos por fonte, sem scraping cego, sem ingestao e com diagnostico de fonte preliminar.
- proximo ciclo: P0 auditoria robots/termos e contratos de fonte por tipo de dado.
- riscos: a regra de `/goal` aumenta autonomia operacional, mas continua limitada por bloqueio P0 real, seguranca, LGPD, robots.txt, termos de uso e ausencia de scraping cego.

## 2026-06-09 — Ciclo 10 — P0 banco leve de termos e URL oficial do Portal da Legislacao

- data/hora local conferida antes do commit: `2026-06-09 09:15:10 -03`.
- ciclo: 10
- prioridade: P0
- objetivo: continuar sem parada; transformar ingestao de termos em banco leve proprio e separado; pesquisar em fonte publica confiavel a URL oficial do Portal da Legislacao; auditar robots/termos sem mascarar pendencias.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; termos juridicos podem iniciar rascunhos, mas conteudo publico/indexavel segue bloqueado ate fonte, revisao, qualidade, SEO, CTA e checkpoint.
- fontes pesquisadas no ciclo: pagina publica `gov.br` `https://www.gov.br/pt-br/servicos/pesquisa-de-legislacao-portal-da-legislacao`; resultados oficiais gov.br que citam `https://legislacao.presidencia.gov.br/`; auditoria HTTP local em `https://legislacao.presidencia.gov.br/robots.txt`, `https://legislacao.presidencia.gov.br/atos/?tipo=LEI&numero=10406&ano=2002`, `https://www.planalto.gov.br/robots.txt` e `https://www4.planalto.gov.br/robots.txt`.
- entregas registradas no ciclo: contrato de banco leve em `content/storage_contract.json`; camadas JSONL separadas em `data/`; modulo `internal/storage`; check `storage-contract`; auditoria live `cmd/audit-sources-live` com append em `data/source-audit/robots_terms.jsonl`; URL oficial do Planalto corrigida para `https://legislacao.presidencia.gov.br/`; contratos `AGENTS.md` e `GOAL.md` reforcados contra parar, mascarar pendencia e misturar fonte/rascunho/conteudo.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `content/source_registry.json`, `content/storage_contract.json`, `data/`, `docs/`, `internal/storage/`, `internal/sources/`, `internal/checks/`, `internal/architecture/`, `internal/contract/`, `cmd/audit-sources-live/`, `tools/audit-sources-live`, `tools/check-storage-contract`, `tools/lab-cycle`.
- decisoes: ingestao de termos e valida somente como `draft_only`; banco leve fica em JSONL proprio, stdlib-only, com `term_seeds`, `source_audits`, `source_snapshots`, `editorial_drafts` e `published_manifest`; a URL oficial de pesquisa do Portal da Legislacao e `https://legislacao.presidencia.gov.br/`, confirmada publicamente por `gov.br`; fonte auditada nao autoriza scraping nem publicacao.
- comandos executados: `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `./tools/check-storage-contract`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/storage`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./cmd/audit-sources-live`; `./tools/audit-sources-live` com rede escalada; pesquisas web oficiais; `curl -L -I --max-time 20` para URLs oficiais do Planalto/Portal da Legislacao; `./tools/lab-cycle`; `date '+%Y-%m-%d %H:%M:%S %Z'`; `git diff --stat`; `git status --short`.
- resultados: `./tools/lab-cycle` passou; `go test -count=1 ./...` passou; `check-all` passou com `storage-contract`; `check-storage-contract: pass`; build gerou `generated_pages=3 indexable_pages=1 output_dir=public`; `./tools/audit-sources-live` gravou eventos JSONL e falhou corretamente com exit 1 porque Planalto/Portal da Legislacao ainda tem `robots=network_error`; Camara, LexML e STJ ficaram `http_200` em robots/termos, mas ingestao permanece `false`.
- falhas: RED inicial mostrou ausencia de `internal/storage`; RED de append mostrou `AppendJSONL` inexistente; RED de JSONL mostrou que o validador nao lia linhas; auditoria live do STF precisou diferenciar TLS fallback; `www.planalto.gov.br` resetou conexao; `legislacao.presidencia.gov.br` e `www4.planalto.gov.br` tiveram timeout/reset em auditoria local; `curl -I` na pagina gov.br retornou 403, mas GET da auditoria registrou `http_200`.
- correcoes: `internal/storage` implementado; `AppendJSONL` implementado com limite por camada; `Validate` passou a parsear JSONL e reprovar linha invalida/record pesado; auditoria live grava eventos em `data/source-audit/robots_terms.jsonl`; Planalto foi corrigido documentalmente para `legislacao.presidencia.gov.br` sem desbloquear ingestao; contratos persistentes exigem rede escalada quando necessaria e proíbem mascarar pendencia.
- provas: `storage_test.go` prova append e rejeicao de record pesado/JSONL invalido; `storage_test.go` de contrato prova separacao de camadas; `sources_test.go` exige data de robots/termos e decisao de auditoria; `cmd/audit-sources-live/main_test.go` prova classificacao de TLS/restricao; `data/source-audit/robots_terms.jsonl` contem eventos reais de auditoria; `content/source_registry.json` mantem `ingestion_enabled=false`.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 no ciclo 11; criar camada inicial de `term_seeds` com esquema/validador de termo juridico natural, fonte e estado `draft_only`; usar fontes com auditoria documental acessivel como Camara/LexML/STJ para sementes de laboratorio, sem publicar pagina; manter Planalto bloqueado ate robots/termos ficarem auditaveis por meio seguro.
- proximo ciclo: P0 esquema de termos juridicos, validador de `term_seeds`, e laboratorio de rascunho natural PT-BR sem exposicao ao Googlebot.
- riscos: URL oficial do Portal da Legislacao foi confirmada publicamente, mas `robots.txt` nao foi alcançado pela rede local; portanto qualquer ingestao do Planalto segue proibida. `data/source-audit/robots_terms.jsonl` e append-only e pode crescer, entao o proximo ciclo deve limitar/particionar eventos se o volume aumentar.

## 2026-06-09 — Ciclo 11 — P0 term seeds como pauta draft_only

- data/hora local conferida antes do commit: `2026-06-09 09:18:08 -03`.
- ciclo: 11
- prioridade: P0
- objetivo: continuar sem parada; criar esquema e validador para `term_seeds` como pauta editorial, sem publicar pagina, sem CTA e sem indexacao.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo ciclo deve transformar seed valida em rascunho temporario de laboratorio PT-BR natural, ainda fora de `content/pages.json`.
- entregas registradas no ciclo: `internal/terms`; teste `internal/contract/term_seed_test.go`; primeira seed `responsabilidade-civil` em `data/terms/legal_terms.jsonl`; check `./tools/check-term-seeds`; integracao do check ao `cmd/check all` e `./tools/lab-cycle`; contratos e docs atualizados para seed como pauta, nao pagina.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `data/terms/legal_terms.jsonl`, `docs/CONTENT_QUALITY.md`, `docs/DATA_SOURCES.md`, `docs/DECISIONS.md`, `internal/checks/checks.go`, `internal/contract/term_seed_test.go`, `internal/terms/terms.go`, `tools/check-term-seeds`, `tools/lab-cycle`.
- decisoes: seed juridica deve ser `draft_only`, PT-BR, com fonte, URL oficial, data de verificacao e intencao editorial; seed invalida bloqueia laboratorio; seed nao vira pagina publica automaticamente.
- comandos executados: `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/terms`; `./tools/check-term-seeds`; `./tools/lab-cycle`; `date '+%Y-%m-%d %H:%M:%S %Z'`.
- resultados: `./tools/lab-cycle` passou; `term-seeds: pass`; `go test -count=1 ./...` passou; `check-all` passou com `term-seeds`; build continuou com `generated_pages=3 indexable_pages=1 output_dir=public`; nenhuma pagina nova foi publicada.
- falhas: RED inicial apontou ausencia de `internal/terms`; teste negativo provou que seed `published` ou sem fonte deve falhar.
- correcoes: `internal/terms.ValidateSeeds` e `ValidateSeed` implementados; check dedicado criado; seed inicial ficou em `draft_only` com LexML como fonte e nota de laboratorio.
- provas: `term_seed_test.go` exige contrato P0; `ValidateSeed` reprova `term_seed_not_draft_only`, `term_seed_without_source` e `term_seed_without_checked_at`; `lab-cycle` inclui `term-seeds: pass`.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 no ciclo 12; criar gerador de rascunho temporario em `/tmp` a partir de seed valida, com texto natural PT-BR, fonte citada, aviso informativo e reprova mecanico/raso antes de qualquer exposicao publica.
- proximo ciclo: P0 laboratorio de rascunho editorial temporario a partir de `term_seeds`, sem publicacao.
- riscos: a seed usa LexML como referencia acessivel, mas ainda nao autoriza conteudo indexavel; proximo ciclo deve manter rascunho em `/tmp` e fora do manifesto publico.

## 2026-06-09 — Ciclo 12 — P0 rascunho temporario em laboratorio

- data/hora local conferida antes do commit: `2026-06-09 09:22:19 -03`.
- ciclo: 12
- prioridade: P0
- objetivo: continuar sem parada; gerar rascunho temporario natural a partir de `term_seeds`, somente em `/tmp`, sem publicar pagina e sem alterar sitemap.
- natureza do checkpoint: rastreabilidade operacional para continuar; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo ciclo deve tornar o rascunho persistivel em camada editorial `editorial_drafts` sem virar pagina publica.
- entregas registradas no ciclo: `internal/draftlab`; teste `internal/contract/term_draft_lab_test.go`; comando `cmd/term-draft-lab`; ferramenta `./tools/lab-term-draft`; integracao ao `./tools/lab-cycle`; docs e contratos para rascunho `draft/noindex` em `/tmp`.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `cmd/term-draft-lab/`, `docs/CONTENT_QUALITY.md`, `docs/DECISIONS.md`, `docs/LAB_VALIDATION.md`, `internal/contract/term_draft_lab_test.go`, `internal/draftlab/`, `internal/terms/terms.go`, `tools/lab-term-draft`, `tools/lab-cycle`.
- decisoes: seed valida pode gerar rascunho temporario, mas nao recebe `PublicPath`, nao altera `content/pages.json`, nao entra em sitemap, nao recebe CTA e precisa passar `quality.AnalyzeText`.
- comandos executados: `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `./tools/lab-term-draft`; `./tools/lab-cycle`; `git diff -- content/pages.json`; `date '+%Y-%m-%d %H:%M:%S %Z'`.
- resultados: `./tools/lab-cycle` passou; `lab-term-draft` gerou `draft=responsabilidade-civil status=draft index=noindex source=lexml path=/tmp/.../responsabilidade-civil.txt`; `content/pages.json` sem diff; build continuou `generated_pages=3 indexable_pages=1 output_dir=public`.
- falhas: RED inicial apontou ausencia de `internal/draftlab`; patch de docs falhou por contexto antigo em `docs/LAB_VALIDATION.md`.
- correcoes: `internal/draftlab.Build` implementado com texto proprio PT-BR e gate de qualidade; comando grava apenas em `/tmp`; patch documental reaplicado em partes menores.
- provas: `term_draft_lab_test.go` exige `draft/noindex`, ausencia de rota publica, fonte LexML, aviso informativo e qualidade natural; `lab-cycle` inclui `lab-term-draft`; `git diff -- content/pages.json` sem output.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 no ciclo 13; criar append validado de `editorial_drafts` para persistir rascunho aprovado pelo laboratorio, mantendo `noindex` e sem rota publica; depois validar fluxo seed -> draft -> quality -> editorial state.
- proximo ciclo: P0 persistencia de rascunho editorial em `data/editorial/drafts.jsonl`, sem publicacao.
- riscos: texto de laboratorio ainda e rascunho generico; antes de conteudo publico real, precisa pesquisa de fonte especifica, autoria/revisao e ajuste humano natural por tema.

## 2026-06-09 — Ciclo 13 — P0 persistencia de rascunho editorial sem publicacao

- data/hora local conferida antes do commit: `2026-06-09 09:30:04 -03`.
- ciclo: 13
- prioridade: P0
- objetivo: continuar sem parada; persistir rascunho validado em `data/editorial/drafts.jsonl`, mantendo `draft/noindex`, sem URL publica, sem sitemap, sem CTA e sem marcar `/goal` como completo.
- natureza do checkpoint: rastreabilidade operacional para continuar; nenhum ciclo e final e parada; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo ciclo deve transformar o fluxo seed -> draft em fila editorial revisavel, ainda sem pagina publica.
- entregas registradas no ciclo: regra anti-finalizacao do `/goal` reforcada em `AGENTS.md`, `GOAL.md`, `docs/DECISIONS.md` e `continuity_test.go`; `internal/editorialdrafts`; teste `internal/contract/editorial_draft_store_test.go`; comando `cmd/persist-term-drafts`; ferramentas `./tools/persist-term-drafts` e `./tools/check-editorial-drafts`; draft `responsabilidade-civil` persistido em `data/editorial/drafts.jsonl`.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `cmd/persist-term-drafts/`, `data/editorial/drafts.jsonl`, `docs/CONTENT_QUALITY.md`, `docs/DATA_SOURCES.md`, `docs/DECISIONS.md`, `docs/LAB_VALIDATION.md`, `internal/architecture/architecture.go`, `internal/checks/checks.go`, `internal/contract/continuity_test.go`, `internal/contract/editorial_draft_store_test.go`, `internal/editorialdrafts/`, `tools/check-editorial-drafts`, `tools/persist-term-drafts`, `tools/lab-cycle`.
- decisoes: `/goal` nao pode ser marcado completo por ciclo parcial; conclusao exige no minimo 10 mil paginas publicas aprovadas e verificadas; rascunho persistido continua separado de `published_manifest` e nao e publicacao.
- comandos executados: `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/editorialdrafts`; `./tools/persist-term-drafts`; `./tools/check-editorial-drafts`; `./tools/lab-cycle`; `git diff -- content/pages.json public/sitemap.xml public/sitemaps/pages-0001.xml`; `date '+%Y-%m-%d %H:%M:%S %Z'`.
- resultados: `./tools/lab-cycle` passou; `editorial-drafts: pass`; `persisted=responsabilidade-civil status=draft index=noindex source=lexml`; diff publico de `content/pages.json`, `public/sitemap.xml` e `public/sitemaps/pages-0001.xml` sem output; build continuou `generated_pages=3 indexable_pages=1 output_dir=public`.
- falhas: RED inicial apontou ausencia de `internal/editorialdrafts`; tentativa anterior de finalizar goal foi corrigida por contrato, teste e objetivo ativo confirmado.
- correcoes: `Append`, `AppendIfMissing`, `Validate` e `ValidateRecord` implementados; ferramenta de persistencia ficou idempotente; check semantico garante `draft/noindex`, `public_path` vazio, fonte e qualidade textual.
- provas: `editorial_draft_store_test.go` persiste fixture e exige `status=draft`, `index_policy=noindex`, `public_path=""` e fonte; `continuity_test.go` exige nao marcar `/goal` como completo ate 10 mil paginas publicas verificadas; `lab-cycle` prova que draft persistido nao entrou no build publico.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 no ciclo 14; criar fila editorial revisavel com estados e historico de revisao/autoria para drafts persistidos, ainda sem publicar; depois preparar contrato de promocao de draft para `needs_review`.
- proximo ciclo: P0 fila editorial e historico de revisao para drafts, sem publicacao.
- riscos: ha apenas um draft persistido e ele ainda e generico; nao satisfaz meta publica, nao autoriza goal completo e nao deve sair do laboratorio sem pesquisa/revisao especifica.

## 2026-06-09 — Ciclo 14 — P0 fila editorial needs_review sem publicacao

- data/hora local conferida antes do commit: `2026-06-09 09:33:59 -03`.
- ciclo: 14
- prioridade: P0
- objetivo: continuar sem parada; criar fila editorial revisavel com autoria, motivo e historico para drafts persistidos, mantendo publicacao bloqueada.
- natureza do checkpoint: rastreabilidade operacional para continuar; nenhum ciclo e final e parada; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo ciclo deve criar contrato de promocao controlada de `needs_review` para `approved` sem ainda publicar URL.
- entregas registradas no ciclo: camada `review_queue` em `content/storage_contract.json`; arquivo `data/editorial/review_queue.jsonl`; modulo `internal/reviewqueue`; teste `internal/contract/review_queue_test.go`; comando `cmd/queue-editorial-review`; tools `./tools/queue-editorial-review` e `./tools/check-review-queue`; integracao ao `./tools/lab-cycle`; docs de fila editorial.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `cmd/queue-editorial-review/`, `content/storage_contract.json`, `data/editorial/review_queue.jsonl`, `docs/ARCHITECTURE.md`, `docs/CONTENT_QUALITY.md`, `docs/DATA_SOURCES.md`, `docs/DECISIONS.md`, `docs/LAB_VALIDATION.md`, `internal/architecture/architecture.go`, `internal/checks/checks.go`, `internal/contract/editorial_draft_store_test.go`, `internal/contract/review_queue_test.go`, `internal/editorialdrafts/`, `internal/reviewqueue/`, `tools/check-review-queue`, `tools/queue-editorial-review`, `tools/lab-cycle`.
- decisoes: fila editorial e uma camada separada de draft e manifesto publicado; estado inicial e `needs_review`; `publication_allowed=false`; `public_path` deve ficar vazio; autoria, motivo e historico sao obrigatorios.
- comandos executados: `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/reviewqueue`; `./tools/queue-editorial-review`; `./tools/check-review-queue`; `./tools/lab-cycle`; `git diff -- content/pages.json public/sitemap.xml public/sitemaps/pages-0001.xml public/index.html`; `date '+%Y-%m-%d %H:%M:%S %Z'`.
- resultados: `./tools/lab-cycle` passou; `review-queue: pass`; `queued=responsabilidade-civil status=needs_review publication_allowed=false`; build continuou `generated_pages=3 indexable_pages=1 output_dir=public`; diff publico sem output.
- falhas: RED inicial apontou ausencia de `internal/reviewqueue`; storage ainda nao tinha camada `review_queue`.
- correcoes: camada JSONL separada criada; `reviewqueue.ValidateRecord` reprova publicacao, rota publica, fonte ausente, metadados editoriais ausentes e historico ausente; `queue-editorial-review` ficou idempotente.
- provas: `review_queue_test.go` exige `needs_review`, `publication_allowed=false`, `public_path=""`, autoria e evento `queued_for_review`; teste negativo reprova registro publicavel; `lab-cycle` cobre `review-queue`.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 no ciclo 15; criar contrato de aprovacao editorial `approved` em camada separada, ainda sem publicacao, exigindo revisor, motivo, fonte e qualidade; depois preparar manifest de publicacao bloqueado.
- proximo ciclo: P0 aprovacao editorial sem publicacao e sem URL publica.
- riscos: fila tem um item e ainda nao tem revisao juridica real; aprovacao futura precisa continuar bloqueando publicacao ate fonte especifica, revisao e escala segura.

## 2026-06-09 — Ciclo 15 — P0 aprovacao editorial ainda sem publicacao

- data/hora local conferida antes do commit: `2026-06-09 09:38:15 -03`.
- ciclo: 15
- prioridade: P0
- objetivo: continuar sem parada; criar aprovacao editorial separada de publicacao, com revisor/motivo/historico e `publication_allowed=false`.
- natureza do checkpoint: rastreabilidade operacional para continuar; nenhum ciclo e final e parada; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo ciclo deve preparar manifesto de publicacao bloqueado e validar que aprovacao editorial nao entra em rotas publicas.
- entregas registradas no ciclo: camada `approved_drafts` em `content/storage_contract.json`; arquivo `data/editorial/approved_drafts.jsonl`; modulo `internal/approvals`; teste `internal/contract/editorial_approval_test.go`; comando `cmd/approve-editorial-review`; tools `./tools/approve-editorial-review` e `./tools/check-approvals`; integracao ao `./tools/lab-cycle`; docs de aprovacao sem publicacao.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `cmd/approve-editorial-review/`, `content/storage_contract.json`, `data/editorial/approved_drafts.jsonl`, `docs/ARCHITECTURE.md`, `docs/CONTENT_QUALITY.md`, `docs/DATA_SOURCES.md`, `docs/DECISIONS.md`, `docs/LAB_VALIDATION.md`, `internal/approvals/`, `internal/architecture/architecture.go`, `internal/checks/checks.go`, `internal/contract/editorial_approval_test.go`, `internal/contract/editorial_draft_store_test.go`, `tools/approve-editorial-review`, `tools/check-approvals`, `tools/lab-cycle`.
- decisoes: aprovacao editorial de laboratorio nao e publicacao; `approved_drafts` continua `noindex`, `public_path=""`, `publication_allowed=false`; publicacao futura exige contrato separado, fonte especifica, revisao juridica completa, SEO, CTA e escala segura.
- comandos executados: `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/approvals`; `./tools/approve-editorial-review`; `./tools/check-approvals`; `./tools/lab-cycle`; `git diff -- content/pages.json public/sitemap.xml public/sitemaps/pages-0001.xml public/index.html`; `date '+%Y-%m-%d %H:%M:%S %Z'`.
- resultados: `./tools/lab-cycle` passou; `approvals: pass`; `approved=responsabilidade-civil editorial_status=approved publication_allowed=false`; build continuou `generated_pages=3 indexable_pages=1 output_dir=public`; diff publico sem output.
- falhas: RED inicial apontou ausencia de `internal/approvals`; primeira tentativa de executar `./tools/approve-editorial-review` falhou por corrida com `chmod` em chamada paralela, corrigida por execucao sequencial.
- correcoes: `approvals.Approve`, `Validate`, `ValidateRecord` implementados; ferramenta idempotente criada; `lab-cycle` passou a cobrir `approvals`.
- provas: `editorial_approval_test.go` exige `editorial_status=approved`, `publication_allowed=false`, `public_path=""`, revisor e evento `approved_editorial_only`; teste negativo reprova registro publicavel; `lab-cycle` cobre `approvals`.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 no ciclo 16; criar `publication_manifest_blocked` ou contrato equivalente para preparar promocao publica sem liberar URL, medindo requisitos ausentes para publicar em escala; depois avançar para blueprint finito de lotes aprovaveis.
- proximo ciclo: P0 manifesto de publicacao bloqueado e lista de requisitos faltantes para URL publica.
- riscos: aprovacao e de laboratorio, nao revisao juridica real; nao satisfaz meta de 10 mil paginas publicas; goal permanece ativo.

## 2026-06-09 — Ciclo 16 — P0 manifesto de publicacao bloqueada

- data/hora local conferida antes do commit: `2026-06-09 09:45:23 -03`.
- ciclo: 16
- prioridade: P0
- objetivo: continuar sem parada; criar manifesto de publicacao bloqueada para provar que aprovacao editorial de laboratorio nao vira URL publica, sitemap, CTA ou indexacao.
- natureza do checkpoint: rastreabilidade operacional para continuar; nenhum ciclo e final e parada; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo ciclo deve pesquisar e registrar termos juridicos mais buscados por humanos, com alta intencao de contratacao juridica 100% digital, antes de qualquer conteudo publico.
- entregas registradas no ciclo: camada `publication_blockers` em `content/storage_contract.json`; arquivo `data/editorial/publication_blockers.jsonl`; modulo `internal/publicationblockers`; teste `internal/contract/publication_blockers_test.go`; comando `cmd/block-publication`; ferramentas `./tools/block-publication` e `./tools/check-publication-blockers`; integracao ao `./tools/lab-cycle`; contratos reforcados para migracao segura do laboratorio para host/produção controlada e priorizacao de termos com contratacao online.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `cmd/block-publication/`, `content/storage_contract.json`, `data/editorial/publication_blockers.jsonl`, `docs/ARCHITECTURE.md`, `docs/CONTENT_QUALITY.md`, `docs/DATA_SOURCES.md`, `docs/DECISIONS.md`, `docs/LAB_VALIDATION.md`, `internal/architecture/architecture.go`, `internal/checks/checks.go`, `internal/contract/editorial_draft_store_test.go`, `internal/contract/publication_blockers_test.go`, `internal/publicationblockers/`, `tools/block-publication`, `tools/check-publication-blockers`, `tools/lab-cycle`.
- decisoes: aprovacao editorial aprovada em laboratorio deve gerar `publication_status=blocked`, `publication_allowed=false` e `public_path=""` ate faltar zero requisito; laboratorio nao e destino final; alta intencao comercial significa contratacao juridica online, por WhatsApp e envio remoto de documentos, nao fluxo predominantemente presencial.
- comandos executados: `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/publicationblockers`; `./tools/block-publication`; `./tools/check-publication-blockers`; `./tools/lab-cycle`; `git diff -- content/pages.json public/sitemap.xml public/sitemaps/pages-0001.xml public/index.html`; `git diff --check`; `date '+%Y-%m-%d %H:%M:%S %Z'`; `git diff --stat`; `git status --short`.
- resultados: testes especificos passaram; `./tools/block-publication` retornou `existing=responsabilidade-civil publication_allowed=false public_path=none`; `./tools/check-publication-blockers` retornou `publication-blockers: pass`; `./tools/lab-cycle` passou com `publication-blockers: pass`, build `generated_pages=3 indexable_pages=1 output_dir=public`; diff publico sem output; `git diff --check` sem output.
- falhas: RED inicial do ciclo apontou ausencia de `internal/publicationblockers`; o contrato anterior permitia aprovacao editorial sem manifesto explicito de bloqueio; os documentos ainda nao distinguiam alta intencao digital de fluxo presencial.
- correcoes: `internal/publicationblockers` implementado com validacao semantica; check dedicado integrado; manifesto idempotente criado; contratos atualizados para exigir migracao segura e termos com contratacao online 100% digital como prioridade inicial.
- provas: `publication_blockers_test.go` exige que draft aprovado fique fora de rota publica; teste negativo reprova bloqueador sem requisitos faltantes; `data/editorial/publication_blockers.jsonl` registra `responsabilidade-civil` bloqueado; `lab-cycle` cobre `publication-blockers`; diff publico prova que nada entrou em `content/pages.json` nem sitemap.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 no ciclo 17; pesquisar em fontes publicas confiaveis e atuais termos juridicos humanos de alta intencao digital, separar evidencia de demanda, adequacao a atendimento online, risco juridico e fonte oficial; criar banco leve/validador para candidatos de termos sem publicar conteudo.
- proximo ciclo: P0 pesquisa e camada de candidatos de termos juridicos de alta intencao online, sem publicacao.
- riscos: ainda ha apenas um termo de laboratorio; nenhum conteudo publico juridico novo foi liberado; pesquisa de demanda deve usar fonte atual e nao scraping cego; goal permanece ativo.

## 2026-06-09 — Ciclo 17 — P0 candidatos de termos de alta intencao digital

- data/hora local conferida antes do commit: `2026-06-09 09:55:19 -03`.
- ciclo: 17
- prioridade: P0
- objetivo: continuar sem parada; registrar caminhos seguros de pesquisa de demanda humana e criar banco leve de candidatos de termos juridicos com alta intencao de contratacao 100% digital, sem publicar conteudo.
- natureza do checkpoint: rastreabilidade operacional para continuar; nenhum ciclo e final e parada; nao e aceite final, nao e ordem de parada.
- continuidade: continuar P0; proximo ciclo deve transformar os candidatos mais fortes em seeds/drafts priorizados com fonte juridica especifica, ainda sem publicar, e preparar migracao segura para host/produção controlada quando os gates permanecerem verdes.
- fontes/caminhos seguros registrados: Google Trends Explore Brasil `https://trends.google.com.br/trends/explore?geo=BR`; Ajuda do Google Trends `https://support.google.com/trends/answer/4359550?hl=pt-BR`; FAQ Google Trends `https://support.google.com/trends/answer/4365533?hl=pt-br`; Google Search Central `https://developers.google.com/search/docs/monitor-debug/trends-start`; Datajud `https://datajud-wiki.cnj.jus.br/api-publica/acesso/`; Camara Dados Abertos `https://dadosabertos.camara.leg.br/swagger/api.html`; STJ Dados Abertos `https://dadosabertos.web.stj.jus.br/`; CNJ/e-Notariado `https://www.cnj.jus.br/e-notariado-completa-tres-anos-com-mais-de-15-milhao-de-atos-online/`; fontes gov.br/Planalto/ANS/MJ usadas como lastro dos candidatos.
- entregas registradas no ciclo: `docs/data-sources/search-demand.md`; camada `term_intent_candidates` em `content/storage_contract.json`; arquivo `data/terms/intent_candidates.jsonl` com 10 candidatos; modulo `internal/termintents`; teste `internal/contract/term_intent_candidates_test.go`; ferramenta `./tools/check-term-intent-candidates`; integracao ao `./tools/lab-cycle`; contratos atualizados com Google Trends/Search Central como sinal direcional, nao fonte juridica.
- arquivos alterados: `AGENTS.md`, `GOAL.md`, `CHECKPOINT.md`, `content/storage_contract.json`, `data/terms/intent_candidates.jsonl`, `docs/ARCHITECTURE.md`, `docs/CONTENT_QUALITY.md`, `docs/DATA_SOURCES.md`, `docs/DECISIONS.md`, `docs/LAB_VALIDATION.md`, `docs/data-sources/search-demand.md`, `internal/architecture/architecture.go`, `internal/checks/checks.go`, `internal/contract/editorial_draft_store_test.go`, `internal/contract/term_intent_candidates_test.go`, `internal/storage/storage.go`, `internal/termintents/`, `tools/check-term-intent-candidates`, `tools/lab-cycle`.
- decisoes: Google Trends e Search Central sao caminhos seguros para demanda/metodologia, mas nao volume absoluto, nao fonte juridica e nao autorizacao de publicacao; candidato de termo precisa `digital_only`, `whatsapp_cta_intent=high`, fonte oficial, evidencia de demanda e `publication_allowed=false`; termos predominantemente presenciais ficam fora da prioridade inicial.
- comandos executados: `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract` como RED; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/contract`; `env GOCACHE=/tmp/wiki-go-build-cache go test -count=1 ./internal/termintents`; `./tools/check-term-intent-candidates`; `./tools/check-storage-contract`; `./tools/lab-cycle`; `git diff -- content/pages.json public/sitemap.xml public/sitemaps/pages-0001.xml public/index.html`; `git diff --check`; `wc -l data/terms/intent_candidates.jsonl`; `date '+%Y-%m-%d %H:%M:%S %Z'`; `git status --short`.
- resultados: RED falhou por pacote ausente `portaljuridico/internal/termintents`; apos implementacao, `./tools/check-term-intent-candidates` retornou `term-intent-candidates: pass`; `./tools/check-storage-contract` retornou `storage-contract: pass`; `./tools/lab-cycle` passou com `term-intent-candidates: pass` e build `generated_pages=3 indexable_pages=1 output_dir=public`; `wc -l` confirmou 10 candidatos; diff publico sem output; `git diff --check` sem output.
- falhas: o primeiro wrapper `tools/check-term-intent-candidates` falhou com `go: no packages loaded from ./cmd/check` por nao exportar `GOCACHE`; a camada ainda nao existia no RED; demanda por Google Trends e direcional e nao deve ser superinterpretada.
- correcoes: wrapper corrigido para o padrao do repo com `GOCACHE`; `internal/termintents` implementado com validacao de URL segura de demanda, fonte oficial, modo digital, CTA alto e bloqueio de publicacao; storage contract passou a exigir a camada `term_intent_candidates`.
- provas: `term_intent_candidates_test.go` exige ao menos 8 candidatos digitais sem publicacao; teste negativo reprova termo presencial/sem demanda/publicavel; `data/terms/intent_candidates.jsonl` guarda 10 candidatos; `lab-cycle` cobre o novo check; diff publico prova que nenhum candidato virou pagina.
- commit: este ciclo deve ser persistido em Git apos validacao final desta entrada.
- plano de continuidade: continuar P0 no ciclo 18; ranquear candidatos por risco, adequacao digital e capacidade de fonte; promover um lote pequeno de candidatos para `term_seeds`/drafts priorizados com fonte especifica e bloqueio de publicacao; preparar check para exigir que qualquer seed gerada de candidato preserve a evidencia de demanda e nao vire CTA/publicacao.
- proximo ciclo: P0 ranking/promocao controlada de candidatos para seeds/drafts priorizados, sem publicacao.
- riscos: Google Trends nao fornece volume absoluto; algumas fontes oficiais sao lastro inicial e ainda exigem pesquisa juridica especifica antes de conteudo; nenhum candidato satisfaz publicacao sozinho; goal permanece ativo.
