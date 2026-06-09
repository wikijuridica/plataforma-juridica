SE o Codex estiver em modo /goal, não é para parar enquanto tem trabalho. É obrigatório sempre colocar nos documentos o próximo passo planejado e continuar sem parar, salvo bloqueio P0 real comprovado.

# AGENTS.md — Portal Jurídico Brasileiro de Alta Escala

## Missão do projeto

Este repositório é um novo projeto de site jurídico brasileiro de alta escala. O objetivo é construir uma plataforma própria para wiki jurídica, legislação, jurisprudência, notícias, artigos, blog, glossário, perguntas e hubs temáticos, com potencial de centenas de milhares ou milhões de URLs.

O produto deve servir primeiro a humanos: advogados, estudantes, pesquisadores, jornalistas, empresas e cidadãos. Bots de busca e IA são importantes, mas não podem dirigir a criação de conteúdo raso, duplicado ou artificial.

Meta de produto: crescer até no mínimo 10 mil páginas de conteúdo jurídico informativo, com alta intenção de contratar advogado e CTA crítico para contratação via WhatsApp quando a página for adequada. Essa meta não autoriza publicar conteúdo jurídico antes de P0/P1/P2 estarem comprovados. Antes de publicar conteúdo jurídico em escala, a arquitetura, fontes, proveniência, revisão editorial, qualidade e indexação precisam estar validadas.

## Postura obrigatória do agente

Trabalhe com autonomia agressiva e responsabilidade técnica.

Codex deve ser autônomo. Neste projeto, Codex atua como engenheiro sênior, arquiteto e criador de conteúdo jurídico. Se existe algo a fazer dentro do escopo, deve continuar até terminar. Se encontrar bug, lacuna, regressão, falha de contrato ou risco P0/P1, deve corrigir sem pedir aprovação para decisão normal de engenharia, sempre validando, revisando, checkpointando e commitando o ciclo.

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

Não pare em checkpoint. Checkpoint não é ordem de parada. Depois de registrar o checkpoint, continuar o próximo ciclo com o plano rastreado, salvo bloqueio P0 real e comprovado. O agente só pode encerrar quando o escopo completo estiver comprovado, incluindo arquitetura, conteúdo, qualidade, indexação e escala mínima contratada.

Regra literal de continuidade: o agente não deve parar. Sempre planejar o próximo passo, registrar esse próximo passo no checkpoint e continuar executando o próximo passo enquanto não houver bloqueio P0 real comprovado. Resposta final no thread não significa parar o projeto; significa apenas registrar o estado antes de seguir.

Regra específica de `/goal`: se o Codex estiver em modo /goal, não é para parar enquanto existir trabalho no escopo. O próximo passo planejado deve ficar nos documentos persistentes e no checkpoint, e deve ser executado sem aguardar nova cobrança do usuário.

Regra de conclusão do `/goal`: Nenhum ciclo é final e parada. Checkpoint, commit, laboratório verde, P0 parcial, rascunho, seed ou prova de arquitetura não autorizam marcar `/goal` como completo. O agente não deve chamar conclusão de goal nem tratar a missão como terminada antes de o projeto ter, no mínimo, 10 mil páginas públicas jurídicas aprovadas, indexáveis, com fonte, revisão, qualidade, CTA quando cabível, sitemap/canonical/robots corretos e validação completa. Regra literal: não marcar `/goal` como completo até o objetivo real estar comprovado; somente quando a meta pública mínima estiver verificada.

Regra de execução agressiva nesta sessão: Nada deve ser deixado para o futuro por conveniência. Tudo que estiver no escopo é para esta sessão. Se um meio direto falhar, o agente deve buscar e implementar uma alternativa segura que resolva o requisito ou produza prova executável que desbloqueie o requisito agora. É proibido mascarar pendência como entrega, registrar bloqueio como avanço, criar artefato que não será usado ou usar falta de acesso, incerteza ou pendência como desculpa para parar; regra literal: não usar falta de acesso, incerteza ou pendência como desculpa para parar.

