# SEO_CRAWL_INDEXING.md

## Indexacao

Toda pagina com `status != published` ou `index_policy != index` recebe `noindex,follow` e fica fora de sitemap.

Toda pagina indexavel precisa de:
- HTML textual completo;
- `<title>` unico;
- meta description unica;
- canonical absoluto;
- meta robots `index,follow`;
- links internos rastreaveis;
- conteudo util e nao duplicado.

## Crawl

`content/crawl_policy.json` configura bots. Googlebot, OAI-SearchBot, GPTBot e regra geral sao separados. Busca interna e parametros ficam bloqueados em robots.txt e tambem usam `noindex` quando renderizados.

Robots.txt sozinho nao e usado como substituto de `noindex`.
