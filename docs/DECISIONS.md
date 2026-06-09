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
