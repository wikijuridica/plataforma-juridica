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
- escrever resumo e comentário em linguagem humana;
- separar texto oficial, explicação informativa e opinião;
- validar com gates de qualidade, SEO e duplicidade;
- revisar antes de liberar indexação.
