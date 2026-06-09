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

Agentes auxiliares podem acelerar pesquisa de fontes oficiais e rascunhos bloqueados, mas a saída deles é referência. Antes de virar prova de conteúdo, fonte ou gate, deve estar registrada em `.agents/agent_context_ledger.jsonl`, com evidência, riscos e decisão de integração validada pelo Codex principal.

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

`./tools/check-batch-drafts` valida amostras de drafts em lote: mínimo por família, score humano aplicado, reescrita automática quando houve falha inicial, baixa similaridade e bloqueio total de publicação.

`./tools/check-batch-draft-generation` valida que o gerador/refinador em lote produz amostras determinísticas, com score calculado, reescrita automática, CTA WhatsApp com origem e métricas agregadas persistidas. `./tools/generate-batch-drafts` deve escrever amostras temporárias em `/tmp` durante o laboratório; quando a massa passa nos gates e contém informação jurídica útil, ela deve ser arquivada no repo como `batch_draft_expansion_archive`, sem render, sitemap, publicação ou `public_path`.

`./tools/check-batch-draft-expansion-archive` valida o arquivo permanente bloqueado de rascunhos massivos. O gate exige pelo menos 600 registros, mínimo de 100 por família, score de rascunho válido, reescrita automática comprovada, `source_matrix_id`, baixa similaridade e flags públicas falsas. Arquivo validado não é descartável; ele só pode ser removido ou rebaixado com prova registrada e novo gate.

`./tools/check-batch-candidate-expansion-readiness` valida `data/editorial/batch_candidate_expansion_readiness.jsonl`. O gate atravessa os 600 rascunhos permanentes e registra alvos de expansão de pelo menos 30 candidatos por família, mas mantém tudo bloqueado quando falta `batch_paid_intent_gates` por intenção, quando o paid intent existe mas foi reprovado, quando há fonte específica pendente ou qualquer flag pública. Readiness não publica: exige `index_policy=noindex`, `manifest_allowed=false`, `render_allowed=false`, `sitemap_allowed=false`, `publication_allowed=false`, `public_path=""` e blockers acionáveis. Status ausente (`batch_candidate_expansion_blocked_paid_gate_missing`) e status reprovado (`batch_candidate_expansion_blocked_paid_gate_failed`) são diagnósticos diferentes e não devem ser mascarados. Quando uma família não tem paid blocker nem source blocker antigo, o blocker acionável correto é `batch_candidate_gate_pending`, para orientar a próxima seleção sem fingir publicação pronta.

`./tools/check-batch-candidate-gates` valida a seleção bloqueada de candidatos a partir do arquivo permanente. O gate exige 6 famílias, pelo menos 3 intenções selecionadas por família, existência de cada intenção no arquivo, 100 registros mínimos por lote, CTA contextual, fonte matricial, base URL oficial configurada e `render_allowed=false`, `sitemap_allowed=false`, `publication_allowed=false`, `public_path=""`.

`./tools/check-batch-candidate-reviews` valida `data/editorial/batch_candidate_reviews.jsonl`. Cada intenção selecionada deve ter revisão jurídico-editorial bloqueada, fonte matricial auditada, CTA WhatsApp contextual com `Origem`, `Gate` e `Intent`, notas úteis e correções exigidas. Essa camada não publica: reprova domínio dentro de path candidato, promessa de resultado, CTA raso, fonte sem auditoria, `render_allowed=true`, `sitemap_allowed=true`, `publication_allowed=true` e qualquer `public_path`.

`./tools/check-batch-prepublication-gates` valida `data/editorial/batch_prepublication_gates.jsonl`. Cada candidato revisado deve ter canonical oficial em `https://wikijuridica.com.br`, `candidate_robots=noindex,follow`, title/meta dentro do orçamento interno e fonte final marcada como pendente. Essa camada não renderiza, não entra em sitemap e não publica.

`./tools/check-batch-source-specificity` valida `data/editorial/batch_source_specificity_resolutions.jsonl`. Cada candidato de pré-publicação precisa ter resolução de fonte: `final_source_locked_reference_only` quando a URL oficial auditada é específica o bastante para referência final, ou `final_source_blocked_needs_specific_url` quando a fonte ainda é ampla. O gate impede mascarar fonte genérica como pronta e mantém `scraping_allowed=false`, `ingestion_allowed=false`, `render_allowed=false`, `sitemap_allowed=false`, `publication_allowed=false` e `public_path=""`.

