# LAB_VALIDATION.md

## Regra

O projeto deve operar como laboratorio: testar, validar, refinar, testar novamente e somente entao registrar checkpoint. Nada de chute.

## Comando principal

`./tools/lab-cycle`

Esse comando combina:
- `go test -count=1 ./...`;
- `./tools/check-all`;
- `./tools/check-sources`;
- `go run ./cmd/build public`;
- `go list -m all`;
- `git diff --check`;
- busca por residuos `.py` e `.pyc`.

O gate `./tools/check-performance-budget` deve reprovar HTML publico pesado, `<script>`, runtime cliente, bundle JavaScript, WebAssembly, mapas, `modulepreload`, import map, marcadores de hidratacao e CSS inline excessivo. Leveza e parte da prova de indexacao para Googlebot e bots valiosos.

## Politica

Nenhum script isolado e prova suficiente para mudanca P0/P1. Use o laboratorio como conjunto minimo e inspecione os artefatos quando a mudanca afetar HTML, sitemap, robots, canonical, indexacao, qualidade ou conteudo.

Validar nao basta. Todo ciclo deve revisar o diff, os artefatos gerados, os riscos e o contrato antes de commitar. A revisao deve procurar atalhos, spam, conteudo mecanico, regressao de P0 e dependencia indevida.

Se qualquer validacao falhar:
1. identificar causa raiz;
2. corrigir;
3. rodar novamente;
4. registrar falha e correcao no checkpoint;
5. commitar checkpoint e artefatos do ciclo;
6. continuar o proximo ciclo quando nao houver bloqueio P0 real.

## Continuidade

O agente nao deve parar ao registrar checkpoint. Todo ciclo deve terminar com proximo passo acionavel, e esse proximo passo deve ser executado em seguida quando nao houver bloqueio P0 real comprovado.

## Commit por ciclo

Depois de validacoes relevantes passarem, fazer commit do checkpoint e dos artefatos do ciclo. O projeto nao deve depender de chat, contexto compactado ou memoria externa para continuidade.

Antes do commit, executar `date`, conferir a hora local e registrar o ciclo numerado no checkpoint. Isso mantem a ordem de auditoria mesmo apos compactacao de contexto.
