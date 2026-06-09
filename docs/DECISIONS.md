# DECISIONS.md

## 2026-06-09 — Go como base propria do portal

Decisao: usar Go e biblioteca padrao como base do projeto.

Motivos:
- linguagem compilada, tipada e adequada a servicos com concorrencia;
- biblioteca padrao suficiente para HTTP, XML, HTML escaping, arquivos e testes;
- binario proprio e operacao sem framework frontend;
- melhor alinhamento com geracao on demand propria e cache controlado pelo projeto.

Evidencia documental consultada:
- `https://go.dev/doc/`: documenta Go como linguagem eficiente, compilada, tipada, com concorrencia e toolchain oficial.
- `https://go.dev/doc/modules/layout`: recomenda `internal/` e `cmd/` para projetos de servidor.
- `https://pkg.go.dev/net/http`: biblioteca padrao para cliente e servidor HTTP.
- `https://pkg.go.dev/encoding/xml`: biblioteca padrao para XML, usada em sitemaps.
- `https://pkg.go.dev/testing`: biblioteca padrao para testes automatizados.

Alternativas avaliadas:
- Python com biblioteca padrao: suficiente para prototipar validadores, mas menos alinhado ao objetivo de portal massivo com binario proprio e geracao HTTP concorrente;
- Next.js ou frameworks frontend: proibidos pelo contrato P0.

Consequencia: qualquer artefato Python criado no ciclo foi descartado; o contrato e a implementacao passam a ser Go-first.

## 2026-06-09 — Geracao on demand propria obrigatoria

Decisao: o portal deve gerar paginas sob demanda com codigo proprio em `internal/ondemand`, sem Next.js, ISR terceirizado ou framework equivalente.

Motivos:
- o projeto quer autonomia sobre roteamento, cache e politica de indexacao;
- rotas futuras podem chegar a milhoes de URLs e nao devem depender de build total;
- o primeiro response de pagina publica deve conter HTML textual completo.

Consequencia: build estatico e permitido como ferramenta operacional, mas nao substitui o gerador on demand.

## 2026-06-09 — Checkpoint nao encerra trabalho

Decisao: checkpoint e rastreabilidade operacional para continuar, nao entrega final nem ordem de parada.

Motivos:
- o projeto tem escopo continuo e grande;
- cada ciclo precisa deixar contexto, provas, falhas e proximo plano;
- o agente deve continuar trabalhando enquanto nao houver bloqueio P0 real.

Consequencia: `CHECKPOINT.md` deve registrar plano de continuidade explicito e nao pode ser usado como substituto da definicao de pronto.

Adendo operacional: sempre planejar o proximo passo e continuar executando. Uma resposta no thread ou um checkpoint nao encerram a missao.

Adendo de persistencia: cada ciclo deve ser commitado apos validacao, incluindo `CHECKPOINT.md`, para preservar continuidade em Git.
Antes do commit, registrar hora local e numero do ciclo no checkpoint.

Adendo de autonomia: Codex atua como engenheiro senior, arquiteto e criador de conteudo juridico. Se houver trabalho no escopo ou bug identificado, deve continuar e corrigir sem pedir aprovacao para decisao normal de engenharia.

## 2026-06-09 — Meta de 10 mil paginas com CTA subordinado a qualidade

Decisao: a meta de produto inclui no minimo 10 mil paginas juridicas informativas, com alta intencao de contratar advogado e CTA proprio de WhatsApp quando apropriado.

Motivos:
- o portal deve ter escala nacional e alta intencao comercial;
- a arquitetura precisa suportar grande volume sem URLs infinitas;
- CTA comercial nao pode sacrificar fonte, revisao, utilidade e seguranca juridica.

Consequencia: durante P0 o projeto cria plano de blueprints e politica de CTA, mas nao publica 10 mil paginas juridicas. A publicacao em escala depende dos gates P0/P1/P2.

Adendo de intencao digital: alta intencao comercial significa potencial de contratacao juridica 100% online. O primeiro lote de termos deve favorecer problemas juridicos que podem ser triados, documentados e contratados por canais digitais, especialmente WhatsApp. Demanda alta com dependencia presencial predominante nao deve ser tratada como prioridade inicial.

## 2026-06-09 — Laboratorio antes de conteudo e mudancas P0

Decisao: toda mudanca P0/P1 deve passar por laboratorio de testes, validacao, refinamento e reteste.

