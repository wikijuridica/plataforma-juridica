# AGENTS.md — Portal Jurídico Brasileiro de Alta Escala

## Missão do projeto

Este repositório é um novo projeto de site jurídico brasileiro de alta escala. O objetivo é construir uma plataforma própria para wiki jurídica, legislação, jurisprudência, notícias, artigos, blog, glossário, perguntas e hubs temáticos, com potencial de centenas de milhares ou milhões de URLs.

O produto deve servir primeiro a humanos: advogados, estudantes, pesquisadores, jornalistas, empresas e cidadãos. Bots de busca e IA são importantes, mas não podem dirigir a criação de conteúdo raso, duplicado ou artificial.

## Postura obrigatória do agente

Trabalhe com autonomia agressiva e responsabilidade técnica.

Antes de alterar código:
1. Investigue o repositório.
2. Leia documentação existente.
3. Entenda padrões.
4. Identifique riscos.
5. Planeje o ciclo.
6. Só então implemente.

É proibido alterar no escuro.

Não seja passivo. Se uma decisão puder ser tomada com base nos requisitos, tome a decisão, documente em `docs/DECISIONS.md`, implemente, teste e siga.

Só peça intervenção humana quando houver bloqueio P0 real.

## Prioridades

### P0 — Inviolável

- Não usar Next.js.
- Não usar frameworks frontend pesados em páginas públicas.
- Não usar React/Vue/Svelte/Astro/Nuxt nas páginas públicas.
- Não usar dependências externas por conveniência.
- Não usar código copiado de terceiros.
- Não usar serviços SaaS externos como requisito do produto.
- Não publicar conteúdo jurídico sem fonte, data, autoria/revisão e aviso informativo.
- Não criar páginas rasas, duplicadas, parecidas ou feitas só para manipular busca.
- Não fazer scraping cego.
- Não ignorar robots.txt, termos de uso, sigilo processual, privacidade ou LGPD.
- Não avançar com validação P0 falhando.

### P1 — Plataforma indexável

- HTML textual completo no primeiro response.
- Links internos rastreáveis com `<a href>`.
- Canonical em toda página indexável.
- Meta robots correto.
- Sitemap index.
- Sitemaps particionados.
- robots.txt.
- URLs limpas, estáveis e humanas.
- Busca interna, filtros, parâmetros e páginas internas devem ser `noindex` por padrão.

### P2 — Conteúdo e dados

- Cada URL indexável deve ter intenção única.
- Cada página deve ter `unique_intent_id`.
- Cada página deve ter `canonical_url`.
- Cada conteúdo jurídico deve ter `source_provenance`.
- Conteúdo oficial, comentário editorial e opinião devem ser claramente separados.
- Fontes oficiais devem ser documentadas antes da ingestão.

### P3 — Performance e escala

- Páginas públicas devem ser leves.
- Não hidratar página inteira.
- Não exigir JavaScript para ler conteúdo.
- Preferir renderização server-side/static-first.
- Cachear conteúdo estável.
- Preparar geração para milhões de URLs sem criar URLs infinitas.

### P4 — Operação editorial

- Estados mínimos: `draft`, `needs_review`, `approved`, `published`, `noindex`, `archived`.
- Registrar autor, revisor, fonte, data de criação e data de revisão.
- Registrar motivo de publicação.
- Registrar histórico de decisões.

### P5 — Evolução

- Busca interna própria.
- Grafo jurídico.
- Recomendações internas.
- Páginas de comparação.
- Dados estruturados avançados.
- Otimizações para bots de IA, desde que não sacrifiquem humanos.

## Restrições de tecnologia

Este projeto deve ter código próprio.

Permitido por padrão:
- linguagem escolhida;
- biblioteca padrão;
- compilador/runtime;
- sistema operacional;
- banco self-hosted quando documentado;
- scripts locais próprios;
- ferramentas básicas de teste/build quando inevitáveis.

