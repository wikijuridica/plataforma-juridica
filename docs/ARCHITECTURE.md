# ARCHITECTURE.md

## Decisao de base

O projeto usa Go e biblioteca padrao como base tecnica. A escolha privilegia binario proprio, tipagem estatica, concorrencia nativa, `net/http`, `encoding/xml`, `html` e `testing`, sem SDK externo e sem framework frontend.

## Modulos

- `cmd/build`: gera artefatos publicos em `public/`.
- `cmd/check`: executa validadores locais.
- `cmd/server`: entrega paginas com geracao on demand propria.
- `internal/content`: modelos e carregamento do manifesto de paginas.
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
- `internal/legal`: controle de conteudo juridico.
- `internal/sources`: contratos de fontes oficiais.
- `internal/provenance`: contrato de proveniencia por payload antes de qualquer conteudo.
- `internal/architecture`: validacao de estrutura e proibicoes P0.
- `tools`: scripts locais obrigatorios.

## Geracao on demand obrigatoria

Rotas publicas devem poder ser geradas sob demanda por `internal/ondemand`, sem Next.js e sem mecanismo terceirizado de ISR/SSR. O gerador resolve uma rota canonica, renderiza HTML completo no primeiro response e grava cache local. A segunda chamada pode servir do cache do proprio projeto.

O build estatico continua permitido como artefato operacional, mas nao substitui o requisito de geracao on demand propria.

## Escala antes de conteudo

Durante P0, a plataforma deve preparar escala, nao fabricar paginas juridicas. `content/scale_plan.json` define blueprints finitos para pelo menos 10 mil paginas planejadas e `published_pages_during_p0` deve permanecer zero. O desbloqueio de conteudo exige P2: fonte oficial documentada, proveniencia, autoria, revisao, intencao unica, qualidade e indexacao coerente.

O plano de escala e um contrato de capacidade, nao um gerador de spam para Google. O modulo `scale` deve ajudar a medir alcance futuro e bloquear publicacao em massa enquanto as fontes e validadores nao estiverem maduros.

## CTA WhatsApp

`content/cta_policy.json` registra a arquitetura de CTA proprio por WhatsApp para paginas de alta intencao de contratar advogado. O CTA nao pode aparecer como atalho para publicar conteudo sem fonte ou sem revisao; ele depende de pagina juridica aprovada, proveniencia e revisao editorial.

CTA WhatsApp e critico para o produto: as paginas devem ser informativas, mas desenhadas para alta intencao de contratar advogado quando o contexto for adequado. A criticidade comercial nao remove os gates juridicos.

## Proveniencia por payload

Antes de qualquer dado oficial virar conteudo, o payload precisa de registro com fonte, URL oficial, data de acesso, hash SHA-256, snapshot de robots, snapshot de termos, campos usados e finalidade. Fonte pesquisada nao equivale a conteudo aprovado.

## Laboratorio

Toda mudanca P0/P1 deve passar por ciclo de laboratorio: escrever ou ajustar teste, rodar validacao, refinar, testar novamente e inspecionar artefatos. `tools/lab-cycle` combina `go test -count=1 ./...`, `tools/check-all`, build, auditoria de dependencias, diff check e busca por residuos Python.
