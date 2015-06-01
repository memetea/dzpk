package dzpk

import "net"

//玩家状态
type PlayerStatus int

const (
	//玩家已连接， 尚未开始
	PlayerConnected PlayerStatus = 0
	//玩家准备好， 可以开始玩游戏了
	PlayerReady PlayerStatus = 1
	//玩家游戏中
	PlayerInPlay PlayerStatus = 2
	//玩家选择放弃已投筹码， 不继续玩
	PlayerGaveup PlayerStatus = 3
)

type Player struct {
	//玩家姓名
	Name string

	//玩家筹码
	Counter int

	//玩家是否准备好开始
	Status PlayerStatus

	//玩家的牌
	Cards []*Card

	//客户端连接
	Conn net.Conn
}

//发牌给玩家
func (p *Player) SendCard(c *Card) error {
	return nil
}

//发公共牌给玩家
func (p *Player) SendCommCard(c *Card) error {
	return nil
}

func NewPlayer(name string, conn net.Conn) *Player {
	return &Player{
		Name:   name,
		Status: PlayerConnected,
		Conn:   conn,
	}
}