Motivos:
- scripts isolados podem passar e ainda deixar falha de contrato;
- conteudo juridico exige fonte correta e escrita natural;
- escala de 10 mil paginas sem laboratorio vira risco de spam, duplicidade e baixa qualidade.

Consequencia: `tools/lab-cycle` passa a ser o comando minimo de laboratorio, e conteudo juridico em escala permanece bloqueado ate pesquisa de fonte, documentacao, revisao e gates.

## 2026-06-09 — Fontes como referencia, nao scraping

Decisao: fontes oficiais e APIs publicas registradas sao referencia/proveniencia, nao alvo de scraping, clone ou espelhamento.

Motivos:
- o portal deve ser unico e editorialmente proprio;
- copiar massa de dados ou estrutura oficial cria risco de spam e baixa qualidade;
- conteudo juridico precisa ser natural, humano e contextualizado.

Consequencia: qualquer uso de fonte oficial exige pesquisa critica, registro de proveniencia e redacao propria. Ingestao automatica permanece bloqueada em P0.

## 2026-06-09 — HTML publico leve para bots valiosos

Decisao: pagina publica deve permanecer leve por contrato, com HTML textual completo, CSS minimo e sem runtime cliente.

Motivos:
- Googlebot, OAI-SearchBot e bots valiosos precisam rastrear conteudo textual com baixo custo;
- escala massiva amplifica qualquer excesso de bytes, bundle ou hidratacao;
- framework frontend, payload JavaScript e ferramenta pesada no HTML publico contrariam P0/P1.

Consequencia: `./tools/check-performance-budget` reprova HTML acima do orcamento, `<script>`, referencias a `.js/.mjs/.wasm`, assets pesados, `modulepreload`, import map, marcadores de hidratacao e sinais de runtime/framework. UI publica deve ser resolvida com HTML semantico e CSS minimo.

## 2026-06-09 — Algoritmos auditaveis e explicaveis

Decisao: gates e algoritmos do projeto devem ser inteligentes, auditaveis e explicaveis.

Motivos:
- escala massiva torna heuristica burra perigosa e cara;
- qualidade juridica depende de decisoes rastreaveis, nao de aprovacoes opacas;
- falsos positivos e falsos negativos precisam virar melhoria de algoritmo, nao excecao manual permanente.

Consequencia: quando um algoritmo estiver ingenuo, caro, opaco, permissivo demais ou agressivo demais, o agente deve escrever teste que reproduza a falha, melhorar a regra, manter mensagens especificas de reprovação e registrar a decisao. Heuristicas simples sao permitidas como etapa inicial, mas devem evoluir quando houver sinal melhor disponivel.

## 2026-06-09 — Conteudo publico em PT-BR correto

Decisao: todo conteudo visivel ao publico deve ser escrito em PT-BR com grafia correta, acentuacao correta, pontuacao clara e linguagem natural.

Motivos:
- o portal serve primeiro a humanos no Brasil;
- texto publico sem acento ou com grafia tecnica empobrece confianca editorial;
- qualidade juridica depende de clareza, naturalidade e revisao humana.

Consequencia: rascunhos tecnicos internos podem usar texto operacional sem polimento, mas paginas, metadados visiveis, CTA, navegacao e rodape publicos devem ser PT-BR correto. O algoritmo de normalizacao de qualidade deve preservar letras acentuadas para nao degradar os gates em portugues.

## 2026-06-09 — Orcamento SERP conservador baseado em pesquisa Google

Decisao: o projeto usa limites internos conservadores para `title`, metadescricao e snippet, sem afirmar que sao limites oficiais fixos do Google.

Pesquisa oficial feita no dia da sessao:
- `https://developers.google.com/search/docs/appearance/title-link?hl=pt-BR`
- `https://developers.google.com/search/docs/appearance/snippet?hl=pt-br`
- `https://developers.google.com/search/docs/essentials/technical`

Motivos:
- Google informa que nao ha limite fixo oficial para `<title>` e metadescricoes;
- os textos podem ser truncados conforme a largura do dispositivo;
- conteudo indexavel exige acesso do Googlebot, HTTP 200 e conteudo que nao viole politicas de spam.

