package dzpk

import (
	"net"
	"sync"
)

//it's just a game
const (
	MAX_GROUP_PLAYER int = 10
)

//游戏状态
type GameStatus int

const (
	//创建了， 但是还没开始
	GameNew GameStatus = 0

	//游戏当中
	GameInPlay GameStatus = 1

	//游戏结束
	GameOver GameStatus = 2

	//
	GameToBeDestroy GameStatus = 3
)

//represents a Game
type Game struct {
	Id int

	//1到多个玩家
	Players []*Player

	//随机排序的牌
	ShuffledCards []*Card

	//当前游戏的状态
	Status GameStatus

	//公共牌
	Cards []*Card

	//specify who is the dealer button. the dealer button rotates clockwise after each hand
	//the one next to the dealer button is small blind. and the one next to small blind is big blind
	Dealer int

	playerAdded <-chan bool

	playerReady <-chan bool

	playerLeft <-chan bool

	lock *sync.Mutex
}

func (g *Game) acceptClient(name string, conn net.Conn) bool {
	if g.Status == GameToBeDestroy {
		// g.Status = g.GameNew
		// g.Players = append(g.Players, NewPlayer(name, conn))
		return false
	}

	if len(g.Players) >= MAX_GROUP_PLAYER {
		return false
	}
	g.Players = append(g.Players, NewPlayer(name, conn))
	//g.playerAdded <- true
	return true
}

//新的一局游戏开始， 每人发两张牌
func (g *Game) onGameStart() {
	g.ShuffledCards = genRandCards()
	for i := 0; i < 2; i++ {
		for i := 0; i < len(g.Players); i++ {
			p := g.Players[(g.Dealer+i)%len(g.Players)]
			c := g.ShuffledCards[0]
			g.ShuffledCards = g.ShuffledCards[1:]
			err := p.SendCard(c)
			if err != nil {
				syslog.Error("Sendcard err:%v", err)
				continue
			}
		}
	}
}

func (g *Game) onGameOver() {
	g.Dealer = (g.Dealer + 1) / len(g.Players)
}

//发三张公共牌
func (g *Game) onFlop() {
	for _, p := range g.Players {
		p.SendCommCard(g.ShuffledCards[0])
	}
	g.ShuffledCards = g.ShuffledCards[1:]
}

//发第四张公共牌
func (g *Game) onTurn() {
	for _, p := range g.Players {
		p.SendCommCard(g.ShuffledCards[0])
	}
	g.ShuffledCards = g.ShuffledCards[1:]
}

//发第五张公共牌
func (g *Game) onRiver() {
	for _, p := range g.Players {
		p.SendCommCard(g.ShuffledCards[0])
	}
	g.ShuffledCards = g.ShuffledCards[1:]
}

func (g *Game) onPlayerAction(p *Player) {

}

func (g *Game) Run() {
GAME_RUN:
	for {
		//wait until there're enough players and all players are ready to start the game
		for {
			//all players left. destroy the game
			if len(g.Players) == 0 {
				g.Status = GameToBeDestroy
				break GAME_RUN
			}

			playersReady := false
			if len(g.playerAdded) >= 3 {
				playersReady = true
				for _, p := range g.Players {
					if p.Status != PlayerReady {
						playersReady = false
						break
					}
				}
			}
			if !playersReady {
				select {
				case <-g.playerAdded:
				case <-g.playerReady:
				case <-g.playerLeft:
					if len(g.Players) == 0 {
						g.Status = GameToBeDestroy
						break GAME_RUN
					}
				}
			}
		}

		//game start
		g.onGameStart()

		//wait flop
		g.onFlop()

		//wait turn
		g.onTurn()

		//wait river
		g.onRiver()

		//
		g.onGameOver()
	}
}

var gameManager struct {
	//to be used for clients lookup
	clients map[string]*Game

	games []*Game
}

func init() {
	gameManager.clients = make(map[string]*Game)
	//gameManager.games = make([]*Game)
}

func newGame() *Game {
	return &Game{
		//ShuffledCards: genRandCards(),
		Status: GameNew,
		lock:   &sync.Mutex{},
	}
}

func JoinGame(userName string, conn net.Conn) (*Game, error) {
	if userGame, ok := gameManager.clients[userName]; ok {
		if userGame.Status != GameToBeDestroy {
			return userGame, nil
		}
	}

	//todo:添加到玩家数为2的Game中， 这样可以凑成更多可以开局的Game
	for _, g := range gameManager.games {
		if g.acceptClient(userName, conn) {
			return g, nil
		}
	}

	newGame := newGame()
	newGame.Players = append(newGame.Players, NewPlayer(userName, conn))
	go newGame.Run()

	return newGame, nil
}
