// cliente.go
package main

import (
	"net/rpc"   // Biblioteca de chamadas remotas (Remote Procedure Call)
	"time"      // Usado para controlar pausas entre iterações (sleep)
)

// Função principal que representa a lógica do cliente do jogo
func executarCliente(clienteID string) {
	// Estabelece conexão com o servidor RPC rodando localmente na porta 1234
	conn, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err) // Encerra o programa caso a conexão falhe
	}
	defer conn.Close() // Fecha a conexão RPC ao final da função

	// Inicializa a interface gráfica/terminal (provavelmente usando termbox ou similar)
	interfaceIniciar()
	defer interfaceFinalizar() // Garante que a interface será finalizada corretamente no fim

	// Cria um mapa fictício apenas para teste inicial do jogo
	var jogo Jogo = criarMapaMockado()

	// Variável que armazenará o estado atual do jogo recebido do servidor
	var estado EstadoJogo

	// Número de sequência dos comandos enviados pelo cliente (usado para garantir ordem)
	seq := 0

	// Loop principal do cliente: lê entrada do usuário, envia comandos e atualiza o jogo
	for {
		// Lê evento de teclado (input do jogador)
		ev := interfaceLerEventoTeclado()

		// Se o tipo do evento for "sair", finaliza o loop e encerra o jogo
		if ev.Tipo == "sair" {
			break
		}

		// Se o jogador pressionou uma tecla de movimento:
		if ev.Tipo == "mover" {
			seq++ // Incrementa o número de sequência do comando

			// Cria um comando representando o movimento
			cmd := Comando{
				ClienteID:      clienteID, // ID único do cliente
				Tipo:           "mover",   // Tipo do comando
				Tecla:          ev.Tecla,  // Tecla pressionada (w, a, s, d)
				SequenceNumber: seq,       // Número de sequência
			}

			// Envia o comando para o servidor e recebe uma string de resposta (ignorada aqui)
			var resposta string
			conn.Call("Servidor.ProcessarComando", cmd, &resposta)
		}

		// Solicita ao servidor o estado atual do jogo para este cliente
		conn.Call("Servidor.ObterEstado", clienteID, &estado)

		// Atualiza a posição do jogador com base na resposta do servidor
		pos := estado.Posicoes[clienteID]
		jogo.PosX = pos[0]
		jogo.PosY = pos[1]

		// Atualiza a mensagem de status do jogador
		jogo.StatusMsg = estado.Status[clienteID]

		// Desenha o jogo no terminal/interface gráfica
		interfaceDesenharJogo(&jogo)

		// Aguarda 50 milissegundos antes de repetir o loop (limita a taxa de atualização)
		time.Sleep(50 * time.Millisecond)
	}
}

// Função temporária que deveria retornar um mapa do jogo para testes
// Atualmente apenas lança um pânico para lembrar que precisa ser implementada
func criarMapaMockado() Jogo {
	panic("unimplemented") // Gera um erro se chamada — precisa ser implementada
}
