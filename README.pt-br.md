# Chatgraph-Go

[![Go Reference](https://pkg.go.dev/badge/github.com/irissonnlima/chatgraph-go.svg)](https://pkg.go.dev/github.com/irissonnlima/chatgraph-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/irissonnlima/chatgraph-go)](https://goreportcard.com/report/github.com/irissonnlima/chatgraph-go)

Um framework leve e flexível para criação de chatbots em Go, com fluxo de conversação baseado em rotas, tratamento de timeout e proteção contra loops.

[English](README.md) | Português

## Funcionalidades

- 🚀 **API Simples** - Um único import, registro de rotas intuitivo
- 🔄 **Fluxo Baseado em Rotas** - Defina fluxos de conversação com rotas nomeadas
- ⏱️ **Tratamento de Timeout** - Timeout automático com duração configurável e rotas de fallback
- 🔁 **Proteção contra Loops** - Previne loops de redirecionamento infinitos automaticamente
- 📦 **Observations Genéricas** - Armazene dados de sessão customizados com type safety
- 🔌 **Adaptadores Plugáveis** - RabbitMQ para entrada, REST API para saída (facilmente extensível)
- 📄 **Suporte a Arquivos** - Upload, envio e download de arquivos com deduplicação via SHA256
- ✅ **Bem Testado** - Cobertura de testes abrangente para os pacotes principais

## Instalação

```bash
go get github.com/irissonnlima/chatgraph-go/chat@latest
```

## Início Rápido

```go
package main

import (
    "github.com/irissonnlima/chatgraph-go/chat"
)

// Defina seu tipo de observation para dados de sessão
type Obs struct {
    Contador int `json:"contador"`
}

func main() {
    // Criar adaptadores
    rabbit := chat.NewRabbitMQ[Obs]("user", "pass", "host", "vhost", "queue")
    router := chat.NewRouterApi("http://api-url", "user", "pass")
    engine := chat.NewEngine[Obs]()
    
    // Criar app
    app := chat.NewApp(engine, rabbit, router)
    
    // Registrar rotas
    engine.RegisterRoute("start", func(ctx *chat.Context[Obs]) chat.RouteReturn {
        ctx.SendTextMessage("Olá! Digite algo:")
        r := ctx.NextRoute("echo")
        return &r
    })
    
    engine.RegisterRoute("echo", func(ctx *chat.Context[Obs]) chat.RouteReturn {
        ctx.SendTextMessage("Você disse: " + ctx.Message.EntireText())
        r := ctx.NextRoute("start")
        return &r
    })
    
    // Obrigatório: handlers de timeout e loop
    engine.RegisterRoute("timeout_route", func(ctx *chat.Context[Obs]) chat.RouteReturn {
        ctx.SendTextMessage("Sessão expirada!")
        r := ctx.NextRoute("start")
        return &r
    })
    
    engine.RegisterRoute("loop_route", func(ctx *chat.Context[Obs]) chat.RouteReturn {
        return &chat.RedirectResponse{TargetRoute: "start"}
    })
    
    app.Start()
}
```

## Arquitetura

O Chatgraph segue o padrão de **arquitetura hexagonal**, separando a lógica de domínio dos adaptadores externos:

```
┌─────────────────────────────────────────────────────────────────┐
│                        Aplicação                                │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │
│  │  RabbitMQ   │───▶│  ChatbotApp │───▶│    Router API       │  │
│  │  (Entrada)  │    │  (Serviço)  │    │    (Saída)          │  │
│  └─────────────┘    └──────┬──────┘    └─────────────────────┘  │
│                            │                                    │
│                     ┌──────▼──────┐                             │
│                     │   Rotas     │                             │
│                     │  (Handlers) │                             │
│                     └─────────────┘                             │
└─────────────────────────────────────────────────────────────────┘
```

### Conceitos Principais

#### 1. Rotas

Rotas são estados de conversação nomeados. Cada rota possui uma função handler que:

- Recebe um `Context` com estado do usuário e mensagem
- Envia mensagens para o usuário
- Retorna a próxima ação (próxima rota, redirecionamento, encerrar sessão, etc.)

```go
app.RegisterRoute("saudacao", func(ctx *chat.Context[Obs]) chat.RouteReturn {
    ctx.SendTextMessage("Olá!")
    return ctx.NextRoute("menu")  // Próxima mensagem do usuário vai para a rota "menu"
})
```

#### 2. Retornos de Rota

Handlers podem retornar diferentes ações:

| Tipo de Retorno | Comportamento |
|-----------------|---------------|
| `ctx.NextRoute("nome")` | Define a rota para a **próxima** mensagem do usuário |
| `&RedirectResponse{TargetRoute: "nome"}` | Executa **imediatamente** outra rota |
| `EndAction{ID: "motivo"}` | Encerra a sessão de conversação |
| `TransferToMenu{MenuID: 1}` | Transfere usuário para um menu diferente |
| `nil` | Permanece na rota atual |

**NextRoute vs Redirect:**

```go
// NextRoute: Aguarda entrada do usuário, depois executa "menu"
return ctx.NextRoute("menu")

// Redirect: Executa "menu" imediatamente sem aguardar
return &chat.RedirectResponse{TargetRoute: "menu"}
```

#### 3. Context

O `Context` fornece acesso a:

```go
func handler(ctx *chat.Context[Obs]) chat.RouteReturn {
    // Informações do usuário
    ctx.UserState.User.Name      // Nome do usuário
    ctx.UserState.ChatID         // Identificador do chat
    ctx.UserState.Route          // Histórico de navegação
    
    // Mensagem recebida
    ctx.Message.EntireText()     // Texto completo da mensagem
    ctx.Message.TextMessage      // Texto estruturado (Title, Detail, Footer)
    ctx.Message.Buttons          // Respostas de botões
    ctx.Message.File             // Arquivos anexados
    
    // Observations da sessão (dados customizados)
    obs := ctx.GetObservation()  // Obter observation tipada
    ctx.SetObservation(obs)      // Atualizar observation
    
    // Enviar mensagens
    ctx.SendTextMessage("Olá!")
    ctx.SendMessage(chat.Message{...})
    
    // Operações com arquivos
    ctx.LoadFile("caminho/do/arquivo")      // Upload do disco
    ctx.LoadFileBytes("nome.txt", []byte)   // Upload de bytes
    
    return ctx.NextRoute("proxima")
}
```

#### 4. Tratamento de Timeout

Cada rota possui um timeout configurável. Quando excedido:

1. A execução do handler é **cancelada** via context
2. O usuário é **redirecionado** para a rota de timeout
3. Nenhuma mensagem adicional é enviada do handler que deu timeout

```go
// Padrão: 5 minutos, redireciona para "timeout_route"
app.RegisterRoute("tarefa_lenta", handler)

// Timeout customizado: 30 segundos, redireciona para "timeout_customizado"
app.RegisterRoute("tarefa_rapida", handler, chat.RouterHandlerOptions{
    Timeout: &chat.TimeoutRouteOps{
        Duration: 30 * time.Second,
        Route:    "timeout_customizado",
    },
})
```

**Como funciona internamente:**

```
Mensagem ──▶ Handler Inicia ──▶ [timeout de 5 min]
                  │
                  ├── Handler completa ──▶ Processa resultado
                  │
                  └── Timeout excedido ──▶ Cancela context
                                          └── Redireciona para timeout_route
```

#### 5. Proteção contra Loops

Previne loops de redirecionamento infinitos contando visitas consecutivas à mesma rota:

```go
// Padrão: 3 visitas consecutivas, redireciona para "loop_route"
// Se a rota "A" redireciona para "A" 3 vezes, usuário vai para "loop_route"
```

**Como funciona:**

```
A → A → A → A (4ª vez) ──▶ Redireciona para loop_route
    │   │   │
    └───┴───┴── CurrentRepeated() = 3 > limite
```

#### 6. Mensagens com Botões

Envie mensagens interativas com botões clicáveis:

```go
ctx.SendMessage(chat.Message{
    TextMessage: chat.TextMessage{
        Title:  "Escolha uma opção",
        Detail: "Por favor, selecione:",
    },
    Buttons: []chat.Button{
        {Type: chat.POSTBACK, Title: "Opção A", Detail: "opcao_a"},
        {Type: chat.POSTBACK, Title: "Opção B", Detail: "opcao_b"},
        {Type: chat.URL, Title: "Visitar Site", Detail: "https://exemplo.com"},
    },
})
```

Tipos de botão:

- `POSTBACK`: Envia o valor de `Detail` de volta como mensagem do usuário
- `URL`: Abre a URL no navegador do usuário

#### 7. Observations (Dados de Sessão)

Armazene dados tipados customizados que persistem entre mensagens:

```go
type Obs struct {
    Etapa     int    `json:"etapa"`
    DadosUser string `json:"dados_user"`
}

func handler(ctx *chat.Context[Obs]) chat.RouteReturn {
    obs := ctx.GetObservation()
    obs.Etapa++
    obs.DadosUser = ctx.Message.EntireText()
    ctx.SetObservation(obs)
    
    return ctx.NextRoute("proxima")
}
```

#### 8. Manipulação de Arquivos

Upload e envio de arquivos:

```go
// Upload do disco
file, err := ctx.LoadFile("documento.pdf")
if err == nil && file != nil {
    ctx.SendMessage(chat.Message{File: *file})
}

// Upload de bytes (ex: conteúdo gerado)
conteudo := []byte("Olá, Mundo!")
file, err := ctx.LoadFileBytes("saudacao.txt", conteudo)
```

Arquivos são deduplicados usando hash SHA256 - fazer upload do mesmo conteúdo duas vezes retorna o arquivo em cache.

**Download do conteúdo do arquivo:**

```go
// Baixar bytes do arquivo a partir da URL
if !ctx.Message.File.IsEmpty() {
    bytes, err := ctx.Message.File.Bytes()
    if err == nil {
        // Processar os bytes do arquivo
        fmt.Printf("Baixados %d bytes\n", len(bytes))
    }
}
```

## Configuração

### Opções Padrão

```go
app := chat.NewApp(rabbit, router, chat.RouterHandlerOptions{
    Timeout: &chat.TimeoutRouteOps{
        Duration: 10 * time.Minute,  // Timeout padrão para todas as rotas
        Route:    "timeout_route",
    },
    LoopCount: &chat.LoopCountRouteOps{
        Count: 5,                    // Permite 5 visitas consecutivas à mesma rota
        Route: "loop_route",
    },
})
```

### Opções Por Rota

```go
app.RegisterRoute("sensivel", handler, chat.RouterHandlerOptions{
    Timeout: &chat.TimeoutRouteOps{
        Duration: 1 * time.Minute,
        Route:    "timeout_sensivel",
    },
})
```

## Exemplos

Veja o diretório [examples/](./examples/) para exemplos completos:

- **basic/** - Chatbot simples com observations
- **buttons/** - Demo de botões interativos
- **files/** - Upload e download de arquivos
- **timeout/** - Configuração de timeout customizado

## Estrutura do Projeto

```
chatgraph-go/
├── chat/                # Pacote público unificado da API
│   └── chatgraph.go     # Type aliases e construtores
├── adapters/
│   ├── input/queue/     # Consumidor de mensagens RabbitMQ
│   └── output/router_api/  # Cliente REST API
├── core/
│   ├── domain/          # Modelos de domínio
│   │   ├── action/      # Ações de retorno de rota
│   │   ├── context/     # Contexto do chat
│   │   ├── message/     # Tipos de mensagem
│   │   ├── route/       # Histórico de navegação
│   │   ├── router/      # Opções de handler
│   │   └── user/        # Estado do usuário
│   ├── ports/adapters/  # Interfaces de adaptadores
│   └── service/         # Serviço da aplicação
└── examples/            # Exemplos de uso
```

## Testes

Execute os testes com cobertura:

```bash
go test ./... -cover

# Ou use o script de cobertura para relatório HTML
./coverage.sh
# Abra coverage/coverage.html no seu navegador
```

## Licença

Licença MIT - veja [LICENSE](LICENSE) para detalhes.
