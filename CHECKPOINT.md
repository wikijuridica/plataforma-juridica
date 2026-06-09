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
