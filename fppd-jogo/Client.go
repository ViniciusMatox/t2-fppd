package main

import (
    "bufio"
    "fmt"
    "net"
    "github.com/nsf/termbox-go"
)

func executarCliente(id string) {
    conn, err := net.Dial("tcp", "127.0.0.1:1234")
    if err != nil {
        fmt.Println("Erro ao conectar:", err)
        return
    }
    defer conn.Close()

    fmt.Fprintln(conn, id)

    mensagens := make(chan string, 100)
    go receberAtualizacoes(conn, mensagens)
    go iniciarInterface(id, mensagens)

    // Inicia captura de teclas em tempo real com termbox
    err = termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    for {
        ev := termbox.PollEvent()
        if ev.Type == termbox.EventKey {
            switch ev.Ch {
            case 'w', 'W':
                fmt.Fprintln(conn, "MOVE UP")
            case 'a', 'A':
                fmt.Fprintln(conn, "MOVE LEFT")
            case 's', 'S':
                fmt.Fprintln(conn, "MOVE DOWN")
            case 'd', 'D':
                fmt.Fprintln(conn, "MOVE RIGHT")
            case 'q', 'Q':
                return
            }
        }
    }
}

func receberAtualizacoes(conn net.Conn, mensagens chan<- string) {
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        mensagens <- scanner.Text()
    }
}
