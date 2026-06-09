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

## 2026-06-09 — Proibido finalizar goal por ferramenta em P0

Decisao: o agente nao pode chamar `update_goal` com `status=complete` enquanto o projeto estiver em P0/laboratório ou antes de 10 mil paginas publicas juridicas aprovadas e verificadas.

Motivos:
- o goal ativo e a meta publica minima sao maiores que qualquer checkpoint de laboratorio;
- resposta final no thread, commit, check verde ou arquivo persistido nao provam o objetivo real;
- marcar completion por ferramenta encerraria a continuidade operacional contra o contrato.

Consequencia: a primeira linha de `AGENTS.md` e `GOAL.md` explicita a proibicao operacional. `internal/contract/continuity_test.go` reprova se a regra sumir dos contratos. O agente deve deixar o goal ativo, registrar o proximo ciclo e continuar trabalhando enquanto nao houver bloqueio P0 real comprovado.

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

## 2026-06-09 — Gerador/refinador de batch drafts com métricas agregadas

Decisao: implementar `batchdraftgen` como gerador/refinador determinístico de rascunhos de lote, produzindo amostras temporárias e persistindo métricas agregadas bloqueadas.

Motivos:
- a fábrica massiva precisa gerar e validar lote por comando, não depender de amostras escritas uma a uma;
- reescrita automática deve ser comprovada por falha inicial, score final e métrica agregada;
- métricas por família permitem crescer volume sem mascarar similaridade, IA-like ou baixa especificidade;
- o laboratório pode usar CPU para gerar/refinar, mas nada deve escapar para HTML, sitemap, `public_path` ou publicação;
- CTA WhatsApp precisa nascer com origem de `unique_intent_id` para triagem digital.

Consequencia: `internal/batchdraftgen`, `cmd/generate-batch-drafts`, `./tools/generate-batch-drafts` e `./tools/check-batch-draft-generation` entram no laboratório. `data/editorial/batch_generation_metrics.jsonl` registra 6 métricas de geração bloqueada; o comando gera 30 rascunhos temporários, cinco por família, todos reescritos e com similaridade máxima 0.27 no laboratório. O próximo ciclo deve aumentar o volume por família, cruzar fonte específica por subtema e preparar gate de pré-publicação bloqueada por lote sem publicar.

## 2026-06-09 — Escala semântica com matriz de fontes por subtema

Decisao: ampliar o gerador para 10 amostras por família e criar `batch_source_matrix` como matriz leve de fontes oficiais por subtema, mantendo uso apenas referencial e sem scraping.

Motivos:
- escala maior revelou que similaridade pode subir quando o algoritmo repete vocabulário operacional de laboratório;
- o projeto exige boa semântica por tema/subtema, não mecanização de palavras;
- cada draft de lote precisa carregar fonte oficial específica, documento, risco, ação digital e CTA de origem;
- métricas de escala precisam registrar cobertura de fonte, risco estrutural e custo estimado de laboratório;
- publicar sem matriz de fonte específica criaria risco jurídico e risco de conteúdo raso.

Consequencia: `internal/batchsourcematrix`, `data/editorial/batch_source_matrix.jsonl` e `./tools/check-batch-source-matrix` entram no laboratório. `batchdraftgen` passa a gerar 60 drafts temporários com `source_matrix_id`, cobertura de matriz, risco estrutural e estimativa de CPU de laboratório; a similaridade máxima validada caiu para 0.52 após refinamento semântico. O próximo ciclo deve ampliar a matriz e o gerador para centenas de amostras por família, mantendo fonte específica e publicação bloqueada.

## 2026-06-09 — Auditoria URL-level e centenas de rascunhos por família

Decisao: criar `batch_source_url_audits` e ampliar o gerador para validar 100 amostras por família no laboratório, mantendo publicação bloqueada e sem afrouxar score, similaridade ou fonte.

