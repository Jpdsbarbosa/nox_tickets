# NOX Tickets

Sistema de gerenciamento de tickets desenvolvido em Go.

## Descrição
NOX Tickets é um sistema de gerenciamento de tickets que permite o controle e acompanhamento de demandas internas da equipe. O sistema oferece funcionalidades como criação de tickets, atribuição de responsáveis, acompanhamento de status e notificações.

## Arquitetura do Sistema

Para entender melhor nossa arquitetura, imagine um restaurante:

### 🏪 Restaurante = Nossa Aplicação
- Assim como um restaurante tem diferentes áreas e regras de funcionamento, nossa aplicação é organizada em camadas com responsabilidades específicas.

### 👨‍💼 Garçons = Use Cases (Casos de Uso)
- São os coordenadores das operações
- Recebem os pedidos (inputs)
- Coordenam entre cozinha e despensa
- Entregam os resultados (outputs)

Exemplos:
- `CriarTicketUseCase` = Garçom que pega novos pedidos
- `BuscarTicketUseCase` = Garçom que verifica status dos pedidos
- `ListarTicketsUseCase` = Garçom que mostra o cardápio/pedidos do dia
- `AtualizarStatusUseCase` = Garçom que atualiza situação dos pedidos

### 🧑‍🍳 Cozinha = Domain (Regras de Negócio)
- Contém as regras de como as coisas devem ser feitas
- Define o que pode ou não pode ser feito
- Tem suas próprias regras e validações
- É o coração do sistema

### 🗄️ Despensa = Repository (Banco de Dados)
- Onde os dados são armazenados
- Os garçons (use cases) precisam passar aqui para pegar ou guardar informações
- Interface com o banco de dados

### 📝 Comanda = Input/Output (Dados)
- Input = O pedido do cliente (dados de entrada)
- Output = A confirmação do pedido (resultado da operação)

### Fluxo de uma Operação
1. Cliente faz pedido (Input)
2. Garçom (Use Case) recebe
3. Garçom consulta despensa (Repository)
4. Garçom leva para cozinha (Domain)
5. Cozinha processa seguindo as regras
6. Garçom guarda resultado na despensa (Repository)
7. Garçom retorna ao cliente (Output)

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