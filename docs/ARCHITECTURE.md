# ARCHITECTURE.md

## Decisao de base

O projeto usa Go e biblioteca padrao como base tecnica. A escolha privilegia binario proprio, tipagem estatica, concorrencia nativa, `net/http`, `encoding/xml`, `html` e `testing`, sem SDK externo e sem framework frontend.

## Modulos

- `cmd/build`: gera artefatos publicos em `public/`.
- `cmd/check`: executa validadores locais.
- `cmd/server`: entrega paginas com geracao on demand propria.
- `internal/content`: modelos e carregamento do manifesto de paginas.
- `internal/manualresearch`: banco leve de pesquisa editorial manual de termos de alta intencao digital.
- `internal/contentbriefs`: briefs editoriais iniciais, nao-template e nao publicaveis.
- `internal/scale`: plano finito de blueprints para pelo menos 10 mil paginas sem publicar conteudo durante P0.
- `internal/cta`: politica propria de CTA WhatsApp subordinada a fonte, revisao e aprovacao editorial.
- `internal/router`: rotas canonicas limpas e mapeamento para cache/saida.
- `internal/render`: HTML completo textual, sem hidratacao.
- `internal/ondemand`: geracao sob demanda e cache local controlado pelo projeto.
- `internal/seo`: title, meta description, canonical e robots.
- `internal/crawl`: robots.txt e politica configuravel de bots.
- `internal/sitemap`: sitemap index e sitemap particionado.
- `internal/quality`: gates contra duplicidade, conteudo raso, falta de fonte e falta de revisao.
- `internal/editorial`: estados editoriais e politica index/noindex.
- `internal/editorialdrafts`: persistencia de rascunhos em banco leve, sempre sem rota publica.
- `internal/reviewqueue`: fila editorial `needs_review` com historico e publicacao bloqueada.
- `internal/approvals`: aprovacao editorial separada de publicacao, ainda sem URL publica.
- `internal/publicationblockers`: manifesto de requisitos faltantes antes de qualquer URL publica.
- `internal/sourceblockers`: bloqueios por termo quando a fonte atual ainda nao e especifica o suficiente para aprovacao/publicacao.
- `internal/legal`: controle de conteudo juridico.
- `internal/sources`: contratos de fontes oficiais.
- `internal/storage`: banco leve proprio em JSONL para termos juridicos, auditorias, snapshots, rascunhos e manifesto publicado.
- `internal/termintents`: candidatos de termos juridicos com demanda humana, fonte oficial e adequacao a contratacao 100% digital.
- `internal/termpromotion`: ranking e promocao controlada de candidatos para seeds `draft_only`, com diversidade de areas e penalizacao de sinais fracos.
- `cmd/refresh-editorial-drafts`: regeneracao segura de rascunhos persistidos quando o algoritmo de escrita e refinado.
- `internal/provenance`: contrato de proveniencia por payload antes de qualquer conteudo.
- `internal/architecture`: validacao de estrutura e proibicoes P0.
- `tools`: scripts locais obrigatorios.

## Geracao on demand obrigatoria

Rotas publicas devem poder ser geradas sob demanda por `internal/ondemand`, sem Next.js e sem mecanismo terceirizado de ISR/SSR. O gerador resolve uma rota canonica, renderiza HTML completo no primeiro response e grava cache local. A segunda chamada pode servir do cache do proprio projeto.

O build estatico continua permitido como artefato operacional, mas nao substitui o requisito de geracao on demand propria.

## HTML publico leve

Leveza e requisito de indexacao. O renderizador deve entregar HTML textual completo, com CSS minimo e sem JavaScript, bundle, WebAssembly, import map, `modulepreload`, payload de framework ou marcador de hidratacao. Isso protege Googlebot, OAI-SearchBot e outros bots valiosos, alem de reduzir custo operacional em escala massiva.

`internal/checks` reprova HTML publico acima do orcamento, CSS inline excessivo, referencias a runtime cliente e marcadores de Next.js, React, Vue/Svelte/Astro/Angular, Vite ou Webpack. Qualquer excecao exigiria ADR de dependencia e continuaria bloqueada para pagina publica indexavel enquanto P0 estiver ativo.

## Escala antes de conteudo

Durante P0, a plataforma deve preparar escala, nao fabricar paginas juridicas. `content/scale_plan.json` define blueprints finitos para pelo menos 10 mil paginas planejadas e `published_pages_during_p0` deve permanecer zero. O desbloqueio de conteudo exige P2: fonte oficial documentada, proveniencia, autoria, revisao, intencao unica, qualidade e indexacao coerente.

O plano de escala e um contrato de capacidade, nao um gerador de spam para Google. O modulo `scale` deve ajudar a medir alcance futuro e bloquear publicacao em massa enquanto as fontes e validadores nao estiverem maduros.

## CTA WhatsApp

`content/cta_policy.json` registra a arquitetura de CTA proprio por WhatsApp para paginas de alta intencao de contratar advogado. O CTA nao pode aparecer como atalho para publicar conteudo sem fonte ou sem revisao; ele depende de pagina juridica aprovada, proveniencia e revisao editorial.

CTA WhatsApp e critico para o produto: as paginas devem ser informativas, mas desenhadas para alta intencao de contratar advogado quando o contexto for adequado. A criticidade comercial nao remove os gates juridicos.

## Proveniencia por payload

Antes de qualquer dado oficial virar conteudo, o payload precisa de registro com fonte, URL oficial, data de acesso, hash SHA-256, snapshot de robots, snapshot de termos, campos usados e finalidade. Fonte pesquisada nao equivale a conteudo aprovado.

## Banco leve de termos e ingestao

Ingestao de termos juridicos e valida para iniciar conteudos somente como semente de rascunho. O armazenamento e proprio, leve e separado em `content/storage_contract.json`, usando arquivos JSONL em `data/` e biblioteca padrao Go.

Camadas obrigatorias:
- `term_seeds`: termos juridicos para iniciar rascunhos, sem texto oficial bruto e sem texto editorial publico;
- `term_intent_candidates`: candidatos priorizados por demanda humana e contratacao online, ainda sem publicacao;
- `manual_keyword_research`: pesquisa editorial manual de alta intencao digital, com Trends como orientacao e fontes oficiais como autoridade;
- `term_seeds` promovidos: seeds com `candidate_id`, evidencia de demanda, modo `digital_only` e CTA alto, mas ainda `draft_only`;
- `source_audits`: auditoria de robots, termos de uso, alcance HTTP e decisao de bloqueio;
- `source_snapshots`: snapshots autorizados, pequenos, com hash e proveniencia;
- `editorial_drafts`: texto editorial proprio em PT-BR, sempre noindex ate aprovacao;
- `content_briefs`: brief inicial natural e especifico por termo, sem URL publica;
- `source_specificity_blockers`: manifesto que impede aprovacao/publicacao quando o termo ainda precisa fonte primaria, norma especifica ou recorte juridico;
- `published_manifest`: manifesto leve de conteudo aprovado, sem substituir o renderizador.

Regra P0: termos podem iniciar `draft_only`; nenhuma linha do banco vira pagina indexavel sem fonte, revisao, qualidade, SEO, intencao unica e checkpoint.

## Laboratorio

Toda mudanca P0/P1 deve passar por ciclo de laboratorio: escrever ou ajustar teste, rodar validacao, refinar, testar novamente e inspecionar artefatos. `tools/lab-cycle` combina `go test -count=1 ./...`, `tools/check-all`, build, auditoria de dependencias, diff check e busca por residuos Python.