Motivos:
- matriz de fonte por subtema ainda nao prova auditoria de cada URL oficial usada pelo lote;
- centenas de amostras por família exigem algoritmo semântico, não repetição de três sufixos;
- similaridade precisa diferenciar faceta e subtema, sem confundir metadado bruto com texto editorial;
- testes podem usar CPU no laboratório, mas o runtime público continua leve;
- nenhum rascunho de lote pode virar render, sitemap, `public_path` ou página indexável no P0.

Consequencia: `internal/batchsourceaudit`, `data/source-audit/batch_source_urls.jsonl` e `./tools/check-batch-source-url-audits` entram no laboratório. `batchdraftgen` passa a gerar 600 drafts temporários em teste de contrato, com 100 por família, facetas semânticas distribuídas por subtema, contexto de área, auditoria URL-level e similaridade máxima abaixo do limite de 0.64. `batchdrafts.MaximumPairSimilarity` foi otimizado para pré-computar sinais semânticos e pondera subtema/faceta sem afrouxar o limite. O próximo ciclo deve transformar essa massa temporária em gate de lote candidato, ainda bloqueado, com amostra persistida controlada e pré-publicação sem URL pública.

## 2026-06-09 — Laboratorio aprovado vira arquivo permanente bloqueado

Decisao: rascunhos massivos gerados em `/tmp` que passam nos gates e contêm informação jurídica útil devem ser preservados no repositório como `batch_draft_expansion_archive`, não descartados.

Motivos:
- checkpoint nao pode depender de diretorio temporario ou contexto compactado;
- rascunhos validados podem virar base permanente de páginas futuras depois de expansão, fonte e revisão;
- excluir dados jurídicos úteis sem prova atrasa a fabrica de conteúdo e reduz rastreabilidade;
- persistir no repo nao significa publicar, renderizar, criar sitemap ou liberar CTA público.

Consequencia: `data/editorial/batch_draft_expansion_archive.jsonl`, `internal/batchdraftarchive` e `./tools/check-batch-draft-expansion-archive` entram no laboratório. O arquivo exige 600 rascunhos, 100 por família, `source_matrix_id`, reescrita automática, baixa similaridade e bloqueio total de render/sitemap/publicação. Remoção ou rebaixamento de rascunho validado exige prova em checkpoint e gate próprio.

## 2026-06-09 — URL oficial do projeto ainda nao esta travada

Decisao: tratar `content/site.json` como fonte configuravel da base de canonical/sitemap/robots e marcar a base atual como placeholder de laboratorio, nao URL oficial do produto.

Motivos:
- o projeto ainda esta em P0 e a URL oficial publica nao foi definida;
- testes rigidos por dominio quebrariam a migracao futura sem melhorar SEO;
- canonical continua obrigatorio, mas deve ser validado por base configurada, HTTPS e path limpo;
- pre-publicacao bloqueada pode planejar canonical candidato sem criar URL publica.

Consequencia: `content/site.json` passa a declarar `base_url_mode`, `official_url_status` e `official_url_locked`. `internal/prepublication` valida canonical candidato contra a base carregada de `content/site.json`, e o teste `TestPrepublicationGateAcceptsConfigurableProjectBaseURL` prova que outro dominio HTTPS pode ser aceito sem alterar algoritmo. `portal-juridico.example` so pode ser tratado como placeholder enquanto `base_url_mode="lab_placeholder"`.

## 2026-06-09 — Gate candidato bloqueado a partir do arquivo permanente

Decisao: criar `batch_candidate_gates` como etapa intermediaria entre arquivo permanente de rascunhos e pre-publicacao, selecionando candidatos reais sem publicar.

Motivos:
- o arquivo de 600 rascunhos precisa virar continuidade operacional, nao ficar apenas como massa bruta;
- seleção de candidato nao pode ser confundida com URL publica, sitemap ou CTA visivel;
- a URL oficial ainda nao esta travada, entao o gate precisa respeitar `base_url_mode`;
- cada candidato deve existir no arquivo permanente, carregar CTA contextual e continuar vinculado a fonte matricial.