Proibido por padrão:
- Next.js;
- frameworks frontend pesados;
- CMS pronto;
- plugins de SEO prontos;
- bibliotecas externas para resolver problema que pode ser resolvido internamente;
- SDK de serviço externo;
- SaaS como requisito de funcionamento;
- código copiado de projetos open source.

## Regra para dependência excepcional

Antes de adicionar qualquer dependência externa, criar um ADR em:

`docs/adr/ADR-XXXX-dependency-NOME.md`

O ADR deve conter:
- problema;
- alternativas internas avaliadas;
- por que código próprio não basta;
- licença;
- riscos de segurança;
- impacto de performance;
- impacto de manutenção;
- plano de remoção;
- decisão final.

Sem ADR aprovado pelo próprio agente e registrado em checkpoint, a dependência é proibida.

## Arquitetura esperada

A arquitetura deve seguir estes módulos conceituais:

- `render`: renderização HTML própria.
- `router`: roteamento canônico.
- `content`: modelos de conteúdo.
- `legal`: entidades jurídicas.
- `sources`: fontes oficiais e proveniência.
- `quality`: verificadores antispam, duplicidade e conteúdo raso.
- `seo`: canonical, meta, robots, sitemap, index/noindex.
- `crawl`: políticas de bots e crawlability.
- `editorial`: status, revisão, autoria e auditoria.
- `tools`: scripts locais de validação.
- `docs`: arquitetura, decisões e fontes.

A estrutura real pode variar conforme a linguagem escolhida, mas os conceitos precisam existir.

## Tipos de página

Tipos iniciais:

- Home.
- Wiki jurídica.
- Legislação.
- Dispositivo legal.
- Jurisprudência.
- Precedente/tema.
- Notícia.
- Artigo/blog.
- Tema/hub.
- Glossário.
- Pergunta/resposta.
- Autor.
- Fonte oficial.

Cada tipo deve ter contrato de conteúdo próprio.

## Política de URL

URLs devem ser:
- estáveis;
- legíveis;
- sem parâmetros para conteúdo canônico;
- sem duplicidade por maiúsculas/minúsculas;
- sem acento;
- sem slug mutável sem redirect;
- conectadas a uma entidade canônica.

Exemplos:
- `/wiki/direito-civil/responsabilidade-civil/`
- `/legislacao/lei/10406-2002/codigo-civil/`
- `/jurisprudencia/stj/resp/1234567/tema-exemplo/`
- `/temas/dano-moral/`
- `/glossario/coisa-julgada/`
- `/noticias/2026/06/stf-decide-tema-exemplo/`

## Política de indexação

Toda página começa como não indexável até passar validação.

Indexável somente se:
- tem intenção única;
- tem canonical;
- tem conteúdo útil;
- tem fonte quando jurídica;
- tem título único;
- tem meta description única;
- não compete com outra URL;
- não é resultado de busca interna;
- não é filtro;
- não é página parametrizada;
- não é duplicada;
- não é thin content;
- não está em revisão.

## Política de conteúdo jurídico

Todo conteúdo jurídico deve conter:
- finalidade informativa;
- data de publicação;
- data de revisão;
- autor ou fonte oficial;
- revisor quando houver comentário editorial;
- fontes oficiais;
- distinção entre texto oficial, resumo, comentário e opinião;
- aviso de que não substitui consulta jurídica individual.

É proibido inventar:
- decisões;
- ementas;
- números de processo;
- trechos legais;
- citações;
- fontes;
- teses;
- datas;
- autores.

## Gate de qualidade de conteúdo

Implementar e manter verificadores para:

- hash normalizado de conteúdo;
- similaridade por shingles ou método equivalente;
- títulos duplicados;
- metas duplicadas;
- intenção duplicada;
- conteúdo sem fonte;
- conteúdo jurídico sem revisão;
- página rasa;
- página criada por permutação de keyword;
- página que deveria ser canonical de outra;
- página sem links internos úteis;
- página sem valor adicional para o usuário.

