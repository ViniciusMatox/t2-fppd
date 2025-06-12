// servidor.go
package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// Estrutura para representar comandos enviados pelos clientes
type Comando struct {
	ClienteID       string
	Tecla           rune
	SequenceNumber  int
}

// Estrutura do estado do jogo que será enviado aos clientes
type EstadoJogo struct {
	Posicoes map[string][2]int
	Status   map[string]string
}

// Estrutura principal do servidor, contendo o estado compartilhado
type Servidor struct {
	posicoes           map[string][2]int          // Posições dos jogadores
	comandosProcessados map[string]map[int]bool   // Controle de execução única
}

// Construtor do servidor
func novoServidor() *Servidor {
	return &Servidor{
		posicoes:           make(map[string][2]int),
		comandosProcessados: make(map[string]map[int]bool),
	}
}

// Processa um comando enviado por um cliente
func (s *Servidor) ProcessarComando(cmd Comando, resposta *string) error {
	fmt.Printf("Recebido comando do cliente %s: tecla=%c, seq=%d\n", cmd.ClienteID, cmd.Tecla, cmd.SequenceNumber)
	// Inicializa o controle de sequenceNumber do cliente, se necessário
	if _, ok := s.comandosProcessados[cmd.ClienteID]; !ok {
		s.comandosProcessados[cmd.ClienteID] = make(map[int]bool)
	}

	// Garante execução única (exactly-once)
	if s.comandosProcessados[cmd.ClienteID][cmd.SequenceNumber] {
		*resposta = "Comando duplicado ignorado"
		return nil
	}

	// Marca o comando como processado
	s.comandosProcessados[cmd.ClienteID][cmd.SequenceNumber] = true

	// Atualiza a posição do jogador
	pos := s.posicoes[cmd.ClienteID]
	switch cmd.Tecla {
	case 'w':
		pos[1]--
	case 's':
		pos[1]++
	case 'a':
		pos[0]--
	case 'd':
		pos[0]++
	}
	s.posicoes[cmd.ClienteID] = pos
	*resposta = fmt.Sprintf("Movido para (%d, %d)", pos[0], pos[1])
	return nil
}

// Fornece o estado atual do jogo ao cliente
func (s *Servidor) ObterEstado(clienteID string, estado *EstadoJogo) error {
	fmt.Printf("Estado solicitado por %s\n", clienteID)

	// Se é a primeira vez que esse cliente se conecta, define posição inicial
	if _, ok := s.posicoes[clienteID]; !ok {
		s.posicoes[clienteID] = [2]int{5, 5} // posição inicial segura
		fmt.Printf("Cliente %s iniciado na posição (5,5)\n", clienteID)
	}

	// Preenche o estado retornado ao cliente
	estado.Posicoes = s.posicoes
	estado.Status = map[string]string{
		clienteID: "Estado sincronizado com sucesso",
	}
	return nil
}


// Inicia o servidor RPC
func iniciarServidor() {
	srv := novoServidor()
	rpc.Register(srv)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Servidor iniciado na porta 1234")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// main chama iniciarServidor()