Consequencia: `data/editorial/batch_candidate_gates.jsonl`, `internal/batchcandidategates` e `./tools/check-batch-candidate-gates` entram no laboratorio. O gate exige 6 famílias, 3 intenções selecionadas por família, 100 registros mínimos no arquivo por lote, similaridade <=0.64, score humano mínimo, base URL flexível e flags públicas falsas. O próximo ciclo deve transformar candidatos selecionados em revisão jurídico-editorial por candidato, ainda sem render público.

## 2026-06-09 — URL oficial travada em wikijuridica.com.br

Decisao: travar a URL oficial do projeto como `https://wikijuridica.com.br` e remover o placeholder `portal-juridico.example` dos canonicals e gates vivos.

Motivos:
- o usuario definiu `wikijuridica.com.br` como dominio oficial do projeto;
- canonical, sitemap e robots precisam de base real antes de evoluir pre-publicacao;
- testes continuam lendo `content/site.json` para evitar algoritmo hardcoded, mas o contrato agora exige `official_configured`;
- URL oficial travada nao e autorizacao para publicar candidatos de lote ou rascunhos.

Consequencia: `content/site.json` passa para `base_url_mode="official_configured"`, `official_url_status="locked"` e `official_url_locked=true`. `content/pages.json`, `data/editorial/prepublication_gates.jsonl` e `data/editorial/batch_candidate_gates.jsonl` acompanham a base oficial. O proximo ciclo deve manter candidatos bloqueados e preparar pre-publicacao em lote apenas depois de fonte e revisao especificas.

## 2026-06-09 — Revisao juridico-editorial bloqueada de candidatos de lote

Decisao: criar `batch_candidate_reviews` como camada obrigatoria entre `batch_candidate_gates` e qualquer pre-publicacao em lote.

Motivos:
- selecionar candidato nao basta para preparar pagina juridica publica;
- cada candidato precisa de revisao juridico-editorial, CTA WhatsApp contextual e fonte matricial auditada;
- escala massiva nao pode depender de revisao manual pagina a pagina, mas a revisao algoritmica precisa deixar rastro por candidato;
- URL oficial travada nao remove os gates de fonte, qualidade, SEO e publicacao;
- promessa de resultado, CTA raso, path com dominio e fonte nao auditada precisam reprovar antes de qualquer render.

Consequencia: `data/editorial/batch_candidate_reviews.jsonl`, `internal/batchcandidatereviews`, `./tools/check-batch-candidate-reviews`, `internal/checks` e `tools/lab-cycle` entram no laboratorio. O gate exige 18 revisoes bloqueadas, uma por intencao selecionada, com `Origem`, `Gate` e `Intent` na mensagem de WhatsApp, matriz auditada, notas especificas e flags publicas falsas. O proximo ciclo deve transformar essas revisoes em pre-publication gates de lote com canonical oficial e `noindex`, ainda sem sitemap/publicacao.

## 2026-06-09 — Validacao global proporcional ao risco

Decisao: tratar `go test -count=1 ./...`, `./tools/check-all`, `./tools/lab-cycle` e equivalentes completos como validacao global de alto custo, nao como ritual automatico para toda alteracao pequena.

Motivos:
- validacao global consome tempo e pode atrasar ciclos localizados;
- engenharia agressiva exige prova suficiente, nao excesso de ritual;
- mudancas pequenas podem ser comprovadas com teste focado, check especifico, diff check e inspecao direta;
- mudancas amplas ou criticas ainda exigem prova global para evitar regressao em varias camadas.

