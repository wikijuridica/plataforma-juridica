# DATA_SOURCES.md

Antes de ingestao, cada fonte oficial precisa de documento em `docs/data-sources/`.

Estado inicial:
- Portal da Legislacao / Planalto: URL oficial confirmada como `https://legislacao.presidencia.gov.br/` por pagina publica `gov.br`, documentada em `docs/data-sources/planalto.md`, sem ingestao automatica por falha de alcance HTTP local.
- Camara Dados Abertos: documentado preliminarmente em `docs/data-sources/camara-dados-abertos.md`, sem ingestao automatica.
- Senado Dados Abertos: documentado preliminarmente em `docs/data-sources/senado-dados-abertos.md`, sem ingestao automatica.
- LexML: documentado preliminarmente em `docs/data-sources/lexml.md`, sem ingestao automatica.
- CNJ/Datajud: documentado preliminarmente em `docs/data-sources/cnj-datajud.md`, sem ingestao automatica.
- STF: documentado preliminarmente em `docs/data-sources/stf.md`, sem ingestao automatica.
- STJ: documentado preliminarmente em `docs/data-sources/stj.md`, sem ingestao automatica.
- CJF: documentado preliminarmente em `docs/data-sources/cjf.md`, sem ingestao automatica.
- APIs auxiliares candidatas: documentadas em `docs/data-sources/api-candidates.md`, sem ingestao automatica.

Nenhum scraping cego esta autorizado neste ciclo.

Antes de qualquer conteudo juridico em escala, a fonte correta deve ser pesquisada e documentada. A ausencia de fonte documentada impede publicacao indexavel, mesmo que exista demanda comercial ou CTA.

O registro operacional fica em `content/source_registry.json`. Durante P0, todas as fontes devem manter `ingestion_enabled=false`.

As fontes listadas sao referencia e proveniencia. Elas nao autorizam scraping, clonagem, espelhamento ou criacao de paginas mecanicas. O portal deve escrever conteudo proprio, natural e unico a partir de pesquisa critica.

## Banco leve separado

Ingestao de termos juridicos pode iniciar a producao de conteudo apenas como rascunho (`draft_only`). Esses termos ficam em `data/terms/legal_terms.jsonl`, separados da auditoria de fonte, snapshots oficiais, rascunho editorial e manifesto publicado.

O contrato de armazenamento fica em `content/storage_contract.json` e e validado por `./tools/check-storage-contract`. Nenhuma fonte auditada, termo juridico ou snapshot autorizado pode ser tratado como conteudo publico pronto. O caminho correto e: termo com proveniencia -> rascunho editorial proprio -> revisao -> qualidade -> SEO -> publicacao.

As sementes iniciais ficam em `data/terms/legal_terms.jsonl` e sao validadas por `./tools/check-term-seeds`. Isso permite iniciar laboratorio editorial por termos sem gerar paginas para Googlebot.

Rascunhos persistidos ficam em `data/editorial/drafts.jsonl`, separados de snapshots oficiais e de paginas publicas. Eles devem permanecer `draft/noindex` ate fonte, revisao, qualidade e decisao editorial completa.

A fila de revisao fica em `data/editorial/review_queue.jsonl`. Ela registra autoria, motivo, historico e estado `needs_review`, mas nao libera publicacao nem cria URL.
