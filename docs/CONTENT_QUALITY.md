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

## Estratégia massiva sem spam

O projeto não deve operar como redação manual página a página. A escala correta é geração em lote com validação automática agressiva, reescrita automática e bloqueio de lote inteiro quando houver padrão mecânico. O objetivo é preparar milhares, centenas de milhares e milhões de páginas possíveis, mas cada URL indexável precisa ter intenção única, fonte, utilidade concreta, texto próprio e CTA contextual quando cabível.

Conteúdo massivo exige engenharia agressiva inteligente: planejar a família jurídica, gerar muitas intenções únicas, validar em massa, refinar o algoritmo e repetir. Não é aceitável liberar lote por confiança em um script isolado, por revisão manual lenta ou por texto que apenas parece diferente.

Regras de lote:
- criar famílias de alta intenção de contratação jurídica 100% digital;
- derivar intenções únicas por problema, documento, fonte, risco, etapa e cenário, não por permutação de palavra-chave;
- gerar conteúdo autoral com seções e exemplos diferentes por tema;
- evitar linguagem comercial agressiva no corpo: página é informativa, com intenção de contratar e CTA contextual separado;
- pontuar naturalidade, risco IA-like, spam, repetição estrutural, similaridade entre itens do lote e densidade de keyword;
- reescrever automaticamente itens abaixo do score mínimo e validar novamente;
- bloquear o lote inteiro se a amostra ou qualquer métrica agregada indicar template, thin content, fonte fraca, CTA sem origem ou páginas parecidas.

Revisão humana não pode ser gargalo para milhões de páginas. O papel do laboratório é transformar validação e revisão em algoritmo: score, motivos de reprovação, reescrita, amostragem e auditoria. Se o algoritmo estiver fraco, o agente deve melhorá-lo e continuar, não transferir correção repetitiva ao usuário.

Validação massiva é obrigatória antes de publicação massiva. O lote precisa provar diversidade real em agregados e amostras: intenção, fonte, estrutura, abertura, exemplos, documentos, CTA contextual, densidade de termos, similaridade e utilidade. Falha de lote exige reescrita/refinamento e reteste, não publicação parcial por conveniência.

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

Mudanca de estrategia: o projeto nao deve gastar ciclo tentando "descobrir termo" por script fraco. A base principal passa a ser pesquisa editorial manual na web, registrada em `data/research/high_intent_terms.jsonl` e validada por `./tools/check-manual-keyword-research`. Google Trends orienta demanda; fontes oficiais dao autoridade; a decisao editorial escolhe termos com potencial real de contratacao digital.

`./tools/check-content-briefs` valida o inicio dos conteudos em `data/editorial/content_briefs.jsonl`. Brief nao e pagina publica: ele deve ter angulo unico, problema do leitor, razao de CTA digital, fontes oficiais e secoes especificas. Titulo ou secao generica como "O que e", "Documentos necessarios" ou "Quando procurar advogado" deve reprovar para evitar template detectavel por humanos e Googlebot.

`./tools/check-authorial-content-drafts` valida rascunhos autorais em `data/editorial/authorial_drafts.jsonl`. Rascunho autoral ainda nao e pagina publica: precisa estar em PT-BR, ligado a um brief, com fonte oficial, abertura natural, secoes especificas, CTA WhatsApp contextual e `publication_allowed=false`. O validador deve reprovar abertura repetida, formato de secoes reaproveitado, heading generico, CTA raso e qualquer `public_path`.

`./tools/promote-term-candidates` promove um lote pequeno de candidatos para `term_seeds` apenas como `draft_only`. `./tools/check-promoted-term-seeds` exige `candidate_id`, evidencia de demanda, modo `digital_only`, CTA WhatsApp alto e publicacao bloqueada. O algoritmo de ranking deve ser refinado: nao confiar em um unico sinal, penalizar fonte nao oficial, penalizar fluxo presencial e manter diversidade de areas para evitar tunel de um unico nicho.

Rascunho editorial tambem deve preservar grafia natural em PT-BR. Mesmo em laboratorio, termos como `divorcio online`, `auxilio doenca negado` e `negativa cobertura plano saude` devem virar texto natural com acentuacao e conectivos corretos antes de qualquer persistencia editorial. O script `./tools/refresh-editorial-drafts` atualiza drafts persistidos quando o algoritmo de escrita e refinado.

`./tools/check-source-specificity-blockers` exige que termos priorizados tenham manifesto de fonte especifica antes de qualquer aprovacao. Fonte institucional ampla pode iniciar laboratorio, mas nao basta para conteudo publico quando falta norma, artigo, regra administrativa, requisito ou limite juridico especifico. Enquanto o bloqueio existir, `approval_allowed=false`, `publication_allowed=false` e `public_path=""`.

`./tools/check-source-specificity-resolutions` valida resolucoes de fonte especifica em `data/editorial/source_resolutions.jsonl`. Resolucao de fonte permite revisar rascunho autoral com mais seguranca, mas nao publica: exige lei primaria quando cabivel, regra de cobertura, fontes oficiais especificas, `index_policy=noindex`, `publication_allowed=false` e `public_path=""`.

`./tools/check-prepublication-gates` valida `data/editorial/prepublication_gates.jsonl`. Esse gate pode registrar title, meta description, canonical e path candidato, mas deve manter `render_allowed=false`, `sitemap_allowed=false`, `publication_allowed=false`, `public_path=""` e `candidate_robots=noindex,follow` ate todos os gates passarem.

`./tools/check-legal-editorial-reviews` valida `data/editorial/legal_reviews.jsonl`. A revisao pode preparar CTA WhatsApp contextual, mas deve reprovar promessa de liminar, prazo, causa ganha ou resultado garantido. O CTA fica `draft_contextual_not_public` ate revisao final, sem render, sitemap ou publicacao.

CTA WhatsApp sem mensagem contextual de origem deve reprovar. A mensagem precisa identificar a pagina/rota candidata, o termo ou intencao unica e os documentos esperados, para nao perder contexto comercial e juridico na triagem.

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

## Score humano e IA-like

O gate `human_content_score` é obrigatório: pontuação de naturalidade e risco IA-like/mecânico para texto jurídico em massa. Esse score deve ser usado para revisar e reescrever automaticamente, não para criar aparência artificial. Métricas mínimas:
- diversidade lexical e de frases;
- variação de abertura, headings e conclusão;
- presença de detalhes concretos do problema jurídico;
- documentos e próximos passos digitais plausíveis;
- fonte/proveniência conectada ao tema;
- ausência de promessas, superlativos e linguagem de venda;
- baixa similaridade com outros textos do mesmo lote;
- CTA contextual com origem, termo e documentos esperados.

Publicação só pode avançar quando o texto atingir score alto, motivos de reprovação estiverem zerados e o lote provar diversidade real.

`./tools/check-human-content-score` valida os registros de score. `./tools/check-scalable-content-batches` valida lotes massivos planejados, bloqueados para publicação e com escala real de produção, sem criar HTML ou sitemap.
