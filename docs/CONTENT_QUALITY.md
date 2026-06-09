# CONTENT_QUALITY.md

Toda pagina indexavel passa por `internal/quality`.

Gates implementados no ciclo P0/P1:
- URL limpa, minuscula, sem parametros e com barra final quando aplicavel;
- canonical HTTPS absoluto apontando para a propria rota;
- `unique_intent_id`, titulo e meta description obrigatorios;
- titulo, meta description, canonical e intencao sem duplicidade entre paginas indexaveis;
- hash normalizado de conteudo;
- similaridade por shingles;
- minimo textual para paginas indexaveis;
- minimo de links internos uteis;
- conteudo juridico indexavel com fonte, revisao e aviso informativo;
- motivo de publicacao.

Conteudo que falha permanece `draft`, `needs_review`, `noindex` ou `archived`, e nao entra em sitemap.

CTA comercial por WhatsApp e subordinado a este gate. Pagina sem fonte, sem revisao, sem intencao unica ou sem valor informativo nao pode usar CTA como justificativa para indexacao.

## Escrita natural

Antes de criar conteudo juridico, pesquisar a fonte correta e documentar a proveniencia. A redacao deve ser natural, clara e util para humanos. E proibido publicar texto mecanico, permutacao de termos, paginas quase iguais ou conteudo criado apenas para atrair busca.

Fonte oficial nao e alvo para scraping, clone ou espelho. A fonte e referencia/proveniencia para produzir conteudo proprio, natural e unico. O texto nao deve copiar estrutura, massa de dados ou conteudo oficial de forma mecanica.

Roteiro minimo antes de qualquer pagina juridica:
- identificar a fonte oficial correta;
- documentar URL, data, limites, riscos e estrategia de proveniencia;
- escrever resumo e comentario em linguagem humana;
- separar texto oficial, explicacao informativa e opiniao;
- validar com gates de qualidade, SEO e duplicidade;
- revisar antes de liberar indexacao.
