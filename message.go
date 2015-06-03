package dzpk

const (
	RESULT_OK = 0
)

type MessageType int

const (
	//server actions. server -> player
	HoleCards MessageType = 20     //2 cards face down to each player
	Flop      MessageType = 20 + 1 //show 3 common cards. all players choose bet or fold after hole cards
	Turn      MessageType = 20 + 2 //show 1 common card after flop
	River     MessageType = 20 + 3 //show 1 common card after turn

	//player actions. player -> server
	Check MessageType = 30     //player check
	Bet   MessageType = 30 + 1 //player bet
	Raise MessageType = 30 + 2 //player raise
	Fold  MessageType = 30 + 3 //player fold
)

type RequestMessage struct {
	MsgId   int         `json:"msgid"`
	MsgType MessageType `json:"msgtype"`
	Payload interface{} `json:"payload"`
}

type ResponseMessage struct {
	RequestId int         `json:"requestid"`
	MsgType   MessageType `json:"msgtype"`
	Result    int         `json:"result"`
	Payload   interface{} `json:"payload"`
}