Consequencia: `title` deve ter 20 a 65 caracteres Unicode, metadescricao deve ter 70 a 160 caracteres Unicode e paginas indexaveis recebem `max-snippet:160`. Quando faltarem dados atuais sobre Googlebot ou superficie de busca, pesquisar a Central da Pesquisa Google no dia da sessao.

## 2026-06-09 — Laboratorio contra conteudo mecanico antes do Googlebot

Decisao: conteudo raso, mecanico ou criado por permutacao de palavras-chave deve ser detectado em laboratorio antes de qualquer exposicao ao Googlebot.

Motivos:
- a meta de escala nao autoriza spam;
- humanos e Googlebot precisam receber conteudo natural, util e especifico;
- scripts isolados podem ser enganados por texto longo, mas repetitivo.

Consequencia: `./tools/lab-content-quality` cria textos temporarios em `/tmp`, aprova texto natural e reprova texto mecanico. `./tools/check-mechanical-content` bloqueia paginas indexaveis com sinais de thin content, keyword stuffing, baixa diversidade lexical, frases repetidas ou permutacao mecanica. Falsos positivos e falsos negativos devem virar melhoria de algoritmo.

## 2026-06-09 — CPU baixo no runtime publico

Decisao: CPU deve ser reservada para trafego legitimo, Googlebot, OAI-SearchBot e bots valiosos no runtime publico.

Motivos:
- producao precisa atender muitas URLs e bots sem desperdiçar CPU;
- renderizacao publica deve ser previsivel, local e barata;
- laboratorio pode ser mais pesado quando necessario para provar qualidade, mas isso nao pode virar dependencia do atendimento publico.

Consequencia: `./tools/check-cpu-budget` escaneia caminhos de runtime publico/producao e reprova execucao externa, chamadas de rede, sleeps ou loops sem limite. Testes, build, laboratorio e auditorias podem usar comandos mais pesados quando forem necessarios e registrados.

## 2026-06-09 — Modo /goal nao permite parada operacional

Decisao: quando o Codex estiver em modo `/goal`, nao deve parar enquanto houver trabalho no escopo e nao houver bloqueio P0 real comprovado.

Motivos:
- checkpoint e resposta no thread sao rastreabilidade, nao ordem de parada;
- o projeto exige continuidade autonoma e documentada;
- o proximo passo deve sobreviver a compactacao de contexto, troca de sessao e perda de chat.

Consequencia: a primeira linha de `AGENTS.md` registra a regra literal de `/goal`; `GOAL.md`, `CHECKPOINT.md` e documentos de decisao devem manter proximo passo planejado e continuidade explicita. O agente deve executar o proximo passo documentado sem aguardar nova cobranca do usuario, salvo bloqueio P0 real comprovado.

## 2026-06-09 — Sem mascarar pendencia e com rede quando necessaria

Decisao: nenhuma pendencia pode ser mascarada como entrega, e rede deve ser usada quando for necessaria para resolver o escopo.

Motivos:
- registrar bloqueio ou "pendente" como avanço desperdiça ciclo e nao resolve o projeto;
- fontes oficiais, robots.txt, termos de uso, APIs publicas e documentos atuais exigem verificacao real quando impactam decisao;
- sandbox sem rede nao e desculpa para parar quando existe mecanismo de escalonamento pela ferramenta.

Consequencia: Nada deve ser deixado para o futuro por conveniencia. Tudo que estiver no escopo e para esta sessao. Alternativa so e valida quando resolve o requisito ou produz prova executavel que desbloqueia o requisito; artefato que nao sera usado nao conta como entrega. Se a rede do sandbox falhar e a rede for necessaria, o agente deve repetir o comando com `sandbox_permissions=require_escalated` e justificativa objetiva pela ferramenta, sem perguntar no chat. Se o escalonamento for negado, a negativa vira bloqueio real com evidencia, nao pendencia mascarada.

## 2026-06-09 — Banco leve separado para ingestao de termos

Decisao: ingestao de termos juridicos pode iniciar conteudos apenas como semente de rascunho, usando banco leve proprio em JSONL e camadas separadas.

Motivos:
- termos juridicos sao uma boa unidade inicial para organizar pautas e glossario sem criar spam;
- fonte auditada, termos de uso, snapshot oficial, rascunho editorial e conteudo publicado tem riscos e contratos diferentes;
- usar Go e arquivos JSONL evita dependencia externa, SDK, SaaS ou banco pesado durante P0.