`./tools/check-batch-public-manifest-gates` valida `data/editorial/batch_public_manifest_gates.jsonl`. O manifesto bloqueado cobre todos os candidatos com resolução de fonte: candidatos com `final_source_locked_reference_only` podem ficar em `public_manifest_blocked_seo_review_pending`, enquanto candidatos com fonte ampla ficam em `public_manifest_blocked_source_specificity`. O gate reprova manifesto permitido, render, sitemap, publicação, `public_path`, robots indexável e tentativa de revisão SEO antes de fonte travada.

`./tools/check-batch-final-authorial-drafts` valida `data/editorial/batch_final_authorial_drafts.jsonl`. Apenas candidatos com `public_manifest_blocked_seo_review_pending` podem receber rascunho autoral final bloqueado. O gate exige fonte travada, PT-BR natural, score humano alto, CTA contextual com origem/intenção/documentos, aviso informativo, canonical oficial e flags públicas falsas.

`./tools/check-paid-intent` valida `data/editorial/batch_paid_intent_gates.jsonl`, a camada bloqueada de intenção comercial paga. O algoritmo deve reprovar gratuidade, não pagamento, defensoria/justiça gratuita, curiosidade, estudo acadêmico e modelo pronto; deve exigir sinal de contratação particular online, honorários, orçamento, consulta/triagem paga e contexto jurídico-econômico do caso, como valor envolvido, documentos, urgência, contrato, banco, INSS ou plano de saúde. O gate separa `paid_signals` do corpo e `cta_paid_signals` do WhatsApp; sinal pago apenas no CTA deve bloquear como `paid_intent_blocked_cta_only_paid_signal`. BPC/LOAS, CadÚnico, renda familiar, baixa renda e cumprimento de exigência/autoatendimento administrativo devem ser roteados para bloqueio comercial quando o risco de gratuidade ou self-service for dominante; `vulnerabilidade` isolada não basta para inferir assistência pública. Isso é filtro de negócio por intenção textual, não perfil pessoal.

`./tools/check-paid-intent-refinements` valida `data/editorial/batch_paid_intent_refinements.jsonl`. O refinamento em lote só pode atuar sobre `paid_intent_blocked_missing_paid_signal` e `paid_intent_blocked_cta_only_paid_signal`; ele deve mover sinal de contratação paga para o corpo informativo quando a frase for natural, recalcular score humano/IA-like e manter `noindex`, render, sitemap, publicação e `public_path` bloqueados. Temas de gratuidade, assistência pública dominante e autoatendimento administrativo não podem ser resgatados pelo refinador. `./tools/refine-paid-intent-drafts` deve ser idempotente: depois do ciclo aplicado, no-op com `refinements=0` é sucesso operacional, não falha.

`./tools/check-batch-source-matrix` valida a matriz de fonte oficial por subtema. Cada subtema precisa ter pelo menos duas URLs oficiais, tipos de fonte, score de especificidade, revisão de robots exigida, política `reference_only_no_scraping` e publicação bloqueada. O gerador só pode escalar amostras quando cada draft tiver `source_matrix_id` coberto por essa matriz.

`./tools/check-batch-source-url-audits` valida a auditoria URL-level da matriz. Cada URL oficial usada por subtema precisa ter hash, robots/termos revisados, status bloqueado, política `reference_only_no_scraping_no_ingestion`, `scraping_allowed=false`, `ingestion_allowed=false`, `publication_allowed=false` e vínculo com todos os `matrix_id` que a usam.

Similaridade e score devem operar por semântica de intenção, não por troca mecânica de palavras. O algoritmo deve separar sinais fortes de tema, subtema, faceta, documento, risco e fonte de termos operacionais de laboratório. Para escala em centenas por família, subtema e faceta semântica devem pesar no comparador; se uma amostra falhar por repetição, a correção correta é enriquecer faceta e contexto do subtema, não baixar o limite.
