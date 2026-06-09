# LAB_VALIDATION.md

## Regra

O projeto deve operar como laboratorio: testar, validar, refinar, testar novamente e somente entao registrar checkpoint. Nada de chute.

## Comando principal

`./tools/lab-cycle`

Esse comando combina:
- `go test -count=1 ./...`;
- `./tools/check-all`;
- `./tools/check-sources`;
- `./tools/check-storage-contract`;
- `./tools/check-term-seeds`;
- `./tools/check-editorial-drafts`;
- `./tools/check-review-queue`;
- `./tools/check-approvals`;
- `./tools/check-publication-blockers`;
- `./tools/lab-term-draft`;
- `./tools/check-google-search-appearance`;
- `./tools/check-mechanical-content`;
- `./tools/check-cpu-budget`;
- `./tools/lab-content-quality`;
- `go run ./cmd/build public`;
- `go list -m all`;
- `git diff --check`;
- busca por residuos `.py` e `.pyc`.

O gate `./tools/check-performance-budget` deve reprovar HTML publico pesado, `<script>`, runtime cliente, bundle JavaScript, WebAssembly, mapas, `modulepreload`, import map, marcadores de hidratacao e CSS inline excessivo. Leveza e parte da prova de indexacao para Googlebot e bots valiosos.

`./tools/lab-content-quality` usa arquivos temporarios em `/tmp` para validar texto natural versus texto mecanico. Isso e laboratorio, nao publicacao. Ele deve detectar conteudo raso, keyword stuffing e permutacao antes que qualquer pagina seja exposta ao Googlebot.

`./tools/check-cpu-budget` vale para runtime publico/producao. Ele nao proibe testes, build, laboratorio ou auditorias mais pesadas quando forem necessarias para provar qualidade; ele impede que caminhos de atendimento publico gastem CPU com execucao externa, rede, sleeps ou loops sem limite.

`./tools/lab-term-draft` usa seeds juridicas aprovadas pelo contrato de `term_seeds` para criar rascunho temporario em `/tmp`. O rascunho deve ser `draft/noindex`, nao pode escrever em `content/pages.json`, nao pode entrar em sitemap e nao pode receber CTA.

`./tools/persist-term-drafts` e ferramenta de ciclo para gravar rascunho validado na camada `editorial_drafts`; `./tools/check-editorial-drafts` e o gate permanente no laboratorio para garantir que o draft persistido continua fora da publicacao.

`./tools/queue-editorial-review` move rascunhos persistidos para fila de revisao `needs_review`, sem publicacao. `./tools/check-review-queue` exige `publication_allowed=false`, `public_path` vazio, fonte, autoria e historico.

`./tools/approve-editorial-review` registra aprovacao editorial de laboratorio; `./tools/check-approvals` garante que essa aprovacao nao cria rota publica nem libera indexacao.

`./tools/block-publication` cria manifesto de publicacao bloqueada; `./tools/check-publication-blockers` garante que aprovacao editorial ainda nao libera URL publica.

Quando o laboratorio passar com seguranca, o ciclo seguinte deve migrar para host/produção controlada e priorizar termos juridicos pesquisados por humanos e com alta intencao de contratacao online, sem spam.

Alta intencao de contratacao online exige que o fluxo possa ser 100% digital: WhatsApp, envio remoto de documentos, triagem remota e contratacao de advogado sem depender de atendimento presencial como caminho principal.

Os caminhos seguros de pesquisa de demanda ficam em `docs/data-sources/search-demand.md`. O laboratorio deve validar `term_intent_candidates` antes de transformar qualquer termo pesquisado em seed, rascunho ou pauta publica.

Para termos de alta intencao, a estrategia preferida e pesquisa editorial manual na web, nao script de descoberta. O laboratorio deve validar `manual-keyword-research`, `content-briefs` e `authorial-content-drafts`, revisando se os briefs e rascunhos possuem angulo unico, nao repetem molde e nao tentam publicar. Google Trends orienta demanda, mas nao decide sozinho.

Promocao de candidato para seed deve passar por ranking refinavel, nao por lista fixa. O score deve considerar area, risco, fonte oficial, adequacao digital, CTA WhatsApp e penalidades para fluxo presencial ou evidencia fraca. Quando falso positivo ou falso negativo aparecer, escrever teste e ajustar o algoritmo antes de continuar.

Rascunho autoral deve passar por algoritmo anti-template antes de qualquer novo lote: abertura repetida, conjunto de secoes reaproveitado, heading generico, CTA raso, fonte ausente ou path publico devem reprovar. Contagem de palavras isolada nao basta para liberar conteudo.

Se a inspeção mostrar grafia mecanica ou sem acento em rascunho PT-BR, tratar como falha de laboratorio. Refinar `internal/draftlab`, regenerar `data/editorial/drafts.jsonl`, atualizar `data/editorial/review_queue.jsonl` e rodar novamente os checks antes do commit.

Antes de aprovar qualquer draft priorizado, validar `source_specificity_blockers`. O laboratorio deve bloquear termo que usa apenas fonte ampla, noticia institucional ou canal administrativo quando ainda falta norma, artigo, regra ou recorte juridico especifico.

Quando uma fonte especifica for encontrada, registrar `source_specificity_resolutions` e rodar `./tools/check-source-specificity-resolutions`. A resolucao aproxima o rascunho da revisao, mas permanece `noindex`, sem URL publica e sem CTA publico ate os demais gates.

Antes de renderizar uma rota candidata, registrar `prepublication_gates` e rodar `./tools/check-prepublication-gates`. O gate deve provar title/meta/canonical dentro do orcamento do Google Search, mas manter render, sitemap e publicacao bloqueados.

## Politica

Nenhum script isolado e prova suficiente para mudanca P0/P1. Use o laboratorio como conjunto minimo e inspecione os artefatos quando a mudanca afetar HTML, sitemap, robots, canonical, indexacao, qualidade ou conteudo.

Validar nao basta. Todo ciclo deve revisar o diff, os artefatos gerados, os riscos e o contrato antes de commitar. A revisao deve procurar atalhos, spam, conteudo mecanico, regressao de P0 e dependencia indevida.

Se qualquer validacao falhar:
1. identificar causa raiz;
2. corrigir;
3. rodar novamente;
4. registrar falha e correcao no checkpoint;
5. commitar checkpoint e artefatos do ciclo;
6. continuar o proximo ciclo quando nao houver bloqueio P0 real.

## Continuidade

O agente nao deve parar ao registrar checkpoint. Todo ciclo deve terminar com proximo passo acionavel, e esse proximo passo deve ser executado em seguida quando nao houver bloqueio P0 real comprovado.

## Commit por ciclo

Depois de validacoes relevantes passarem, fazer commit do checkpoint e dos artefatos do ciclo. O projeto nao deve depender de chat, contexto compactado ou memoria externa para continuidade.

Antes do commit, executar `date`, conferir a hora local e registrar o ciclo numerado no checkpoint. Isso mantem a ordem de auditoria mesmo apos compactacao de contexto.
