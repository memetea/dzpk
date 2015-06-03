package dzpk

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type PlayerConn struct {
	conn net.Conn
}

func (pc *PlayerConn) DoRequest(req *RequestMessage) (*ResponseMessage, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	//write bytes count
	err = binary.Write(pc.conn, binary.LittleEndian, uint32(len(b)))
	if err != nil {
		return nil, err
	}

	//write bytes
	n, err := pc.conn.Write(b)
	if err != nil {
		return nil, err
	}
	if n != len(b) {
		return nil, fmt.Errorf("Short write. %d != %d", n, len(b))
	}

	//read response
	var respLen uint32
	err = binary.Read(pc.conn, binary.LittleEndian, &respLen)
	if err != nil {
		return nil, err
	}
	respBuf := make([]byte, respLen)
	n, err = pc.conn.Read(respBuf)
	if err != nil {
		return nil, err
	}
	var resp ResponseMessage
	err = json.Unmarshal(respBuf, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

//玩家状态
type PlayerStatus int

const (
	PlayerConnected PlayerStatus = 0 //玩家已连接， 尚未开始
	PlayerReady     PlayerStatus = 1 //玩家准备好， 可以开始玩游戏了
	PlayerInPlay    PlayerStatus = 2 //玩家游戏中
	PlayerGaveup    PlayerStatus = 3 //玩家选择放弃已投筹码， 不继续玩
)

type Player struct {
	UserId  int    //player id. unique
	Name    string //player name
	msgId   int
	Counter int          //palyer counter
	Status  PlayerStatus //玩家是否准备好开始
	Cards   []*Card      //玩家的牌
	Conn    *PlayerConn  //客户端连接
}

//发牌给玩家
func (p *Player) SendHoleCard(c []*Card) error {
	if len(c) != 2 {
		return fmt.Errorf("HoleCard must be 2")
	}
	req := &RequestMessage{
		MsgId:   p.msgId,
		MsgType: HoleCards,
		Payload: c,
	}
	resp, err := p.Conn.DoRequest(req)
	if err != nil {
		return err
	}
	if resp.Result != RESULT_OK {
		return fmt.Errorf("User response result err:%v", resp)
	}
	p.msgId++
	return nil
}

//发公共牌给玩家
func (p *Player) SendCommCard(c []*Card) error {
	req := &RequestMessage{
		MsgId:   p.msgId,
		MsgType: Flop,
		Payload: c,
	}
	resp, err := p.Conn.DoRequest(req)
	if err != nil {
		return err
	}
	if resp.Result != RESULT_OK {
		return fmt.Errorf("User response result err:%v", resp)
	}
	p.msgId++
	return nil
}

func NewPlayer(name string, conn net.Conn) *Player {
	return &Player{
		Name:   name,
		Status: PlayerConnected,
		Conn:   &PlayerConn{conn},
	}
}