Consequencia: `content/storage_contract.json` define `term_seeds`, `source_audits`, `source_snapshots`, `editorial_drafts` e `published_manifest`. `./tools/check-storage-contract` reprova camada ausente, caminho compartilhado, arquivo inexistente, record pesado, camada diretamente indexavel ou mistura entre fonte bruta e texto editorial. Termos juridicos so podem iniciar `draft_only` ate fonte, revisao, qualidade, SEO, CTA e checkpoint passarem.

## 2026-06-09 — URL oficial do Portal da Legislacao

Decisao: para o registro Planalto, a URL oficial de pesquisa do Portal da Legislacao passa a ser `https://legislacao.presidencia.gov.br/`, com atos em `https://legislacao.presidencia.gov.br/atos/?...`.

Motivos:
- pagina publica confiavel do proprio `gov.br` para o servico "Realizar pesquisa de legislacao no Portal da Legislacao" aponta o canal Web para `legislacao.presidencia.gov.br`;
- o resultado tambem descreve que a pesquisa acessa atos normativos federais de hierarquia superior e a base REFLEGIS;
- `www.planalto.gov.br` e `www4.planalto.gov.br` sao referencias historicas/auxiliares, mas falharam na auditoria HTTP local desta sessao.

Consequencia: `content/source_registry.json` registra `https://legislacao.presidencia.gov.br/` como base oficial e mantem ingestao bloqueada. A auditoria local teve timeout em `legislacao.presidencia.gov.br`, reset em `www.planalto.gov.br` e timeout em `www4.planalto.gov.br`; portanto a URL foi confirmada documentalmente, mas coleta automatica continua proibida.

## 2026-06-09 — Term seeds como pauta, nao pagina

Decisao: `data/terms/legal_terms.jsonl` pode conter sementes de termos juridicos para laboratorio editorial, desde que cada registro seja `draft_only`, PT-BR, tenha fonte, URL oficial, data de verificacao e intencao editorial.

Motivos:
- termos sao unidade leve para iniciar pauta e glossario sem criar pagina automaticamente;
- uma seed sem fonte vira risco de conteudo juridico inventado;
- um termo publicado direto seria spam ou thin content.

Consequencia: `internal/terms` e `./tools/check-term-seeds` validam as seeds. O laboratorio pode usar essas seeds para rascunhos temporarios, mas nenhuma seed vira pagina publica, CTA ou URL indexavel sem passar pelos demais gates.

## 2026-06-09 — Rascunho temporario antes de pagina publica

Decisao: seed valida pode gerar rascunho temporario em `/tmp` por `./tools/lab-term-draft`, sempre `draft/noindex` e fora do manifesto publico.

Motivos:
- permite testar escrita natural, fonte e aviso informativo antes de tocar o pipeline publico;
- impede que uma seed vire pagina mecanica ou thin content;
- preserva baixo CPU e HTML leve em producao, deixando o experimento no laboratorio.

Consequencia: `internal/draftlab` gera rascunho proprio em PT-BR e roda `quality.AnalyzeText`. O rascunho nao recebe URL publica, nao entra em sitemap, nao recebe CTA e nao altera `content/pages.json`.

## 2026-06-09 — Goal nao termina em ciclo parcial

Decisao: nenhum ciclo e final e parada. `/goal` nao pode ser marcado como completo por checkpoint, commit, laboratorio verde, P0 parcial, seed ou rascunho.

Motivos:
- o objetivo contratado inclui plataforma grande e minimo de 10 mil paginas publicas juridicas aprovadas;
- o projeto ainda esta em P0/laboratorio e nao saiu para publicacao em escala;
- marcar goal como completo em ciclo parcial distorce a missao e causa parada indevida.

Consequencia: `AGENTS.md`, `GOAL.md` e testes de contrato exigem nao marcar `/goal` como completo ate a meta publica minima estar verificada: no minimo 10 mil paginas publicas aprovadas, indexaveis, com fonte, revisao, qualidade, CTA quando cabivel, sitemap/canonical/robots corretos e validacao completa.

## 2026-06-09 — Draft editorial persistido ainda nao e publicacao

Decisao: rascunho validado pode ser persistido em `data/editorial/drafts.jsonl`, mas continua `draft/noindex`, sem rota publica, sem sitemap e sem CTA.

