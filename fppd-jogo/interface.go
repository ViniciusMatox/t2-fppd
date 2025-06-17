package main

import (
    "fmt"
    "github.com/nsf/termbox-go"
    "strings"
    "sync"
    "time"
)

// Define um tipo Cor compatível com termbox
type Cor = termbox.Attribute

// Cores utilizadas
const (
    CorPadrao         = termbox.ColorDefault
    CorCinzaEscuro    = termbox.ColorWhite
    CorVermelho       = termbox.ColorRed
    CorVerde          = termbox.ColorGreen
    CorParede         = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
    CorFundoParede    = termbox.ColorDarkGray
    CorTexto          = termbox.ColorDarkGray
)

// EventoTeclado representa uma ação detectada do teclado
type EventoTeclado struct {
    Tipo  string
    Tecla rune
}

func interfaceCapturarEvento() EventoTeclado {
    ev := termbox.PollEvent()
    if ev.Type == termbox.EventKey {
        switch ev.Key {
        case termbox.KeyEsc:
            return EventoTeclado{Tipo: "sair"}
        case termbox.KeySpace:
            return EventoTeclado{Tipo: "interagir"}
        default:
            return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
        }
    }
    return EventoTeclado{Tipo: ""}
}

// Estrutura de jogador para multiplayer
type Jogador struct {
    ID string
    X  int
    Y  int
}

var (
    jogadores   = make(map[string]Jogador)
    jogadoresMu sync.Mutex
    meuID       string
    jogoLocal   Jogo
)

// iniciarInterface carrega o mapa local e exibe o jogo
func iniciarInterface(id string, mensagens <-chan string) {
    meuID = id
    _ = jogoCarregarMapa("mapa.txt", &jogoLocal)

    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    go func() {
        for msg := range mensagens {
            processarMensagem(msg)
        }
    }()

    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {
        desenharTela()
    }
}

// processarMensagem trata mensagens recebidas do servidor
func processarMensagem(msg string) {
    partes := strings.Fields(msg)
    if len(partes) == 0 {
        return
    }

    jogadoresMu.Lock()
    defer jogadoresMu.Unlock()

    switch partes[0] {
    case "NOVO", "POS":
        if len(partes) == 4 {
            id := partes[1]
            x := atoi(partes[2])
            y := atoi(partes[3])
            jogadores[id] = Jogador{ID: id, X: x, Y: y}
        }
    case "SAIU":
        if len(partes) == 2 {
            delete(jogadores, partes[1])
        }
    }
}

// desenharTela exibe o mapa e os jogadores no terminal
func desenharTela() {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    // desenha o mapa do jogo
    for y, linha := range jogoLocal.Mapa {
        for x, elem := range linha {
            termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
        }
    }

    // desenha os jogadores por cima do mapa
    jogadoresMu.Lock()
    for _, j := range jogadores {
        cor := termbox.ColorRed
        if j.ID == meuID {
            cor = termbox.ColorGreen
        }
        termbox.SetCell(j.X, j.Y, '\u263A', cor, termbox.ColorDefault)
    }
    jogadoresMu.Unlock()

    termbox.Flush()
}

// atoi converte string para int (sem erro)
func atoi(s string) int {
    var i int
    fmt.Sscanf(s, "%d", &i)
    return i
}
