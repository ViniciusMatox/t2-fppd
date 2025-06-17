package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "sync"
)

type Player struct {
    ID   string
    Conn net.Conn
    X, Y int
}

var (
    players   = make(map[string]*Player)
    playersMu sync.Mutex
)

func iniciarServidor() {
    ln, err := net.Listen("tcp", "0.0.0.0:1234")
    if err != nil {
        fmt.Println("Erro ao iniciar servidor:", err)
        return
    }
    fmt.Println("Servidor rodando em 0.0.0.0:1234")

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Erro ao aceitar conexão:", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)
    idLine, _ := reader.ReadString('\n')
    id := strings.TrimSpace(idLine)

    playersMu.Lock()
    players[id] = &Player{ID: id, Conn: conn, X: 5, Y: len(players)*2 + 1}
    playersMu.Unlock()

    fmt.Printf("[+] Novo jogador conectado: %s\n", id)
    fmt.Fprintf(conn, "ID registrado: %s\n", id)
    broadcast(fmt.Sprintf("NOVO %s %d %d\n", id, players[id].X, players[id].Y), id)

    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            break
        }

        msg = strings.TrimSpace(msg)

        if strings.HasPrefix(msg, "MOVE ") {
            direcao := strings.TrimPrefix(msg, "MOVE ")
            logTecla(id, direcao)
            movePlayer(id, direcao)

            playersMu.Lock()
            pos := players[id]
            playersMu.Unlock()

            broadcast(fmt.Sprintf("POS %s %d %d\n", id, pos.X, pos.Y), "")
        }
    }

    playersMu.Lock()
    delete(players, id)
    playersMu.Unlock()
    fmt.Printf("[-] Jogador desconectado: %s\n", id)
    broadcast(fmt.Sprintf("SAIU %s\n", id), id)
}

func movePlayer(id, dir string) {
    playersMu.Lock()
    defer playersMu.Unlock()

    if p, ok := players[id]; ok {
        switch dir {
        case "UP":
            p.Y--
        case "DOWN":
            p.Y++
        case "LEFT":
            p.X--
        case "RIGHT":
            p.X++
        }
    }
}

func broadcast(msg, except string) {
    playersMu.Lock()
    defer playersMu.Unlock()

    for id, p := range players {
        if id != except {
            fmt.Fprint(p.Conn, msg)
        }
    }
}

func logTecla(id, direcao string) {
    tecla := map[string]string{
        "UP":    "W",
        "DOWN":  "S",
        "LEFT":  "A",
        "RIGHT": "D",
    }[direcao]

    fmt.Printf("[%s] pressionou: %s (%s)\n", id, tecla, "MOVE "+direcao)
}