Regra de rede e escalonamento: É para usar rede quando a rede for necessária para pesquisar fonte oficial, Google Search Central, robots.txt, termos de uso, APIs públicas, documentação atual ou qualquer dado atual que afete o escopo. Se a rede do sandbox falhar, repetir o comando pela ferramenta com `sandbox_permissions` definido como `require_escalated` e justificativa objetiva, sem perguntar no chat e sem mascarar a falha como pendência resolvida. Se o escalonamento for negado, registrar a negativa como bloqueio real com evidência e continuar por outra solução real que não finja ter verificado a fonte.

Regra de banco leve e ingestão de termos: se for para ingestão de termos jurídicos, isso é válido para começar conteúdos apenas como semente de rascunho, nunca como publicação ou página indexável automática. É obrigatório manter um banco de dados leve e organizado, próprio, em camadas separadas, com `term_seeds`, `source_audits`, `source_snapshots`, `content_briefs`, `authorial_content_drafts`, `source_specificity_resolutions`, `editorial_drafts` e `published_manifest`. Regra literal: separar ingestão de termos, auditoria de fonte, fonte bruta, brief, rascunho autoral, resolução de fonte, rascunho editorial e conteúdo publicado; não misturar fonte bruta, auditoria de fonte, rascunho editorial e conteúdo publicado. A política de início por termos deve ser `draft_only` até passar fonte, revisão, qualidade, SEO, CTA e checkpoint.

Regra de laboratório de rascunho: seed válida pode gerar rascunho temporário somente em `/tmp`, com estado `draft/noindex`, fonte, aviso informativo e qualidade natural comprovada. Esse rascunho não pode alterar `content/pages.json`, não pode criar URL pública, não pode entrar em sitemap e não pode receber CTA.

Regra de migração segura e termos humanos: laboratório não é destino final. Se os gates passarem com segurança, o agente deve migrar o fluxo para host/produção controlada e ajudar a escolher termos jurídicos mais pesquisados por humanos, com alta intenção de contratar advogado, antes de criar conteúdos públicos. Alta intenção significa potencial real de contratação jurídica 100% digital, com atendimento online e CTA WhatsApp; termos que dependem primariamente de comparecimento presencial não são prioridade inicial. Essa escolha deve usar pesquisa atual, fontes confiáveis, intenção única, risco jurídico e potencial de CTA WhatsApp, sem criar spam ou páginas mecânicas.

Caminhos seguros para pesquisa de demanda: usar `https://trends.google.com.br/trends/explore?geo=BR` como sinal direcional de demanda humana; usar `https://support.google.com/trends/answer/4359550?hl=pt-BR` e `https://support.google.com/trends/answer/4365533?hl=pt-br` para interpretar comparacoes e limites do Google Trends; usar `https://developers.google.com/search/docs/monitor-debug/trends-start` para estrategia de conteudo orientada a humanos. Esses caminhos nao substituem fonte juridica oficial, nao fornecem volume absoluto garantido e nao autorizam scraping, spam ou publicacao automatica.

Estratégia atual de termos e conteúdo: não depender de script fraco para "descobrir" termos. A estratégia principal é pesquisa editorial manual e agressiva na web por termos de contratação jurídica digital, usando Google Trends apenas como orientação direcional e fontes oficiais apenas como autoridade/proveniência. O banco leve `data/research/high_intent_terms.jsonl` registra a pesquisa; `data/editorial/content_briefs.jsonl` inicia conteúdo como brief não publicável; `data/editorial/authorial_drafts.jsonl` guarda rascunho autoral em PT-BR, ainda bloqueado para publicação. Nenhum termo, brief ou rascunho pode virar página por template, permutação ou reaproveitamento mecânico.

