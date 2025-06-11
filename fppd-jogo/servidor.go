// servidor.go
package main

import (
	"fmt"       // Usado para imprimir mensagens no terminal
	"net"       // Usado para criar o listener de conexão TCP
	"net/rpc"   // Biblioteca para implementar RPC (Remote Procedure Call)
)

// Struct que representa o servidor e armazena o estado do jogo
type Servidor struct {
	posicoes map[string][2]int // Mapa que associa cada jogador (por ID) à sua posição (X, Y)
}

// Função construtora para criar um novo servidor com o mapa de posições inicializado
func novoServidor() *Servidor {
	return &Servidor{
		posicoes: make(map[string][2]int), // Inicializa o mapa de posições
	}
}

// Método chamado remotamente via RPC para processar comandos enviados pelos clientes
func (s *Servidor) ProcessarComando(cmd Comando, resposta *string) error {
	// Obtém a posição atual do jogador
	pos := s.posicoes[cmd.ClienteID]

	// Atualiza a posição de acordo com a tecla pressionada
	switch cmd.Tecla {
	case 'w': // mover para cima
		pos[1]--
	case 's': // mover para baixo
		pos[1]++
	case 'a': // mover para a esquerda
		pos[0]--
	case 'd': // mover para a direita
		pos[0]++
	}

	// Atualiza a posição do jogador no mapa
	s.posicoes[cmd.ClienteID] = pos

	// Resposta informando a nova posição do jogador
	*resposta = fmt.Sprintf("Movido para (%d, %d)", pos[0], pos[1])
	return nil // Nenhum erro ocorreu
}

// Método chamado remotamente via RPC para fornecer o estado atual do jogo ao cliente
func (s *Servidor) ObterEstado(clienteID string, estado *EstadoJogo) error {
	// Envia as posições de todos os jogadores
	estado.Posicoes = s.posicoes

	// Mensagem de status para o cliente específico
	estado.Status = map[string]string{
		clienteID: "Status do jogador atualizado",
	}
	return nil // Nenhum erro ocorreu
}

// Função que inicia o servidor e escuta conexões RPC
func iniciarServidor() {
	srv := novoServidor() // Cria uma nova instância do servidor

	// Registra o servidor para permitir chamadas RPC nos seus métodos
	rpc.Register(srv)

	// Inicia o listener TCP na porta 1234
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err) // Se não conseguir escutar na porta, o programa para com erro
	}
	defer listener.Close() // Garante que o listener será fechado ao final

	fmt.Println("Servidor iniciado na porta 1234")

	// Loop principal do servidor: aceita conexões dos clientes
	for {
		conn, err := listener.Accept() // Aceita uma nova conexão
		if err != nil {
			continue // Se houver erro, ignora e tenta de novo
		}

		// Inicia uma nova goroutine para tratar a conexão RPC com o cliente
		go rpc.ServeConn(conn)
	}
}
