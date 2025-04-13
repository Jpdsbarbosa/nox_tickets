# NOX Tickets

Sistema de gerenciamento de tickets desenvolvido em Go.

## Descrição
NOX Tickets é um sistema de gerenciamento de tickets que permite o controle e acompanhamento de demandas internas da equipe. O sistema oferece funcionalidades como criação de tickets, atribuição de responsáveis, acompanhamento de status e notificações.

## Funcionalidades (Fase 1)
- Criação, leitura, atualização e exclusão de tickets (CRUD)
- Sistema de autenticação de usuários
- Comentários em tickets
- Histórico de alterações
- Status básicos (Aberto, Em Andamento, Resolvido)

## Tecnologias Utilizadas
- Go
- Echo (Framework Web)
- PostgreSQL (Banco de dados)
- JWT (Autenticação)

## Estrutura do Projeto
```
project/
├── cmd/
│   └── server/         # Ponto de entrada da aplicação
├── internal/
│   ├── domain/         # Regras de negócio e entidades
│   ├── infrastructure/ # Implementações de infraestrutura
│   ├── application/    # Casos de uso da aplicação
│   ├── interfaces/     # Interfaces HTTP e mensageria
│   └── shared/         # Código compartilhado
```

## Como Executar
1. Clone o repositório
2. Configure as variáveis de ambiente
3. Execute `go run cmd/server/main.go`

## Próximos Passos
- Implementação de notificações
- Integração com Google Chat
- Sistema de upload de arquivos
- Dashboard e relatórios 