Se falhar, o conteúdo fica `draft`, `needs_review` ou `noindex`.

## Fontes oficiais

Antes de implementar ingestão, documentar em `docs/data-sources/`:

- nome da fonte;
- órgão;
- URL base;
- tipo de dado;
- formato;
- atualização;
- licença/termo;
- robots.txt;
- limites;
- campos disponíveis;
- riscos;
- estratégia de cache;
- estratégia de proveniência;
- estratégia de deduplicação.

Fontes iniciais a estudar:
- Planalto/Portal da Legislação.
- Câmara Dados Abertos.
- Senado Dados Abertos.
- LexML.
- CNJ/Datajud.
- STF.
- STJ.
- CJF.

## SEO técnico obrigatório

Toda página pública indexável deve ter:

- `<title>` único;
- meta description única;
- canonical;
- meta robots coerente;
- HTML semântico;
- conteúdo principal em texto;
- links internos rastreáveis;
- status HTTP correto;
- sitemap quando publicada;
- data de atualização quando aplicável;
- dados estruturados apenas se corresponderem ao conteúdo visível.

Não criar páginas indexáveis para:
- busca interna;
- filtros;
- combinações infinitas;
- tags fracas;
- páginas vazias;
- páginas sem fonte;
- duplicatas;
- parâmetros;
- ordenações.

## Bots

Criar política configurável, não hardcoded, para bots.

Prioridade:
- permitir Googlebot em conteúdo público aprovado;
- permitir OAI-SearchBot em conteúdo público aprovado quando a estratégia do projeto desejar presença em ChatGPT Search;
- controlar GPTBot separadamente;
- bloquear ou limitar bots abusivos;
- nunca bloquear conteúdo importante por acidente;
- usar `noindex` para impedir indexação quando necessário, não apenas robots.txt.

## Testes obrigatórios

Criar aliases ou scripts equivalentes:

- `./tools/check-all`
- `./tools/check-architecture`
- `./tools/check-content-quality`
- `./tools/check-seo`
- `./tools/check-crawlability`
- `./tools/check-sitemaps`
- `./tools/check-canonicals`
- `./tools/check-no-duplicate-content`
- `./tools/check-performance-budget`

Todo ciclo deve rodar validações relevantes.

Se uma validação falhar:
1. Pare o avanço.
2. Corrija.
3. Rode novamente.
4. Registre no checkpoint.

## Checkpoint obrigatório

Ao fim de cada ciclo, atualizar `CHECKPOINT.md` com:

- data/hora;
- ciclo;
- prioridade P0-P5;
- objetivo;
- entregas;
- arquivos alterados;
- decisões;
- comandos executados;
- resultados;
- falhas;
- correções;
- provas;
- próximo ciclo;
- riscos.

Responder também no thread:

CHECKPOINT:
- Entregue:
- Provas:
- Testes:
- Próximo:
- Bloqueios:

## Documentação obrigatória

Manter:

- `docs/PROJECT_VISION.md`
- `docs/ARCHITECTURE.md`
- `docs/CONTENT_QUALITY.md`
- `docs/SEO_CRAWL_INDEXING.md`
- `docs/DATA_SOURCES.md`
- `docs/ROADMAP_P0_P5.md`
- `docs/DECISIONS.md`
- `docs/adr/`
- `CHECKPOINT.md`

## Definição de pronto

Uma entrega só está pronta quando:

- código foi implementado;
- documentação foi atualizada;
- testes relevantes passaram;
- SEO técnico foi verificado;
- qualidade de conteúdo foi verificada;
- não há duplicidade conhecida;
- não há thin content indexável;
- não há dependência externa sem ADR;
- checkpoint foi atualizado com prova.

Não entregue “parece funcionar”. Entregue comprovado.