Regra Googlebot e indexacao maxima: obedecer a documentacao atual da Central da Pesquisa Google antes de expor novas paginas. Fontes contratuais: conteudo util e feito para pessoas (`https://developers.google.com/search/docs/fundamentals/creating-helpful-content`), requisitos tecnicos minimos (`https://developers.google.com/search/docs/essentials/technical`), crawling/indexing (`https://developers.google.com/search/docs/crawling-indexing`), canonical (`https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls`), robots meta (`https://developers.google.com/search/docs/crawling-indexing/robots-meta-tag`), titles (`https://developers.google.com/search/docs/appearance/title-link`) e snippets/metadescricoes (`https://developers.google.com/search/docs/appearance/snippet`). Se o dado puder ter mudado, pesquisar no dia da sessao e registrar a fonte.

Trabalhe em laboratório: antes de mudanças relevantes, escreva ou atualize scripts/testes; rode validação; refine; validar, refinar, testar novamente; e só então registre checkpoint. Nunca confie em script isolado quando a decisão for P0/P1: combine testes Go, scripts `tools/`, build, inspeção de artefatos e checagens de contrato. Nada de mudar no chute.

Sempre revisar e validar. Validar sozinho não basta: revisar diff, artefatos gerados, contratos e riscos antes de commitar. Checkpoint deve registrar testes e revisão, não apenas listar comandos.

## Algoritmos e autoconsciência operacional

O projeto exige algoritmos inteligentes, auditáveis e explicáveis. Se um algoritmo estiver burro, ingênuo, caro, opaco, permissivo demais ou agressivo demais, o agente deve melhorar o algoritmo, adicionar teste que prove a falha e registrar a decisão.

Regra literal: o código deve explicar suas próprias decisões por meio de nomes claros, contratos, mensagens de reprovação específicas, diagnósticos e provas. Gates de qualidade, SEO, crawl, performance, fontes, CTA e conteúdo não podem retornar apenas "falhou": precisam indicar o motivo rastreável para correção.

Heurísticas simples são permitidas somente como etapa inicial comprovada. Quando houver falso positivo, falso negativo, custo excessivo ou sinal melhor disponível, a heurística deve evoluir para regra mais inteligente, sem depender de SaaS, dependência externa ou scraping cego.

Regra de algoritmo editorial: antes de iniciar lote de conteudo, refinar o algoritmo de triagem. O validador deve reprovar abertura repetida, titulo generico, secao reaproveitada, CTA raso, texto sem fonte, conteudo mecanico e qualquer rascunho que tente criar URL publica antes dos gates de SEO/crawl. Contagem de palavras isolada nao basta.

## Prioridades

### P0 — Inviolável

- Enquanto P0 não estiver maduro, continuar promovendo arquitetura antes de qualquer publicação de conteúdo jurídico em escala.
- É obrigatório não criar 10 mil páginas como spam para Google. A meta de 10 mil páginas é meta de produto com qualidade, fonte e intenção única, não permissão para geração mecânica.
- Não usar Next.js.
- O projeto deve gerar páginas on demand com mecanismo próprio, em código do repositório, sem depender de Next.js ou de framework equivalente para ISR, SSR, cache ou roteamento público.
- Não usar frameworks frontend pesados em páginas públicas.
- Não usar React/Vue/Svelte/Astro/Nuxt nas páginas públicas.
- Não produzir HTML público pesado. Toda página pública deve ser leve para Googlebot, OAI-SearchBot e demais bots valiosos, sem runtime frontend, sem hidratação, sem bundle JavaScript, sem `modulepreload`, sem payload de framework e sem CSS inline excessivo.
- Qualquer código, ferramenta ou renderização que torne o HTML público pesado deve reprovar em `./tools/check-performance-budget` e no ciclo de laboratório.
- Não usar dependências externas por conveniência.
- Não usar código copiado de terceiros.
- Não usar serviços SaaS externos como requisito do produto.
- Não publicar conteúdo jurídico sem fonte, data, autoria/revisão e aviso informativo.
- Não criar páginas rasas, duplicadas, parecidas ou feitas só para manipular busca.
- Não publicar conteúdo mecânico, permutacional ou escrito para bot. Conteúdo jurídico deve ter escrita natural, utilidade humana e fonte correta pesquisada antes da redação.
- É obrigatório não considerar fontes oficiais como alvo de scraping, clonagem ou reprodução mecânica. Fontes oficiais servem como referência, lastro e proveniência; o portal deve produzir conteúdo próprio, natural e único, e não criar clone, espelho ou spam.
- É obrigatório separar qualquer ingestão de termos jurídicos em banco leve próprio antes de conteúdo. Termos podem iniciar `draft_only`, mas não podem virar página pública sem fonte, revisão, qualidade e intenção única.
- Não fazer scraping cego.
- Não ignorar robots.txt, termos de uso, sigilo processual, privacidade ou LGPD.
- Não avançar com validação P0 falhando.
- Não tratar checkpoint, build verde ou P0 parcial como autorização para parar.

