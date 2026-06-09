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
