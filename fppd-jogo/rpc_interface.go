// rpc_interface.go
package main

// Estrutura que representa um comando enviado de um cliente para o servidor
type Comando struct {
	ClienteID      string // Identificador único do cliente (por exemplo, nome ou UUID)
	SequenceNumber int    // Número de sequência do comando (para controle de ordem e evitar duplicação)
	Tipo           string // Tipo do comando, como "mover" ou "interagir"
	Tecla          rune   // Tecla pressionada pelo jogador (ex: 'w', 'a', 's', 'd')
}

// Estrutura que representa o estado do jogo retornado pelo servidor para o cliente
type EstadoJogo struct {
	Posicoes map[string][2]int // Mapa com a posição (x, y) de cada cliente no jogo
	Status   map[string]string // Mapa com mensagens de status personalizadas por cliente
}