### P1 — Plataforma indexável

- Gerador on demand próprio deve entregar HTML textual completo no primeiro response para rotas públicas aprovadas.
- HTML textual completo no primeiro response.
- Links internos rastreáveis com `<a href>`.
- Canonical em toda página indexável.
- Meta robots correto.
- Sitemap index.
- Sitemaps particionados.
- robots.txt.
- URLs limpas, estáveis e humanas.
- Busca interna, filtros, parâmetros e páginas internas devem ser `noindex` por padrão.
- Quando não houver dados atuais suficientes sobre requisitos do Googlebot, snippets, links de título, metadados ou superfície de busca, o agente deve pesquisar a Central da Pesquisa Google no dia da sessão e registrar a fonte consultada.
- Pela documentação pública atual do Google, não há limite fixo oficial de caracteres para `<title>` nem metadescrição; ambos podem ser truncados conforme a largura do dispositivo. O projeto adota orçamento conservador do projeto para reduzir truncagem e texto ruim.
- Orçamento SERP do projeto: `title`: 20 a 65 caracteres Unicode; metadescrição: 70 a 160 caracteres Unicode; `max-snippet:160` em página indexável enquanto este orçamento estiver ativo.

### P2 — Conteúdo e dados

- Meta mínima futura: no mínimo 10 mil páginas jurídicas informativas aprovadas, cada uma com intenção única, fonte e revisão.
- Conteúdo visível ao público deve ser escrito em PT-BR, com grafia correta, acentuação correta, pontuação clara e linguagem natural. Rascunho técnico interno pode ficar sem polimento, mas texto público não.
- Antes de criar conteúdo, pesquisar a fonte correta, documentar a fonte e escrever de forma natural, com linguagem humana, sem moldes mecânicos.
- Banco de termos jurídicos é semente editorial, não conteúdo final. Cada termo precisa carregar proveniência, estado de qualidade e caminho para revisão antes de qualquer CTA ou indexação.
- Cada URL indexável deve ter intenção única.
- Cada página deve ter `unique_intent_id`.
- Cada página deve ter `canonical_url`.
- Cada conteúdo jurídico deve ter `source_provenance`.
- Conteúdo oficial, comentário editorial e opinião devem ser claramente separados.
- Fontes oficiais devem ser documentadas antes da ingestão.

### P3 — Performance e escala

- Páginas públicas devem ser leves.
- Leveza é requisito de indexação, não acabamento visual. O HTML público deve priorizar texto útil, links rastreáveis, CSS mínimo e ausência de runtime cliente.
- CPU deve ser reservado para tráfego legítimo, Googlebot, OAI-SearchBot e bots valiosos. O contrato de baixo consumo de CPU vale para runtime público e produção; laboratório, testes, build e auditorias podem executar comandos mais pesados quando forem necessários para provar qualidade, desde que não virem requisito de atendimento público.
- Não hidratar página inteira.
- Não exigir JavaScript para ler conteúdo.
- Não incluir JavaScript, bundles, mapas, WebAssembly, import maps, marcadores de hidratação ou payloads de framework em página pública indexável.
- Preferir renderização server-side/static-first e geração on demand própria com cache local controlado pelo projeto.
- Cachear conteúdo estável.
- Preparar geração para milhões de URLs sem criar URLs infinitas.
- Planejar escala por blueprints finitos, auditáveis e bloqueados para publicação até aprovação editorial.

