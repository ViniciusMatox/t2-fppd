package main

import (
	"fmt"
	"net/rpc"
	"os"
	"time"
)

// Evento representa entrada do usuário
type Evento struct {
	Tipo  string // "mover" ou "sair"
	Tecla rune   // 'w', 'a', 's', 'd'
}

// Jogoh representa estado local do cliente
type Jogoh struct {
	PosX      int
	PosY      int
	StatusMsg string
}

// Função principal do cliente
func executarCliente(clienteID string) {
	conn, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Inicia o jogo no terminal
interfaceIniciar()
defer interfaceFinalizar()

// Carrega o mapa localmente apenas para obter layout visual
mapaFile := "mapa.txt"
if len(os.Args) > 1 {
	mapaFile = os.Args[1]
}

jogo := jogoNovo()
if _, err := os.Stat(mapaFile); err != nil {
	mapaFile = "mapa.txt"
}
if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
	panic(err)
}

// Variáveis do loop RPC
var estado EstadoJogo
seq := 0

for {
	// Lê evento do teclado (w/a/s/d ou sair)
	evento := interfaceLerEventoTeclado()
	if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}

	if evento.Tipo == "sair" {
		break
	}

	if evento.Tipo == "mover" {
		seq++
		cmd := Comando{
			ClienteID:      clienteID,
			Tecla:          evento.Tecla,
			SequenceNumber: seq,
		}

		var resposta string
		err := conn.Call("Servidor.ProcessarComando", cmd, &resposta)
		if err != nil {
			fmt.Println("Erro ao enviar comando:", err)
		}
	}

	// Pega estado atualizado do servidor
	err := conn.Call("Servidor.ObterEstado", clienteID, &estado)
	if err != nil {
		fmt.Println("Erro ao obter estado:", err)
		continue
	}

	// Atualiza o estado do jogo local (posição e status)
	if pos, ok := estado.Posicoes[clienteID]; ok {

		jogo.PosX = pos[0]
		jogo.PosY = pos[1]
		jogo.StatusMsg = estado.Status[clienteID]
	}
	// Desenha o jogo com a posição do personagem
	interfaceDesenharJogo(&jogo)

	time.Sleep(50 * time.Millisecond)
}

}
