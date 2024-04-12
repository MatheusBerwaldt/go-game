package main

import (
    "bufio"
    "github.com/nsf/termbox-go"
    "os"
    "fmt"
    "time"
    "math/rand"
)

type Elemento struct {
    simbolo rune
    cor termbox.Attribute
    corFundo termbox.Attribute
    tangivel bool
}
	
var personagem = Elemento{
	simbolo: '☺',
	cor: termbox.ColorWhite,
	corFundo: termbox.ColorDefault,
	tangivel: true,
}

var parede = Elemento{
	simbolo: '▤',
	cor: termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim,
	corFundo: termbox.ColorDarkGray,
	tangivel: true,
}

var vazio = Elemento{
    simbolo: ' ',
    cor: termbox.ColorDefault,
    corFundo: termbox.ColorDefault,
    tangivel: false,
}

var inimigo = Elemento{
    simbolo: 'Ω',
    cor: termbox.ColorRed,
    corFundo: termbox.ColorDefault,
    tangivel: true,
}

var estrela = Elemento{
    simbolo: '★',
    cor: termbox.ColorYellow,
    corFundo: termbox.ColorDefault,
    tangivel: false,
}

var portal = Elemento{
    simbolo: '0',
    cor: termbox.ColorGreen,
    corFundo: termbox.ColorDefault,
    tangivel: false,
}

var mapa [][]Elemento

var statusMsg string

var posX, posY int
var posXinimigo, posYinimigo int
var posXestrela, posYestrela int

func main(){
	err := termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

	carregarMapa("mapa.txt")
	desenhaTudo()

    go moverInimigo()
    go moverEstrela()

	for {
        event := termbox.PollEvent();
		
		if event.Type == termbox.EventKey {
			switch event.Key {
			case termbox.KeyEsc:
				return
			case 'e':
				interagir()
			default:
				mover(event.Ch)
			}

            desenhaTudo()
        }
    }
}

func carregarMapa(nomeArquivo string) {
    arquivo, err := os.Open(nomeArquivo)

    if err != nil {
        panic(err)
    }
    defer arquivo.Close()

    scanner := bufio.NewScanner(arquivo)
    y := 0

    for scanner.Scan() {
        linhaTexto := scanner.Text()
        var linhaElementos []Elemento
        for x, char := range linhaTexto {
            elementoAtual := vazio
            switch char {
            case parede.simbolo:
                elementoAtual = parede
            case personagem.simbolo:
                posX, posY = x, y
                elementoAtual = personagem
			case inimigo.simbolo:
                posXinimigo, posYinimigo = x, y
                elementoAtual = inimigo
            case estrela.simbolo:
                posXestrela, posYestrela = x, y
                elementoAtual = estrela
            case portal.simbolo:
                elementoAtual = portal
            }
            linhaElementos = append(linhaElementos, elementoAtual)
        }
        mapa = append(mapa, linhaElementos)
        y++
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }
}

func desenhaTudo() {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    for y, linha := range mapa {
        for x, elem := range linha {
			termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
        }
    }

    desenhaBarraDeStatus()

    termbox.Flush()
}

func desenhaBarraDeStatus() {
    for i, c := range statusMsg {
        termbox.SetCell(i, len(mapa)+1, c, termbox.ColorBlack, termbox.ColorDefault)
    }
    msg := "Use WASD para mover e E para interagir. ESC para sair."
    for i, c := range msg {
        termbox.SetCell(i, len(mapa)+3, c, termbox.ColorBlack, termbox.ColorDefault)
    }
}

func mover(comando rune) {
    dx, dy := 0, 0
    switch comando {
    case 'w':
        dy = -1
    case 'a':
        dx = -1
    case 's':
        dy = 1
    case 'd':
        dx = 1
    }
    novaPosX, novaPosY := posX+dx, posY+dy

    // Fora dos limites
    if dentroDosLimites(novaPosX, novaPosY) == false {
        return
    }

    // Conflito
    if mapa[novaPosY][novaPosX].tangivel {
        if mapa[novaPosY][novaPosX].simbolo == inimigo.simbolo {
            // morrer()
        }
        return
    }

    mapa[posY][posX] = vazio
    posX, posY = novaPosX, novaPosY
    mapa[posY][posX] = personagem
}

func interagir() {
    statusMsg = fmt.Sprintf("Interagindo em (%d, %d)", posX, posY)
}

func moverInimigo() {
    for {
        time.Sleep(time.Duration(500) * time.Millisecond)

        direcoes := []int{-1, 0, 1}

        dirX := direcoes[rand.Intn(len(direcoes))]
        dirY := direcoes[rand.Intn(len(direcoes))]

        novaPosX := posXinimigo + dirX
        novaPosY := posYinimigo + dirY

        if dentroDosLimites(novaPosX, novaPosY) == false {
            return
        }


        if mapa[novaPosY][novaPosX].tangivel {
            if mapa[novaPosY][novaPosX].simbolo == inimigo.simbolo {
                //morrer()
            }
        }

        mapa[posYinimigo][posXinimigo] = vazio
        posXinimigo = novaPosX
        posYinimigo = novaPosY
        mapa[posYinimigo][posXinimigo] = inimigo

        desenhaTudo()
    }
}

func moverEstrela() {

    for {
        time.Sleep(time.Duration(10000) * time.Millisecond)

        optionsY := make([]int, l)
        optionsX := make([]int, len(mapa[0]))

        dirY := rand.Intn(len(mapa))
        dirX := rand.Intn(len(mapa[0]))

        novaPosX := posXinimigo + dirX
        novaPosY := posYinimigo + dirY

        if dentroDosLimites(novaPosX, novaPosY) == false {
            return
        }


        if mapa[novaPosY][novaPosX].tangivel {
            if mapa[novaPosY][novaPosX].simbolo == inimigo.simbolo {
                //morrer()
            }
        }

        mapa[posYinimigo][posXinimigo] = vazio
        posXinimigo = novaPosX
        posYinimigo = novaPosY
        mapa[posYinimigo][posXinimigo] = inimigo

        desenhaTudo()

    }
}

func dentroDosLimites(x int, y int) bool {
    return y >= 0 && y < len(mapa) && x >= 0 && x < len(mapa)
}