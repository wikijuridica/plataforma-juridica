# CONTENT_QUALITY.md

Toda página indexável passa por `internal/quality`.

Gates implementados no ciclo P0/P1:
- URL limpa, minúscula, sem parâmetros e com barra final quando aplicável;
- canonical HTTPS absoluto apontando para a própria rota;
- `unique_intent_id`, título e meta description obrigatórios;
- título, meta description, canonical e intenção sem duplicidade entre páginas indexáveis;
- hash normalizado de conteúdo;
- similaridade por shingles;
- mínimo textual para páginas indexáveis;
- mínimo de links internos úteis;
- conteúdo jurídico indexável com fonte, revisão e aviso informativo;
- motivo de publicação.

Conteúdo que falha permanece `draft`, `needs_review`, `noindex` ou `archived`, e não entra em sitemap.

CTA comercial por WhatsApp é subordinado a este gate. Página sem fonte, sem revisão, sem intenção única ou sem valor informativo não pode usar CTA como justificativa para indexação.

## Escrita natural

Conteúdo visível ao público deve ser escrito em PT-BR, com grafia correta, acentuação correta, pontuação clara e linguagem natural. Rascunho técnico interno pode ficar sem polimento, mas texto público não.

Antes de criar conteúdo jurídico, pesquisar a fonte correta e documentar a proveniência. A redação deve ser natural, clara e útil para humanos. É proibido publicar texto mecânico, permutação de termos, páginas quase iguais ou conteúdo criado apenas para atrair busca.

Fonte oficial não é alvo para scraping, clone ou espelho. A fonte é referência/proveniência para produzir conteúdo próprio, natural e único. O texto não deve copiar estrutura, massa de dados ou conteúdo oficial de forma mecânica.

Roteiro mínimo antes de qualquer página jurídica:
- identificar a fonte oficial correta;
- documentar URL, data, limites, riscos e estratégia de proveniência;
- registrar termo jurídico, quando usado como semente, no banco leve separado e manter o estado `draft_only`;
- escrever resumo e comentário em linguagem humana;
- separar texto oficial, explicação informativa e opinião;
- validar com gates de qualidade, SEO e duplicidade;
- revisar antes de liberar indexação.

## Termos juridicos como semente

Termos juridicos podem ser usados para iniciar pauta e rascunho, mas nao podem virar pagina automaticamente. A camada `term_seeds` exige proveniencia e estado de qualidade; a camada `editorial_drafts` guarda texto proprio em PT-BR; e o manifesto publicado registra apenas conteudo que ja passou por fonte, revisao, qualidade, SEO e decisao editorial.

`./tools/check-term-seeds` valida que cada semente esta em PT-BR, `draft_only`, tem fonte, URL oficial, data de verificacao e indicacao de intencao editorial. Seed invalida bloqueia o laboratorio antes de qualquer rascunho.

`./tools/check-term-intent-candidates` valida candidatos de termos pesquisados por humanos antes de qualquer conteudo publico. A evidencia segura inicial e Google Trends no Brasil (`https://trends.google.com.br/trends/explore?geo=BR`), interpretado conforme a Ajuda do Google Trends e a Central da Pesquisa Google. Essa evidencia e direcional: nao e volume absoluto, nao e fonte juridica e nao autoriza publicar sem fonte oficial, revisao e qualidade.

`./tools/promote-term-candidates` promove um lote pequeno de candidatos para `term_seeds` apenas como `draft_only`. `./tools/check-promoted-term-seeds` exige `candidate_id`, evidencia de demanda, modo `digital_only`, CTA WhatsApp alto e publicacao bloqueada. O algoritmo de ranking deve ser refinado: nao confiar em um unico sinal, penalizar fonte nao oficial, penalizar fluxo presencial e manter diversidade de areas para evitar tunel de um unico nicho.

Rascunho editorial tambem deve preservar grafia natural em PT-BR. Mesmo em laboratorio, termos como `divorcio online`, `auxilio doenca negado` e `negativa cobertura plano saude` devem virar texto natural com acentuacao e conectivos corretos antes de qualquer persistencia editorial. O script `./tools/refresh-editorial-drafts` atualiza drafts persistidos quando o algoritmo de escrita e refinado.

`./tools/lab-term-draft` transforma seed valida em rascunho temporario dentro de `/tmp`, com estado `draft/noindex`, fonte citada e aviso informativo. Esse laboratorio prova linguagem natural sem escrever em `content/pages.json` e sem expor a seed ao Googlebot.

`./tools/persist-term-drafts` persiste rascunhos validados em `data/editorial/drafts.jsonl` de forma idempotente. `./tools/check-editorial-drafts` reprova qualquer draft que tenha rota publica, indexacao, fonte ausente ou qualidade textual insuficiente.

`./tools/queue-editorial-review` coloca drafts persistidos em `data/editorial/review_queue.jsonl` com `needs_review`, historico e `publication_allowed=false`. `./tools/check-review-queue` bloqueia qualquer fila que crie URL publica, permita publicacao ou perca autoria/historico.

`./tools/approve-editorial-review` registra aprovacao editorial em `data/editorial/approved_drafts.jsonl`, mas ainda com `publication_allowed=false`, `public_path` vazio e `noindex`. Aprovacao editorial nao e publicacao.

`./tools/block-publication` cria `data/editorial/publication_blockers.jsonl` com requisitos faltantes antes de qualquer URL publica. Esse manifesto e obrigatorio para nao confundir laboratorio aprovado com conteudo pronto para Googlebot.

Quando os gates de laboratorio forem seguros, o proximo passo nao e inventar conteudo: e migrar para host/produção controlada e pesquisar termos juridicos reais, mais buscados por humanos e com alta intencao de contratar advogado online, mantendo fonte, unicidade e CTA responsavel.

Alta intencao, para este projeto, significa contratacao juridica 100% digital: a pessoa pode entender o problema, enviar documentos, conversar por WhatsApp e contratar atendimento juridico online. Termo cujo caminho normal exige comparecimento presencial, diligencia local obrigatoria ou baixa conversao digital deve perder prioridade inicial, mesmo que tenha volume de busca.

## Conteúdo mecânico

Conteúdo raso ou mecânico deve ser detectado antes de qualquer exposição ao Googlebot. O laboratório `./tools/lab-content-quality` cria textos temporários em `/tmp`: um texto natural em PT-BR e um texto mecânico de permutação de palavras-chave. O detector precisa aprovar o texto natural e reprovar o mecânico.

Sinais iniciais:
- mínimo textual para conteúdo indexável;
- baixa diversidade lexical;
- excesso de repetição de palavra relevante;
- frase repetida;
- sentença repetida;
- combinação de sinais que indique permutação de keyword.

O algoritmo deve ser melhorado sempre que gerar falso positivo ou falso negativo. Exemplo já registrado: repetição normal de marca/título/heading não deve ser confundida com conteúdo mecânico; por isso a análise de página usa o corpo editorial, não o conjunto inteiro de title/meta/heading.