Motivos:
- persistir rascunho reduz perda de contexto entre ciclos sem expor conteudo incompleto;
- separar `editorial_drafts` de `published_manifest` impede que laboratorio seja confundido com pagina publica;
- validacao de qualidade deve acompanhar a persistencia, nao ficar apenas no texto temporario.

Consequencia: `internal/editorialdrafts`, `./tools/persist-term-drafts` e `./tools/check-editorial-drafts` controlam a camada editorial. O fluxo atual e seed -> draft temporario -> draft persistido; ainda nao existe publicacao indexavel desse conteudo.

## 2026-06-09 — Fila editorial antes de promocao

Decisao: draft persistido deve entrar em `data/editorial/review_queue.jsonl` como `needs_review`, com autoria, motivo e historico, ainda com `publication_allowed=false`.

Motivos:
- revisao precisa ser rastreavel antes de qualquer promocao;
- separar fila de revisao de draft e manifesto publicado evita confundir laboratorio com publicacao;
- a meta de 10 mil paginas exige processo repetivel, nao improviso por pagina.

Consequencia: `internal/reviewqueue`, `./tools/queue-editorial-review` e `./tools/check-review-queue` validam a fila. Nenhuma entrada da fila cria URL publica, sitemap, canonical ou CTA.

## 2026-06-09 — Aprovacao editorial ainda nao publica

Decisao: aprovacao editorial de laboratorio fica em `data/editorial/approved_drafts.jsonl` e continua com `publication_allowed=false`, `public_path` vazio e `noindex`.

Motivos:
- aprovacao de texto nao equivale a publicacao tecnica;
- publicacao publica exige fonte especifica, revisao juridica completa, SEO, CTA, sitemap, canonical, robots e escala segura;
- separar aprovacao editorial de manifesto publicado evita salto indevido do laboratorio para Googlebot.

Consequencia: `internal/approvals`, `./tools/approve-editorial-review` e `./tools/check-approvals` validam aprovacao sem publicacao. O proximo contrato deve preparar manifest de publicacao bloqueado antes de qualquer URL publica.

## 2026-06-09 — Manifesto de publicacao bloqueada

Decisao: aprovacao editorial deve gerar manifesto de publicacao bloqueada antes de qualquer URL publica.

Motivos:
- explicitar requisitos faltantes evita mascarar laboratorio como producao;
- publicacao em escala exige lista objetiva de pendencias por termo;
- o fluxo precisa saber quando migrar do laboratorio para host/produção controlada.

Consequencia: `internal/publicationblockers`, `./tools/block-publication` e `./tools/check-publication-blockers` registram `publication_status=blocked`, `publication_allowed=false`, `public_path=""` e requisitos como pesquisa de fonte especifica, revisor juridico real, intencao unica, canonical, SEO, CTA e capacidade de lote.

## 2026-06-09 — Termos humanos de alta intencao antes de conteudo publico

Decisao: quando os gates passarem com seguranca, o agente deve migrar o fluxo para host/produção controlada e priorizar termos juridicos mais pesquisados por humanos e com alta intencao de contratar advogado online.

Motivos:
- o objetivo e portal juridico util e comercial, nao laboratorio infinito;
- termos reais de busca humana reduzem risco de pagina artificial;
- CTA WhatsApp so faz sentido em temas com intencao de contratacao e utilidade juridica clara;
- servicos juridicos digitais sao prioridade: triagem, envio de documentos e contratacao remota devem ser viaveis sem atendimento presencial como padrao.

Consequencia: antes de criar conteudo publico, deve haver pesquisa atual de demanda/intencao, registro no banco leve, fonte confiavel, intencao unica e bloqueio contra spam. Essa regra nao libera publicacao automatica; ela define o proximo movimento apos gates seguros.

## 2026-06-09 — Caminhos seguros para demanda humana

Decisao: usar Google Trends e Google Search Central como caminhos seguros para orientar demanda humana e estrategia, sem tratar esses caminhos como fonte juridica ou como autorizacao de publicacao.

Caminhos registrados:
- Google Trends Explore Brasil: `https://trends.google.com.br/trends/explore?geo=BR`
- Ajuda do Google Trends sobre comparacao: `https://support.google.com/trends/answer/4359550?hl=pt-BR`
- FAQ de dados do Google Trends: `https://support.google.com/trends/answer/4365533?hl=pt-br`
- Google Search Central sobre Trends: `https://developers.google.com/search/docs/monitor-debug/trends-start`

