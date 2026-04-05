# BarberFlow API 💈 

## Visão do Produto
* **Para:** Donos de pequenas barbearias e seus clientes.
* **Que:** Sofrem com a desorganização de agendamentos via WhatsApp e desistências em cima da hora.
* **O "BarberFlow API" é um:** Sistema de gestão de agendamentos e serviços.
* **Que:** Permite a reserva de horários online com visualização de disponibilidade em tempo real e gestão de catálogo de serviços.
* **Diferente de:** Agendas de papel ou conversas soltas no chat.
* **Nosso produto:** Oferece uma integração simplificada e uma API performática feita em Go, focada apenas na jornada de reserva.

---

## Product Backlog 

| Prioridade | User Story | Critérios de Aceitação | Sprint |
| :--- | :--- | :--- | :--- |
| **P1** | Como cliente, quero me cadastrar... | Email único, senha criptografada (bcrypt). | 1 |
| **P1** | Como barbeiro, quero listar serviços.. | Nome, preço e duração. Retorno em JSON. | 1 |
| **P2** | Como cliente, quero agendar um serviço... | Impedir agendamento em horário já ocupado. | 2 |
| **P2** | Como admin, quero ver a agenda do dia... | Filtro por data e ID do profissional. | 3 |

---

## Stack Tecnológico

* **Go:** Escolhido pela alta performance e concorrência nativa, sendo ideal para sistemas de agendamento com muitos acessos.
* **Chi:** Um roteador leve, idiomático e compatível com `net/http`, o que mantém a base de código limpa.
* **sqlc:** Ferramenta que gera código Go a partir de SQL puro, garantindo segurança de tipos (type-safety) e evitando erros comuns de ORMs pesados.
* **PostgreSQL:** Banco de dados relacional robusto para garantir a consistência dos dados de agendamento através das propriedades ACID.
