# DATA_SOURCES.md

Antes de ingestao, cada fonte oficial precisa de documento em `docs/data-sources/`.

Estado inicial:
- Planalto: documentado preliminarmente em `docs/data-sources/planalto.md`, sem ingestao automatica.
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
