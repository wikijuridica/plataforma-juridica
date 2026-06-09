# SEO_CRAWL_INDEXING.md

## Indexacao

Toda pagina com `status != published` ou `index_policy != index` recebe `noindex,follow` e fica fora de sitemap.

Toda pagina indexavel precisa de:
- HTML textual completo;
- HTML leve, sem JavaScript, sem runtime cliente, sem hidratacao, sem bundle e sem CSS inline excessivo;
- `<title>` unico;
- meta description unica;
- canonical absoluto;
- meta robots `index,follow`;
- links internos rastreaveis;
- conteudo util e nao duplicado.

## Orcamento de HTML

Pagina publica indexavel deve ser facil de rastrear. O contrato atual reprova HTML publico acima de 50 KB, `<script>`, referencias `.js/.mjs/.wasm`, `modulepreload`, import maps, payloads de framework e marcadores de hidratacao. O objetivo e manter o primeiro response barato, textual e previsivel para Googlebot, OAI-SearchBot e bots valiosos.

## Crawl

`content/crawl_policy.json` configura bots. Googlebot, OAI-SearchBot, GPTBot e regra geral sao separados. Busca interna e parametros ficam bloqueados em robots.txt e tambem usam `noindex` quando renderizados.

Robots.txt sozinho nao e usado como substituto de `noindex`.
