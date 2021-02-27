# Owner Terminal -> Owner With Request
Quando o Node recebe um pedido.
- Definição da Queue principal, em que a cabeça da lista é o Node que realizou o pedido que chegou.

# Waiter With Request -> Owner With Request
- Remoção do Node da sua Queue

# Waiter Terminal -> Owner Terminal
Quando um Node recebe o acesso ao objeto.
- Eliminação da Queue que contém o Node que se transformou (Head).

# Waiter Terminal -> Waiter With Request
Quando o Node recebe um pedido.
- Concatenação das duas Queues, a primeira sendo a Queue do Node, a segunda sendo a Queue do Node que realizou o pedido.

#  Idle -> WaiterTerminal
Quando o Node decide realizar um pedido.
- Forma-se uma Queue de um único elemento.