Consequencia: o proximo ciclo deve escolher validacao proporcional. Rodar validacao global quando houver alteracao ampla, contrato central, risco P0/P1 critico, HTML/sitemap/canonical/robots/indexacao/performance, gerador em massa, preparacao de publicacao ou falha transversal. Em ciclos localizados, registrar no checkpoint os checks focados usados e por que eles cobrem o risco.

## 2026-06-09 — Pre-publicacao bloqueada para candidatos revisados

Decisao: criar `batch_prepublication_gates` como etapa posterior a `batch_candidate_reviews`, registrando canonical oficial, `noindex,follow`, title/meta e pendencias finais sem publicar.

Motivos:
- revisao juridico-editorial de candidato ainda nao equivale a pagina publica;
- a URL oficial ja esta travada e deve aparecer no canonical candidato;
- Googlebot nao deve receber candidatos enquanto fonte final, revisao SEO, manifesto publico e render/sitemap nao forem aprovados;
- pre-publicacao em lote precisa ser validada por candidato, nao por suposicao global.

Consequencia: `data/editorial/batch_prepublication_gates.jsonl`, `internal/batchprepublication`, `./tools/check-batch-prepublication-gates`, `internal/checks` e `tools/lab-cycle` entram no laboratorio. O ciclo usa validacao focada por politica proporcional: teste do contrato novo, tool especifica, storage contract, checks internos e diff check; validacao global fica reservada para alteracao ampla ou risco critico.

## 2026-06-09 — Especificidade de fonte por candidato pre-publicado

Decisao: criar `batch_source_specificity_resolutions` como camada obrigatoria depois de `batch_prepublication_gates`, cobrindo cada candidato com fonte final travada como referencia ou bloqueio explicito por fonte ampla.

Motivos:
- matriz de fonte auditada nao basta para dizer que todo candidato esta pronto;
- fonte institucional ampla nao pode ser mascarada como fonte final;
- candidatos com URL oficial especifica podem avancar para o proximo gate bloqueado sem scraping, ingestao ou publicacao;
- candidatos com fonte ampla precisam registrar motivo, detalhe necessario e permanecer fora de render/sitemap/publicacao.

Consequencia: `data/editorial/batch_source_specificity_resolutions.jsonl`, `internal/batchsourcespecificity`, `./tools/check-batch-source-specificity`, `internal/checks`, `content/storage_contract.json` e `tools/lab-cycle` entram no laboratorio. O gate exige 18 resolucoes, uma por candidato pre-publicado, fonte URL-level auditada, politica `reference_only_no_scraping_no_ingestion`, `candidate_robots=noindex,follow`, canonical oficial e flags publicas falsas. O ciclo usa validacao proporcional focada; `check-all` e `lab-cycle` ficam reservados para alteracao ampla ou risco transversal.

## 2026-06-09 — Manifesto publico bloqueado por candidato de lote

Decisao: criar `batch_public_manifest_gates` como camada bloqueada depois de `batch_source_specificity_resolutions`, cobrindo todos os candidatos e permitindo avanço interno para SEO/conteudo final apenas quando a fonte esta travada.

Motivos:
- fonte travada ainda nao autoriza publicacao, render ou sitemap;
- candidatos com fonte ampla nao podem entrar em revisao SEO como se estivessem prontos;
- o pipeline precisa separar backlog de fonte de backlog de SEO/conteudo;
- manifesto publico real deve ser posterior e mais restrito que este gate bloqueado.

Consequencia: `data/editorial/batch_public_manifest_gates.jsonl`, `internal/batchpublicmanifest`, `./tools/check-batch-public-manifest-gates`, `internal/checks`, `content/storage_contract.json` e `tools/lab-cycle` entram no laboratorio. O gate exige 18 registros, 7 com `public_manifest_blocked_seo_review_pending` e 11 com `public_manifest_blocked_source_specificity`, mantendo `index_policy=noindex`, `manifest_allowed=false`, `render_allowed=false`, `sitemap_allowed=false`, `publication_allowed=false` e `public_path=""`.

