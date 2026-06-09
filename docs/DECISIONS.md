# DECISIONS.md

## 2026-06-09 — Go como base propria do portal

Decisao: usar Go e biblioteca padrao como base do projeto.

Motivos:
- linguagem compilada, tipada e adequada a servicos com concorrencia;
- biblioteca padrao suficiente para HTTP, XML, HTML escaping, arquivos e testes;
- binario proprio e operacao sem framework frontend;
- melhor alinhamento com geracao on demand propria e cache controlado pelo projeto.

Evidencia documental consultada:
- `https://go.dev/doc/`: documenta Go como linguagem eficiente, compilada, tipada, com concorrencia e toolchain oficial.
- `https://go.dev/doc/modules/layout`: recomenda `internal/` e `cmd/` para projetos de servidor.
- `https://pkg.go.dev/net/http`: biblioteca padrao para cliente e servidor HTTP.
- `https://pkg.go.dev/encoding/xml`: biblioteca padrao para XML, usada em sitemaps.
- `https://pkg.go.dev/testing`: biblioteca padrao para testes automatizados.

Alternativas avaliadas:
- Python com biblioteca padrao: suficiente para prototipar validadores, mas menos alinhado ao objetivo de portal massivo com binario proprio e geracao HTTP concorrente;
- Next.js ou frameworks frontend: proibidos pelo contrato P0.

Consequencia: qualquer artefato Python criado no ciclo foi descartado; o contrato e a implementacao passam a ser Go-first.

## 2026-06-09 — Geracao on demand propria obrigatoria

Decisao: o portal deve gerar paginas sob demanda com codigo proprio em `internal/ondemand`, sem Next.js, ISR terceirizado ou framework equivalente.

Motivos:
- o projeto quer autonomia sobre roteamento, cache e politica de indexacao;
- rotas futuras podem chegar a milhoes de URLs e nao devem depender de build total;
- o primeiro response de pagina publica deve conter HTML textual completo.

Consequencia: build estatico e permitido como ferramenta operacional, mas nao substitui o gerador on demand.

## 2026-06-09 — Checkpoint nao encerra trabalho

Decisao: checkpoint e rastreabilidade operacional para continuar, nao entrega final nem ordem de parada.

Motivos:
- o projeto tem escopo continuo e grande;
- cada ciclo precisa deixar contexto, provas, falhas e proximo plano;
- o agente deve continuar trabalhando enquanto nao houver bloqueio P0 real.

Consequencia: `CHECKPOINT.md` deve registrar plano de continuidade explicito e nao pode ser usado como substituto da definicao de pronto.

Adendo operacional: sempre planejar o proximo passo e continuar executando. Uma resposta no thread ou um checkpoint nao encerram a missao.

Adendo de persistencia: cada ciclo deve ser commitado apos validacao, incluindo `CHECKPOINT.md`, para preservar continuidade em Git.

## 2026-06-09 — Meta de 10 mil paginas com CTA subordinado a qualidade

Decisao: a meta de produto inclui no minimo 10 mil paginas juridicas informativas, com CTA proprio de WhatsApp para contratar advogado quando apropriado.

Motivos:
- o portal deve ter escala nacional e alta intencao comercial;
- a arquitetura precisa suportar grande volume sem URLs infinitas;
- CTA comercial nao pode sacrificar fonte, revisao, utilidade e seguranca juridica.

Consequencia: durante P0 o projeto cria plano de blueprints e politica de CTA, mas nao publica 10 mil paginas juridicas. A publicacao em escala depende dos gates P0/P1/P2.

## 2026-06-09 — Laboratorio antes de conteudo e mudancas P0

Decisao: toda mudanca P0/P1 deve passar por laboratorio de testes, validacao, refinamento e reteste.

Motivos:
- scripts isolados podem passar e ainda deixar falha de contrato;
- conteudo juridico exige fonte correta e escrita natural;
- escala de 10 mil paginas sem laboratorio vira risco de spam, duplicidade e baixa qualidade.

Consequencia: `tools/lab-cycle` passa a ser o comando minimo de laboratorio, e conteudo juridico em escala permanece bloqueado ate pesquisa de fonte, documentacao, revisao e gates.

## 2026-06-09 — Fontes como referencia, nao scraping

Decisao: fontes oficiais e APIs publicas registradas sao referencia/proveniencia, nao alvo de scraping, clone ou espelhamento.

Motivos:
- o portal deve ser unico e editorialmente proprio;
- copiar massa de dados ou estrutura oficial cria risco de spam e baixa qualidade;
- conteudo juridico precisa ser natural, humano e contextualizado.

Consequencia: qualquer uso de fonte oficial exige pesquisa critica, registro de proveniencia e redacao propria. Ingestao automatica permanece bloqueada em P0.