Motivos:
- o Google Trends permite comparar termos e observar interesse de busca, mas seus dados sao normalizados e direcionais;
- a Central da Pesquisa Google orienta usar Trends para estrategia de conteudo sem escrever apenas porque algo esta em alta;
- o projeto precisa escolher termos com demanda humana real antes de criar pauta publica.

Consequencia: `data/terms/intent_candidates.jsonl` deve registrar URL de comparacao, fonte juridica oficial, adequacao a contratacao 100% digital e bloqueio de publicacao. Nenhum candidato vira pagina publica sem novo ciclo de fonte, revisao, qualidade, SEO e CTA.

## 2026-06-09 — Ranking refinavel de termos, sem confiar no primeiro sinal

Decisao: promover candidatos de alta intencao para `term_seeds` por ranking refinavel e testado, nao por lista fixa nem confianca cega em Google Trends.

Motivos:
- um unico sinal de demanda pode ser enganoso, sazonal ou amplo demais;
- termo presencial pode ter volume, mas baixa adequacao ao produto 100% digital;
- fonte nao oficial ou fraca aumenta risco juridico e editorial;
- concentrar todos os termos em uma area reduz aprendizado e escala.

Consequencia: `internal/termpromotion` pontua candidatos, penaliza fonte nao oficial e modo presencial, exige diversidade de areas e preserva evidencia de demanda na seed. `./tools/promote-term-candidates` e `./tools/check-promoted-term-seeds` mantem seeds como `draft_only`, sem URL publica e sem CTA publico.

## 2026-06-09 — Rascunho PT-BR natural tambem no laboratorio

Decisao: rascunhos editoriais persistidos devem usar grafia natural em PT-BR, mesmo antes de publicacao.

Motivos:
- laboratorio e onde erro mecanico deve aparecer, nao no Googlebot;
- texto sem acento ou sem conectivos naturais e sinal de algoritmo burro;
- persistir rascunho ruim aumenta risco de promover conteudo fraco depois.

Consequencia: `internal/draftlab` aplica termo de exibicao natural; `./tools/refresh-editorial-drafts` regenera drafts persistidos apos refinamento; fila de revisao atualiza registros existentes sem liberar publicacao.

## 2026-06-09 — Fonte especifica bloqueia aprovacao

Decisao: termo priorizado com fonte ampla, institucional ou generica nao pode ser aprovado nem publicado ate haver manifesto de fonte especifica resolvido.

Motivos:
- demanda humana e CTA alto nao substituem base juridica correta;
- fonte institucional pode orientar pesquisa, mas texto publico precisa regra, norma, artigo, requisito ou limite juridico aplicavel;
- bloquear aprovacao evita que rascunho de laboratorio vire conteudo publico por engano.

Consequencia: `data/editorial/source_blockers.jsonl` registra `approval_allowed=false`, `publication_allowed=false`, `public_path=""`, requisitos faltantes e proxima pesquisa por termo. `./tools/check-source-specificity-blockers` entra no laboratorio.

## 2026-06-09 — Pesquisa editorial manual antes de scripts de termo

Decisao: a estrategia principal para escolher termos de alta intencao passa a ser pesquisa editorial manual na web, nao script conservador ou gerador de termos.

Motivos:
- termos juridicos de contratacao digital dependem de leitura de intencao humana, nao apenas heuristica;
- Google Trends e util como orientacao direcional, mas nao substitui julgamento editorial;
- fontes oficiais dao autoridade, mas nao sao fonte de demanda nem autorizam clone;
- a meta de 10 mil paginas exige uma base leve, organizada e escalavel sem template mecanico.

Consequencia: `data/research/high_intent_terms.jsonl` vira o banco leve de pesquisa manual; `data/editorial/content_briefs.jsonl` inicia conteudos como briefs nao publicaveis; `./tools/check-manual-keyword-research` e `./tools/check-content-briefs` entram no laboratorio. Nenhum brief vira pagina publica sem fonte especifica, revisao, qualidade, SEO, CTA e nova validacao.

## 2026-06-09 — Rascunhos autorais antes de qualquer pagina de alta intencao

Decisao: briefs pesquisados podem virar rascunhos autorais no banco leve, mas esses rascunhos continuam bloqueados para publicacao ate passarem por fonte especifica, revisao editorial, SEO/crawl, qualidade, CTA e checkpoint.