## 2026-06-09 — Rascunho autoral final bloqueado e intenção paga

Decisao: criar `batch_final_authorial_drafts` apenas para os candidatos com fonte travada e manifesto SEO pendente, e adicionar `paid-intent` como gate de negócio para impedir funil de gratuidade, curiosidade ou baixa intenção de contratação.

Motivos:
- rascunho final não é publicação, mas precisa virar artefato permanente para escalar conteúdo;
- CTA WhatsApp deve carregar origem, intenção, documentos e sinal de contratação particular;
- serviço jurídico comercial precisa priorizar busca com honorários/orçamento, valor envolvido, urgência e documentos concretos;
- termos de gratuidade, defensoria, justiça gratuita, estudo acadêmico, modelo pronto ou curiosidade não devem alimentar o lote comercial;
- inferência aceitável é textual e jurídico-econômica do termo/caso, não perfil pessoal sensível.

Consequencia: `data/editorial/batch_final_authorial_drafts.jsonl`, `internal/batchfinaldrafts`, `internal/paidintent`, `./tools/check-batch-final-authorial-drafts`, `./tools/check-paid-intent`, `internal/checks`, `content/storage_contract.json` e `tools/lab-cycle` entram no laboratorio. Todos os 7 rascunhos seguem `noindex`, sem render, sem sitemap, sem publicação e sem `public_path`; o gate pago deve ser refinado quando surgir falso positivo/negativo antes de escalar.

## 2026-06-09 — Agentes auxiliares sem concorrência crítica

Decisao: no próximo ciclo e nos seguintes, usar agentes auxiliares quando houver trabalho independente, especialmente pesquisa de fontes oficiais, matriz de proveniência, testes, conteúdo bloqueado e debugging, mas sem concorrência no estado do repositório.

Motivos:
- a meta massiva exige acelerar pesquisa e produção sem perder validação;
- fontes oficiais e conteúdo por subtema podem ser divididos por área/fonte;
- concorrência no mesmo arquivo, gate, commit, fonte jurídica ou decisão crítica aumenta risco de conflito e mascaramento;
- o Codex principal deve manter responsabilidade por arquitetura, P0/P1, integração, validação, checkpoint e commit.

Consequencia: agentes devem produzir pesquisa, evidência ou rascunho em escopo isolado. Se houver escrita, ela precisa ser disjunta e só entra no repo após validação do Codex principal. O Codex principal valida evidência, revisa o diff, roda os checks relevantes, registra checkpoint e não publica nada sem gate. Para não perder contexto em compactação, todo agente usado em ciclo deve virar registro em `.agents/agent_context_ledger.jsonl` antes do checkpoint/commit.

## 2026-06-09 — Ledger persistente contra perda de contexto de agentes

Decisao: criar `.agents/agent_context_ledger.jsonl` e `./tools/check-agent-context-ledger` como contrato operacional para preservar contexto de subagentes entre compactações.

Motivos:
- compactação pode remover IDs, achados, riscos e decisões de integração dos agentes;
- pesquisa de fonte oficial e revisão de conteúdo precisam sobreviver ao próximo ciclo;
- agente auxiliar não pode virar prova invisível nem autorização implícita;
- o Codex principal precisa conseguir retomar sem refazer pesquisa ou confiar em memória solta.

Consequencia: cada agente usado deve registrar ciclo, ID, apelido, tipo de tarefa, escopo, status, política de uso, resumo, evidências, riscos, decisão de integração e flags `repo_write_allowed=false`, `codex_validation_required=true`, `closed_before_checkpoint=true`. `internal/agentcontext` valida o ledger; `check-all` e `lab-cycle` passam a incluir esse gate.

## 2026-06-09 — Timing e otimização obrigatórios antes de commit

