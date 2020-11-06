# Nodes

### Tipos 
- Owners  
    - Terminal
    - With Request

- Idle Nodes

- Waiters 
    - Terminal
    - With Request


### Atributos Comuns 
- MyChan - Channel do Node.
- find - Channel onde Recebe Pedidos.

---

# Channels

### Tipos

- Access Request - #1
- Give Access - #2

### O que é comunicado

Access Request:
    - Um Channel que identifica quem está a fazer o Request. #2
    - Um Channel que identifica a nova ligação. #1

Give Access:
    - O acesso ao Objeto.

---

# Comportamento dos Nodes

## Owner Terminal:
1. Recebe um pedido - Transforma-se em **Owner with Request** e atualiza o **Link**.

#### Comportamento 1
1. Recebe no **find**:
    - Channel de quem fez o pedido de acesso (WaiterChan).
    - Channel de quem fez chegar o pedido de acesso (NewLink).

2. Transforma-se em **Owner with Request**

`OwnerWithRequest(find, MyChan, WaiterChan, NewLink, Obj)`


#### Atributos do Owner Terminal:
- **find** - Channel onde recebe pedidos.
- **MyChan** - Channel do Node.
- **Obj** - o Objecto. 


---
## Owner With Request:
1. Chegada de pedido - Reorganiza as ligações, reencaminhando o pedido para o parent **Node** inicial. 
2. Cedência do Objeto - Transforma-se em **Idle** e envia mensagem para o proximo da fila.

#### Comportamento 1

1. Recebe no **find**:
    - Channel de quem fez o pedido de acesso (WaiterChan)
    - Channel de quem fez chegar o pedido de acesso (NewLink)

2. Envia no **Link**:
    - WaiterChan
    - find

3. Continua a ser **Owner With Request** mas atualiza o **Link** para **NewLink**

`OwnerWithRequest(find, MyChan, Obj, NewLink, WaiterChan)`

#### Comportamento 2

1. Envia no **WaiterChan**.
    - **Obj**
2. Transforma-se em **Idle**.

`Idle(find, MyChan, Link)`

#### Atributos do Owner With Request:
- **find** - Channel onde recebe pedidos.
- **MyChan** - Channel do Node.
- **Obj** - o Objecto. 
- **Link** - Ligação para o child **Node**.
- **WaiterChan** - Channel de quem fez o pedido de acesso.

---

## Idle: 
1.  Chegada de pedido - Reorganiza as ligações, reencaminhando o pedido para o parent Node inicial.
2. O próprio **Node** decide pedir o **Obj** - Muda para Waiter Terminal e avisa a rede do seu pedido.


#### Comportamento 1

1. Recebe no **find**:
    - Channel de quem fez o pedido de acesso (WaiterChan)
    - Channel de quem fez chegar o pedido de acesso (NewLink)

2. Envia no **Link**:
    - **WaiterChan**
    - **find**

3. Continua a ser **Idle** mas atualiza o **Link** para **NewLink**.

`Idle(find, MyChan, NewLink)`

#### Comportamento 2
1. Envia no **Link**:
    - **MyChan** 
    - **find**
2. Mudança de **Idle** para **Waiter Terminal**
    (MyChan dá acesso ao Obj)

`WaiterTerminal(find, MyChan)`

#### Atributos do Idle Node:
- **find** - Channel onde recebe pedidos.
- **MyChan** - Channel do Node.
- **Link** - Ligação para o child **Node**.


-------
## Waiter Terminal:
1. Recebe um pedido - Transforma-se em Waiter with request e atualiza a ligação.
2. Recebe o acesso ao Objeto.


#### Comportamento 1 
1. Recebe no **find**:
    - Channel de quem fez o pedido de acesso (WaiterChan)
    - Channel de quem fez chegar o pedido de acesso (NewLink)

2. Transforma-se em **WaiterWithRequest**.

`WaiterWithRequest(find, MyChan, NewLink, WaiterChan)`

#### Comportamento 2 
1. Recebe no **MyChan**:
    - **Obj**

2. Transforma-se em **Owner Terminal**.

`OwnerTerminal(find, MyChan, Obj)`

#### Atributos do Waiter Terminal:
- **find** - Channel onde recebe pedidos.
- **MyChan** - Channel do Node.

------
## Waiter with Request:

1. Recebe o acesso ao Objeto. 
2. Chegada de pedido - Reorganiza as ligações, reencaminhando o pedido para o parent Node inicial.


#### Comportamento 1
1. Receber no **MyChan**:
    - **Obj**

2. Transforma-se em **Owner With Request**.

`OwnerWithRequest(find, MyChan, Obj, Link, WaiterChan)`

#### Comportamento 2

1. Recebe no **find**:
    - Channel de quem fez o pedido de acesso (WaiterChan)
    - Channel de quem fez chegar o pedido de acesso (NewLink)

2. Envia no **Link**:
    - **WaiterChan**
    - **find**

3. Continua a ser **Waiter With Request** mas atualiza o **Link** para **NewLink**.

`WaiterWithRequest(find, MyChan, Obj, NewLink, WaiterChan)`


#### Atributos do Waiter Terminal:
- **find** - Channel onde recebe pedidos.
- **MyChan** - Channel do Node.
- **Obj** - o Objecto. 
- **Link** - Ligação para o child **Node**.
- **WaiterChan** - Channel de quem fez o pedido de acesso.