Motivos:
- o projeto precisa comecar a construir conteudo sem gerar spam, clone ou pagina mecanica;
- Googlebot e humanos devem receber apenas paginas com valor proprio, nao texto de molde;
- CTA WhatsApp e critico, mas deve ser contextual e responsavel;
- a escala de 10 mil paginas exige banco organizado antes de renderizacao publica.

Consequencia: `data/editorial/authorial_drafts.jsonl` vira camada propria de rascunho autoral; `internal/authorialdrafts` e `./tools/check-authorial-content-drafts` reprovam abertura repetida, shape de secoes reaproveitado, heading generico, CTA raso, fonte ausente, publicacao permitida e path publico.

## 2026-06-09 — Resolução de fonte específica antes de pré-publicação

Decisao: fonte específica resolvida deve virar manifesto próprio antes de qualquer contrato de publicação, sem remover automaticamente bloqueadores nem criar URL pública.

Motivos:
- o projeto precisa diferenciar fonte ampla de fonte realmente útil para revisar o texto;
- fontes oficiais devem dar autoridade e proveniência, não texto copiado;
- Googlebot só deve ver página depois de fonte, revisão, qualidade, SEO/crawl e CTA passarem;
- a escala de 10 mil páginas exige rastreabilidade por termo.

Consequencia: `data/editorial/source_resolutions.jsonl` registra a primeira resolução para `negativa-cobertura-plano-saude`; `internal/sourceresolutions` e `./tools/check-source-specificity-resolutions` exigem lei primária, regra de cobertura, fontes oficiais específicas, score mínimo, `noindex`, `publication_allowed=false` e `public_path=""`.

## 2026-06-09 — Gate SEO/crawl de pré-publicação bloqueada

Decisao: uma rota candidata pode ter title, meta description, canonical e path planejados antes de publicar, mas esse planejamento deve ficar em gate bloqueado e nao pode renderizar HTML publico.

Motivos:
- Googlebot deve ver somente paginas realmente aprovadas;
- SEO tecnico deve ser validado antes de render publico, nao depois;
- fonte resolvida nao remove sozinha bloqueio editorial, CTA e revisao juridica;
- a escala de 10 mil paginas exige paths finitos e canonicals planejados sem criar URL prematura.

Consequencia: `data/editorial/prepublication_gates.jsonl` registra a primeira rota candidata para `negativa-cobertura-plano-saude`; `internal/prepublication` e `./tools/check-prepublication-gates` exigem fonte resolvida, blocker ativo, path limpo, canonical HTTPS, `noindex,follow`, title/meta dentro do orcamento, `render_allowed=false`, `sitemap_allowed=false`, `publication_allowed=false` e `public_path=""`.

## 2026-06-09 — Revisão jurídico-editorial bloqueada com CTA contextual

Decisao: CTA WhatsApp e parte critica do produto, mas deve passar por revisao juridico-editorial antes de ficar visivel.

Motivos:
- paginas informativas podem ter alta intencao de contratacao sem prometer resultado;
- CTA agressivo sem fonte e revisao vira risco juridico e conteudo ruim para humanos;
- o fluxo digital precisa pedir documentos relevantes sem induzir expectativa falsa;
- publicar CTA antes dos gates poderia contaminar a pagina e prejudicar confianca.

Consequencia: `data/editorial/legal_reviews.jsonl` registra a primeira revisao bloqueada; `internal/legalreviews` e `./tools/check-legal-editorial-reviews` exigem rascunho autoral, fonte resolvida, gate de pre-publicacao, notas juridicas, correcoes pendentes, CTA WhatsApp sem promessa, mensagem contextual com origem da pagina/rota candidata, `render_allowed=false`, `sitemap_allowed=false`, `publication_allowed=false` e `public_path=""`. `content/cta_policy.json` tambem exige template global com `{path}`, `{unique_intent_id}` e `{title}`.

## 2026-06-09 — Mudança para fábrica massiva de conteúdo único

Decisao: o projeto nao deve depender de revisão humana página a página nem limitar a produção a uma página ou poucas dezenas de rascunhos. A estratégia passa a ser geração massiva por lotes, com validação automática agressiva, score humano/IA-like, reescrita automática e bloqueio de lote quando houver spam, template ou baixa utilidade.

