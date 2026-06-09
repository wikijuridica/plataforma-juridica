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
