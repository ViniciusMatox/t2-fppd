
package main

import (
	
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "servidor" {
		iniciarServidor()
		return
	}
		executarCliente("Player1") // Inicia o cliente com ID "Player1"
		// Inicia o cliente do jog

	// Inicializa a interface (termbox)
	
}