Motivos:
- a meta real é portal jurídico massivo, com potencial para milhões de páginas;
- revisão manual repetitiva não escala;
- conteúdo único e natural precisa ser propriedade do algoritmo, não exceção artesanal;
- Googlebot deve encontrar páginas informativas, úteis e leves, não templates;
- CTA WhatsApp deve ser contextual e lucrativo para advogado, sem transformar o texto em anúncio.

Consequencia: próximos ciclos devem implementar `human_content_score` e `scalable_content_batches`. Cada lote deve gerar muitas intenções únicas de contratação jurídica 100% digital, validar fonte, CTA contextual, score de naturalidade, similaridade intra-lote e bloqueio de publicação. Itens abaixo do score devem ser reescritos automaticamente e revalidados, não entregues ao usuário para correção manual.

## 2026-06-09 — Engenharia agressiva inteligente e autocrítica pré-commit

Decisao: Codex deve operar com engenharia agressiva inteligente: planejar a hipótese, executar sem passividade, validar em massa quando o escopo for massa, refinar algoritmo e registrar autocrítica antes de commitar.

Motivos:
- o contrato do projeto ja define plataforma juridica massiva, alta intencao e contratação digital;
- perguntas desnecessarias e decisões medrosas atrasam o P0;
- commit por ciclo é rastreabilidade, não conclusão do `/goal`;
- massa sem validação em massa vira spam, e validação tímida não prova milhões de páginas;
- o agente precisa apontar o que pode melhorar e o próximo ciclo antes de preservar o checkpoint.

Consequencia: antes de cada commit, `CHECKPOINT.md` deve registrar o que foi resolvido, provas, autocrítica, pendências reais, melhorias possíveis, próximo ciclo e a frase operacional de que o commit não encerra o `/goal`. O laboratório pode usar CPU agressivamente para testes, score, reescrita, auditoria e validação em massa; o baixo consumo de CPU continua obrigatório no runtime público/produção.

## 2026-06-09 — Score humano e lotes massivos bloqueados

Decisao: implementar `human_content_score` e `scalable_content_batches` como primeiras camadas executáveis da fábrica massiva de conteúdo jurídico único, mantendo tudo bloqueado para render, sitemap, indexação e publicação.

Motivos:
- o projeto precisa planejar milhões de páginas sem criar spam nem páginas públicas prematuras;
- validação massiva precisa ser artefato executável, não promessa em contrato;
- score humano/IA-like deve orientar reescrita automática e bloqueio de lote;
- CTA WhatsApp contextual precisa nascer no lote com origem e documentos esperados;
- produção continua leve, enquanto laboratório pode usar CPU para score e validação em massa.

Consequencia: `data/editorial/scalable_content_batches.jsonl` registra 1.020.000 páginas planejadas em seis famílias jurídicas digitais, todas bloqueadas. `data/editorial/human_content_scores.jsonl` registra scores humanos/naturalidade sem publicar conteúdo. `internal/humanscore`, `internal/scalablebatches`, `./tools/check-human-content-score` e `./tools/check-scalable-content-batches` entram no laboratório. O próximo ciclo deve gerar drafts em lote e testar reescrita automática de falhas, sem exposição pública.

## 2026-06-09 — Batch drafts com score e reescrita bloqueada

Decisao: implementar `batch_drafts` como camada de rascunhos de amostra por lote massivo, com texto editorial próprio, score humano calculado, prova de reescrita automática em falhas iniciais, baixa similaridade e publicação bloqueada.

Motivos:
- manifestos de milhão de páginas precisam virar amostras editoriais verificáveis antes de qualquer página pública;
- validar em massa sem amostra textual ainda não prova naturalidade, especificidade ou CTA contextual;
- reescrita automática precisa deixar evidência de falha inicial e correção;
- a escala deve evoluir por algoritmo, não por revisão manual página a página;
- Googlebot não deve ver rascunho enquanto score, fonte, revisão, SEO e publicação não estiverem completos.

Consequencia: `data/editorial/batch_drafts.jsonl` registra 18 rascunhos, três por família jurídica de lote, todos `batch_draft_scored_blocked`. `internal/batchdrafts` valida score via `internal/humanscore`, similaridade máxima, reescritas, origem de lote e bloqueio de render/sitemap/publicação. `./tools/check-batch-drafts` entra no laboratório. O próximo ciclo deve transformar amostras em geração programática ampliada e medição agregada por lote.