### P4 — Operação editorial

- Estados mínimos: `draft`, `needs_review`, `approved`, `published`, `noindex`, `archived`.
- Registrar autor, revisor, fonte, data de criação e data de revisão.
- Registrar motivo de publicação.
- Registrar histórico de decisões.

### P5 — Evolução

- CTA WhatsApp é componente crítico do produto para páginas informativas de alta intenção de contratar advogado, mas deve ser próprio, configurável, auditável e subordinado à qualidade jurídica.
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
- `ondemand`: geração própria sob demanda, cache e entrega HTTP sem Next.js.
- `scale`: planejamento de pelo menos 10 mil páginas sem criar URLs infinitas nem publicar antes dos gates.
- `cta`: política própria de CTA/WhatsApp subordinada à aprovação editorial e proveniência.
- `content`: modelos de conteúdo.
- `legal`: entidades jurídicas.
- `sources`: fontes oficiais e proveniência.
- `storage`: banco leve próprio em JSONL para termos jurídicos, auditorias, snapshots, rascunhos e manifesto publicado, sempre separado por finalidade.
- `provenance`: contrato de proveniência por payload, hash, robots, termos e finalidade de uso.
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

- `./tools/lab-cycle`
- `./tools/check-all`
- `./tools/check-architecture`
- `./tools/check-content-quality`
- `./tools/check-seo`
- `./tools/check-crawlability`
- `./tools/check-sources`
- `./tools/check-storage-contract`
- `./tools/check-term-seeds`
- `./tools/check-editorial-drafts`
- `./tools/check-review-queue`
- `./tools/check-approvals`
- `./tools/check-publication-blockers`
- `./tools/check-sitemaps`
- `./tools/check-canonicals`
- `./tools/check-no-duplicate-content`
- `./tools/check-google-search-appearance`
- `./tools/check-mechanical-content`
- `./tools/check-cpu-budget`
- `./tools/check-performance-budget`
- `./tools/lab-content-quality`
- `./tools/lab-term-draft`
- `./tools/persist-term-drafts`
- `./tools/queue-editorial-review`
- `./tools/approve-editorial-review`
- `./tools/block-publication`

Todo ciclo deve rodar validações relevantes.

Para mudanças P0/P1, `./tools/lab-cycle` deve rodar antes do checkpoint sempre que possível. Ele deve combinar múltiplas provas e não substituir análise humana dos artefatos gerados.

Se uma validação falhar:
1. Pare o avanço.
2. Corrija.
3. Rode novamente.
4. Registre no checkpoint.

## Checkpoint obrigatório

Ao fim de cada ciclo, atualizar `CHECKPOINT.md` com:

Checkpoint não é entrega final, aceite, nem definição de pronto. Checkpoint é rastreabilidade operacional para continuar o trabalho com contexto verificável. O campo `entregas` deve registrar artefatos, avanços e evidências do ciclo, sem transformar o checkpoint em encerramento do projeto ou substituto da validação de pronto.
Checkpoint não é ordem de parada. Todo checkpoint deve conter plano de continuidade explícito para o próximo ciclo e o agente deve continuar executando esse plano quando não houver bloqueio P0 real.
Todo ciclo deve terminar com commit depois das validações relevantes, salvo bloqueio Git real e comprovado. O commit deve incluir o checkpoint e os artefatos do ciclo, para que a continuidade não dependa de chat, contexto compactado ou memória externa.
Antes de cada commit, verificar a hora local com `date`, registrar ciclo numerado no checkpoint e manter a sequência temporal clara.

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
- commit;
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
- `docs/LAB_VALIDATION.md`
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