Decisao: teste lento não pode ser tratado como normal sem medição. Criar `cmd/profile-tests`, `internal/testprofile`, `./tools/profile-contract-tests` e fazer `tools/lab-cycle` medir cada etapa com `TIMING`, removendo duplicação de `check-all` e checks individuais já cobertos por `cmd/check all`.

Motivos:
- `internal/contract` chegou a 207,365s em medição por `go test -json`;
- gargalos reais eram revalidação editorial upstream e similaridade de 600 drafts, não falta de vontade de usar CPU;
- `lab-cycle` repetia `go test`, `check-all` e vários checks individuais;
- otimização precisa preservar contrato, não mascarar gate.

Consequencia: antes de commit com teste/gate lento, rodar perfil ou registrar timing. O ciclo otimizado usa `go test -count=1 ./...`, `go run ./cmd/check all` e ferramentas de laboratório não cobertas pelo check-all. `batchdrafts` evita alocação de mapa de união em Jaccard e validadores finais deixam de reconstruir índice upstream duas vezes no mesmo caminho. Medição pós-otimização de `./tools/profile-contract-tests`: `total_observed_seconds=17.02`, `slow_tests=0` com threshold de 5s.

## 2026-06-09 — Gate pago persistente e fonte geral como fonte ampla

Decisao: intenção comercial paga agora fica registrada em `batch_paid_intent_gates`, e fontes legais gerais como Código Civil, CDC e CLT compilada não destravam recortes específicos sozinhas.

Motivos:
- CTA com honorários não basta quando o tema dominante indica assistência pública, gratuidade provável ou autoatendimento administrativo;
- BPC/LOAS, CadÚnico, renda familiar/vulnerabilidade e cumprimento de exigência precisam de bloqueio comercial explícito antes de escalar conteúdo;
- fonte ampla pode ser referência oficial, mas não substitui ato, súmula, regra, serviço ou orientação específica para o subtema;
- agentes auxiliares podem acelerar pesquisa oficial e rascunho bloqueado, mas o Codex principal precisa validar evidência, integração e gates críticos.

Consequencia: `./tools/check-paid-intent` valida o arquivo permanente `data/editorial/batch_paid_intent_gates.jsonl`; candidatos comerciais fortes continuam bloqueados para publicação, e candidatos com risco de assistência pública ou self-service ficam roteados para bloqueio comercial. `./tools/check-batch-source-specificity` trata `codigo_civil`, `codigo_consumidor` e `clt_compilada` como fontes amplas quando o recorte precisa de fonte específica. O próximo ciclo deve criar prontidão de expansão de candidatos, usando agentes sem concorrência crítica para pesquisar fontes oficiais e ampliar conteúdo bloqueado com validação pelo Codex principal.

## 2026-06-09 — Prontidao de expansao bloqueada por paid gate

Decisao: criar `batch_candidate_expansion_readiness` como camada entre o arquivo permanente de 600 rascunhos e a seleção candidata ampliada, com alvo inicial de 30 candidatos por família e bloqueio explícito quando faltar paid-intent por intenção.

Motivos:
- o projeto precisa sair de 18 candidatos rumo a dezenas/centenas por família sem publicar spam;
- `batch_draft_expansion_archive` já prova massa útil, mas não prova que cada intenção expandida tem paid gate, fonte específica e CTA aptos;
- CTA colado não pode salvar candidato fraco, e paid-intent ausente deve ser blocker acionável;
- readiness precisa sobreviver ao checkpoint e orientar agentes auxiliares sem concorrência no mesmo gate.

Consequencia: `data/editorial/batch_candidate_expansion_readiness.jsonl`, `internal/batchcandidateexpansion` e `./tools/check-batch-candidate-expansion-readiness` entram no laboratório. O gate registra 6 famílias com 30 alvos cada, mas mantém status `batch_candidate_expansion_blocked_paid_gate_missing` enquanto as intenções expandidas não tiverem `batch_paid_intent_gates`. Nenhum registro permite manifesto, render, sitemap, publicação ou `public_path`.
