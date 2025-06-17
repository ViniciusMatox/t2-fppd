package main

import (
    "os"
)

func main() {
    if len(os.Args) > 1 && os.Args[1] == "servidor" {
        iniciarServidor()
    } else if len(os.Args) > 1 {
        executarCliente(os.Args[1]) // passa apenas o ID
    } else {
        println("Uso:")
        println("  go run . servidor")
        println("  go run . <SeuID>")
    }
}
